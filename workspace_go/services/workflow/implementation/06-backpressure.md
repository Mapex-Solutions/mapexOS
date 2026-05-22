# Backpressure — MongoDB Write Latency Tracking

Decisão: A6

---

## Problema

A lib MongoDB atual (`packages/infrastructure/mongodb/`) não tem:

- Write latency tracking (só mede ping/read)
- BulkWrite no Model[T] (só InsertMany)
- Backpressure — se MongoDB está lento, o caller fica bloqueado sem informação

O Archiver do workflow faz BulkWrite de 5000 documentos. Se MongoDB estiver degradado (alta latência), o Archiver deveria reduzir o batch para aliviar a pressão.

---

## Como Estamos Resolvendo

### Implementação atual

**Princípio:** A lib expõe informação (modo de backpressure), o caller decide o que fazer. Zero auto-throttle.

```
╔═══════════╦═══════════════════════════╦══════════════════════════════════════════╗
║ Modo      ║ Trigger                   ║ O que o caller faz                       ║
╠═══════════╬═══════════════════════════╬══════════════════════════════════════════╣
║ Normal    ║ P99 < 150ms               ║ Nada — tudo normal                       ║
║ Throttled ║ P99 > 150ms por 3 janelas ║ Reduz batch (5000 → 1000)               ║
║ Backoff   ║ P99 > 500ms por 3 janelas ║ Reduz mais (1000 → 500) + pausa 2s      ║
╚═══════════╩═══════════════════════════╩══════════════════════════════════════════╝
```

---

## Como Implementar

### 1. Config extension (opt-in flag)

```go
// manager/types.go
type Config struct {
    // ... campos existentes ...

    // Backpressure (opt-in, default: false = backward compatible)
    EnableBackpressure   bool          // Habilita tracking de write latency
    BackpressureWindow   int           // Amostras no circular buffer (default: 1000)
    ThrottledThresholdMs int64         // P99 > isso → Throttled (default: 150ms)
    BackoffThresholdMs   int64         // P99 > isso → Backoff (default: 500ms)
}
```

### 2. API pública

```go
// manager/methods.go — novos métodos
func (m *MongoManager) GetBackpressureMode() BackpressureMode   // Normal|Throttled|Backoff
func (m *MongoManager) WriteP99() int64                          // P99 em ms (última janela)
func (m *MongoManager) RecordWriteLatency(d time.Duration)       // Caller registra após operação
```

### 3. Internals (circular buffer + goroutine)

```go
// manager/backpressure.go (~80 linhas)

type backpressureTracker struct {
    samples    []int64          // circular buffer (1000 amostras)
    pos        atomic.Int64     // posição no buffer
    p99        atomic.Int64     // último P99 calculado (ms)
    mode       atomic.Int32     // 0=Normal, 1=Throttled, 2=Backoff
    thresholds [2]int64         // [throttled, backoff] em ms
    aboveCount int              // janelas consecutivas acima do threshold
}

// Background goroutine (a cada 5s):
//   1. Copia samples (lock-free, atomic reads)
//   2. Ordena, calcula P99
//   3. Compara com thresholds
//   4. Se P99 > threshold por 3 janelas consecutivas → muda modo
//   5. Se P99 < threshold → reset para Normal
//   6. Armazena modo (atomic store)
```

### 4. Custo de performance

```
RecordWriteLatency():  atomic store num slot do buffer     → ~50ns
GetBackpressureMode(): atomic load do modo                 → ~1ns
Background goroutine:  sort 1000 amostras a cada 5s        → ~50μs (fora do hot path)
Total overhead por write: ~51ns (negligível)
```

### 5. Uso no Archiver

```go
func (a *ArchiverService) ProcessStateBatch(messages []*natsModel.Message) {
    mode := a.mongoManager.GetBackpressureMode()

    // Ajusta batch size baseado no modo
    batchSize := 5000
    switch mode {
    case mongoManager.Throttled:
        batchSize = 1000
        logger.Warn("[WORKFLOW:ARCHIVER] MongoDB throttled, reducing batch to 1000")
    case mongoManager.Backoff:
        batchSize = 500
        time.Sleep(2 * time.Second)
        logger.Warn("[WORKFLOW:ARCHIVER] MongoDB backoff, reducing batch to 500, pausing 2s")
    }

    // Processa em sub-batches se mensagens > batchSize
    for i := 0; i < len(messages); i += batchSize {
        end := min(i+batchSize, len(messages))
        chunk := messages[i:end]

        start := time.Now()
        a.bulkWrite(chunk)
        a.mongoManager.RecordWriteLatency(time.Since(start))
    }
}
```

### 6. Backward compatible

```go
// Outros serviços (events, triggers, etc.) — ZERO mudança
mgr, _ := mongoManager.New(mongoManager.Config{
    URI:      "mongodb://...",
    Database: "dev-events",
    // EnableBackpressure NÃO setado → false → tracking desligado
})
// GetBackpressureMode() sempre retorna Normal
// RecordWriteLatency() é no-op
// WriteP99() retorna 0
```

---

## Futuro Enterprise

```
╔═══════════════════════════════════╦══════════════╦══════════════════════════════════════════════╗
║ Feature                           ║ Prioridade   ║ Descrição                                    ║
╠═══════════════════════════════════╬══════════════╬══════════════════════════════════════════════╣
║ 1. Adaptive Batch Sizing          ║ P1           ║ Lib ajusta batch automaticamente baseado     ║
║                                   ║              ║ no P99. Caller configura min/max batch size.  ║
║ 2. Per-Collection Metrics         ║ P1           ║ P99/P95/P50 por collection, não global.      ║
║ 3. Prometheus Histograms          ║ P1           ║ mongo_write_duration_seconds{collection,op}  ║
║ 4. Write Concern Profiles         ║ P2           ║ Normal: w=majority, Backoff: w=1              ║
║ 5. Circuit Breaker                ║ P2           ║ P99 > 2s por 10 janelas → rejeita non-crit   ║
║ 6. Auto-Throttle (opt-in)         ║ P2           ║ Rate limiter interno na lib                   ║
║ 7. Connection Pool Metrics        ║ P3           ║ InUse, Available, WaitCount, WaitDuration     ║
║ 8. Slow Query Log                 ║ P3           ║ Operações > 200ms logadas automaticamente     ║
║ 9. Read/Write Splitting           ║ P3           ║ Reads → Secondary Preferred, Writes → Primary ║
╚═══════════════════════════════════╩══════════════╩══════════════════════════════════════════════╝
```

**Evolução:** Atual → Prometheus + per-collection (baseado em métricas reais) → Circuit breaker + auto-throttle (enterprise).

### Checklist de implementação

```
Na lib MongoDB (packages/infrastructure/mongodb/manager/):
  □ backpressure.go          — circular buffer, P99 calc, mode transitions (~80 linhas)
  □ types.go                 — 4 novos campos no Config + BackpressureMode type
  □ methods.go               — 3 novos métodos: GetBackpressureMode, WriteP99, RecordWriteLatency
  □ manager.go               — inicializar tracker no New() se EnableBackpressure=true

No workflow service (services/workflow/):
  □ bootstrap/mongo.go       — EnableBackpressure: true no Config
  □ modules/archiver/        — ler modo antes de cada batch, ajustar batch size
```

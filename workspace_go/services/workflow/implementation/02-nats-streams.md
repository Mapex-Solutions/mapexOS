# NATS Streams — 2 Streams + Double-buffer

Decisões: A3, A5

---

## Problema

O workflow engine precisa de filas NATS para:

1. **Criar novas instâncias** (triggers externos)
2. **Continuar instâncias existentes** (callbacks, signals, timers)

Se tudo estiver na mesma fila, um burst de 100K triggers novos atrasa callbacks de workflows que estão parados esperando resposta. Além disso, fanout nodes que criam múltiplas branches não devem sobrecarregar a fila.

---

## Como Estamos Resolvendo

### 2 streams separados por tipo de entrada (A3)

- **WORKFLOW-TRIGGER**: instâncias NOVAS (criar e começar a rodar)
- **WORKFLOW-RESUME**: instâncias que JÁ ESTÃO rodando (callbacks, signals, timers)

Streams separados = consumers independentes = zero contenção.

**Fanout NÃO vai pro NATS:** Branches executam via goroutines inline no mesmo worker. Só vai pro WORKFLOW-RESUME se uma branch suspender num async (callback/signal/timer resume).

### Double-buffer (A5)

Já implementado no NATS consumer (`asyncFetch` em `packages/infrastructure/nats/consumer.go`). O workflow service herda automaticamente ao usar `BatchMessageHandlerV2`.

- Prefetch batch N+1 enquanto processa batch N
- `MaxAckPending = BatchSize × 2` (1 processando + 1 prefetch)
- Zero idle time entre batches

---

## Como Implementar

### Streams definidos

```
Stream: WORKFLOW-TRIGGER
  Subjects: workflow.trigger.*
  Propósito: criar instância nova + iniciar execução
  Volume: alto (bursts de triggers, campanhas, IoT)
  BatchSize: 5000
  FetchTimeout: 500ms

Stream: WORKFLOW-RESUME
  Subjects: workflow.resume.callback.*    → resposta de JS-Executor, trigger, subworkflow
            workflow.resume.signal.*      → sinal externo (usuário, webhook)
            workflow.resume.timer.*       → delay expirou, poll timer
            workflow.resume.transition.*  → re-entrada após checkpoint
  Propósito: continuar execução de instância existente
  Volume: mais baixo (só pontos de re-entrada async)
  BatchSize: 1000
  FetchTimeout: 200ms (menor latência)
```

### Quem publica em qual stream

```
                                    WORKFLOW-TRIGGER    WORKFLOW-RESUME
                                    ────────────────    ───────────────
Trigger externo (API, NATS, MQTT)        ✓
Callback (JS-Executor, trigger)                              ✓
Signal (usuário, webhook)                                    ✓
Timer expirado (delay, poll)                                 ✓
Subworkflow retorno                                          ✓
Fanout (branches)                   goroutines inline (NÃO vai pro NATS)
```

### Fanout inline via goroutines

```
Fanout node tem 3 branches:

  Worker spawna 3 goroutines (estado já em memória):

  goroutine 1: set_state → condition → end
    → Inline, completa em ~1ms ✓

  goroutine 2: set_state → code(ASYNC)
    → Inline até code → suspende → checkpoint

  goroutine 3: set_state → trigger_event(ASYNC)
    → Inline até trigger_event → suspende → checkpoint

  Worker espera todas terminarem/suspenderem:
    → NATS KV Put (state atualizado)
    → Checkpoint pro Archiver
    → ACK da transition original

  Só vai pro WORKFLOW-RESUME quando a branch ASYNC receber callback/signal.
```

### Consumer configs

```go
// Consumer 1: Novos workflows
bus.StartConsumer(natsModel.ConsumerOptions{
    Stream:       "WORKFLOW-TRIGGER",
    Subject:      "workflow.trigger.>",
    Durable:      fmt.Sprintf("%s-trigger", serviceName),
    QueueGroup:   fmt.Sprintf("%s-TRIGGER-GROUP", serviceName),
    BatchSize:    5000,
    FetchTimeout: 500 * time.Millisecond,
    RetryPolicy: &natsModel.RetryPolicy{
        MaxRetries: 5,
        Backoff:    []time.Duration{1*time.Second, 5*time.Second, 30*time.Second, 2*time.Minute, 10*time.Minute},
        AckWait:    30 * time.Second,
    },
    DLQPolicy: &natsModel.DLQPolicy{
        Stream:      "MAPEXOS-DLQ",
        Subject:     "dlq.mapexos",
        ServiceName: serviceName,
        ServiceType: "workflow",
        EventType:   "trigger",
    },
    BatchMessageHandlerV2: func(messages []*natsModel.Message) {
        runtimeService.ProcessTriggerBatch(messages)
    },
})

// Consumer 2: Workflows existentes que precisam continuar
bus.StartConsumer(natsModel.ConsumerOptions{
    Stream:       "WORKFLOW-RESUME",
    Subject:      "workflow.resume.>",
    Durable:      fmt.Sprintf("%s-resume", serviceName),
    QueueGroup:   fmt.Sprintf("%s-RESUME-GROUP", serviceName),
    BatchSize:    1000,
    FetchTimeout: 200 * time.Millisecond,
    RetryPolicy: &natsModel.RetryPolicy{
        MaxRetries: 5,
        Backoff:    []time.Duration{1*time.Second, 5*time.Second, 30*time.Second, 2*time.Minute, 10*time.Minute},
        AckWait:    30 * time.Second,
    },
    DLQPolicy: &natsModel.DLQPolicy{
        Stream:      "MAPEXOS-DLQ",
        Subject:     "dlq.mapexos",
        ServiceName: serviceName,
        ServiceType: "workflow",
        EventType:   "resume",
    },
    BatchMessageHandlerV2: func(messages []*natsModel.Message) {
        runtimeService.ProcessResumeBatch(messages)
    },
})
```

### Mapa completo de streams do workflow service

```
╔════════════════════════╦══════════════════════════════╦═══════════════════╗
║ Stream                 ║ Quem produz                  ║ Quem consome      ║
╠════════════════════════╬══════════════════════════════╬═══════════════════╣
║ WORKFLOW-TRIGGER       ║ API / Trigger Service / NATS ║ Runtime (trigger) ║
║ WORKFLOW-RESUME        ║ Runtime (callback, signal)   ║ Runtime (resume)  ║
║ WORKFLOW-STATE         ║ Runtime (lifecycle events)   ║ Archiver          ║
║ WORKFLOW-LOGS          ║ Runtime (step logs)          ║ Logs-writer       ║
║ WORKFLOW-RECONCILER    ║ Reconciler + Archiver        ║ Reconciler        ║
╚════════════════════════╩══════════════════════════════╩═══════════════════╝
```

### NATS init (docker-compose)

Adicionar ao `nats-init`:

```bash
nats stream add WORKFLOW-TRIGGER \
  --subjects="workflow.trigger.*" \
  --storage=file --retention=work --defaults

nats stream add WORKFLOW-RESUME \
  --subjects="workflow.resume.callback.*,workflow.resume.signal.*,workflow.resume.timer.*,workflow.resume.transition.*" \
  --storage=file --retention=work --defaults

nats stream add WORKFLOW-STATE \
  --subjects="workflow.state.*" \
  --storage=file --retention=work --defaults

nats stream add WORKFLOW-LOGS \
  --subjects="workflow.logs.*" \
  --storage=file --retention=work --defaults

nats stream add WORKFLOW-RECONCILER \
  --subjects="workflow.reconciler.*" \
  --storage=file --retention=work --defaults
```

### Checklist de implementação

```
No nats-init (docker-compose):
  □ Adicionar 5 streams ao script de inicialização

No workflow service:
  □ Consumer WORKFLOW-TRIGGER (BatchSize: 5000, FetchTimeout: 500ms)
  □ Consumer WORKFLOW-RESUME (BatchSize: 1000, FetchTimeout: 200ms)
  □ Consumer WORKFLOW-STATE → Archiver
  □ Consumer WORKFLOW-LOGS → Logs-writer
  □ Consumer WORKFLOW-RECONCILER → Reconciler
```

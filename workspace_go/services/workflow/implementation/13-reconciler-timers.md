# Reconciler Timers — MongoDB + NATS Self-Schedule

Decisao: A11

---

## Problema

O Reconciler e responsavel por acordar workflows suspensos em timers (delay, timeout de wait_signal, timeout de code/subworkflow). O design V1 usava:

- Min-heap in-memory (container/heap) no pod
- Goroutine ticker a cada 1s
- Scan de todas as keys do NATS KV no startup (recoverFromKV)

**Problemas do V1:**

| Problema | Impacto |
|----------|---------|
| Pod morre → heap perdido | Timers somem, workflows ficam presos |
| KV scan no startup | O(N) sobre TODAS as keys, nao escala a 1M+ |
| Estado local no pod | Scale up/down perde timers |
| Ticker a cada 1s | Polling constante mesmo sem timers |
| QueueGroup no timer | Timer vai pra 1 pod, mas outro pod pode ter o heap |

---

## Como Estamos Resolvendo

### MongoDB como alarm clock + NATS self-schedule

**Principio:** MongoDB sabe QUANDO acordar (index no `timerExpiresAt`). KV sabe O QUE fazer (NodeStates completo). NATS agenda a entrega precisa.

### Novo campo no WorkflowInstance

```go
type WorkflowInstance struct {
    // ... campos existentes ...
    TimerExpiresAt *time.Time `bson:"timerExpiresAt,omitempty"`
    // ... StartedAt, CompletedAt, Created, Updated ...
}
```

Unico campo novo. `omitempty` = nao aparece em 95% dos documentos (workflows sem timer). Limpo automaticamente no terminal Upsert.

### MongoDB lifecycle — 4 writes (via Archiver, BULK)

```
╔═══════════════╦══════════════════════════════════════════════════════════╦════════╗
║ Evento        ║ MongoDB Operation (via Archiver, BULK)                  ║ ~Size  ║
╠═══════════════╬══════════════════════════════════════════════════════════╬════════╣
║               ║                                                          ║        ║
║ 1. START      ║ InsertOne lightweight stub                               ║ ~200B  ║
║               ║ { status:"running", activeNodeIds, workflowName, ... }   ║        ║
║               ║                                                          ║        ║
║ 2. AWAIT      ║ UpdateOne (SOMENTE async nodes)                          ║ ~150B  ║
║    START      ║ { status:"waiting", timerExpiresAt: earliest,            ║        ║
║               ║   activeNodeIds: ["n-delay-1"] }                         ║        ║
║               ║                                                          ║        ║
║ 3. AWAIT      ║ UpdateOne (SOMENTE async nodes)                          ║ ~100B  ║
║    FINISH     ║ { status:"running", timerExpiresAt: null }               ║        ║
║               ║                                                          ║        ║
║ 4. COMPLETE   ║ Upsert FULL (terminal: completed/failed/cancelled)       ║ ~5-25K ║
║               ║ Sobrescreve TUDO. KV Delete.                             ║        ║
║               ║                                                          ║        ║
╚═══════════════╩══════════════════════════════════════════════════════════╩════════╝
```

**Workflow inline (sem async):** so 1 + 4 (Start + Complete). Zero writes intermediarios.

**Workflow com async:** 1 + 2 + 3 + 4. Se 2 async nodes: 1 + 2 + 3 + 2 + 3 + 4.

### NATS subjects no WORKFLOW-STATE

```
workflow.state.created    → Step 1 (START)         — ja existente
workflow.state.waiting    → Step 2 (AWAIT START)    — NOVO
workflow.state.resumed    → Step 3 (AWAIT FINISH)   — NOVO
workflow.state.completed  → Step 4 (COMPLETE)       — ja existente
workflow.state.failed     → Step 4 (COMPLETE)       — ja existente
workflow.state.cancelled  → Step 4 (COMPLETE)       — ja existente
```

---

## Short Timer Fast-Path (< 1 minuto)

Timers curtos (delay 5s, 10s, 30s) nao devem esperar o sweep de 1 minuto. O Archiver decide na hora:

```
Archiver processa AWAIT START event:

  1. MongoDB UpdateOne (sempre — para safety net + listagem)

  2. Para CADA node no array TimerNodes do evento:
       SE expiresAt <= now + 1min:
         → NATS scheduled message direto ao WORKFLOW-RESUME
         → subject: workflow.resume.timer.{instanceId}
         → delivery: expiresAt exato desse node
         → payload: { instanceId, nodeId, reason: "timer_expired" }
         → NATS entrega no tempo certo → Runtime resume ESSE node

  3. Nodes com expiresAt > 1min:
       → Nada. Reconciler sweep pega no futuro.

  Exemplo fanout com 4 delays (5s, 10s, 30s, 2h):
    → MongoDB: timerExpiresAt = T+5s (earliest)
    → NATS schedula 3 msgs (5s, 10s, 30s — todos < 1min) com nodeId individual
    → delay(2h) fica no MongoDB → sweep pega em ~1h59min
```

**Resultado:**
- `delay(10s)` → Archiver → NATS schedule 10s → resume. Zero latencia extra.
- `delay(1h)` → Archiver → MongoDB → sweep pega em ~59min → NATS schedule ~1min → resume.
- NATS nunca fica lotado com timers distantes (so proximos 60s).

---

## NATS Self-Schedule — Sweep a cada 1 minuto

### Bootstrap: OnConnect callback + NATS dedup

O Reconciler precisa disparar o primeiro sweep quando o pod conecta (ou reconecta) ao NATS. Isso requer um **callback opcional no `packages/infrastructure/nats`**.

#### O que precisa mudar no package NATS

```go
// packages/infrastructure/nats/types.go — NOVO campo no Config
type Config struct {
    Options    Options
    OnConnect  func(c *Client) // NOVO — callback opcional, chamado apos connect E reconnect
}
```

```go
// packages/infrastructure/nats/nats.go — chamar callback apos connect
func New(c Config) (*Client, error) {
    nc, err := c.Options.Connect()
    if err != nil {
        return nil, err
    }

    js, err := nc.JetStream()
    if err != nil {
        return nil, err
    }

    client := &Client{nc: nc, js: js}

    logger.Info("[INFRA:NATS] Connected to server")

    // Chamar OnConnect callback se configurado
    if c.OnConnect != nil {
        c.OnConnect(client)
    }

    // Registrar reconnect handler para chamar OnConnect novamente
    // nats.Options.ReconnectedCB ja existe na lib nats.go
    if c.OnConnect != nil {
        nc.SetReconnectHandler(func(conn *nats.Conn) {
            logger.Info("[INFRA:NATS] Reconnected — calling OnConnect callback")
            c.OnConnect(client)
        })
    }

    return client, nil
}
```

**Por que no package e nao no service:**
- O `bootstrap/nats.go` de QUALQUER service pode usar (nao so workflow)
- O campo e `opcional` — services que nao precisam nao sao afetados
- Segue o padrao da lib nats.go (`ConnectedCB`, `ReconnectedCB`) mas com acesso ao nosso `*Client`
- Mantem encapsulamento: o service nao precisa saber detalhes de `nats.Conn`

#### Como o workflow service usa

```go
// bootstrap/nats.go — workflow service
natsCoreCfg := config.GetNatsConfig()
natsCoreCfg.OnConnect = func(client *natsModel.Client) {
    // Publicar primeiro sweep do Reconciler
    // MsgId com timestamp garante dedup de N pods conectando ao mesmo tempo
    bus := natsModel.NewBus(client)
    bus.Publish(natsModel.PublishConfig{
        Subject: "workflow.reconciler.sweep",
        Data:    map[string]string{"trigger": "connect"},
    })
    logger.Info("[APP:BOOTSTRAP] Reconciler sweep triggered via OnConnect")
}

nc, err := natsModel.New(natsCoreCfg)
```

#### Cenario: 3 pods ligam simultaneamente

```
Tempo  Evento
─────  ─────────────────────────────────────────────────────
T=0    Pod-A conecta ao NATS
       → OnConnect → Publish "workflow.reconciler.sweep"
       → Mensagem entra no JetStream stream WORKFLOW-RECONCILER

T=0    Pod-B conecta ao NATS (simultaneo)
       → OnConnect → Publish "workflow.reconciler.sweep"
       → NATS verifica DuplicateWindow (15min padrao)
       → Aceita (MsgId diferente OU sem MsgId → 2 mensagens no stream)

T=0    Pod-C conecta ao NATS (simultaneo)
       → OnConnect → Publish "workflow.reconciler.sweep"
       → 3 mensagens no stream (pior caso)

T=0.01 Consumer WORKFLOW-RECONCILER (QueueGroup):
       → Pod-A recebe msg 1 → executa sweep → MongoDB query → agenda timers → ACK
       → Pod-B recebe msg 2 → executa sweep → MongoDB query retorna mesmos docs
         → KV Get para cada → alguns ja foram agendados por Pod-A
         → Agenda os que faltam (idempotente) → ACK
       → Pod-C recebe msg 3 → executa sweep → tudo ja agendado → ACK (noop)

T=60   Pod-A (ou B ou C) termina sweep:
       → Self-schedule: Publish "workflow.reconciler.sweep" com delay 1min
       → Agora so 1 mensagem pendente no stream (loop estavel)
```

**Resultado:** Apos o burst inicial de connect, o loop estabiliza em 1 sweep por minuto. Sweeps extras sao inofensivos (KV idempotencia + MongoDB query barata ~1ms).

#### Cenario: pod reconecta apos queda de rede

```
T=0     Pod-A perde conexao com NATS
        → nats.go detecta disconnect
        → Consumers param de receber mensagens

T=30s   Rede volta
        → nats.go reconecta automaticamente (reconnect handler da lib)
        → nc.SetReconnectHandler dispara
        → OnConnect callback executa
        → Publish "workflow.reconciler.sweep"
        → Sweep roda, pega timers que podem ter acumulado durante os 30s de queda
```

### Como funciona o sweep

```
1. Pod conecta/reconecta ao NATS:
     → OnConnect callback → Publish "workflow.reconciler.sweep"
     → Se N pods conectam ao mesmo tempo: N mensagens, todas processadas, idempotente

2. UM pod recebe a mensagem (QueueGroup — round-robin entre pods):
     → Query MongoDB:
         { timerExpiresAt: { $lte: ISODate("now + 1min") } }
         .projection({ _id:1, activeNodeIds:1, timerExpiresAt:1 })
         .limit(1000)

3. Para cada instancia encontrada:
     → KV Get → verifica que ainda esta waiting (idempotencia)
     → SE KV diz completed/cancelled/running → SKIP (stale)
     → Itera NodeStates procurando nodes com waitType + expiresAt:
       SE expiresAt <= now:
         → JA EXPIROU
         → Publish resume no WORKFLOW-RESUME:
             { instanceId, nodeId, reason: "timer_expired" }
       SE now < expiresAt <= now + 1min:
         → EXPIRA EM BREVE
         → NATS scheduled message para expiresAt exato
             { instanceId, nodeId, reason: "timer_expired" }

4. Self-schedule proximo sweep:
     → Publica "workflow.reconciler.sweep" com delay de 1min
     → NATS entrega em 1min → volta ao passo 2

5. Pod morre:
     → Mensagem de sweep fica no JetStream (file-backed)
     → QueueGroup reentrega para outro pod
     → Zero perda de timers

6. TODOS os pods morrem:
     → Mensagem de sweep fica no JetStream
     → Primeiro pod que subir → OnConnect → publish sweep → sweep roda
     → SE nenhum sweep pendente no stream (expired): OnConnect garante que um novo e publicado
```

### Frequencia: por que 1 minuto

```
Muito rapido (1s):   86,400 queries/dia, desperdicio em 95% dos casos (sem timer novo)
Muito lento (10min): Timer de 1min espera ate 10min no pior caso
1 minuto:            1,440 queries/dia, latencia maxima +1min, query ~1ms

Para timers < 1min: Archiver agenda NATS scheduled message DIRETO (zero espera por sweep)
Para timers >= 1min: Sweep pega no proximo ciclo. Latencia maxima = 1min.
Resultado: 99% dos timers curtos tem latencia zero. Timers longos tem +1min aceitavel.
```

### Por que funciona sem estado no pod

```
Estado no pod:     ZERO (tudo em NATS JetStream + MongoDB)
Pod descartavel:   Mata e levanta — OnConnect publica sweep, outro pod assume
Scale up:          Novo pod → OnConnect → sweep extra (idempotente) → QueueGroup distribui
Scale down:        Pod removido → NATS reentrega sweep pendente para restantes
Todos os pods caem: Sweep fica no JetStream → primeiro pod que subir retoma via OnConnect
Reconnect:         nats.go reconecta auto → SetReconnectHandler → OnConnect → sweep
```

### Idempotencia do resume

O resume handler (Runtime) sempre faz KV verify antes de processar:

```
Reconciler publica resume
  → Runtime recebe
  → KV Get(instanceId)
  → SE instance.NodeStates[nodeId]["waitType"] ainda existe:
      → Processa resume (continua executando)
  → SE ja foi processado (NodeStates limpo):
      → Skip (idempotente, sem efeito)
```

Duplicatas de sweep ou scheduled messages sao inofensivas.

---

## Mensagem de Resume (NATS scheduled ou Reconciler)

Toda mensagem de resume **DEVE incluir nodeId**. Sem ele, o Runtime nao sabe qual branch/node resumir.

```go
// Publicado pelo Archiver (short timer) ou Reconciler (sweep)
// Subject: workflow.resume.timer.{instanceId}
type TimerResumeMessage struct {
    InstanceID string `json:"instanceId"`
    NodeID     string `json:"nodeId"`      // QUAL node está sendo resumido
    Reason     string `json:"reason"`      // "timer_expired"
}
```

---

## Fanout com multiplos timers (race mode)

Fanout e race mode — quem terminar primeiro leva. Todos os timers sao gravados.

### StateEvent para waiting com multiplos nodes

```go
type StateEvent struct {
    InstanceID     string       `json:"instanceId"`
    Status         string       `json:"status"`
    ActiveNodeIDs  []string     `json:"activeNodeIds,omitempty"`
    TimerExpiresAt *time.Time   `json:"timerExpiresAt,omitempty"` // earliest de todos
    TimerNodes     []TimerNode  `json:"timerNodes,omitempty"`     // cada node com timer
    // ... outros campos omitidos ...
}

type TimerNode struct {
    NodeID    string    `json:"nodeId"`
    ExpiresAt time.Time `json:"expiresAt"`
}
```

### Flow — fanout com 4 delays

```
Runtime executa fanout com 4 branches:
  Branch 1: delay(5s)
  Branch 2: delay(10s)
  Branch 3: delay(30s)
  Branch 4: delay(2h)

  Todas suspendem. Runtime monta:
    ActiveNodeIDs: ["n-d1", "n-d2", "n-d3", "n-d4"]
    NodeStates: {
      "n-d1": {"waitType":"timer", "expiresAt":"T+5s"},
      "n-d2": {"waitType":"timer", "expiresAt":"T+10s"},
      "n-d3": {"waitType":"timer", "expiresAt":"T+30s"},
      "n-d4": {"waitType":"timer", "expiresAt":"T+2h"}
    }

  KV Put (estado completo)
  NATS Pub "workflow.state.waiting":
    timerExpiresAt: T+5s (earliest)
    timerNodes: [
      {nodeId:"n-d1", expiresAt:"T+5s"},
      {nodeId:"n-d2", expiresAt:"T+10s"},
      {nodeId:"n-d3", expiresAt:"T+30s"},
      {nodeId:"n-d4", expiresAt:"T+2h"}
    ]

Archiver processa AWAIT START:
  MongoDB UpdateOne: { status:"waiting", timerExpiresAt: T+5s, activeNodeIds: [...] }
  n-d1: T+5s  < 1min → NATS schedule { instanceId, nodeId:"n-d1", expiresAt:T+5s }
  n-d2: T+10s < 1min → NATS schedule { instanceId, nodeId:"n-d2", expiresAt:T+10s }
  n-d3: T+30s < 1min → NATS schedule { instanceId, nodeId:"n-d3", expiresAt:T+30s }
  n-d4: T+2h  > 1min → nada (sweep pega depois)

T+5s: NATS entrega resume { nodeId:"n-d1" }
  Runtime: KV Get → confirma n-d1 waiting → resume branch 1
  Branch 1 executa → hit merge → merge precisa de mais branches? NAO (race mode)
  → Merge resolve com branch 1 → continua apos merge → end
  → Workflow COMPLETA
  NATS Pub: "workflow.state.completed"

T+5s: Archiver: terminal ReplaceOne → sobrescreve TUDO → timerExpiresAt desaparece

T+10s: NATS entrega resume { nodeId:"n-d2" }
  Runtime: KV Get("inst:abc-123") → KEY NOT FOUND (Archiver ja deletou)
  → SKIP (idempotente)

T+30s: NATS entrega resume { nodeId:"n-d3" } → KEY NOT FOUND → SKIP
T+2h:  Reconciler sweep encontra... nada (timerExpiresAt ja foi limpo pelo terminal)
```

### Fanout SEM merge (branches independentes)

Se nao tem merge node, cada branch roda independentemente. Quando branch 1 termina
e nao tem merge, o workflow continua esperando as outras branches:

```
T+5s:  Resume n-d1 → branch 1 completa (hit end da branch, nao do workflow)
       Workflow continua waiting (3 branches restam)
       Runtime publica workflow.state.waiting DE NOVO:
         timerExpiresAt: T+10s (proximo earliest)
         timerNodes: [{n-d2, T+10s}, {n-d3, T+30s}, {n-d4, T+2h}]
       Archiver: MongoDB UpdateOne { timerExpiresAt: T+10s }

T+10s: Resume n-d2 → branch 2 completa
       2 restam → workflow.state.waiting: timerExpiresAt: T+30s

T+30s: Resume n-d3 → branch 3 completa
       1 resta → workflow.state.waiting: timerExpiresAt: T+2h

T+2h:  Reconciler sweep → resume n-d4 → branch 4 completa
       Todas completaram → workflow COMPLETA
       workflow.state.resumed + workflow.state.completed
```

---

## MongoDB Index

```js
// Partial index — so inclui documentos COM timerExpiresAt
// Com 1M instances e ~50K com timer ativo: index ~5MB em RAM
db.workflow_instances.createIndex(
  { timerExpiresAt: 1 },
  {
    name: "idx_timer_expires",
    partialFilterExpression: { timerExpiresAt: { $type: "date" } }
  }
)
```

**Query do sweep:** `{ timerExpiresAt: { $lte: now + 1min } }` → IXSCAN no partial index → ~1ms.

---

## Reconciler Consumer

```go
// WORKFLOW-RECONCILER stream — self-schedule
bus.StartConsumer(natsModel.ConsumerOptions{
    Stream:       "WORKFLOW-RECONCILER",
    Subject:      "workflow.reconciler.>",
    Durable:      fmt.Sprintf("%s-reconciler", serviceName),
    QueueGroup:   fmt.Sprintf("%s-RECONCILER-GROUP", serviceName),
    BatchSize:    1,          // 1 sweep por vez
    FetchTimeout: 60 * time.Second,  // Espera ate 1min por proximo sweep

    RetryPolicy: &natsModel.RetryPolicy{
        MaxRetries: 3,
        Backoff:    []time.Duration{5 * time.Second, 15 * time.Second, 30 * time.Second},
        AckWait:    60 * time.Second,
    },

    DLQPolicy: &natsModel.DLQPolicy{
        Stream:      "MAPEXOS-DLQ",
        Subject:     "dlq.mapexos",
        ServiceName: serviceName,
        ServiceType: "workflow",
        EventType:   "reconciler.sweep",
    },

    BatchMessageHandlerV2: func(messages []*natsModel.Message) {
        reconcilerService.HandleSweep(messages[0])
    },
})

// Ao inicializar (bootstrap): publicar primeiro sweep
natsPublisher.Publish("workflow.reconciler.sweep", []byte(`{"trigger":"startup"}`))
```

---

## Archiver — Novos batch types

```go
func (s *ArchiverService) ProcessStateBatch(messages []*natsModel.Message) {
    var createdStubs     []LightweightInstance
    var waitingUpdates   []TimerUpdate
    var resumedUpdates   []string  // instanceIds
    var terminalFull     []*entities.WorkflowInstance

    for _, msg := range messages {
        switch {
        case isCreatedEvent(msg.Subject):   // existente
            // acumula stub
        case isWaitingEvent(msg.Subject):   // NOVO
            // acumula { instanceId, status, timerExpiresAt, activeNodeIds }
            // SE timerExpiresAt <= now + 1min → NATS scheduled message direto
        case isResumedEvent(msg.Subject):   // NOVO
            // acumula instanceId → set status:"running", timerExpiresAt: null
        case isTerminalEvent(msg.Subject):  // existente
            // KV Get → acumula full instance
        }
    }

    // Batch 1: InsertMany (created)        — existente
    // Batch 2: BulkUpdateOne (waiting)     — NOVO
    // Batch 3: BulkUpdateOne (resumed)     — NOVO
    // Batch 4: BulkUpsert (terminal)       — existente
}
```

### TimerUpdate struct

```go
type TimerUpdate struct {
    InstanceID     string
    Status         string     // "waiting"
    ActiveNodeIDs  []string
    TimerExpiresAt time.Time  // earliest de todos os timer nodes
}
```

### Short timer scheduling (dentro do case waiting)

```go
case isWaitingEvent(msg.Subject):
    update := buildTimerUpdate(event)
    waitingUpdates = append(waitingUpdates, update)
    refs = append(refs, archiverTypes.MsgRef{Msg: msg, Batch: "waiting"})

    // Short timer fast-path: agendar cada node com timer < 1min
    for _, tn := range event.TimerNodes {
        if tn.ExpiresAt.Before(time.Now().Add(1 * time.Minute)) {
            s.deps.Publisher.PublishScheduled(
                fmt.Sprintf("workflow.resume.timer.%s", event.InstanceID),
                TimerResumeMessage{
                    InstanceID: event.InstanceID,
                    NodeID:     tn.NodeID,
                    Reason:     "timer_expired",
                },
                tn.ExpiresAt, // NATS entrega neste horario exato
            )
        }
    }
```

---

## Principio: Lightweight vs Full Snapshot

Async nodes (delay, code, subworkflow, trigger_event, wait_signal) criam **snapshots lightweight** no MongoDB. Somente eventos terminais (completed, failed, cancelled) gravam o estado completo.

```
╔═══════════════════════╦════════════════════════════════════════════════════════╦════════╗
║ Tipo de Evento        ║ O que grava no MongoDB                                ║ ~Size  ║
╠═══════════════════════╬════════════════════════════════════════════════════════╬════════╣
║                       ║                                                        ║        ║
║ AWAIT START           ║ UpdateOne PARCIAL:                                     ║ ~150B  ║
║ (async suspendeu)     ║   - status: "waiting"                                  ║        ║
║                       ║   - timerExpiresAt: earliest timer                     ║        ║
║                       ║   - activeNodeIds: quais nodes estao waiting           ║        ║
║                       ║                                                        ║        ║
║                       ║ NAO grava:                                             ║        ║
║                       ║   × state (variaveis do workflow)                      ║        ║
║                       ║   × executionPath (historico de execucao)              ║        ║
║                       ║   × nodeOutputs (saidas dos nodes)                    ║        ║
║                       ║   × nodeStates (estado completo de cada node)         ║        ║
║                       ║   × logs, errors, etc.                                ║        ║
║                       ║                                                        ║        ║
║ AWAIT FINISH          ║ UpdateOne PARCIAL:                                     ║ ~100B  ║
║ (async retomou)       ║   - status: "running"                                  ║        ║
║                       ║   - timerExpiresAt: null                               ║        ║
║                       ║                                                        ║        ║
║                       ║ Mesma coisa: NAO grava state, path, outputs...         ║        ║
║                       ║                                                        ║        ║
║ TERMINAL              ║ Upsert COMPLETO:                                       ║~5-25KB ║
║ (completed/failed/    ║   - TUDO: state, executionPath, nodeOutputs,           ║        ║
║  cancelled)           ║     nodeStates, logs, errors, completedAt, etc.        ║        ║
║                       ║   - Sobrescreve documento inteiro (ReplaceOne)         ║        ║
║                       ║   - KV Delete (hot state removido)                     ║        ║
║                       ║                                                        ║        ║
╚═══════════════════════╩════════════════════════════════════════════════════════╩════════╝
```

**Por que funciona:**
- O estado completo SEMPRE esta no NATS KV durante a execucao (fonte da verdade)
- MongoDB no AWAIT START/FINISH serve SOMENTE para: (1) listagem no frontend com status correto, (2) timerExpiresAt para o Reconciler poder fazer sweep
- Se o pod morrer, o KV tem tudo. Se o KV morrer, o NATS JetStream tem as mensagens pendentes
- O custo de 150B por update vs 5-25KB por upsert = economia de ~99% nos writes intermediarios

**Consequencia:** Um workflow com 10 async nodes faz: 1 InsertOne (200B) + 10 UpdateOne (150B) + 10 UpdateOne (100B) + 1 Upsert (5-25KB) = ~2.7KB total no MongoDB durante toda a vida. Sem os lightweight snapshots seriam 21 upserts de 5-25KB = ~100-525KB.

---

## Archiver Module — Alteracoes Detalhadas

O Archiver e o UNICO modulo que escreve no MongoDB. Toda alteracao de persistencia passa por aqui.

### Visao geral das mudancas

```
ANTES (V1):
  Archiver processa 2 tipos de evento:
    created   → BulkInsertLightweight
    terminal  → KV Get FULL → BulkUpsertFull → KV Delete

DEPOIS (V2):
  Archiver processa 4 tipos de evento:
    created   → BulkInsertLightweight              (existente)
    waiting   → BulkUpdateWaiting + short timers    (NOVO)
    resumed   → BulkUpdateResumed                   (NOVO)
    terminal  → KV Get FULL → BulkUpsertFull → KV Delete  (existente)

  NOVA dependencia:
    Publisher (NATS) para agendar short timers via scheduled messages
```

### Arquivo por arquivo

**1. `interfaces/message/types.go`** — StateEvent struct

```go
// ANTES
type StateEvent struct {
    InstanceID    string `json:"instanceId"`
    WorkflowID    string `json:"workflowId"`
    OrgID         string `json:"orgId"`
    WorkflowName  string `json:"workflowName"`
    Status        string `json:"status"`
    ActiveNodeID  string `json:"activeNodeId,omitempty"`   // singular
    Version       int    `json:"version"`
}

// DEPOIS
type StateEvent struct {
    InstanceID     string      `json:"instanceId"`
    WorkflowID     string      `json:"workflowId,omitempty"`
    OrgID          string      `json:"orgId,omitempty"`
    WorkflowName   string      `json:"workflowName,omitempty"`
    Status         string      `json:"status"`
    ActiveNodeIDs  []string    `json:"activeNodeIds,omitempty"`   // PLURAL
    Version        int         `json:"version,omitempty"`
    TimerExpiresAt *time.Time  `json:"timerExpiresAt,omitempty"` // NOVO — earliest timer
    TimerNodes     []TimerNode `json:"timerNodes,omitempty"`     // NOVO — todos os nodes com timer
}

// NOVO
type TimerNode struct {
    NodeID    string    `json:"nodeId"`
    ExpiresAt time.Time `json:"expiresAt"`
}
```

**2. `domain/repositories/types.go`** — LightweightInstance + novos types

```go
// ANTES
type LightweightInstance struct {
    // ... campos atuais ...
    ActiveNodeID string `bson:"activeNodeId"`   // singular
}

// DEPOIS
type LightweightInstance struct {
    // ... campos atuais ...
    ActiveNodeIDs []string `bson:"activeNodeIds"` // PLURAL
}

// NOVO — parametro para BulkUpdateWaiting
type WaitingUpdate struct {
    InstanceID     string
    Status         string     // "waiting"
    ActiveNodeIDs  []string
    TimerExpiresAt time.Time  // earliest de todos os timer nodes
}

// NOVO — parametro para BulkUpdateResumed (somente instanceId e necessario)
// Usa []string direto no metodo
```

**3. `domain/repositories/archive_repository.go`** — 2 novos metodos

```go
// ANTES
type ArchiveRepository interface {
    BulkInsertLightweight(ctx context.Context, stubs []LightweightInstance) error
    BulkUpsertFull(ctx context.Context, instances []*entities.WorkflowInstance) error
}

// DEPOIS
type ArchiveRepository interface {
    BulkInsertLightweight(ctx context.Context, stubs []LightweightInstance) error
    BulkUpsertFull(ctx context.Context, instances []*entities.WorkflowInstance) error

    // NOVO — UpdateOne por instanceId: { status:"waiting", timerExpiresAt, activeNodeIds }
    BulkUpdateWaiting(ctx context.Context, updates []WaitingUpdate) error

    // NOVO — UpdateOne por instanceId: { status:"running", timerExpiresAt: null }
    BulkUpdateResumed(ctx context.Context, instanceIDs []string) error
}
```

**4. `infrastructure/persistence/mongo/collection/archive_repository.go`** — implementar 2 novos metodos

```go
// BulkUpdateWaiting — MongoDB BulkWrite com UpdateOne operations
func (r *ArchiveRepository) BulkUpdateWaiting(ctx context.Context, updates []WaitingUpdate) error {
    var models []mongo.WriteModel
    for _, u := range updates {
        filter := bson.M{"_id": u.InstanceID}
        update := bson.M{"$set": bson.M{
            "status":         u.Status,
            "activeNodeIds":  u.ActiveNodeIDs,
            "timerExpiresAt": u.TimerExpiresAt,
            "updated":        time.Now(),
        }}
        models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
    }
    _, err := r.collection.BulkWrite(ctx, models)
    return err
}

// BulkUpdateResumed — MongoDB BulkWrite com UpdateOne operations
func (r *ArchiveRepository) BulkUpdateResumed(ctx context.Context, instanceIDs []string) error {
    var models []mongo.WriteModel
    for _, id := range instanceIDs {
        filter := bson.M{"_id": id}
        update := bson.M{
            "$set":   bson.M{"status": "running", "updated": time.Now()},
            "$unset": bson.M{"timerExpiresAt": ""},
        }
        models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
    }
    _, err := r.collection.BulkWrite(ctx, models)
    return err
}
```

**5. `application/di/archiver_service_di.go`** — nova dependencia Publisher

```go
// ANTES
type ArchiverServiceDependenciesInjection struct {
    dig.In
    ArchiveRepo repositories.ArchiveRepository
    KVStore     natsModel.KeyValueStore
}

// DEPOIS
type ArchiverServiceDependenciesInjection struct {
    dig.In
    ArchiveRepo repositories.ArchiveRepository
    KVStore     natsModel.KeyValueStore
    Publisher   natsModel.Publisher              // NOVO — para scheduled messages (short timers)
}
```

**6. `application/services/archiver_service.go`** — 2 novos cases + 2 novos batches

Mudancas no `ProcessStateBatch`:

```go
// NOVOS acumuladores (alem dos existentes):
var waitingUpdates  []repositories.WaitingUpdate
var resumedIDs      []string

// NOVOS cases no switch:
case isWaitingEvent(msg.Subject):
    update := buildWaitingUpdate(event)
    waitingUpdates = append(waitingUpdates, update)
    refs = append(refs, archiverTypes.MsgRef{Msg: msg, Batch: "waiting"})

    // Short timer fast-path: agendar NATS scheduled msg para cada timer < 1min
    for _, tn := range event.TimerNodes {
        if tn.ExpiresAt.Before(time.Now().Add(1 * time.Minute)) {
            s.deps.Publisher.PublishScheduled(
                fmt.Sprintf("workflow.resume.timer.%s", event.InstanceID),
                TimerResumeMessage{
                    InstanceID: event.InstanceID,
                    NodeID:     tn.NodeID,
                    Reason:     "timer_expired",
                },
                tn.ExpiresAt,
            )
        }
    }

case isResumedEvent(msg.Subject):
    resumedIDs = append(resumedIDs, event.InstanceID)
    refs = append(refs, archiverTypes.MsgRef{Msg: msg, Batch: "resumed"})

// NOVOS batches (apos os existentes):
// Batch 2: Waiting updates (status + timer)
if len(waitingUpdates) > 0 {
    if err := s.deps.ArchiveRepo.BulkUpdateWaiting(ctx, waitingUpdates); err != nil {
        nackBatch(refs, "waiting")
    } else {
        ackBatch(refs, "waiting")
    }
}

// Batch 3: Resumed updates (status + clear timer)
if len(resumedIDs) > 0 {
    if err := s.deps.ArchiveRepo.BulkUpdateResumed(ctx, resumedIDs); err != nil {
        nackBatch(refs, "resumed")
    } else {
        ackBatch(refs, "resumed")
    }
}

// NOVAS helper functions:
func isWaitingEvent(subject string) bool {
    return subject == "workflow.state.waiting"
}

func isResumedEvent(subject string) bool {
    return subject == "workflow.state.resumed"
}

func buildWaitingUpdate(event archiverMsg.StateEvent) repositories.WaitingUpdate {
    return repositories.WaitingUpdate{
        InstanceID:     event.InstanceID,
        Status:         event.Status,
        ActiveNodeIDs:  event.ActiveNodeIDs,
        TimerExpiresAt: *event.TimerExpiresAt,
    }
}
```

**7. `interfaces/message/consumers/workflow_state/consumer.go`** — sem mudancas

O consumer ja consome `workflow.state.>` (wildcard), entao `workflow.state.waiting` e `workflow.state.resumed` ja sao capturados automaticamente. Nenhuma alteracao necessaria.

**8. `module.go`** — sem mudancas

O DIG injeta Publisher automaticamente se disponivel no container. Basta registrar no bootstrap.

### Resumo de impacto

```
Arquivos MODIFICADOS (6):
  1. interfaces/message/types.go         — +TimerExpiresAt, +TimerNodes, ActiveNodeID→ActiveNodeIDs
  2. domain/repositories/types.go        — ActiveNodeID→ActiveNodeIDs, +WaitingUpdate struct
  3. domain/repositories/archive_repository.go — +BulkUpdateWaiting, +BulkUpdateResumed
  4. infrastructure/.../archive_repository.go  — implementar 2 novos metodos
  5. application/di/archiver_service_di.go     — +Publisher dependencia
  6. application/services/archiver_service.go  — +2 cases, +2 batches, +short timer scheduling

Arquivos SEM mudancas (3):
  7. interfaces/message/consumers/.../consumer.go — wildcard ja cobre novos subjects
  8. application/constants/archiver.constant.go   — BatchSize 5000 permanece
  9. module.go                                     — DIG resolve Publisher automaticamente
```

---

## Runtime — Publicar eventos

### No suspend (async node)

```go
// lifecycle.go — apos checkpoint KV
func (s *RuntimeService) publishWaitingEvent(instance *entities.WorkflowInstance) {
    var earliest *time.Time
    var timerNodes []archiverMsg.TimerNode

    // Iterar TODOS os nodes waiting — pode ter multiplos em fanout
    for nodeId, ns := range instance.NodeStates {
        if exp, ok := ns["expiresAt"].(time.Time); ok {
            timerNodes = append(timerNodes, archiverMsg.TimerNode{
                NodeID:    nodeId,
                ExpiresAt: exp,
            })
            if earliest == nil || exp.Before(*earliest) {
                earliest = &exp
            }
        }
    }

    event := archiverMsg.StateEvent{
        InstanceID:     instance.ID.Hex(),
        Status:         "waiting",
        ActiveNodeIDs:  instance.ActiveNodeIDs,    // TODOS os nodes waiting
        TimerExpiresAt: earliest,                   // earliest para MongoDB index
        TimerNodes:     timerNodes,                  // cada node com seu expiresAt
    }
    s.deps.RuntimePublisher.PublishStateEvent("workflow.state.waiting", event)
}
```

### No resume

```go
// lifecycle.go — apos KV Get e antes de executar
func (s *RuntimeService) publishResumedEvent(instance *entities.WorkflowInstance) {
    event := archiverMsg.StateEvent{
        InstanceID: instance.ID.Hex(),
        Status:     "running",
    }
    s.deps.RuntimePublisher.PublishStateEvent("workflow.state.resumed", event)
}
```

---

## Flow completo — delay(10s)

```
T=0.000s  Runtime executa start → delay(10s) → SUSPENDE
          KV Put: { status:"waiting", nodeStates:{"n-2":{"waitType":"timer","expiresAt":"T+10s"}} }
          NATS Pub: "workflow.state.waiting" { instanceId, timerExpiresAt:"T+10s", activeNodeIds:["n-2"] }
          Runtime: 0 CPU, 0 RAM. Goroutine encerrada. ACK.

T=0.050s  Archiver (bulk):
          MongoDB UpdateOne { status:"waiting", timerExpiresAt: T+10s, activeNodeIds: ["n-2"] }
          TimerNodes[0]: n-2 expiresAt T+10s < 1min → NATS scheduled message:
            subject: "workflow.resume.timer.{instanceId}"
            delivery: T+10s
            payload: { instanceId, nodeId:"n-2", reason:"timer_expired" }

T=10.00s  NATS entrega scheduled message → Runtime resume consumer
          Payload inclui nodeId:"n-2" → Runtime sabe qual node resumir
          KV Get → instance com nodeStates["n-2"] = waiting
          Confirma waitType presente → processa resume do node n-2
          NATS Pub: "workflow.state.resumed" { instanceId }
          Continua executando: set_state → condition → end
          NATS Pub: "workflow.state.completed" { instanceId }

T=10.05s  Archiver (bulk):
          MongoDB UpdateOne { status:"running", timerExpiresAt: null }    ← AWAIT FINISH
          MongoDB Upsert FULL { ...todo o estado... }                     ← COMPLETE
          KV Delete(instanceId)
```

## Flow completo — delay(2h)

```
T=0.000s  Runtime: start → delay(2h) → SUSPENDE
          NATS Pub: "workflow.state.waiting" { timerExpiresAt:"T+2h" }

T=0.050s  Archiver:
          MongoDB UpdateOne { status:"waiting", timerExpiresAt: T+2h }
          timerExpiresAt <= now + 1min? NAO → nada mais (sweep pega depois)

T=1h59m   Reconciler sweep (a cada 1min):
          Query: { timerExpiresAt: { $lte: now + 1min } }
          Encontra timer para T+2h (agora esta dentro de 1min)
          NATS scheduled message: delivery T+2h (exato)

T=2h00m   NATS entrega → Runtime resume → continua executando
          Archiver: AWAIT FINISH + COMPLETE + KV Delete
```

---

## Numeros a 1M+

```
MongoDB writes extras (AWAIT START + AWAIT FINISH):
  ~4K/s × 150B + ~4K/s × 100B = ~1MB/s
  vs existente: ~244MB/s
  Aumento: +0.4% (desprezivel)

Sweep queries:
  1/min com partial index = ~1ms por query
  1,440 queries/dia (vs 86,400 com ticker 1s)

NATS scheduled messages:
  Short timers: ~4K/s (capacidade NATS: 10M+)
  Long timers: ~0 (ficam no MongoDB ate ultimo minuto)

Partial index RAM:
  ~50K timer docs × ~100B = ~5MB
```

---

## Comparacao com grandes players

```
╔═══════════════╦══════════════════════════════════════════════╗
║ Sistema       ║ Abordagem                                    ║
╠═══════════════╬══════════════════════════════════════════════╣
║ Temporal.io   ║ Timer Queue em Cassandra + Transfer Queue    ║
║               ║ Poll DB para timers + dispatch imediato      ║
║               ║ → NOSSO: MongoDB partial index + NATS direct ║
╠═══════════════╬══════════════════════════════════════════════╣
║ Google Cloud  ║ Spanner para timer persistence               ║
║ Workflows     ║ Internal scheduler para delivery precisa     ║
║               ║ → NOSSO: MongoDB + NATS scheduled messages   ║
╠═══════════════╬══════════════════════════════════════════════╣
║ Netflix       ║ Redis/Dynomite (1000 nodes!)                 ║
║ Conductor     ║ Sweep polling                                ║
║               ║ → NOSSO: 1 MongoDB + 1 NATS (2 componentes) ║
╚═══════════════╩══════════════════════════════════════════════╝
```

---

## Validacao enterprise

| Criterio | Status | Detalhe |
|----------|--------|---------|
| Escala 1M+ | ✅ | Partial index 5MB, query 1ms, +0.4% writes |
| Barato | ✅ | Zero infra nova, +1MB/s sobre existente |
| Seguro | ✅ | Dupla persistencia (NATS+MongoDB), KV idempotencia |
| Enterprise pattern | ✅ | Mesmo design Temporal/Google (DB timers + queue dispatch) |
| Scale UP pods | ✅ | QueueGroup distribui sweep + scheduled messages |
| Scale DOWN pods | ✅ | NATS JetStream reentrega, zero estado local |

---

## Checklist de implementacao

```
No packages/infrastructure/nats (task #13):
  □ types.go → Config: adicionar campo OnConnect func(c *Client) (opcional)
  □ nats.go → New(): chamar OnConnect apos connect
  □ nats.go → New(): registrar nc.SetReconnectHandler para chamar OnConnect no reconnect
  □ Manter backward-compatible: OnConnect = nil → comportamento atual (nada muda)

No entity (task #10):
  □ workflow_instance.go → adicionar TimerExpiresAt *time.Time

No Fanout (task #11):
  □ fanout.go → gravar TODOS os branches waiting (nao so o primeiro)
  □ Remover break no loop de branches waiting

No Archiver (task #12):
  □ interfaces/message/types.go:
    - ActiveNodeID string → ActiveNodeIDs []string
    - Adicionar TimerExpiresAt *time.Time
    - Adicionar TimerNodes []TimerNode
    - Adicionar struct TimerNode { NodeID, ExpiresAt }
  □ domain/repositories/:
    - Adicionar BulkUpdateWaiting(ctx, []WaitingUpdate) error
    - Adicionar BulkUpdateResumed(ctx, []string) error
  □ application/di/archiver_service_di.go:
    - Adicionar Publisher natsModel.Publisher (para scheduled messages)
  □ application/services/archiver_service.go:
    - Case isWaitingEvent: acumula TimerUpdate + short timer NATS schedule com nodeId
    - Case isResumedEvent: acumula instanceId para BulkUpdateResumed
    - Batch 2: BulkUpdateWaiting
    - Batch 3: BulkUpdateResumed

No Runtime (task #10):
  □ lifecycle.go → publishWaitingEvent() com TimerNodes array
  □ lifecycle.go → publishResumedEvent() no HandleResume

No Reconciler (task #10):
  □ Reescrever reconciler_service.go:
    - Remover: TimerQueue (heap), sweepLoop(), recoverFromKV()
    - Adicionar: HandleSweep() com MongoDB query
    - KV Get para verificar + extrair nodeIds dos NodeStates
    - NATS scheduled messages com nodeId individual
  □ Remover: application/types/timer_queue.go (heap inteiro)

No Bootstrap (task #10):
  □ nats.go → Config.OnConnect: publicar "workflow.reconciler.sweep" (trigger inicial)
  □ nats.go → consumer WORKFLOW-RECONCILER (self-schedule)
  □ MongoDB index: idx_timer_expires (partial, timerExpiresAt)

No docker-compose nats-init:
  □ Verificar WORKFLOW-RECONCILER stream criado

Resume message (todos os publishers):
  □ SEMPRE incluir nodeId no payload: { instanceId, nodeId, reason }
```

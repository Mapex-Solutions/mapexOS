# State Persistence — NATS KV + MongoDB + Archiver

Decisões: A1, A7, A9

---

## Problema

O workflow engine precisa persistir o estado de instâncias em execução de forma que:

1. **Runtime não bloqueie em I/O lento** — se MongoDB estiver lento, os workflows continuam executando
2. **Recovery automático** — se um worker crashar, outro retoma sem perda de dados
3. **Escala para 1M+ instâncias simultâneas** — sem bottleneck de RAM ou disco
4. **Separação de concerns** — quem executa não grava no banco, quem grava não executa

A maioria dos sistemas modernos (Temporal, n8n) usa DB como primary e evita Redis como state store. Netflix Conductor usou Redis como primary, mas precisou de 1000 nodes Dynomite.

---

## Como Estamos Resolvendo

### Escola 2 — Pointer mínimo, DB como dono

**NATS KV** guarda o estado COMPLETO durante a execução (~1-5KB por instância). File-backed (persiste em disco), CAS nativo via revision number. Delete quando o workflow termina.

**NATS KV é persistente.** Pode reiniciar o processo NATS quantas vezes quiser — dados intactos no disco. Single node ou cluster, dados sobrevivem a restart. Único cenário de perda é corrupção física do disco (mesmo risco que MongoDB single node).

**MongoDB** recebe writes em **4 momentos** (decisão A10 + A11):
1. **START** → InsertOne leve (~200B) — workflow aparece na listagem imediata
2. **AWAIT START** → UpdateOne (~150B) — SOMENTE async nodes. Status "waiting" + timerExpiresAt
3. **AWAIT FINISH** → UpdateOne (~100B) — SOMENTE async nodes. Status "running" + timerExpiresAt null
4. **TERMINAL** (completed/failed/cancelled) → Upsert FULL (~5-25KB) + KV Delete (cleanup)

**Writes intermediárias mínimas.** Steps 2 e 3 só ocorrem para async nodes (~20% dos workflows). Workflows inline (sem async) continuam com apenas 2 writes (START + TERMINAL). KV é a verdade durante execução. Os writes intermediários existem exclusivamente para o Reconciler encontrar timers no MongoDB e para a listagem mostrar status "waiting" corretamente.

Ver `13-reconciler-timers.md` para design completo do Reconciler.

**Redis** NÃO é usado. NATS KV substitui completamente (CAS nativo + disco barato).

### Queries — Lógica híbrida (decisão A10)

```
GetInstances (listagem):
  → MongoDB puro (paginado, rápido)
  → Mostra: workflowName, status, currentNodeId, created, orgId
  → Status vem do Created (mostra "running") — não atualiza durante execução
  → NÃO bate no KV (seria N queries, não escala)

GetInstanceById (detalhe):
  → Se status TERMINAL (completed/failed/cancelled) → MongoDB (documento permanente FULL)
  → Se status NÃO terminal → KV Get (estado fresco: DAG position, state atual, executionPath)
    + ClickHouse (logs detalhados por step — igual Temporal.io)
  → Fallback: se KV falhar ou key não existe → retorna MongoDB

Motivo: KV é a "verdade" durante execução. MongoDB só tem o stub leve do Created.
Para detalhe individual (click na instância), o usuário espera ver o estado REAL — vem do KV.
```

### GetInstanceById — Hybrid KV/MongoDB (implementado)

```
Arquivo: instances/application/services/instances_service.go

GetInstanceById(ctx, instanceId):
  ┌─ 1. MongoDB FindById(instanceId)
  │     → Sempre retorna algo (Archiver grava InsertOne no Created)
  │     → Documento pode ser: stub leve (200B) ou FULL (5-25KB)
  │
  ├─ 2. Se status é TERMINAL (completed, failed, cancelled):
  │     → Retorna MongoDB (documento FULL, permanente, KV já foi deletado)
  │
  ├─ 3. Se status NÃO é terminal (running, waiting):
  │     → KV Get("inst:{instanceId}")
  │     → Se KV hit: retorna KV (estado fresco: DAG, state, executionPath)
  │     → Se KV miss: retorna MongoDB (fallback — KV indisponível)
  │
  └─ 4. Se MongoDB retornou nil → retorna nil (instância não existe)
```

**Cascata de prioridade:**

| Status | Fonte | Motivo |
|--------|-------|--------|
| `completed` / `failed` / `cancelled` | MongoDB | KV deletado após terminal. MongoDB tem doc FULL. |
| `running` / `waiting` | NATS KV | KV é a "verdade" durante execução. MongoDB só tem stub leve do Created (ou lightweight update do waiting). |
| `running` / `waiting` (KV miss) | MongoDB | Fallback: KV temporariamente indisponível. MongoDB mostra stub (menos dados, mas funcional). |

**Por que NÃO usar MongoDB como fonte única:**
- Durante execução, MongoDB só tem o stub do Created (~200B) + possible waiting update.
- KV tem o estado COMPLETO (state, executionPath, nodeOutputs) — atualizado per-step.
- Retornar MongoDB para running/waiting mostraria dados incompletos no frontend.

**Por que NÃO usar KV como fonte única:**
- KV é deletado após terminal (cleanup do Archiver).
- MongoDB guarda o documento FULL permanente para histórico e audit.

### Workflows long-running (decisão A10)

Cenário: `start → wait_signal → process → set_state → wait_signal (loop back)`

Este workflow **nunca completa naturalmente**. Como fica:
- MongoDB tem o stub leve do Created (~200B) → aparece na listagem como "running"
- KV mantém o hot state completo (nunca limpo enquanto ativo) → persistente em disco
- Frontend vê na listagem via MongoDB (status="running")
- Detalhe individual busca KV → estado fresco (DAG, state, executionPath)
- KV é file-backed → dados intactos mesmo com restart do NATS

**Status na listagem:** Com o AWAIT START (decisão A11), MongoDB agora mostra `status: "waiting"` corretamente para workflows suspensos em async nodes. A listagem reflete o estado real sem necessidade de sweep periódico.

### Por que NÃO gravar checkpoint FULL no MongoDB (análise de scale A10)

A decisão anterior (checkpoint FULL a cada async) foi reavaliada por impacto em 1M+ scale:

```
COM checkpoint FULL no MongoDB:
  Created:    22K/s × 200B  = 4.4MB/s
  Checkpoint: 66K/s × 15KB  = 990MB/s    ← PROBLEMA: ~165MB/s comprimido
  Terminal:   16K/s × 15KB  = 240MB/s
  Total:      ~210MB/s comprimido
  Archiver KV Gets: ~82K/s

COM lightweight AWAIT START/FINISH (decisão A11):
  Created:        22K/s × 200B  = 4.4MB/s
  Await Start:    ~4K/s × 150B  = 600KB/s   ← LEVE: só status + timerExpiresAt
  Await Finish:   ~4K/s × 100B  = 400KB/s   ← LEVE: só status + null timer
  Terminal:       16K/s × 15KB  = 240MB/s
  Total:          ~42MB/s comprimido (+1MB/s = +0.4%)
  Archiver KV Gets: ~16K/s (só terminals)
```

Checkpoint FULL multiplicaria volume MongoDB por 5x e exigiria 2-3 shards vs 1.
Lightweight AWAIT START/FINISH adiciona apenas +0.4% e resolve o problema do Reconciler.
GetInstanceById ainda resolve long-running via KV read.
KV é persistente (file-backed) — mesma segurança que MongoDB.

### Módulos separados (A9)

```
services/workflow/src/modules/
  ├── runtime/       → Executa workflows (consome WORKFLOW-TRIGGER + WORKFLOW-RESUME)
  │                    Lê/escreve: NATS KV
  │                    Publica: NATS stream WORKFLOW-STATE
  │                    MongoDB: NUNCA
  │
  └── archiver/      → Persiste no MongoDB (consome WORKFLOW-STATE)
                       Lê: NATS stream WORKFLOW-STATE + NATS KV (só para terminals)
                       Escreve: MongoDB BulkWrite (Created=InsertOne leve, Terminal=Upsert FULL)
                       Deleta: NATS KV (cleanup após terminal — completed/failed/cancelled)
```

### Archiver = Writer simplificado (A7, A10)

O Archiver recebe **2 tipos de evento**: Created (leve) e Terminal (completo). Para terminals, faz KV Get para obter o estado FULL antes de gravar no MongoDB. BulkWrite direto, sem patches incrementais, sem ordering complexo.

### Workers separados por responsabilidade (A9)

| Consumer | Stream | Responsabilidade |
|---|---|---|
| Runtime (trigger + resume) | WORKFLOW-TRIGGER + WORKFLOW-RESUME | Executa nodes, lê/escreve NATS KV |
| Archiver (state writer) | WORKFLOW-STATE | BulkWrite MongoDB, cleanup NATS KV |
| Logs-writer | WORKFLOW-LOGS | Batch insert ClickHouse |
| Reconciler | WORKFLOW-RECONCILER | Sweep timers |

Se ClickHouse lento → logs-writer acumula no NATS → Runtime não para.
Se MongoDB lento → Archiver acumula no NATS → Runtime não para.

---

## Como Implementar

### NATS Stream: WORKFLOW-STATE

```
Stream: WORKFLOW-STATE
Subjects:
  workflow.state.created      → instância criada (MongoDB InsertOne leve)
  workflow.state.waiting      → async node suspendeu (MongoDB UpdateOne: status + timerExpiresAt)
  workflow.state.resumed      → async node acordou (MongoDB UpdateOne: status running + timer null)
  workflow.state.completed    → terminou com sucesso (MongoDB Upsert FULL + KV Delete)
  workflow.state.failed       → falhou (MongoDB Upsert FULL + KV Delete)
  workflow.state.cancelled    → cancelado (MongoDB Upsert FULL + KV Delete)
Storage: file
Retention: work (removida após ACK)

NOTA: "waiting" e "resumed" são lightweight UpdateOne (~150B e ~100B).
NÃO são checkpoint FULL — somente status, timerExpiresAt e activeNodeIds.
KV per-step continua sendo para crash recovery (interno entre Runtime workers).
```

### Consumer Archiver

```go
bus.StartConsumer(natsModel.ConsumerOptions{
    Stream:       "WORKFLOW-STATE",
    Subject:      "workflow.state.>",
    Durable:      fmt.Sprintf("%s-archiver", serviceName),
    QueueGroup:   fmt.Sprintf("%s-WORKFLOW-STATE-GROUP", serviceName),
    BatchSize:    5000,
    FetchTimeout: 500 * time.Millisecond,

    RetryPolicy: &natsModel.RetryPolicy{
        MaxRetries: 5,
        Backoff: []time.Duration{
            1 * time.Second,
            5 * time.Second,
            30 * time.Second,
            2 * time.Minute,
            10 * time.Minute,
        },
        AckWait: 30 * time.Second,
    },

    DLQPolicy: &natsModel.DLQPolicy{
        Stream:      "MAPEXOS-DLQ",
        Subject:     "dlq.mapexos",
        ServiceName: serviceName,
        ServiceType: "workflow",
        EventType:   "state.archiver",
    },

    BatchMessageHandlerV2: func(messages []*natsModel.Message) {
        archiverService.ProcessStateBatch(messages)
    },
})
```

### Flow completo — exemplo concreto

Workflow: `start → set_state → condition → code(ASYNC) → set_state → end`

```
═══════ PASSO 1: TRIGGER ═══════

Runtime:
  1. Recebe trigger do NATS stream WORKFLOW-TRIGGER
  2. NATS KV Put("inst:123", {state:{}, executionPath:[], version:0})
  3. NATS Publish("workflow.state.created", {
       instanceId: "123",
       workflowId: "wf-abc",
       orgId: "org-1",
       status: "running",
       currentNodeId: "start",
       version: 0
     })
  4. Começa execução inline

Archiver (batch):
  → Recebe msg "workflow.state.created"
  → Acumula no batch com outras msgs
  → MongoDB BulkWrite:
      InsertOne({ _id:"123", workflowId:"wf-abc", orgId:"org-1",
                  status:"running", currentNodeId:"start", version:0 })
  → 200 bytes gravados. ACK.


═══════ PASSO 2: EXECUÇÃO INLINE (KV per-step) ═══════

Runtime:
  1. Executa: start → KV Put → set_state → KV Put → condition → KV Put (mesma goroutine)
  2. Cada step: aplica resultado, NATS KV Put("inst:123", estado atualizado)
     → ~1ms overhead por step (NATS KV é file-backed, Put síncrono)
     → CAS: cada Put usa revision do Get/Put anterior
  3. Step logs: NATS Publish fire-and-forget → ClickHouse
  4. NÃO publica no WORKFLOW-STATE por step (inline, não é checkpoint para Archiver)
  5. Chega no "code" → ASYNC → PARA

  Motivo do KV per-step:
    Operações como increment, append, remove são DESTRUTIVAS.
    Se o pod crashar no step 4 e re-executar do step 1,
    um increment(+1) executado 2x produz resultado errado (valor +2).
    KV per-step garante que após crash, o recovery continua do último step
    completado — nunca re-executa steps já persistidos.


═══════ PASSO 3: ASYNC NODE (suspende execução) ═══════

Runtime:
  1. KV já tem o estado atualizado até o step anterior (per-step Put)
  2. Último KV Put inclui status="waiting", currentNodeId="code"
     KV revision: 4 (avançou a cada step: 0 → 1 → 2 → 3 → 4)
  3. Publica request pro JS-Executor via NATS
  4. Publica "workflow.state.waiting" no WORKFLOW-STATE (Archiver atualiza MongoDB)
  5. ACK da transition original

Archiver (batch):
  → Recebe msg "workflow.state.waiting"
  → MongoDB UpdateOne: { status:"waiting", timerExpiresAt: ..., activeNodeIds: ["code"] }
  → SE timerExpiresAt <= now + 1min: NATS scheduled message direto ao WORKFLOW-RESUME
  → ~150 bytes gravados. ACK.

  MongoDB agora mostra status="waiting" na listagem.
  GetInstanceById para esta instância → busca KV → mostra estado completo.


═══════ PASSO 4: CALLBACK ═══════

Runtime:
  1. Recebe callback do JS-Executor via NATS
  2. NATS KV Get("inst:123") → state completo (~1ms)
  3. CAS: revision do NATS KV bate? ✓
  4. Publica "workflow.state.resumed" no WORKFLOW-STATE
  5. Executa inline: code_result → set_state → end
  6. Step logs: NATS Publish fire-and-forget
  7. Workflow completou!

Archiver (batch):
  → Recebe msg "workflow.state.resumed"
  → MongoDB UpdateOne: { status:"running", timerExpiresAt: null }
  → ~100 bytes gravados. ACK.


═══════ PASSO 5: COMPLETION ═══════

Runtime:
  1. NATS KV Put("inst:123", {
       state: {counter:10, result:"ok"},
       executionPath: [ ...6 entries completas... ],
       version: 2
     })
  2. NATS Publish("workflow.state.completed", {
       instanceId: "123",
       version: 2
     })
  3. ACK

Archiver (batch):
  → Recebe msg "workflow.state.completed"
  → NATS KV Get("inst:123") → state + executionPath completo
  → MongoDB BulkWrite:
      UpdateOne({_id:"123"}, {$set: {
        status: "completed",
        state: {counter:10, result:"ok"},
        executionPath: [ ...6 entries... ],
        version: 2,
        completedAt: now
      }})
  → ~5-25KB gravados (documento completo, permanente)
  → NATS KV Delete("inst:123") ← CLEANUP
  → ACK
```

### Análise de resiliência — cenários de falha

```
╔══════════════════════════════════════╦═════════════════════════════════════════════╗
║ Cenário de falha                     ║ O que acontece                              ║
╠══════════════════════════════════════╬═════════════════════════════════════════════╣
║                                      ║                                             ║
║ 1. Runtime pod crash                 ║ Transition no NATS não teve ACK             ║
║    (mid-execution inline)            ║ → NATS reentrega para outro worker          ║
║                                      ║ → Worker lê state do NATS KV                ║
║                                      ║ → KV tem o estado do ÚLTIMO STEP completado ║
║                                      ║ → Continua do próximo step (não re-executa) ║
║                                      ║ → CAS (revision) previne duplicação         ║
║                                      ║ → RECOVERY AUTOMÁTICO                       ║
║                                      ║                                             ║
╠══════════════════════════════════════╬═════════════════════════════════════════════╣
║                                      ║                                             ║
║ 2. Archiver pod crash               ║ BulkWrite não completou → msgs sem ACK      ║
║    (mid-batch)                       ║ → NATS reentrega batch para outro Archiver  ║
║                                      ║ → BulkWrite é idempotente (upsert by _id)  ║
║                                      ║ → RECOVERY AUTOMÁTICO                       ║
║                                      ║                                             ║
╠══════════════════════════════════════╬═════════════════════════════════════════════╣
║                                      ║                                             ║
║ 3. NATS reinicia                     ║ JetStream persiste em disco → dados OK      ║
║                                      ║ → NATS KV persiste em disco → state OK      ║
║                                      ║ → Após restart: consumers reconectam        ║
║                                      ║ → Msgs unACKed são reentregues              ║
║                                      ║ → RECOVERY AUTOMÁTICO                       ║
║                                      ║                                             ║
╠══════════════════════════════════════╬═════════════════════════════════════════════╣
║                                      ║                                             ║
║ 4a. NATS restart                     ║ NATS KV é FILE-BACKED (storage=file)        ║
║     (processo reinicia, disco OK)    ║ → Lê dados do disco → state intacto        ║
║                                      ║ → Consumers reconectam → msgs reentregues  ║
║                                      ║ → RECOVERY AUTOMÁTICO                       ║
║                                      ║                                             ║
╠══════════════════════════════════════╬═════════════════════════════════════════════╣
║                                      ║                                             ║
║ 4b. NATS disco corrompido           ║ Cluster (3 réplicas): 1 nó perde disco     ║
║     (falha física/corrupção)         ║ → 2 réplicas intactas → zero perda         ║
║                                      ║ → Nó corrompido re-sincroniza do cluster   ║
║                                      ║ → RECOVERY AUTOMÁTICO                       ║
║                                      ║ Single node: disco corrompido = dados lost ║
║                                      ║ → Mesmo risco que MongoDB single node      ║
║                                      ║ → Mitigação padrão: RAID + backups         ║
║                                      ║                                             ║
╠══════════════════════════════════════╬═════════════════════════════════════════════╣
║                                      ║                                             ║
║ 5. MongoDB cai                       ║ Runtime: ZERO impacto (não toca MongoDB)    ║
║                                      ║ Archiver: acumula msgs no NATS stream       ║
║                                      ║ Frontend: listagem para de atualizar        ║
║                                      ║ Quando MongoDB volta: Archiver drena fila   ║
║                                      ║ → RUNTIME NÃO PARA                          ║
║                                      ║ → RECOVERY AUTOMÁTICO                       ║
║                                      ║                                             ║
╠══════════════════════════════════════╬═════════════════════════════════════════════╣
║                                      ║                                             ║
║ 6. MongoDB lento (alta latência)     ║ Runtime: ZERO impacto                       ║
║                                      ║ Archiver: BulkWrite demora mais             ║
║                                      ║ → Processa menos batches/s                  ║
║                                      ║ → Acumula no NATS (buffer natural)          ║
║                                      ║ → Listagem frontend fica stale (5-30s)      ║
║                                      ║ → Quando normaliza: drena acumulado         ║
║                                      ║ → DEGRADAÇÃO GRACEFUL                       ║
║                                      ║                                             ║
╠══════════════════════════════════════╬═════════════════════════════════════════════╣
║                                      ║                                             ║
║ 7. Mensagem "envenenada"             ║ RetryPolicy: 5 tentativas com backoff       ║
║    (bad payload, bug)                ║ → [1s, 5s, 30s, 2min, 10min]                ║
║                                      ║ → Após 5 falhas: msg vai pro MAPEXOS-DLQ    ║
║                                      ║ → DLQ consumer armazena para debug          ║
║                                      ║ → Outras msgs continuam normalmente         ║
║                                      ║ → ISOLAMENTO DE FALHA                       ║
║                                      ║                                             ║
╠══════════════════════════════════════╬═════════════════════════════════════════════╣
║                                      ║                                             ║
║ 8. Network partition                 ║ KV Put falha → step execution para          ║
║    (Runtime ↔ NATS)                  ║ → Runtime detecta erro no KV Put            ║
║                                      ║ → Não faz ACK da transition                 ║
║                                      ║ → NATS reentrega após partition resolver    ║
║                                      ║ → Worker retoma do último step persistido   ║
║                                      ║ → RECOVERY AUTOMÁTICO                       ║
║                                      ║                                             ║
╚══════════════════════════════════════╩═════════════════════════════════════════════╝
```

**NATS KV é file-backed (JetStream storage=file).** Não é volátil como Redis — dados persistem em disco. Reinício do NATS = dados intactos. O único cenário de perda é corrupção física do disco — mesmo risco do MongoDB.

**Contrato de persistência:**
- **Persistência:** garantida pelo disco (file-backed). Single node ou cluster, dados sobrevivem a restart.
- **Disponibilidade:** single node = downtime durante restart (NATS fora = Runtime parado). Cluster (3 nós) = HA, zero downtime em falha de 1 nó.
- **Redundância de disco:** cluster replica dados entre nós. Single node = single point of failure de disco.

**Produção:** NATS cluster (3 nós) obrigatório para **disponibilidade** e **redundância de disco**.
**Standalone (dev/staging):** single node aceitável — dados em disco, mas sem HA nem redundância.

### Números para 1M+ (premissas-alvo, validar com benchmark)

```
Premissas: 1M workflows ativos | 15 nodes média | lifecycle 60s | KV Put per-step

NOTA: Números abaixo são estimativas de dimensionamento para guiar sizing de infra.
Validar com benchmark real antes de produção. Capacidades de NATS/MongoDB dependem
de hardware, rede e configuração.

NATS KV:
  1M × 5KB = 5GB disco
  1M × 15 steps / 60s = ~250K Puts/s + ~16K Gets/s (callbacks) = ~266K ops/s
  Index em memória: ~200MB
  (Referência NATS: benchmarks oficiais reportam 10M+ msgs/s em hardware adequado)
  (250K Puts/s ≈ 25% capacidade de um cluster NATS 3 nós em hardware modesto)

MongoDB (via Archiver BulkWrite — Created + Await + Terminal):
  Writes LEVES (created):     ~22K/s × 200 bytes = ~4.4MB/s
  Writes AWAIT START:         ~4K/s  × 150 bytes = ~600KB/s   (async only)
  Writes AWAIT FINISH:        ~4K/s  × 100 bytes = ~400KB/s   (async only)
  Writes FULL (terminal):     ~16K/s × 15KB avg  = ~240MB/s bruto → com zstd: ~40MB/s
  Total:                      ~45MB/s comprimido (+1MB/s sobre anterior)
  BulkWrite 5000:             ~8 BulkWrites/s (created + terminal) + ~1/s (await)
  Archiver KV Gets:           ~16K/s (só terminals)

  Single replica set estimado suficiente até lifecycle 30s

NATS streams:
  WORKFLOW-TRIGGER + RESUME: ~22K triggers/s + ~44K resumes/s
  WORKFLOW-STATE:       ~22K creates + ~16K terminals = ~38K msgs/s (sem checkpoints)
  WORKFLOW-LOGS:        ~330K step logs/s
  Total: ~444K msgs/s
```

---

## Evolução Futura: Event Sourcing

### Atual: KV Per-Step Snapshot

O NATS KV guarda o **estado completo** a cada step. Após crash, o worker lê o último snapshot e continua. Simples, rápido, confiável.

Limitação: **não há histórico de mudanças**. Se quiser saber "o que mudou entre step 3 e step 7", ou fazer replay/audit, não temos essa informação — cada Put sobrescreve o anterior.

### Futuro: Event Sourcing via NATS JetStream

NATS JetStream **já é um event store** — append-only, ordenado, persistido em disco. A evolução é natural:

```
Atual (snapshot only):
  step → apply → KV Put(estado completo)

Com Event Sourcing (snapshot + event log):
  step → apply → KV Put(estado completo) + JetStream Publish(evento do step)
                 ↑ hot state (query)       ↑ history (audit/replay)
```

**Stream de eventos:**
```
Stream: WORKFLOW-HISTORY
Subject: workflow.history.inst.{instanceId}
Storage: file
Retention: limits (por tempo ou por número de mensagens)
MaxAge: 30d (configurável)
```

**Evento por step:**
```json
{
  "instanceId": "123",
  "stepIndex": 3,
  "nodeId": "set_state_1",
  "nodeType": "set_state",
  "status": "success",
  "durationMs": 0,
  "statePatch": {"counter": {"op": "set", "value": 5}},
  "outputHandle": "default",
  "timestamp": "2026-03-07T..."
}
```

**Mudança no runtime:**
```go
/* Atual */
r.kvStore.Put(instance)

/* Com Event Sourcing — adiciona 1 linha */
r.kvStore.Put(instance)
go r.eventStream.Publish(stepEvent)  // fire-and-forget, async
```

Uma única linha de código. O `go publish()` é assíncrono — se falhar, o KV snapshot já está salvo. Eventos são best-effort para audit, não críticos para execução.

### O que Event Sourcing habilita

| Capacidade | Como |
|---|---|
| **Replay** | Ler eventos do JetStream, re-aplicar em ordem → reconstituir estado em qualquer ponto |
| **Audit trail** | Consultar todos os eventos de uma instância por subject filter |
| **Time-travel debug** | "Mostre o estado da instância no step 5" → replay eventos 1-5 |
| **Analytics** | Consumir stream com analytics consumer → métricas por node type, duração, etc. |
| **Compensação** | Em caso de bug, identificar instâncias afetadas pelo evento errado |

### Por que NÃO fazer Event Sourcing agora

1. **Complexidade zero justificada** — Snapshot já resolve crash recovery, o caso crítico
2. **Sem demanda** — Audit trail e replay são features enterprise
3. **Evolução sem rewrite** — Quando necessário, é literalmente 1 linha de código no runtime
4. **NATS JetStream já está lá** — Zero infra adicional

### Decisão

Snapshot em produção. Event Sourcing quando houver demanda real por audit/replay. A arquitetura suporta a evolução sem modificar o runtime loop — apenas adicionar o publish.

### Checklist de implementação

```
No workflow service (services/workflow/):
  □ modules/archiver/          — Consumer WORKFLOW-STATE, BulkWrite MongoDB, cleanup NATS KV
  □ modules/runtime/           — Consumer WORKFLOW-TRIGGER + WORKFLOW-RESUME, lê/escreve NATS KV
  □ modules/logswriter/        — Consumer WORKFLOW-LOGS, batch insert ClickHouse
  □ modules/reconciler/        — Consumer WORKFLOW-RECONCILER, sweep timers
  □ bootstrap/nats.go          — Inicializar NATS KV bucket para workflow instances
```

### Referências

- Estudo de concorrentes: `services/workflow/ESTUDO_CONCORRENTES.md`
- Escola 2 (Pointer mínimo) = padrão Temporal, n8n, maioria moderna
- DLQ: padrão centralizado `MAPEXOS-DLQ` com `DLQPolicy` + `RetryPolicy`
- BulkWrite: driver Go MongoDB nativo, wrapper precisa expor

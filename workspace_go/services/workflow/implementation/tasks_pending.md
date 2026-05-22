# Tasks Pendentes — Workflow Service

Status atual: **servico compila e roda**. Todos os itens abaixo sao melhorias e integracoes, nao bloqueiam o core.

---

## P0 — Necessario para execucao completa de nodes async

### ~~1. Publicar requests apos async dispatch~~ DONE

Implementado: `dispatchByNodeType()` no RuntimeService.
- `core/code` → publica `CodeExecutionRequest` no subject `workflow.js.code` (stream `WORKFLOW-JS-CODE`)
- `core/subworkflow` → publica `SubworkflowTrigger` no subject `workflow.trigger.subworkflow.{id}` (stream `WORKFLOW-TRIGGER`)
- `core/trigger_event` → publica `TriggerEventRequest` no subject `workflow.trigger.event`
- Variables defaults armazenadas em `StateDefaults` na instancia
- CallbackSubject derivado de `instance.ID.Hex()`

---

### ~~2. Verificar criacao dos 3 streams no bootstrap~~ NAO NECESSARIO

Deployment (nats-init no docker-compose) cria todos os streams. Package NATS faz ensure via `createOrGetConsumer()` como fallback. Ver task #16 para adicionar streams faltantes ao deployment.

---

### ~~8. TieredCache para WorkflowDefinition no Runtime~~ DONE

Implementado: `bootstrap/cache.go` (L0+L1+L2 MinIO), `definition_loader.go` (cache integration), FANOUT invalidation via `definition_service.go`.

---

### ~~10. Reconciler redesign: MongoDB timers + NATS self-schedule~~ DONE

Implementado: `reconciler_service.go` com HandleSweep + FindExpiredTimers + KV validation. MongoDB partial index para timers. Doc: `13-reconciler-timers.md`.

---

### ~~13. packages/infrastructure/nats: OnConnect callback~~ DONE

Implementado: `nats.go` L37-49 (OnConnect + SetReconnectHandler), `types.go` L26-33 (Config field). Campo opcional, backward-compatible.

---

### ~~11. Fanout: gravar TODOS os branches waiting (bug)~~ DONE

Corrigido: `fanout.go` L91-100 agora coleta ALL waiting nodes sem `break`. Todos os branches async sao rastreados em `ActiveNodeIDs`.

---

### ~~12. Archiver: processar waiting/resumed events + short timer fast-path~~ DONE

Implementado: `archiver_service.go` processa 4 lifecycle writes (created/waiting/resumed/terminal). Short timer fast-path para timers < 1min. DI com Publisher para scheduled messages.

---

### ~~14. Centralizar message types em src/shared/types/~~ DONE

Implementado: `shared/types/state_event.go`, `resume_message.go`, `execution_error.go`. Modulos usam type aliases (`type StateEvent = sharedTypes.StateEvent`). Single source of truth.

---

### ~~15. Graceful shutdown no packages/microservices~~ DONE

Implementado: `packages/microservices/shutdown/` (types.go, methods.go, internals.go, shutdown.go). ShutdownManager com hooks por prioridade. Workflow service integrado via `bootstrap/shutdown.go`. 10 unit tests + 5 integration tests.

---

### 16. Deployment: adicionar stream WORKFLOW-RECONCILER ao nats-init

**Problema**: 18 de 19 streams ja criados no docker-compose. Falta apenas WORKFLOW-RECONCILER.

**Arquivos a atualizar**:
```
deployment/docker-compose/standalone/infra/docker-compose.yml
deployment/docker-compose/services_required/docker-compose.yml
```

**Stream a adicionar**:
```bash
nats stream add WORKFLOW-RECONCILER --subjects="workflow.reconciler.>" --storage=file --retention=work --defaults
```

**Complexidade**: Baixa

---

### ~~17. Documentar GetInstanceById hybrid KV/MongoDB~~ DONE

Documentado: `01-state-persistence.md` L59-96 (cascading logic, priority table) + `05-instance-listing.md` L42-57 (GetInstances vs GetInstanceById).

---

## P1 — Performance e observabilidade

### ~~3. MongoDB Indexes~~ DONE

Implementado: covering indexes para instances e definitions + partial timer index. Criados via `EnsureIndexes()` no bootstrap.

---

### 4. WORKFLOW-LOGS consumer → ClickHouse

**Problema**: O stream `WORKFLOW-LOGS` esta definido nas constants. O `LogExecutor` emite `LogEntry` mas o RuntimeService **nao publica os logs no stream**. Nao existe consumer para persistir no ClickHouse.

**O que falta**:
1. RuntimeService: apos cada step, se `result.LogEntries` nao vazio → publicar no WORKFLOW-LOGS
2. Novo modulo ou consumer no Archiver: consumir WORKFLOW-LOGS → batch insert no ClickHouse
3. Schema ClickHouse para logs (ver `workspace_go/packages/seeds/clickhouse/events/v1`)

**Complexidade**: Media

---

### ~~9. ExecutionPath: atualizar status de nodes async apos resume~~ DONE

Corrigido: `runtime_service.go` L217-225 — loop reverso marca waiting→completed com ExitedAt no HandleResume.

---

## P2 — Robustez e testes

### ~~5. Testes adicionais~~ DONE

Implementado: `graph_builder_test.go` (12 tests), `archiver_service_test.go` (14 tests), `reconciler_service_test.go` (8 tests). Total: 34 novos tests com hand-written mocks.

---

### 6. Backpressure no Archiver (A6)

**Problema**: Documentado em `06-backpressure.md` — 3 modos baseados em P99 de write latency MongoDB. O Archiver faz BulkWrite mas nao tem throttling.

**Modos**:
- Normal (P99 < 50ms): batch completo
- Degraded (50ms < P99 < 200ms): reduzir batch size
- Critical (P99 > 200ms): pausar consumer, backoff

**Complexidade**: Media

---

## P3 — Features enterprise

### 7. Error handling per-node

**Problema**: Retry policy per-node. Hoje um erro num node falha o workflow inteiro.

**O que falta**:
- Campo `retryPolicy` no config de cada node
- RuntimeService: ao falhar um node, verificar retry policy antes de marcar como failed
- Error boundary nodes (catch errors de um grupo de nodes)

**Complexidade**: Alta

---

## Checklist Rapido

| # | Item | Prioridade | Status |
|---|------|-----------|--------|
| ~~1~~ | ~~Publicar code/trigger_event/subworkflow requests~~ | ~~P0~~ | ~~DONE~~ |
| ~~2~~ | ~~Streams no bootstrap~~ | ~~P0~~ | ~~NAO NECESSARIO~~ |
| ~~8~~ | ~~TieredCache L0+L1+L2 para definitions~~ | ~~P0~~ | ~~DONE~~ |
| ~~9~~ | ~~ExecutionPath status apos resume~~ | ~~P0~~ | ~~DONE~~ |
| ~~10~~ | ~~Reconciler: MongoDB timers + NATS self-schedule~~ | ~~P0~~ | ~~DONE~~ |
| ~~11~~ | ~~Fanout: gravar TODOS os branches waiting~~ | ~~P0~~ | ~~DONE~~ |
| ~~12~~ | ~~Archiver: waiting/resumed events + short timer~~ | ~~P0~~ | ~~DONE~~ |
| ~~13~~ | ~~NATS OnConnect callback~~ | ~~P0~~ | ~~DONE~~ |
| ~~14~~ | ~~Centralizar message types em shared/types~~ | ~~P0~~ | ~~DONE~~ |
| 16 | Deployment: stream WORKFLOW-RECONCILER | P0 | **PENDENTE** |
| ~~17~~ | ~~Documentar GetInstanceById hybrid~~ | ~~P0~~ | ~~DONE~~ |
| ~~3~~ | ~~MongoDB indexes~~ | ~~P1~~ | ~~DONE~~ |
| 4 | WORKFLOW-LOGS → ClickHouse | P1 | PENDENTE |
| ~~5~~ | ~~Testes adicionais~~ | ~~P2~~ | ~~DONE~~ |
| ~~15~~ | ~~Graceful shutdown~~ | ~~P1~~ | ~~DONE~~ |
| 6 | Backpressure no Archiver | P2 | PENDENTE |
| 7 | Error handling per-node | P3 | PENDENTE |

---

## O que JA funciona

- Criar/editar/deletar workflow definitions (CRUD completo)
- Trigger de workflows via NATS (consumer WORKFLOW-TRIGGER)
- Execucao inline de todo o DAG (start → condition → set_state → log → end)
- Parallel fanout com merge (goroutines) — ALL branches tracked
- Loop, sequence, goto nodes
- Delay (timer via Reconciler com MongoDB sweep)
- Wait signal (HTTP endpoint para enviar signal)
- Wait for (condition polling via Reconciler)
- Per-step KV checkpoint (crash recovery)
- Archiver persiste no MongoDB (4 lifecycle writes: created/waiting/resumed/terminal)
- Short timer fast-path (< 1min) no Archiver
- TieredCache (L0 RAM + L1 NVMe + L2 MinIO) para definitions
- FANOUT cache invalidation across pods
- GetInstances (MongoDB) e GetInstanceById (hibrido KV/MongoDB)
- Cancel instance
- ExecutionPath corretamente atualizado apos resume
- Graceful shutdown com ShutdownManager (P0 HTTP → P5 connections)
- NATS OnConnect callback para Reconciler sweep
- ~103 unit/integration tests passando

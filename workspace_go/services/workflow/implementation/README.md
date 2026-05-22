# Workflow Engine — Implementation Docs

Cada arquivo contém: o problema, como resolvemos, e como implementar.

## Arquivos

| # | Arquivo | Decisões | Conteúdo |
|---|---------|----------|----------|
| 01 | [state-persistence](./01-state-persistence.md) | A1, A7, A9 | NATS KV + MongoDB + Archiver, workers separados, resiliência |
| 02 | [nats-streams](./02-nats-streams.md) | A3, A5 | 2 streams (TRIGGER + RESUME), double-buffer, consumer configs |
| 03 | [concurrency](./03-concurrency.md) | A2, E6 | Idempotência, CAS conflict handling, DISCARD vs RETRY |
| 04 | [definition-cache](./04-definition-cache.md) | A4 | TieredCache L1/L2 MinIO + NATS FANOUT invalidation |
| 05 | [instance-listing](./05-instance-listing.md) | A10 | Covering index, data obrigatória, workflowName denormalizado |
| 06 | [backpressure](./06-backpressure.md) | A6 | MongoDB write latency tracking, P99, 3 modos |
| 07 | [inline-execution](./07-inline-execution.md) | E2 | Inline até async/end/error, max 500 steps, 30s timeout |
| 08 | [structs](./08-structs.md) | B1-B7 | Todas as Go structs: entities, configs, interfaces |
| 09 | [executors](./09-executors.md) | C1-C17 | 17 node executors: 7 inline, 5 async, 5 control |
| 10 | [execution-graph](./10-execution-graph.md) | E1, E3-E5 | BuildGraph, goto pairs, edge cases |
| 13 | [reconciler-timers](./13-reconciler-timers.md) | A11 | MongoDB timers + NATS self-schedule, 4-write lifecycle |

## Status

```
PARTE A: 10/10 ✅ (Decisões arquiteturais)
PARTE B:  7/7  ✅ (Go structs)
PARTE C: 17/17 ✅ (Node executors)
PARTE D:  0/5  ⏸️ (Schemas/infra — definir durante implementação)
PARTE E:  6/6  ✅ (Algoritmos e edge cases)
```

## PARTE D: Pendências (definir durante implementação)

### D1. MongoDB: Schema completo + indexes

**workflow_definitions:**
- Indexes para queries comuns: `{orgId: 1, enabled: 1}`, `{name: 1, orgId: 1}`

**workflow_instances:**
- Covering index para listing: ver `05-instance-listing.md`
- Indexes adicionais a definir com dados reais

### D2. NATS: Subject patterns exatos

- Subjects definidos em alto nível em `02-nats-streams.md`
- Detalhar patterns exatos durante implementação (ex: wildcards, consumer group names)

### D3. NATS: Consumer configs detalhados

- BatchSize, FetchTimeout, MaxAckPending por consumer
- Configs base definidos em `02-nats-streams.md`
- Ajustar com testes de carga

### D4. ClickHouse: Log retention tiers

- Tabela raw (TTL 3d ou 90d?)
- Rollups (1min, 1h) — precisa ou overkill?

### D5. Logs-writer consumer separado

- Stream WORKFLOW-LOGS para step logs → ClickHouse
- Reusa pattern do events service (batch insert)

## Outros documentos

- `ESTUDO_CONCORRENTES.md` — Estudo de concorrentes (Redis, Temporal, Conductor, n8n, etc.)

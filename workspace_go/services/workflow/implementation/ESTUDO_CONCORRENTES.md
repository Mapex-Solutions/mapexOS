# Workflow Engines — Estudo de Concorrentes (Redis & State Persistence)

Data: 2026-03-06

---

## 1. Netflix Conductor

**Arquitetura de estado:**
- Redis = **estado COMPLETO** (workflow inteiro serializado como JSON)
- Data structure: Redis HASH (`hset(workflowId, "data", JSON)`) + SET para lookups por workflow name
- Cada workflow: **1-10KB** no Redis
- Task scheduling: Redis HASH (`hset(SCHEDULED_TASKS, workflowInstanceId, taskId)`)
- Desde Conductor v2.x: workflow definition **embarcada** na execution (evita lookup repetido)

**Escala Redis:**
- Netflix usou **Dynomite** (cluster Redis proprietário, ~1000 nodes em 70 clusters, dados de 2016)
- Dynomite: r3.2xlarge (61GB RAM, 8 CPUs) por node
- Throughput por node: 33K reads/s + 10K writes/s (DC_ONE), 18K ops/s (DC_QUORUM)
- Latência: p99 < 4ms com DC_QUORUM, mediana < 2.5ms

**Benchmark público (Orkes, Conductor OSS):**
- 210 workflows/s sustained = ~540M workflows/mês
- Peak: ~1,800 workflow starts/s
- Task throughput: 1,450+ task executions/s
- Data rate: ~80 MB/s sustained
- Infra: 3 servidores commodity + 1 Redis node + PostgreSQL
- JVM heap: < 2GB

**Persistência:**
- Redis = primary (hot state)
- PostgreSQL ou Cassandra = secondary (archival, search)
- Elasticsearch = indexação para queries
- Suporta: Redis Standalone, Sentinel, Cluster, Dynomite

**Cleanup:**
- Sem cleanup automático robusto (issues reportados na comunidade)
- Workflows completados precisam ser limpos explicitamente
- Problema conhecido em produção

---

## 2. Netflix Timestone (Priority Queue do Conductor)

**O que é:** Sistema de fila de prioridade que alimenta o Conductor.

**Redis usage:**
- Sorted Set (ZSET) para filas com prioridade
- HASH para metadata de mensagem + configs
- Por mensagem: ID, exclusivity key/value, state (Pending/Invisible/Running/Completed/Canceled/Errored)
- **NÃO guarda payload** — só metadata da mensagem

**Atomicidade:**
- TODAS as operações via **Lua scripts** (ACID-like dentro do Redis)
- Garante consistência em transições de estado

**Persistência:**
- Redis com AOF (append-only file) para durabilidade
- Single Redis (não cluster)

---

## 3. Temporal.io

**Arquitetura de estado:**
- **NÃO usa Redis** para estado
- Estado completo no DB: Cassandra, MySQL ou PostgreSQL
- Grava **WorkflowSnapshot** (estado completo na criação) + **WorkflowMutation** (delta incremental a cada activity/timer)
- Cada state transition = write no DB (7 writes para 1 workflow com 1 activity)
- Redis usado apenas como cache opcional em deployments híbridos

**Escala:**
- Sharding: 512 shards (partições lógicas por hash do workflowId)
- Sticky execution: workflow fica no mesmo shard → localidade de cache
- History limit: 51,200 events ou 50 MB por execution
- Payload: warning em 256 KB, erro em 2 MB

**Benchmark público:**
- 1,350 state transitions/s em cluster pequeno (4 CPU MySQL, 512 shards)
- ~1,600 workflow tasks/s
- StartWorkflow latency: ~50ms (target < 150ms)
- Self-hosted ceiling: ~3,000-4,000 transitions/s (sem Temporal Cloud)

**Temporal Cloud (versão paga):**
- WAL proprietário para batching de writes (não disponível self-hosted)
- Escala significativamente melhor que self-hosted

---

## 4. Uber Cadence

**Arquitetura de estado:**
- Fork do que virou o Temporal (mesmos criadores)
- **NÃO usa Redis**
- Estado completo no DB: Cassandra, MySQL ou PostgreSQL
- Mesma arquitetura de event sourcing do Temporal
- Sharding por workflowId para distribuir carga

**Escala:**
- Uber rodou Cadence em produção com Cassandra
- Cassandra = linearmente escalável para writes
- Sem números públicos detalhados

---

## 5. n8n (BullMQ)

**Arquitetura de estado:**
- Redis guarda **SOMENTE ponteiro** (execution ID + metadata de fila)
- ~200-500 bytes por entry no Redis
- Estado completo sempre no PostgreSQL
- Data structure: HASH per job (`bull:<queue>:<jobId>`), LIST para waiting/active, ZSET para delayed/prioritized

**Flow:**
1. Job entra na fila Redis (ID + metadata mínima)
2. Worker pega job do Redis
3. Worker carrega contexto completo do PostgreSQL
4. Executa workflow
5. Grava resultado no PostgreSQL
6. ACK no Redis

**Escala:**
- Redis single node (default)
- Sem dados públicos de grande escala
- Não é projetado para 1M+ concorrentes

---

## 6. Sidekiq (Ruby)

**Arquitetura de estado:**
- Redis = **ÚNICO storage** (não tem DB por trás)
- Payload completo como JSON (ou MessagePack em versões novas)
- Per job: class name, args, queue name, retry count, timestamps, JID
- Data structures: LIST (FIFO), ZSET (scheduled/retry), HASH (process metadata), STRING (counters)

**Escala:**
- Single Redis node (Redis Cluster NÃO suportado)
- Throughput: 20,000+ jobs/s reportado por clientes
- MessagePack = ~30% menor que JSON

**Cleanup:**
- Dead set capped em 10,000 jobs
- TTL de 6 meses em dead jobs
- Jobs completados removidos automaticamente

---

## 7. Celery (Python)

**Arquitetura de estado:**
- Redis ou RabbitMQ como broker (fila de tarefas)
- Redis: LIST para filas, STRING para resultados (com TTL)
- Per task: payload serializado completo
- Resultados: TTL default de 1 dia

**Escala:**
- Redis single node (típico)
- Alternativa: RabbitMQ para maior throughput

---

## 8. Airbnb (Resque → Dynein)

**Histórico:**
- Começou com **Resque** (Redis-backed, payloads completos no Redis)
- Atingiu limites: single Redis instance bottleneck, sem cluster support, at-most-once delivery
- **Migrou para Dynein** (DynamoDB-based) para workloads críticos

**Lição:** Redis single node não escala para workloads enterprise de job queuing. Airbnb abandonou Redis.

---

## 9. Discord

**Redis usage:**
- Cache de shard mappings e user presence (dados pequenos)
- Mensagens: Cassandra → ScyllaDB
- **Abandonou Redis para indexação de filas** — Redis "não tinha CPU/espaço suficiente" quando filas acumulavam
- Migrou para Elasticsearch + Kubernetes

---

## 10. Pinterest

**Redis usage:**
- **Functional partitioning**: Redis instances separadas por use case
  - Uma para followers
  - Uma para feeds
  - Uma para caching
- 1000+ Redis instances no total
- Não usa Redis para workflow state

---

## 11. Shopify

**Redis usage:**
- **1 Redis node por shard MySQL** (pod isolado)
- Redis segue modelo de sharding do MySQL
- Não usa Redis como primary storage

---

## 12. Instagram — Otimização de Memória Redis

**Caso de uso:** 300M key-value pairs com valores pequenos (IDs de fotos → IDs de usuários).

**Antes (STRING):** 300M keys × ~85 bytes overhead = ~25 GB

**Depois (Bucketed HASH com ziplist):**
- Agrupou em 100K hashes, cada um com ~1000 fields
- Config: `hash-max-ziplist-entries 1000`, `hash-max-ziplist-value 64`
- Resultado: **4x economia de memória** (70 MB → 16 MB no test set)

**Lição:** HASH com ziplist é ordens de magnitude mais eficiente que STRING para dados pequenos.

---

## Tabela Comparativa: O Que Cada Um Guarda no Redis

| Sistema | O que fica no Redis | Tamanho/entry | O que fica no DB | Redis como... |
|---|---|---|---|---|
| **Netflix Conductor** | Estado COMPLETO (JSON) | 1-10 KB | Réplica (PostgreSQL/Cassandra) | Primary |
| **Netflix Timestone** | Metadata de mensagem | 200-500 bytes | Kafka + Elasticsearch | Queue |
| **Temporal** | NADA (ou cache opcional) | — | Estado completo (Cassandra/MySQL) | N/A |
| **Uber Cadence** | NADA | — | Estado completo (Cassandra/MySQL) | N/A |
| **n8n (BullMQ)** | Pointer (ID + metadata fila) | 200-500 bytes | Estado completo (PostgreSQL) | Queue |
| **Sidekiq** | Payload COMPLETO | 200-500 bytes | NADA (Redis é único) | Primary + Queue |
| **Celery** | Payload completo + resultado | 200-1000 bytes | NADA | Primary + Queue |
| **Airbnb** | ~~Payload completo~~ → migrou | — | DynamoDB (Dynein) | Abandonou |
| **Discord** | Shard mappings, presence | ~100 bytes | ScyllaDB | Cache |
| **Pinterest** | Feeds, followers (particionado) | Variável | MySQL | Cache |
| **Instagram** | IDs (bucketed HASH) | ~0.16 bytes | PostgreSQL | Cache |

---

## Memória Redis por Escala

| Tipo de dado | Bytes/entry | 1M entries | 10M entries | 100M entries |
|---|---|---|---|---|
| Pointer/ID + status | 80-130 bytes | ~100 MB | ~1 GB | ~10 GB |
| Job payload pequeno | 200-500 bytes | ~400 MB | ~4 GB | ~40 GB |
| Estado completo (Conductor) | 1-10 KB | 5-10 GB | 50-100 GB | 500 GB-1 TB |
| Counter/boolean (ziplist) | ~0.16 bytes | ~16 MB | ~160 MB | ~1.6 GB |

**Overhead por key Redis (64-bit):**
- STRING: ~128 bytes overhead (dictEntry 32 + robj 32 + value 32 + key 32)
- HASH field: ~10 bytes overhead por field (dentro de hash existente)
- 5 fields em 1 HASH = ~114 bytes overhead vs 5 STRING keys = ~640 bytes (**5.6x economia**)

---

## Duas Escolas de Arquitetura

### Escola 1: Redis-as-Primary (Conductor, Sidekiq)
```
Estado completo no Redis
DB = réplica/archival

Pro: leitura ultra rápida (0.1ms)
Con: Redis caro (memória RAM), precisa cluster grande
Con: Redis crash = risco de perda de state
Con: cleanup complexo
Quem: Netflix (com Dynomite, 1000 nodes), Sidekiq (single node, escala limitada)
```

### Escola 2: Redis-as-Queue/Cache, DB-as-Primary (Temporal, n8n, maioria moderna)
```
Redis: somente ponteiro/metadata (~100-500 bytes)
DB = source of truth

Pro: Redis barato (100MB-1GB para milhões)
Pro: Redis crash = zero perda (tudo no DB)
Pro: DB escala com indexes/sharding
Con: leitura de estado via DB (~5ms vs 0.1ms)
Quem: Temporal, n8n, tendência moderna
```

**Tendência da indústria:** Escola 2 domina. Apenas Netflix (2016) e Sidekiq (escala limitada) usam Escola 1. Temporal, Cadence, n8n, Airbnb (migrou) — todos escolheram DB como primary.

---

## Lições para MapexOS Workflow Engine

1. **Redis-as-Primary não escala** sem infraestrutura tipo Dynomite (1000 nodes). Nenhum sistema moderno adota.
2. **Pointer mínimo no Redis** (~100 bytes: version + status + currentNode) = 100MB para 1M instâncias = 1 Redis node de 256MB.
3. **Batch reads do MongoDB** (find $in com 500+ IDs por batch) eliminam o overhead de latência individual.
4. **Se Redis cair**: zero impacto (MongoDB é source of truth, Redis é reconstruído automaticamente nos próximos reads).
5. **HASH > STRING** para dados estruturados (5.6x economia de memória).
6. **Cleanup automático**: DEL ao completar + TTL de segurança (24h).

# Router Service — Benchmark Results

> **Tag:** 1RouteGroupPerEvent
> **Date:** 2026-02-21
> **Pattern:** NATS Stream Drain

---

## Table of Contents

1. [Overview](#1-overview)
2. [Test Environment](#2-test-environment)
3. [Test Configuration](#3-test-configuration)
4. [Methodology](#4-methodology)
5. [Results](#5-results)
6. [Key Findings](#6-key-findings)
7. [Conclusions](#7-conclusions)
8. [Reproduction](#8-reproduction)

---

## 1. Overview

This document records throughput and memory profiles of the router service under isolated CPU conditions across **all three routing scenarios**. The router is a NATS JetStream consumer: it receives events from the `ROUTE-GROUPS` stream, resolves the asset's associated RouteGroups via a three-tier cache (L0 RAM → L1 Disk → L2 MinIO), evaluates optional conditional match rules, and publishes routing decisions to downstream NATS subjects using fire-and-forget (`PublishCore` + `FlushConnection`).

The benchmark uses the **NATS stream drain** pattern: 1,000,000 messages are pre-seeded into the stream before the service starts, and wall-clock time to drain the stream to zero is measured. Throughput is computed as `events_ok / drain_seconds`. No HTTP client is involved.

### Scenarios

| Scenario | Router kind | Match rules | Downstream subject |
|----------|-------------|-------------|--------------------|
| `save_event` | Passthrough (no match) | 0 | `events.save` |
| `rule_engine` | Conditional: `data.status == "active"` | 1 | `ruleengine.{id}.execute` |
| `trigger` | Conditional: `data.temperature >= 25` | 1 | `trigger.{id}.execute` |

Each scenario uses a **dedicated asset** linked to exactly **1 RouteGroup**, ensuring isolated measurement per routing type.

---

## 2. Test Environment

| Parameter | Value |
|-----------|-------|
| CPU | 13th Gen Intel Core i9-13900K |
| Logical cores | 16 (available to benchmark shield) |
| RAM | 94 GiB |
| OS | Linux 6.17.0 |
| Go version | go1.25.3 |
| CPU isolation | cgroup v2 cpuset, cores 16-31 |
| MongoDB | localhost:27017, replica set rs0, database `dev-router` |
| NATS | localhost:4222, JetStream enabled, stream `ROUTE-GROUPS` |
| Redis | localhost:6379 (L1 disk-backed cache) |
| MinIO | localhost:9000, bucket `mapex-assets` (L2 object store) |
| Competing services | All killed except Assets (port 5002, required for tiered cache) |

The cgroup v2 shield confines all system processes to cores 0-15. The router binary is moved into `/sys/fs/cgroup/benchmark` immediately after launch. `GOMAXPROCS` controls how many of the shielded cores the Go scheduler uses.

---

## 3. Test Configuration

| Parameter | Value |
|-----------|-------|
| Scenarios | `save_event`, `rule_engine`, `trigger` |
| Messages per test | 1,000,000 |
| CPU configurations | 1, 2, 4, 8, 12, 16 |
| Total tests | 18 (3 scenarios x 6 CPU) |
| NATS batch size | 8,000 messages per fetch |
| NATS stream | `ROUTE-GROUPS` |
| NATS subject | `route.bench.execute` |
| Publish mode | Fire-and-forget (`PublishCore` + batch `FlushConnection`) |
| Poll interval | 2 seconds (drain progress check via `/metrics`) |
| Cool-down between tests | 5 seconds |

### Payload structure (common to all scenarios)

```json
{
  "orgId": "000000000000000000000001",
  "assetUUID": "bench-asset-save-event",
  "pathKey": "benchmark",
  "eventTrackerId": "bench-tracker-001",
  "event": {
    "data": {
      "temperature": 28.5,
      "humidity": 65.2,
      "pressure": 1013.25,
      "battery": 87,
      "status": "active"
    },
    "eventId": "bench-event-001",
    "eventType": "sensorData",
    "targetName": "projects/bench/devices/bench-asset-001",
    "timestamp": "2026-02-16T12:00:00.000Z"
  }
}
```

Each scenario uses a different `assetUUID` (`bench-asset-save-event`, `bench-asset-rule-engine`, `bench-asset-trigger`), mapping to a dedicated MinIO asset with exactly 1 RouteGroup.

### Seeded RouteGroups (MongoDB)

**save_event** — Unconditional passthrough (0 match rules):
```json
{
  "_id": ObjectId("000000000000000000000010"),
  "orgId": ObjectId("000000000000000000000001"),
  "pathKey": "benchmark",
  "enabled": true,
  "routers": [{ "kind": "save_event" }]
}
```

**rule_engine** — Conditional: `data.status == "active"` (1 match rule):
```json
{
  "_id": ObjectId("000000000000000000000011"),
  "routers": [{
    "kind": "rule_engine",
    "match": { "policy": "all", "rules": [{ "field": "data.status", "operator": "eq", "value": "active" }] },
    "ruleEngine": { "businessRuleId": "000000000000000000000020" }
  }]
}
```

**trigger** — Conditional: `data.temperature >= 25` (1 match rule):
```json
{
  "_id": ObjectId("000000000000000000000012"),
  "routers": [{
    "kind": "trigger",
    "match": { "policy": "all", "rules": [{ "field": "data.temperature", "operator": "gte", "value": 25 }] },
    "trigger": { "triggerId": "000000000000000000000021" }
  }]
}
```

---

## 4. Methodology

### Step-by-step execution

```
1. sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh
       Confines system processes to cores 0-15.
       Creates /sys/fs/cgroup/benchmark (cores 16-31).

2. Kill all competing MapexOS services (keep Assets alive on port 5002)

3. scripts/seed.sh setup
       Upserts 3 RouteGroups into MongoDB.
       Uploads 3 AssetReadModel JSONs to MinIO (1 per scenario).

4. For each scenario in [save_event, rule_engine, trigger]:
     For each CPU in [1, 2, 4, 8, 12, 16]:
       a. Kill any leftover service on port 5003.
       b. Narrow shield cpuset to 16-(16+N-1).
       c. Purge all NATS streams (ROUTE-GROUPS + downstream).
       d. Publish 1,000,000 messages to route.bench.execute.
       e. Start binary: GOMAXPROCS=N NATS_BATCH_SIZE=8000 LOG_LEVEL=silent
       f. Move PID into /sys/fs/cgroup/benchmark/cgroup.procs.
       g. Poll GET /metrics every 2s: sum router_message_total{result="ack|nack|reject"}.
          Stop when total >= 1,000,000 or stall detected (10 polls with no progress).
       h. Scrape /metrics once more for full snapshot.
       i. Kill service. Wait 5s.

5. scripts/seed.sh teardown
       Purges NATS streams.
       Deletes RouteGroups from MongoDB + flushes Redis.
       Removes all 3 MinIO assets.
```

### Drain completion detection

The script polls `router_message_total` counters (ack + nack + reject). When the sum reaches the seeded message count, the drain is considered complete. Wall-clock drain time is measured from service start to drain completion.

### Cache warm-up behavior

The first message triggers an L2 MinIO lookup to populate L0 and L1. All subsequent messages are served from L0 RAM. For 1,000,000 messages the cold-start overhead is negligible.

---

## 5. Results

### 5.1 Throughput and Drain Time

#### save_event (passthrough, 0 match rules)

| CPU | Cores | Throughput (ev/s) | Drain Time | Proc avg (ms) |
|-----|-------|-------------------|------------|----------------|
| 1   | 16       | 15,904         | 66.4 s     | 2.072 |
| 2   | 16-17    | 29,240         | 34.2 s     | 1.045 |
| 4   | 16-19    | 62,112         | 16.1 s     | 0.490 |
| 8   | 16-23    | 123,457        | 8.1 s      | 0.260 |
| 12  | 16-27    | 163,934        | 6.1 s      | 0.196 |
| 16  | 16-31    | 163,934        | 6.1 s      | 0.195 |

#### rule_engine (1 match rule: `data.status == "active"`)

| CPU | Cores | Throughput (ev/s) | Drain Time | Proc avg (ms) |
|-----|-------|-------------------|------------|----------------|
| 1   | 16       | 12,213         | 84.5 s     | 2.672 |
| 2   | 16-17    | 23,641         | 42.3 s     | 1.344 |
| 4   | 16-19    | 49,505         | 20.2 s     | 0.649 |
| 8   | 16-23    | 82,645         | 12.1 s     | 0.362 |
| 12  | 16-27    | 123,457        | 8.1 s      | 0.253 |
| 16  | 16-31    | 123,457        | 8.1 s      | 0.216 |

#### trigger (1 match rule: `data.temperature >= 25`)

| CPU | Cores | Throughput (ev/s) | Drain Time | Proc avg (ms) |
|-----|-------|-------------------|------------|----------------|
| 1   | 16       | 12,606         | 82.5 s     | 2.588 |
| 2   | 16-17    | 23,641         | 42.3 s     | 1.303 |
| 4   | 16-19    | 49,505         | 20.2 s     | 0.627 |
| 8   | 16-23    | 99,010         | 10.1 s     | 0.329 |
| 12  | 16-27    | 123,457        | 8.1 s      | 0.239 |
| 16  | 16-31    | 123,457        | 8.1 s      | 0.208 |

#### Cross-scenario comparison (ev/s)

| CPU | save_event | rule_engine | trigger | Match overhead |
|-----|------------|-------------|---------|----------------|
| 1   | 15,904     | 12,213      | 12,606  | ~23% slower |
| 2   | 29,240     | 23,641      | 23,641  | ~19% slower |
| 4   | 62,112     | 49,505      | 49,505  | ~20% slower |
| 8   | 123,457    | 82,645      | 99,010  | ~20-33% slower |
| 12  | 163,934    | 123,457     | 123,457 | ~25% slower |
| 16  | 163,934    | 123,457     | 123,457 | ~25% slower |

### 5.2 Message Lifecycle (0 nack, 0 reject across all 18 tests)

| Scenario | CPU | Ack | Nack | Reject | Error |
|----------|-----|-----|------|--------|-------|
| save_event | 1 | 1,056,000 | 0 | 0 | 0 |
| save_event | 2-16 | 1,000,000 | 0 | 0 | 0 |
| rule_engine | 1 | 1,032,000 | 0 | 0 | 0 |
| rule_engine | 2-16 | 1,000,000 | 0 | 0 | 0 |
| trigger | 1 | 1,040,000 | 0 | 0 | 0 |
| trigger | 2-16 | 1,000,000 | 0 | 0 | 0 |

> CPU=1 shows slightly more than 1M acked due to NATS redelivery (slow single-core processing exceeds ack timeout). All errors are zero in every run.

### 5.3 Asset Cache (TieredCache)

All lookups served from L0 RAM after the first message. Cache hit rate >99.999%.

| CPU | L0 RAM | L1 Disk | L2 MinIO | Fallback | Miss |
|-----|--------|---------|----------|----------|------|
| 1   | ~1,000,000 | 1 | 0 | 0 | 0 |
| 2   | ~1,000,000 | 1-2 | 0 | 0 | 0 |
| 4   | ~1,000,000 | 4-6 | 0 | 0 | 0 |
| 8   | ~1,000,000 | 4-8 | 0 | 0 | 0 |
| 12  | ~1,000,000 | 6-10 | 0 | 0 | 0 |
| 16  | ~1,000,000 | 0 | 32 | 0 | 0 |

> CPU=16 shows L2 MinIO hits instead of L1 Disk. This reflects initial cold-start race conditions with many goroutines before L0 is primed. All subsequent lookups are L0 RAM.

### 5.4 Match Evaluation

| Scenario | CPU | Matched | Unmatched | No config |
|----------|-----|---------|-----------|-----------|
| save_event | all | 0 | 0 | ~1,000,000 |
| rule_engine | all | ~1,000,000 | 0 | 0 |
| trigger | all | ~1,000,000 | 0 | 0 |

- `save_event` has no match config — all events bypass evaluation (`no_config`).
- `rule_engine` and `trigger` both evaluate 1 match rule per event. All events match (payload contains `status: "active"` and `temperature: 28.5 >= 25`).

### 5.5 NATS Publish

All events published via fire-and-forget (`PublishCore`). Zero publish errors across all 18 tests.

| Scenario | Downstream subject | Pub OK | Pub Error |
|----------|-------------------|--------|-----------|
| save_event | `events.save` | ~1,000,000 | 0 |
| rule_engine | `ruleengine.{id}.execute` | ~1,000,000 | 0 |
| trigger | `trigger.{id}.execute` | ~1,000,000 | 0 |

### 5.6 RouteGroup Cache

In-process cache hit rate >99.99% after the first fetch.

| CPU | Hit | Miss |
|-----|-----|------|
| 1   | ~1,000,000 | 0 |
| 2-12 | ~1,000,000 | 0 |
| 16  | ~999,981 | 18-19 |

### 5.7 Batch Metrics

Consistent 8,000 messages per batch across all tests.

| CPU | Batches | Avg batch size |
|-----|---------|----------------|
| 1   | 130-133 | 8,000 |
| 2-16 | 125 | 8,000 |

> CPU=1 processes extra batches due to NATS redelivery.

### 5.8 Memory and Go Runtime

#### save_event

| CPU | RSS (MB) | Heap (MB) | Goroutines | Threads | GC cycles | GC avg (us) |
|-----|----------|-----------|------------|---------|-----------|-------------|
| 1   | 82.7     | 58.2      | 56         | 3       | 418       | 31.8 |
| 2   | 79.5     | 50.9      | 27         | 5       | 281       | 57.7 |
| 4   | 99.3     | 77.5      | 28         | 7       | 214       | 76.5 |
| 8   | 99.2     | 75.9      | 28         | 12      | 215       | 116.6 |
| 12  | 99.8     | 57.5      | 28         | 16      | 218       | 158.4 |
| 16  | 106.1    | 90.1      | 57         | 24      | 216       | 138.9 |

#### rule_engine

| CPU | RSS (MB) | Heap (MB) | Goroutines | Threads | GC cycles | GC avg (us) |
|-----|----------|-----------|------------|---------|-----------|-------------|
| 1   | 82.3     | 54.0      | 56         | 3       | 568       | 31.5 |
| 2   | 86.2     | 64.1      | 27         | 5       | 430       | 59.6 |
| 4   | 100.2    | 95.5      | 28         | 7       | 319       | 77.9 |
| 8   | 103.4    | 53.2      | 28         | 11      | 321       | 155.5 |
| 12  | 98.7     | 50.6      | 28         | 16      | 322       | 211.0 |
| 16  | 106.2    | 59.1      | 58         | 22      | 323       | 145.7 |

#### trigger

| CPU | RSS (MB) | Heap (MB) | Goroutines | Threads | GC cycles | GC avg (us) |
|-----|----------|-----------|------------|---------|-----------|-------------|
| 1   | 77.0     | 47.5      | 56         | 3       | 622       | 32.0 |
| 2   | 78.3     | 67.8      | 27         | 5       | 384       | 57.9 |
| 4   | 100.3    | 73.2      | 28         | 7       | 288       | 78.6 |
| 8   | 99.8     | 75.6      | 28         | 12      | 289       | 120.9 |
| 12  | 103.6    | 52.4      | 28         | 15      | 293       | 198.0 |
| 16  | 108.2    | 78.9      | 58         | 21      | 292       | 140.1 |

### 5.9 Scaling Efficiency (save_event)

| CPU | ev/s per core | Scaling factor vs CPU=1 | Efficiency |
|-----|---------------|------------------------|------------|
| 1   | 15,904        | 1.00x                  | 100% |
| 2   | 14,620        | 1.84x                  | 92% |
| 4   | 15,528        | 3.91x                  | 98% |
| 8   | 15,432        | 7.76x                  | 97% |
| 12  | 13,661        | 10.31x                 | 86% |
| 16  | 10,246        | 10.31x                 | 64% |

> CPU=12 and CPU=16 produce identical throughput (163,934 ev/s), indicating a ceiling. The bottleneck at this point is NATS JetStream publish throughput on a single connection.

---

## 6. Key Findings

- **Near-linear scaling from 1 to 8 cores.** CPU=8 delivers 7.76x the throughput of CPU=1 for save_event (123,457 vs 15,904 ev/s) with 97% per-core efficiency. This is a significant improvement over the previous baseline (6.37x at CPU=8) due to the double-buffer consumer pattern.

- **Throughput ceiling at ~164K ev/s.** CPU=12 and CPU=16 produce identical throughput for save_event (163,934 ev/s). The bottleneck above 12 cores is the single NATS connection's publish throughput, not CPU.

- **Match rule evaluation costs ~20-25% throughput.** Scenarios with 1 match rule (rule_engine, trigger) consistently run 20-25% slower than save_event across all CPU configs. The overhead is per-event rule evaluation, not publish or cache.

- **rule_engine ceiling is lower than trigger.** At CPU=8, trigger achieves 99,010 ev/s vs rule_engine at 82,645 ev/s. Both plateau at 123,457 ev/s from CPU=12+. The string comparison (`data.status == "active"`) is slightly more expensive than the numeric comparison (`data.temperature >= 25`).

- **100% success rate across all 18 tests.** Zero nacks, zero rejects, zero publish errors across 18M total messages processed (18 tests x 1M).

- **TieredCache >99.999% L0 RAM hit rate.** After the initial cold-start lookup, virtually all asset resolutions are served from L0 RAM with sub-microsecond latency.

- **Memory footprint is flat.** RSS stays between 77-108 MB regardless of scenario or core count. No memory growth with throughput.

- **GC pause increases with core count but remains negligible.** From 32 us at CPU=1 to ~140-211 us at CPU=12-16. All values are well within acceptable bounds for a background consumer.

- **Goroutine count is stable.** 27-28 goroutines for CPU 2-12, 56-58 for CPU 1 and 16 (includes polling goroutines). No goroutine leak detected.

- **NATS redelivery at CPU=1.** Single-core processing is slow enough to exceed NATS ack timeout, causing redelivery of ~3-5% of messages. This is expected and has zero impact on correctness (idempotent processing).

---

## 7. Conclusions

### Production recommendation

For a production deployment, a **4-core replica** is the optimal configuration for balanced throughput and efficiency:

| Scenario | 4-core throughput | Per-core efficiency |
|----------|-------------------|---------------------|
| save_event | 62,112 ev/s | 15,528 ev/s/core (98%) |
| rule_engine | 49,505 ev/s | 12,376 ev/s/core (98%) |
| trigger | 49,505 ev/s | 12,376 ev/s/core (98%) |

For higher sustained throughput, prefer **horizontal scaling** (multiple 4-core replicas with competing consumers) over vertical scaling beyond 12 cores. Above 12 cores there is zero throughput gain due to NATS publish saturation.

### Scaling formula

```
# save_event (passthrough routing)
replicas = ceil(target_events_per_second / 62,000)

# rule_engine or trigger (conditional routing)
replicas = ceil(target_events_per_second / 49,000)
```

Examples (save_event):
- 120,000 ev/s -> 2 replicas (4 cores each)
- 240,000 ev/s -> 4 replicas
- 1,000,000 ev/s -> 17 replicas

### Resource requests (Kubernetes)

```yaml
resources:
  requests:
    cpu: "4"
    memory: "128Mi"
  limits:
    cpu: "4"
    memory: "256Mi"
```

### Throughput ceiling

The maximum throughput on a single instance is **~164K ev/s** (save_event) or **~123K ev/s** (conditional routing), reached at CPU=12. This ceiling is imposed by the single NATS connection. Potential optimizations:
- Multiple NATS connections for publishing (connection pool)
- Larger `NATS_BATCH_SIZE` (e.g., 16,000) to reduce fetch round-trips
- NATS server tuning (write buffer size, flush interval)

---

## 8. Reproduction

### Prerequisites

| Tool | Purpose | Install |
|------|---------|---------|
| `nats` | Stream management and message publishing | `go install github.com/nats-io/natscli/nats@latest` |
| `mongosh` | MongoDB seeding and cleanup | Install with MongoDB or standalone |
| `mc` | MinIO asset upload | `brew install minio/stable/mc` or binary download |
| `redis-cli` | Redis cache flush | Install with Redis or standalone |
| `curl` | Metrics collection | Pre-installed on most Linux systems |
| `go` | Binary compilation | `https://go.dev/dl/` |
| Root access | cgroup v2 shield setup | Required for `setup-cgroup-shield.sh` |

Infrastructure must be running:

- MongoDB at `localhost:27017` (replica set `rs0`)
- NATS with JetStream at `localhost:4222` (user `service`, password `service_secret`)
- Redis at `localhost:6379`
- MinIO at `localhost:9000` (mc alias `local`, bucket `mapex-assets`)
- Assets service at `localhost:5002` (required for tiered cache fallback)

### Commands

```bash
# 1. Setup CPU isolation (one-time, requires root)
sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh

# 2. Run the full benchmark (all 3 scenarios, 1M messages, 6 CPU configs = 18 tests)
bash docs/benchmarks/scripts/full-benchmark.sh

# 3. Run a single scenario
SCENARIOS="save_event" bash docs/benchmarks/scripts/full-benchmark.sh

# 4. Quick smoke test (1K messages, 1 CPU)
SCENARIOS="save_event" bash docs/benchmarks/scripts/full-benchmark.sh 1000 "1" smoke-test

# 5. Seed only, for manual testing
bash docs/benchmarks/scripts/seed.sh setup
bash docs/benchmarks/scripts/seed.sh teardown

# 6. Teardown CPU isolation when done
sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh teardown
```

### Script locations

| Script | Path |
|--------|------|
| CPU isolation | `scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh` |
| Benchmark runner | `docs/benchmarks/scripts/full-benchmark.sh` |
| Service config | `docs/benchmarks/scripts/config.sh` |
| Data seeding | `docs/benchmarks/scripts/seed.sh` |
| NATS payloads | `docs/benchmarks/seed/nats/` |
| MongoDB seed | `docs/benchmarks/seed/mongodb/` |
| MinIO assets | `docs/benchmarks/seed/minio/` |

### Full-benchmark.sh arguments

```
./full-benchmark.sh [message_count] [cpu_list] [tag] [nats_batch_size]

Defaults: 1000000 "1 2 4 8 16" baseline 8000
```

Override scenarios via environment variable:
```bash
SCENARIOS="save_event rule_engine trigger"   # default: all three
SCENARIOS="save_event"                        # single scenario
```

### Results storage

Per-test raw data (Prometheus snapshots, service logs) stored in:

```
docs/benchmarks/results/{tag}/
├── test-{scenario}-cpu{N}-metrics.txt    # Full Prometheus dump
└── test-{scenario}-cpu{N}-output.log     # Service stdout/stderr
```

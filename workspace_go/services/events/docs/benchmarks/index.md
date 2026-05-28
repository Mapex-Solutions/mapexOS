# Benchmarks

## Overview

The events service is the **only MapexOS service that writes to ClickHouse**. It consumes NATS JetStream messages from 7 streams and performs bulk inserts into ClickHouse tables. Each consumer processes a different event type with varying complexity — from lightweight DLQ passthrough to heavy business rule evaluation with 5x JSON marshaling.

The benchmark uses the **NATS stream drain pattern**: a fixed number of messages are pre-seeded into each stream before the service starts. The service is then started under cgroup v2 CPU isolation and the time to drain all messages to zero pending is measured via Prometheus metrics polling. This isolates pure processing throughput from network and client latency.

### Consumers

| Consumer | NATS Stream | ClickHouse Table | Complexity |
|---|---|---|---|
| `save_raw_event` | `EVENTS-RAW` | `events_raw` | Light — unmarshal + validate + map |
| `save_jsexec_event` | `EVENTS-JSEXEC` | `events_jsexecutor` | Medium — flatten DTO + marshal payload |
| `save_router_event` | `EVENTS-ROUTER` | `events_router` | Medium — counter calc + marshal routers array |
| `save_businessrule_event` | `EVENTS-BUSINESSRULE` | `events_businessrule` | **Heavy** — 5x JSON marshal per message |
| `save_trigger_event` | `EVENTS-TRIGGER` | `events_trigger` | Light — direct mapping |
| `save_event` | `EVENTS` | `events` | Medium — TieredCache template lookup + EVA field resolution |
| `save_dlq_event` | `MAPEXOS-DLQ` | `events_dlq` | Lightest — passthrough (never nacks to avoid loop) |

### Processing Architecture

All consumers use the **Three-Phase** pattern with parallel worker pools:

```
Phase 1 (Parallel): Worker pool (NumCPU*2) parses/validates/maps each message
Phase 2 (Bulk I/O):  Collect valid entities → single ClickHouse INSERT
Phase 3 (ACK/Nack):  Reject invalid, ACK valid (or Nack all on insert failure)
```

---

## Test Environment

| Property | Value |
|---|---|
| CPU | Intel Core i9-13900K (8 P-cores HT + 16 E-cores = 32 logical) |
| RAM | 64 GB DDR5 |
| OS | Ubuntu 24.04 LTS, kernel 6.17 |
| Go | 1.25.x |
| Fiber | v2.52.9 |
| Isolation | cgroup v2 cpuset — benchmark cgroup pinned to cores 16-31 |
| ClickHouse | 24.x, local (native TCP port 9440) |
| NATS | 2.10.x with JetStream, local |
| MongoDB | 7.x, replicaSet=rs0, local |
| Redis | 7.x, local |
| MinIO | RELEASE.2024-x, local (TieredCache L2 for templates) |

---

## Test Configuration

| Parameter | Value |
|---|---|
| Messages per test | 1,000,000 |
| CPU configurations | 1, 2, 4, 8, 16 cores |
| NATS batch size | 10,000 (`NATS_BATCH_SIZE`) |
| Scenarios | All 7 consumers (default) |
| Service port | 5004 |
| ClickHouse pool | MaxOpenConns=20, MaxIdleConns=10 |
| ClickHouse compression | LZ4 |
| Protocol | Native TCP (not HTTP) |
| LOG_LEVEL | `warn` (suppress I/O noise) |

### Payload Sizes

| Scenario | Payload | Size |
|---|---|---|
| `save_raw_event` | IoT sensor data (temperature, humidity, pressure, battery, status) | ~280 bytes |
| `save_jsexec_event` | JS executor debug log with flat DTO | ~350 bytes |
| `save_router_event` | Router execution with routers array | ~400 bytes |
| `save_businessrule_event` | Business rule evaluation with conditions + actions | ~500 bytes |
| `save_trigger_event` | Trigger execution log | ~300 bytes |
| `save_event` | EventStore with 7 dynamic fields (EVA resolution via template) | ~350 bytes |
| `save_dlq_event` | Dead letter envelope with original headers | ~400 bytes |

---

## Methodology

The benchmark follows the [general benchmark standard](../../../../../documentations/architecture/microservices/general/benchmarks/benchmark_standard_pattern.md) and uses common modules from `scripts/benchmarks/common/`.

### Step-by-step

1. **CPU isolation** — run `sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh` once per session. Pins all OS processes to cores 0-15, creates `/sys/fs/cgroup/benchmark` on cores 16-31.

2. **Seed** — `full-benchmark.sh` calls `./seed.sh setup` automatically (idempotent):
   - ClickHouse: creates all 7 event tables (IF NOT EXISTS)
   - MinIO: uploads benchmark template for TieredCache L2 (`save_event` EVA resolution)
   - MongoDB: upserts retention policy for benchmark org
   - Redis: flushes app cache
   - NATS: purges all 7 event streams

3. **Build** — compiles a fresh binary (`events_bench`) with `GOWORK=off CGO_ENABLED=0 go build`.

4. **Per-scenario x CPU loop** — for each scenario in `$SCENARIOS` × each CPU in `$CPU_LIST`:
   - Write `cpuset.cpus` on the benchmark cgroup to pin N cores.
   - Purge the scenario's NATS stream + delete consumer.
   - Seed 1M messages into the stream via `nats pub --count`.
   - Start service with `GOMAXPROCS=N LOG_LEVEL=warn NATS_BATCH_SIZE=10000`.
   - Move process into cgroup. Wait up to 60s for `/metrics` readiness.
   - **Drain poll**: every 2s, sum `events_message_total{result="ack|nack|reject"}`. Stop when total >= N (or stall detected after 10 polls with no progress).
   - Dump full Prometheus metrics to `test-{scenario}-cpu{N}-metrics.txt`.
   - Stop service (SIGTERM), 5s cool-down.

5. **Teardown** — `./seed.sh teardown` (purge NATS, remove MinIO template, delete MongoDB seed, flush Redis).

### Drain Detection

The benchmark polls the service's own `/metrics` endpoint (not `nats stream info`). It sums:

```
events_message_total{consumer="<label>", result="ack"}
+ events_message_total{consumer="<label>", result="nack"}
+ events_message_total{consumer="<label>", result="reject"}
```

When this sum >= seeded message count, the drain is complete. This is more accurate than broker-side pending counts because it tracks service-side processing.

---

## Results

> **Run**: 2026-02-26 | **Messages**: 1,000,000 per scenario × CPU | **Batch**: 10,000 | **Reliability**: 100% (zero errors, zero nacks, zero rejects across all 35 tests)

### save_raw_event — Raw events from HTTP/MQTT gateways

| CPU | Throughput (ev/s) | Proc Avg (ms) | CH Insert (ms) | CH Batch | RSS (MB) | Heap (MB) | Goroutines | GC |
|-----|-------------------|----------------|-----------------|----------|----------|-----------|------------|----|
| 1 | 15,397 | 649.5 | 38.8 | 10,000 | 103.9 | 52.6 | 85 | 384 |
| 2 | 27,270 | 366.7 | 34.0 | 10,000 | 115.5 | 87.4 | 48 | 305 |
| **4** | **32,657** | **306.2** | 31.2 | 10,000 | 123.8 | 83.6 | 48 | 234 |
| 8 | 31,944 | 313.0 | 31.5 | 10,000 | 133.4 | 75.8 | 49 | 243 |
| 16 | 31,609 | 316.4 | 31.2 | 10,000 | 148.8 | 83.9 | 48 | 244 |

### save_jsexec_event — JS executor debug logs

| CPU | Throughput (ev/s) | Proc Avg (ms) | CH Insert (ms) | CH Batch | RSS (MB) | Heap (MB) | Goroutines | GC |
|-----|-------------------|----------------|-----------------|----------|----------|-----------|------------|----|
| 1 | 15,582 | 641.8 | 19.8 | 10,000 | 102.6 | 69.0 | 83 | 390 |
| 2 | 28,139 | 355.4 | 20.1 | 10,000 | 110.9 | 67.8 | 48 | 306 |
| **4** | **33,566** | **297.9** | 16.8 | 10,000 | 125.9 | 114.4 | 48 | 231 |
| 8 | 33,242 | 300.8 | 16.8 | 10,000 | 131.3 | 86.5 | 48 | 235 |
| 16 | 32,409 | 308.6 | 17.3 | 10,000 | 127.8 | 77.5 | 48 | 234 |

### save_router_event — Router execution history

| CPU | Throughput (ev/s) | Proc Avg (ms) | CH Insert (ms) | CH Batch | RSS (MB) | Heap (MB) | Goroutines | GC |
|-----|-------------------|----------------|-----------------|----------|----------|-----------|------------|----|
| 1 | 16,784 | 595.8 | 39.6 | 10,000 | 207.6 | 131.6 | 85 | 281 |
| 2 | 29,590 | 338.0 | 32.2 | 10,000 | 207.9 | 99.2 | 48 | 240 |
| **4** | **32,513** | **307.6** | 30.8 | 10,000 | 180.1 | 118.1 | 49 | 205 |
| 8 | 32,136 | 311.2 | 32.5 | 10,000 | 120.4 | 59.4 | 48 | 200 |
| 16 | 31,408 | 318.4 | 32.6 | 10,000 | 189.8 | 146.6 | 49 | 201 |

### save_businessrule_event — Business rule evaluation (heaviest)

| CPU | Throughput (ev/s) | Proc Avg (ms) | CH Insert (ms) | CH Batch | RSS (MB) | Heap (MB) | Goroutines | GC |
|-----|-------------------|----------------|-----------------|----------|----------|-----------|------------|----|
| 1 | 11,410 | 876.4 | 74.3 | 10,000 | 275.3 | 235.1 | 89 | 378 |
| 2 | 20,567 | 486.2 | 59.0 | 9,940 | 247.6 | 101.3 | 48 | 348 |
| 4 | 29,280 | 341.5 | 55.0 | 10,000 | 284.9 | 190.2 | 49 | 288 |
| **8** | **29,661** | **337.1** | 53.9 | 10,000 | 271.0 | 162.8 | 49 | 287 |
| 16 | 28,733 | 348.0 | 54.8 | 10,000 | 268.6 | 169.4 | 48 | 287 |

### save_trigger_event — Trigger execution logs

| CPU | Throughput (ev/s) | Proc Avg (ms) | CH Insert (ms) | CH Batch | RSS (MB) | Heap (MB) | Goroutines | GC |
|-----|-------------------|----------------|-----------------|----------|----------|-----------|------------|----|
| 1 | 20,427 | 489.5 | 26.8 | 9,931 | 100.8 | 44.5 | 48 | 296 |
| **2** | **34,887** | **286.6** | 25.2 | 10,000 | 139.3 | 96.9 | 48 | 242 |
| 4 | 33,259 | 300.7 | 24.7 | 10,000 | 153.7 | 122.6 | 48 | 205 |
| 8 | 33,107 | 302.1 | 25.3 | 10,000 | 163.1 | 101.8 | 48 | 203 |
| 16 | 32,318 | 309.4 | 25.8 | 10,000 | 163.3 | 106.9 | 48 | 203 |

### save_event — EventStore with EVA field resolution

| CPU | Throughput (ev/s) | Proc Avg (ms) | CH Insert (ms) | CH Batch | RSS (MB) | Heap (MB) | Goroutines | GC |
|-----|-------------------|----------------|-----------------|----------|----------|-----------|------------|----|
| 1 | 15,226 | 656.8 | 27.5 | 10,000 | 123.1 | 62.0 | 86 | 411 |
| 2 | 27,524 | 363.3 | 24.4 | 10,000 | 137.7 | 113.4 | 48 | 314 |
| **4** | **33,621** | **297.4** | 23.1 | 10,000 | 141.8 | 78.6 | 48 | 256 |
| 8 | 33,520 | 298.3 | 23.6 | 10,000 | 101.9 | 52.7 | 49 | 259 |
| 16 | 32,687 | 305.9 | 23.6 | 10,000 | 158.3 | 115.4 | 49 | 257 |

> **Template cache**: 1M hits, 0 misses, 0 errors. **EVA fields**: 7M mapped (7 fields × 1M events). TieredCache L0 (RAM) served all lookups after initial L2 (MinIO) load.

### save_dlq_event — Dead Letter Queue (lightest)

| CPU | Throughput (ev/s) | Proc Avg (ms) | CH Insert (ms) | CH Batch | RSS (MB) | Heap (MB) | Goroutines | GC |
|-----|-------------------|----------------|-----------------|----------|----------|-----------|------------|----|
| 1 | 82,091 | 121.8 | 33.4 | 10,000 | 137.1 | 59.8 | 50 | 120 |
| 2 | 139,617 | 71.6 | 27.2 | 10,000 | 120.8 | 53.5 | 50 | 134 |
| 4 | 202,511 | 49.4 | 24.0 | 10,000 | 123.7 | 21.9 | 49 | 162 |
| 8 | 233,807 | 42.8 | 26.6 | 10,000 | 122.7 | 53.4 | 50 | 177 |
| **16** | **263,393** | **38.0** | 25.8 | 10,000 | 67.5 | 23.8 | 50 | 165 |

> DLQ is special: it **always ACKs** (never nacks/rejects to avoid infinite loop). Minimal processing — just passthrough mapping.

### Stream Drain Time — 1M messages

Total time to consume and persist all 1,000,000 messages per stream. **Bold** = best time per stream.

| Stream | CPU=1 | CPU=2 | CPU=4 | CPU=8 | CPU=16 |
|---|---|---|---|---|---|
| `save_raw_event` | 66.9s | 36.7s | **30.6s** | 31.3s | 31.6s |
| `save_jsexec_event` | 66.7s | 35.5s | **29.8s** | 30.1s | 30.9s |
| `save_router_event` | 63.2s | 33.8s | **30.8s** | 31.1s | 31.8s |
| `save_businessrule_event` | 91.1s | 49.1s | 34.2s | **33.7s** | 34.8s |
| `save_trigger_event` | 49.9s | **28.7s** | 30.1s | 30.2s | 30.9s |
| `save_event` | 69.0s | 36.3s | **29.7s** | 29.8s | 30.6s |
| `save_dlq_event` | 12.2s | 7.2s | 4.9s | 4.3s | **3.8s** |

### Performance Summary

| Scenario | Best CPU | Peak (ev/s) | Proc Avg (ms) | CH Insert (ms) | Bottleneck |
|---|---|---|---|---|---|
| **save_dlq_event** | 16 | **263,393** | 38.0 | 25.8 | ClickHouse insert (processing is trivial) |
| **save_trigger_event** | 2 | **34,887** | 286.6 | 25.2 | CPU (direct mapping, light) |
| **save_event** | 4 | **33,621** | 297.4 | 23.1 | CPU (template cache + EVA mapping) |
| **save_jsexec_event** | 4 | **33,566** | 297.9 | 16.8 | CPU (flatten DTO + marshal) |
| **save_raw_event** | 4 | **32,657** | 306.2 | 31.2 | CPU (unmarshal + validate + map) |
| **save_router_event** | 4 | **32,513** | 307.6 | 30.8 | CPU (counter calc + marshal array) |
| **save_businessrule_event** | 8 | **29,661** | 337.1 | 53.9 | CPU + ClickHouse (5x JSON marshal + large rows) |

---

## Key Findings

### Processing Performance

1. **Sweet spot at 4 CPUs** — Five of seven consumers peak at 4 CPUs (~30-34K ev/s). The Three-Phase worker pool (NumCPU×2 workers) saturates at 8 goroutines. Beyond 4 cores, context switching overhead negates parallelism gains.

2. **BusinessRule is the heaviest consumer** — 876ms/batch at CPU=1 (5x JSON marshal per message), needs 8 CPUs to peak at 29.6K ev/s. It also has the largest CH insert latency (54-74ms) due to bigger row payloads.

3. **DLQ scales linearly with CPUs** — 82K→263K ev/s from 1→16 CPUs. Minimal per-message processing makes it purely ClickHouse-insert-bound at lower CPU counts.

4. **Trigger peaks early at 2 CPUs** — The simplest standard consumer (direct mapping, no marshaling). Two worker goroutines are enough; more adds overhead.

5. **All standard consumers converge at ~30-34K ev/s** — Despite varying message complexity (raw vs router vs jsexec), the ClickHouse insert becomes the common bottleneck once parsing is parallelized.

### ClickHouse Insert Patterns

6. **Insert latency varies by payload size** — JsExec has the smallest rows (17ms avg), BusinessRule the largest (54ms avg). Raw and Router are mid-range (31-40ms). Batch size is consistently 10,000 across all scenarios.

7. **100 inserts per 1M messages** — With batch size 10,000, exactly 100 ClickHouse INSERT calls per test. Zero insert failures across all 35 tests.

8. **CH insert is ~10% of batch time** — For standard consumers at 4 CPUs, insert takes ~20-30ms out of ~300ms total. The remaining ~270ms is Phase 1 (parallel parse/validate/map). This confirms parallelization was the right optimization.

### Template Cache (save_event)

9. **TieredCache delivers 100% hit rate** — 1M template cache hits, 0 misses, 0 errors. After the first L2 (MinIO) load, all subsequent lookups are served from L0 (in-process RAM, ~50µs). **Without template cache, save_event drops from 33.6K/s to 5.4K/s** (measured in previous runs with un-seeded MinIO).

10. **7M EVA fields mapped** — 7 dynamic fields × 1M events, all resolved via cached template. Zero skipped fields.

### Memory and Runtime

11. **RSS stable across message volume** — RSS ranges from 68-285MB depending on consumer complexity (not message count). BusinessRule and Router allocate the most (marshal buffers for large JSON arrays). DLQ is the lightest.

12. **Goroutines drop from ~85 to ~48 at 2+ CPUs** — At CPU=1, the Go runtime creates more goroutines to compensate for the single thread. At 2+ CPUs, goroutine count stabilizes at ~48-50.

13. **GC pressure decreases with more CPUs** — CPU=1 runs 280-411 GC cycles vs 200-290 at 4+ CPUs. Faster batch processing means less heap accumulation between GC pauses.

### Scalability for Enterprise IoT

14. **Standard consumers handle 30K+ events/s at 4 CPUs** — For a fleet of 100K IoT devices reporting every 30 seconds (3,333 ev/s sustained), a single instance has 10x headroom.

15. **All 7 consumers run in a single process** — The 7 consumers share the same connection pool (20 open / 10 idle). In production, all consumers process concurrently, sharing CPU and CH bandwidth. The benchmark isolates each consumer to measure individual throughput.

---

## Conclusions

### Production Recommendations

1. **Events service: 4 CPUs** — Optimal for 6 of 7 consumers. BusinessRule benefits from 8 CPUs if it's the dominant workload.

2. **NATS batch size: 10,000** — Matches ClickHouse optimal insert range (10K-100K rows). Larger batches (50K) would reduce insert round-trips but increase per-batch latency and memory.

3. **ClickHouse pool: 20/10** — 7 concurrent consumers + queries need headroom beyond the old default of 10/5.

4. **Template cache: always seed** — `save_event` depends on TieredCache for EVA field resolution. Without a cached template, throughput drops 6x. Ensure MinIO has templates before starting the service.

5. **DLQ capacity** — At 263K ev/s, the DLQ consumer can absorb massive failure bursts without backpressure. No tuning needed.

### Capacity Planning

| Fleet Size | Events/s (sustained) | Events/s (burst 10x) | CPUs needed | Instances |
|---|---|---|---|---|
| 10K devices @ 30s | 333 | 3,330 | 1 | 1 |
| 100K devices @ 30s | 3,333 | 33,330 | 4 | 1 |
| 500K devices @ 30s | 16,667 | 166,670 | 4 | 5-6 |
| 1M devices @ 30s | 33,333 | 333,330 | 4 | 10-12 |

> Based on `save_raw_event` throughput (~32K ev/s at 4 CPUs) as the representative hot path. `save_event` (EventStore) performs identically with template cache warm.

---

## Reproduction

### Prerequisites

```bash
# Tools
go install github.com/nats-io/natscli/nats@latest
sudo apt install mongosh redis-tools jq bc

# MinIO Client (required for template cache seeding)
# https://min.io/docs/minio/linux/reference/minio-mc.html
mc alias set local http://localhost:9000 mapex_admin mapex_admin_secret_change_me
```

### Infrastructure

- **ClickHouse** — local, native TCP port 9440 (docker maps 9440→9000)
- **NATS** — `nats://localhost:4222` with JetStream, user `service` / `service_secret`
- **MongoDB** — `mongodb://localhost:27017/?replicaSet=rs0`
- **Redis** — `localhost:6379`
- **MinIO** — `localhost:9000` with bucket `mapex-templates`

### Steps

```bash
# 1. CPU isolation (one-time per session, requires root)
sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh

# 2. Full benchmark — all 7 scenarios, 1M messages, CPU 1 2 4 8 16
cd services/events
bash docs/benchmarks/scripts/full-benchmark.sh

# 3. Custom message count and CPU list
bash docs/benchmarks/scripts/full-benchmark.sh 1000000 "16 12 8 4 2 1"

# 4. Smoke test (10K messages, 2 CPUs)
bash docs/benchmarks/scripts/full-benchmark.sh 10000 "1 2"

# 5. Run specific scenarios via SCENARIOS env var
SCENARIOS="save_raw_event" bash docs/benchmarks/scripts/full-benchmark.sh 10000 "1 2"
SCENARIOS="save_businessrule_event save_event" bash docs/benchmarks/scripts/full-benchmark.sh 100000 "4 8"

# 6. Help — shows all env vars and usage
bash docs/benchmarks/scripts/full-benchmark.sh --help

# 7. Teardown cgroup shield
sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh teardown

# 8. Manual seed teardown
bash docs/benchmarks/scripts/seed.sh teardown
```

### Script Locations

| Script | Purpose |
|---|---|
| `scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh` | CPU isolation via cgroup v2 (run once as root) |
| `docs/benchmarks/scripts/config.sh` | Service constants + common module bootstrap |
| `docs/benchmarks/scripts/seed.sh` | ClickHouse tables + MongoDB + MinIO template + Redis + NATS |
| `docs/benchmarks/scripts/full-benchmark.sh` | Main benchmark runner (calls seed.sh automatically) |
| `docs/benchmarks/seed/payloads/*.json` | 7 payload files (one per scenario) |
| `docs/benchmarks/seed/payloads/bench-template.json` | Template for TieredCache (save_event EVA resolution) |

### Output Artifacts

```
docs/benchmarks/results/<timestamp>/
├── test-<scenario>-cpu<N>-metrics.txt   # Full Prometheus metrics dump
├── test-<scenario>-cpu<N>-output.log    # Service stdout/stderr
└── ...                                  # 70 files for full run (7 scenarios × 5 CPUs × 2 files)
```

### Environment Variables

```bash
# Scenarios (default: all 7)
SCENARIOS="save_raw_event save_event"

# Infrastructure
MONGO_URI="mongodb://localhost:27017/?replicaSet=rs0"
MONGO_DB="events"
REDIS_HOST="localhost"
REDIS_PORT="6379"
REDIS_APP_DB="0"
NATS_USER="service"
NATS_PASS="service_secret"
CLICKHOUSE_HOST="localhost"
CLICKHOUSE_PORT="9440"
CLICKHOUSE_DATABASE="mapexos"
CLICKHOUSE_USERNAME="mapexos_user"
CLICKHOUSE_PASSWORD="mapexos_password"
MINIO_ALIAS="local"
MINIO_TEMPLATES_BUCKET="mapex-templates"

# Runtime
GO_ENV="dev"
SUDO_PASS="<password>"
```

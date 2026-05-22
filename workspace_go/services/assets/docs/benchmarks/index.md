# Benchmarks

> **Status — MQTT control-plane scenarios not yet benchmarked.** The
> earlier NATS Auth Callout scenarios (`auth-cache-hit` /
> `auth-cache-miss`) and the `auth-bench` MQTT tool are gone — the leaf
> path they exercised was retired when MQTT auth moved into the
> mapex-mqtt-broker plugin as a local decision off the AssetReadModel
> served by the existing read-model fallback (`GET /internal/assets/:uuid`).
> There is no separate auth callout endpoint to benchmark. A dedicated
> harness for `POST /api/v1/mqtt_certs` (cert issuance) is pending.
> Until then, only the HTTP CRUD scenarios below are exercised end to end.

## Overview

The assets service exposes a standard CRUD HTTP API built on Fiber, backed by MongoDB for persistence and a TieredCache (L0 = in-process RAM via Ristretto, L1 = disk, L2 = MinIO) for read acceleration. Two modules are measured here: `assets` (device registry) and `assettemplates` (template catalogue).

Benchmarks use `hey` to fire HTTP requests at the running service. Prometheus metrics are collected from `GET /metrics` after each run.

### HTTP Scenarios

| Scenario | Endpoint | Key behaviour |
|---|---|---|
| `list-assets` | `GET /api/v1/assets?page=1&perPage=20` | MongoDB aggregation with $lookup, multi-tenant PathKey filter |
| `get-asset` | `GET /api/v1/assets/:id` | TieredCache lookup (L0 → L1 → L2 MinIO → DB fallback) |
| `create-asset` | `POST /api/v1/assets` | MongoDB insert + NATS FANOUT publish + MinIO read-model write |
| `get-template` | `GET /api/v1/asset_templates/:id` | Template read, Redis cache with 24h TTL |

---

## Test Environment

| Property | Value |
|---|---|
| CPU | Intel Core i9-13900K (8 P-cores HT + 16 E-cores = 32 logical) |
| RAM | 64 GB DDR5 |
| OS | Ubuntu 24.04 LTS, kernel 6.17 |
| Go | 1.25.x |
| Isolation | cgroup v2 cpuset — benchmark cgroup pinned to cores 16-31 |
| MongoDB | 7.x, replicaSet=rs0, local |
| Redis | 7.x, local |
| NATS | 2.10.x with JetStream (port 4222), local |
| MinIO | RELEASE.2024-x, local |

---

## Test Configuration

| Parameter | Value |
|---|---|
| Request count (read scenarios) | 100,000 |
| Request count (create scenario) | 10,000 |
| Concurrency | 200 parallel workers |
| CPU configurations | 1, 2, 4, 8, 16 cores |
| Warmup | 1,000 requests at 50 concurrency — discarded |
| Cool-down between runs | 5 seconds |
| Seed HTTP assets | 1,000 documents (deterministic IDs) |
| Seed template | 1 document (ID `000000000000000000000010`) |

---

## Methodology

The benchmark follows the [general benchmark standard](../../../../../documentations/architecture/microservices/general/benchmarks/benchmark_standard_pattern.md) and uses common modules from `scripts/benchmarks/common/`.

### HTTP Scenarios — Step-by-step

1. **CPU isolation** — run `sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh` once per session. This pins all OS processes to cores 0-15 and creates `/sys/fs/cgroup/benchmark` on cores 16-31.

2. **Seed** — `full-benchmark.sh` calls `./seed.sh setup` automatically (idempotent). This inserts 1,000 HTTP assets + 1 template into MongoDB and uploads 100 AssetReadModel JSON objects to MinIO.

3. **Build** — the script compiles a fresh binary (`assets_bench`) with `GOWORK=off CGO_ENABLED=0 go build`.

4. **JWT generation** — a 24-hour HS256 JWT is generated using the default `auth_secret` from `config.go` so the auth middleware passes without database calls.

5. **Per-CPU loop** — for each CPU value in `[1, 2, 4, 8, 16]`:
   - Write `cpuset.cpus` on the benchmark cgroup to pin N cores.
   - Start the service with `GOMAXPROCS=N LOG_LEVEL=silent`.
   - Move the process into the cgroup.
   - Wait up to 60 s for `/metrics` to respond (readiness check).
   - **Warmup**: 1,000 requests at 50 concurrency (results discarded).
   - **Measurement**: fire the full request count with `hey`.
   - **Metrics dump**: `curl /metrics` → `.txt` file.
   - Stop service. Wait 5 s cool-down.

6. **Teardown** — `full-benchmark.sh` calls `./seed.sh teardown` at the end. This deletes all seed documents (regex `^0{10,}`), flushes Redis AppCache DB 0, purges the NATS FANOUT stream, and removes MinIO objects.

### What is measured on the client side (hey)

- `Requests/sec` — throughput under sustained load.
- Latency: avg, p50, p95, p99 — all reported in milliseconds.
- HTTP status code distribution — 2xx/4xx/5xx totals to detect errors during the run.

### What is measured on the server side (/metrics)

| Metric | What it tells you |
|---|---|
| `assets_asset_operations_total{status="success"}` | Successful CRUD completions |
| `assets_asset_operation_duration_seconds` | Server-side op latency (excludes network) |
| `assets_asset_cache_total{result="hit\|miss"}` | TieredCache effectiveness |
| `assets_template_operations_total` | Template CRUD completions |
| `assets_template_cache_total` | Template cache hit/miss |
| `process_resident_memory_bytes` | RSS — total memory footprint |
| `go_memstats_heap_alloc_bytes` | Live heap allocation |
| `go_goroutines` | Goroutine count — detects leaks under load |
| `go_gc_duration_seconds` | GC pause distribution |

---

## Results

> **Test environment**: 100K requests per scenario (10K for create-asset), 200 concurrency, cgroup v2 CPU isolation.

### list-assets — GET /api/v1/assets

| CPU | Cores | Throughput (req/s) | Latency avg | Latency p50 | Latency p95 | Latency p99 | RSS (MB) |
|-----|-------|--------------------|-------------|-------------|-------------|-------------|----------|
| 1 | 16 | 987 | 202.5ms | 195.3ms | 260.9ms | 311.6ms | 56.6 |
| 2 | 16-17 | **1,065** | **187.7ms** | **182.0ms** | 226.5ms | 316.2ms | 71.0 |
| 4 | 16-19 | 1,062 | 188.2ms | 182.1ms | 227.1ms | 321.9ms | 99.4 |
| 8 | 16-23 | 1,052 | 190.0ms | 183.1ms | 230.8ms | 329.8ms | 145.7 |
| 16 | 16-31 | 1,043 | 191.7ms | 184.8ms | 233.7ms | 333.1ms | 212.6 |

> **Bottleneck: MongoDB** — aggregation pipeline with `$lookup` is the limiting factor. Throughput is flat at ~1K req/s regardless of CPU count.

### get-asset — GET /api/v1/assets/:id

| CPU | Cores | Throughput (req/s) | Latency avg | Latency p50 | Latency p95 | Latency p99 | RSS (MB) |
|-----|-------|--------------------|-------------|-------------|-------------|-------------|----------|
| 1 | 16 | 4,976 | 40.2ms | 40.0ms | 44.8ms | 45.7ms | 58.3 |
| 2 | 16-17 | 10,130 | 19.7ms | 18.9ms | 24.2ms | 26.3ms | 70.3 |
| 4 | 16-19 | 20,340 | 9.8ms | 9.5ms | 12.9ms | 14.4ms | 96.7 |
| 8 | 16-23 | **30,982** | **6.4ms** | **6.1ms** | 9.7ms | 13.3ms | 144.1 |
| 16 | 16-31 | 30,258 | 6.6ms | 6.2ms | 10.1ms | 13.3ms | 235.8 |

> **Scales linearly up to 8 CPUs** — 4,976 → 30,982 req/s (6.2x improvement). TieredCache (L0 Ristretto → L1 disk → L2 MinIO) handles the read path efficiently. Diminishing returns at 16 CPUs.

### create-asset — POST /api/v1/assets

| CPU | Cores | Throughput (req/s) | Latency avg | Latency p50 | Latency p95 | Latency p99 | RSS (MB) |
|-----|-------|--------------------|-------------|-------------|-------------|-------------|----------|
| 1 | 16 | 6,377 | 31.3ms | 29.7ms | 37.5ms | 38.9ms | 60.1 |
| 2 | 16-17 | 12,467 | 16.0ms | 15.0ms | 20.8ms | 32.1ms | 71.7 |
| 4 | 16-19 | 14,965 | 12.9ms | 6.3ms | 44.2ms | 111.4ms | 94.1 |
| 8 | 16-23 | 15,115 | 12.9ms | 6.2ms | 44.6ms | 108.6ms | 138.5 |
| 16 | 16-31 | **15,580** | **12.4ms** | **6.6ms** | 45.1ms | 95.7ms | 225.0 |

> **Note**: P95/P99 latency spikes at 4+ CPUs indicate contention in the write path (MongoDB insert + NATS FANOUT + MinIO).

### get-template — GET /api/v1/asset_templates/:id

| CPU | Cores | Throughput (req/s) | Latency avg | Latency p50 | Latency p95 | Latency p99 | RSS (MB) |
|-----|-------|--------------------|-------------|-------------|-------------|-------------|----------|
| 1 | 16 | 6,838 | 29.2ms | 27.8ms | 34.6ms | 35.5ms | 57.9 |
| 2 | 16-17 | 13,835 | 14.4ms | 13.8ms | 18.5ms | 20.3ms | 69.6 |
| 4 | 16-19 | 27,778 | 7.2ms | 7.0ms | 9.7ms | 11.2ms | 91.1 |
| 8 | 16-23 | **44,600** | **4.5ms** | **4.1ms** | 7.2ms | 10.5ms | 141.7 |
| 16 | 16-31 | 43,429 | 4.6ms | 4.3ms | 7.5ms | 10.1ms | 227.3 |

> **Fastest scenario** — 44,600 req/s at 8 CPUs with 4.1ms P50. Redis-cached templates with 24h TTL. Scales 6.5x from 1→8 CPUs.

---

## Key Findings

1. **get-template is the fastest** — 44,600 req/s at 8 CPUs, 4.1ms P50. Redis-cached with 24h TTL, near-linear scaling up to 8 cores.

2. **get-asset scales well** — 30,982 req/s at 8 CPUs, 6.1ms P50. TieredCache read path (Ristretto → disk → MinIO) is efficient. 6.2x improvement from 1→8 CPUs.

3. **list-assets is MongoDB-bound** — flat at ~1,050 req/s regardless of CPU count. The `$lookup` aggregation pipeline is the bottleneck. Adding CPUs only increases memory without improving throughput.

4. **Diminishing returns at 16 CPUs** — Both get-asset and get-template show slight regression going from 8→16 cores (contention overhead exceeds parallelism gains).

5. **RSS scales with CPU** — 57MB at 1 CPU → 230MB at 16 CPUs for HTTP scenarios (goroutine stacks, GC buffers).

---

## Conclusions

### Performance Summary

| Scenario | Best CPU | Peak RPS | P50 Latency | Bottleneck |
|---|---|---|---|---|
| get-template | 8 | **44,600** | 4.1ms | CPU (scales linearly) |
| get-asset | 8 | **30,982** | 6.1ms | CPU (scales linearly) |
| create-asset | 16 | 15,580 | 6.6ms | MongoDB writes + NATS fanout |
| list-assets | 2 | 1,065 | 182.0ms | MongoDB aggregation ($lookup) |

### Production Recommendations

1. **Assets service: 4-8 CPUs** — optimal for the HTTP workload. Beyond 8 cores returns diminish on read scenarios; create scenario keeps modest gains to 16.
2. **list-assets optimization** — MongoDB `$lookup` aggregation is the bottleneck at ~1K req/s. Consider denormalization or read-model projection for high-throughput listing.

---

## Reproduction

### Prerequisites

```bash
# Tools
go install github.com/rakyll/hey@latest
go install github.com/nats-io/natscli/nats@latest
sudo apt install mongosh redis-tools jq bc

# MinIO Client (optional — used by seed.sh for L2 cache population)
mc alias set local http://localhost:9000 mapexos_admin mapexos_admin_secret_change_me
```

### Infrastructure requirements

- MongoDB `mongodb://localhost:27017/?replicaSet=rs0` — replica set required for transactions
- Redis `localhost:6379`
- NATS `nats://localhost:4222` with JetStream enabled, user `service` / `service_secret`
- MinIO `localhost:9000` with buckets `mapex-assets` and `mapex-templates`

### Steps

```bash
# 1. CPU isolation (one-time per session, requires root)
sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh

# 2. Full benchmark — runs HTTP scenarios, 100K requests, CPU 1 2 4 8 16
cd services/assets
bash docs/benchmarks/scripts/full-benchmark.sh

# 3. Custom request count and CPU list
bash docs/benchmarks/scripts/full-benchmark.sh 1000000 "16 12 8 4 2 1"

# 4. Smoke test (10K requests, 1 CPU)
bash docs/benchmarks/scripts/full-benchmark.sh 10000 "1"

# 5. Run only specific scenarios via SCENARIOS env var
SCENARIOS="get-asset" bash docs/benchmarks/scripts/full-benchmark.sh 50000 "4 8"

# 6. Help — shows all env vars and usage
bash docs/benchmarks/scripts/full-benchmark.sh --help

# 7. When done with all benchmarking sessions
sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh teardown
```

### Script locations

| Script | Purpose |
|---|---|
| `scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh` | CPU isolation via cgroup v2 (run once as root) |
| `docs/benchmarks/scripts/config.sh` | Service constants + common module bootstrap |
| `docs/benchmarks/scripts/seed.sh` | MongoDB + MinIO seed / teardown |
| `docs/benchmarks/scripts/full-benchmark.sh` | Main benchmark runner (calls seed.sh automatically) |
| `docs/benchmarks/seed/payloads/create-asset.json` | POST body for create-asset scenario |
| `docs/benchmarks/seed/payloads/update-asset.json` | PATCH body (available for manual testing) |
| `docs/benchmarks/seed/payloads/create-asset-template.json` | POST body for template creation (manual testing) |

### Output artifacts

Results are written to `docs/benchmarks/results/` inside the service directory. After a run, move them to the documentation monorepo for permanent storage:

```
documentations/architecture/microservices/services/assets/results/<bench_tag>/
├── test-<scenario>-cpu<N>-hey.txt       # Raw hey output
├── test-<scenario>-cpu<N>-metrics.txt   # Prometheus metrics dump
├── benchmark-<scenario>-<timestamp>.txt # Human-readable report
└── benchmark-<scenario>-<timestamp>.csv # Machine-readable CSV
```

### CSV columns

| Column | Description |
|---|---|
| `cpu` | GOMAXPROCS value |
| `cores` | cpuset range |
| `concurrency` | hey -c value |
| `total_reqs` | Total requests fired |
| `duration_sec` | hey total duration |
| `rps` | Requests per second |
| `lat_avg_ms` | Average latency |
| `lat_p50_ms` | p50 latency |
| `lat_p95_ms` | p95 latency |
| `lat_p99_ms` | p99 latency |
| `asset_ops_ok` | `assets_asset_operations_total{status="success"}` |
| `asset_ops_err` | `assets_asset_operations_total{status="error"}` |
| `asset_op_avg_ms` | Derived from `assets_asset_operation_duration_seconds` sum/count |
| `asset_cache_hit` | `assets_asset_cache_total{result="hit"}` |
| `asset_cache_miss` | `assets_asset_cache_total{result="miss"}` |
| `template_ops_ok` | `assets_template_operations_total{status="success"}` |
| `template_ops_err` | `assets_template_operations_total{status="error"}` |
| `template_op_avg_ms` | Derived from `assets_template_operation_duration_seconds` sum/count |
| `template_cache_hit` | `assets_template_cache_total{result="hit"}` |
| `template_cache_miss` | `assets_template_cache_total{result="miss"}` |
| `rss_mb` | `process_resident_memory_bytes / 1MiB` |
| `heap_mb` | `go_memstats_heap_alloc_bytes / 1MiB` |
| `goroutines` | `go_goroutines` |
| `gc_count` | `go_gc_duration_seconds_count` |
| `go_threads` | `go_threads` |
| `http_2xx` | hey 2xx response count |
| `http_4xx` | hey 4xx response count |
| `http_5xx` | hey 5xx response count |

---

## Pending — MQTT control-plane benchmarks

- **`POST /api/v1/mqtt_certs`** (`mqttcerts` module) — device cert issuance. Local ECDSA P-256 signing + Mongo insert + fanout publish. Target: measure issuance rate + p95 latency at fixed concurrency; verify `assets_mqttcerts_ca_ready` stays at `1` (no `OnMount` retry storms under load).
- **`GET /internal/assets/:assetUUID`** (L3 read-model fallback) — already exercised by the existing `get-asset` HTTP scenario from a JWT-gated angle. A complementary benchmark on the internal X-API-Key path would isolate the read-only payload the broker plugin consumes. Not strictly necessary while warm-path L1/L2 hit rates stay healthy.

There is NO MQTT auth-callout endpoint to benchmark — the broker plugin decides every CONNECT locally off the AssetReadModel returned by L1/L2/L3, with bcrypt + cert-serial-equality on the broker thread.

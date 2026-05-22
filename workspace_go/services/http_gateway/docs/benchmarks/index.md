# HTTP Gateway - Performance Benchmarks

> **Date**: 2026-02-21
> **Version**: Double-buffer (Redis cache + lazy parse + PubAckFuture)
> **Tag**: `double-buffer`

---

## Overview

This benchmark measures the throughput and latency of the `http_gateway` service under sustained HTTP load across four authentication strategies (JWT, API Key, IP Whitelist, and None) and five CPU configurations (1, 2, 4, 8, and 16 cores). The service receives IoT webhook events via `POST /api/v1/events`, validates the incoming request against a cached DataSource configuration, and publishes the event to a NATS JetStream subject. All tests use cgroup v2 cpuset isolation to eliminate scheduler noise. Results confirm near-linear CPU scaling up to 4 cores (~10K to 45K req/s per auth type), a throughput ceiling of ~140K req/s at 16 cores where NATS/Redis I/O becomes the bottleneck, and a 100% success rate across all 20 million total requests benchmarked.

---

## 1. Test Environment

| Parameter | Value |
|-----------|-------|
| **Machine** | 16 cores, 94Gi RAM |
| **CPU** | 13th Gen Intel Core i9-13900K |
| **Go** | go1.25.3 |
| **Isolation** | cgroup v2 cpuset (cores 16-31) |
| **MongoDB** | dev-http_gateway (local replica set) |
| **Redis** | local (DB 0, service-private AppCache) |
| **NATS** | local JetStream (stream: PROCESSOR-JS-EXECUTE) |

### Test Configuration

| Parameter | Value |
|-----------|-------|
| **Requests per test** | 1,000,000 |
| **Concurrency** | 200 workers |
| **Payload** | 577 bytes (IoT networkStatus event) |
| **CPU configurations** | 1, 2, 4, 8, 16 cores |
| **Auth types tested** | JWT, API Key, IP Whitelist, None |
| **Load tool** | [hey](https://github.com/rakyll/hey) |
| **Total test runs** | 20 (4 auth types x 5 CPU configs) |
| **Total requests** | 20,000,000 |

### CPU Isolation

Tests use **cgroup v2 cpuset** to isolate the service on specific CPU cores (16-31), preventing OS scheduler noise. The `GOMAXPROCS` value matches the number of assigned cores.

---

## 2. Optimizations Applied

Three optimizations were implemented before running these benchmarks:

### 2.1 Redis GetOrSetEx Cache for DataSource (TTL 24h)

**Problem**: Every incoming event required a MongoDB `FindById` to fetch the DataSource configuration (~10.93ms average).

**Solution**: Cache-aside pattern using `GetOrSetEx` from the infrastructure package. On cache miss, fetches from MongoDB and stores in Redis. Cache is invalidated on Create/Update/Delete.

**Impact**: 99.95% cache hit rate. Throughput +45-59% at 1 CPU.

### 2.2 Lazy Body Parse

**Problem**: The request body was parsed twice — once in `CustomAuthMiddleware` (to include in auth failure events) and again in the event handler. Since 99.9% of requests pass auth, the middleware parse was wasted work.

**Solution**: Body is only parsed lazily on the error path via `parseEventBody()`. The happy path skips the middleware body parse entirely.

**Impact**: Additional +5-6% throughput.

### 2.3 Per-Message PubAckFuture

**Problem**: `PublishAsyncComplete()` creates a convoy effect — all 200 goroutines block waiting for ALL pending NATS messages instead of just their own. Processing time was ~7.4ms due to this wait.

**Solution**: Each goroutine now waits on its own `PubAckFuture.Ok()/Err()` channel. Implemented in the shared infrastructure package (`packages/infrastructure/nats/publish.go`) so all Go services benefit.

**Impact**: Processing time dropped from 7.4ms to 0.66ms. p99 latency from 35.7ms to 24.7ms at 1 CPU.

---

## 3. Results by Auth Type

### 3.1 JWT (HMAC-SHA256)

| CPU | Cores | Req/s | Lat avg | Lat p50 | Lat p95 | Lat p99 | RSS | Heap | 2xx |
|-----|-------|-------|---------|---------|---------|---------|-----|------|-----|
| 1 | 16 | 9,881/s | 20.20ms | 19.50ms | 24.20ms | 25.20ms | 52MB | 19MB | 1,000,000 |
| 2 | 16-17 | 19,552/s | 10.20ms | 9.90ms | 13.10ms | 14.70ms | 63MB | 42MB | 1,000,000 |
| 4 | 16-19 | 39,277/s | 5.10ms | 5.00ms | 7.70ms | 9.10ms | 84MB | 61MB | 1,000,000 |
| 8 | 16-23 | 76,880/s | 2.60ms | 2.50ms | 4.30ms | 5.40ms | 130MB | 131MB | 1,000,000 |
| 16 | 16-31 | 126,666/s | 1.60ms | 1.40ms | 2.70ms | 4.40ms | 216MB | 260MB | 1,000,000 |

Scaling: 1-2=2.0x  2-4=2.0x  4-8=2.0x  8-16=1.6x

### 3.2 API Key (Header Comparison)

| CPU | Cores | Req/s | Lat avg | Lat p50 | Lat p95 | Lat p99 | RSS | Heap | 2xx |
|-----|-------|-------|---------|---------|---------|---------|-----|------|-----|
| 1 | 16 | 11,099/s | 18.00ms | 17.40ms | 21.80ms | 22.70ms | 52MB | 31MB | 1,000,000 |
| 2 | 16-17 | 21,855/s | 9.10ms | 8.90ms | 11.70ms | 13.40ms | 62MB | 37MB | 1,000,000 |
| 4 | 16-19 | 43,691/s | 4.60ms | 4.50ms | 6.90ms | 8.20ms | 84MB | 52MB | 1,000,000 |
| 8 | 16-23 | 86,783/s | 2.30ms | 2.20ms | 3.70ms | 4.80ms | 129MB | 109MB | 1,000,000 |
| 16 | 16-31 | 140,112/s | 1.40ms | 1.30ms | 2.40ms | 3.50ms | 212MB | 182MB | 1,000,000 |

Scaling: 1-2=2.0x  2-4=2.0x  4-8=2.0x  8-16=1.6x

### 3.3 IP Whitelist (CIDR Check)

| CPU | Cores | Req/s | Lat avg | Lat p50 | Lat p95 | Lat p99 | RSS | Heap | 2xx |
|-----|-------|-------|---------|---------|---------|---------|-----|------|-----|
| 1 | 16 | 10,928/s | 18.30ms | 17.70ms | 22.10ms | 23.00ms | 52MB | 19MB | 1,000,000 |
| 2 | 16-17 | 21,571/s | 9.30ms | 9.00ms | 11.90ms | 13.50ms | 62MB | 34MB | 1,000,000 |
| 4 | 16-19 | 43,390/s | 4.60ms | 4.50ms | 6.90ms | 8.20ms | 82MB | 85MB | 1,000,000 |
| 8 | 16-23 | 85,237/s | 2.30ms | 2.20ms | 3.80ms | 4.80ms | 128MB | 115MB | 1,000,000 |
| 16 | 16-31 | 139,143/s | 1.40ms | 1.30ms | 2.40ms | 3.60ms | 213MB | 262MB | 1,000,000 |

Scaling: 1-2=2.0x  2-4=2.0x  4-8=2.0x  8-16=1.6x

### 3.4 None (No Auth)

| CPU | Cores | Req/s | Lat avg | Lat p50 | Lat p95 | Lat p99 | RSS | Heap | 2xx |
|-----|-------|-------|---------|---------|---------|---------|-----|------|-----|
| 1 | 16 | 11,384/s | 17.60ms | 17.00ms | 21.30ms | 22.20ms | 52MB | 23MB | 1,000,000 |
| 2 | 16-17 | 22,592/s | 8.80ms | 8.60ms | 11.40ms | 12.80ms | 63MB | 29MB | 1,000,000 |
| 4 | 16-19 | 45,202/s | 4.40ms | 4.30ms | 6.70ms | 8.00ms | 83MB | 52MB | 1,000,000 |
| 8 | 16-23 | 88,981/s | 2.20ms | 2.10ms | 3.60ms | 4.70ms | 128MB | 89MB | 1,000,000 |
| 16 | 16-31 | 144,015/s | 1.40ms | 1.30ms | 2.30ms | 3.40ms | 215MB | 189MB | 1,000,000 |

Scaling: 1-2=2.0x  2-4=2.0x  4-8=2.0x  8-16=1.6x

---

## 4. Cross-Auth Comparison

### 4.1 Throughput (Req/s)

| CPU | JWT | API Key | IP Whitelist | None |
|-----|-----|---------|--------------|------|
| 1 | 9,881 | 11,099 | 10,928 | 11,384 |
| 2 | 19,552 | 21,855 | 21,571 | 22,592 |
| 4 | 39,277 | 43,691 | 43,390 | 45,202 |
| 8 | 76,880 | 86,783 | 85,237 | 88,981 |
| 16 | 126,666 | 140,112 | 139,143 | 144,015 |

### 4.2 Auth Overhead (CPU=1, vs None baseline)

| Auth Type | Req/s | Delta | Overhead | Reason |
|-----------|-------|-------|----------|--------|
| None | 11,384 | — | Baseline | — |
| API Key | 11,099 | -285 | -2.5% | String comparison |
| IP Whitelist | 10,928 | -456 | -4.0% | CIDR parsing |
| JWT | 9,881 | -1,503 | -13.2% | HMAC-SHA256 crypto |

At 16 CPUs: None=144K, apiKey=140K, ipWhiteList=139K, JWT=127K. JWT overhead persists at scale (~12% slower) due to crypto cost per request.

### 4.3 Latency p99 (ms)

| CPU | JWT | API Key | IP Whitelist | None |
|-----|-----|---------|--------------|------|
| 1 | 25.2ms | 22.7ms | 23.0ms | 22.2ms |
| 2 | 14.7ms | 13.4ms | 13.5ms | 12.8ms |
| 4 | 9.1ms | 8.2ms | 8.2ms | 8.0ms |
| 8 | 5.4ms | 4.8ms | 4.8ms | 4.7ms |
| 16 | 4.4ms | 3.5ms | 3.6ms | 3.4ms |

---

## 5. Scaling Analysis

### 5.1 CPU Scaling Efficiency

| Transition | JWT | API Key | IP Whitelist | None |
|------------|-----|---------|--------------|------|
| 1-2 cores | 2.0x | 2.0x | 2.0x | 2.0x |
| 2-4 cores | 2.0x | 2.0x | 2.0x | 2.0x |
| 4-8 cores | 2.0x | 2.0x | 2.0x | 2.0x |
| 8-16 cores | 1.6x | 1.6x | 1.6x | 1.6x |

**Analysis:**
- **1-8 CPUs**: Near-perfect linear scaling (2.0x per doubling) — CPU bound
- **8-16 CPUs**: Diminishing returns (1.6x) — NATS/Redis I/O becoming the bottleneck

### 5.2 Memory Scaling

| CPU | RSS (avg) | Heap (avg) | Goroutines | GC Max Pause |
|-----|-----------|------------|------------|--------------|
| 1 | ~52MB | ~23MB | 220 | 0.05ms |
| 2 | ~63MB | ~36MB | 220 | 0.15ms |
| 4 | ~83MB | ~63MB | 220 | 0.17ms |
| 8 | ~129MB | ~111MB | 220 | 0.51ms |
| 16 | ~214MB | ~223MB | 219 | 0.80ms |

**Analysis:**
- RSS grows linearly with CPU count (~11MB per additional CPU)
- Heap grows faster at high CPU counts due to more concurrent allocations
- Goroutine count stays constant (~220) — controlled by concurrency parameter
- GC max pause stays under 1ms across all configurations (excellent)

---

## 6. Production Capacity Planning

### 6.1 Throughput Estimates

Based on the benchmarks, expected production throughput (with real-world overhead ~20% reduction):

| CPU Allocation | Expected Throughput | Latency p99 |
|----------------|--------------------:|-------------|
| 1 core | ~8,000-9,000 req/s | ~25ms |
| 2 cores | ~16,000-18,000 req/s | ~14ms |
| 4 cores | ~32,000-36,000 req/s | ~9ms |
| 8 cores | ~64,000-72,000 req/s | ~5ms |

### 6.2 Sizing Recommendations

| Traffic Tier | Daily Events | Peak req/s | Recommended CPU | Memory |
|-------------|-------------|------------|-----------------|--------|
| Small | <10M | <500 | 1 core | 128MB |
| Medium | 10M-100M | 500-5,000 | 2 cores | 256MB |
| Large | 100M-500M | 5,000-20,000 | 4 cores | 512MB |
| Enterprise | 500M+ | 20,000+ | 8+ cores | 1GB |

### 6.3 Industry Context

For a Go HTTP service performing real work (auth + Redis cache + NATS JetStream publish + Prometheus metrics), these numbers are **above average**:

| Benchmark Context | Typical Range | Our Result |
|-------------------|--------------|------------|
| Go HTTP frameworks (no I/O) | 100K-500K req/s | N/A (not comparable) |
| Go services with DB + cache | 5K-30K req/s (1 CPU) | 10K-11K req/s |
| Go services with messaging | 3K-20K req/s (1 CPU) | 10K-11K req/s |
| Go API gateways (production) | 20K-80K req/s (multi-core) | 144K req/s (16 CPU) |

---

## 7. Key Findings

1. **100% success rate** — Zero errors, zero timeouts across all 20 test runs (20M total requests)
2. **DataSource cache**: 99.95% hit rate (Redis GetOrSetEx, TTL 24h)
3. **NATS reliability**: Per-message PubAckFuture guarantees delivery without convoy effect
4. **Body parse**: Lazy parsing saves ~5-6% throughput on the happy path
5. **I/O ceiling**: At 16 CPUs (~140K req/s), the bottleneck shifts from CPU to NATS/Redis I/O
6. **JWT overhead**: 13.2% slower than no-auth at 1 CPU due to HMAC-SHA256 computation; persists at ~12% even at 16 CPUs
7. **Linear scaling**: Perfect 2.0x scaling up to 8 CPUs; diminishing 1.6x at 8-16
8. **GC excellence**: Max GC pause under 1ms even at 16 CPUs processing 1M requests

---

## 8. Benchmark Methodology

### 8.1 Script

The benchmark is fully automated via `scripts/full-benchmark.sh`. A single command runs all 4 auth types across all CPU configurations:

```bash
# Run all auth types (default), 1M requests, CPU 16 8 4 2 1
./scripts/full-benchmark.sh 1000000 "16 8 4 2 1" double-buffer

# Override auth types via env var
AUTH_TYPES="jwt none" ./scripts/full-benchmark.sh 1000000 "16 8 4 2 1" custom-tag

# Quick smoke test
AUTH_TYPES="none" ./scripts/full-benchmark.sh 1000 "1" smoke-test
```

### 8.2 Process Per Test

1. **Purge NATS stream** — Clean state, prevent backpressure from previous run
2. **Start service** — `GOMAXPROCS=N`, `LOG_LEVEL=silent` for zero noise
3. **Move to cgroup** — Isolate service on specific CPU cores via cgroup v2 cpuset
4. **Benchmark** — 1,000,000 requests at 200 concurrency via `hey`
5. **Collect metrics** — Scrape `/metrics` endpoint (Prometheus format)
6. **Stop service** — Kill process, extract results

### 8.3 What Is Measured

**Client-side** (via `hey`):
- Requests/second, total duration
- Latency: average, p50, p95, p99
- HTTP status code distribution (2xx, 4xx, 5xx)

**Server-side** (via Prometheus `/metrics`):
- Auth success/failure counts and duration
- Events processed/published counts
- NATS publish success/error counts
- Processing duration (handler-only, excludes auth)
- RSS, heap, goroutines, GC cycles, OS threads

### 8.4 Output Files

```
docs/benchmarks/results/{tag}/
├── test-{auth}-cpu{N}-hey.txt              # Raw hey output per test
└── test-{auth}-cpu{N}-metrics.txt          # Raw Prometheus metrics per test
```

### 8.5 Reproducibility

**Prerequisites:**
1. cgroup v2 shield configured on cores 16-31
2. MongoDB, NATS, Redis running locally
3. `hey` installed: `go install github.com/rakyll/hey@latest`
4. Benchmark data sources created (script handles this automatically)

**Isolation guarantees:**
- Service runs in dedicated cgroup with pinned CPU cores
- Each test starts with a fresh process (no state carryover)
- NATS stream is purged between tests
- 2-second cooldown after each test for metric collection

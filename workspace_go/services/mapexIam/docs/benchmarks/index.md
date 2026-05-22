# MapexOS Service — Benchmark Results

> **Date:** 2026-02-23
> **Pattern:** HTTP Flood (100,000 requests per test)

---

## Overview

Throughput and latency of the MapexOS authentication endpoints across 1 to 16 CPU cores. Two scenarios cover opposite performance profiles: **login** (CPU-bound, bcrypt) and **coverage** (I/O-bound, Redis-cached).

| Scenario | Endpoint | What it measures |
|----------|----------|------------------|
| `auth_login` | `POST /auth/login` | bcrypt password verify + JWT generation + Redis session |
| `auth_coverage` | `GET /auth/users/me/coverage` | Org tree resolution via direct and group-inherited memberships |

All tests ran on an i9-13900K (cores 16-31 isolated via cgroup v2), 100,000 requests per test, 1,000 seeded users rotating per request. **10 tests total, 1,000,000 requests, 100% success rate.**

---

## Results

### auth_login — CPU-bound (bcrypt cost 10)

| CPU | Req/s | Latency avg | Latency p50 | Latency p99 | Total time |
|----:|------:|------------:|------------:|------------:|-----------:|
| 1   | 26    | 963 ms      | 1,007 ms    | 1,116 ms    | 64 min     |
| 2   | 52    | 960 ms      | 957 ms      | 1,307 ms    | 32 min     |
| 4   | 104   | 963 ms      | 956 ms      | 1,354 ms    | 16 min     |
| 8   | 207   | 966 ms      | 959 ms      | 1,371 ms    | 8 min      |
| 16  | 366   | 546 ms      | 535 ms      | 981 ms      | 4.5 min    |

### auth_coverage — I/O-bound (Redis + MongoDB)

| CPU | Req/s | Latency avg | Latency p50 | Latency p99 | Total time |
|----:|------:|------------:|------------:|------------:|-----------:|
| 1   | 40,444  | 4.9 ms  | 4.7 ms  | 8.6 ms  | 2.5 s  |
| 2   | 71,697  | 2.8 ms  | 2.7 ms  | 4.6 ms  | 1.4 s  |
| 4   | 108,077 | 1.8 ms  | 1.7 ms  | 4.0 ms  | 0.9 s  |
| 8   | 124,283 | 1.5 ms  | 1.4 ms  | 3.7 ms  | 0.8 s  |
| 16  | 121,479 | 1.5 ms  | 1.3 ms  | 4.6 ms  | 0.8 s  |

---

## Scaling

### auth_login — Perfect linear scaling up to 8 cores

| CPU | Req/s | Per core | Scale factor | Efficiency |
|----:|------:|---------:|-------------:|-----------:|
| 1   | 26    | 26       | 1.0x         | 100%       |
| 2   | 52    | 26       | 2.0x         | 100%       |
| 4   | 104   | 26       | 4.0x         | 100%       |
| 8   | 207   | 26       | 8.0x         | 100%       |
| 16  | 366   | 23       | 14.1x        | 88%        |

Throughput doubles exactly with each CPU doubling. The slight drop at 16 cores is due to concurrency being capped at 200 (not enough in-flight requests to saturate all cores).

### auth_coverage — Peaks at 8 cores, then plateaus

| CPU | Req/s   | Per core | Scale factor | Efficiency |
|----:|--------:|---------:|-------------:|-----------:|
| 1   | 40,444  | 40,444   | 1.0x         | 100%       |
| 2   | 71,697  | 35,849   | 1.8x         | 89%        |
| 4   | 108,077 | 27,019   | 2.7x         | 67%        |
| 8   | 124,283 | 15,535   | 3.1x         | 38%        |
| 16  | 121,479 | 7,592    | 3.0x         | 19%        |

Beyond 8 cores, Redis (single-threaded) becomes the bottleneck — more Go cores don't help.

### Side-by-side

| CPU | auth_login | auth_coverage | Ratio |
|----:|-----------:|--------------:|------:|
| 1   | 26         | 40,444        | 1,557x |
| 4   | 104        | 108,077       | 1,039x |
| 8   | 207        | 124,283       | 600x   |
| 16  | 366        | 121,479       | 332x   |

The 1,500x difference at CPU=1 confirms bcrypt is the sole bottleneck in the login pipeline — not MongoDB, Redis, JWT signing, or Fiber routing.

---

## Key Findings

1. **Login scales linearly with cores.** 26 req/s per core, 100% efficiency from 1 to 8 CPUs. This is the theoretical maximum for bcrypt cost 10 (~38ms per hash).

2. **Coverage delivers 40K-124K req/s.** After warmup, Redis serves cached org trees in under 5ms. The I/O bottleneck means diminishing returns beyond 4 cores.

3. **Latency is predictable from queue theory.** With 25 concurrent requests per core and 38ms bcrypt time: `38ms * 25 = 950ms` — matching the observed ~960ms average. At 16 cores (12.5 per core): `38ms * 12.5 = 475ms` — matching the observed ~546ms.

4. **Zero errors across 1M requests.** No timeouts, no 4xx, no 5xx. The 30s context timeout prevents request cancellation under bcrypt queue pressure.

5. **Both membership paths exercised.** 500 users resolve access via direct memberships, 500 via group-inherited memberships, and all 1,000 via both paths.

---

## Conclusions

### Production sizing

**Login throughput** is fixed at **~26 req/s per core** (bcrypt cost 10). Scale horizontally:

| Target logins/s | Replicas (4 cores each) |
|----------------:|------------------------:|
| 100             | 1                       |
| 500             | 5                       |
| 1,000           | 10                      |

**Coverage throughput** plateaus at **~108K req/s** with 4 cores:

| Target req/s | Replicas (4 cores each) |
|-------------:|------------------------:|
| 100,000      | 1                       |
| 500,000      | 5                       |

### Kubernetes resources

```yaml
resources:
  requests: { cpu: "4", memory: "128Mi" }
  limits:   { cpu: "8", memory: "256Mi" }
```

### How to increase login throughput

| Strategy | Effect |
|----------|--------|
| More replicas | Linear increase, no ceiling |
| More cores per replica | Linear up to ~8, diminishing after |
| Lower bcrypt cost (10 → 9) | 2x throughput, reduces security margin |
| Rate limiting | Reduces peak load, protects from brute force |

---

## Reproduction

```bash
# 1. CPU isolation (one-time, requires root)
sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh

# 2. Full benchmark (100K requests, 5 CPU configs, 10 tests)
bash docs/benchmarks/scripts/full-benchmark.sh

# 3. Smoke test (10K requests, 1 CPU)
bash docs/benchmarks/scripts/full-benchmark.sh 10000 "1"

# 4. Single scenario
SCENARIOS="auth_login" bash docs/benchmarks/scripts/full-benchmark.sh 100000 "4 8"

# 5. Teardown
sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh teardown
```

Raw results: `docs/benchmarks/results/{timestamp}/`
Scripts and seed data: `docs/benchmarks/scripts/`

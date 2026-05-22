# Benchmarks — JS-Workflow-Executor

## Table of Contents

1. [Overview](#1-overview)
2. [Test Environment](#2-test-environment)
3. [Test Configuration](#3-test-configuration)
4. [Methodology](#4-methodology)
5. [Results](#5-results)
6. [Key Findings](#6-key-findings)
7. [Conclusions](#7-conclusions)
8. [Production Capacity Planning](#8-production-capacity-planning)
9. [Reproduction](#9-reproduction)

---

## 1. Overview

This benchmark measures the end-to-end throughput and memory footprint of the JS-Workflow-Executor service under controlled, CPU-isolated conditions. The service consumes NATS JetStream messages from the `WORKFLOW-JS-CODE` stream, runs workflow code node scripts via a Piscina worker thread pool (V8 isolates with `isolated-vm`), and publishes callback results to the `WORKFLOW-RESUME` stream.

Each message represents a workflow code node execution request: the service fetches the script source from TieredCache (L0/L1/L2), optionally loads cached V8 bytecode, executes the script in a sandboxed V8 isolate, and publishes the result (output + statePatch) back to the Go workflow runtime.

Key differences from `js-executor`:
- **No MongoDB/Redis dependencies** — only NATS + MinIO
- **Single script per event** — no decode/validate/transform pipeline
- **Bytecode caching** — V8 compiled bytecode stored in L1/L2 for faster cold starts
- **Callback pattern** — publishes result to `workflow.resume.callback.{instanceId}`

---

## 2. Test Environment

| Property | Value |
|---|---|
| CPU | Intel Core i9-13900K (8 P-cores + 16 E-cores = 32 logical cores) |
| Benchmark cores | E-cores 16–31 (4.3 GHz, no Hyper-Threading) |
| System cores | P-cores 0–15 (confined via cgroup v2 cpuset) |
| RAM | 94 GB |
| OS | Linux (cgroup v2 enabled) |
| Node.js | v24.0.0 |
| NATS | JetStream enabled, local instance |
| LOG_LEVEL | `silent` (zero I/O overhead during measurement) |
| Isolation | cgroup v2 cpuset shield — service and all Piscina workers confined to benchmark cores |

### CPU Configuration Mapping

| CPU_LIMIT | Piscina Workers | Batch Size | Cores Used |
|---|---|---|---|
| 1 | 1 (minimum) | 500 | 16 |
| 2 | 1 | 1,000 | 16–17 |
| 4 | 3 | 2,000 | 16–19 |
| 8 | 7 | 4,000 | 16–23 |
| 16 | 15 | 8,000 | 16–31 |

---

## 3. Test Configuration

| Parameter | Value |
|---|---|
| Events per test | 1,000,000 |
| Message payload | `WorkflowScriptInput` (code node execution) |
| Stream | `WORKFLOW-JS-CODE` |
| Subject | `workflow.js.code` |
| Consumer | `js-workflow-executor-code` (durable pull) |
| NATS credentials | `service` / `service_secret` |
| Metrics endpoint | `http://localhost:8001/metrics` |
| Warmup | Stream pre-populated before service start |
| Cool-down | 5 seconds between CPU configurations |
| Results storage | `docs/benchmarks/results/<tag>/` |

---

## 4. Methodology

### 4.1 Pattern: NATS Stream Drain

This service has no HTTP endpoints to benchmark. The pattern is:

- **Pre-seed** N messages into the `WORKFLOW-JS-CODE` stream (run `seed.sh setup` before the benchmark).
- **Start the service** with `CPU_LIMIT=N` inside the cgroup v2 shield.
- **Measure time** from first message dequeued until stream pending count reaches zero.
- **Throughput** = total messages / drain time (seconds).

### 4.2 Script Execution Flow

Every event processed by the benchmark runs the full production pipeline:

```
NATS fetch (batch)
  └── For each message in batch:
        → JSON parse WorkflowScriptInput
        → TieredCache lookup: script source (L0 RAM / L1 Disk / L2 MinIO)
        → BytecodeCache lookup: V8 bytecode (L1 / L2)
  └── Dispatch to Piscina worker (V8 isolate):
        → Compile script (with bytecode fast path if cached)
        → Build sandbox context: { event, state, variables, nodes }
        → Execute script in isolated-vm (32MB heap, 10s timeout)
        → Extract result: { output, statePatch }
        → Return newBytecode if first compile
  └── Post-execution:
        → Store newBytecode in BytecodeCache (fire-and-forget)
        → Build WorkflowScriptCallback
        → Publish callback to WORKFLOW-RESUME
        → ACK message
```

### 4.3 Measurement Steps

The `full-benchmark.sh` script executes the following for each CPU configuration:

- Kill any leftover service process on port 8001.
- Purge ALL NATS streams (including downstream) to ensure clean state.
- Wait for `seed.sh` to populate the `WORKFLOW-JS-CODE` stream (polls until >= 90% of target and stable for 15 seconds).
- Start the service with `CPU_LIMIT=N LOG_LEVEL=silent` and move it (plus all Piscina worker processes) into the cgroup v2 benchmark cgroup.
- Poll `nats stream info` every 2 seconds; record throughput until pending count reaches 0.
- Wait 5 seconds after drain for metrics to flush, then scrape `GET /metrics`.
- Parse and record: drain time, throughput, event duration (avg/p50/p99), RSS, heap, GC counts, event loop p99.
- Stop the service, cool down 5 seconds.

---

## 5. Results

> **TODO:** Run benchmarks and populate results.

### 5.1 Baseline (bench tag: `baseline`)

**Date:** — | **Events per test:** 1,000,000 | **Node.js:** v24.0.0

| CPU | Workers | Batch | Throughput | Drain Time | RSS | Heap Used | EL p99 | GC Minor | GC Major |
|-----|---------|-------|------------|------------|-----|-----------|--------|----------|----------|
| 16 | 15 | 8,000 | — | — | — | — | — | — | — |
| 8 | 7 | 4,000 | — | — | — | — | — | — | — |
| 4 | 3 | 2,000 | — | — | — | — | — | — | — |
| 2 | 1 | 1,000 | — | — | — | — | — | — | — |
| 1 | 1 | 500 | — | — | — | — | — | — | — |

---

## 6. Key Findings

> **TODO:** Populate after benchmark run.

Expected differences from `js-executor`:
- **Lower throughput per event** — each event runs a full V8 isolate script (not a lightweight pipeline stage)
- **BytecodeCache impact** — first run is slower (compile + produce bytecode), subsequent runs benefit from cached bytecode
- **Simpler pipeline** — no TieredCache for assets/templates, only script source lookup
- **Callback overhead** — each event publishes a callback to WORKFLOW-RESUME (additional NATS write)

---

## 7. Conclusions

> **TODO:** Populate after benchmark run.

Expected recommendation: `CPU_LIMIT=1` with horizontal scaling (same pattern as `js-executor`), since each Piscina worker owns its own V8 isolate and context recycling is per-worker.

---

## 8. Production Capacity Planning

> **TODO:** Populate with actual numbers after benchmark run.

### Scaling Formula

```
Target throughput:    T ev/s
Per-instance rate:    ~X ev/s  (CPU=1, measured)
Replicas needed:      ceil(T / X)
Memory per replica:   ~Y GB request / ~Z GB limit
max_ack_pending:      replicas × 1000
```

### Kubernetes Resource Recommendations

```yaml
resources:
  requests:
    cpu: "1000m"
    memory: "1500Mi"
  limits:
    cpu: "1500m"
    memory: "2000Mi"
env:
  - name: CPU_LIMIT
    value: "1"
  - name: NATS_CONSUMER_MAX_ACK_PENDING
    value: "6000"   # adjust: replicas × 1000
```

---

## 9. Reproduction

### Prerequisites

| Requirement | Notes |
|---|---|
| NATS with JetStream | Stream `WORKFLOW-JS-CODE` must exist |
| MinIO | Workflow script source seeded in `mapex-workflows` bucket |
| Node.js v24 | Via nvm: `nvm use 24` |
| `nats` CLI | `go install github.com/nats-io/natscli/nats@latest` |
| Linux with cgroup v2 | Default on modern distros |
| Root access | Required for `setup-cgroup-shield.sh` |
| 32 logical cores | Adjust `BENCHMARK_CORES` in the shield script for fewer cores |

### Full Benchmark Run

```bash
# From service root: workspace_js/services/js-workflow-executor/

# One-time: set up CPU isolation (requires root)
sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh setup

# Populate MinIO + NATS stream with 1M seed messages
bash docs/benchmarks/scripts/seed.sh setup

# Run full benchmark suite
bash docs/benchmarks/scripts/full-benchmark.sh 1000000 "16 8 4 2 1" baseline

# Teardown
bash docs/benchmarks/scripts/seed.sh teardown

# Optional: tear down CPU isolation
sudo bash scripts/benchmarks/common/cgroup/setup-cgroup-shield.sh teardown
```

### Script Locations

```
docs/benchmarks/
├── index.md                                  # This file
├── scripts/
│   ├── config.sh                             # Service-specific constants
│   ├── full-benchmark.sh                     # Main benchmark runner
│   └── seed.sh                               # Stream population and teardown
└── seed/
    ├── minio/
    │   └── workflow-script-source.js          # Sample workflow code node script
    └── nats/
        └── workflow-code-execute.json         # WorkflowScriptInput payload
```

### Results Location

Prometheus metrics snapshots per CPU configuration:

```
docs/benchmarks/results/
└── <bench-tag>/
    ├── test-cpu<N>-metrics.txt   (full Prometheus dump)
    └── test-cpu<N>-output.log    (service output)
```

# Workflow Service - Performance Benchmarks

> **Status**: Pending — benchmarks not yet executed.

## Overview

This document will contain performance benchmarks for the workflow execution engine,
covering trigger throughput, inline execution latency, and concurrent instance scaling.

## Planned Test Scenarios

1. **Trigger throughput** — messages/sec on `WORKFLOW-TRIGGER` stream
2. **Inline execution latency** — end-to-end for simple DAGs (start → condition → set_state → end)
3. **Fanout performance** — parallel branch execution with varying branch counts
4. **KV checkpoint overhead** — per-step Put latency under load
5. **Archiver batch throughput** — MongoDB BulkWrite with varying batch sizes
6. **Concurrent instances** — scaling behavior with 100/1K/10K simultaneous instances

## How to Run

```bash
cd docs/benchmarks/scripts
./full-benchmark.sh
```

## Results

Results will be published in `results/` after execution.

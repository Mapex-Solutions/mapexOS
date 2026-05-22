# JS‑Executor Documentation

## Overview
JS‑Executor runs dynamic scripts against incoming events. It consumes NATS events, resolves assets/templates, executes scripts in isolated V8 workers, and emits processed payloads. Its key competitive value is the **multi‑tier cache (L0/L1/L2 + fallback)** that scales horizontally while keeping latency low. L2 is MinIO/S3, used as a **distributed read cache**, and is repopulated on‑demand via fallback when missing.

## Responsibilities
- Execute decode/validation/transform scripts in isolated workers
- Maintain tiered cache for assets, templates, and bytecode
- Provide script testing and sample payload endpoints

## Non‑Responsibilities
- Event ingestion (HTTP Gateway)
- Rule evaluation (RuleEngine)
- Event persistence (Events service)

## Primary Data Flow
1. NATS consumer receives execution event
2. Resolve asset/template via TieredCache (L0 → L1 → L2)
3. If L2 miss, fallback calls Assets Service and **repopulates L2**
4. Execute scripts in isolated workers
5. Publish results downstream

## Execution Pipeline (Deep‑Dive)
- **Script sources**: templates stored in MinIO (L2) are the ready‑to‑run script set (decode/validation/transform).
- **Tiered cache**:
  - L0 RAM: hottest templates/assets (short TTL)
  - L1 Disk: warm cache (longer TTL)
  - L2 MinIO/S3: shared read model
  - Fallback: internal HTTP to Assets Service when L2 misses; response writes back to L2
- **Bytecode**:
  - Bytecode is the compiled V8 representation of JavaScript.
  - JS‑Executor caches bytecode in L1/L2 (MinIO) to avoid recompilation.
  - L0 for bytecode is optional; ScriptRegistry already caches compiled scripts in RAM.
- **Isolation + parallelism**:
  - Script execution runs inside isolated V8 contexts (isolated‑vm).
  - Piscina worker pool executes tasks concurrently with CPU‑aware auto‑tuning.

## Cache Strategy (Competitive Differentiator)
- **L0 (RAM)**: fastest, short TTL, keeps hottest items
- **L1 (Disk/NVME)**: medium speed, longer TTL
- **L2 (MinIO/S3)**: shared cache across services (distributed read model)
- **Fallback**: HTTP call to Assets Service on L2 miss; response **writes into L2**
- **Fanout invalidation**: consuming services invalidate L0/L1 on template/asset changes

## Docs Map
- [Architecture](architecture/index.md)
- [Endpoints](endpoints/index.md)
- [Configuration](configuration/index.md)
- [Operations](operations/index.md)
- [Observability](observability/index.md)
- [Tests](tests/index.md)
- [Benchmarks](benchmarks/index.md)

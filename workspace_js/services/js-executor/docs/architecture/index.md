# Architecture

## Patterns
- DDD + Hexagonal (Ports & Adapters)
- Dependency Injection via `tsyringe`

## Project Structure
```
src/
├── bootstrap/                 # Config, logger, NATS, Redis, cache, metrics, Express
├── modules/
│   ├── app/                   # App routes (health)
│   ├── scripts/               # Script API + script registry
│   │   ├── application/        # Services, DTOs, DI, ports
│   │   ├── domain/             # Script domain types
│   │   ├── infrastructure/     # Cache adapters
│   │   └── interfaces/http/    # HTTP routes + handlers
│   ├── engine/                # Execution engine
│   │   ├── domain/             # Script execution core
│   │   └── infrastructure/     # Piscina workers + bytecode cache
│   └── events/                # NATS consumers
└── shared/
    └── configuration/          # Config definitions + module registry
```

## Module Responsibilities
- `scripts`: script test endpoints and template utilities
- `engine`: isolated execution and worker pool
- `events`: consume NATS execution events

## Execution Flow
1. Consume event from NATS
2. Resolve asset/template + scripts (tiered cache)
3. Execute scripts in isolated workers
4. Emit results downstream

## Script Execution Internals
- [Services Layer](services/index.md)
- [Compression & Bytecode](compression/index.md)

## Tiered Cache (L0/L1/L2 + Fallback)
JS‑Executor uses the shared TieredCache implementation:
- **L0 RAM**: fastest, short TTL for hot data
- **L1 Disk**: NVME/disk cache for warm data
- **L2 MinIO/S3**: distributed cache shared across services
- **Fallback**: HTTP call to Assets Service on L2 miss and **repopulates L2**

### Cache targets
- **Assets**: JSON read models in `MINIO_ASSETS_BUCKET`
- **Templates**: JSON in `MINIO_TEMPLATES_BUCKET`
- **Bytecode**: compiled scripts in `MINIO_BYTECODE_BUCKET`

### CQRS Read Model Flow
1. Assets/Templates are written to MongoDB (source of truth)
2. Assets Service writes the read model to L2 (MinIO)
3. JS‑Executor loads L2 to build L1/L0 locally
4. If L2 miss → fallback fetches from Assets Service and **writes back to L2**

### Invalidation
Template/asset updates publish fanout events; consumers invalidate L0/L1 so the next read reloads from L2.

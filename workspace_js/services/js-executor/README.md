# JS‑Executor Service

## Summary
Executes MapexOS scripts in isolated V8 workers. Consumes NATS events, resolves assets/templates, and runs decode/validation/transform pipelines.

## Responsibilities
- Execute scripts safely and efficiently
- Cache assets/templates/bytecode
- Expose script testing endpoints

## Architecture
Node.js service using DDD + Hexagonal patterns with `scripts`, `engine`, and `events` modules.

## Documentation
Deep‑dive documentation (architecture, endpoints, configuration, observability, tests, benchmarks):
- [docs/index.md](docs/index.md)

## How to run
```bash
# dev
npm run dev

# build
npm run build

# start
npm run start
```

## Observability
- Health: `GET /health`
- Metrics: `GET /metrics`

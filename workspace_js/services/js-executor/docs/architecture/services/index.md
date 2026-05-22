# Services: ScriptService and ScriptEngineService

## Overview
The services layer orchestrates script execution, caching, and compilation. Two services cooperate:

```
ScriptService (orchestration)
├── Script storage & retrieval
├── Pipeline management
└── NATS publishing

ScriptEngineService (execution)
├── Isolated‑VM management
├── Bytecode caching
├── Script compilation
└── Error sanitization
```

## Cache Strategy (service-level)
- **Script Source Cache**: `SCRIPT:{assetId}:SCRIPTS` (JSON with decode/validation/transform)
- **Bytecode Cache**: `SCRIPT:{assetId}:DECODE|VALIDATION|TRANSFORM`

## ScriptService
### Purpose
Orchestrates the end‑to‑end pipeline (NATS → cache → execution → publish).

### Dependencies
- `ScriptEngineService`
- `RedisService`
- `NatsBus`
- `Logger`

### Key Methods
#### executeScripts(message: ScriptProcessorMessage)
- Resolves asset ID
- Loads scripts from cache or API
- Runs decode → validation → transform
- Publishes results to NATS

#### scripsTest(payload, scripts)
- Runs the pipeline without publishing
- Used by HTTP test endpoint

#### fetchAssetScripts(assetId)
- Loads from cache, else fallback
- Stores compressed scripts with TTL

## ScriptEngineService
### Purpose
Executes scripts inside **isolated‑vm** and manages bytecode caching.

### Dependencies
- `RedisService`
- `Logger`

### Key Methods
#### runScriptPipeline(rawPayload, userScripts)
- Executes decode → validation → transform
- Enforces timeouts and isolation
- Returns structured execution result

## Why this design
- Clear separation of orchestration vs execution
- Safe isolation of untrusted scripts
- Bytecode caching reduces compilation cost

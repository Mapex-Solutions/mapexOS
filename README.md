# MapexOS

> **IoT-first, but not limited to IoT.**
> MapexOS doesn't see devices or sensors — it sees **Assets**.
> Any source. Any protocol. One abstraction.
>
> **Connect. Automate. Scale.** — The open platform for data integration
> and intelligent automation.

```
   Sources                       MapexOS                         Destinations
   ───────                       ───────                         ────────────
   Devices ──┐                                              ┌── Webhooks / APIs
   Gateways ─┤   Ingest → Validate → Transform → Route →    ├── Slack / Teams / Email
   APIs ─────┼──        Store / Notify / Automate           ├── NATS / MQTT
   Apps ─────┤                                              └── Custom plugins
   3rd-party ┘
```

This repository holds the source for every MapexOS service, the
Vue 3 frontend, and the shared Go and TypeScript packages.
The runnable distribution — Docker Compose plus pre-built multi-arch
images — lives in
[`mapexOSDeploy`](https://github.com/Mapex-Solutions/mapexOSDeploy).

[Versão em português](./README_pt.md) · [Documentation site](https://mapexos.io)

---

## What MapexOS gives you

| | |
|---|---|
| **Telemetry ingestion** | Authenticated HTTP and MQTT. The MQTT broker (a Mosquitto plugin) makes every CONNECT/PUBLISH decision locally off a three-tier cache (Pebble → MinIO → HTTP fallback) — no round-trips on the hot path. |
| **Generic data conversion** | Per-asset *preprocess → validate → convert* pipelines written in JavaScript and executed in V8 isolates. Accept any device payload and normalize it into the platform's event schema without touching the source code. |
| **Dynamic event routing** | Match rules over arbitrary event fields — `payload.temperature > 30`, `device.tag in [...]`, JSONPath, regex. Rules live in MongoDB and reload at runtime; no schema migrations, no redeploys. |
| **Workflow engine** | A DAG runtime inspired by [Temporal.io](https://temporal.io) — deterministic execution, retries, timers, sub-workflows, idempotent triggers. State persists across process restarts. |
| **Plugin UI (n8n-style)** | Custom workflow nodes ship as plugins served from a CDN-style manifest registry. New HTTP / MQTT / NATS / Slack / database connectors land in the editor without rebuilding the frontend. |
| **Multi-tenant by design** | Organization hierarchy (parent → children), per-org isolation, shared templates that propagate to descendant orgs. |
| **RBAC + groups** | Fine-grained permissions and group memberships, evaluated centrally by `mapex-iam`. Every cross-service call carries an identity context. |
| **Self-hostable, multi-arch** | One `docker compose up -d` boots the whole stack. Images are published for `linux/amd64` and `linux/arm64` — runs on Linux servers, Apple Silicon, Windows + Docker Desktop, and Raspberry Pi 4/5. |

---

## The MapexOS ecosystem

MapexOS is split across four open repositories. Most users only need
the deploy repo; the others are for contributors, broker operators,
and Go integrators.

| Repository | Role |
|---|---|
| **[mapexOS](https://github.com/Mapex-Solutions/mapexOS)** *(this repo)* | Source for the eleven backend services and the Vue 3 frontend. |
| **[mapexOSDeploy](https://github.com/Mapex-Solutions/mapexOSDeploy)** | Docker Compose distribution that pulls pre-built images from Docker Hub. **Start here to run the platform.** |
| **[mapexMQTTBroker](https://github.com/Mapex-Solutions/mapexMQTTBroker)** | The production MQTT broker — Eclipse Mosquitto v2 plus the in-house Go plugin that handles auth, ACL, presence, and ingress in a single `.so`. |
| **[mapexGoKit](https://github.com/Mapex-Solutions/mapexGoKit)** | Shared Go libraries used by every Go service — HTTP middleware, NATS helpers, observability, validation, contracts. |

---

## What's inside this repo

```
mapexOS/
├── workspace_go/                # Go services (DDD + hexagonal)
│   ├── services/
│   │   ├── mapexIam/            # users, organizations, roles, RBAC, auth
│   │   ├── http_gateway/        # webhook ingestion, datasource registry
│   │   ├── assets/              # IoT assets, templates, EVA fields
│   │   ├── router/              # event routing, match rules
│   │   ├── events/              # ClickHouse storage, 7 NATS consumers
│   │   ├── triggers/            # 8 executors (HTTP, MQTT, NATS, NATS-JS, NATS-KV, NATS-OBJ, NATS-RPC, Webhook)
│   │   ├── workflow/            # Temporal-inspired DAG engine + plugins + credentials
│   │   └── mapexVault/          # credential vault, PKI authority
│   └── packages/
│       └── contracts/           # cross-service DTOs (the single source of truth)
│
├── workspace_js/                # Node services + the frontend
│   ├── services/
│   │   ├── js-executor/             # V8 isolates for IoT event scripts
│   │   ├── js-workflow-executor/    # V8 isolates for workflow code nodes
│   │   └── plugin-marketplace-mock/ # static plugin manifest CDN (dev)
│   ├── apps/
│   │   └── mapexOS/             # Vue 3 + Quasar SPA
│   └── packages/                # Zod schemas, typed API wrappers, infrastructure clients
│
└── docs/                        # finalized public documentation
```

---

## Quick start

To **run** the platform, head to `mapexOSDeploy`:

```bash
git clone https://github.com/Mapex-Solutions/mapexOSDeploy.git
cd mapexOSDeploy
docker compose up -d
```

Open <http://localhost> and log in as `admin@mapex.local` /
`mapex@123`. The full walkthrough — including a temperature sensor
quickstart over HTTP or MQTT — lives in
[`mapexOSDeploy/quickstart/`](https://github.com/Mapex-Solutions/mapexOSDeploy/tree/main/quickstart).

---

## Why this exists

Most "IoT platforms" stop at ingestion and dashboards. The hard part
starts after the data lands: routing it to the right downstream system,
running stateful workflows that survive process restarts, letting
operators build new integrations without redeploying the UI, and doing
all of the above for many tenants on a single stack.

MapexOS treats those four problems as first-class — not afterthoughts
bolted onto a time-series database. The result is a platform you can
actually run a fleet on, not just a demo with a temperature widget.

---

## License

MapexOS is distributed under the **Business Source License 1.1** —
see [`LICENSE`](./LICENSE) for the full terms.

In short:

- **You can** self-host, modify, and run MapexOS in production for
  your own organization.
- **You cannot** offer MapexOS as a hosted commercial service to
  third parties without a separate commercial agreement.
- On the **Change Date**, the license converts to **Apache 2.0**.

For commercial licensing inquiries, contact **Mapex Solutions**.

# MapexOS E2E Tests

End-to-end test suite for the MapexOS platform. Two kinds of tests live
here:

- **Module e2e** (`services/<svc>/<module>/`) — CRUD and contract tests
  scoped to a single module. Every endpoint of each module is exercised
  with real fixtures against the running stack.
- **Saga journeys** (`journey/<context>/<journey_name>/phaseN_*/`) —
  cross-module flows organised as ordered phases. Each phase is a saga
  with steps + asserts + rollback. Gated by the `saga` build tag.

The whole suite runs against the canonical stack from
[`mapexOSDeploy`](../../mapexOSDeploy/) and uses the seed admin user
(`admin@mapex.local` / `mapex@123`) as the bootstrap actor.

## Prerequisites

1. Stack up via the canonical docker-compose:
   ```bash
   cd ../mapexOSDeploy
   docker compose up -d
   ```
2. The Go services running and listening on the default ports — see
   [Service ports](#service-ports) below.
3. Go 1.25+ on the host running the tests.

The seed admin user, role, organization, and recursive membership are
provisioned by the `mongodb-init` container on first boot of mongo.
Tests that need additional users provision them at runtime via the
public onboarding orchestrator (`POST /api/v1/onboarding/users`) — the
seed is for real/production data only and is never modified by tests.

## Layout

```
e2e_tests/
├── common/                     # Shared test code
│   ├── constants/              # URLs, ids, credentials
│   ├── types/                  # StandardResponse, ErrorResponse, ...
│   └── utils/                  # Login, assertions, setup
│
├── services/                   # Module e2e — one test package per module
│   ├── assets/{assets,assettemplates}/
│   ├── events/events/
│   ├── http_gateway/datasources/
│   ├── mapexIam/{auth,onboarding,organizations,roles}/
│   ├── mapexos/{auth,groups,lists,memberships,organizations,roles,users}/
│   ├── router/routegroups/
│   ├── triggers/triggers/
│   └── workflow/{definitions,instances}/
│
└── journey/                    # Saga journeys (build tag: saga)
    ├── automations/            # 8 trigger journeys (http, email, websocket,
    │   │                       #   slack, teams, mqtt, nats, rabbitmq)
    │   └── trigger_<type>/{phase1_connectivity,phase2_event_pipeline}/
    └── iot/
        ├── connectivity_actions_{http,mqtt}/{phase1_workflow,phase2_trigger}/
        └── mqtt_broker_auth/{phase0_iam_bootstrap..phase3_cascade}/
```

## How to run

Every command runs from `e2e_tests/`.

```bash
cd e2e_tests

# All module e2e tests (saga tag NOT set)
go test ./services/...

# A single module
go test ./services/mapexos/organizations -v

# A single test
go test ./services/mapexos/organizations -v -run TestCreateOrganization_Customer

# All saga journeys (saga tag REQUIRED)
go test -tags=saga ./journey/...

# A single journey context / journey / phase
go test -tags=saga ./journey/automations/...
go test -tags=saga ./journey/automations/trigger_http/...
go test -tags=saga ./journey/automations/trigger_http/phase1_connectivity
```

For longer-running suites, set `-timeout`:

```bash
go test -tags=saga -timeout 15m ./journey/...
```

## Service ports

| Service        | Port  | Required for                       |
|----------------|-------|------------------------------------|
| mapexos / iam  | 5000  | All tests (auth + org CRUD)        |
| http_gateway   | 5001  | datasources tests, saga phase 2    |
| assets         | 5002  | assets / assettemplates / saga IoT |
| router         | 5003  | routegroups + saga                 |
| events         | 5004  | event-pipeline saga phases         |
| triggers       | 5006  | trigger journeys                   |
| workflow       | 5007  | workflow tests + IoT actions       |

Saga journeys also bind in-process sinks on the host: `11010` (HTTP),
`11025` (SMTP), `11026` (WebSocket). Make sure those ports are free.

## Environment overrides

The defaults match the canonical stack. Override only if your stack is
on different hosts/ports:

```bash
export MAPEXOS_URL=http://localhost:5000
export GATEWAY_URL=http://localhost:5001
export ASSETS_URL=http://localhost:5002
export ROUTER_URL=http://localhost:5003
export EVENTS_URL=http://localhost:5004
export TRIGGERS_URL=http://localhost:5006
export WORKFLOW_URL=http://localhost:5007
```

## Conventions

- **Tests use only the public API.** Internal routes (`/internal/*`) are
  cache-rebuild fallbacks and are never invoked by tests.
- **The seed admin is the bootstrap actor.** Any additional users a test
  needs are provisioned at runtime through the onboarding orchestrator
  (`POST /api/v1/onboarding/users`). The seed JSON is never modified by
  tests.
- **Fixtures live next to the test** in a `fixtures/` folder and use the
  canonical seed ids (`0000000000000000000aa001` for the root org,
  `0000000000000000000aa201` for the SuperAdmin role).
- **Cleanup is mandatory.** Every mutating test registers `t.Cleanup` or
  `defer` to delete what it created.

## Documentation

Each module e2e package and each saga journey carries a
`README.md` + `README_pt.md` describing scope, fixtures, and how to run
it standalone. Start at the directory you care about — the per-folder
READMEs are the source of truth for each test surface.

See also: [`journey/README.md`](./journey/README.md) for the saga
journey hierarchy rules.

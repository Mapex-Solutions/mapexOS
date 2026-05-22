# Journey hierarchy

Tests scale into the thousands as the platform grows. To keep them findable
and to keep cross-journey reuse honest, the directory tree is organised as:

```
journey/
├── README.md                     # this file (global hierarchy + rules)
└── {context}/                    # business domain (iot, workflow, vault, platform, ...)
    ├── README.md                 # context narrative + run commands for every journey here
    └── {journey_name}/           # one named end-to-end story (snake_case)
        ├── README.md             # what the journey covers + outcome at a glance
        └── {phaseN_descriptor}/  # one phase folder per stage of the journey
            ├── journey.go        # journey-local helpers + Run + (optional) ItemsForCompose
            └── journey_test.go   # //go:build saga; one TestPhaseN_<Descriptor>_Saga func
```

## Rules

- **Context** groups journeys by business domain. Examples: `iot/`, `workflow/`, `vault/`, `iam/`, `platform/` (cross-cutting). Add a new context only when an existing one would mislead a reader looking for the journey.

- **Journey** is the named story (snake_case). One journey = one cohesive
  end-to-end narrative. Examples: `mqtt_full_pipeline`, `http_heartbeat`,
  `credential_refresh`, `dlq_collects_failed_messages`.

- **Phase** is a named stage of the journey. Phases run in order (Phase 0
  before Phase 1, etc). Each phase folder is a Go package; package name is
  short (`phase0`, `phase1`) so consumers alias as `phase0 ".../phase0_iam_bootstrap"`.
  Phase folders carry a descriptor suffix (`phase0_iam_bootstrap`) so the
  folder name alone tells you what the phase does.

- **Cross-journey reuse** never happens by importing another journey's
  phase. Reuse lives in `services/{svc}/{mod}/{steps,asserts,payloads}` —
  the building blocks every phase composes. Within a single journey,
  phases CAN import each other (PhaseN+1 typically composes
  PhaseN.BootstrapItems).

- **Each phase MUST carry an Outcome block** in both `journey.go` (package
  godoc) and `journey_test.go` (test func godoc) describing what passing
  the phase proves and what a failure typically points at. The Outcome
  blocks are how a dev figures out what a journey covers without opening
  source files.

- **Every context README MUST carry a "How to run" block** with the
  go test commands scoped at the context, journey, and phase level.
  Devs land on a context, see the run commands, and run the relevant
  scope without searching docs.

- **Documentation never bakes in a developer-specific absolute path.**
  Run commands assume the working directory is the repository root and
  use relative paths (e.g. `cd workspace_go/packages/e2eTests`). Anyone
  who clones the monorepo can copy-paste without editing.

## Currently registered journeys

| Context | Journey                    | Phases                                  |
|---------|----------------------------|-----------------------------------------|
| iot     | mqtt_full_pipeline         | phase0_iam_bootstrap, phase1_iot_setup  |

## How to run

Every command runs from the e2eTests package root.

```bash
cd workspace_go/packages/e2eTests

# All journeys in every context
go test -tags=saga -v ./journey/...

# All journeys in one context (see the per-context README for finer scopes)
go test -tags=saga -v ./journey/iot/...
```

The `saga` build tag gates these tests: `go test ./...` (no tag) skips
them; only `go test -tags=saga` walks the journey folders.

## How to add a new journey

1. Pick the context (or create one if no existing context fits).
2. If creating a context, add `journey/{context}/README.md` with the
   narrative, a registered-journeys table, and a "How to run" block.
3. Create `journey/{context}/{journey_name}/README.md` with the
   narrative and a phase index.
4. Create `phaseN_<descriptor>/journey.go` and `journey_test.go`. Copy
   the layout from an existing phase for consistency.
5. Update this README's registered-journeys table and the context
   README's table.
6. PhaseN reuses Phase0..N-1 by importing them directly within the same
   journey folder. Cross-journey reuse only at the building-block level.

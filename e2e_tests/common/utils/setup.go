package utils

import (
	"sync"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/infra"
)

var setupOnce sync.Once

// SetupE2EEnvironment ensures the whole stack is provisioned before a test package
// runs. It delegates to infra.EnsureAll — the same two-stack (infra + mapex),
// reuse-or-bring-up path the saga suite uses in its TestMain — so EVERY e2e test
// (module e2e AND saga journeys) provisions the environment the same way, whether you
// run a single test or all of them. Idempotent (sync.Once per package): a service
// already answering its health probe is reused.
//
// The teardown EnsureAll returns is intentionally NOT run here: a module e2e package
// leaves the stack up for reuse across runs (the common local flow), and the saga
// runner tears down only what it started. Modes are env-driven (MAPEX_E2E_INFRA_MODE /
// MAPEX_E2E_MAPEX_MODE); the default is infra=docker, mapex=local. It never errors
// (EnsureAll is non-fatal and logs a per-stack summary); a genuinely-down dependency
// surfaces as the first failing HTTP call, with EnsureAll's "start it" log above it.
func SetupE2EEnvironment() error {
	setupOnce.Do(func() {
		_ = infra.EnsureAll()
	})
	return nil
}

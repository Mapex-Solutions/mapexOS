package di

import "go.uber.org/dig"

// EngineServiceDependenciesInjection aggregates all dependencies for EngineService.
// Engine module is pure computation — no external dependencies (repos, NATS, metrics).
type EngineServiceDependenciesInjection struct {
	dig.In
}

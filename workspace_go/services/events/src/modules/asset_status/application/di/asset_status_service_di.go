package di

import (
	"events/src/modules/asset_status/domain/repositories"

	"go.uber.org/dig"
)

// AssetStatusServiceDependenciesInjection aggregates the dependencies
// required by the AssetStatusService. The dig.In tag enables automatic
// resolution by the DIG container at InitServices phase.
//
// Dependencies are declared as ports (interfaces), never concrete types —
// Hexagonal boundary preserved.
type AssetStatusServiceDependenciesInjection struct {
	dig.In

	// Repository persists + queries asset connectivity events in ClickHouse.
	Repository repositories.AssetStatusRepository
}

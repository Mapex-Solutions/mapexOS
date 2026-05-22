package ports

import (
	"assets/src/modules/assettemplates/domain/entities"
	"assets/src/modules/assettemplates/domain/repositories"
)

// Port-level type aliases — expose domain types through the port boundary.
// Other modules import these types from ports, NEVER from domain/ directly.

type AssetTemplateRepository = repositories.AssetTemplateRepository

// Assettemplate is the cross-module alias for the domain entity. External
// modules (e.g. assets) MUST import this alias instead of reaching into
// assettemplates/domain/entities directly.
type Assettemplate = entities.Assettemplate

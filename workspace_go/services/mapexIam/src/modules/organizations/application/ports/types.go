package ports

import (
	"mapexIam/src/modules/organizations/domain/entities"
)

// Organization is the public type alias for the domain entity.
// Cross-module consumers import this instead of domain/entities directly.
type Organization = entities.Organization

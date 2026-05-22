package ports

import (
	"mapexIam/src/modules/memberships/domain/entities"
)

// Membership is the public type alias for the domain entity.
// Cross-module consumers import this instead of domain/entities directly.
type Membership = entities.Membership

package ports

import (
	"mapexIam/src/modules/users/domain/entities"
)

// User is the public type alias for the domain entity.
// Cross-module consumers import this instead of domain/entities directly.
type User = entities.User

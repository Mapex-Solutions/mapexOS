package services

import (
	"mapexIam/src/modules/cache_invalidation/application/di"
)

// CacheInvalidationService orchestrates reaction to IAM domain change events and invalidates
// authorization/coverage caches. It is a pure application-layer service; message Ack/Nack
// is done here based on handler outcome (the consumer layer is only wiring).
type CacheInvalidationService struct {
	deps di.CacheInvalidationServiceDependenciesInjection
}

// UserOrgPair represents a unique user + organization combination for cache invalidation.
// Used as map key to deduplicate pairs when resolving memberships (user and group).
type UserOrgPair struct {
	UserID string
	OrgID  string
}

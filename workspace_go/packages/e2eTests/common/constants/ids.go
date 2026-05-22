package constants

// Fixed ObjectIDs for deterministic E2E tests.
// These IDs match the canonical mongodb-init seed shipped with the
// services_required and standalone/infra docker-compose stacks
// (deployment/docker-compose/.../mongodb/seed/mapex_iam/*.json).

// Organization IDs.
const (
	// MapexosOrgID is the seed root organization (vendor type, depth=0).
	// Sagas use it as parent when creating scratch tenant orgs.
	MapexosOrgID = "0000000000000000000aa001"
)

// Role IDs.
const (
	// SuperAdminRoleID is the role bound to the seed admin user. Carries
	// the wildcard "mapex.*" permission so every endpoint accepts it.
	SuperAdminRoleID = "0000000000000000000aa201"

	// Legacy aliases preserved so older tests in this package keep
	// compiling. They all point at the same seeded SuperAdmin role —
	// finer-grained roles are no longer part of the default seed.
	RootRoleID        = SuperAdminRoleID
	VendorAdminRoleID = SuperAdminRoleID
	ViewerRoleID      = SuperAdminRoleID
)

// User IDs.
const (
	// AdminUserID is the seed admin user used by GetRootToken when a saga
	// needs an authenticated bootstrap actor before the test user it
	// provisions takes over.
	AdminUserID = "0000000000000000000aa101"

	// RootUserID is preserved as an alias for AdminUserID for legacy
	// tests; the seed only ships one super-admin user.
	RootUserID = AdminUserID
)

// User credentials.
//
// The legacy "Root" prefix is preserved so journeys and helpers that
// already call utils.GetRootToken keep working; the values point at the
// admin@mapex.local user the seed actually creates. Sagas that need a
// pre-authenticated bootstrap actor reuse this single seeded user.
const (
	// RootUserEmail / RootUserPassword authenticate the seed admin user.
	// Bound to SuperAdminRoleID inside MapexosOrgID with scope=recursive.
	RootUserEmail    = "admin@mapex.local"
	RootUserPassword = "mapex@123"

	// AdminUserEmail / AdminUserPassword alias the same seed admin user.
	// Kept for backward compatibility with legacy tests that reach for
	// "Admin*" constants explicitly.
	AdminUserEmail    = "admin@mapex.local"
	AdminUserPassword = "mapex@123"
)

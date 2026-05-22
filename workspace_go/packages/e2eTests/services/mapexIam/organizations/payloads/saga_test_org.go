// Package payloads holds canonical OrganizationCreate fixtures for the
// IAM organizations saga module.
//
// Each function returns a fluent builder that wraps the contract DTO
// straight from packages/contracts/services/mapexIam/organizations so the
// payload stays in lockstep with the API contract — never duplicating the
// struct shape. Builders never mutate their canonical defaults; each call
// returns a fresh copy ready for chaining With* overrides.
package payloads

import (
	"fmt"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/organizations"

	"github.com/Mapex-Solutions/mapexGoKit/utils/zerovalue"
)

// SagaTestCustomerOrgParentID is the deterministic seed root org under which
// every saga creates its scratch customer org. Mirrors
// constants.MapexosOrgID so payloads stay independent of the constants
// package; using a stable parent keeps hierarchy assertions predictable
// across runs.
const SagaTestCustomerOrgParentID = "0000000000000000000aa001"

// OrganizationCreateBuilder wraps contracts.OrganizationCreate so tests can
// override individual fields without redeclaring the canonical baseline.
type OrganizationCreateBuilder struct {
	spec contracts.OrganizationCreate
}

// Build returns the contracts payload ready to send to
// POST /api/v1/organizations.
func (b *OrganizationCreateBuilder) Build() contracts.OrganizationCreate {
	return b.spec
}

// WithName overrides the canonical name. Used when a journey needs a stable
// identifier the test asserts on.
func (b *OrganizationCreateBuilder) WithName(name string) *OrganizationCreateBuilder {
	b.spec.Name = name
	return b
}

// WithParentOrgID overrides the parent. Useful when a saga needs to create
// the org under a different baseline parent than the canonical seed root.
func (b *OrganizationCreateBuilder) WithParentOrgID(id string) *OrganizationCreateBuilder {
	b.spec.ParentOrgID = zerovalue.Ptr(id)
	return b
}

// SagaTestCustomerOrg returns the canonical fixture for a customer-type
// organization parented at the seed root. The name carries runID so cleanup
// by prefix and parallel runs do not collide.
//
// Defaults:
//   - Type:         customer
//   - ParentOrgID:  SagaTestCustomerOrgParentID
//   - Enabled:      true
//   - AuthConfig:   internal provider (no external IDP)
//   - AccessPolicy: strict role policy + local default scope
func SagaTestCustomerOrg(runID string) *OrganizationCreateBuilder {
	return &OrganizationCreateBuilder{
		spec: contracts.OrganizationCreate{
			Name:        fmt.Sprintf("saga-customer-%s", runID),
			Type:        "customer",
			ParentOrgID: zerovalue.Ptr(SagaTestCustomerOrgParentID),
			Enabled:     true,
			AuthConfig: contracts.AuthConfig{
				ProviderType: "internal",
			},
			AccessPolicy: contracts.AccessPolicy{
				RolePolicy:   "strict",
				DefaultScope: "local",
			},
		},
	}
}

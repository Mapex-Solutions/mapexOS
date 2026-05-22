// Package payloads holds canonical CreateUserWithMemberships fixtures for
// the IAM onboarding orchestrator module.
//
// The onboarding orchestrator endpoint creates user + group + membership in
// a single atomic call, removing the need for the saga to coordinate three
// separate POSTs. Saga journeys reuse SagaIoTAdminUser as the canonical
// actor authenticated against the freshly created scratch org.
package payloads

import (
	"fmt"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/onboarding"

	"github.com/Mapex-Solutions/mapexGoKit/utils/zerovalue"
)

// SagaIoTAdminUserPassword is the deterministic password assigned to every
// saga test user. Hard-coded here so the auth step can sign in without
// passing the secret through the bag.
const SagaIoTAdminUserPassword = "saga-test-password-1234"

// CreateUserWithMembershipsBuilder wraps the contract DTO with fluent
// overrides callers use when a journey needs a tailored variant.
type CreateUserWithMembershipsBuilder struct {
	spec contracts.CreateUserWithMemberships
}

// Build returns the contracts payload ready for
// POST /api/v1/onboarding/users.
func (b *CreateUserWithMembershipsBuilder) Build() contracts.CreateUserWithMemberships {
	return b.spec
}

// WithEmail overrides the email address; useful for negative cases that
// assert duplicate-email handling.
func (b *CreateUserWithMembershipsBuilder) WithEmail(email string) *CreateUserWithMembershipsBuilder {
	b.spec.Email = email
	return b
}

// SagaIoTAdminUser returns the canonical onboarding payload that creates a
// new user inside a brand-new group bound to the saga IoT admin role. The
// runID-stamped email keeps the user unique across parallel runs while the
// password is the shared SagaIoTAdminUserPassword every saga uses to sign
// in.
//
// Inputs:
//   - runID    journey identifier; embedded into email and group name
//   - roleID   IAM role id created earlier in the saga (bag: iam.roleID)
func SagaIoTAdminUser(runID string, roleID string) *CreateUserWithMembershipsBuilder {
	pwd := SagaIoTAdminUserPassword
	return &CreateUserWithMembershipsBuilder{
		spec: contracts.CreateUserWithMemberships{
			Email:                   fmt.Sprintf("saga-iot-admin-%s@test.local", runID),
			Password:                zerovalue.Ptr(pwd),
			ChangePasswordNextLogin: false,
			FirstName:               "Saga",
			LastName:                "Admin",
			Enabled:                 true,
			Groups: []contracts.GroupAccessData{
				{
					NewGroup: &contracts.NewGroupData{
						Name:    fmt.Sprintf("saga-iot-admin-group-%s", runID),
						RoleIds: []string{roleID},
					},
				},
			},
		},
	}
}

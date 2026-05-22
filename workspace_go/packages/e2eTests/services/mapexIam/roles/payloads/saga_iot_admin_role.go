// Package payloads holds canonical RoleCreate fixtures for the IAM roles
// saga module.
//
// The fixtures bundle the permissions needed to drive the IoT pipeline from
// a single role: org/role/group/user CRUD, asset and template CRUD, route
// group + trigger + workflow CRUD, plus events read for the verification
// step. Every saga that needs an end-user actor in IoT scenarios should
// reuse SagaIoTAdminRole rather than crafting a custom permission list.
package payloads

import (
	"fmt"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/roles"

	assetsPerms "github.com/Mapex-Solutions/MapexOS/permissions/assets"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/mapexos"
	routerPerms "github.com/Mapex-Solutions/MapexOS/permissions/router"
)

// RoleCreateBuilder wraps contracts.RoleCreate with fluent overrides.
type RoleCreateBuilder struct {
	spec contracts.RoleCreate
}

// Build returns the contracts payload ready to send to POST /api/v1/roles.
func (b *RoleCreateBuilder) Build() contracts.RoleCreate { return b.spec }

// WithName overrides the role name. Tests that assert by name use a
// deterministic override; otherwise the runID-stamped default suffices.
func (b *RoleCreateBuilder) WithName(name string) *RoleCreateBuilder {
	b.spec.Name = name
	return b
}

// SagaIoTAdminRole returns the canonical IoT-admin role for a saga test
// user: every permission needed to bootstrap the IoT pipeline end to end.
//
// Defaults:
//   - Scope:      local (org-scoped; user only acts inside the saga org)
//   - IsSystem:   false (regular role, attached to the saga org)
//   - IsTemplate: false
//   - Permissions: org list/read, role list/read, user list/read, group
//     list/read, asset/template/route/trigger/workflow CRUD,
//     events read.
func SagaIoTAdminRole(runID string) *RoleCreateBuilder {
	return &RoleCreateBuilder{
		spec: contracts.RoleCreate{
			Name:       fmt.Sprintf("saga-iot-admin-%s", runID),
			Scope:      "local",
			IsSystem:   false,
			IsTemplate: false,
			Permissions: []string{
				// IAM read paths the saga test user navigates after login.
				perms.OrganizationList, perms.OrganizationRead,
				perms.RoleList, perms.RoleRead,
				perms.UserCreate, perms.UserList, perms.UserRead, perms.UserUpdate,
				perms.GroupList, perms.GroupRead,

				// Asset templates and assets — full CRUD so subsequent
				// phases can provision and tear down IoT entities under
				// the saga org without elevating to ROOT.
				assetsPerms.AssetTemplateAll,
				assetsPerms.AssetAll,

				// Route groups — required by every asset (RouteGroupIds
				// is mandatory) and by router-side configuration tests.
				routerPerms.RouteGroupAll,
			},
		},
	}
}

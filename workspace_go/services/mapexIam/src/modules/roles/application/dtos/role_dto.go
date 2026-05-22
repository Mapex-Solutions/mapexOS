package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/roles"
)

type (
	RoleIdDto      = v1.RoleId
	CreateRoleDto  = v1.RoleCreate
	UpdateRoleDto  = v1.RoleUpdate
	RoleQueryDto   = v1.RoleQuery
	RoleResponse   = v1.RoleResponse
)

package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/groups"
)

type (
	GroupIdDto           = v1.GroupId
	CreateGroupDto       = v1.GroupCreate
	UpdateGroupDto       = v1.GroupUpdate
	GroupQueryDto        = v1.GroupQuery
	GroupResponse        = v1.GroupResponse
	GroupMemberResponse  = v1.GroupMemberResponse
	GroupMembersQueryDto = v1.GroupMembersQuery
	GroupMemberAddDto    = v1.GroupMemberAddDto
	GroupMemberIdDto     = v1.GroupMemberIdDto
)

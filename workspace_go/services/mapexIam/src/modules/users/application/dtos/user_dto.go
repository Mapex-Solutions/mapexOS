package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/users"
)

type (
	UserCreateDTO      = v1.UserCreate
	UserUpdateDTO      = v1.UserUpdate
	UserResponse       = v1.UserResponse
	UserIdDTO          = v1.UserId
	UserQueryDto       = v1.UserQuery
	UserGroupInfo      = v1.UserGroupInfo
	UserMembershipInfo = v1.UserMembershipInfo
)

package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/onboarding"
)

type (
	MembershipData                  = v1.MembershipData
	ExistingGroupData               = v1.ExistingGroupData
	NewGroupData                    = v1.NewGroupData
	GroupAccessData                 = v1.GroupAccessData
	CreateUserWithMembershipsDto    = v1.CreateUserWithMemberships
	UpdateUserWithAccessParamsDto   = v1.UpdateUserWithAccessParams
	UpdateUserWithAccessDto         = v1.UpdateUserWithAccess
	UserOnboardingResponse          = v1.UserOnboardingResponse
)

package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/memberships"
)

type (
	MembershipIdDto       = v1.MembershipId
	CreateMembershipDto   = v1.MembershipCreate
	UpdateMembershipDto   = v1.MembershipUpdate
	MembershipQueryDto    = v1.MembershipQuery
	MembershipResponse    = v1.MembershipResponse
	CustomerCoverage      = v1.CustomerCoverage
	MeCoverageResponse    = v1.MeCoverageResponse
)

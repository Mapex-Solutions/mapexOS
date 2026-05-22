package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/events/retention"
)

type (
	RetentionPolicyUpsertDTO  = v1.RetentionPolicyUpsert
	RetentionPolicyQueryDTO   = v1.RetentionPolicyQuery
	RetentionPolicyResponse   = v1.RetentionPolicyResponse
	RetentionPolicyParamsDTO  = v1.RetentionPolicyParams
	RetentionPolicyLimits     = v1.RetentionPolicyLimits
)

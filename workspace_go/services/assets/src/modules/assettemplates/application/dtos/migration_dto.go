package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets_templates"

	query "github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

// Contract DTOs are canonical in packages/contracts; re-aliased here so the
// module depends on the shared contract, never a local redefinition.
type (
	MigrationPlanCreateRequest = v1.MigrationPlanCreateRequest
	MigrationPlanUpdateRequest = v1.MigrationPlanUpdateRequest
	MigrationPlanResponse      = v1.MigrationPlanResponse
	MigrationExecutionResponse = v1.MigrationExecutionResponse
)

// MigrationPlanQueryDTO is the query for the plans /find endpoint.
type MigrationPlanQueryDTO struct {
	query.BaseQueryDTO
	Name   *string `query:"name"`
	Status *string `query:"status"`
}

// MigrationExecutionQueryDTO is the query for a plan's executions /find endpoint.
type MigrationExecutionQueryDTO struct {
	query.BaseQueryDTO
	Status  *string `query:"status"`
	AssetID *string `query:"assetId"`
}

// MigrationPlanIdDto is the :id path parameter.
type MigrationPlanIdDto struct {
	Id string `params:"id" validate:"required,mongoid"`
}

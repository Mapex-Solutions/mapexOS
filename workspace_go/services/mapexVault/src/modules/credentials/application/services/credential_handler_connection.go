package services

import (
	"fmt"
	"time"

	"mapexVault/src/modules/credentials/application/dtos"
	"mapexVault/src/modules/credentials/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

// buildConnectionEntity assembles a Connection entity for Create from the
// inbound DTO, applying org/path scoping from the request context.
func (s *CredentialService) buildConnectionEntity(rc *reqCtx.RequestContext, dto *dtos.CreateConnectionDTO) (*entities.Connection, error) {
	credObjId, err := model.ToObjectID(dto.CredentialId)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Invalid credentialId: %w", err)
	}
	entity := &entities.Connection{
		Provider:     dto.Provider,
		AccountId:    dto.AccountId,
		AccountName:  dto.AccountName,
		Status:       entities.ConnectionStatusActive,
		CredentialId: credObjId,
		Scopes:       dto.Scopes,
		ConnectedAt:  time.Now(),
		Created:      time.Now(),
		Updated:      time.Now(),
	}
	if rc.OrgContext != nil && *rc.OrgContext != "" {
		if orgObjectId, err := model.ToObjectID(*rc.OrgContext); err == nil {
			entity.OrgId = &orgObjectId
		}
	}
	if rc.OrgContextData != nil {
		entity.PathKey = rc.OrgContextData.PathKey
	}
	return entity, nil
}

// buildConnectionListFilters builds the Mongo filter for GetConnections.
// Org scope comes from the request context; provider/status come from query.
func (s *CredentialService) buildConnectionListFilters(rc *reqCtx.RequestContext, query *dtos.ConnectionQueryDTO) model.Map {
	filters := model.Map{}
	if rc.OrgContext != nil && *rc.OrgContext != "" {
		if orgId, err := model.ToObjectID(*rc.OrgContext); err == nil {
			filters["orgId"] = orgId
		}
	}
	if query.Provider != nil {
		filters["provider"] = *query.Provider
	}
	if query.Status != nil {
		filters["status"] = *query.Status
	}
	return filters
}

// buildConnectionListPagination derives pagination opts from the query DTO.
func (s *CredentialService) buildConnectionListPagination(query *dtos.ConnectionQueryDTO) *model.PaginationOpts {
	var page, perPage int64 = 1, 20
	if query.Page != nil {
		page = int64(*query.Page)
	}
	if query.PerPage != nil {
		perPage = int64(*query.PerPage)
	}
	return &model.PaginationOpts{Page: page, PerPage: perPage}
}

// mapConnectionList converts a paginated entity result into the response DTO.
func (s *CredentialService) mapConnectionList(result *model.PaginatedResult[entities.Connection]) *model.PaginatedResult[dtos.ConnectionResponse] {
	responses := make([]dtos.ConnectionResponse, len(result.Items))
	for i, conn := range result.Items {
		responses[i] = *toConnectionResponse(&conn)
	}
	return &model.PaginatedResult[dtos.ConnectionResponse]{
		Items:      responses,
		Pagination: result.Pagination,
	}
}

// buildUpsertConnectionEntity builds the Connection entity for the upsert
// flow and returns the org id pointer used as part of the upsert key.
func (s *CredentialService) buildUpsertConnectionEntity(rc *reqCtx.RequestContext, dto *dtos.UpsertConnectionDTO) (*entities.Connection, *model.ObjectId, error) {
	credObjId, err := model.ToObjectID(dto.CredentialId)
	if err != nil {
		return nil, nil, fmt.Errorf("[SERVICE:Credential] Invalid credentialId: %w", err)
	}
	var orgId *model.ObjectId
	pathKey := ""
	if rc.OrgContext != nil && *rc.OrgContext != "" {
		if oid, err := model.ToObjectID(*rc.OrgContext); err == nil {
			orgId = &oid
		}
	}
	if rc.OrgContextData != nil {
		pathKey = rc.OrgContextData.PathKey
	}
	entity := &entities.Connection{
		AccountName:  dto.AccountName,
		Status:       entities.ConnectionStatusActive,
		CredentialId: credObjId,
		PathKey:      pathKey,
		Scopes:       dto.Scopes,
		ConnectedAt:  time.Now(),
		Created:      time.Now(),
		Updated:      time.Now(),
	}
	return entity, orgId, nil
}

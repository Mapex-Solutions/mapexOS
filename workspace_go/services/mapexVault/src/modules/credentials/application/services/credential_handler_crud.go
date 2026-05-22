package services

import (
	"context"
	"fmt"
	"time"

	"mapexVault/src/modules/credentials/application/dtos"
	"mapexVault/src/modules/credentials/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/utils/envelope"
)

// buildCredentialEntity assembles the Credential entity for Create from the
// inbound DTO + envelope-encrypted blob, applying org/path scoping from the
// request context.
func (s *CredentialService) buildCredentialEntity(rc *reqCtx.RequestContext, dto *dtos.CreateCredentialDTO, env *envelope.EncryptedEnvelope) *entities.Credential {
	entity := &entities.Credential{
		Name:            dto.Name,
		Type:            dto.Type,
		PluginId:        dto.PluginId,
		CredentialDefId: dto.CredentialDefId,
		IsTemplate:      dto.IsTemplate,
		Status:          entities.CredentialStatusActive,
		EncryptedDEK:    env.EncryptedDEK,
		DEKNonce:        env.DEKNonce,
		EncryptedData:   env.EncryptedData,
		DataNonce:       env.DataNonce,
		ProviderConfig:  dto.ProviderConfig,
		Created:         time.Now(),
		Updated:         time.Now(),
	}
	if rc.OrgContext != nil && *rc.OrgContext != "" {
		if orgObjectId, err := model.ToObjectID(*rc.OrgContext); err == nil {
			entity.OrgId = &orgObjectId
		}
	}
	if rc.OrgContextData != nil {
		entity.PathKey = rc.OrgContextData.PathKey
	}
	return entity
}

// buildCredentialUpdateMap composes the $set map for UpdateCredentialById,
// re-running envelope encryption only when the patch carries new secrets.
func (s *CredentialService) buildCredentialUpdateMap(dto *dtos.UpdateCredentialDTO) (model.Map, error) {
	update := model.Map{"updated": time.Now()}
	if dto.Name != nil {
		update["name"] = *dto.Name
	}
	if dto.IsTemplate != nil {
		update["isTemplate"] = *dto.IsTemplate
	}
	if dto.ProviderConfig != nil {
		update["providerConfig"] = dto.ProviderConfig
	}
	if dto.Data != nil {
		env, err := encryptData(s.deps.Encryption, dto.Data)
		if err != nil {
			return nil, err
		}
		update["encryptedDEK"] = env.EncryptedDEK
		update["dekNonce"] = env.DEKNonce
		update["encryptedData"] = env.EncryptedData
		update["dataNonce"] = env.DataNonce
	}
	return update, nil
}

// buildCredentialListFilters builds the Mongo filter for GetCredentials.
// Org scope comes from the request context; pluginId/type/status come from
// the query DTO.
func (s *CredentialService) buildCredentialListFilters(rc *reqCtx.RequestContext, query *dtos.CredentialQueryDTO) model.Map {
	filters := model.Map{}
	if rc.OrgContext != nil && *rc.OrgContext != "" {
		if orgId, err := model.ToObjectID(*rc.OrgContext); err == nil {
			filters["orgId"] = orgId
		}
	}
	if query.PluginId != nil {
		filters["pluginId"] = *query.PluginId
	}
	if query.Type != nil {
		filters["type"] = *query.Type
	}
	if query.Status != nil {
		filters["status"] = *query.Status
	}
	return filters
}

// buildCredentialListPagination derives the Mongo pagination opts from the
// query DTO using the shared per-page defaults.
func (s *CredentialService) buildCredentialListPagination(query *dtos.CredentialQueryDTO) *model.PaginationOpts {
	var page, perPage int64 = 1, 20
	if query.Page != nil {
		page = int64(*query.Page)
	}
	if query.PerPage != nil {
		perPage = int64(*query.PerPage)
	}
	return &model.PaginationOpts{Page: page, PerPage: perPage}
}

// mapCredentialList converts a paginated entity result into the response
// DTO, stripping the encrypted blobs at the boundary.
func (s *CredentialService) mapCredentialList(result *model.PaginatedResult[entities.Credential]) *model.PaginatedResult[dtos.CredentialResponse] {
	responses := make([]dtos.CredentialResponse, len(result.Items))
	for i, cred := range result.Items {
		responses[i] = *toCredentialResponse(&cred)
	}
	return &model.PaginatedResult[dtos.CredentialResponse]{
		Items:      responses,
		Pagination: result.Pagination,
	}
}

// fetchCredentialOrError loads a credential and surfaces a contextualized
// not-found error so callers can short-circuit cleanly. The op label keeps
// log lines in sync with the calling orchestration.
func (s *CredentialService) fetchCredentialOrError(ctx context.Context, id, op string) (*entities.Credential, error) {
	cred, err := s.deps.CredentialRepo.FindById(ctx, &id)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Not found for %s: %w", op, err)
	}
	if cred == nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Credential %s not found", id)
	}
	return cred, nil
}

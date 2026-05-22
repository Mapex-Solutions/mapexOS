package services

import (
	ctx "context"
	"time"

	"http_gateway/src/modules/datasources/application/dtos"
	"http_gateway/src/modules/datasources/domain/entities"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// applyCreateRequestContext populates the org-related fields of the create DTO
// (OrgID and PathKey) from the caller RequestContext when available.
func (s *DataSourceService) applyCreateRequestContext(requestContext *reqCtx.RequestContext, dto *dtos.DataSourceCreateDTO) {
	if requestContext.OrgContext != nil && *requestContext.OrgContext != "" {
		if orgObjectId, err := model.ToObjectID(*requestContext.OrgContext); err == nil {
			dto.OrgID = &orgObjectId
		}
	}

	if requestContext.OrgContextData != nil && requestContext.OrgContextData.PathKey != "" {
		pathKey := requestContext.OrgContextData.PathKey
		dto.PathKey = &pathKey
	}
}

// persistNewDataSource maps the create DTO into an entity and delegates creation to the repository.
func (s *DataSourceService) persistNewDataSource(c ctx.Context, dto *dtos.DataSourceCreateDTO) (*entities.DataSource, error) {
	entity, _ := mapper.DtoToEntity[dtos.DataSourceCreateDTO, entities.DataSource](dto)
	return s.deps.DataSourceRepo.Create(c, entity)
}

// buildUpdateFields converts the update DTO into a map of fields to update and stamps "updated".
func (s *DataSourceService) buildUpdateFields(dto *dtos.DataSourceUpdateDTO) map[string]any {
	fieldsToUpdate, _ := mapper.DtoToMap(dto)
	fieldsToUpdate["updated"] = time.Now()
	return fieldsToUpdate
}

// applyDataSourceUpdate invokes the repository update and returns the refreshed entity.
// Preserves legacy behavior of swallowing repository errors — callers rely on ID.IsZero()
// to detect not-found results rather than on the error value.
func (s *DataSourceService) applyDataSourceUpdate(
	c ctx.Context,
	dataSourceId *string,
	fieldsToUpdate map[string]any,
) *entities.DataSource {
	entity, _ := s.deps.DataSourceRepo.FindByIdAndUpdate(c, dataSourceId, fieldsToUpdate)
	return entity
}

// deleteDataSourceCache removes the cached entry for the given DataSource id.
func (s *DataSourceService) deleteDataSourceCache(c ctx.Context, dataSourceId *string) {
	cacheKey := s.deps.CacheKeyBuilder.BuildKey(*dataSourceId)
	s.deps.CacheRepo.Del(c, cacheKey)
}

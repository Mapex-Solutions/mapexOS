package services

import (
	ctx "context"
	"strings"

	"mapexIam/src/modules/lists/application/dtos"
	"mapexIam/src/modules/lists/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// findResolvedListId returns the id of the org-scoped list matching
// (type, value=slug, orgId[, parentId]), or "" when none exists. Filtering by a
// concrete orgId inherently excludes global/system rows (they carry no orgId), so
// a resolve never returns a shared/system list.
func (s *ListService) findResolvedListId(c ctx.Context, req dtos.ListResolveRequest) (string, error) {
	orgID, err := model.ToObjectID(req.OrgId)
	if err != nil {
		return "", &customErrors.ServerCustomError{Code: httpStatus.BAD_REQUEST, Errors: []string{"invalid orgId"}}
	}

	filters := model.Map{"type": req.Type, "value": req.Slug, "orgId": orgID}
	if req.ParentId != nil && *req.ParentId != "" {
		if pid, perr := model.ToObjectID(*req.ParentId); perr == nil {
			filters["parentId"] = pid
		}
	}

	result, err := s.deps.Repo.FindWithFilters(c, filters, &model.PaginationOpts{Page: 1, PerPage: 1}, nil)
	if err != nil {
		return "", err
	}
	if len(result.Items) == 0 {
		return "", nil
	}
	return result.Items[0].ID.Hex(), nil
}

// createResolvedList inserts the org-scoped list for the request — always
// IsSystem=false, so an install never adds a public/global classification. A
// duplicate-key race (a concurrent resolve created the same list first) is
// treated as success by re-resolving.
func (s *ListService) createResolvedList(c ctx.Context, req dtos.ListResolveRequest) (string, error) {
	orgID, err := model.ToObjectID(req.OrgId)
	if err != nil {
		return "", &customErrors.ServerCustomError{Code: httpStatus.BAD_REQUEST, Errors: []string{"invalid orgId"}}
	}

	dto := &dtos.ListCreateDTO{
		Type:       req.Type,
		Name:       req.Name,
		Value:      req.Slug,
		Enabled:    true,
		ParentId:   req.ParentId,
		IsSystem:   false,
		IsTemplate: req.ShareWithChildren,
		OrgID:      &orgID,
		PathKey:    req.PathKey,
	}

	entity, _ := mapper.DtoToEntityWithOptions[dtos.ListCreateDTO, entities.List](dto, mapper.MapperOptions{StringToObjectId: true})
	if entity.Scope == "" {
		entity.Scope = "local"
	}

	created, err := s.deps.Repo.Create(c, entity)
	if err != nil {
		if isDuplicateKeyError(err) {
			return s.findResolvedListId(c, req)
		}
		return "", err
	}
	return created.ID.Hex(), nil
}

// isDuplicateKeyError reports whether err is a MongoDB unique-index violation,
// detected by the driver's E11000 marker so it survives the model wrapper.
func isDuplicateKeyError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "E11000")
}

package services

import (
	"context"
	"time"

	"workflow/src/modules/definitions/application/di"
	"workflow/src/modules/definitions/application/dtos"
	"workflow/src/modules/definitions/application/ports"
	domainConstants "workflow/src/modules/definitions/domain/constants"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// Compile-time check to ensure DefinitionService implements DefinitionServicePort
var _ ports.DefinitionServicePort = (*DefinitionService)(nil)

// New creates and returns a new DefinitionService.
func New(deps di.DefinitionServiceDependenciesInjection) ports.DefinitionServicePort {
	return &DefinitionService{deps: deps}
}

// CreateDefinition validates the workflow graph, stamps tenant + plugin
// status, persists the definition, and uploads code-node scripts to MinIO
// (L2 cache). All steps are tracked via the create metric pair.
func (s *DefinitionService) CreateDefinition(ctx context.Context, requestContext *reqCtx.RequestContext, dto *dtos.DefinitionCreateDTO) (*dtos.DefinitionResponse, error) {
	start := time.Now()
	defEntity, err := s.prepareDefinitionForCreate(ctx, requestContext, dto)
	if err != nil {
		return nil, err
	}
	created, err := s.persistDefinitionForCreate(ctx, defEntity, start)
	if err != nil {
		return nil, err
	}
	return s.toDefinitionResponseTracked(created, "create", start)
}

// UpdateDefinitionById applies a partial update with graph re-validation when
// nodes/edges change, recomputes plugin status, bumps definition version on
// structural change, persists, re-uploads scripts to MinIO, and publishes a
// FANOUT cache invalidation for js-workflow-executor.
func (s *DefinitionService) UpdateDefinitionById(ctx context.Context, definitionId *string, dto *dtos.DefinitionUpdateDTO) (*dtos.DefinitionResponse, error) {
	start := time.Now()
	fieldsToUpdate, err := s.prepareDefinitionUpdate(ctx, definitionId, dto, start)
	if err != nil {
		return nil, err
	}
	updated, err := s.persistDefinitionUpdate(ctx, definitionId, fieldsToUpdate, start)
	if err != nil {
		return nil, err
	}
	if err := s.uploadAndInvalidateUpdatedScripts(ctx, updated, start); err != nil {
		return nil, err
	}
	return s.toDefinitionResponseTracked(updated, "update", start)
}

// DeleteDefinitionById removes the definition, then performs best-effort
// cleanup of MinIO node-data and a FANOUT invalidation. The pre-delete fetch
// is best-effort: a fetch failure logs a warning but does not block the
// delete itself, since the row removal is the operation users care about.
func (s *DefinitionService) DeleteDefinitionById(ctx context.Context, definitionId *string) (map[string]bool, error) {
	start := time.Now()
	def := s.fetchDefinitionForDelete(ctx, definitionId)

	if err := s.deps.DefinitionRepo.DeleteById(ctx, definitionId); err != nil {
		s.trackDefinitionMetrics("delete", "error", start)
		return nil, err
	}
	s.cleanupDefinitionAfterDelete(ctx, def)
	s.trackDefinitionMetrics("delete", "success", start)
	return map[string]bool{"success": true}, nil
}

// GetDefinitionById retrieves a definition by id. A missing or zero-id row
// surfaces as a 404 ServerCustomError so HTTP callers map it cleanly. All
// outcomes are observed by the read metric pair.
func (s *DefinitionService) GetDefinitionById(ctx context.Context, definitionId *string) (*dtos.DefinitionResponse, error) {
	start := time.Now()
	defEntity, err := s.deps.DefinitionRepo.FindById(ctx, definitionId)
	if err != nil || defEntity == nil || defEntity.ID.IsZero() {
		s.trackDefinitionMetrics("read", "error", start)
		return nil, &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"Workflow definition not found"}}
	}
	return s.toDefinitionResponseTracked(defEntity, "read", start)
}

// GetDefinitions returns a paginated list scoped by org context, with
// module-specific predicates (name, enabled, status, isTemplate, version)
// merged into the filter. Pagination + projection are derived from the DTO.
func (s *DefinitionService) GetDefinitions(ctx context.Context, requestContext *reqCtx.RequestContext, query *dtos.DefinitionQueryDTO) (*model.PaginatedResult[dtos.DefinitionResponse], error) {
	start := time.Now()
	filters, err := s.buildDefinitionListFilters(requestContext, query, start)
	if err != nil {
		return nil, err
	}
	pagination := s.buildDefinitionPagination(query)
	projection := orgfilter.BuildProjection(query.Projection)

	result, err := s.deps.DefinitionRepo.FindWithFilters(ctx, filters, pagination, projection)
	if err != nil {
		s.trackDefinitionMetrics("list", "error", start)
		return nil, err
	}
	return s.mapDefinitionPaginatedResultTracked(result, start)
}

// CountDefinitions returns the count of definitions visible to the request
// org scope. Errors are logged at the service layer (not just bubbled) so
// the dashboard endpoint can degrade gracefully with a 0 fallback.
func (s *DefinitionService) CountDefinitions(ctx context.Context, requestContext *reqCtx.RequestContext) (int64, error) {
	filters, err := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{
		ReqContext: requestContext,
	})
	if err != nil {
		logger.Error(err, "[SERVICE:Definition] Failed to build org filter for counter")
		return 0, err
	}
	count, err := s.deps.DefinitionRepo.CountDocuments(ctx, filters)
	if err != nil {
		logger.Error(err, "[SERVICE:Definition] Failed to count definitions")
		return 0, err
	}
	return count, nil
}

// GetNodeScript serves a code-node script as the HTTP fallback when the L2
// (MinIO) cache misses. After serving the script, an async write
// repopulates L2 in the background so subsequent requests skip this path.
// A non-code node, missing definition, or empty script all surface as 404.
func (s *DefinitionService) GetNodeScript(ctx context.Context, definitionId string, nodeId string) (string, error) {
	def, err := s.deps.DefinitionRepo.FindById(ctx, &definitionId)
	if err != nil || def == nil || def.ID.IsZero() {
		return "", &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"Workflow definition not found"}}
	}
	node, ok := findNode(def.Nodes, nodeId)
	if !ok || node.Type != domainConstants.NodeTypeCode {
		return "", &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"Code node not found"}}
	}
	script := getNodeScript(node)
	if script == "" {
		return "", &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"Script source is empty"}}
	}
	s.repopulateNodeScriptL2Async(def, definitionId, nodeId, script)
	return script, nil
}


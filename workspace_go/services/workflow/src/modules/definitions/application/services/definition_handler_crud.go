package services

import (
	"context"
	"fmt"
	"time"

	"workflow/src/modules/definitions/application/dtos"
	"workflow/src/modules/definitions/domain/entities"
	domainServices "workflow/src/modules/definitions/domain/services"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// prepareDefinitionForCreate orchestrates the pre-persistence pipeline:
// build the entity from the DTO + tenant context, run the domain validators,
// stamp the plugin-status snapshot, and apply the create-only defaults.
// Returns the entity ready for persistence or the first validation/build
// error encountered.
func (s *DefinitionService) prepareDefinitionForCreate(ctx context.Context, requestContext *reqCtx.RequestContext, dto *dtos.DefinitionCreateDTO) (*entities.WorkflowDefinition, error) {
	defEntity, err := s.buildDefinitionEntityForCreate(requestContext, dto)
	if err != nil {
		return nil, err
	}
	if err := s.validateDefinitionForCreate(defEntity); err != nil {
		return nil, err
	}
	s.applyDefinitionPluginStatus(ctx, defEntity, defEntity.Nodes)
	s.applyCreateDefaults(defEntity)
	return defEntity, nil
}

// prepareDefinitionUpdate orchestrates the pre-persistence pipeline for
// updates: fetch the pre-update state, run validations against the
// incoming DTO, and build the partial-update payload (with plugin-status
// recompute when nodes change). Returns the merge map ready for
// persistDefinitionUpdate, or the first validation/build error.
func (s *DefinitionService) prepareDefinitionUpdate(ctx context.Context, definitionId *string, dto *dtos.DefinitionUpdateDTO, start time.Time) (map[string]any, error) {
	defBefore, err := s.fetchDefinitionForUpdate(ctx, definitionId, start)
	if err != nil {
		return nil, err
	}
	if err := s.validateDefinitionUpdate(dto, defBefore); err != nil {
		return nil, err
	}
	return s.buildDefinitionUpdatePayload(ctx, dto, defBefore)
}

// buildDefinitionEntityForCreate maps the create DTO to a fresh entity and
// stamps it with the org context so the persisted document carries the same
// multi-tenant labels as future queries.
func (s *DefinitionService) buildDefinitionEntityForCreate(requestContext *reqCtx.RequestContext, dto *dtos.DefinitionCreateDTO) (*entities.WorkflowDefinition, error) {
	defEntity, err := mapper.DtoToEntity[dtos.DefinitionCreateDTO, entities.WorkflowDefinition](dto)
	if err != nil {
		return nil, fmt.Errorf("failed to map create DTO to entity: %w", err)
	}
	if requestContext.OrgContext != nil && *requestContext.OrgContext != "" {
		if orgObjectId, err := model.ToObjectID(*requestContext.OrgContext); err == nil {
			defEntity.OrgID = &orgObjectId
		}
	}
	if requestContext.OrgContextData != nil && requestContext.OrgContextData.PathKey != "" {
		defEntity.PathKey = requestContext.OrgContextData.PathKey
	}
	return defEntity, nil
}

// validateDefinitionForCreate runs the domain validators that must hold for
// every brand-new definition: no tight cycles without async pause points,
// and every node has a valid config. Errors surface as 400 ValidationError.
func (s *DefinitionService) validateDefinitionForCreate(defEntity *entities.WorkflowDefinition) error {
	if tightCycleNodeIDs := domainServices.DetectTightCycles(defEntity.Nodes, defEntity.Edges); len(tightCycleNodeIDs) > 0 {
		return &customErrors.ValidationError{Errors: []string{
			fmt.Sprintf("workflow contains tight cycles without async pause points (nodes: %v)", tightCycleNodeIDs),
		}}
	}
	if nodeErrors := domainServices.ValidateNodes(defEntity.Nodes); len(nodeErrors) > 0 {
		return &customErrors.ValidationError{Errors: collectNodeValidationErrors(nodeErrors)}
	}
	return nil
}

// applyDefinitionPluginStatus computes installedPlugins, status, and
// missingPlugins from the supplied node set. When no plugins are required
// the definition is implicitly StatusValid — there is nothing to be missing.
func (s *DefinitionService) applyDefinitionPluginStatus(ctx context.Context, defEntity *entities.WorkflowDefinition, nodes []entities.WorkflowNode) {
	requiredPlugins := domainServices.ExtractRequiredPlugins(nodes)
	defEntity.InstalledPlugins = requiredPlugins
	if len(requiredPlugins) == 0 {
		defEntity.Status = string(entities.StatusValid)
		return
	}
	pluginStatus, missingPlugins := s.computePluginStatusForRequired(ctx, requiredPlugins)
	defEntity.Status = string(pluginStatus)
	defEntity.MissingPlugins = missingPlugins
}

// computePluginStatusForRequired loads the enabled-plugin list and asks the
// domain to classify the definition (valid / plugin_missing / invalid)
// against the required set. A loader failure logs a warning and falls
// through with the existing manifest list (best-effort).
func (s *DefinitionService) computePluginStatusForRequired(ctx context.Context, requiredPlugins []string) (entities.DefinitionStatus, []string) {
	enabledManifests, err := s.deps.PluginLoader.GetAllEnabled(ctx)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Definition] Failed to load enabled plugins: %s", err))
	}
	enabledIDs := make([]string, len(enabledManifests))
	for i, m := range enabledManifests {
		enabledIDs[i] = m.PluginID
	}
	return domainServices.ComputeDefinitionStatus(requiredPlugins, enabledIDs)
}

// applyCreateDefaults stamps the immutable create-time fields onto the new
// definition: version 1 and a paired created/updated timestamp.
func (s *DefinitionService) applyCreateDefaults(defEntity *entities.WorkflowDefinition) {
	defEntity.DefinitionVersion = 1
	now := time.Now()
	defEntity.Created = now
	defEntity.Updated = now
}

// persistDefinitionForCreate runs the two-phase create — insert the row,
// then upload code-node scripts to MinIO. A persist failure rolls up the
// raw repo error; an upload failure surfaces a wrapped error so the client
// knows the row exists but the L2 cache is missing scripts.
func (s *DefinitionService) persistDefinitionForCreate(ctx context.Context, defEntity *entities.WorkflowDefinition, start time.Time) (*entities.WorkflowDefinition, error) {
	created, err := s.deps.DefinitionRepo.Create(ctx, defEntity)
	if err != nil {
		s.trackDefinitionMetrics("create", "error", start)
		return nil, err
	}
	if err := s.writeNodeScripts(ctx, created); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Definition] Script upload failed for definition %s", created.ID.Hex()))
		s.trackDefinitionMetrics("create", "error", start)
		return nil, scriptUploadFailedError("create", err)
	}
	return created, nil
}

// fetchDefinitionForUpdate loads the pre-update snapshot. A missing row
// surfaces as 404 immediately so the caller never sees a partial update.
func (s *DefinitionService) fetchDefinitionForUpdate(ctx context.Context, definitionId *string, start time.Time) (*entities.WorkflowDefinition, error) {
	defBefore, err := s.deps.DefinitionRepo.FindById(ctx, definitionId)
	if err != nil || defBefore == nil || defBefore.ID.IsZero() {
		s.trackDefinitionMetrics("update", "error", start)
		return nil, notFoundDefinition()
	}
	return defBefore, nil
}

// validateDefinitionUpdate applies graph-structure checks only when the DTO
// touches nodes or edges. Tight-cycle detection always runs (using the
// merged before/after node + edge sets); per-node config validation only
// runs when nodes are explicitly being updated.
func (s *DefinitionService) validateDefinitionUpdate(dto *dtos.DefinitionUpdateDTO, defBefore *entities.WorkflowDefinition) error {
	if dto.Nodes == nil && dto.Edges == nil {
		return nil
	}
	nodesToValidate, edgesToValidate := s.resolveValidationGraph(dto, defBefore)
	if tightCycleNodeIDs := domainServices.DetectTightCycles(nodesToValidate, edgesToValidate); len(tightCycleNodeIDs) > 0 {
		return &customErrors.ValidationError{Errors: []string{
			fmt.Sprintf("workflow contains tight cycles without async pause points (nodes: %v)", tightCycleNodeIDs),
		}}
	}
	if dto.Nodes != nil {
		if nodeErrors := domainServices.ValidateNodes(nodesToValidate); len(nodeErrors) > 0 {
			return &customErrors.ValidationError{Errors: collectNodeValidationErrors(nodeErrors)}
		}
	}
	return nil
}

// resolveValidationGraph picks the effective node/edge slices for graph
// validation: caller-supplied if present, otherwise the existing persisted
// values. Lets a partial update still validate the merged graph shape.
func (s *DefinitionService) resolveValidationGraph(dto *dtos.DefinitionUpdateDTO, defBefore *entities.WorkflowDefinition) ([]entities.WorkflowNode, []entities.WorkflowEdge) {
	nodes := defBefore.Nodes
	edges := defBefore.Edges
	if dto.Nodes != nil {
		nodes = contractNodesToEntity(*dto.Nodes)
	}
	if dto.Edges != nil {
		edges = contractEdgesToEntity(*dto.Edges)
	}
	return nodes, edges
}

// buildDefinitionUpdatePayload converts the partial DTO to a Mongo $set map,
// recomputes plugin status against the effective node set, bumps the
// definition version on structural changes, and stamps the updated time.
func (s *DefinitionService) buildDefinitionUpdatePayload(ctx context.Context, dto *dtos.DefinitionUpdateDTO, defBefore *entities.WorkflowDefinition) (map[string]any, error) {
	fieldsToUpdate, err := mapper.DtoToMap(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to map update DTO to map: %w", err)
	}
	effectiveNodes := defBefore.Nodes
	if dto.Nodes != nil {
		effectiveNodes = contractNodesToEntity(*dto.Nodes)
	}
	s.applyPluginStatusToUpdatePayload(ctx, fieldsToUpdate, effectiveNodes)
	if dto.Nodes != nil || dto.Edges != nil {
		fieldsToUpdate["definitionVersion"] = defBefore.DefinitionVersion + 1
	}
	fieldsToUpdate["updated"] = time.Now()
	return fieldsToUpdate, nil
}

// applyPluginStatusToUpdatePayload writes installedPlugins, status, and
// missingPlugins into the update map. When no plugins are required, the
// missing list is forced to empty so a previous plugin_missing flag clears.
func (s *DefinitionService) applyPluginStatusToUpdatePayload(ctx context.Context, fieldsToUpdate map[string]any, effectiveNodes []entities.WorkflowNode) {
	requiredPlugins := domainServices.ExtractRequiredPlugins(effectiveNodes)
	fieldsToUpdate["installedPlugins"] = requiredPlugins
	if len(requiredPlugins) == 0 {
		fieldsToUpdate["status"] = string(entities.StatusValid)
		fieldsToUpdate["missingPlugins"] = []string{}
		return
	}
	pluginStatus, missingPlugins := s.computePluginStatusForRequired(ctx, requiredPlugins)
	fieldsToUpdate["status"] = string(pluginStatus)
	fieldsToUpdate["missingPlugins"] = missingPlugins
}

// persistDefinitionUpdate writes the update map and validates the updated
// row exists. A nil/zero-id result means the row vanished mid-flight (race
// with delete) and surfaces as 404 — same envelope the pre-fetch uses.
func (s *DefinitionService) persistDefinitionUpdate(ctx context.Context, definitionId *string, fieldsToUpdate map[string]any, start time.Time) (*entities.WorkflowDefinition, error) {
	updated, _ := s.deps.DefinitionRepo.FindByIdAndUpdate(ctx, definitionId, fieldsToUpdate)
	if updated == nil || updated.ID.IsZero() {
		s.trackDefinitionMetrics("update", "error", start)
		return nil, notFoundDefinition()
	}
	return updated, nil
}

// uploadAndInvalidateUpdatedScripts re-uploads every code-node script to
// MinIO (covering script body changes the DTO might carry) and publishes a
// FANOUT message so js-workflow-executor invalidates its in-process L0/L1
// caches for the affected nodes.
func (s *DefinitionService) uploadAndInvalidateUpdatedScripts(ctx context.Context, updated *entities.WorkflowDefinition, start time.Time) error {
	if err := s.writeNodeScripts(ctx, updated); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Definition] Script upload failed for definition %s", updated.ID.Hex()))
		s.trackDefinitionMetrics("update", "error", start)
		return scriptUploadFailedError("update", err)
	}
	codeNodeIds := getCodeNodeIds(updated)
	if len(codeNodeIds) > 0 {
		s.publishDefinitionInvalidate(ctx, updated, codeNodeIds)
	}
	return nil
}

// fetchDefinitionForDelete is a best-effort pre-delete read: a fetch failure
// logs a warning but does not block the delete itself, since the row removal
// is the primary intent of the call.
func (s *DefinitionService) fetchDefinitionForDelete(ctx context.Context, definitionId *string) *entities.WorkflowDefinition {
	def, fetchErr := s.deps.DefinitionRepo.FindById(ctx, definitionId)
	if fetchErr != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Definition] Failed to fetch definition %s before delete (MinIO cleanup skipped): %v", *definitionId, fetchErr))
	}
	return def
}

// cleanupDefinitionAfterDelete removes the MinIO node-data and publishes a
// FANOUT cache invalidation. Nil-safe so a failed pre-delete fetch simply
// skips cleanup rather than crashing the request.
func (s *DefinitionService) cleanupDefinitionAfterDelete(ctx context.Context, def *entities.WorkflowDefinition) {
	if def == nil {
		return
	}
	codeNodeIds := getCodeNodeIds(def)
	s.deleteAllNodeData(ctx, def, codeNodeIds)
	if len(codeNodeIds) > 0 {
		s.publishDefinitionInvalidate(ctx, def, codeNodeIds)
	}
}

// collectNodeValidationErrors flattens the per-node errors into the flat
// string list accepted by ValidationError, prefixing each with the offending
// node id and type for easy client diagnosis.
func collectNodeValidationErrors(nodeErrors []domainServices.NodeValidationError) []string {
	errs := make([]string, 0, len(nodeErrors))
	for _, ne := range nodeErrors {
		for _, e := range ne.Errors {
			errs = append(errs, fmt.Sprintf("node %s (%s): %s", ne.NodeID, ne.NodeType, e))
		}
	}
	return errs
}

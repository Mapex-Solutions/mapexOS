package minio

import (
	"context"
	"encoding/json"
	"fmt"

	"assets/src/modules/assettemplates/application/ports"
	"assets/src/modules/assettemplates/domain/entities"

	minioModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/minio"
)

// NewTemplateStorageAdapter creates a new TemplateStorageAdapter.
//
// Parameters:
//   - client: Configured MinIO client for object storage operations
//
// Returns:
//   - ports.TemplateStoragePort: The port interface implementation
func NewTemplateStorageAdapter(client *minioModel.MinIOClient) ports.TemplateStoragePort {
	return &TemplateStorageAdapter{
		client: client,
	}
}

// Compile-time check to ensure TemplateStorageAdapter implements TemplateStoragePort interface.
var _ ports.TemplateStoragePort = (*TemplateStorageAdapter)(nil)

// WriteScripts writes template scripts to MinIO (L2 cache).
//
// Key format: {orgId}/{templateId}.json
// If IsSystem=true, uses "mapexos_public" as orgId (accessible to all orgs)
// Example: mapexos_public/507f1f77bcf86cd799439011.json (system template)
// Example: 507f1f77bcf86cd799439011/607f1f77bcf86cd799439022.json (org template)
func (a *TemplateStorageAdapter) WriteScripts(ctx context.Context, template *entities.Assettemplate) error {
	if template == nil || template.ID.IsZero() {
		return nil
	}

	// Build full template payload (scripts + dynamicFields)
	payload := a.buildPayload(template)

	// Serialize to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to serialize template scripts: %w", err)
	}

	// Determine orgId based on IsSystem flag
	orgId := a.getOrgId(template)

	// Write to MinIO with orgId prefix for tenant isolation
	// Key format: {orgId}/{templateId}.json
	key := orgId + "/" + template.ID.Hex() + ".json"
	if err := a.client.PutJSON(ctx, key, data); err != nil {
		return fmt.Errorf("failed to write scripts to MinIO: %w", err)
	}

	return nil
}

// DeleteScripts removes template scripts from MinIO (L2 cache).
//
// Called when a template is deleted to ensure consuming services get cache miss.
//
// Key format: {orgId}/{templateId}.json
func (a *TemplateStorageAdapter) DeleteScripts(ctx context.Context, orgId string, templateId string) error {
	if orgId == "" || templateId == "" {
		return nil
	}

	// Key format: {orgId}/{templateId}.json
	key := orgId + "/" + templateId + ".json"
	if err := a.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete scripts from MinIO: %w", err)
	}

	return nil
}

// buildPayload converts template entity to TemplatePayload for MinIO storage.
func (a *TemplateStorageAdapter) buildPayload(template *entities.Assettemplate) TemplatePayload {
	payload := TemplatePayload{
		ID:               template.ID.Hex(),
		Name:             template.Name,
		ScriptValidator:  template.ScriptValidator,
		ScriptConversion: template.ScriptConversion,
		NextFieldId:      template.NextFieldId,
	}

	if template.Description != nil {
		payload.Description = *template.Description
	}

	// Handle optional scripts
	if template.ScriptTest != nil {
		payload.ScriptTest = *template.ScriptTest
	}
	if template.ScriptProcessor != nil {
		payload.ScriptProcessor = *template.ScriptProcessor
	}

	// Map DynamicFields from entity to payload
	payload.DynamicFields = make([]DynamicFieldPayload, len(template.DynamicFields))
	for i, f := range template.DynamicFields {
		payload.DynamicFields[i] = DynamicFieldPayload{
			FieldId:       f.FieldId,
			Field:         f.Field,
			Value:         f.Value,
			Type:          f.Type,
			Status:        f.Status,
			LatitudePath:  f.LatitudePath,
			LongitudePath: f.LongitudePath,
		}
	}

	return payload
}

// getOrgId determines the organization ID for storage key based on template visibility.
//
// Logic:
//   - IsSystem=true → returns "mapexos_public" (accessible to all orgs)
//   - IsSystem=false → returns the template's OrgID
//
// This enables proper tenant isolation while supporting global system templates.
func (a *TemplateStorageAdapter) getOrgId(template *entities.Assettemplate) string {
	if template.IsSystem {
		return PublicOrgID
	}

	if template.OrgID != nil {
		return template.OrgID.Hex()
	}

	// Fallback to public if no OrgID (shouldn't happen for non-system templates)
	return PublicOrgID
}

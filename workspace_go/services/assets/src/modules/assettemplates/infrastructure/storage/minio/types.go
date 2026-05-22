package minio

import (
	minioModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/minio"
	"go.uber.org/dig"
)

/* TemplateStorageProviderParams */

// TemplateStorageProviderParams defines the dependencies for TemplateStoragePort provider.
// Uses dig.In to enable automatic dependency injection with named parameters.
type TemplateStorageProviderParams struct {
	dig.In

	// MinIOClient is injected with the "templates" name tag
	// Configured with bucket and key prefix for template scripts
	MinIOClient *minioModel.MinIOClient `name:"templates"`
}

// TemplateStorageAdapter implements TemplateStoragePort using MinIO for object storage.
//
// This adapter follows Hexagonal Architecture by:
//   - Implementing the application port interface
//   - Encapsulating all MinIO/infrastructure details
//   - Keeping the application layer clean from storage concerns
//
// The adapter handles:
//   - Entity to scripts payload conversion
//   - JSON serialization
//   - MinIO operations (put/delete)
type TemplateStorageAdapter struct {
	client *minioModel.MinIOClient
}

/* DynamicFieldPayload */

// DynamicFieldPayload represents a DynamicField in the MinIO payload.
type DynamicFieldPayload struct {
	FieldId       uint16 `json:"fieldId"`
	Field         string `json:"field"`
	Value         string `json:"value,omitempty"`
	Type          string `json:"type"`
	Status        uint8  `json:"status"`
	LatitudePath  string `json:"latitudePath,omitempty"`
	LongitudePath string `json:"longitudePath,omitempty"`
}

// TemplatePayload represents the template data stored in MinIO (L2 cache).
// This payload serves multiple consuming services:
//   - JS-Executor: uses scripts (scriptValidator, scriptConversion, etc.)
//   - Events: uses dynamicFields for EVA field mapping
type TemplatePayload struct {
	ID               string                `json:"_id"`
	Name             string                `json:"name"`
	Description      string                `json:"description,omitempty"`
	ScriptTest       string                `json:"scriptTest,omitempty"`
	ScriptProcessor  string                `json:"scriptProcessor,omitempty"`
	ScriptValidator  string                `json:"scriptValidator"`
	ScriptConversion string                `json:"scriptConversion"`
	DynamicFields    []DynamicFieldPayload `json:"dynamicFields"`
	NextFieldId      uint16                `json:"nextFieldId"`
}

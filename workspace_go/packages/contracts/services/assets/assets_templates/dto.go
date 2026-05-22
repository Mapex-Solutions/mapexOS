package assetstemplate

import (
	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type AssetTemplateId struct {
	AssetTemplateId string `params:"assetTemplateId" validate:"required,mongoid"`
}

/**
 * SHARED READ MODEL (CQRS Pattern)
 *
 * TemplateScriptsReadModel represents the shared read model for template scripts.
 * This is a CQRS Read Model - a denormalized projection optimized for reads.
 *
 * OWNERSHIP: Assets Service (write only)
 * CONSUMERS: JS-Executor, Events (read only)
 *
 * Key Format: {orgId}/{templateId}.json (or mapexos_public/{templateId}.json for system templates)
 * Storage: MinIO (L2 cache)
 */
type TemplateScriptsReadModel struct {
	ScriptTest       string `json:"scriptTest,omitempty"`
	ScriptProcessor  string `json:"scriptProcessor,omitempty"`
	ScriptValidator  string `json:"scriptValidator"`
	ScriptConversion string `json:"scriptConversion"`
}

// DynamicField represents a field mapping for dynamic/typed storage
// Used for efficient event querying and filtering
type DynamicField struct {
	FieldId       uint16 `json:"fieldId"`
	Field         string `json:"field" validate:"required"`
	Value         string `json:"value,omitempty" validate:"omitempty"`
	Type          string `json:"type" validate:"required,oneof=string number bool date geo"`
	Status        uint8  `json:"status"`
	LatitudePath  string `json:"latitudePath,omitempty" validate:"omitempty"`
	LongitudePath string `json:"longitudePath,omitempty" validate:"omitempty"`
}

type AssetTemplateCreate struct {
	Name        string  `json:"name" validate:"required,min=1"`
	Enabled     bool    `json:"enabled" validate:"required"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`

	// Asset Classification - IDs (source of truth) + Names (denormalized for performance)
	// IDs come as strings in JSON and are converted to ObjectId in service layer
	CategoryId       *string `json:"categoryId,omitempty" validate:"omitempty,mongoid"`
	CategoryName     *string `json:"categoryName,omitempty" validate:"omitempty,max=254"`
	ManufacturerId   *string `json:"manufacturerId,omitempty" validate:"omitempty,mongoid"`
	ManufacturerName *string `json:"manufacturerName,omitempty" validate:"omitempty,max=254"`
	ModelId          *string `json:"modelId,omitempty" validate:"omitempty,mongoid"`
	ModelName        *string `json:"modelName,omitempty" validate:"omitempty,max=254"`
	Version          *string `json:"version,omitempty" validate:"omitempty,max=100"`

	// Template visibility flags
	IsSystem   bool `json:"isSystem,omitempty" validate:"omitempty"`
	IsTemplate bool `json:"isTemplate,omitempty" validate:"omitempty"`

	AssetIDPath string `json:"assetIdPath" validate:"required,min=1"`

	ScriptProcessor  *string `json:"scriptProcessor" validate:"omitempty"`
	ScriptValidator  *string `json:"scriptValidator" validate:"omitempty"`
	ScriptConversion string  `json:"scriptConversion" validate:"required"`
	ScriptTest       string  `json:"scriptTest,omitempty" validate:"omitempty"`

	// Available Fields - for Rule autocomplete support
	AvailableFields []string `json:"availableFields,omitempty" validate:"omitempty"`

	// Dynamic Fields - for typed event storage and querying
	DynamicFields []DynamicField `json:"dynamicFields,omitempty" validate:"omitempty,dive"`

	// Multi-tenant fields (populated automatically by coverage middleware)
	OrgID   *model.ObjectId `json:"orgId,omitempty" validate:"omitempty"`
	PathKey *string         `json:"pathKey,omitempty" validate:"omitempty"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

type AssetTemplateUpdate struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1"`
	Enabled     *bool   `json:"enabled,omitempty" validate:"omitempty"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`

	// Asset Classification - IDs (source of truth) + Names (denormalized for performance)
	// IDs come as strings in JSON and are converted to ObjectId in service layer
	CategoryId       *string `json:"categoryId,omitempty" validate:"omitempty,mongoid"`
	CategoryName     *string `json:"categoryName,omitempty" validate:"omitempty,max=254"`
	ManufacturerId   *string `json:"manufacturerId,omitempty" validate:"omitempty,mongoid"`
	ManufacturerName *string `json:"manufacturerName,omitempty" validate:"omitempty,max=254"`
	ModelId          *string `json:"modelId,omitempty" validate:"omitempty,mongoid"`
	ModelName        *string `json:"modelName,omitempty" validate:"omitempty,max=254"`
	Version          *string `json:"version,omitempty" validate:"omitempty,max=100"`

	// Template visibility flags
	IsSystem   *bool `json:"isSystem,omitempty" validate:"omitempty"`
	IsTemplate *bool `json:"isTemplate,omitempty" validate:"omitempty"`

	AssetIDPath *string `json:"assetIdPath" validate:"required,min=1"`

	ScriptTest       *string `json:"scriptTest,omitempty" validate:"omitempty"`
	ScriptProcessor  *string `json:"scriptProcessor" validate:"omitempty"`
	ScriptValidator  *string `json:"scriptValidator" validate:"omitempty"`
	ScriptConversion *string `json:"scriptConversion" validate:"omitempty"`

	// Available Fields - for Rule autocomplete support
	AvailableFields []string `json:"availableFields,omitempty" validate:"omitempty"`

	// Dynamic Fields - for typed event storage and querying
	DynamicFields []DynamicField `json:"dynamicFields,omitempty" validate:"omitempty,dive"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

type AssetTemplateResponse struct {
	ID          *common.ObjectID `json:"id,omitempty"`
	Name        *string          `json:"name,omitempty"`
	Enabled     *bool            `json:"enabled,omitempty"`
	Description *string          `json:"description,omitempty"`

	// Asset Classification - IDs (source of truth) + Names (denormalized for performance)
	CategoryId       *model.ObjectId `json:"categoryId,omitempty"`
	CategoryName     *string         `json:"categoryName,omitempty"`
	ManufacturerId   *model.ObjectId `json:"manufacturerId,omitempty"`
	ManufacturerName *string         `json:"manufacturerName,omitempty"`
	ModelId          *model.ObjectId `json:"modelId,omitempty"`
	ModelName        *string         `json:"modelName,omitempty"`
	Version          *string         `json:"version,omitempty"`

	AssetIDPath *string         `json:"assetIdPath,omitempty"`
	OrgId       *model.ObjectId `json:"orgId,omitempty"`

	// Template visibility flags
	IsSystem   *bool `json:"isSystem,omitempty"`
	IsTemplate *bool `json:"isTemplate,omitempty"`

	ScriptTest       *string `json:"scriptTest,omitempty"`
	ScriptProcessor  *string `json:"scriptProcessor,omitempty"`
	ScriptValidator  *string `json:"scriptValidator,omitempty"`
	ScriptConversion *string `json:"scriptConversion,omitempty"`

	// Available Fields - for Rule autocomplete support
	AvailableFields []string `json:"availableFields,omitempty"`

	// Dynamic Fields - for typed event storage and querying
	DynamicFields []DynamicField `json:"dynamicFields,omitempty"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

func (d *AssetTemplateResponse) SetCreated(t *common.NullTime) { d.Created = t }
func (d *AssetTemplateResponse) SetUpdated(t *common.NullTime) { d.Updated = t }

// AssetTemplateQuery represents query parameters for listing asset templates.
// Embeds BaseQueryDTO for standard pagination, sorting, and hierarchy support.
//
// Standard fields (from BaseQueryDTO):
//   - Projection: comma-separated fields to return
//   - Page: page number (default: 1)
//   - PerPage: items per page (default: 20)
//   - Sort: sort order (default: "created:desc")
//   - IncludeChildren: include child orgs hierarchically (default: false)
//
// Module-specific filters:
//   - Name: filter by template name (partial match)
//   - Enabled: filter by enabled status (true/false)
//   - IsSystem: filter system templates
//   - Manufacture: filter by manufacturer
//   - Model: filter by model
//
// Organization filtering is handled automatically via RequestContext:
//   - No manual orgId/customerId needed
//   - Context-aware filtering via X-Org-Context header
//   - Hierarchical queries via includeChildren parameter
type AssetTemplateQuery struct {
	query.BaseQueryDTO

	// Module-specific filters
	Name           *string `query:"name" validate:"omitempty,max=150"`
	Enabled        *bool   `query:"enabled" validate:"omitempty"`
	IsSystem       *bool   `query:"isSystem" validate:"omitempty"`
	IsTemplate     *bool   `query:"isTemplate" validate:"omitempty"`
	ManufacturerId *string `query:"manufacturerId" validate:"omitempty,mongoid"`
	ModelId        *string `query:"modelId"`
}

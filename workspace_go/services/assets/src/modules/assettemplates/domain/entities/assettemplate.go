package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// DynamicField represents a field mapping for dynamic/typed storage in ClickHouse.
// Used for efficient event querying and filtering with MAP<INT, TYPE> storage.
//
// The FieldId is used for EVA (Entity-Value-Attribute) storage optimization:
//   - FieldId (uint16) instead of Field (string) for 3-4x faster queries
//   - Better compression (no repeated field names in storage)
//   - Query rewrite: "temperature > 25" → "eva_number[1] > 25"
//
// Status values:
//   - 1 = active (visible in UI, used for new events)
//   - 0 = deprecated (hidden in UI, preserved for historical queries)
//
// Rules:
//   - FieldId is IMMUTABLE once assigned (never reuse)
//   - Maximum 200 active fields per template
//   - Deleted fields get Status=0, never actually deleted
type DynamicField struct {
	FieldId       uint16 `bson:"fieldId"`                          // Unique numeric ID for MAP storage (1-65535, auto-increment)
	Field         string `bson:"field"`                              // Human-readable field name for display/debug
	Value         string `bson:"value,omitempty"`          // JSON path to extract value from event payload
	Type          string `bson:"type"`                                // Data type: "number", "string", "bool", "date", "geo"
	Status        uint8  `bson:"status"`                            // 1=active, 0=deprecated
	LatitudePath  string `bson:"latitudePath,omitempty"`  // For geo type: path to latitude
	LongitudePath string `bson:"longitudePath,omitempty"` // For geo type: path to longitude
}

type Assettemplate struct {
	ID          model.ObjectId `bson:"_id,omitempty"`
	Name        string         `bson:"name"`
	Enabled     bool           `bson:"enabled"`
	Description *string        `bson:"description,omitempty"`

	// Asset Classification - IDs (source of truth) + Names (denormalized for performance)
	CategoryId       *model.ObjectId `bson:"categoryId,omitempty"`
	CategoryName     *string         `bson:"categoryName,omitempty"`
	ManufacturerId   *model.ObjectId `bson:"manufacturerId,omitempty"`
	ManufacturerName *string         `bson:"manufacturerName,omitempty"`
	ModelId          *model.ObjectId `bson:"modelId,omitempty"`
	ModelName        *string         `bson:"modelName,omitempty"`
	Version          *string         `bson:"version,omitempty"`

	AssetIDPath string `bson:"assetIdPath"`

	ScriptTest       *string `bson:"scriptTest,omitempty"`
	ScriptProcessor  *string `bson:"scriptProcessor,omitempty"`
	ScriptValidator  string  `bson:"scriptValidator"`
	ScriptConversion string  `bson:"scriptConversion"`

	// Available Fields - for Rule autocomplete support
	AvailableFields []string `bson:"availableFields,omitempty"`

	// Dynamic Fields - for typed event storage and querying (EVA pattern)
	DynamicFields []DynamicField `bson:"dynamicFields,omitempty"`

	// NextFieldId is the auto-increment counter for assigning FieldId to new DynamicFields.
	// Starts at 1, increments for each new field, NEVER decrements or reuses values.
	NextFieldId uint16 `bson:"nextFieldId,omitempty"`

	// Template visibility flags
	IsSystem   bool `bson:"isSystem"`   // true = visible to everyone (MAPEX global templates)
	IsTemplate bool `bson:"isTemplate"` // true = shared template (vendor/customer only)

	// Multi-tenant fields
	OrgID   *model.ObjectId `bson:"orgId,omitempty"`   // null for system, org for template/local
	PathKey *string         `bson:"pathKey,omitempty"` // null for system, pathKey for template/local

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

// Keep these as requested
func (u *Assettemplate) GetCreated() time.Time { return u.Created }
func (u *Assettemplate) GetUpdated() time.Time { return u.Updated }

// PATCH/UPDATE payload (every field optional)
type AssetTemplateUpdate struct {
	ID          model.ObjectId `bson:"_id,omitempty"`
	Name        string         `bson:"name,omitempty"`
	Enabled     bool           `bson:"enabled,omitempty"`
	Description *string        `bson:"description,omitempty"`

	// Asset Classification - IDs (source of truth) + Names (denormalized for performance)
	CategoryId       *model.ObjectId `bson:"categoryId,omitempty"`
	CategoryName     *string         `bson:"categoryName,omitempty"`
	ManufacturerId   *model.ObjectId `bson:"manufacturerId,omitempty"`
	ManufacturerName *string         `bson:"manufacturerName,omitempty"`
	ModelId          *model.ObjectId `bson:"modelId,omitempty"`
	ModelName        *string         `bson:"modelName,omitempty"`
	Version          *string         `bson:"version,omitempty"`

	AssetIDPath string `bson:"assetIdPath,omitempty"`

	ScriptTest       *string `bson:"scriptTest,omitempty"`
	ScriptProcessor  *string `bson:"scriptProcessor,omitempty"`
	ScriptValidator  string  `bson:"scriptValidator,omitempty"`
	ScriptConversion string  `bson:"scriptConversion,omitempty"`

	// Available Fields - for Rule autocomplete support
	AvailableFields []string `bson:"availableFields,omitempty"`

	// Dynamic Fields - for typed event storage and querying (EVA pattern)
	DynamicFields []DynamicField `bson:"dynamicFields,omitempty"`

	// NextFieldId counter for DynamicFields
	NextFieldId *uint16 `bson:"nextFieldId,omitempty"`

	// Auditing (usually only Updated is set by service/repo)
	Created *time.Time `bson:"created,omitempty"`
	Updated *time.Time `bson:"updated,omitempty"`
}

func (u *AssetTemplateUpdate) GetCreated() *time.Time { return u.Created }
func (u *AssetTemplateUpdate) GetUpdated() *time.Time { return u.Updated }

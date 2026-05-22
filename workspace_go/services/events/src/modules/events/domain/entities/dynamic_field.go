package entities

// DynamicField represents a field mapping from AssetTemplate.
// Used for EVA field resolution: field_name → fieldId → typed ClickHouse MAP column.
//
// Field types and their EVA column mapping:
//   - "number", "float", "int", "integer" → eva_number (MAP<UInt16, Float64>)
//   - "string", "text" → eva_string (MAP<UInt16, String>)
//   - "boolean", "bool" → eva_bool (MAP<UInt16, UInt8>)
//   - "date", "datetime", "timestamp" → eva_date (MAP<UInt16, DateTime>)
//   - "geo" → eva_number (lat) + eva_string (lat,lng) using LatitudePath/LongitudePath
type DynamicField struct {
	FieldId       uint16
	Field         string
	Value         string
	Type          string
	Status        uint8
	LatitudePath  string
	LongitudePath string
}

// CachedTemplate contains the template data needed for EVA field resolution.
// This is a minimal view of AssetTemplate, containing only what's needed for events.
type CachedTemplate struct {
	ID            string
	Name          string
	Description   string
	DynamicFields []DynamicField
	NextFieldId   uint16
}

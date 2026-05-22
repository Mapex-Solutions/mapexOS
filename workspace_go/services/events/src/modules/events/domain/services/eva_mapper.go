package domainservices

import (
	"fmt"
	"time"

	"events/src/modules/events/domain/entities"

	"github.com/Mapex-Solutions/mapexGoKit/utils/flatten"
	"github.com/Mapex-Solutions/mapexGoKit/utils/typeconv"
)

// MapEvaFields populates EVA typed maps on the event based on DynamicFields configuration.
// This is a Domain Service: pure business logic with no I/O, no infrastructure dependencies.
//
// For each active DynamicField:
//   - Extracts the value from eventData using the field's JSON path (Value)
//   - Converts to the appropriate type (number, string, boolean, date)
//   - Stores in the event's EVA MAP using the field's FieldId as key
//
// Parameters:
//   - event: Event entity to populate EVA fields
//   - fields: DynamicFields from AssetTemplate (already fetched by caller)
//   - eventData: Map containing the event payload
func MapEvaFields(event *entities.Event, fields []entities.DynamicField, eventData map[string]interface{}) {
	for _, field := range fields {
		// Skip deprecated fields
		if field.Status != 1 {
			continue
		}

		// Skip fields without FieldId (shouldn't happen, but safety check)
		if field.FieldId == 0 {
			continue
		}

		// Handle geo type specially (has latitude and longitude paths)
		if field.Type == "geo" {
			mapGeoField(event, field, eventData)
			continue
		}

		// Skip fields without a value path
		if field.Value == "" {
			continue
		}

		// Extract value using JSON path (supports nested paths like "data.value")
		value, multi, err := flatten.GetValueByPath(eventData, field.Value)
		if err != nil || multi {
			continue // Path not found or multi-value (array) - skip
		}

		// Map value to appropriate EVA type based on field type
		setEvaFieldByType(event, field.FieldId, field.Type, value)
	}
}

// mapGeoField handles geo type fields which have separate latitude and longitude paths.
// Stores latitude as eva_number and "lat,lng" as eva_string using the field's FieldId.
func mapGeoField(event *entities.Event, field entities.DynamicField, eventData map[string]interface{}) {
	// Extract latitude
	if field.LatitudePath != "" {
		latVal, multi, err := flatten.GetValueByPath(eventData, field.LatitudePath)
		if err == nil && !multi {
			if lat, ok := typeconv.TryFloat64(latVal); ok {
				event.SetEvaNumber(field.FieldId, lat)
			}
		}
	}

	// Extract longitude (use fieldId + 1 by convention for geo pairs)
	if field.LongitudePath != "" {
		lngVal, multi, err := flatten.GetValueByPath(eventData, field.LongitudePath)
		if err == nil && !multi {
			if lng, ok := typeconv.TryFloat64(lngVal); ok {
				latVal, _, _ := flatten.GetValueByPath(eventData, field.LatitudePath)
				event.SetEvaString(field.FieldId, fmt.Sprintf("%v,%v", latVal, lng))
			}
		}
	}
}

// setEvaFieldByType sets the EVA field based on the declared type.
func setEvaFieldByType(event *entities.Event, fieldId uint16, fieldType string, value interface{}) {
	switch fieldType {
	case "number", "float", "int", "integer":
		if numVal, ok := typeconv.TryFloat64(value); ok {
			event.SetEvaNumber(fieldId, numVal)
		}

	case "string", "text":
		if strVal, ok := value.(string); ok {
			event.SetEvaString(fieldId, strVal)
		} else {
			event.SetEvaString(fieldId, fmt.Sprintf("%v", value))
		}

	case "boolean", "bool":
		if boolVal, ok := typeconv.TryBool(value); ok {
			event.SetEvaBool(fieldId, boolVal)
		}

	case "date", "datetime", "timestamp":
		if timeVal, ok := typeconv.TryTime(value, ""); ok {
			event.SetEvaDate(fieldId, timeVal)
		}

	default:
		inferAndSetEvaField(event, fieldId, value)
	}
}

// inferAndSetEvaField attempts to infer the EVA type from the Go value and set it.
func inferAndSetEvaField(event *entities.Event, fieldId uint16, value interface{}) {
	switch v := value.(type) {
	case float64:
		event.SetEvaNumber(fieldId, v)
	case float32:
		event.SetEvaNumber(fieldId, float64(v))
	case int:
		event.SetEvaNumber(fieldId, float64(v))
	case int64:
		event.SetEvaNumber(fieldId, float64(v))
	case int32:
		event.SetEvaNumber(fieldId, float64(v))
	case bool:
		event.SetEvaBool(fieldId, v)
	case string:
		event.SetEvaString(fieldId, v)
	case time.Time:
		event.SetEvaDate(fieldId, v)
	default:
		event.SetEvaString(fieldId, fmt.Sprintf("%v", v))
	}
}

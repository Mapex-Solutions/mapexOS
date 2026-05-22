package services

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	triggerDtos "triggers/src/modules/triggers/application/dtos"

	eventsContracts "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
	triggersContracts "github.com/Mapex-Solutions/MapexOS/contracts/services/triggers/triggers"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// placeholderRegex matches {{path.to.field}} tokens for placeholder resolution.
var placeholderRegex = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// publishTriggerEvent publishes trigger execution event to events service for logging.
func (s *EventService) publishTriggerEvent(
	ctx context.Context,
	event *triggerDtos.TriggerExecuteEvent,
	triggerName, triggerType, category string,
	startTime time.Time,
	success bool,
	errorMsg string,
	requestData map[string]interface{},
) {
	// Calculate duration
	durationMs := time.Since(startTime).Milliseconds()

	// Convert request data to JSON string
	requestDataJSON := ""
	if requestData != nil {
		if jsonBytes, err := json.Marshal(requestData); err == nil {
			requestDataJSON = string(jsonBytes)
		}
	}

	// Build trigger event DTO
	triggerEventDTO := eventsContracts.TriggerEventDTO{
		Created:        time.Now().UTC(),
		EventTrackerId: event.EventTrackerId, // UUID for end-to-end event tracking across services
		OrgId:          event.OrgID,
		PathKey:        event.PathKey,
		TriggerId:      event.TriggerID,
		TriggerName:    triggerName,
		TriggerType:    triggerType,
		Category:       category,
		Source:         event.Source,
		Success:        success,
		DurationMs:     durationMs,
		Error:          errorMsg,
		RequestData:    requestDataJSON,
		ResponseData:   "", // Could be populated by executor in future
	}

	// MsgId = {eventTrackerId}-triggerlog for JetStream dedup
	var msgId string
	if event.EventTrackerId != "" {
		msgId = fmt.Sprintf("%s-triggerlog", event.EventTrackerId)
	}

	// Publish to events service using CorePublisher (fire-and-forget)
	publishStart := time.Now()
	if err := s.deps.NatsBus.PublishCore(natsModel.PublishCoreConfig{
		Subject: triggersContracts.SubjectTriggerEvents,
		Data:    triggerEventDTO,
		MsgId:   msgId,
	}); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] Failed to publish trigger event: %s", event.TriggerID))
		s.deps.Metrics.EventsPublished.WithLabelValues("error").Inc()
		s.deps.Metrics.PublishDuration.Observe(time.Since(publishStart).Seconds())
	} else {
		logger.Debug(fmt.Sprintf("[SERVICE:Event] Published trigger event: triggerId=%s, eventTrackerId=%s, success=%t, duration=%dms",
			event.TriggerID, event.EventTrackerId, success, durationMs))
		s.deps.Metrics.EventsPublished.WithLabelValues("ok").Inc()
		s.deps.Metrics.PublishDuration.Observe(time.Since(publishStart).Seconds())
	}
}

// navigatePath traverses a nested map/slice using dot-notation path.
func navigatePath(data interface{}, path string) interface{} {
	if path == "" {
		return data
	}
	parts := splitDotPath(path)
	current := data
	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[part]
		default:
			return current
		}
	}
	return current
}

// splitDotPath splits a dot-notation path into parts.
func splitDotPath(path string) []string {
	result := []string{}
	current := ""
	for _, c := range path {
		if c == '.' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// ResolvePlaceholders replaces all {{placeholder}} patterns in the input string
// with their corresponding values from the data map.
//
// Parameters:
//   - input: String containing placeholders (e.g., "Temperature is {{payload.temp}}°C")
//   - data: Map containing the data to extract values from
//
// Returns:
//   - string: Input string with all placeholders replaced
//   - error: If a placeholder references a non-existent field
func ResolvePlaceholders(input string, data map[string]interface{}) (string, error) {
	// Find all placeholders in the input
	matches := placeholderRegex.FindAllStringSubmatch(input, -1)

	result := input

	// Process each placeholder
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		fullPlaceholder := match[0] // e.g., "{{payload.temperature}}"
		fieldPath := match[1]       // e.g., "payload.temperature"

		// Extract value from data using field path
		value, err := extractValue(fieldPath, data)
		if err != nil {
			return "", fmt.Errorf("failed to resolve placeholder %s: %w", fullPlaceholder, err)
		}

		// Convert value to string
		valueStr := fmt.Sprintf("%v", value)

		// Replace placeholder with value
		result = strings.Replace(result, fullPlaceholder, valueStr, 1)
	}

	return result, nil
}

// ResolvePlaceholdersInMap resolves all placeholders in a map recursively.
//
// This function is useful for resolving placeholders in trigger config objects
// where values can be nested maps or arrays.
//
// Parameters:
//   - input: Map containing values with placeholders
//   - data: Map containing the data to extract values from
//
// Returns:
//   - map[string]interface{}: Input map with all placeholders resolved
//   - error: If any placeholder resolution fails
func ResolvePlaceholdersInMap(input map[string]interface{}, data map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for key, value := range input {
		switch v := value.(type) {
		case string:
			// Resolve placeholders in string value
			resolved, err := ResolvePlaceholders(v, data)
			if err != nil {
				return nil, err
			}
			result[key] = resolved

		case map[string]interface{}:
			// Recursively resolve nested map
			resolved, err := ResolvePlaceholdersInMap(v, data)
			if err != nil {
				return nil, err
			}
			result[key] = resolved

		case []interface{}:
			// Resolve array elements
			resolvedArray := make([]interface{}, len(v))
			for i, item := range v {
				switch itemVal := item.(type) {
				case string:
					resolved, err := ResolvePlaceholders(itemVal, data)
					if err != nil {
						return nil, err
					}
					resolvedArray[i] = resolved
				case map[string]interface{}:
					resolved, err := ResolvePlaceholdersInMap(itemVal, data)
					if err != nil {
						return nil, err
					}
					resolvedArray[i] = resolved
				default:
					resolvedArray[i] = item
				}
			}
			result[key] = resolvedArray

		default:
			// Non-string, non-map, non-array values pass through unchanged
			result[key] = value
		}
	}

	return result, nil
}

// TriggerConfigToMap converts a TriggerConfig struct to map[string]interface{}.
// This is needed for placeholder resolution which works with maps.
//
// Parameters:
//   - config: Pointer to TriggerConfig struct
//
// Returns:
//   - map[string]interface{}: Converted config map
//   - error: If conversion fails
func TriggerConfigToMap(config *triggersContracts.TriggerConfig) (map[string]interface{}, error) {
	if config == nil {
		return make(map[string]interface{}), nil
	}

	// Convert to JSON and back to map[string]interface{}
	// This handles all nested structures correctly
	jsonData, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TriggerConfig: %w", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal to map: %w", err)
	}

	return result, nil
}

// extractValue navigates a nested map using a dot-separated path and returns the value.
//
// Parameters:
//   - path: Dot-separated field path (e.g., "payload.sensor.temperature")
//   - data: Map to extract value from
//
// Returns:
//   - interface{}: The extracted value
//   - error: If the path doesn't exist in the data
func extractValue(path string, data map[string]interface{}) (interface{}, error) {
	// Split path by dots
	parts := strings.Split(path, ".")

	// Navigate through nested maps
	var current interface{} = data
	for _, part := range parts {
		// Trim whitespace
		part = strings.TrimSpace(part)

		// Check if current value is a map
		currentMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("field '%s' is not a map, cannot navigate to '%s'", part, path)
		}

		// Get value from map
		value, exists := currentMap[part]
		if !exists {
			return nil, fmt.Errorf("field '%s' not found in path '%s'", part, path)
		}

		current = value
	}

	return current, nil
}

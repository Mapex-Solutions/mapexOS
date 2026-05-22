package services

import (
	"fmt"
	"strings"

	"workflow/src/modules/fetch_options/application/types"
)

// extractOptions extracts [{label, value}] from a JSON response using dot-path expressions.
func extractOptions(data interface{}, dataPath, valuePath, labelPath string) ([]types.FetchOptionsItem, error) {
	arr := data
	if dataPath != "" {
		arr = navigatePath(data, dataPath)
	}

	slice, ok := arr.([]interface{})
	if !ok {
		return nil, fmt.Errorf("dataPath '%s' did not resolve to an array", dataPath)
	}

	items := make([]types.FetchOptionsItem, 0, len(slice))
	for _, item := range slice {
		value := navigatePath(item, valuePath)
		label := navigatePath(item, labelPath)
		items = append(items, types.FetchOptionsItem{
			Label: fmt.Sprintf("%v", label),
			Value: value,
		})
	}
	return items, nil
}

// navigatePath traverses a dot-separated path through nested maps.
func navigatePath(data interface{}, path string) interface{} {
	if path == "" {
		return data
	}
	parts := strings.Split(path, ".")
	current := data
	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current, ok = m[part]
		if !ok {
			return nil
		}
	}
	return current
}

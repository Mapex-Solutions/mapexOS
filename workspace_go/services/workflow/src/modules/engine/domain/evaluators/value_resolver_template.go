package evaluators

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/*
 * VALUE RESOLVER — TEMPLATE INTERPOLATION
 *
 * Helper for the literal source kind: when a FieldValue.Value contains
 * {{namespace.path}} placeholders, this file resolves them against the four
 * runtime context maps (event, state, input, output).
 *
 * Semantics differ from the per-source Resolve branches: this path is
 * best-effort and never returns an error. Edge cases degrade gracefully:
 *   - Missing path                  → empty string at the placeholder.
 *   - Scalar (number, bool, etc.)   → fmt.Sprintf("%v", v).
 *   - Map / slice                   → encoding/json.Marshal + warn log.
 *   - Marshal failure               → empty string + warn log.
 *   - Malformed `{{` without `}}`   → verbatim return.
 *
 * Namespaces are UI-aligned: event / state / input / output. Intentionally
 * distinct from the plugin executor's namespace shape (wf.state, wf.input,
 * etc.) which stays unchanged for plugin manifest back-compat.
 */

// templatePlaceholderPattern matches {{namespace.path.subpath}} placeholders.
// The path starts with a letter and contains alphanumerics, underscores, and
// dots only. Whitespace around the path is tolerated.
var templatePlaceholderPattern = regexp.MustCompile(`\{\{\s*([a-zA-Z][a-zA-Z0-9_.]*?)\s*\}\}`)

// renderTemplate resolves {{namespace.path}} placeholders in value against the
// four runtime context maps. Best-effort: never returns an error; missing or
// malformed placeholders degrade to empty string or verbatim.
func renderTemplate(
	value string,
	eventPayload map[string]interface{},
	state map[string]interface{},
	externalInputs map[string]interface{},
	nodeOutputs map[string]interface{},
) string {
	if !strings.Contains(value, "{{") {
		return value
	}

	contexts := map[string]interface{}{
		"event":  eventPayload,
		"state":  state,
		"input":  externalInputs,
		"output": nodeOutputs,
	}

	return templatePlaceholderPattern.ReplaceAllStringFunc(value, func(match string) string {
		// match = "{{ event.user.name }}" → trim braces + whitespace.
		path := strings.TrimSpace(match[2 : len(match)-2])

		resolved, found := navigatePath(path, contexts)
		if !found || resolved == nil {
			return ""
		}

		switch v := resolved.(type) {
		case string:
			return v
		case map[string]interface{}, []interface{}:
			return stringifyObject(path, v)
		default:
			return fmt.Sprintf("%v", v)
		}
	})
}

// navigatePath walks the dot-notation path through the contexts map.
// Supports map-key access ("a.b.c") and 0-based array indexing ("a.0.b" — the
// "0" segment indexes into a []interface{} value). Returns (value, true) on
// full match, (nil, false) when any segment is missing, when a non-leaf segment
// is neither a map nor a slice (or the slice index is out of range / not a
// non-negative integer), or when the top-level namespace value is nil.
func navigatePath(path string, contexts map[string]interface{}) (interface{}, bool) {
	parts := strings.Split(path, ".")
	if len(parts) == 0 || parts[0] == "" {
		return nil, false
	}

	current, ok := contexts[parts[0]]
	if !ok || current == nil {
		return nil, false
	}

	for _, part := range parts[1:] {
		if part == "" {
			return nil, false
		}
		switch v := current.(type) {
		case map[string]interface{}:
			next, exists := v[part]
			if !exists {
				return nil, false
			}
			current = next
		case []interface{}:
			idx, err := strconv.Atoi(part)
			if err != nil || idx < 0 || idx >= len(v) {
				return nil, false
			}
			current = v[idx]
		default:
			return nil, false
		}
	}

	return current, true
}

// stringifyObject converts a map or slice value to its JSON representation.
// Emits a warn-level log line so authors notice unintended object templating.
// Marshal failure → empty string.
func stringifyObject(path string, value interface{}) string {
	data, err := json.Marshal(value)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:ValueResolver] template object marshal failed for path %q: %v", path, err))
		return ""
	}
	logger.Warn(fmt.Sprintf("[SERVICE:ValueResolver] template object stringified to JSON for path %q", path))
	return string(data)
}

package services

import (
	"testing"

	triggers "github.com/Mapex-Solutions/MapexOS/contracts/services/triggers/triggers"
)

/**
 * ResolvePlaceholders Tests
 */

func TestResolvePlaceholders_Simple(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		data     map[string]interface{}
		expected string
		hasError bool
	}{
		{
			name:  "single placeholder",
			input: "Hello {{name}}!",
			data: map[string]interface{}{
				"name": "World",
			},
			expected: "Hello World!",
			hasError: false,
		},
		{
			name:  "multiple placeholders",
			input: "{{greeting}} {{name}}!",
			data: map[string]interface{}{
				"greeting": "Hello",
				"name":     "World",
			},
			expected: "Hello World!",
			hasError: false,
		},
		{
			name:     "no placeholders",
			input:    "Hello World!",
			data:     map[string]interface{}{},
			expected: "Hello World!",
			hasError: false,
		},
		{
			name:     "empty string",
			input:    "",
			data:     map[string]interface{}{},
			expected: "",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ResolvePlaceholders(tt.input, tt.data)

			if tt.hasError && err == nil {
				t.Error("ResolvePlaceholders() expected error, got nil")
			}
			if !tt.hasError && err != nil {
				t.Errorf("ResolvePlaceholders() unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("ResolvePlaceholders() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestResolvePlaceholders_NestedPaths(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		data     map[string]interface{}
		expected string
		hasError bool
	}{
		{
			name:  "single level nested",
			input: "Value: {{payload.temperature}}",
			data: map[string]interface{}{
				"payload": map[string]interface{}{
					"temperature": 42.5,
				},
			},
			expected: "Value: 42.5",
			hasError: false,
		},
		{
			name:  "two level nested",
			input: "Sensor: {{payload.sensor.name}}",
			data: map[string]interface{}{
				"payload": map[string]interface{}{
					"sensor": map[string]interface{}{
						"name": "SENSOR-001",
					},
				},
			},
			expected: "Sensor: SENSOR-001",
			hasError: false,
		},
		{
			name:  "three level nested",
			input: "Location: {{payload.device.location.city}}",
			data: map[string]interface{}{
				"payload": map[string]interface{}{
					"device": map[string]interface{}{
						"location": map[string]interface{}{
							"city": "São Paulo",
						},
					},
				},
			},
			expected: "Location: São Paulo",
			hasError: false,
		},
		{
			name:  "mixed root and nested",
			input: "Event {{eventId}}: Sensor {{payload.sensor.id}} at {{payload.timestamp}}",
			data: map[string]interface{}{
				"eventId": "EVT-123",
				"payload": map[string]interface{}{
					"sensor": map[string]interface{}{
						"id": "S-001",
					},
					"timestamp": "2024-01-01T00:00:00Z",
				},
			},
			expected: "Event EVT-123: Sensor S-001 at 2024-01-01T00:00:00Z",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ResolvePlaceholders(tt.input, tt.data)

			if tt.hasError && err == nil {
				t.Error("ResolvePlaceholders() expected error, got nil")
			}
			if !tt.hasError && err != nil {
				t.Errorf("ResolvePlaceholders() unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("ResolvePlaceholders() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestResolvePlaceholders_DifferentTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		data     map[string]interface{}
		expected string
	}{
		{
			name:  "string value",
			input: "{{value}}",
			data: map[string]interface{}{
				"value": "hello",
			},
			expected: "hello",
		},
		{
			name:  "integer value",
			input: "Count: {{value}}",
			data: map[string]interface{}{
				"value": 42,
			},
			expected: "Count: 42",
		},
		{
			name:  "float value",
			input: "Temperature: {{value}}°C",
			data: map[string]interface{}{
				"value": 36.5,
			},
			expected: "Temperature: 36.5°C",
		},
		{
			name:  "boolean value",
			input: "Active: {{value}}",
			data: map[string]interface{}{
				"value": true,
			},
			expected: "Active: true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ResolvePlaceholders(tt.input, tt.data)

			if err != nil {
				t.Errorf("ResolvePlaceholders() unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("ResolvePlaceholders() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestResolvePlaceholders_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		data  map[string]interface{}
	}{
		{
			name:  "field not found",
			input: "{{nonexistent}}",
			data:  map[string]interface{}{},
		},
		{
			name:  "nested field not found",
			input: "{{payload.nonexistent}}",
			data: map[string]interface{}{
				"payload": map[string]interface{}{},
			},
		},
		{
			name:  "intermediate path not a map",
			input: "{{payload.sensor.id}}",
			data: map[string]interface{}{
				"payload": map[string]interface{}{
					"sensor": "not a map",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ResolvePlaceholders(tt.input, tt.data)

			if err == nil {
				t.Error("ResolvePlaceholders() expected error, got nil")
			}
		})
	}
}

/**
 * ResolvePlaceholdersInMap Tests
 */

func TestResolvePlaceholdersInMap_SimpleMap(t *testing.T) {
	input := map[string]interface{}{
		"message": "Alert: {{payload.message}}",
		"count":   42,
		"active":  true,
	}

	data := map[string]interface{}{
		"payload": map[string]interface{}{
			"message": "Temperature exceeded threshold",
		},
	}

	result, err := ResolvePlaceholdersInMap(input, data)

	if err != nil {
		t.Fatalf("ResolvePlaceholdersInMap() unexpected error: %v", err)
	}

	if result["message"] != "Alert: Temperature exceeded threshold" {
		t.Errorf("ResolvePlaceholdersInMap() message = %q, want 'Alert: Temperature exceeded threshold'", result["message"])
	}

	// Non-string values should pass through unchanged
	if result["count"] != 42 {
		t.Errorf("ResolvePlaceholdersInMap() count = %v, want 42", result["count"])
	}
	if result["active"] != true {
		t.Errorf("ResolvePlaceholdersInMap() active = %v, want true", result["active"])
	}
}

func TestResolvePlaceholdersInMap_NestedMap(t *testing.T) {
	input := map[string]interface{}{
		"http": map[string]interface{}{
			"endpoint": "https://api.example.com/{{payload.endpoint}}",
			"body": map[string]interface{}{
				"sensor":  "{{payload.sensorId}}",
				"value":   "{{payload.value}}",
				"static":  "unchanged",
				"numeric": 100,
			},
		},
	}

	data := map[string]interface{}{
		"payload": map[string]interface{}{
			"endpoint": "webhook",
			"sensorId": "SENSOR-001",
			"value":    42.5,
		},
	}

	result, err := ResolvePlaceholdersInMap(input, data)

	if err != nil {
		t.Fatalf("ResolvePlaceholdersInMap() unexpected error: %v", err)
	}

	http := result["http"].(map[string]interface{})
	if http["endpoint"] != "https://api.example.com/webhook" {
		t.Errorf("endpoint = %q, want 'https://api.example.com/webhook'", http["endpoint"])
	}

	body := http["body"].(map[string]interface{})
	if body["sensor"] != "SENSOR-001" {
		t.Errorf("sensor = %q, want 'SENSOR-001'", body["sensor"])
	}
	if body["static"] != "unchanged" {
		t.Errorf("static = %q, want 'unchanged'", body["static"])
	}
	if body["numeric"] != 100 {
		t.Errorf("numeric = %v, want 100", body["numeric"])
	}
}

func TestResolvePlaceholdersInMap_ArrayValues(t *testing.T) {
	input := map[string]interface{}{
		"recipients": []interface{}{
			"{{payload.email1}}",
			"{{payload.email2}}",
			"static@example.com",
		},
	}

	data := map[string]interface{}{
		"payload": map[string]interface{}{
			"email1": "user1@example.com",
			"email2": "user2@example.com",
		},
	}

	result, err := ResolvePlaceholdersInMap(input, data)

	if err != nil {
		t.Fatalf("ResolvePlaceholdersInMap() unexpected error: %v", err)
	}

	recipients := result["recipients"].([]interface{})
	if len(recipients) != 3 {
		t.Fatalf("recipients length = %d, want 3", len(recipients))
	}

	if recipients[0] != "user1@example.com" {
		t.Errorf("recipients[0] = %q, want 'user1@example.com'", recipients[0])
	}
	if recipients[1] != "user2@example.com" {
		t.Errorf("recipients[1] = %q, want 'user2@example.com'", recipients[1])
	}
	if recipients[2] != "static@example.com" {
		t.Errorf("recipients[2] = %q, want 'static@example.com'", recipients[2])
	}
}

func TestResolvePlaceholdersInMap_ArrayWithMaps(t *testing.T) {
	input := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{
				"name":  "{{payload.item1.name}}",
				"value": "{{payload.item1.value}}",
			},
			map[string]interface{}{
				"name":  "{{payload.item2.name}}",
				"value": "{{payload.item2.value}}",
			},
		},
	}

	data := map[string]interface{}{
		"payload": map[string]interface{}{
			"item1": map[string]interface{}{
				"name":  "Item One",
				"value": "100",
			},
			"item2": map[string]interface{}{
				"name":  "Item Two",
				"value": "200",
			},
		},
	}

	result, err := ResolvePlaceholdersInMap(input, data)

	if err != nil {
		t.Fatalf("ResolvePlaceholdersInMap() unexpected error: %v", err)
	}

	items := result["items"].([]interface{})
	item1 := items[0].(map[string]interface{})
	item2 := items[1].(map[string]interface{})

	if item1["name"] != "Item One" {
		t.Errorf("item1.name = %q, want 'Item One'", item1["name"])
	}
	if item2["value"] != "200" {
		t.Errorf("item2.value = %q, want '200'", item2["value"])
	}
}

/**
 * TriggerConfigToMap Tests
 */

func TestTriggerConfigToMap_NilConfig(t *testing.T) {
	result, err := TriggerConfigToMap(nil)

	if err != nil {
		t.Fatalf("TriggerConfigToMap() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("TriggerConfigToMap() returned nil, want empty map")
	}

	if len(result) != 0 {
		t.Errorf("TriggerConfigToMap() returned map with %d items, want 0", len(result))
	}
}

func TestTriggerConfigToMap_HttpConfig(t *testing.T) {
	endpoint := "https://api.example.com"
	method := "POST"

	config := &triggers.TriggerConfig{
		Http: &triggers.HttpConfig{
			Endpoint: endpoint,
			Method:   method,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: map[string]interface{}{
				"message": "test",
			},
		},
	}

	result, err := TriggerConfigToMap(config)

	if err != nil {
		t.Fatalf("TriggerConfigToMap() unexpected error: %v", err)
	}

	http, ok := result["http"].(map[string]interface{})
	if !ok {
		t.Fatal("TriggerConfigToMap() should contain 'http' field")
	}

	if http["endpoint"] != endpoint {
		t.Errorf("endpoint = %v, want %q", http["endpoint"], endpoint)
	}
	if http["method"] != method {
		t.Errorf("method = %v, want %q", http["method"], method)
	}
}

func TestTriggerConfigToMap_EmailConfig(t *testing.T) {
	to := "test@example.com"
	subject := "Test Subject"

	config := &triggers.TriggerConfig{
		Email: &triggers.EmailConfig{
			To:      to,
			Subject: subject,
		},
	}

	result, err := TriggerConfigToMap(config)

	if err != nil {
		t.Fatalf("TriggerConfigToMap() unexpected error: %v", err)
	}

	email, ok := result["email"].(map[string]interface{})
	if !ok {
		t.Fatal("TriggerConfigToMap() should contain 'email' field")
	}

	if email["to"] != to {
		t.Errorf("to = %v, want %q", email["to"], to)
	}
	if email["subject"] != subject {
		t.Errorf("subject = %v, want %q", email["subject"], subject)
	}
}

/**
 * extractValue Tests
 */

func TestExtractValue_RootLevel(t *testing.T) {
	data := map[string]interface{}{
		"name":  "test",
		"value": 42,
	}

	result, err := extractValue("name", data)

	if err != nil {
		t.Fatalf("extractValue() unexpected error: %v", err)
	}

	if result != "test" {
		t.Errorf("extractValue() = %v, want 'test'", result)
	}
}

func TestExtractValue_NestedLevel(t *testing.T) {
	data := map[string]interface{}{
		"payload": map[string]interface{}{
			"sensor": map[string]interface{}{
				"temperature": 42.5,
			},
		},
	}

	result, err := extractValue("payload.sensor.temperature", data)

	if err != nil {
		t.Fatalf("extractValue() unexpected error: %v", err)
	}

	if result != 42.5 {
		t.Errorf("extractValue() = %v, want 42.5", result)
	}
}

func TestExtractValue_NotFound(t *testing.T) {
	data := map[string]interface{}{
		"name": "test",
	}

	_, err := extractValue("nonexistent", data)

	if err == nil {
		t.Fatal("extractValue() expected error for non-existent field")
	}
}

func TestExtractValue_IntermediateNotMap(t *testing.T) {
	data := map[string]interface{}{
		"payload": "not a map",
	}

	_, err := extractValue("payload.field", data)

	if err == nil {
		t.Fatal("extractValue() expected error when intermediate path is not a map")
	}
}

/**
 * Benchmarks
 */

func BenchmarkResolvePlaceholders_Simple(b *testing.B) {
	input := "Hello {{name}}!"
	data := map[string]interface{}{
		"name": "World",
	}

	for i := 0; i < b.N; i++ {
		ResolvePlaceholders(input, data)
	}
}

func BenchmarkResolvePlaceholders_Nested(b *testing.B) {
	input := "Sensor {{payload.sensor.id}} reported {{payload.sensor.temperature}}°C"
	data := map[string]interface{}{
		"payload": map[string]interface{}{
			"sensor": map[string]interface{}{
				"id":          "SENSOR-001",
				"temperature": 42.5,
			},
		},
	}

	for i := 0; i < b.N; i++ {
		ResolvePlaceholders(input, data)
	}
}

func BenchmarkResolvePlaceholdersInMap(b *testing.B) {
	input := map[string]interface{}{
		"http": map[string]interface{}{
			"endpoint": "https://api.example.com/{{payload.endpoint}}",
			"body": map[string]interface{}{
				"sensor": "{{payload.sensorId}}",
				"value":  "{{payload.value}}",
			},
		},
	}

	data := map[string]interface{}{
		"payload": map[string]interface{}{
			"endpoint": "webhook",
			"sensorId": "SENSOR-001",
			"value":    42.5,
		},
	}

	for i := 0; i < b.N; i++ {
		ResolvePlaceholdersInMap(input, data)
	}
}

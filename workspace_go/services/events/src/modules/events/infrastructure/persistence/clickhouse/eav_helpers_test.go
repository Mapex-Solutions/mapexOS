package clickhouseRepo

import (
	"strings"
	"testing"
	"time"
)

// TestEAVQueryBuilder_NumberFilters tests number field filtering
func TestEAVQueryBuilder_NumberFilters(t *testing.T) {
	tests := []struct {
		name           string
		builderFunc    func(*EAVQueryBuilder)
		expectedClause string
		expectedParams int
	}{
		{
			name: "NumberEquals",
			builderFunc: func(b *EAVQueryBuilder) {
				b.NumberEquals("temperature", 25.5)
			},
			expectedClause: "arrayExists(x -> x.1 = 'temperature' AND x.2 = ?, numberFields)",
			expectedParams: 1,
		},
		{
			name: "NumberGreaterThan",
			builderFunc: func(b *EAVQueryBuilder) {
				b.NumberGreaterThan("temperature", 20)
			},
			expectedClause: "arrayExists(x -> x.1 = 'temperature' AND x.2 > ?, numberFields)",
			expectedParams: 1,
		},
		{
			name: "NumberBetween",
			builderFunc: func(b *EAVQueryBuilder) {
				b.NumberBetween("temperature", 20, 30)
			},
			expectedClause: "arrayExists(x -> x.1 = 'temperature' AND x.2 >= ? AND x.2 <= ?, numberFields)",
			expectedParams: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewEAVQueryBuilder()
			tt.builderFunc(builder)

			whereClause, params := builder.Build()

			if whereClause != tt.expectedClause {
				t.Errorf("Expected clause: %s, got: %s", tt.expectedClause, whereClause)
			}

			if len(params) != tt.expectedParams {
				t.Errorf("Expected %d params, got %d", tt.expectedParams, len(params))
			}
		})
	}
}

// TestEAVQueryBuilder_StringFilters tests string field filtering
func TestEAVQueryBuilder_StringFilters(t *testing.T) {
	tests := []struct {
		name           string
		builderFunc    func(*EAVQueryBuilder)
		expectedClause string
		expectedParams int
	}{
		{
			name: "StringEquals",
			builderFunc: func(b *EAVQueryBuilder) {
				b.StringEquals("location", "warehouse_A")
			},
			expectedClause: "arrayExists(x -> x.1 = 'location' AND x.2 = ?, stringFields)",
			expectedParams: 1,
		},
		{
			name: "StringContains",
			builderFunc: func(b *EAVQueryBuilder) {
				b.StringContains("deviceName", "sensor")
			},
			expectedClause: "arrayExists(x -> x.1 = 'deviceName' AND positionCaseInsensitive(x.2, ?) > 0, stringFields)",
			expectedParams: 1,
		},
		{
			name: "StringIn",
			builderFunc: func(b *EAVQueryBuilder) {
				b.StringIn("status", []string{"online", "active"})
			},
			expectedClause: "arrayExists(x -> x.1 = 'status' AND x.2 IN (?, ?), stringFields)",
			expectedParams: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewEAVQueryBuilder()
			tt.builderFunc(builder)

			whereClause, params := builder.Build()

			if whereClause != tt.expectedClause {
				t.Errorf("Expected clause: %s, got: %s", tt.expectedClause, whereClause)
			}

			if len(params) != tt.expectedParams {
				t.Errorf("Expected %d params, got %d", tt.expectedParams, len(params))
			}
		})
	}
}

// TestEAVQueryBuilder_BoolFilters tests boolean field filtering
func TestEAVQueryBuilder_BoolFilters(t *testing.T) {
	builder := NewEAVQueryBuilder()
	builder.BoolEquals("isOnline", true)

	whereClause, params := builder.Build()

	expected := "arrayExists(x -> x.1 = 'isOnline' AND x.2 = ?, boolFields)"
	if whereClause != expected {
		t.Errorf("Expected clause: %s, got: %s", expected, whereClause)
	}

	if len(params) != 1 {
		t.Errorf("Expected 1 param, got %d", len(params))
	}

	if params[0].(uint8) != 1 {
		t.Errorf("Expected param to be 1, got %v", params[0])
	}
}

// TestEAVQueryBuilder_DateFilters tests date field filtering
func TestEAVQueryBuilder_DateFilters(t *testing.T) {
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)

	builder := NewEAVQueryBuilder()
	builder.DateAfter("lastSeen", oneHourAgo)

	whereClause, params := builder.Build()

	expected := "arrayExists(x -> x.1 = 'lastSeen' AND x.2 > ?, dateFields)"
	if whereClause != expected {
		t.Errorf("Expected clause: %s, got: %s", expected, whereClause)
	}

	if len(params) != 1 {
		t.Errorf("Expected 1 param, got %d", len(params))
	}
}

// TestEAVQueryBuilder_ExistenceChecks tests field existence checking
func TestEAVQueryBuilder_ExistenceChecks(t *testing.T) {
	builder := NewEAVQueryBuilder()
	builder.HasNumberField("temperature").
		HasStringField("location")

	whereClause, params := builder.Build()

	if !strings.Contains(whereClause, "arrayExists(x -> x.1 = 'temperature', numberFields)") {
		t.Errorf("Expected temperature existence check in clause: %s", whereClause)
	}

	if !strings.Contains(whereClause, "arrayExists(x -> x.1 = 'location', stringFields)") {
		t.Errorf("Expected location existence check in clause: %s", whereClause)
	}

	if len(params) != 0 {
		t.Errorf("Expected 0 params for existence checks, got %d", len(params))
	}
}

// TestEAVQueryBuilder_ChainedConditions tests chaining multiple conditions
func TestEAVQueryBuilder_ChainedConditions(t *testing.T) {
	builder := NewEAVQueryBuilder()
	builder.NumberBetween("temperature", 20, 30).
		StringEquals("location", "warehouse_A").
		BoolEquals("isOnline", true)

	whereClause, params := builder.Build()

	// Should contain all three conditions joined by AND
	if !strings.Contains(whereClause, "temperature") {
		t.Error("Missing temperature condition")
	}
	if !strings.Contains(whereClause, "location") {
		t.Error("Missing location condition")
	}
	if !strings.Contains(whereClause, "isOnline") {
		t.Error("Missing isOnline condition")
	}

	// Check AND separators
	andCount := strings.Count(whereClause, " AND ")
	if andCount < 2 {
		t.Errorf("Expected at least 2 AND operators, got %d", andCount)
	}

	// Should have 4 params: 2 for temperature range, 1 for location, 1 for bool
	if len(params) != 4 {
		t.Errorf("Expected 4 params, got %d", len(params))
	}
}

// TestEAVQueryBuilder_BuildWithPrefix tests BuildWithPrefix method
func TestEAVQueryBuilder_BuildWithPrefix(t *testing.T) {
	// Test with conditions
	builder := NewEAVQueryBuilder()
	builder.NumberGreaterThan("temperature", 25)

	whereClause, params := builder.BuildWithPrefix()

	if !strings.HasPrefix(whereClause, "WHERE ") {
		t.Error("Expected clause to start with 'WHERE '")
	}

	if len(params) != 1 {
		t.Errorf("Expected 1 param, got %d", len(params))
	}

	// Test without conditions
	emptyBuilder := NewEAVQueryBuilder()
	emptyClause, emptyParams := emptyBuilder.BuildWithPrefix()

	if emptyClause != "" {
		t.Errorf("Expected empty clause, got: %s", emptyClause)
	}

	if emptyParams != nil {
		t.Error("Expected nil params for empty builder")
	}
}

// TestExtractFieldFunctions tests field extraction helper functions
func TestExtractFieldFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function func(string) string
		field    string
		expected string
	}{
		{
			name:     "ExtractNumberField",
			function: ExtractNumberField,
			field:    "temperature",
			expected: "arrayFirst(x -> x.1 = 'temperature', numberFields).2",
		},
		{
			name:     "ExtractStringField",
			function: ExtractStringField,
			field:    "location",
			expected: "arrayFirst(x -> x.1 = 'location', stringFields).2",
		},
		{
			name:     "ExtractBoolField",
			function: ExtractBoolField,
			field:    "isOnline",
			expected: "arrayFirst(x -> x.1 = 'isOnline', boolFields).2",
		},
		{
			name:     "ExtractDateField",
			function: ExtractDateField,
			field:    "lastSeen",
			expected: "arrayFirst(x -> x.1 = 'lastSeen', dateFields).2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(tt.field)
			if result != tt.expected {
				t.Errorf("Expected: %s, got: %s", tt.expected, result)
			}
		})
	}
}

// TestAggregationBuilder tests the aggregation builder
func TestAggregationBuilder(t *testing.T) {
	aggBuilder := NewAggregationBuilder("temperature", "number")
	aggBuilder.GroupByStringField("location").
		Avg("avg_temp").
		Min("min_temp").
		Max("max_temp").
		Count("event_count")

	selectClause := aggBuilder.BuildSelect()
	groupByClause := aggBuilder.BuildGroupBy()

	// Check SELECT clause contains all aggregations
	if !strings.Contains(selectClause, "avg_temp") {
		t.Error("Missing avg_temp in SELECT")
	}
	if !strings.Contains(selectClause, "min_temp") {
		t.Error("Missing min_temp in SELECT")
	}
	if !strings.Contains(selectClause, "max_temp") {
		t.Error("Missing max_temp in SELECT")
	}
	if !strings.Contains(selectClause, "event_count") {
		t.Error("Missing event_count in SELECT")
	}

	// Check GROUP BY clause
	if !strings.HasPrefix(groupByClause, "GROUP BY ") {
		t.Error("GROUP BY clause should start with 'GROUP BY '")
	}

	if !strings.Contains(groupByClause, "location") {
		t.Error("GROUP BY should contain location field")
	}
}

// TestAggregationBuilder_EmptyGroupBy tests aggregation without GROUP BY
func TestAggregationBuilder_EmptyGroupBy(t *testing.T) {
	aggBuilder := NewAggregationBuilder("temperature", "number")
	aggBuilder.Avg("avg_temp").Count("total")

	selectClause := aggBuilder.BuildSelect()
	groupByClause := aggBuilder.BuildGroupBy()

	// Should have SELECT with aggregations
	if !strings.Contains(selectClause, "avg_temp") {
		t.Error("Missing avg_temp in SELECT")
	}

	// Should have empty GROUP BY
	if groupByClause != "" {
		t.Errorf("Expected empty GROUP BY, got: %s", groupByClause)
	}
}

// BenchmarkEAVQueryBuilder benchmarks the query builder
func BenchmarkEAVQueryBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		builder := NewEAVQueryBuilder()
		builder.NumberBetween("temperature", 20, 30).
			StringEquals("location", "warehouse_A").
			BoolEquals("isOnline", true)
		builder.Build()
	}
}

// BenchmarkExtractNumberField benchmarks field extraction
func BenchmarkExtractNumberField(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractNumberField("temperature")
	}
}

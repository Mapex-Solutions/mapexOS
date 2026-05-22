package clickhouseRepo

import (
	"fmt"
	"strings"
	"time"
)

// NewEAVQueryBuilder creates a new query builder for EAV fields.
func NewEAVQueryBuilder() *EAVQueryBuilder {
	return &EAVQueryBuilder{
		conditions: make([]string, 0),
		params:     make([]interface{}, 0),
	}
}

/**
 * Number Field Filters
 */

// NumberEquals adds a condition for a number field equals exact value.
// Example: NumberEquals("temperature", 25.5)
func (b *EAVQueryBuilder) NumberEquals(field string, value float64) *EAVQueryBuilder {
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s' AND x.2 = ?, numberFields)", field)
	b.conditions = append(b.conditions, condition)
	b.params = append(b.params, value)
	return b
}

// NumberGreaterThan adds a condition for a number field > value.
// Example: NumberGreaterThan("temperature", 20)
func (b *EAVQueryBuilder) NumberGreaterThan(field string, value float64) *EAVQueryBuilder {
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s' AND x.2 > ?, numberFields)", field)
	b.conditions = append(b.conditions, condition)
	b.params = append(b.params, value)
	return b
}

// NumberLessThan adds a condition for a number field < value.
// Example: NumberLessThan("temperature", 30)
func (b *EAVQueryBuilder) NumberLessThan(field string, value float64) *EAVQueryBuilder {
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s' AND x.2 < ?, numberFields)", field)
	b.conditions = append(b.conditions, condition)
	b.params = append(b.params, value)
	return b
}

// NumberBetween adds a condition for a number field in range [min, max].
// Example: NumberBetween("temperature", 20, 30)
func (b *EAVQueryBuilder) NumberBetween(field string, min, max float64) *EAVQueryBuilder {
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s' AND x.2 >= ? AND x.2 <= ?, numberFields)", field)
	b.conditions = append(b.conditions, condition)
	b.params = append(b.params, min, max)
	return b
}

/**
 * String Field Filters
 */

// StringEquals adds a condition for a string field equals exact value.
// Example: StringEquals("location", "warehouse_A")
func (b *EAVQueryBuilder) StringEquals(field, value string) *EAVQueryBuilder {
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s' AND x.2 = ?, stringFields)", field)
	b.conditions = append(b.conditions, condition)
	b.params = append(b.params, value)
	return b
}

// StringContains adds a condition for a string field contains substring (case-insensitive).
// Example: StringContains("deviceName", "sensor")
func (b *EAVQueryBuilder) StringContains(field, substring string) *EAVQueryBuilder {
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s' AND positionCaseInsensitive(x.2, ?) > 0, stringFields)", field)
	b.conditions = append(b.conditions, condition)
	b.params = append(b.params, substring)
	return b
}

// StringIn adds a condition for a string field in list of values.
// Example: StringIn("status", []string{"online", "active"})
func (b *EAVQueryBuilder) StringIn(field string, values []string) *EAVQueryBuilder {
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
		b.params = append(b.params, values[i])
	}
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s' AND x.2 IN (%s), stringFields)",
		field, strings.Join(placeholders, ", "))
	b.conditions = append(b.conditions, condition)
	return b
}

/**
 * Boolean Field Filters
 */

// BoolEquals adds a condition for a boolean field.
// Example: BoolEquals("isOnline", true)
func (b *EAVQueryBuilder) BoolEquals(field string, value bool) *EAVQueryBuilder {
	boolValue := uint8(0)
	if value {
		boolValue = 1
	}
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s' AND x.2 = ?, boolFields)", field)
	b.conditions = append(b.conditions, condition)
	b.params = append(b.params, boolValue)
	return b
}

/**
 * Date Field Filters
 */

// DateAfter adds a condition for a date field > timestamp.
// Example: DateAfter("lastSeen", time.Now().Add(-24*time.Hour))
func (b *EAVQueryBuilder) DateAfter(field string, timestamp time.Time) *EAVQueryBuilder {
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s' AND x.2 > ?, dateFields)", field)
	b.conditions = append(b.conditions, condition)
	b.params = append(b.params, timestamp)
	return b
}

// DateBefore adds a condition for a date field < timestamp.
// Example: DateBefore("lastSeen", time.Now())
func (b *EAVQueryBuilder) DateBefore(field string, timestamp time.Time) *EAVQueryBuilder {
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s' AND x.2 < ?, dateFields)", field)
	b.conditions = append(b.conditions, condition)
	b.params = append(b.params, timestamp)
	return b
}

// DateBetween adds a condition for a date field in range [start, end].
// Example: DateBetween("created", startTime, endTime)
func (b *EAVQueryBuilder) DateBetween(field string, start, end time.Time) *EAVQueryBuilder {
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s' AND x.2 >= ? AND x.2 <= ?, dateFields)", field)
	b.conditions = append(b.conditions, condition)
	b.params = append(b.params, start, end)
	return b
}

/**
 * Existence Checks
 */

// HasNumberField checks if event has a specific number field (any value).
// Example: HasNumberField("temperature")
func (b *EAVQueryBuilder) HasNumberField(field string) *EAVQueryBuilder {
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s', numberFields)", field)
	b.conditions = append(b.conditions, condition)
	return b
}

// HasStringField checks if event has a specific string field (any value).
// Example: HasStringField("location")
func (b *EAVQueryBuilder) HasStringField(field string) *EAVQueryBuilder {
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s', stringFields)", field)
	b.conditions = append(b.conditions, condition)
	return b
}

// HasBoolField checks if event has a specific bool field (any value).
// Example: HasBoolField("isActive")
func (b *EAVQueryBuilder) HasBoolField(field string) *EAVQueryBuilder {
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s', boolFields)", field)
	b.conditions = append(b.conditions, condition)
	return b
}

// HasDateField checks if event has a specific date field (any value).
// Example: HasDateField("lastUpdate")
func (b *EAVQueryBuilder) HasDateField(field string) *EAVQueryBuilder {
	condition := fmt.Sprintf("arrayExists(x -> x.1 = '%s', dateFields)", field)
	b.conditions = append(b.conditions, condition)
	return b
}

/**
 * Builder Finalization
 */

// Build returns the WHERE clause and parameters for the query.
// Returns empty string if no conditions were added.
func (b *EAVQueryBuilder) Build() (whereClause string, params []interface{}) {
	if len(b.conditions) == 0 {
		return "", nil
	}
	whereClause = strings.Join(b.conditions, " AND ")
	return whereClause, b.params
}

// BuildWithPrefix returns the WHERE clause with "WHERE" prefix if conditions exist.
func (b *EAVQueryBuilder) BuildWithPrefix() (whereClause string, params []interface{}) {
	clause, params := b.Build()
	if clause == "" {
		return "", params
	}
	return "WHERE " + clause, params
}

/**
 * Field Extraction Helpers
 */

// ExtractNumberField returns SQL to extract a number field value.
// Example: ExtractNumberField("temperature") -> "arrayFirst(x -> x.1 = 'temperature', numberFields).2"
func ExtractNumberField(field string) string {
	return fmt.Sprintf("arrayFirst(x -> x.1 = '%s', numberFields).2", field)
}

// ExtractStringField returns SQL to extract a string field value.
// Example: ExtractStringField("location") -> "arrayFirst(x -> x.1 = 'location', stringFields).2"
func ExtractStringField(field string) string {
	return fmt.Sprintf("arrayFirst(x -> x.1 = '%s', stringFields).2", field)
}

// ExtractBoolField returns SQL to extract a bool field value.
// Example: ExtractBoolField("isOnline") -> "arrayFirst(x -> x.1 = 'isOnline', boolFields).2"
func ExtractBoolField(field string) string {
	return fmt.Sprintf("arrayFirst(x -> x.1 = '%s', boolFields).2", field)
}

// ExtractDateField returns SQL to extract a date field value.
// Example: ExtractDateField("lastSeen") -> "arrayFirst(x -> x.1 = 'lastSeen', dateFields).2"
func ExtractDateField(field string) string {
	return fmt.Sprintf("arrayFirst(x -> x.1 = '%s', dateFields).2", field)
}

/**
 * Aggregation Helpers
 */

// NewAggregationBuilder creates a new aggregation builder.
func NewAggregationBuilder(field, fieldType string) *AggregationBuilder {
	return &AggregationBuilder{
		field:      field,
		fieldType:  fieldType,
		groupBy:    make([]string, 0),
		aggregates: make([]string, 0),
	}
}

// GroupByStringField adds a string field to GROUP BY clause.
func (a *AggregationBuilder) GroupByStringField(field string) *AggregationBuilder {
	a.groupBy = append(a.groupBy, ExtractStringField(field))
	return a
}

// Avg adds AVG aggregation (only for number fields).
func (a *AggregationBuilder) Avg(alias string) *AggregationBuilder {
	if a.fieldType == "number" {
		a.aggregates = append(a.aggregates, fmt.Sprintf("avg(%s) as %s", ExtractNumberField(a.field), alias))
	}
	return a
}

// Min adds MIN aggregation (only for number fields).
func (a *AggregationBuilder) Min(alias string) *AggregationBuilder {
	if a.fieldType == "number" {
		a.aggregates = append(a.aggregates, fmt.Sprintf("min(%s) as %s", ExtractNumberField(a.field), alias))
	}
	return a
}

// Max adds MAX aggregation (only for number fields).
func (a *AggregationBuilder) Max(alias string) *AggregationBuilder {
	if a.fieldType == "number" {
		a.aggregates = append(a.aggregates, fmt.Sprintf("max(%s) as %s", ExtractNumberField(a.field), alias))
	}
	return a
}

// Count adds COUNT aggregation.
func (a *AggregationBuilder) Count(alias string) *AggregationBuilder {
	a.aggregates = append(a.aggregates, fmt.Sprintf("count(*) as %s", alias))
	return a
}

// BuildSelect returns the SELECT clause for aggregations.
func (a *AggregationBuilder) BuildSelect() string {
	parts := make([]string, 0, len(a.groupBy)+len(a.aggregates))

	// Add GROUP BY fields to SELECT
	for _, field := range a.groupBy {
		parts = append(parts, field)
	}

	// Add aggregates to SELECT
	parts = append(parts, a.aggregates...)

	return strings.Join(parts, ", ")
}

// BuildGroupBy returns the GROUP BY clause.
func (a *AggregationBuilder) BuildGroupBy() string {
	if len(a.groupBy) == 0 {
		return ""
	}
	return "GROUP BY " + strings.Join(a.groupBy, ", ")
}

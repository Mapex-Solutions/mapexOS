package clickhouseRepo

import (
	"context"
	"fmt"
	"time"

	"events/src/modules/events/domain/entities"
)

// NewIoTHelpers creates a new IoT helpers instance.
func NewIoTHelpers(repo *EventRepositoryClickHouse) *IoTHelpers {
	return &IoTHelpers{repo: repo}
}

/**
 * Last Value Queries
 */

// GetLastNumberValue retrieves the most recent value for a numeric field.
// Returns the value, created, and error if any.
//
// Example: GetLastNumberValue(ctx, "sensor_123", "temperature")
func (h *IoTHelpers) GetLastNumberValue(ctx context.Context, assetId, field string) (*LastValueResult, error) {
	query := fmt.Sprintf(`
		SELECT
			asset_id,
			%s as value,
			created
		FROM events
		WHERE asset_id = ?
		  AND arrayExists(x -> x.1 = ?, numberFields)
		ORDER BY created DESC
		LIMIT 1
	`, ExtractNumberField(field))

	row := h.repo.conn.QueryRow(ctx, query, assetId, field)

	var result LastValueResult
	var value float64
	err := row.Scan(&result.AssetId, &value, &result.Timestamp)
	if err != nil {
		return nil, err
	}

	result.Value = value
	return &result, nil
}

// GetLastStringValue retrieves the most recent value for a string field.
func (h *IoTHelpers) GetLastStringValue(ctx context.Context, assetId, field string) (*LastValueResult, error) {
	query := fmt.Sprintf(`
		SELECT
			asset_id,
			%s as value,
			created
		FROM events
		WHERE asset_id = ?
		  AND arrayExists(x -> x.1 = ?, stringFields)
		ORDER BY created DESC
		LIMIT 1
	`, ExtractStringField(field))

	row := h.repo.conn.QueryRow(ctx, query, assetId, field)

	var result LastValueResult
	var value string
	err := row.Scan(&result.AssetId, &value, &result.Timestamp)
	if err != nil {
		return nil, err
	}

	result.Value = value
	return &result, nil
}

// GetLastBoolValue retrieves the most recent value for a boolean field.
func (h *IoTHelpers) GetLastBoolValue(ctx context.Context, assetId, field string) (*LastValueResult, error) {
	query := fmt.Sprintf(`
		SELECT
			asset_id,
			%s as value,
			created
		FROM events
		WHERE asset_id = ?
		  AND arrayExists(x -> x.1 = ?, boolFields)
		ORDER BY created DESC
		LIMIT 1
	`, ExtractBoolField(field))

	row := h.repo.conn.QueryRow(ctx, query, assetId, field)

	var result LastValueResult
	var value uint8
	err := row.Scan(&result.AssetId, &value, &result.Timestamp)
	if err != nil {
		return nil, err
	}

	result.Value = value == 1
	return &result, nil
}

/**
 * Time Window Aggregations
 */

// GetAvgInLastHours calculates average of a numeric field in the last N hours.
//
// Example: GetAvgInLastHours(ctx, "sensor_123", "temperature", 24)
// Returns average temperature in the last 24 hours.
func (h *IoTHelpers) GetAvgInLastHours(ctx context.Context, assetId, field string, hours int) (float64, error) {
	query := fmt.Sprintf(`
		SELECT avg(%s) as avg_value
		FROM events
		WHERE asset_id = ?
		  AND created >= now() - INTERVAL ? HOUR
		  AND arrayExists(x -> x.1 = ?, numberFields)
	`, ExtractNumberField(field))

	row := h.repo.conn.QueryRow(ctx, query, assetId, hours, field)

	var avgValue float64
	err := row.Scan(&avgValue)
	return avgValue, err
}

// GetMaxInLastHours finds the maximum value of a numeric field in the last N hours.
//
// Example: GetMaxInLastHours(ctx, "sensor_123", "temperature", 24)
func (h *IoTHelpers) GetMaxInLastHours(ctx context.Context, assetId, field string, hours int) (float64, error) {
	query := fmt.Sprintf(`
		SELECT max(%s) as max_value
		FROM events
		WHERE asset_id = ?
		  AND created >= now() - INTERVAL ? HOUR
		  AND arrayExists(x -> x.1 = ?, numberFields)
	`, ExtractNumberField(field))

	row := h.repo.conn.QueryRow(ctx, query, assetId, hours, field)

	var maxValue float64
	err := row.Scan(&maxValue)
	return maxValue, err
}

// GetMinInLastHours finds the minimum value of a numeric field in the last N hours.
//
// Example: GetMinInLastHours(ctx, "sensor_123", "temperature", 24)
func (h *IoTHelpers) GetMinInLastHours(ctx context.Context, assetId, field string, hours int) (float64, error) {
	query := fmt.Sprintf(`
		SELECT min(%s) as min_value
		FROM events
		WHERE asset_id = ?
		  AND created >= now() - INTERVAL ? HOUR
		  AND arrayExists(x -> x.1 = ?, numberFields)
	`, ExtractNumberField(field))

	row := h.repo.conn.QueryRow(ctx, query, assetId, hours, field)

	var minValue float64
	err := row.Scan(&minValue)
	return minValue, err
}

// GetStatsInLastHours retrieves all statistics (avg, min, max, count) in one query.
// More efficient than calling individual methods.
//
// Example: GetStatsInLastHours(ctx, "sensor_123", "temperature", 24)
func (h *IoTHelpers) GetStatsInLastHours(ctx context.Context, assetId, field string, hours int) (*TimeWindowStats, error) {
	now := time.Now()
	startTime := now.Add(-time.Duration(hours) * time.Hour)

	query := fmt.Sprintf(`
		SELECT
			avg(%s) as avg_value,
			min(%s) as min_value,
			max(%s) as max_value,
			count(*) as count
		FROM events
		WHERE asset_id = ?
		  AND created >= ?
		  AND arrayExists(x -> x.1 = ?, numberFields)
	`, ExtractNumberField(field), ExtractNumberField(field), ExtractNumberField(field))

	row := h.repo.conn.QueryRow(ctx, query, assetId, startTime, field)

	stats := &TimeWindowStats{
		Field:     field,
		StartTime: startTime,
		EndTime:   now,
	}

	err := row.Scan(&stats.Avg, &stats.Min, &stats.Max, &stats.Count)
	return stats, err
}

/**
 * Time Series Aggregations
 */

// GetHourlyAggregation retrieves hourly aggregations for the last N hours.
// Returns a time series with avg, min, max, and count per hour.
//
// Example: GetHourlyAggregation(ctx, "sensor_123", "temperature", 24)
// Returns 24 data points, one per hour.
func (h *IoTHelpers) GetHourlyAggregation(ctx context.Context, assetId, field string, hours int) ([]*TimeSeriesPoint, error) {
	query := fmt.Sprintf(`
		SELECT
			toStartOfHour(created) as hour,
			avg(%s) as avg_value,
			min(%s) as min_value,
			max(%s) as max_value,
			count(*) as count
		FROM events
		WHERE asset_id = ?
		  AND created >= now() - INTERVAL ? HOUR
		  AND arrayExists(x -> x.1 = ?, numberFields)
		GROUP BY hour
		ORDER BY hour ASC
	`,
		ExtractNumberField(field),
		ExtractNumberField(field),
		ExtractNumberField(field),
	)

	rows, err := h.repo.conn.Query(ctx, query, assetId, hours, field)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []*TimeSeriesPoint
	for rows.Next() {
		var point TimeSeriesPoint
		if err := rows.Scan(&point.Timestamp, &point.Avg, &point.Min, &point.Max, &point.Count); err != nil {
			return nil, err
		}
		points = append(points, &point)
	}

	return points, rows.Err()
}

// GetDailyAggregation retrieves daily aggregations for the last N days.
//
// Example: GetDailyAggregation(ctx, "sensor_123", "temperature", 30)
// Returns 30 data points, one per day.
func (h *IoTHelpers) GetDailyAggregation(ctx context.Context, assetId, field string, days int) ([]*TimeSeriesPoint, error) {
	query := fmt.Sprintf(`
		SELECT
			toStartOfDay(created) as day,
			avg(%s) as avg_value,
			min(%s) as min_value,
			max(%s) as max_value,
			count(*) as count
		FROM events
		WHERE asset_id = ?
		  AND created >= now() - INTERVAL ? DAY
		  AND arrayExists(x -> x.1 = ?, numberFields)
		GROUP BY day
		ORDER BY day ASC
	`,
		ExtractNumberField(field),
		ExtractNumberField(field),
		ExtractNumberField(field),
	)

	rows, err := h.repo.conn.Query(ctx, query, assetId, days, field)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []*TimeSeriesPoint
	for rows.Next() {
		var point TimeSeriesPoint
		if err := rows.Scan(&point.Timestamp, &point.Avg, &point.Min, &point.Max, &point.Count); err != nil {
			return nil, err
		}
		points = append(points, &point)
	}

	return points, rows.Err()
}

// GetMinuteAggregation retrieves minute-by-minute aggregations for the last N minutes.
// Useful for real-time monitoring and dashboards.
//
// Example: GetMinuteAggregation(ctx, "sensor_123", "temperature", 60)
// Returns 60 data points, one per minute.
func (h *IoTHelpers) GetMinuteAggregation(ctx context.Context, assetId, field string, minutes int) ([]*TimeSeriesPoint, error) {
	query := fmt.Sprintf(`
		SELECT
			toStartOfMinute(created) as minute,
			avg(%s) as avg_value,
			min(%s) as min_value,
			max(%s) as max_value,
			count(*) as count
		FROM events
		WHERE asset_id = ?
		  AND created >= now() - INTERVAL ? MINUTE
		  AND arrayExists(x -> x.1 = ?, numberFields)
		GROUP BY minute
		ORDER BY minute ASC
	`,
		ExtractNumberField(field),
		ExtractNumberField(field),
		ExtractNumberField(field),
	)

	rows, err := h.repo.conn.Query(ctx, query, assetId, minutes, field)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []*TimeSeriesPoint
	for rows.Next() {
		var point TimeSeriesPoint
		if err := rows.Scan(&point.Timestamp, &point.Avg, &point.Min, &point.Max, &point.Count); err != nil {
			return nil, err
		}
		points = append(points, &point)
	}

	return points, rows.Err()
}

/**
 * Current State Queries
 */

// GetCurrentState retrieves the most recent event for an asset with all EVA fields.
// Useful for dashboards showing current sensor readings.
//
// Note: This returns raw EVA fields with fieldId as keys. To get human-readable
// field names, the caller should resolve fieldIds using the AssetTemplate.
//
// Example: GetCurrentState(ctx, "sensor_123")
func (h *IoTHelpers) GetCurrentState(ctx context.Context, assetId string) (*CurrentState, error) {
	query := `
		SELECT
			asset_id,
			created,
			eva_number,
			eva_string,
			eva_bool,
			eva_date
		FROM events
		WHERE asset_id = ?
		ORDER BY created DESC
		LIMIT 1
	`

	row := h.repo.conn.QueryRow(ctx, query, assetId)

	var event entities.Event
	err := row.Scan(
		&event.AssetId,
		&event.Created,
		&event.EvaNumber,
		&event.EvaString,
		&event.EvaBool,
		&event.EvaDate,
	)
	if err != nil {
		return nil, err
	}

	// Convert EVA maps to generic Fields map (keyed by fieldId as string)
	state := &CurrentState{
		AssetId:   event.AssetId,
		Timestamp: event.Created,
		Fields:    make(map[string]interface{}),
	}

	// Add number fields (fieldId -> value)
	for fieldId, value := range event.EvaNumber {
		state.Fields[fmt.Sprintf("eva_number_%d", fieldId)] = value
	}

	// Add string fields (fieldId -> value)
	for fieldId, value := range event.EvaString {
		state.Fields[fmt.Sprintf("eva_string_%d", fieldId)] = value
	}

	// Add bool fields (fieldId -> value, converted to bool)
	for fieldId, value := range event.EvaBool {
		state.Fields[fmt.Sprintf("eva_bool_%d", fieldId)] = value == 1
	}

	// Add date fields (fieldId -> value)
	for fieldId, value := range event.EvaDate {
		state.Fields[fmt.Sprintf("eva_date_%d", fieldId)] = value
	}

	return state, nil
}

/**
 * Threshold Monitoring
 */

// GetValuesAboveThreshold finds all readings above a threshold in the last N hours.
// Useful for alarm detection and threshold monitoring.
//
// Example: GetValuesAboveThreshold(ctx, "sensor_123", "temperature", 30.0, 24)
// Returns all temperature readings > 30°C in the last 24 hours.
func (h *IoTHelpers) GetValuesAboveThreshold(ctx context.Context, assetId, field string, threshold float64, hours int) ([]*ThresholdViolation, error) {
	query := fmt.Sprintf(`
		SELECT
			asset_id,
			%s as value,
			created
		FROM events
		WHERE asset_id = ?
		  AND created >= now() - INTERVAL ? HOUR
		  AND arrayExists(x -> x.1 = ? AND x.2 > ?, numberFields)
		ORDER BY created DESC
	`, ExtractNumberField(field))

	rows, err := h.repo.conn.Query(ctx, query, assetId, hours, field, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var violations []*ThresholdViolation
	for rows.Next() {
		var v ThresholdViolation
		v.Threshold = threshold
		if err := rows.Scan(&v.AssetId, &v.Value, &v.Timestamp); err != nil {
			return nil, err
		}
		violations = append(violations, &v)
	}

	return violations, rows.Err()
}

// GetValuesBelowThreshold finds all readings below a threshold in the last N hours.
func (h *IoTHelpers) GetValuesBelowThreshold(ctx context.Context, assetId, field string, threshold float64, hours int) ([]*ThresholdViolation, error) {
	query := fmt.Sprintf(`
		SELECT
			asset_id,
			%s as value,
			created
		FROM events
		WHERE asset_id = ?
		  AND created >= now() - INTERVAL ? HOUR
		  AND arrayExists(x -> x.1 = ? AND x.2 < ?, numberFields)
		ORDER BY created DESC
	`, ExtractNumberField(field))

	rows, err := h.repo.conn.Query(ctx, query, assetId, hours, field, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var violations []*ThresholdViolation
	for rows.Next() {
		var v ThresholdViolation
		v.Threshold = threshold
		if err := rows.Scan(&v.AssetId, &v.Value, &v.Timestamp); err != nil {
			return nil, err
		}
		violations = append(violations, &v)
	}

	return violations, rows.Err()
}

/**
 * Multi-Asset Queries
 */

// GetStatsForMultipleAssets retrieves statistics for multiple assets in one query.
// Useful for fleet monitoring and multi-asset dashboards.
//
// Example: GetStatsForMultipleAssets(ctx, []string{"sensor_1", "sensor_2"}, "temperature", 24)
func (h *IoTHelpers) GetStatsForMultipleAssets(ctx context.Context, assetIds []string, field string, hours int) (map[string]*AssetStats, error) {
	// Build IN clause for asset IDs
	placeholders := make([]interface{}, len(assetIds)+2)
	for i, id := range assetIds {
		placeholders[i] = id
	}
	placeholders[len(assetIds)] = hours
	placeholders[len(assetIds)+1] = field

	query := fmt.Sprintf(`
		SELECT
			asset_id,
			argMax(%s, created) as latest_value,
			max(created) as latest_time,
			avg(%s) as avg_value,
			min(%s) as min_value,
			max(%s) as max_value,
			count(*) as count
		FROM events
		WHERE asset_id IN (?)
		  AND created >= now() - INTERVAL ? HOUR
		  AND arrayExists(x -> x.1 = ?, numberFields)
		GROUP BY asset_id
	`,
		ExtractNumberField(field),
		ExtractNumberField(field),
		ExtractNumberField(field),
		ExtractNumberField(field),
	)

	rows, err := h.repo.conn.Query(ctx, query, placeholders...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]*AssetStats)
	for rows.Next() {
		var s AssetStats
		if err := rows.Scan(&s.AssetId, &s.LatestValue, &s.LatestTime, &s.AvgValue, &s.MinValue, &s.MaxValue, &s.ReadingCount); err != nil {
			return nil, err
		}
		stats[s.AssetId] = &s
	}

	return stats, rows.Err()
}

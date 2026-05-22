package clickhouseRepo

import (
	"context"
	"fmt"
	"time"
)

// Example usage of EAV helpers for building ClickHouse queries.
// This file contains practical examples and can be used as reference.

// ExampleSimpleFilter demonstrates basic filtering by EAV fields.
func ExampleSimpleFilter(repo *EventRepositoryClickHouse, ctx context.Context) error {
	// Build query: Find events where temperature > 25
	builder := NewEAVQueryBuilder()
	builder.NumberGreaterThan("temperature", 25)

	whereClause, params := builder.BuildWithPrefix()

	query := fmt.Sprintf(`
		SELECT
			created,
			asset_id,
			%s as temperature
		FROM events
		%s
		ORDER BY created DESC
		LIMIT 100
	`, ExtractNumberField("temperature"), whereClause)

	// Execute query with params
	rows, err := repo.conn.Query(ctx, query, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Process results...
	return nil
}

// ExampleComplexFilter demonstrates combining multiple EAV conditions.
func ExampleComplexFilter(repo *EventRepositoryClickHouse, ctx context.Context, orgId, assetId string) error {
	// Build query: Find events where:
	// - temperature between 20 and 30
	// - location is "warehouse_A"
	// - isOnline is true
	// - last 24 hours
	builder := NewEAVQueryBuilder()
	builder.NumberBetween("temperature", 20, 30).
		StringEquals("location", "warehouse_A").
		BoolEquals("isOnline", true)

	whereClause, params := builder.Build()

	// Add standard filters (org_id, asset_id, created)
	allParams := []interface{}{orgId, assetId}
	allParams = append(allParams, params...)

	query := fmt.Sprintf(`
		SELECT
			created,
			asset_id,
			org_id,
			%s as temperature,
			%s as location,
			%s as is_online
		FROM events
		WHERE org_id = ?
		  AND asset_id = ?
		  AND created >= now() - INTERVAL 1 DAY
		  AND %s
		ORDER BY created DESC
		LIMIT 1000
	`,
		ExtractNumberField("temperature"),
		ExtractStringField("location"),
		ExtractBoolField("isOnline"),
		whereClause,
	)

	rows, err := repo.conn.Query(ctx, query, allParams...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Process results...
	return nil
}

// ExampleAggregation demonstrates aggregating EAV fields.
func ExampleAggregation(repo *EventRepositoryClickHouse, ctx context.Context, orgId string) error {
	// Build aggregation: Average, min, max temperature grouped by location
	// for the last 7 days
	aggBuilder := NewAggregationBuilder("temperature", "number")
	aggBuilder.GroupByStringField("location").
		Avg("avg_temp").
		Min("min_temp").
		Max("max_temp").
		Count("event_count")

	selectClause := aggBuilder.BuildSelect()
	groupByClause := aggBuilder.BuildGroupBy()

	// Also filter: only events with temperature field
	filterBuilder := NewEAVQueryBuilder()
	filterBuilder.HasNumberField("temperature")
	whereClause, params := filterBuilder.Build()

	allParams := []interface{}{orgId}
	allParams = append(allParams, params...)

	query := fmt.Sprintf(`
		SELECT %s
		FROM events
		WHERE org_id = ?
		  AND created >= now() - INTERVAL 7 DAY
		  AND %s
		%s
	`, selectClause, whereClause, groupByClause)

	rows, err := repo.conn.Query(ctx, query, allParams...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Process aggregated results...
	// Each row: location, avg_temp, min_temp, max_temp, event_count
	return nil
}

// ExampleStringSearch demonstrates searching in string fields.
func ExampleStringSearch(repo *EventRepositoryClickHouse, ctx context.Context, orgId string) error {
	// Find events where:
	// - deviceName contains "sensor"
	// - status is either "online" or "active"
	builder := NewEAVQueryBuilder()
	builder.StringContains("deviceName", "sensor").
		StringIn("status", []string{"online", "active"})

	whereClause, params := builder.Build()
	allParams := []interface{}{orgId}
	allParams = append(allParams, params...)

	query := fmt.Sprintf(`
		SELECT
			created,
			asset_id,
			%s as device_name,
			%s as status
		FROM events
		WHERE org_id = ?
		  AND %s
		ORDER BY created DESC
		LIMIT 100
	`, ExtractStringField("deviceName"), ExtractStringField("status"), whereClause)

	rows, err := repo.conn.Query(ctx, query, allParams...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

// ExampleDateRangeQuery demonstrates querying by date fields.
func ExampleDateRangeQuery(repo *EventRepositoryClickHouse, ctx context.Context, orgId string) error {
	// Find events where lastSeen was in the last hour
	oneHourAgo := time.Now().Add(-1 * time.Hour)

	builder := NewEAVQueryBuilder()
	builder.DateAfter("lastSeen", oneHourAgo)

	whereClause, params := builder.Build()
	allParams := []interface{}{orgId}
	allParams = append(allParams, params...)

	query := fmt.Sprintf(`
		SELECT
			created,
			asset_id,
			%s as last_seen
		FROM events
		WHERE org_id = ?
		  AND %s
		ORDER BY created DESC
	`, ExtractDateField("lastSeen"), whereClause)

	rows, err := repo.conn.Query(ctx, query, allParams...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

// ExampleExistenceCheck demonstrates checking for field existence.
func ExampleExistenceCheck(repo *EventRepositoryClickHouse, ctx context.Context, orgId string) error {
	// Find all events that have a "temperature" field
	// (regardless of value)
	builder := NewEAVQueryBuilder()
	builder.HasNumberField("temperature").
		HasStringField("location")

	whereClause, params := builder.Build()
	allParams := []interface{}{orgId}
	allParams = append(allParams, params...)

	query := fmt.Sprintf(`
		SELECT
			created,
			asset_id,
			%s as temperature,
			%s as location
		FROM events
		WHERE org_id = ?
		  AND created >= now() - INTERVAL 1 DAY
		  AND %s
		ORDER BY created DESC
		LIMIT 1000
	`,
		ExtractNumberField("temperature"),
		ExtractStringField("location"),
		whereClause,
	)

	rows, err := repo.conn.Query(ctx, query, allParams...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

// ExampleLastValue demonstrates getting the most recent value for a field.
func ExampleLastValue(repo *EventRepositoryClickHouse, ctx context.Context, assetId string) error {
	// Get the latest temperature reading for an asset
	query := fmt.Sprintf(`
		SELECT
			created,
			%s as temperature
		FROM events
		WHERE asset_id = ?
		  AND arrayExists(x -> x.1 = 'temperature', numberFields)
		ORDER BY created DESC
		LIMIT 1
	`, ExtractNumberField("temperature"))

	row := repo.conn.QueryRow(ctx, query, assetId)

	var created time.Time
	var temperature float64
	err := row.Scan(&created, &temperature)
	if err != nil {
		return err
	}

	fmt.Printf("Latest temperature: %.2f at %s\n", temperature, created)
	return nil
}

// ExampleTimeSeriesAggregation demonstrates time-series aggregation.
func ExampleTimeSeriesAggregation(repo *EventRepositoryClickHouse, ctx context.Context, assetId string) error {
	// Get hourly average temperature for the last 24 hours
	query := fmt.Sprintf(`
		SELECT
			toStartOfHour(created) as hour,
			avg(%s) as avg_temp,
			count(*) as readings
		FROM events
		WHERE asset_id = ?
		  AND created >= now() - INTERVAL 24 HOUR
		  AND arrayExists(x -> x.1 = 'temperature', numberFields)
		GROUP BY hour
		ORDER BY hour ASC
	`, ExtractNumberField("temperature"))

	rows, err := repo.conn.Query(ctx, query, assetId)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Process time-series data...
	for rows.Next() {
		var hour time.Time
		var avgTemp float64
		var readings uint64

		if err := rows.Scan(&hour, &avgTemp, &readings); err != nil {
			return err
		}

		fmt.Printf("Hour: %s, Avg Temp: %.2f, Readings: %d\n", hour, avgTemp, readings)
	}

	return nil
}

// ExampleMultiFieldAggregation demonstrates aggregating multiple fields.
func ExampleMultiFieldAggregation(repo *EventRepositoryClickHouse, ctx context.Context, orgId string) error {
	// Get statistics for multiple metrics grouped by location
	query := fmt.Sprintf(`
		SELECT
			%s as location,
			avg(%s) as avg_temp,
			avg(%s) as avg_humidity,
			count(*) as event_count
		FROM events
		WHERE org_id = ?
		  AND created >= now() - INTERVAL 1 DAY
		  AND arrayExists(x -> x.1 = 'temperature', numberFields)
		  AND arrayExists(x -> x.1 = 'humidity', numberFields)
		  AND arrayExists(x -> x.1 = 'location', stringFields)
		GROUP BY location
		ORDER BY event_count DESC
	`,
		ExtractStringField("location"),
		ExtractNumberField("temperature"),
		ExtractNumberField("humidity"),
	)

	rows, err := repo.conn.Query(ctx, query, orgId)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Process aggregated results...
	return nil
}

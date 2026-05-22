package clickhouseRepo

import (
	"time"

	"events/src/modules/events/domain/entities"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
)

// EventRepositoryClickHouse implements the EventRepository interface using ClickHouse.
// This follows the Repository pattern from Hexagonal Architecture.
//
// It uses chModel.Table for generic operations on event tables,
// providing automatic JSON marshaling/unmarshaling for map fields.
//
// Repository methods are split across multiple files:
//   - event_repository.go: Struct definition and constructor
//   - event_repository_legacy.go: Legacy Event methods (Save, SaveBatch)
//   - raw_event_repository.go: RawEvent methods (SaveRawEvent, SaveRawEventBatch, Query)
//   - jsexec_event_repository.go: JsExecEvent methods (SaveJsExecEventBatch, Query)
//   - dlq_event_repository.go: DLQEvent methods (SaveDLQEvent, SaveDLQEventBatch, Query)
type EventRepositoryClickHouse struct {
	conn                   driver.Conn
	eventTable             *chModel.Table[entities.Event]
	rawEventTable          *chModel.Table[entities.RawEvent]
	jsExecEventTable       *chModel.Table[entities.JsExecEvent]
	dlqEventTable          *chModel.Table[entities.DLQEvent]
	routerEventTable       *chModel.Table[entities.RouterEvent]
	businessRuleEventTable *chModel.Table[entities.BusinessRuleEvent]
	triggerEventTable      *chModel.Table[entities.TriggerEvent]
	workflowEventTable     *chModel.Table[entities.WorkflowEvent]
}

/**
 * EAVQueryBuilder provides helper functions for building queries with EAV buckets.
 * This simplifies complex ClickHouse queries when searching or aggregating EAV fields.
 */
type EAVQueryBuilder struct {
	conditions []string
	params     []interface{}
}

// AggregationBuilder helps build aggregation queries for EAV fields.
type AggregationBuilder struct {
	field      string
	fieldType  string // "number", "string", "bool", "date"
	groupBy    []string
	aggregates []string
}

/**
 * IoTHelpers provides specialized query helpers for IoT use cases.
 * These methods handle common IoT patterns like:
 * - Last value queries
 * - Time-window aggregations (avg, min, max)
 * - Time-series data
 * - Current state retrieval
 */
type IoTHelpers struct {
	repo *EventRepositoryClickHouse
}

// LastValueResult represents a single field's last value.
type LastValueResult struct {
	Value     interface{}
	Timestamp time.Time
	AssetId   string
}

// TimeWindowStats represents statistical data for a time window.
type TimeWindowStats struct {
	Field     string
	Avg       float64
	Min       float64
	Max       float64
	Count     uint64
	StartTime time.Time
	EndTime   time.Time
}

// TimeSeriesPoint represents a single point in a time series.
type TimeSeriesPoint struct {
	Timestamp time.Time
	Avg       float64
	Min       float64
	Max       float64
	Count     uint64
}

// CurrentState represents the current state of an asset with multiple fields.
type CurrentState struct {
	AssetId   string
	Timestamp time.Time
	Fields    map[string]interface{}
}

// ThresholdViolation represents a reading that violated a threshold.
type ThresholdViolation struct {
	AssetId   string
	Value     float64
	Threshold float64
	Timestamp time.Time
}

// AssetStats represents statistics for a single asset.
type AssetStats struct {
	AssetId      string
	LatestValue  float64
	LatestTime   time.Time
	AvgValue     float64
	MinValue     float64
	MaxValue     float64
	ReadingCount uint64
}

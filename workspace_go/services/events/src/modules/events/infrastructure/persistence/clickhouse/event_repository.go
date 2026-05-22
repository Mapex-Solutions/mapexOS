package clickhouseRepo

import (
	"events/src/modules/events/domain/entities"
	"events/src/modules/events/domain/repositories"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewEventRepository creates a new ClickHouse-based event repository.
// This is registered in the DIG container and injected where needed.
func NewEventRepository(conn driver.Conn) repositories.EventRepository {
	// Initialize the Table wrapper for events (processed events with EVA)
	eventTable, err := chModel.NewTable[entities.Event](conn, "events", chModel.TableConfig{
		TimestampField: "created",
		DefaultOrder:   "DESC",
	})
	if err != nil {
		logger.Error(err, "[REPO:Event] Failed to initialize events table model")
	}

	// Initialize the Table wrapper for events_raw
	rawTable, err := chModel.NewTable[entities.RawEvent](conn, "events_raw", chModel.TableConfig{
		TimestampField: "created",
		DefaultOrder:   "DESC",
	})
	if err != nil {
		logger.Error(err, "[REPO:Event] Failed to initialize events_raw table model")
	}

	// Initialize the Table wrapper for events_jsexecutor
	jsExecTable, err := chModel.NewTable[entities.JsExecEvent](conn, "events_jsexecutor", chModel.TableConfig{
		TimestampField: "created",
		DefaultOrder:   "DESC",
	})
	if err != nil {
		logger.Error(err, "[REPO:Event] Failed to initialize events_jsexecutor table model")
	}

	// Initialize the Table wrapper for events_dlq
	dlqTable, err := chModel.NewTable[entities.DLQEvent](conn, "events_dlq", chModel.TableConfig{
		TimestampField: "created",
		DefaultOrder:   "DESC",
	})
	if err != nil {
		logger.Error(err, "[REPO:Event] Failed to initialize events_dlq table model")
	}

	// Initialize the Table wrapper for events_router
	routerTable, err := chModel.NewTable[entities.RouterEvent](conn, "events_router", chModel.TableConfig{
		TimestampField: "created",
		DefaultOrder:   "DESC",
	})
	if err != nil {
		logger.Error(err, "[REPO:Event] Failed to initialize events_router table model")
	}

	// Initialize the Table wrapper for events_businessrule
	businessRuleTable, err := chModel.NewTable[entities.BusinessRuleEvent](conn, "events_businessrule", chModel.TableConfig{
		TimestampField: "created",
		DefaultOrder:   "DESC",
	})
	if err != nil {
		logger.Error(err, "[REPO:Event] Failed to initialize events_businessrule table model")
	}

	// Initialize the Table wrapper for events_trigger
	triggerTable, err := chModel.NewTable[entities.TriggerEvent](conn, "events_trigger", chModel.TableConfig{
		TimestampField: "created",
		DefaultOrder:   "DESC",
	})
	if err != nil {
		logger.Error(err, "[REPO:Event] Failed to initialize events_trigger table model")
	}

	// Initialize the Table wrapper for events_workflow
	workflowTable, err := chModel.NewTable[entities.WorkflowEvent](conn, "events_workflow", chModel.TableConfig{
		TimestampField: "created",
		DefaultOrder:   "DESC",
	})
	if err != nil {
		logger.Error(err, "[REPO:Event] Failed to initialize events_workflow table model")
	}

	return &EventRepositoryClickHouse{
		conn:                   conn,
		eventTable:             eventTable,
		rawEventTable:          rawTable,
		jsExecEventTable:       jsExecTable,
		dlqEventTable:          dlqTable,
		routerEventTable:       routerTable,
		businessRuleEventTable: businessRuleTable,
		triggerEventTable:      triggerTable,
		workflowEventTable:     workflowTable,
	}
}

// Compile-time check to ensure EventRepositoryClickHouse implements EventRepository interface
var _ repositories.EventRepository = (*EventRepositoryClickHouse)(nil)

package repositories

import (
	"context"
	"time"

	"events/src/modules/events/domain/entities"

	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
)

// BusinessRuleEventRepository defines storage operations for Business Rule events.
type BusinessRuleEventRepository interface {
	SaveBusinessRuleEventBatch(ctx context.Context, events []*entities.BusinessRuleEvent) error

	QueryEventsBusinessRuleCursor(
		ctx context.Context,
		orgFilter chModel.Map,
		eventTrackerId *string,
		threadId *string,
		ruleId *string,
		businessRuleId *string,
		matched *bool,
		startTime *time.Time,
		endTime *time.Time,
		cursorOpts *chModel.TimeCursorOpts,
	) (*chModel.TimeCursorResult[entities.BusinessRuleEvent], error)
}

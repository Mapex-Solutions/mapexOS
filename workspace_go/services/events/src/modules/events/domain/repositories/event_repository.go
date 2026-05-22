package repositories

import (
	"context"

	"events/src/modules/events/domain/entities"
)

// EventRepository composes all event sub-repositories via embedding (ISP).
// This is backwards-compatible: any implementation that satisfies EventRepository
// also satisfies each individual sub-interface.
//
// In FASE 5, each Application Service will depend on the specific sub-interface
// it needs, not the full composite. For now, the DI container still injects this.
type EventRepository interface {
	// Legacy methods (Save single event)
	Save(ctx context.Context, event *entities.Event) error
	SaveBatch(ctx context.Context, events []*entities.Event) error

	// Sub-interfaces by aggregate
	RawEventRepository
	JsExecEventRepository
	DLQEventRepository
	RouterEventRepository
	BusinessRuleEventRepository
	EventStoreRepository
	TriggerEventRepository
	WorkflowEventRepository
}

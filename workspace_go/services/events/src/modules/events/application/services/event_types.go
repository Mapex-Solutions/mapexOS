package services

import (
	"events/src/modules/events/application/di"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// EventService provides methods for processing and storing events.
// It serves as an application service layer that handles domain-level actions
// for events received from NATS and stores them in ClickHouse.
//
// This service implements the EventServicePort interface, following
// Hexagonal Architecture principles by depending on interfaces rather
// than concrete implementations.
type EventService struct {
	deps di.EventServiceDependenciesInjection
}

// tenantContext is the minimal shape used to extract orgId/pathKey/eventTrackerId
// from a raw NATS payload before full DTO unmarshal. Ensures DLQ messages always
// carry tenant context even when the full unmarshal fails downstream.
type tenantContext struct {
	OrgId          string `json:"orgId"`
	PathKey        string `json:"pathKey"`
	EventTrackerId string `json:"eventTrackerId"`
}

// messageResult holds the outcome of processing a single message during Phase 1.
//
// Phase 1 sets:
//   - action="reject" + rejectReason: validation/parse failure -> DLQ
//   - action="pending" + entity: valid entity awaiting Phase 2 bulk insert
//   - action="ack_skip": DLQ consumer parse failure (ACK without entity)
//
// Phase 3 resolves "pending" -> "ack" (insert OK) or "nack" (insert failed).
type messageResult[T any] struct {
	msg          *natsModel.Message
	entity       *T
	action       string // "reject", "pending", "ack_skip"
	rejectReason string
}

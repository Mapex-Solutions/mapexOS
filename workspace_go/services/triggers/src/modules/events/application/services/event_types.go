package services

import (
	"triggers/src/modules/events/application/di"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// EventService handles trigger execution events from NATS.
//
// Following Hexagonal Architecture principles:
// - Application layer service (orchestration logic)
// - Depends on PORTS (interfaces), never on concrete implementations
// - Infrastructure adapters injected via dependency injection
//
// This service implements the EventServicePort interface.
//
// Workflow:
// - Receives TriggerExecuteEvent from NATS consumer (interface layer)
// - Fetches trigger configuration from database/cache (via TriggerService port)
// - Resolves placeholders in trigger config using event payload (application logic)
// - Delegates execution to appropriate executor based on trigger type (via ExecutorRegistry port)
// - Logs result
//
// Note: The service does NOT know about concrete executor implementations.
// It only depends on the ExecutorRegistry interface (port).
type EventService struct {
	deps    di.EventServiceDependenciesInjection
	workers int
}

// messageResult holds the outcome of processing a single trigger message.
// Collected by goroutines during Phase 1, consumed in Phase 3 for ACK/Nack/Reject.
type messageResult struct {
	msg          *natsModel.Message
	action       string // "ack", "nack", "reject"
	nackErr      error  // error for Nack()
	rejectReason string // reason for Reject()
	isDLQ        bool   // true when retries exhausted (track as "dlq" instead of "nack")
}

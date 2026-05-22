package services

import (
	"router/src/modules/events/application/di"
	"router/src/modules/events/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// EventService handles route execution events following Hexagonal Architecture.
type EventService struct {
	deps           di.EventServiceDependenciesInjection
	matchEvaluator ports.MatchEvaluatorPort
}

// messageResult holds the outcome of processing a single route message.
// Collected by goroutines during Phase 1, consumed in Phase 3 for ACK/Nack/Reject.
type messageResult struct {
	msg          *natsModel.Message
	action       string // "ack", "nack", "reject"
	nackErr      error  // error for Nack()
	rejectReason string // reason for Reject()
	status       string // metrics label: "success", "error"
	duration     float64
}

// tierLabels maps TieredCache tier values to Prometheus label values.
var tierLabels = map[int]string{
	0:  "L0_RAM",
	1:  "L1_Disk",
	2:  "L2_MinIO",
	3:  "Fallback_HTTP",
	-1: "MISS",
}

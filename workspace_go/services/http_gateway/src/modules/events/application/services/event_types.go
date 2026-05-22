package services

import (
	"http_gateway/src/modules/events/application/di"
)

// EventService provides methods for managing event-related operations.
//
// It serves as an application service layer that receives events from external
// sources (webhooks, APIs) and publishes them to NATS for asynchronous processing.
//
// This service implements the EventServicePort interface, following
// Hexagonal Architecture principles by depending on interfaces rather
// than concrete implementations.
type EventService struct {
	deps di.EventServiceDependenciesInjection
}

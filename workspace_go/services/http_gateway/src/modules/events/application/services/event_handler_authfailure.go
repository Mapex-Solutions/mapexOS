package services

import (
	"context"
	"time"

	dsDto "http_gateway/src/modules/datasources/application/dtos"
	eventsConstants "http_gateway/src/modules/events/application/constants"

	eventsDto "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// buildAuthFailurePayload composes the RawEventDTO documented as the events.raw
// shape for failed authentication attempts. Field extraction was previously
// inlined in PublishAuthFailure (55-LOC method, /go-arch §3 anti-pattern A);
// extracting here keeps the public method as a thin orchestration skeleton.
func (s *EventService) buildAuthFailurePayload(dataSource *dsDto.DataSourceResponse, event map[string]any, eventTrackerId string, errorMsg string) eventsDto.RawEventDTO {
	now := time.Now().UTC()
	orgId := ""
	if dataSource != nil && dataSource.OrgId != nil {
		orgId = dataSource.OrgId.Hex()
	}
	pathKey := ""
	if dataSource != nil && dataSource.PathKey != nil {
		pathKey = *dataSource.PathKey
	}
	threadId := ""
	if dataSource != nil && dataSource.ID != nil {
		threadId = dataSource.ID.Hex()
	}
	name := ""
	if dataSource != nil && dataSource.Name != nil {
		name = *dataSource.Name
	}
	description := ""
	if dataSource != nil && dataSource.Description != nil {
		description = *dataSource.Description
	}
	return eventsDto.RawEventDTO{
		EventTrackerId: eventTrackerId,
		ThreadId:       threadId,
		OrgId:          orgId,
		PathKey:        pathKey,
		Event:          event,
		Source:         "http_gateway",
		Created:        &now,
		Name:           name,
		Description:    description,
		Success:        false,
		Error:          errorMsg,
	}
}

// publishAuthFailureFireAndForget publishes the RawEventDTO to events.raw in
// a detached goroutine — the auth middleware MUST NOT block on this publish.
// Errors are silently swallowed (the auth response was already sent and the
// security audit trail is best-effort).
func (s *EventService) publishAuthFailureFireAndForget(payload eventsDto.RawEventDTO) {
	go func() {
		_ = s.deps.NatsBus.Publish(natsModel.PublishConfig{
			Ctx:     context.Background(),
			Subject: eventsConstants.EventsRawSubject,
			Data:    payload,
			Headers: nil,
		})
	}()
}

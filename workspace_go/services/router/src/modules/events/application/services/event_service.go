package services

import (
	ctx "context"
	"sync"

	"router/src/modules/events/application/di"
	"router/src/modules/events/application/ports"
	domainServices "router/src/modules/events/domain/services"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// New creates a new EventService instance.
func New(deps di.EventServiceDependenciesInjection) ports.EventServicePort {
	return &EventService{
		deps:           deps,
		matchEvaluator: domainServices.NewMatchEvaluator(),
	}
}

var _ ports.EventServicePort = (*EventService)(nil)

// ProcessEventBatch processes route execution events using a bounded worker pool.
//
// Three-Phase Processing:
//   - Phase 1 (Parallel): Workers process messages concurrently via bounded pool.
//   - Phase 2 (Flush):    Single FlushConnection for all fire-and-forget publishes.
//   - Phase 3 (ACK):      Sequential ACK/Nack/Reject based on collected results.
func (s *EventService) ProcessEventBatch(messages []*natsModel.Message) error {
	if len(messages) == 0 {
		return nil
	}
	s.deps.Metrics.EventsBatchSize.Observe(float64(len(messages)))
	results := s.processBatchPhase1(messages)
	s.processBatchPhase2Flush()
	s.processBatchPhase3Ack(results)
	return nil
}

// ProcessEvent processes a single route execution event (legacy V1).
func (s *EventService) ProcessEvent(data []byte, index int, headers map[string][]string) error {
	orgId, assetId, event, eventTrackerId, eventSource, err := s.parseLegacyEventPayload(data)
	if err != nil {
		return err
	}
	return s.execute(ctx.Background(), orgId, assetId, event, eventTrackerId, eventSource)
}

// ProcessAssetInvalidateBatch fans out one goroutine per FANOUT cache-invalidation
// message and waits for them all to complete. Each message clears its
// {orgId}/{assetUUID} entry in the local L0+L1 asset cache.
func (s *EventService) ProcessAssetInvalidateBatch(messages []*natsModel.Message) {
	var wg sync.WaitGroup
	for i, msg := range messages {
		wg.Add(1)
		go func(idx int, m *natsModel.Message) {
			defer wg.Done()
			s.processAssetInvalidateMessage(idx, m)
		}(i, msg)
	}
	wg.Wait()
}

// ProcessTemplateInvalidateBatch fans out one goroutine per FANOUT template
// cache-invalidation message and waits for them all to complete. Each message
// clears its {orgId}/{templateId} entry in the local L0+L1 template cache.
// L2 source of truth in MinIO is untouched.
func (s *EventService) ProcessTemplateInvalidateBatch(messages []*natsModel.Message) {
	if len(messages) == 0 {
		return
	}
	var wg sync.WaitGroup
	for i, msg := range messages {
		wg.Add(1)
		go func(idx int, m *natsModel.Message) {
			defer wg.Done()
			s.processTemplateInvalidateMessage(idx, m)
		}(i, msg)
	}
	wg.Wait()
}

package services

import (
	"fmt"

	"triggers/src/modules/events/application/constants"
	"triggers/src/modules/events/application/di"
	"triggers/src/modules/events/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// Compile-time check to ensure EventService implements EventServicePort interface.
var _ ports.EventServicePort = (*EventService)(nil)

// New creates and returns a new instance of EventService. The bounded
// worker count is read from config (trigger_executor_workers) with a
// safe default for the parallel batch phase.
func New(deps di.EventServiceDependenciesInjection) ports.EventServicePort {
	workers, _ := config.GetIntValue("trigger_executor_workers")
	if workers <= 0 {
		workers = constants.DefaultWorkerCount
	}
	logger.Info(fmt.Sprintf("[SERVICE:Event] Executor workers: %d", workers))
	return &EventService{
		deps:    deps,
		workers: workers,
	}
}

// ProcessTriggerExecutionBatch processes a batch of trigger execution
// events. Three-phase pipeline: parallel parse + execute (bounded by
// worker count) -> single CorePublisher flush for all fire-and-forget
// publishes -> sequential Ack/Nack/Reject from collected results.
func (s *EventService) ProcessTriggerExecutionBatch(messages []*natsModel.Message) error {
	if len(messages) == 0 {
		return nil
	}
	s.deps.Metrics.TriggersBatchSize.Observe(float64(len(messages)))
	logger.Info(fmt.Sprintf("[SERVICE:Event] Processing trigger execution batch: %d messages, %d workers", len(messages), s.workers))
	results := s.runTriggerBatchPhase1(messages)
	s.flushBatchPublishes()
	s.runTriggerBatchPhase3Ack(results)
	return nil
}

// ProcessTriggerExecution handles a single trigger execution event from
// the legacy non-batch path. Steps: parse the event payload -> fetch and
// validate the trigger -> resolve placeholders against the event -> execute
// via the matching executor.
func (s *EventService) ProcessTriggerExecution(data []byte) error {
	event, err := s.parseTriggerExecuteEvent(data)
	if err != nil {
		return err
	}
	trigger, skip, err := s.fetchTriggerForExecution(event)
	if err != nil || skip {
		return err
	}
	resolvedConfig, triggerType, err := s.resolveTriggerConfig(trigger, event)
	if err != nil {
		return err
	}
	return s.executeResolvedTrigger(triggerType, resolvedConfig, event)
}

// ProcessWorkflowExecutionBatch processes a batch of workflow plugin /
// trigger-as-step execution requests from the workflow service. Same
// three-phase pipeline as ProcessTriggerExecutionBatch, with the resume
// callback published in Phase 1 and flushed in Phase 2.
func (s *EventService) ProcessWorkflowExecutionBatch(messages []*natsModel.Message) error {
	if len(messages) == 0 {
		return nil
	}
	s.deps.Metrics.WorkflowBatchSize.Observe(float64(len(messages)))
	logger.Info(fmt.Sprintf("[SERVICE:Event] Processing workflow execution batch: %d messages", len(messages)))
	results := s.runWorkflowBatchPhase1(messages)
	s.flushBatchPublishes()
	s.runWorkflowBatchPhase3Ack(results)
	return nil
}

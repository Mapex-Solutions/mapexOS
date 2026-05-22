package ports

import (
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// EventServicePort defines the contract for trigger execution event processing operations.
//
// This port interface enables Hexagonal Architecture by decoupling the business
// logic from its implementation. NATS consumers depend on this interface
// rather than the concrete service implementation.
//
// Methods:
//   - ProcessTriggerExecution: Processes a trigger execution event received from NATS (V1 - legacy)
//   - ProcessTriggerExecutionBatch: Processes a batch of events with retry/DLQ support (V2 - recommended)
type EventServicePort interface {
	// ProcessTriggerExecution handles the processing of trigger execution events.
	// LEGACY (V1): This method is called by the old NATS consumer for each message received on trigger.*.execute.
	//
	// Workflow:
	// Deserialize TriggerExecuteEvent from JSON
	// Fetch trigger configuration from database/cache using triggerId
	// Resolve placeholders in trigger config using event payload
	// Execute the trigger based on its type (http, mqtt, email, etc.)
	// Log execution result
	//
	// Parameters:
	//   - data: The message payload containing trigger execution event (JSON)
	//
	// Returns:
	//   - nil: Event processed successfully (will be ACKed by NATS)
	//   - error: Processing failed (will be NAKed and redelivered)
	ProcessTriggerExecution(data []byte) error

	// ProcessTriggerExecutionBatch processes a complete batch of trigger execution events at once.
	// RECOMMENDED (V2): Receives all messages from NATS batch with retry/DLQ support.
	//
	// The service handles ALL message lifecycle decisions:
	//   - msg.Reject(reason): Invalid JSON/validation errors → DLQ immediately
	//   - msg.Nack(err): Processing failed → retry with backoff
	//   - msg.Ack(): Successfully processed → removed from queue
	//
	// Parameters:
	//   - messages: Slice of Message pointers with retry-aware Ack/Nack/Reject methods
	//
	// Returns:
	//   - Always nil (service handles all messages internally)
	ProcessTriggerExecutionBatch(messages []*natsModel.Message) error

	// ProcessWorkflowExecutionBatch processes a batch of execution requests from the Workflow Service.
	// Each message contains a WorkflowTriggerRequest with mode "trigger" or "plugin".
	// mode "trigger": fetch trigger entity config, resolve, execute, publish callback resume.
	// mode "plugin": execute resolved action pipeline (hooks + operation), publish callback resume.
	ProcessWorkflowExecutionBatch(messages []*natsModel.Message) error
}

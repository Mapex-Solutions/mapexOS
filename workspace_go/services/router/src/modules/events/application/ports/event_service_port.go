package ports

import (
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// EventServicePort defines the contract for event processing operations.
//
// This port interface enables Hexagonal Architecture by decoupling the business
// logic from its implementation. NATS consumers depend on this interface
// rather than the concrete service implementation.
//
// Methods:
//   - ProcessEvent: Processes a route execution event received from NATS (V1 - legacy)
//   - ProcessEventBatch: Processes a batch of events with retry/DLQ support (V2 - recommended)
//   - ProcessAssetInvalidateBatch: Invalidates local asset cache via FANOUT messages
//   - ProcessTemplateInvalidateBatch: Invalidates local template cache via FANOUT messages
type EventServicePort interface {
	// ProcessEvent handles the processing of route execution events.
	// LEGACY (V1): This method is called by the old NATS consumer for each message received.
	//
	// Parameters:
	//   - data: The message payload containing route execution data
	//   - index: Position in the batch (for logging/debugging)
	//   - headers: Message headers containing metadata (Nats-Msg-Id, timestamp, source, etc.)
	//
	// Returns:
	//   - nil: Event processed successfully (will be ACKed by NATS)
	//   - error: Processing failed (will be NAKed and redelivered)
	ProcessEvent(data []byte, index int, headers map[string][]string) error

	// ProcessEventBatch processes a complete batch of route execution events at once.
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
	ProcessEventBatch(messages []*natsModel.Message) error

	// ProcessAssetInvalidateBatch handles FANOUT cache invalidation messages.
	// Invalidates local cache (L0+L1) when asset data changes.
	//
	// TieredCache Architecture:
	//   L0 (RAM): Hot cache with ~5min TTL
	//   L1 (Disk): Persistent cache with ~1h TTL
	//   L2 (MinIO): Source of truth (AssetReadModel JSON)
	//
	// FANOUT Pattern:
	//   - Each service instance receives the invalidation message
	//   - Only L0+L1 are cleared (L2 is source of truth, already updated)
	//   - Next request fetches fresh data from L2 → populates L0/L1
	//
	// Parameters:
	//   - messages: Slice of Message pointers with Ack/Nack/Reject methods
	ProcessAssetInvalidateBatch(messages []*natsModel.Message)

	// ProcessTemplateInvalidateBatch handles FANOUT template cache invalidation messages.
	// Invalidates local template cache (L0+L1) when a template is created, updated, or deleted.
	//
	// TieredCache Architecture:
	//   L0 (RAM): Hot cache cleared on invalidation
	//   L1 (Disk): Persistent cache cleared on invalidation
	//   L2 (MinIO): Source of truth (AssetTemplate JSON)
	//
	// FANOUT Pattern:
	//   - Each service instance receives the invalidation message
	//   - Only L0+L1 are cleared (L2 is source of truth, already updated)
	//   - Next request fetches fresh data from L2 → populates L0/L1
	//
	// Parameters:
	//   - messages: Slice of Message pointers carrying the invalidateEvent payload
	ProcessTemplateInvalidateBatch(messages []*natsModel.Message)
}

package types

import natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"

/**
 * INTERNAL SERVICE TYPES
 * Helper types used within the ArchiverService implementation.
 */

/**
 * MsgRef tracks which NATS message belongs to which batch for ACK/NACK routing.
 */
type MsgRef struct {
	// Msg is the original NATS message.
	Msg *natsModel.Message

	// Batch identifies the batch type ("created" or "terminal").
	Batch string
}

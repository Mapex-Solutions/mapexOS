package message

import sharedTypes "workflow/src/shared/types"

/*
 * NATS MESSAGE TYPES
 * Serialized as JSON when published/consumed via NATS JetStream streams.
 *
 * Shared types live in src/shared/types/. This file re-exports them
 * so existing archiver code continues to compile with archiverMsg.StateEvent.
 */

// StateEvent is an alias to the shared type. Published by Runtime, consumed by Archiver.
type StateEvent = sharedTypes.StateEvent

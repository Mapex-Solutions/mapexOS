package message

import (
	definitionsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/definitions"
)

/*
 * NATS MESSAGE TYPES
 * Payload types for messages published by the definitions module.
 *
 * Cross-service payloads (consumed by js-workflow-executor) live in
 * packages/contracts/services/workflow/definitions/. This file re-exposes
 * them to the local module via aliases.
 */

// DefinitionInvalidatePayload is the FANOUT message payload for cache invalidation.
// Published when code nodes are created, updated, or deleted in a definition.
// Consumers (js-workflow-executor) invalidate L0 (RAM) + L1 (Disk) for the specified nodeIds.
type DefinitionInvalidatePayload = definitionsContract.DefinitionInvalidatePayload

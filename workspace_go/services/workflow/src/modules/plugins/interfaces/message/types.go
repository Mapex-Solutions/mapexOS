package message

import (
	"workflow/src/modules/plugins/application/ports"
)

/*
 * NATS MESSAGE TYPES
 * Payload types for messages published by the plugins module.
 *
 * The canonical payload type is defined in application/ports/types.go so that
 * application services do not need to import interfaces/message (see Hexagonal layering).
 * This file re-exposes the type for the interface layer (module.go subscriber).
 */

// PluginInvalidatePayload represents the FANOUT message payload for cache invalidation.
// Published when a plugin manifest is created, updated, or deleted.
// Consumers (all workflow pods) invalidate L0 (RAM) + L1 (Disk) for the specified pluginId.
type PluginInvalidatePayload = ports.PluginInvalidatePayload

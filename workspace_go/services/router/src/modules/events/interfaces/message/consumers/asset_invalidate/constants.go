package asset_invalidate

import (
	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for AssetInvalidate FANOUT consumer
 *
 * FANOUT Pattern:
 * - Each service instance receives a copy of the message
 * - Used for cache invalidation across all replicas
 * - No queue group (each instance processes independently)
 */

// Stream name for asset invalidation events (shared FANOUT stream)
var Stream = contracts.FanoutStreamName

// Subject for asset invalidation events
var Subject = contracts.FanoutAssetSubject

// Durable name for this consumer.
var Durable = config.Durable("router", "asset-invalidate")

// EventType for DLQ metadata
const EventType = "fanout.asset.invalidate"

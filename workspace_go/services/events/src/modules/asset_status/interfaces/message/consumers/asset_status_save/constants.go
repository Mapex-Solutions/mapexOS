package asset_status_save

import (
	hmContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/healthmonitor"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Constants for AssetStatusSave consumer.
 *
 * The wire-level contract (subject/stream) is owned by the assets service
 * healthmonitor module and declared in packages/contracts/services/assets/
 * healthmonitor. These locals are thin aliases kept so the consumer file
 * stays stable when subjects evolve. Durable is local to this consumer.
 */

// Stream name for asset connectivity persistence events (offline/online transitions).
var Stream = hmContract.AssetStatusSaveStream

// Subject the healthmonitor alert publisher emits to on every transition.
var Subject = hmContract.AssetStatusSaveSubject

// Durable name for the asset_status_save consumer.
var Durable = config.Durable("events", "asset-status")

// EventType tags DLQ messages produced by this consumer.
const EventType = hmContract.AssetStatusSaveEventType

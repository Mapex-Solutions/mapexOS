package asset_l2sync

import (
	authContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/auth"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Stream is the durable retry stream for L2 (MinIO) write failures.
// Shared across asset and assettemplate modules — consumers
// distinguish via Subject filter.
var Stream = authContract.L2WritesStreamName

// Subject filters retries specific to the assets module.
var Subject = authContract.L2WritesAssetSubject

// Durable name for this consumer.
var Durable = config.Durable("assets", "l2sync")

// QueueGroup distributes retries across multiple assets-MS replicas.
const QueueGroup = "assets-l2sync"

// EventType identifies this consumer in DLQ metadata.
const EventType = authContract.L2WritesAssetEventType

package template_l2sync

import (
	authContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/auth"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Stream is the durable retry stream for L2 (MinIO) write failures.
// Shared across asset and assettemplate modules — consumers
// distinguish via Subject filter + QueueGroup.
var Stream = authContract.L2WritesStreamName

// Subject filters retries specific to the assettemplates module.
var Subject = authContract.L2WritesTemplateSubject

// Durable name for this consumer.
var Durable = config.Durable("assettemplates", "l2sync")

// QueueGroup distributes retries across replicas. Distinct from the
// assets module's "assets-l2sync" group so the two are independent.
const QueueGroup = "assettemplates-l2sync"

// EventType identifies this consumer in DLQ metadata.
const EventType = authContract.L2WritesTemplateEventType

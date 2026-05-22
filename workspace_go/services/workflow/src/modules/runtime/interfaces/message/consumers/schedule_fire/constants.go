package schedule_fire

import (
	"workflow/src/modules/runtime/application/constants"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Stream is the workflow schedule stream where fired schedules land.
var Stream = constants.ScheduleStreamName

// Subject is the filter for fired schedule messages.
var Subject = constants.ScheduleFiredSubject

// Durable name for this consumer.
var Durable = config.Durable("workflow", "schedule-fire")

// EventType for DLQ metadata.
const EventType = "schedule-fire"

// DLQServiceType is the service-type tag attached to DLQ messages produced
// by this consumer. Identifies the workflow service in cross-service DLQ
// inspection tools.
const DLQServiceType = "workflow"

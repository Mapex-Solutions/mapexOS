package refresh

import (
	"mapexVault/src/modules/credentials/application/constants"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Stream is the vault schedule stream that holds per-credential refresh timers.
var Stream = constants.VaultScheduleStreamName

// Subject is the filter for fired refresh messages.
var Subject = constants.VaultScheduleFiredSubject

// Durable is the consumer's durable name. Resolved at package init —
// e.g. "dev-mapexvault-credential-refresh-consumer".
var Durable = config.Durable("mapexvault", "credential-refresh")

// EventType is used for DLQ metadata.
const EventType = "credential-refresh"

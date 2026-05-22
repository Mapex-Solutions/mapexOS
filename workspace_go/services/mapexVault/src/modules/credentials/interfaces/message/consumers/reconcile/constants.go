package reconcile

import (
	"mapexVault/src/modules/credentials/application/constants"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Stream is the vault reconciler stream that holds the reconcile timer.
var Stream = constants.VaultReconcilerStreamName

// Subject is the filter for fired reconcile messages.
var Subject = constants.VaultReconcileFiredSubject

// Durable is the consumer's durable name. Resolved at package init —
// e.g. "dev-mapexvault-vault-reconcile-consumer".
var Durable = config.Durable("mapexvault", "vault-reconcile")

// EventType is used for DLQ metadata.
const EventType = "vault-reconcile"

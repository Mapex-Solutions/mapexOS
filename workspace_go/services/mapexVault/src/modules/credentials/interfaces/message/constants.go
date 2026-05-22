package message

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * Intra-service NATS Stream Subject Patterns
 *
 * Wildcard patterns used exclusively by the infrastructure wiring layer
 * (src/bootstrap/nats.go) to define the subject space owned by each
 * JetStream stream. They describe the consumer-facing subject topology, not
 * any individual publish target, so they live in this interfaces/message
 * package as consumer-facing topology constants.
 *
 * The concrete per-message subjects (e.g. VaultScheduleFiredSubject) are
 * published BY the application service and therefore remain in
 * application/constants.
 *
 * Stream subject names resolve at package init from GO_ENV.
 */

// VaultScheduleSubjectPattern covers every per-credential refresh timer
// published to the vault schedule stream
// (${env}.mapexos.vault.schedule.{credentialId} and
// ${env}.mapexos.vault.schedule.fired). Resolved at package init —
// e.g. "dev.mapexos.vault.schedule.>".
var VaultScheduleSubjectPattern = config.Subject("vault", "schedule") + ".>"

// VaultReconcilerSubjectPattern covers the single self-republishing
// reconcile timer on the vault reconciler stream
// (${env}.mapexos.vault.reconcile.schedule and
// ${env}.mapexos.vault.reconcile.fired). Resolved at package init —
// e.g. "dev.mapexos.vault.reconcile.>".
var VaultReconcilerSubjectPattern = config.Subject("vault", "reconcile") + ".>"

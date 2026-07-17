package constants

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

/**
 * NATS Stream Configuration
 *
 * All constants below are PUBLISHED BY the application service (credential
 * refresh scheduling, reconciler re-arm, vault event emission). Producer-side
 * intra-service NATS constants live in application/constants. The purely
 * consumer-facing wildcard subject patterns live in
 * interfaces/message/constants.go.
 *
 * Stream and subject names resolve at package init from GO_ENV via the
 * mapexGoKit config helpers so the same binary serves multiple environments
 * on a shared NATS cluster.
 */

// VaultStreamName is the JetStream stream for vault events. Resolved at
// package init — e.g. "DEV-MAPEXOS-MAPEXVAULT-EVENTS".
var VaultStreamName = config.StreamName("MAPEXVAULT", "EVENTS")

// VaultEventsSubject is the single subject for all vault events.
// Event type is discriminated by the payload action field.
// Published by: application (publishVaultEvent). Resolved at package init —
// e.g. "dev.mapexos.vault.events".
var VaultEventsSubject = config.Subject("vault", "events")

/**
 * NATS Schedule Configuration
 *
 * The vault schedule stream holds per-credential refresh timers. The
 * application publishes to VaultScheduleSubjectPrefix.{credentialId} with
 * TargetSubject VaultScheduleFiredSubject; the refresh consumer subscribes
 * to the fired subject and invokes HandleRefreshMessage.
 */

// VaultScheduleStreamName is the JetStream stream for credential refresh schedules.
// AllowMsgSchedules: true — NATS delivers to target subject at @at time.
// Published by: application (PublishScheduled, PurgeStreamSubject). Resolved
// at package init — e.g. "DEV-MAPEXOS-MAPEXVAULT-SCHEDULE".
var VaultScheduleStreamName = config.StreamName("MAPEXVAULT", "SCHEDULE")

// VaultScheduleSubjectPrefix is used to build per-credential subjects:
// ${env}.mapexos.vault.schedule.{credentialID}. Published by application
// during refresh scheduling and reconciler reseed.
var VaultScheduleSubjectPrefix = config.Subject("vault", "schedule")

// VaultScheduleFiredSubject is the target subject for fired schedules.
// MUST be inside the vault schedule stream (covered by the prefix wildcard)
// because NATS requires the @at target to be within the same stream that
// has AllowMsgSchedules: true. Published by: application (as TargetSubject
// in ScheduledPublishConfig).
var VaultScheduleFiredSubject = config.Subject("vault", "schedule.fired")

// RefreshBufferMinutes is how many minutes before token expiry to trigger refresh.
// Business rule, not NATS wiring.
const RefreshBufferMinutes = 15

// VaultReconcileDefaultIntervalSeconds is the default interval between reconcile
// cycles. Overridable via config key "vault_reconcile_interval".
// 1h keeps the worst-case recovery window short enough to react before most
// OAuth tokens expire (typical provider TTL = 1h).
const VaultReconcileDefaultIntervalSeconds = 3600

// VaultLeaderBucket is the NATS KV bucket holding the reconcile leader-election
// lease. Env-scoped so environments never share a lease; one key (the reconcile
// leader). Resolved at package init — e.g. "DEV-MAPEXOS-MAPEXVAULT-LEADER".
var VaultLeaderBucket = config.StreamName("MAPEXVAULT", "LEADER")

// VaultReconcileLeaderKey is the lease key that elects the single pod running the
// credential reseed sweep.
const VaultReconcileLeaderKey = "reconcile"

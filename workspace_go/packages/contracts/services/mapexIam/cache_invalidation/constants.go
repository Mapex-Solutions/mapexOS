package cache_invalidation

/**
 * Cross-service NATS constants for the cache_invalidation bounded context.
 *
 * These constants form the PUBLIC contract used by:
 *   - Publishers (roles, organizations, memberships, groups modules inside mapexIam)
 *   - Consumer (cache_invalidation module inside mapexIam)
 * and potentially any future external service that publishes to this wildcard.
 *
 * Stream and Subject names resolve at package init from GO_ENV so the same
 * binary serves multiple environments on a shared NATS cluster without
 * code changes.
 */

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Stream is the NATS JetStream stream name for cache invalidation events.
// Resolved at package init from GO_ENV — e.g. "DEV-MAPEXOS-MAPEXIAM-CACHE-INVALIDATION".
var Stream = config.StreamName("MAPEXIAM", "CACHE-INVALIDATION")

// subjectPrefix is the env-prefixed base from which Subject (wildcard) and
// the per-event format strings are built. Resolved at package init —
// e.g. "dev.mapexos.cache.invalidation".
var subjectPrefix = config.Subject("cache", "invalidation")

// Subject is the wildcard subject pattern that the cache_invalidation consumer binds to.
// Publishers use concrete subjects under this prefix (e.g.
// `${env}.mapexos.cache.invalidation.role.{roleId}.permissions.changed`).
var Subject = subjectPrefix + ".>"

// EventType is the DLQ classification tag for cache invalidation events.
const EventType = "cache.invalidation"

// Subject format strings for the concrete events that publishers fan out
// under the wildcard Subject. Use fmt.Sprintf with the roleId/groupId/orgId
// to build the final subject. Resolved at package init from GO_ENV.
var (
	RolePermissionsChangedSubjectFormat = subjectPrefix + ".role.%s.permissions.changed"
	RoleDeletedSubjectFormat            = subjectPrefix + ".role.%s.deleted"
	GroupChangedSubjectFormat           = subjectPrefix + ".group.%s.changed"
	GroupDeletedSubjectFormat           = subjectPrefix + ".group.%s.deleted"
	OrgHierarchyChangedSubjectFormat    = subjectPrefix + ".organization.%s.hierarchy.changed"
	OrgAccessPolicyChangedSubjectFormat = subjectPrefix + ".organization.%s.access_policy.changed"
	MembershipChangedSubjectFormat      = subjectPrefix + ".membership.%s.changed"
	MembershipDeletedSubjectFormat      = subjectPrefix + ".membership.%s.deleted"
)

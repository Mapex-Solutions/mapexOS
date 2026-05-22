package constants

import (
	definitionsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/definitions"
)

// FANOUT CONFIGURATION (NATS JetStream)
//
// FANOUT is used for cache invalidation notifications.
// Consuming services subscribe to these subjects and invalidate
// their TieredCache (L0/L1) when definitions are updated/deleted.
//
// Cross-service contracts — re-exported from
// packages/contracts/services/workflow/definitions.

// FanoutStreamName is the JetStream stream for FANOUT messages.
var FanoutStreamName = definitionsContract.FanoutStreamName

// FanoutDefinitionSubject is the subject for definition cache invalidation.
// Published on: Create, Update, Delete.
// Consumers: js-workflow-executor (future).
var FanoutDefinitionSubject = definitionsContract.FanoutDefinitionSubject

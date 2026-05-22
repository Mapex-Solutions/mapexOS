package constants

import (
	pluginsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/plugins"
)

// FANOUT CONFIGURATION (NATS JetStream)
//
// FANOUT is used for cache invalidation notifications across pods.
// When a plugin manifest is created, updated, or deleted, a FANOUT
// message is published so all pods invalidate their TieredCache (L0/L1).
//
// Cross-service contracts — re-exported from
// packages/contracts/services/workflow/plugins.

// FanoutStreamName is the JetStream stream for FANOUT messages.
var FanoutStreamName = pluginsContract.FanoutStreamName

// FanoutPluginSubject is the subject for plugin cache invalidation.
// Published on: Create, Update, Delete.
// Consumers: all workflow pods (self-subscribe via NATS Fanout).
var FanoutPluginSubject = pluginsContract.FanoutPluginSubject

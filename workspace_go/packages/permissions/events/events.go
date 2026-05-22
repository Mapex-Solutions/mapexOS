package events

// Events permissions are organized by event type (table) for granular access control.
// Each event type has its own set of CRUD permissions.

const (
	// Events Raw - Raw events from HTTP/MQTT gateways (debugging, short retention)
	EventsRawList   = "events.raw.list"
	EventsRawRead   = "events.raw.read"
	EventsRawCreate = "events.raw.create"
	EventsRawDelete = "events.raw.delete"

	// Events Processed - Main processed events table (long retention)
	EventsProcessedList   = "events.processed.list"
	EventsProcessedRead   = "events.processed.read"
	EventsProcessedCreate = "events.processed.create"
	EventsProcessedDelete = "events.processed.delete"

	// Events JS Executor - Events from JS script execution
	EventsJsExecutorList   = "events.js_executor.list"
	EventsJsExecutorRead   = "events.js_executor.read"
	EventsJsExecutorCreate = "events.js_executor.create"
	EventsJsExecutorDelete = "events.js_executor.delete"

	// Events Router - Events from router service
	EventsRouterList   = "events.router.list"
	EventsRouterRead   = "events.router.read"
	EventsRouterCreate = "events.router.create"
	EventsRouterDelete = "events.router.delete"

	// Events Business Rule - Events from business rules evaluation
	EventsBusinessRuleList   = "events.business_rule.list"
	EventsBusinessRuleRead   = "events.business_rule.read"
	EventsBusinessRuleCreate = "events.business_rule.create"
	EventsBusinessRuleDelete = "events.business_rule.delete"

	// Events Trigger - Events from trigger executions
	EventsTriggerList   = "events.trigger.list"
	EventsTriggerRead   = "events.trigger.read"
	EventsTriggerCreate = "events.trigger.create"
	EventsTriggerDelete = "events.trigger.delete"

	// Events Workflow - Workflow execution history events
	EventsWorkflowList   = "events.workflow.list"
	EventsWorkflowRead   = "events.workflow.read"
	EventsWorkflowCreate = "events.workflow.create"
	EventsWorkflowDelete = "events.workflow.delete"

	// Events Audit - Audit trail events (compliance, long retention)
	EventsAuditList   = "events.audit.list"
	EventsAuditRead   = "events.audit.read"
	EventsAuditCreate = "events.audit.create"
	EventsAuditDelete = "events.audit.delete"

	// Events Notifications - Notification events
	EventsNotificationsList   = "events.notifications.list"
	EventsNotificationsRead   = "events.notifications.read"
	EventsNotificationsCreate = "events.notifications.create"
	EventsNotificationsDelete = "events.notifications.delete"

	// Events DLQ - Dead Letter Queue events (read-only)
	EventsDLQList = "events.dlq.list"
	EventsDLQRead = "events.dlq.read"

	// Events Asset Status - Asset connectivity history (offline/online transitions)
	EventsAssetStatusList = "events.asset_status.list"
	EventsAssetStatusRead = "events.asset_status.read"

	// General Events - For admin/system-level access to all event types
	EventsList   = "events.list"   // Can list any event type
	EventsRead   = "events.read"   // Can read any event type
	EventsCreate = "events.create" // Can create any event type
	EventsDelete = "events.delete" // Can delete any event type
)

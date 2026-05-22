package permissions

// Trigger Permissions
const (
	// TriggerList - Permission to list all triggers
	TriggerList = "triggers.list"

	// TriggerCreate - Permission to create a new trigger
	TriggerCreate = "triggers.create"

	// TriggerRead - Permission to read a specific trigger
	TriggerRead = "triggers.read"

	// TriggerUpdate - Permission to update a trigger
	TriggerUpdate = "triggers.update"

	// TriggerDelete - Permission to delete a trigger
	TriggerDelete = "triggers.delete"

	// TriggerAll - Wildcard permission for all trigger operations
	TriggerAll = "triggers.*"
)

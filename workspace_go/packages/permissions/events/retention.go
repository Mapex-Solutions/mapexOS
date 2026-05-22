package events

// Retention Permissions
const (
	// RetentionList - Permission to list retention policies
	RetentionList = "retention.list"

	// RetentionRead - Permission to read a specific retention policy
	RetentionRead = "retention.read"

	// RetentionUpdate - Permission to create or update a retention policy
	RetentionUpdate = "retention.update"

	// RetentionAll - Wildcard permission for all retention operations
	RetentionAll = "retention.*"
)

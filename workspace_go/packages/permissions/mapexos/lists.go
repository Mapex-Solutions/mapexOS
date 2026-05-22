package permissions

// List Permissions
const (
	// ListList - Permission to lists all lists
	ListList = "lists.lists"

	// ListCreate - Permission to create a new lists
	ListCreate = "lists.create"

	// ListRead - Permission to read a specific lists
	ListRead = "lists.read"

	// ListUpdate - Permission to update a lists
	ListUpdate = "lists.update"

	// ListDelete - Permission to delete a lists
	ListDelete = "lists.delete"

	// ListAll - Wildcard permission for all lists operations
	ListAll = "lists.*"
)

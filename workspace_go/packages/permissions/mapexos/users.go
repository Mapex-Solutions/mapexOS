package permissions

// User Permissions
const (
	// UserList - Permission to list all users
	UserList = "users.list"

	// UserCreate - Permission to create a new users
	UserCreate = "users.create"

	// UserRead - Permission to read a specific users
	UserRead = "users.read"

	// UserUpdate - Permission to update a users
	UserUpdate = "users.update"

	// UserDelete - Permission to delete a users
	UserDelete = "users.delete"

	// UserAll - Wildcard permission for all users operations
	UserAll = "users.*"
)

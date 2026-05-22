package permissions

// Membership Permissions
const (
	// MembershipList - Permission to list all memberships
	MembershipList = "memberships.list"

	// MembershipCreate - Permission to create a new memberships
	MembershipCreate = "memberships.create"

	// MembershipRead - Permission to read a specific memberships
	MembershipRead = "memberships.read"

	// MembershipUpdate - Permission to update a memberships
	MembershipUpdate = "memberships.update"

	// MembershipDelete - Permission to delete a memberships
	MembershipDelete = "memberships.delete"

	// MembershipAll - Wildcard permission for all memberships operations
	MembershipAll = "memberships.*"
)

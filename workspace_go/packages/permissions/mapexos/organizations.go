package permissions

// Organization Permissions
const (
	// OrganizationList - Permission to list all organizations
	OrganizationList = "organizations.list"

	// OrganizationCreate - Permission to create a new organizations
	OrganizationCreate = "organizations.create"

	// OrganizationRead - Permission to read a specific organizations
	OrganizationRead = "organizations.read"

	// OrganizationUpdate - Permission to update an organizations
	OrganizationUpdate = "organizations.update"

	// OrganizationDelete - Permission to delete an organizations
	OrganizationDelete = "organizations.delete"

	// OrganizationAll - Wildcard permission for all organizations operations
	OrganizationAll = "organizations.*"
)

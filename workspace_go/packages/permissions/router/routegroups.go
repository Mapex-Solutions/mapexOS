package permissions

// RouteGroup Permissions
const (
	// RouteGroupList - Permission to list all route groups
	RouteGroupList = "routegroups.list"

	// RouteGroupCreate - Permission to create a new route group
	RouteGroupCreate = "routegroups.create"

	// RouteGroupRead - Permission to read a specific route group
	RouteGroupRead = "routegroups.read"

	// RouteGroupUpdate - Permission to update a route group
	RouteGroupUpdate = "routegroups.update"

	// RouteGroupDelete - Permission to delete a route group
	RouteGroupDelete = "routegroups.delete"

	// RouteGroupAll - Wildcard permission for all route group operations
	RouteGroupAll = "routegroups.*"
)

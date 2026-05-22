package services

import "router/src/modules/routegroups/application/constants"

// BuildRouteGroupCacheKey constructs a Redis cache key for a RouteGroup entity.
//
// This helper builds a standardized cache key using only the RouteGroup's unique identifier.
// The key follows the pattern: ROUTE_GROUP:{routeGroupId}
//
// Parameters:
//   - routeGroupId: The unique identifier (ObjectID as string) of the RouteGroup.
//
// Returns:
//   - string: A formatted cache key ready to be used with Redis operations.
//
// Example:
//
//	cacheKey := BuildRouteGroupCacheKey("507f1f77bcf86cd799439011")
//	// Returns: "ROUTE_GROUP:507f1f77bcf86cd799439011"
//
// Usage:
//
//	This function is typically used in RouteGroupService methods to:
//	- Store RouteGroup entities in cache after creation
//	- Retrieve RouteGroup entities from cache before DB queries
//	- Invalidate cache entries after updates or deletions
//
// Note: We use only the routeGroupId (without orgId) because:
//   - RouteGroups can be accessed by ID without needing to know the orgId first
//   - This avoids circular dependencies when fetching by ID
//   - The orgId is already part of the RouteGroup entity stored in cache
func BuildRouteGroupCacheKey(routeGroupId string) string {
	return constants.RouteGroupCacheKeyPrefix + ":" + routeGroupId
}

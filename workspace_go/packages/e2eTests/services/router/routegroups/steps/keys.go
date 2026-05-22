package steps

// Bag keys this package writes. Other packages reading these keys import
// the constants from here.
const (
	// BagKeyRouteGroupID is the route group id returned by
	// CreateRouteGroup. Asset creation reads it to bind the new asset
	// to the route group exercised by the saga.
	BagKeyRouteGroupID = "router.routeGroupID"

	// BagKeyOnlineRouteGroupID is the id of the route group fired by
	// the healthmonitor on the offline→online transition. Connectivity
	// action journeys create this via CreateRouteGroupWith and bind it
	// to the asset's HealthMonitor.OnlineRouteGroupIds.
	BagKeyOnlineRouteGroupID = "router.onlineRouteGroupID"

	// BagKeyOfflineRouteGroupID is the id of the route group fired by
	// the healthmonitor on the online→offline transition.
	BagKeyOfflineRouteGroupID = "router.offlineRouteGroupID"
)

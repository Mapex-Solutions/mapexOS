package constants

// eventSource discriminator values for mapexos.route.execute payloads.
// The router reads the top-level `eventSource` field to decide which
// route-group collection on the asset cache to consult.
const (
	// EventSourceAssetEvent is the default — regular IoT event routing.
	// Router uses asset.RouteGroupIds from the asset cache.
	EventSourceAssetEvent = "assetEvent"

	// EventSourceHealthStatus marks a sensor health state transition.
	// Router uses asset.HealthMonitor.OfflineRouteGroupIds or OnlineRouteGroupIds
	// depending on the nested event's eventType ("offline" | "online").
	EventSourceHealthStatus = "healthStatus"
)

// Health status values carried in the nested StandardizedPayload.eventType
// when eventSource == EventSourceHealthStatus.
const (
	HealthStatusOffline = "offline"
	HealthStatusOnline  = "online"
)

// HealthStatusAllowedRouterKinds enumerates the router kinds permitted to
// receive healthStatus events. Routers with any other kind (save_event,
// notification, lake_house, ...) are skipped silently when
// eventSource == EventSourceHealthStatus.
var HealthStatusAllowedRouterKinds = []string{"trigger", "workflow"}

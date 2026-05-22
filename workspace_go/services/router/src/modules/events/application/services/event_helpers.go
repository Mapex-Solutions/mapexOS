package services

import (
	"router/src/modules/events/application/constants"
	domainServices "router/src/modules/events/domain/services"
	routegroupPorts "router/src/modules/routegroups/application/ports"

	assetsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
)

// resolveRouteGroupIds picks the correct route groups from the asset cache
// based on the payload's top-level eventSource discriminator.
//
// Switch behavior:
//   - eventSource == EventSourceHealthStatus: read asset.HealthMonitor and
//     pick Offline/Online route groups based on the nested
//     StandardizedPayload.eventType ("offline" | "online").
//   - default (including EventSourceAssetEvent and empty): use
//     asset.RouteGroupIds.
//
// Unknown permutations (unrecognized healthType, nil HealthMonitor while
// eventSource=healthStatus) return nil — the caller iterates zero groups
// and the publish is skipped.
//
// Source of truth is always the asset cache — the payload NEVER carries
// route group IDs.
func resolveRouteGroupIds(asset *assetsContract.AssetReadModel, eventSource string, event map[string]interface{}) []string {
	switch eventSource {
	case constants.EventSourceHealthStatus:
		if asset.HealthMonitor == nil {
			return nil
		}
		healthType, _ := event["eventType"].(string)
		switch healthType {
		case constants.HealthStatusOffline:
			return asset.HealthMonitor.OfflineRouteGroupIds
		case constants.HealthStatusOnline:
			return asset.HealthMonitor.OnlineRouteGroupIds
		default:
			return nil
		}
	default:
		return asset.RouteGroupIds
	}
}

// isAllowedKindForHealthStatus reports whether a router kind is permitted
// to receive healthStatus events. Used by processRouter to silently skip
// disallowed kinds (save_event, notification, lake_house, ...) when
// eventSource == EventSourceHealthStatus.
func isAllowedKindForHealthStatus(kind string) bool {
	for _, allowed := range constants.HealthStatusAllowedRouterKinds {
		if kind == allowed {
			return true
		}
	}
	return false
}

// toMatchConfig adapts a routegroup MatchConfig (external schema) into the
// events domain MatchConfig (evaluation contract). Keeps the domain layer
// independent from the routegroups module.
func toMatchConfig(src *routegroupPorts.MatchConfig) *domainServices.MatchConfig {
	if src == nil {
		return nil
	}

	rules := make([]domainServices.MatchRule, len(src.Rules))
	for i, r := range src.Rules {
		rules[i] = domainServices.MatchRule{
			Field:    r.Field,
			Operator: r.Operator,
			Value:    r.Value,
		}
	}

	return &domainServices.MatchConfig{
		Policy: src.Policy,
		Rules:  rules,
	}
}

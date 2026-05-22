package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"assets/src/modules/assets/application/constants"
	"assets/src/modules/assets/application/converters"
	"assets/src/modules/assets/application/ports"
	"assets/src/modules/assets/domain/entities"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// generateAlphanumericPassword returns a cryptographically random
// alphanumeric string of the given length, drawn uniformly from the
// platform alphabet. Used by the GenerateMqttPassword endpoint to
// suggest a strong password to the operator without forcing them to
// invent one. crypto/rand + rejection-free uniform sampling via
// rand.Int over the alphabet length keeps the distribution fair.
func generateAlphanumericPassword(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("password length must be positive, got %d", length)
	}
	alphabet := constants.MqttPasswordAlphabet
	max := big.NewInt(int64(len(alphabet)))
	out := make([]byte, length)
	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("read crypto/rand: %w", err)
		}
		out[i] = alphabet[idx.Int64()]
	}
	return string(out), nil
}

// convertHealthMonitor is a thin alias over the package-shared converter.
// Kept for call-site brevity inside the services package; production logic
// lives in application/converters/health_monitor.go.
func convertHealthMonitor(entity *entities.HealthMonitorConfig) *contracts.HealthMonitorConfig {
	return converters.HealthMonitorEntityToContract(entity)
}

// validateHealthMonitorConfig enforces the HealthMonitor invariants at the
// Assets API boundary (Create/Update):
//
// - Empty arrays are valid — monitor-only mode (Redis/Mongo state tracking
//   and persistence to ClickHouse asset_status_history; no router publish).
// - When at least one route group is provided, every router in that group
//   MUST have kind in {trigger, workflow}; otherwise 422.
//
// Returns:
//   - nil when the config is absent, disabled, or passes the kind rule.
//   - *customErrors.ServerCustomError with Code=UNPROCESSABLE_ENTITY (422)
//     on violation; the HTTP layer surfaces this directly to the caller.
//   - Wrapped lookup error when the RouteGroup port fails critically.
//
// Missing route groups (returned as absent from the port's map) are
// treated as "skip validation for this id" — the caller layer (or the
// router's runtime skip) handles the missing-reference case.
func validateHealthMonitorConfig(
	ctx context.Context,
	routeGroupPort ports.RouteGroupPort,
	hm *contracts.HealthMonitorConfig,
) error {
	if hm == nil {
		return nil
	}
	if hm.Enabled == nil || !*hm.Enabled {
		return nil
	}

	// Dedup ids before the port call so we don't lookup the same group twice
	// when an admin puts the same id in both offline and online lists.
	idSet := make(map[string]struct{}, len(hm.OfflineRouteGroupIds)+len(hm.OnlineRouteGroupIds))
	for _, id := range hm.OfflineRouteGroupIds {
		idSet[id] = struct{}{}
	}
	for _, id := range hm.OnlineRouteGroupIds {
		idSet[id] = struct{}{}
	}
	ids := make([]string, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	kindsByGroup, err := routeGroupPort.GetRouterKindsByIds(ctx, ids)
	if err != nil {
		return fmt.Errorf("failed to fetch router kinds for HealthMonitor validation: %w", err)
	}

	for _, id := range ids {
		kinds, found := kindsByGroup[id]
		if !found {
			continue
		}
		for _, kind := range kinds {
			if !constants.HealthStatusAllowedRouterKinds[kind] {
				return &customErrors.ServerCustomError{
					Code: status.UNPROCESSABLE_ENTITY,
					Errors: []string{
						fmt.Sprintf("route group %s contains disallowed router kind %q (only trigger/workflow allowed for health monitoring)", id, kind),
					},
				}
			}
		}
	}

	return nil
}

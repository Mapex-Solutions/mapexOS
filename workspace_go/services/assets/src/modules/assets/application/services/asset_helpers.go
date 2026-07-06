package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strings"
	"time"

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

// validateLorawanDeviceHeartbeatMode rejects explicit heartbeat mode for a
// LoRaWAN end-device: a connectionless device has no explicit liveness
// producer, so only implicit (by-data) is valid. Gateways are exempt — the
// LNS Gateway Server publishes their presence explicitly.
func validateLorawanDeviceHeartbeatMode(
	protocol *contracts.ProtocolType,
	hm *contracts.HealthMonitorConfig,
) error {
	if protocol == nil || protocol.Type != "lorawan" || protocol.Lorawan == nil {
		return nil
	}
	if protocol.Lorawan.Kind != contracts.LorawanKindDevice {
		return nil
	}
	if hm.HeartbeatMode == nil || *hm.HeartbeatMode != "explicit" {
		return nil
	}
	return &customErrors.ServerCustomError{
		Code: status.UNPROCESSABLE_ENTITY,
		Errors: []string{
			"LoRaWAN end-devices support implicit heartbeat only; explicit mode has no liveness source for a connectionless device",
		},
	}
}

// validateHealthMonitorConfig enforces the HealthMonitor invariants at the
// Assets API boundary (Create/Update):
//
//   - A LoRaWAN end-device may NOT use explicit heartbeat mode. LoRaWAN is
//     connectionless at the device level, so liveness can only be inferred from
//     inbound data (implicit); explicit mode has no producer for a device.
//     Gateways are exempt — the LNS Gateway Server drives their presence
//     explicitly.
//   - Empty arrays are valid — monitor-only mode (Redis/Mongo state tracking
//     and persistence to ClickHouse asset_status_history; no router publish).
//   - When at least one route group is provided, every router in that group
//     MUST have kind in {trigger, workflow}; otherwise 422.
//
// Returns:
//   - nil when the config is absent, disabled, or passes every rule.
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
	protocol *contracts.ProtocolType,
	hm *contracts.HealthMonitorConfig,
) error {
	if hm == nil {
		return nil
	}
	if err := validateLorawanDeviceHeartbeatMode(protocol, hm); err != nil {
		return err
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

// validateAssetTopologyRequirements requires an asset template and at least one
// route group unless the asset is a LoRaWAN gateway, which is radio
// infrastructure with no data model and no telemetry of its own to route.
func validateAssetTopologyRequirements(protocol *contracts.ProtocolType, assetTemplateID string, routeGroupIds []string) error {
	if protocol != nil && protocol.Type == "lorawan" && protocol.Lorawan != nil &&
		protocol.Lorawan.Kind == contracts.LorawanKindGateway {
		return nil
	}
	if assetTemplateID == "" || len(routeGroupIds) == 0 {
		return &customErrors.ServerCustomError{
			Code: status.UNPROCESSABLE_ENTITY,
			Errors: []string{
				"asset template and at least one route group are required unless the asset is a LoRaWAN gateway",
			},
		}
	}
	return nil
}

// validateUpdateAssetTopology applies the patch over the existing asset and runs
// validateAssetTopologyRequirements on the resulting state, so a partial update
// that omits the two fields keeps the existing values while explicitly clearing
// them on a non-gateway is rejected.
func validateUpdateAssetTopology(before *entities.Asset, dto *contracts.AssetUpdate) error {
	protocol := dto.Protocol
	if protocol == nil {
		protocol = &contracts.ProtocolType{Type: before.Protocol.Type}
		if before.Protocol.Lorawan != nil {
			protocol.Lorawan = &contracts.LorawanConfig{Kind: before.Protocol.Lorawan.Kind}
		}
	}

	templateID := ""
	if !before.AssetTemplateID.IsZero() {
		templateID = before.AssetTemplateID.Hex()
	}
	if dto.AssetTemplateID != nil {
		templateID = *dto.AssetTemplateID
	}

	routeGroupIds := before.RouteGroupIds
	if dto.RouteGroupIds != nil {
		routeGroupIds = *dto.RouteGroupIds
	}

	return validateAssetTopologyRequirements(protocol, templateID, routeGroupIds)
}

// assetAttributeLabelPattern constrains a custom-attribute label to a compact,
// index-safe character set.
var assetAttributeLabelPattern = regexp.MustCompile(`^[A-Za-z0-9_.:-]+$`)

// validateAssetAttributes checks operator-defined custom attributes: bounded
// count, well-formed unique labels (the reserved mapex. prefix is rejected), a
// known kind, and a value whose type matches the kind. Value validity cannot be
// expressed with struct tags because Value is an any decoded from JSON.
func validateAssetAttributes(attrs []contracts.AssetAttribute) error {
	if len(attrs) == 0 {
		return nil
	}

	var errs []string
	if len(attrs) > contracts.AssetAttributeMaxCount {
		errs = append(errs, fmt.Sprintf("at most %d attributes are allowed", contracts.AssetAttributeMaxCount))
	}

	seen := make(map[string]bool, len(attrs))
	for i, attr := range attrs {
		switch {
		case attr.Label == "":
			errs = append(errs, fmt.Sprintf("attribute %d: label is required", i))
		case len(attr.Label) > contracts.AssetAttributeLabelMaxLen:
			errs = append(errs, fmt.Sprintf("attribute %q: label exceeds %d characters", attr.Label, contracts.AssetAttributeLabelMaxLen))
		case !assetAttributeLabelPattern.MatchString(attr.Label):
			errs = append(errs, fmt.Sprintf("attribute %q: label must match ^[A-Za-z0-9_.:-]+$", attr.Label))
		case strings.HasPrefix(attr.Label, contracts.AssetAttributeReservedPrefix):
			errs = append(errs, fmt.Sprintf("attribute %q: the %q prefix is reserved", attr.Label, contracts.AssetAttributeReservedPrefix))
		case seen[attr.Label]:
			errs = append(errs, fmt.Sprintf("attribute %q: duplicate label", attr.Label))
		}
		if attr.Label != "" {
			seen[attr.Label] = true
		}
		if msg := validateAssetAttributeValue(attr); msg != "" {
			errs = append(errs, msg)
		}
	}

	if len(errs) > 0 {
		return &customErrors.ServerCustomError{Code: status.UNPROCESSABLE_ENTITY, Errors: errs}
	}
	return nil
}

// validateAssetAttributeValue returns a message when attr.Value does not match
// attr.Kind, or "" when it is valid. JSON decodes numbers to float64 and objects
// to map[string]any, which the kind checks account for.
func validateAssetAttributeValue(attr contracts.AssetAttribute) string {
	switch attr.Kind {
	case contracts.AssetAttributeKindString:
		if _, ok := attr.Value.(string); !ok {
			return fmt.Sprintf("attribute %q: value must be a string", attr.Label)
		}
	case contracts.AssetAttributeKindInteger:
		if !isIntegerValue(attr.Value) {
			return fmt.Sprintf("attribute %q: value must be an integer", attr.Label)
		}
	case contracts.AssetAttributeKindBoolean:
		if _, ok := attr.Value.(bool); !ok {
			return fmt.Sprintf("attribute %q: value must be a boolean", attr.Label)
		}
	case contracts.AssetAttributeKindDate:
		s, ok := attr.Value.(string)
		if !ok {
			return fmt.Sprintf("attribute %q: value must be an ISO-8601 date string", attr.Label)
		}
		if _, err := time.Parse(time.RFC3339, s); err != nil {
			return fmt.Sprintf("attribute %q: value must be an ISO-8601 date", attr.Label)
		}
	case contracts.AssetAttributeKindGeo:
		if !isGeoValue(attr.Value) {
			return fmt.Sprintf("attribute %q: value must be an object with lat in [-90,90] and lon in [-180,180]", attr.Label)
		}
	default:
		return fmt.Sprintf("attribute %q: unknown kind %q", attr.Label, attr.Kind)
	}
	return ""
}

// isIntegerValue reports whether v is a whole number. JSON numbers arrive as
// float64; native Go integer types are also accepted.
func isIntegerValue(v any) bool {
	switch n := v.(type) {
	case float64:
		return n == math.Trunc(n)
	case float32:
		return float64(n) == math.Trunc(float64(n))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

// isGeoValue reports whether v is a { lat, lon } object with lat in [-90,90] and
// lon in [-180,180]. JSON objects decode to map[string]any, numbers to float64.
func isGeoValue(v any) bool {
	m, ok := v.(map[string]any)
	if !ok {
		return false
	}
	lat, latOK := attributeFloat(m["lat"])
	lon, lonOK := attributeFloat(m["lon"])
	return latOK && lonOK && lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}

// attributeFloat coerces a JSON-decoded numeric value to float64.
func attributeFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

package converters

import (
	"assets/src/modules/assets/domain/entities"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
)

// HealthMonitorEntityToContract maps the domain HealthMonitorConfig
// (value bool/int fields) to the cross-service contract HealthMonitorConfig
// (pointer fields with omitempty). Returns nil when the entity has no
// HealthMonitor configured.
//
// Both the application response path and the MinIO read-model writer share
// this single converter so the wire shape is byte-identical regardless of
// which layer produced it. The jinzhu/copier mapper is unreliable for
// bool↔*bool and int↔*int conversions, so spell it out manually.
func HealthMonitorEntityToContract(entity *entities.HealthMonitorConfig) *contracts.HealthMonitorConfig {
	if entity == nil {
		return nil
	}
	enabled := entity.Enabled
	threshold := entity.ThresholdMinutes
	required := entity.RequiredMisses
	mode := entity.ResolvedMode()
	hm := &contracts.HealthMonitorConfig{
		Enabled:          &enabled,
		ThresholdMinutes: &threshold,
		RequiredMisses:   &required,
		HeartbeatMode:    &mode,
	}
	if len(entity.OfflineRouteGroupIds) > 0 {
		hm.OfflineRouteGroupIds = append([]string(nil), entity.OfflineRouteGroupIds...)
	}
	if len(entity.OnlineRouteGroupIds) > 0 {
		hm.OnlineRouteGroupIds = append([]string(nil), entity.OnlineRouteGroupIds...)
	}
	return hm
}

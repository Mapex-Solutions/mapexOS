package di

import (
	"assets/src/bootstrap"
	assetPorts "assets/src/modules/assets/application/ports"
	"assets/src/modules/healthmonitor/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// HealthMonitorServiceDI aggregates all dependencies required by the HealthMonitorService.
type HealthMonitorServiceDI struct {
	dig.In

	HealthRepo      ports.HealthRepository
	AlertPublisher  ports.AlertPublisherPort
	AssetRepo       assetPorts.AssetRepository
	Metrics         *bootstrap.AssetsMetrics
	ScheduleManager natsModel.ScheduleManager `name:"core"`
}

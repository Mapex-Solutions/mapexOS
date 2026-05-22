package services

import (
	"assets/src/modules/healthmonitor/application/di"
)

// HealthMonitorService orchestrates sensor health monitoring use cases.
//
// Public surface (port + lifecycle):
//   - OnMount            — bootstraps the periodic scan on startup
//   - HandleHeartbeat    — processes heartbeat NATS messages
//   - RunScan            — stale-sensor scan entry point (scheduled via NATS)
//
// Private work lives in healthmonitor_handler_{heartbeat,scanner,schedule}.go.
type HealthMonitorService struct {
	deps      di.HealthMonitorServiceDI
	batchSize int
}

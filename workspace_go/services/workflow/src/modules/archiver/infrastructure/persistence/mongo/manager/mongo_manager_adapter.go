package manager

import (
	"time"

	"workflow/src/modules/archiver/application/ports"

	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
)

// New creates a MongoManagerAdapter bound to the provided concrete manager
// and returns it as the archiver MongoManagerPort.
func New(mgr *mongoManager.MongoManager) ports.MongoManagerPort {
	return &MongoManagerAdapter{mgr: mgr}
}

// GetBackpressureMode maps the concrete MongoManager backpressure mode
// into the archiver port's BackpressureMode.
func (a *MongoManagerAdapter) GetBackpressureMode() ports.BackpressureMode {
	return toPortMode(a.mgr.GetBackpressureMode())
}

// WriteP99 returns the current P99 write latency in milliseconds.
func (a *MongoManagerAdapter) WriteP99() int64 {
	return a.mgr.WriteP99()
}

// RecordWriteLatency records a single write duration into the backpressure window.
func (a *MongoManagerAdapter) RecordWriteLatency(d time.Duration) {
	a.mgr.RecordWriteLatency(d)
}

// toPortMode translates the concrete manager enum to the port enum.
func toPortMode(m mongoManager.BackpressureMode) ports.BackpressureMode {
	switch m {
	case mongoManager.Throttled:
		return ports.BackpressureThrottled
	case mongoManager.Backoff:
		return ports.BackpressureBackoff
	default:
		return ports.BackpressureNormal
	}
}

package ports

import "time"

/*
 * MONGO MANAGER PORT
 * Abstracts the backpressure-tracking surface of the MongoDB manager.
 * Keeps the application layer decoupled from the concrete
 * infrastructure/mongodb/manager package.
 */

// BackpressureMode represents the current MongoDB write pressure level
// as seen by the Archiver service.
type BackpressureMode int

const (
	// BackpressureNormal indicates MongoDB is responding within acceptable latency.
	BackpressureNormal BackpressureMode = 0

	// BackpressureThrottled indicates elevated P99 write latency.
	BackpressureThrottled BackpressureMode = 1

	// BackpressureBackoff indicates critically high P99 write latency.
	BackpressureBackoff BackpressureMode = 2
)

// MongoManagerPort exposes ONLY the backpressure-tracking methods that the
// archiver service consumes from the underlying MongoDB manager.
type MongoManagerPort interface {
	// GetBackpressureMode returns the current write backpressure level.
	GetBackpressureMode() BackpressureMode

	// WriteP99 returns the current P99 write latency in milliseconds.
	WriteP99() int64

	// RecordWriteLatency records a single write duration sample into the
	// backpressure window for P99 computation.
	RecordWriteLatency(d time.Duration)
}

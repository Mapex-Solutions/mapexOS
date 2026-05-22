package constants

/**
 * Application-layer constants for the events module.
 *
 * These are default fallback values applied when the corresponding
 * configuration keys are absent or invalid.
 */

// DefaultWorkerCount is the default number of concurrent workers used
// by the trigger executor when `trigger_executor_workers` is not set
// or resolves to a non-positive value.
const DefaultWorkerCount = 50

// DefaultBatchSize is the default NATS batch size used by event consumers
// when `nats_batch_size` is not set.
const DefaultBatchSize = 10

// DefaultFetchTimeoutSeconds is the default NATS fetch timeout (in seconds)
// used by event consumers when `nats_fetch_timeout` is not set.
const DefaultFetchTimeoutSeconds = 30

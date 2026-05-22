package constants

import (
	"time"

	archiverContract "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/archiver"
)

/*
 * ARCHIVER BATCH CONFIGURATION
 */

const (
	// BatchSize is the number of WORKFLOW-STATE messages fetched per batch.
	// Tuned for throughput: 5000 messages × ~200B avg = ~1MB per batch.
	BatchSize = 5000

	// MaxBatchRetries is the max retries for a failed batch MongoDB write.
	MaxBatchRetries = 3
)

/*
 * BACKPRESSURE BEHAVIOR
 * Pause durations applied between batches when MongoDB write latency
 * crosses the configured Backoff threshold.
 */

// BackpressureBackoffPause is the pause applied before processing a batch
// when MongoDB is in Backoff mode (P99 above Backoff threshold).
const BackpressureBackoffPause = 2 * time.Second

/*
 * ARCHIVED EXECUTION TTL
 * Terminal executions are kept in MongoDB long enough for the UI to query
 * them after a run, then auto-purged via TTL index on expireAt.
 */

// TerminalExecutionTTL is the time-to-live applied to terminal workflow
// executions persisted to MongoDB. After this duration MongoDB removes
// the document via the expireAt TTL index.
const TerminalExecutionTTL = 3 * 24 * time.Hour

/*
 * BATCH TAGS
 * Tags attached to MsgRef.Batch entries to drive ack/nack routing per
 * archiver state-event batch type.
 */

const (
	BatchTagCreated  = "created"
	BatchTagWaiting  = "waiting"
	BatchTagResumed  = "resumed"
	BatchTagTerminal = "terminal"
)

// EventsWorkflowSubject is the subject for workflow execution history events.
// Cross-service contract — re-exported from packages/contracts/services/workflow/archiver.
var EventsWorkflowSubject = archiverContract.SubjectEventsWorkflow

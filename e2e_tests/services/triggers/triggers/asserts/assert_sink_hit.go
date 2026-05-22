// Package asserts holds saga oracles for the triggers/triggers module.
package asserts

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	triggerSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/steps"
)

// AssertSinkHitEventually polls the in-process trigger sink counter
// until it observes at least `expected` hits or the timeout elapses.
// Used by the connectivity-action phase2_trigger journey to confirm
// the triggers service POSTed against the saga-managed sink without
// depending on the ClickHouse insert pipeline (which carries a batch
// buffer that can push the visibility of a single row past a tight
// polling budget).
//
// Reads (bag):
//   - triggerSteps.BagKeyTriggerSinkHits  *atomic.Int64  set by StartTestSink
//
// The expected count is cumulative across the journey: pass 1 after
// the first transition (online), 2 after the second (offline).
func AssertSinkHitEventually(expected int64) saga.Assert {
	return AssertSinkHitEventuallyWithTimeout(expected, 90*time.Second, 500*time.Millisecond)
}

// AssertSinkHitEventuallyWithTimeout overrides the polling budget.
// The trigger MS appears to retry on transient failures (observed
// ~60s end-to-end before the first POST lands in dev), so the default
// is generous to absorb that latency without flaking.
func AssertSinkHitEventuallyWithTimeout(expected int64, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: fmt.Sprintf("triggers/triggers.AssertSinkHitEventually[%d]", expected),
		Check: func(c *saga.Context) error {
			v, ok := c.Get(triggerSteps.BagKeyTriggerSinkHits)
			if !ok {
				return fmt.Errorf("bag key %q missing — did StartTestSink run?", triggerSteps.BagKeyTriggerSinkHits)
			}
			hits, ok := v.(*atomic.Int64)
			if !ok {
				return fmt.Errorf("bag key %q is not *atomic.Int64 (%T)", triggerSteps.BagKeyTriggerSinkHits, v)
			}

			deadline := time.Now().Add(timeout)
			for {
				got := hits.Load()
				if got >= expected {
					return nil
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("sink hits: want >=%d, got %d after %v", expected, got, timeout)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("sink-hit poll cancelled: %w", c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}

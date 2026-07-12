package saga

import (
	"context"
	"testing"
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
)

// AfterRun, when set, is called once per journey AFTER its items and all their step
// Compensates have run, with that journey's Context (authenticated clients + RunID).
// It is an optional test-observability hook — nil by default; the suite runner sets
// it to a leak check. It never affects the saga result path (a panic/error inside it
// is the hook's own concern).
var AfterRun func(c *Context)

// Run executes a journey: it walks Items in order, records executed Steps, and runs
// every registered Compensate in reverse on completion (success OR failure) so the
// live stack ends where the next journey can start. On completion the optional
// AfterRun hook observes the finished run (e.g. a leak check). The clients are threaded
// through the Context (Run does not own their lifecycle); runID is the per-journey
// unique tag, synthesized from the time when empty.
func Run(t *testing.T, ctx context.Context, runID string, clients ClientSet, items ...Item) {
	t.Helper()
	if runID == "" {
		runID = time.Now().UTC().Format("20060102-150405")
	}
	sctx := newContext(t, ctx, runID, clients)
	if name, err := run(sctx, items...); err != nil {
		t.Errorf("[SAGA] %s failed: %v", name, err)
	}
}

// run is the testable core: it walks items, records executed steps for rollback, and
// RETURNS the first failing name+error instead of calling t.Errorf. Reporting is left
// to Run so tests can drive the failure path without failing the test binary (a
// testing.TB cannot be faked). Deferred order (LIFO): rollback is registered last so it
// runs first, then AfterRun; the rollback closure captures `executed` so the reverse
// walk sees the full list at call time (a past bug captured it at len 0, making
// rollback a silent no-op). On t.Fatal inside an item, Goexit still unwinds both defers.
func run(sctx *Context, items ...Item) (failedName string, failErr error) {
	sctx.T.Logf("[SAGA] start runID=%s items=%d", sctx.RunID, len(items))

	executed := make([]Item, 0, len(items))
	// Registered before rollback so it runs AFTER rollback (LIFO): compensations have
	// run and the stack is still up, so AfterRun can observe leftovers.
	if AfterRun != nil {
		defer func() { AfterRun(sctx) }()
	}
	defer func() { rollback(sctx, executed) }()

	for i, item := range items {
		sctx.T.Logf("[SAGA] %d/%d %s", i+1, len(items), item.GetName())
		if err := item.Execute(sctx); err != nil {
			return item.GetName(), err
		}
		executed = append(executed, item)
	}

	sctx.T.Logf("[SAGA] all %d items passed; running compensations", len(executed))
	return "", nil
}

// rollback walks the executed list in reverse and invokes Rollback on every
// item. Compensation failures are logged but do not abort the rollback —
// the goal is best-effort cleanup so subsequent runs can proceed even when
// one compensation step is unhappy.
func rollback(c *Context, executed []Item) {
	for i := len(executed) - 1; i >= 0; i-- {
		item := executed[i]
		if err := item.Rollback(c); err != nil {
			c.T.Logf("[SAGA] compensation %s failed (non-fatal): %v", item.GetName(), err)
		}
	}
}

// NewClientSet constructs a ClientSet wired against every service URL.
// Auth headers are populated later by the auth step that publishes the JWT
// to the bag and propagates it across every client via ClientSet.SetBearer.
func NewClientSet(urls ClientURLs) ClientSet {
	build := func(baseURL string) *httpclient.HTTPClient {
		return httpclient.New(httpclient.Config{BaseURL: baseURL})
	}
	return ClientSet{
		HTTP:     build(urls.MapexIam),
		Assets:   build(urls.Assets),
		Router:   build(urls.Router),
		Gateway:  build(urls.Gateway),
		Events:   build(urls.Events),
		Triggers: build(urls.Triggers),
		Workflow: build(urls.Workflow),
	}
}

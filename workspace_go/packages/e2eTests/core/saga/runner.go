package saga

import (
	"context"
	"testing"
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
)

// Run executes a journey: walks Items in order, records executed Steps for
// rollback, and runs every registered Compensate in reverse on completion.
// Compensations run on both success and failure so the live stack ends in
// the same state where the next journey can start.
//
// The clients argument carries the HTTP wrappers prepared by the caller.
// Run does not own client lifecycle — it just threads them through the
// Context so each Item can drive the stack.
//
// runID is the per-journey unique tag callers inject into payload
// identifiers. When empty, Run synthesizes one from the current time so
// every execution still gets prefix-based isolation.
func Run(t *testing.T, ctx context.Context, runID string, clients ClientSet, items ...Item) {
	t.Helper()
	if runID == "" {
		runID = time.Now().UTC().Format("20060102-150405")
	}

	sctx := newContext(t, ctx, runID, clients)
	t.Logf("[SAGA] start runID=%s items=%d", runID, len(items))

	executed := make([]Item, 0, len(items))
	failed := false

	defer rollback(sctx, executed)

	for i, item := range items {
		t.Logf("[SAGA] %d/%d %s", i+1, len(items), item.GetName())
		if err := item.Execute(sctx); err != nil {
			t.Errorf("[SAGA] %s failed: %v", item.GetName(), err)
			failed = true
			break
		}
		executed = append(executed, item)
	}

	if !failed {
		t.Logf("[SAGA] all %d items passed; running compensations", len(executed))
	}
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

// ClientURLs groups the per-service base URLs the saga journey speaks to.
// Each saga test derives this from common/constants so the URL configuration
// stays in one place and journeys do not hard-code endpoints.
type ClientURLs struct {
	MapexIam string
	Assets   string
	Router   string
	Gateway  string
	Events   string
	Triggers string
	Workflow string
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

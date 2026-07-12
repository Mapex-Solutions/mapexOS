// Package leakcheck is the suite's post-journey leak detector. Wired as saga.AfterRun
// by the runner, it sweeps the public list endpoints after each journey's Compensates
// and fails the journey if any entity still carries that journey's RunID — a signal
// that a Compensate did not clean up what its Do created.
package leakcheck

import (
	saga "github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// standardSpecs are the public list endpoints checked for RunID-tagged leftovers,
// built from the journey's authenticated clients. It observes ONLY the public API and
// only the first page, so it catches leaks without false positives (it matches this
// run's RunID exactly) but may miss one buried past the first page — a best-effort
// safety net. Extend as more entity kinds gain a public list endpoint.
//
// NOTE: the endpoint paths and name fields below are the expected shapes but have not
// yet been validated against a live stack (the suite could not run end-to-end here).
// A wrong path/field degrades to a skipped sweep (logged), never a false failure —
// but confirm them on the first real run and add the missing entity kinds.
func standardSpecs(c *saga.Context) []saga.LeakSpec {
	// A large limit reduces the first-page-only blind spot when many parallel
	// journeys share the tenant; the sweep still matches this run's RunID exactly.
	return []saga.LeakSpec{
		{Kind: "asset", Client: c.Clients.Assets, ListPath: "/api/v1/assets?limit=500", NameField: "name"},
		{Kind: "trigger", Client: c.Clients.Triggers, ListPath: "/api/v1/triggers?limit=500", NameField: "name"},
	}
}

// Hook is the saga.AfterRun hook: sweep for entities still carrying c.RunID and fail
// the journey if any survived its rollback. SweepLeaks returns any leaks it DID find
// alongside a partial error, so leaks are reported first; a sweep error (endpoint
// unreachable, not authenticated, wrong path) is then logged and skipped — the check
// is a safety net, not a gate that wrongly fails a journey on an infra hiccup.
func Hook(c *saga.Context) {
	leaks, err := saga.SweepLeaks(c, standardSpecs(c))
	if len(leaks) > 0 {
		c.T.Errorf("[LEAK] %d entity(ies) survived rollback for runID=%s: %+v", len(leaks), c.RunID, leaks)
	}
	if err != nil {
		c.T.Logf("[LEAK] sweep partially skipped for runID=%s (non-fatal): %v", c.RunID, err)
	}
}

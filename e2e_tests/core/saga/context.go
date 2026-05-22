package saga

import (
	"context"
	"sync"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
)

// Context is threaded through every Item executed by Run. It carries the
// testing handle, the cancellation context, the typed clients used by the
// stack-facing actions, and a free-form bag where steps publish outputs for
// downstream consumers (e.g. CreateOrganization writes orgID; subsequent
// steps read it without taking the value through closure plumbing).
type Context struct {
	// T is the testing handle. Steps that need to fail the test directly
	// (rare; prefer returning errors from Do/Check) call methods on T.
	T *testing.T

	// Stdctx is the standard library context used by every HTTP / NATS /
	// Mongo call inside steps. Cancellation propagates through the journey.
	Stdctx context.Context

	// Clients are the typed wrappers used to drive the stack. The HTTP
	// client is preconfigured with the JWT and X-Org-Context populated by
	// the auth step so steps that come after authentication do not need to
	// know how to attach headers.
	Clients ClientSet

	// RunID is a per-journey unique tag injected into payload identifiers
	// (org name, user email, asset name) so cleanup-by-prefix works and
	// parallel runs do not collide.
	RunID string

	bag map[string]any
	mu  sync.Mutex
}

// ClientSet groups one HTTPClient per platform service. Steps invoke the
// client matching the service whose endpoint they target. Authentication
// state (JWT and X-Org-Context) is propagated to every client by the auth
// step so callers do not have to remember which client carries the token.
//
// All clients are mapexGoKit httpclient.HTTPClient instances driven through
// the Raw method — saga steps assert on resp.StatusCode directly, so the
// auto-decoding Get/Post entry points are not what they need.
type ClientSet struct {
	// HTTP is the mapexIam-facing HTTPClient. Kept as the primary field
	// because IAM is the bootstrap entry point: every saga starts by
	// driving mapexIam endpoints (org, role, onboarding, auth/login).
	HTTP *httpclient.HTTPClient

	// Assets points at the assets service base URL.
	Assets *httpclient.HTTPClient

	// Router points at the router service base URL.
	Router *httpclient.HTTPClient

	// Gateway points at the http_gateway service base URL.
	Gateway *httpclient.HTTPClient

	// Events points at the events service base URL. Saga journeys that
	// observe pipeline outcomes (presence events, telemetry events) read
	// from this client; they never subscribe to internal NATS subjects.
	Events *httpclient.HTTPClient

	// Triggers points at the triggers service base URL. Connectivity
	// action journeys create triggers via POST /api/v1/triggers and
	// reference the returned id in route groups.
	Triggers *httpclient.HTTPClient

	// Workflow points at the workflow service base URL. Connectivity
	// action journeys create workflow definitions and instances via
	// POST /api/v1/definitions and POST /api/v1/instances; route groups
	// of kind=workflow reference the returned instance id.
	Workflow *httpclient.HTTPClient
}

// All returns every HTTP client in the set so helper code (e.g. the auth
// step) can apply a single change — typically the bearer token and the
// X-Org-Context header — across the full set without enumerating fields.
func (s ClientSet) All() []*httpclient.HTTPClient {
	out := make([]*httpclient.HTTPClient, 0, 7)
	for _, c := range []*httpclient.HTTPClient{s.HTTP, s.Assets, s.Router, s.Gateway, s.Events, s.Triggers, s.Workflow} {
		if c != nil {
			out = append(out, c)
		}
	}
	return out
}

// SetBearer registers the bearer token on every client in the set. Callers
// pass an empty string to clear the header — useful in cleanup paths and
// when a step deliberately drops to anonymous traffic.
func (s ClientSet) SetBearer(token string) {
	value := ""
	if token != "" {
		value = "Bearer " + token
	}
	for _, c := range s.All() {
		c.SetHeader("Authorization", value)
	}
}

// SetOrgContext registers the X-Org-Context header on every client in the
// set. Empty value removes the header.
func (s ClientSet) SetOrgContext(orgID string) {
	for _, c := range s.All() {
		c.SetHeader("X-Org-Context", orgID)
	}
}

// Set publishes a value to the bag for downstream Items. Concurrent calls
// are safe because steps may dispatch goroutines internally; the saga
// runner itself walks Items sequentially.
func (c *Context) Set(key string, val any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bag[key] = val
}

// Get returns the value previously published to key. The bool is false when
// the key was never set.
func (c *Context) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.bag[key]
	return v, ok
}

// MustGetString fetches a string from the bag, failing the test fast when
// the key is missing or holds a non-string value. Steps that consume bag
// inputs use this to express the contract: "I require this key to exist
// before I can execute" — surfaced as a clear test failure rather than a
// nil dereference.
func (c *Context) MustGetString(key string) string {
	v, ok := c.Get(key)
	if !ok {
		c.T.Fatalf("[SAGA] missing required bag key %q (step out of order?)", key)
	}
	s, ok := v.(string)
	if !ok {
		c.T.Fatalf("[SAGA] bag key %q is not a string (got %T)", key, v)
	}
	return s
}

// newContext is invoked by Run to materialize the per-journey Context.
// Kept package-private so callers cannot bypass the runner contract.
func newContext(t *testing.T, ctx context.Context, runID string, clients ClientSet) *Context {
	return &Context{
		T:       t,
		Stdctx:  ctx,
		Clients: clients,
		RunID:   runID,
		bag:     make(map[string]any),
	}
}

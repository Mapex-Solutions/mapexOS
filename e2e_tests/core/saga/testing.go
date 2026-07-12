package saga

import (
	"context"
	"testing"
)

// NewTestContext builds a Context for unit tests of steps and asserts that live
// outside this package. It is the sanctioned way to construct a Context without
// going through Run — production journeys must always go through Run; this exists
// only so a step's Do/Compensate or an assert's Check can be exercised against an
// httptest.Server. RunID defaults to "test" so cleanup-by-prefix helpers have a
// stable value. Build the ClientSet with NewClientSet(ClientURLs{...}) pointed at
// the test server.
func NewTestContext(t *testing.T, clients ClientSet) *Context {
	return newContext(t, context.Background(), "test", clients)
}

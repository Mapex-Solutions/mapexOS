package saga

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSweepLeaks_FindsRunIDTaggedItems proves the sweep reports only entities whose
// name carries the RunID and leaves the rest alone.
func TestSweepLeaks_FindsRunIDTaggedItems(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// One name carries the default test RunID ("test"), one does not.
		_, _ = w.Write([]byte(`{"data":{"items":[{"id":"1","name":"asset-test-abc"},{"id":"2","name":"asset-other"}]}}`))
	}))
	defer srv.Close()

	c := NewTestContext(t, NewClientSet(ClientURLs{Assets: srv.URL}))
	specs := []LeakSpec{{Kind: "asset", Client: c.Clients.Assets, ListPath: "/api/v1/assets", NameField: "name"}}

	leaks, err := SweepLeaks(c, specs)
	if err != nil {
		t.Fatalf("SweepLeaks: %v", err)
	}
	if len(leaks) != 1 {
		t.Fatalf("leaks = %v, want exactly 1 (the RunID-tagged one)", leaks)
	}
	if leaks[0].ID != "1" || leaks[0].Kind != "asset" || leaks[0].Name != "asset-test-abc" {
		t.Errorf("leak = %+v, want {asset 1 asset-test-abc}", leaks[0])
	}
}

// TestSweepLeaks_CleanRun_NoLeaks proves an empty list yields no leaks and no error.
func TestSweepLeaks_CleanRun_NoLeaks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"data":{"items":[]}}`))
	}))
	defer srv.Close()

	c := NewTestContext(t, NewClientSet(ClientURLs{Assets: srv.URL}))
	leaks, err := SweepLeaks(c, []LeakSpec{{Kind: "asset", Client: c.Clients.Assets, ListPath: "/api/v1/assets", NameField: "name"}})
	if err != nil {
		t.Fatalf("SweepLeaks: %v", err)
	}
	if len(leaks) != 0 {
		t.Fatalf("leaks = %v, want none (clean run)", leaks)
	}
}

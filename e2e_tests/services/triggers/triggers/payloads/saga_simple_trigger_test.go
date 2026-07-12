package payloads

import (
	"strings"
	"testing"
)

// TestSagaSimpleTrigger proves the builder writes the passed sinkURL into the HTTP
// endpoint and stamps the runID into the name — no fixed constant is read.
func TestSagaSimpleTrigger(t *testing.T) {
	const runID, sinkURL = "run-abc", "http://localhost:54321"

	p := SagaSimpleTrigger(runID, sinkURL)

	cfg := p["config"].(map[string]any)
	httpCfg := cfg["http"].(map[string]any)
	if got := httpCfg["endpoint"]; got != sinkURL {
		t.Errorf("config.http.endpoint = %v, want %q (the passed sinkURL)", got, sinkURL)
	}
	if name := p["name"].(string); !strings.Contains(name, runID) {
		t.Errorf("name = %q, want it to carry runID %q", name, runID)
	}
}

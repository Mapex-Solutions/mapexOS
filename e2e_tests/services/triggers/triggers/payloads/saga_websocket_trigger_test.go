package payloads

import "testing"

// TestSagaWebsocketTrigger proves the builder writes the passed wsURL into the
// websocket url field.
func TestSagaWebsocketTrigger(t *testing.T) {
	const runID, wsURL = "run-xyz", "ws://localhost:54321/ws"

	p := SagaWebsocketTrigger(runID, wsURL)

	cfg := p["config"].(map[string]any)
	wsCfg := cfg["websocket"].(map[string]any)
	if got := wsCfg["url"]; got != wsURL {
		t.Errorf("config.websocket.url = %v, want %q (the passed wsURL)", got, wsURL)
	}
}

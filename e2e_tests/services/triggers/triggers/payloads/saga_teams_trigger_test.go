package payloads

import (
	"strings"
	"testing"
)

// TestSagaTeamsTrigger proves the builder writes the passed sinkURL into the Teams
// webhookUrl and embeds the runID in the text.
func TestSagaTeamsTrigger(t *testing.T) {
	const runID, sinkURL = "run-teams", "http://localhost:54321"

	p := SagaTeamsTrigger(runID, sinkURL)

	cfg := p["config"].(map[string]any)
	teamsCfg := cfg["teams"].(map[string]any)
	if got := teamsCfg["webhookUrl"]; got != sinkURL {
		t.Errorf("config.teams.webhookUrl = %v, want %q (the passed sinkURL)", got, sinkURL)
	}
	if text := teamsCfg["text"].(string); !strings.Contains(text, runID) {
		t.Errorf("text = %q, want it to carry runID %q", text, runID)
	}
}

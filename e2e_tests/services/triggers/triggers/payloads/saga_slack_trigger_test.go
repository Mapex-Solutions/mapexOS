package payloads

import (
	"strings"
	"testing"
)

// TestSagaSlackTrigger proves the builder writes the passed sinkURL into the Slack
// webhookUrl and embeds the runID in the message.
func TestSagaSlackTrigger(t *testing.T) {
	const runID, sinkURL = "run-slack", "http://localhost:54321"

	p := SagaSlackTrigger(runID, sinkURL)

	cfg := p["config"].(map[string]any)
	slackCfg := cfg["slack"].(map[string]any)
	if got := slackCfg["webhookUrl"]; got != sinkURL {
		t.Errorf("config.slack.webhookUrl = %v, want %q (the passed sinkURL)", got, sinkURL)
	}
	if msg := slackCfg["message"].(string); !strings.Contains(msg, runID) {
		t.Errorf("message = %q, want it to carry runID %q", msg, runID)
	}
}

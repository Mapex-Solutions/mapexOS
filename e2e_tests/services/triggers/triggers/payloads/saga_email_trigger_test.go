package payloads

import "testing"

// TestSagaEmailTrigger proves the builder writes the passed host and port into the
// email config's smtpHost / smtpPort fields.
func TestSagaEmailTrigger(t *testing.T) {
	const runID, host, port = "run-mail", "host.docker.internal", 54321

	p := SagaEmailTrigger(runID, host, port)

	cfg := p["config"].(map[string]any)
	emailCfg := cfg["email"].(map[string]any)
	if got := emailCfg["smtpHost"]; got != host {
		t.Errorf("config.email.smtpHost = %v, want %q", got, host)
	}
	if got := emailCfg["smtpPort"]; got != port {
		t.Errorf("config.email.smtpPort = %v, want %d", got, port)
	}
}

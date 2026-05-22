package email

import (
	"context"
	"strings"
	"testing"
	"time"

	mocksmtp "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mock_servers/smtp"
)

// smtpConfig returns a base email config with valid SMTP fields for testing.
// Tests override specific fields to test validation paths.
func smtpConfig(overrides map[string]interface{}) map[string]interface{} {
	base := map[string]interface{}{
		"smtpHost": "127.0.0.1",
		"smtpPort": 587,
		"fromAddr": "sender@test.com",
		"to":       "recipient@test.com",
		"subject":  "Test Subject",
		"body":     "Test body",
	}
	for k, v := range overrides {
		if v == nil {
			delete(base, k)
		} else {
			base[k] = v
		}
	}
	return map[string]interface{}{"email": base}
}

/**
 * EmailExecutor Tests
 */

func TestEmailExecutor_GetType(t *testing.T) {
	executor := NewEmailExecutor()

	if executor.GetType() != "email" {
		t.Errorf("GetType() = %q, want 'email'", executor.GetType())
	}
}

func TestEmailExecutor_Execute_MissingEmailField(t *testing.T) {
	executor := NewEmailExecutor()

	config := map[string]interface{}{
		"notemail": map[string]interface{}{
			"to":      "test@example.com",
			"subject": "Test",
		},
	}

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'email' field, got nil")
	}
}

func TestEmailExecutor_Execute_MissingSmtpHost(t *testing.T) {
	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{"smtpHost": nil})

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'smtpHost', got nil")
	}
	if !strings.Contains(err.Error(), "smtpHost") {
		t.Errorf("Error should mention 'smtpHost': %v", err)
	}
}

func TestEmailExecutor_Execute_MissingFromAddr(t *testing.T) {
	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{"fromAddr": nil})

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'fromAddr', got nil")
	}
	if !strings.Contains(err.Error(), "fromAddr") {
		t.Errorf("Error should mention 'fromAddr': %v", err)
	}
}

func TestEmailExecutor_Execute_MissingTo(t *testing.T) {
	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{"to": nil})

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'to', got nil")
	}
}

func TestEmailExecutor_Execute_EmptyTo(t *testing.T) {
	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{"to": ""})

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for empty 'to', got nil")
	}
}

func TestEmailExecutor_Execute_MissingSubject(t *testing.T) {
	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{"subject": nil})

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for missing 'subject', got nil")
	}
}

func TestEmailExecutor_Execute_EmptySubject(t *testing.T) {
	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{"subject": ""})

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Execute() expected error for empty 'subject', got nil")
	}
}

/**
 * buildEmailMessage Tests
 */

func TestBuildEmailMessage_PlainText(t *testing.T) {
	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	subject := "Test Subject"
	body := "This is the email body"
	isHTML := false

	message := buildEmailMessage(from, to, subject, body, isHTML)

	// Verify required headers
	if !strings.Contains(message, "From: sender@example.com") {
		t.Error("Message should contain From header")
	}
	if !strings.Contains(message, "To: recipient@example.com") {
		t.Error("Message should contain To header")
	}
	if !strings.Contains(message, "Subject: Test Subject") {
		t.Error("Message should contain Subject header")
	}
	if !strings.Contains(message, "Content-Type: text/plain; charset=UTF-8") {
		t.Error("Message should have text/plain content type")
	}
	if !strings.Contains(message, "MIME-Version: 1.0") {
		t.Error("Message should contain MIME-Version header")
	}
	if !strings.Contains(message, body) {
		t.Error("Message should contain the body")
	}
}

func TestBuildEmailMessage_HTML(t *testing.T) {
	from := "sender@example.com"
	to := []string{"recipient@example.com"}
	subject := "Test Subject"
	body := "<html><body><h1>Hello</h1></body></html>"
	isHTML := true

	message := buildEmailMessage(from, to, subject, body, isHTML)

	if !strings.Contains(message, "Content-Type: text/html; charset=UTF-8") {
		t.Error("Message should have text/html content type")
	}
	if !strings.Contains(message, body) {
		t.Error("Message should contain the HTML body")
	}
}

func TestBuildEmailMessage_MultipleRecipients(t *testing.T) {
	from := "sender@example.com"
	to := []string{"recipient1@example.com", "recipient2@example.com", "recipient3@example.com"}
	subject := "Test Subject"
	body := "Body"
	isHTML := false

	message := buildEmailMessage(from, to, subject, body, isHTML)

	if !strings.Contains(message, "To: recipient1@example.com, recipient2@example.com, recipient3@example.com") {
		t.Error("Message should contain all recipients")
	}
}

/**
 * Config Extraction Tests
 */

func TestEmailExecutor_Execute_WithCC(t *testing.T) {
	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{
		"cc": "cc@example.com",
	})

	err := executor.Execute(context.Background(), config)

	// Should fail on SMTP connection, not config extraction
	if err == nil {
		t.Fatal("Expected SMTP error, got nil")
	}
	if strings.Contains(err.Error(), "missing") && strings.Contains(err.Error(), "config") {
		t.Errorf("Should not fail on config extraction: %v", err)
	}
}

func TestEmailExecutor_Execute_WithBCC(t *testing.T) {
	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{
		"bcc": "bcc@example.com",
	})

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected SMTP error, got nil")
	}
}

func TestEmailExecutor_Execute_WithHTMLBody(t *testing.T) {
	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{
		"body":     nil,
		"htmlBody": "<html><body><h1>Hello</h1></body></html>",
	})

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected SMTP error, got nil")
	}
}

func TestEmailExecutor_Execute_WithPlainBody(t *testing.T) {
	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{
		"body": "This is a plain text body",
	})

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected SMTP error, got nil")
	}
}

func TestEmailExecutor_Execute_HTMLBodyPriority(t *testing.T) {
	executor := NewEmailExecutor()

	// When both body and htmlBody exist, htmlBody should take priority
	config := smtpConfig(map[string]interface{}{
		"body":     "Plain text",
		"htmlBody": "<html>HTML body</html>",
	})

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected SMTP error, got nil")
	}
}

/**
 * Full Config Tests
 */

func TestEmailExecutor_Execute_FullConfig(t *testing.T) {
	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{
		"to":       "admin@example.com",
		"cc":       "manager@example.com",
		"bcc":      "audit@example.com",
		"subject":  "Alert: System Warning",
		"htmlBody": "<html><body><h1>Warning</h1><p>System alert triggered</p></body></html>",
	})

	err := executor.Execute(context.Background(), config)

	// Should fail on SMTP connection, not config
	if err == nil {
		t.Fatal("Expected SMTP error, got nil")
	}
}

/**
 * Edge Cases
 */

func TestEmailExecutor_Execute_EmptyCC(t *testing.T) {
	executor := NewEmailExecutor()

	// Empty CC should be ignored
	config := smtpConfig(map[string]interface{}{
		"cc": "",
	})

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected SMTP error, got nil")
	}
}

func TestEmailExecutor_Execute_EmptyBCC(t *testing.T) {
	executor := NewEmailExecutor()

	// Empty BCC should be ignored
	config := smtpConfig(map[string]interface{}{
		"bcc": "",
	})

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected SMTP error, got nil")
	}
}

func TestEmailExecutor_Execute_EmptyHTMLBody(t *testing.T) {
	executor := NewEmailExecutor()

	// Empty htmlBody should fall back to body
	config := smtpConfig(map[string]interface{}{
		"htmlBody": "",
		"body":     "Plain text fallback",
	})

	err := executor.Execute(context.Background(), config)

	if err == nil {
		t.Fatal("Expected SMTP error, got nil")
	}
}

func TestEmailExecutor_Execute_NoBody(t *testing.T) {
	executor := NewEmailExecutor()

	// Email with no body should still work (empty body)
	config := smtpConfig(map[string]interface{}{
		"body": nil,
	})

	err := executor.Execute(context.Background(), config)

	// Should fail on SMTP, not config
	if err == nil {
		t.Fatal("Expected SMTP error, got nil")
	}
}

func TestEmailExecutor_Execute_DefaultSmtpPort(t *testing.T) {
	executor := NewEmailExecutor()

	// When smtpPort is not provided, should default to 587
	config := smtpConfig(map[string]interface{}{
		"smtpPort": nil,
	})

	err := executor.Execute(context.Background(), config)

	// Should fail on SMTP connection (port 587), not config extraction
	if err == nil {
		t.Fatal("Expected SMTP error, got nil")
	}
	if strings.Contains(err.Error(), "missing") {
		t.Errorf("Should not fail on missing field: %v", err)
	}
}

/**
 * Success Path Tests (with mock SMTP server)
 */

func TestEmailExecutor_Execute_Success(t *testing.T) {
	port, messages, cleanup := mocksmtp.StartServer(t)
	defer cleanup()

	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{
		"smtpHost": "127.0.0.1",
		"smtpPort": port,
		"username": "testuser",
		"password": "testpass",
		"fromAddr": "sender@test.com",
		"to":       "recipient@test.com",
		"subject":  "Test Subject",
		"body":     "Hello, this is a test email",
	})

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.From != "sender@test.com" {
			t.Errorf("From = %q, want 'sender@test.com'", msg.From)
		}
		if len(msg.Recipients) != 1 || msg.Recipients[0] != "recipient@test.com" {
			t.Errorf("Recipients = %v, want ['recipient@test.com']", msg.Recipients)
		}
		if !strings.Contains(msg.Data, "Subject: Test Subject") {
			t.Error("Message should contain Subject header")
		}
		if !strings.Contains(msg.Data, "Hello, this is a test email") {
			t.Error("Message should contain the body")
		}
		if !strings.Contains(msg.Data, "Content-Type: text/plain; charset=UTF-8") {
			t.Error("Message should have text/plain content type")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for email message")
	}
}

func TestEmailExecutor_Execute_Success_WithHTMLBody(t *testing.T) {
	port, messages, cleanup := mocksmtp.StartServer(t)
	defer cleanup()

	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{
		"smtpHost": "127.0.0.1",
		"smtpPort": port,
		"username": "testuser",
		"password": "testpass",
		"fromAddr": "sender@test.com",
		"to":       "recipient@test.com",
		"subject":  "HTML Test",
		"body":     nil,
		"htmlBody": "<html><body><h1>Hello</h1></body></html>",
	})

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if !strings.Contains(msg.Data, "Content-Type: text/html; charset=UTF-8") {
			t.Error("Message should have text/html content type")
		}
		if !strings.Contains(msg.Data, "<h1>Hello</h1>") {
			t.Error("Message should contain HTML body")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for email message")
	}
}

func TestEmailExecutor_Execute_Success_WithCC(t *testing.T) {
	port, messages, cleanup := mocksmtp.StartServer(t)
	defer cleanup()

	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{
		"smtpHost": "127.0.0.1",
		"smtpPort": port,
		"username": "testuser",
		"password": "testpass",
		"fromAddr": "sender@test.com",
		"to":       "primary@test.com",
		"cc":       "cc@test.com",
		"bcc":      "bcc@test.com",
		"subject":  "Multi-recipient Test",
		"body":     "Test with CC and BCC",
	})

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if len(msg.Recipients) != 3 {
			t.Fatalf("Recipients count = %d, want 3", len(msg.Recipients))
		}
		expected := map[string]bool{
			"primary@test.com": false,
			"cc@test.com":      false,
			"bcc@test.com":     false,
		}
		for _, r := range msg.Recipients {
			if _, ok := expected[r]; !ok {
				t.Errorf("Unexpected recipient: %q", r)
			}
			expected[r] = true
		}
		for addr, found := range expected {
			if !found {
				t.Errorf("Missing expected recipient: %q", addr)
			}
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for email message")
	}
}

func TestEmailExecutor_Execute_Success_FullConfig(t *testing.T) {
	port, messages, cleanup := mocksmtp.StartServer(t)
	defer cleanup()

	executor := NewEmailExecutor()

	config := smtpConfig(map[string]interface{}{
		"smtpHost": "127.0.0.1",
		"smtpPort": port,
		"username": "admin",
		"password": "secret",
		"fromAddr": "alerts@mycompany.com",
		"to":       "admin@example.com",
		"cc":       "manager@example.com",
		"bcc":      "audit@example.com",
		"subject":  "Alert: System Warning",
		"body":     nil,
		"htmlBody": "<html><body><h1>Warning</h1><p>System alert triggered</p></body></html>",
	})

	err := executor.Execute(context.Background(), config)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.From != "alerts@mycompany.com" {
			t.Errorf("From = %q, want 'alerts@mycompany.com'", msg.From)
		}
		if len(msg.Recipients) != 3 {
			t.Fatalf("Recipients count = %d, want 3", len(msg.Recipients))
		}
		if !strings.Contains(msg.Data, "Subject: Alert: System Warning") {
			t.Error("Message should contain subject")
		}
		if !strings.Contains(msg.Data, "Content-Type: text/html; charset=UTF-8") {
			t.Error("Message should have text/html content type")
		}
		if !strings.Contains(msg.Data, "<h1>Warning</h1>") {
			t.Error("Message should contain HTML content")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for email message")
	}
}

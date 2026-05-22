package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"triggers/src/modules/events/application/ports"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewEmailExecutor creates a new email trigger executor adapter.
func NewEmailExecutor() ports.TriggerExecutor {
	return &EmailExecutor{}
}

// Execute sends an email based on the trigger configuration.
//
// All SMTP and email settings are read from the trigger config document,
// following the same pattern as HTTP (endpoint, method) and MQTT (broker, port).
//
// Parameters:
//   - ctx: Context for controlling cancellation and timeouts
//   - config: Trigger configuration (email field) with placeholders already resolved
//
// Returns:
//   - error: If email sending fails
func (e *EmailExecutor) Execute(ctx context.Context, config map[string]interface{}) error {
	// Extract email config (the application already extracted this from trigger.config.email)
	emailConfig, ok := config["email"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("email trigger config missing 'email' field")
	}

	// SMTP server configuration (from trigger document)
	smtpHost, ok := emailConfig["smtpHost"].(string)
	if !ok || smtpHost == "" {
		return fmt.Errorf("email trigger config missing required field 'smtpHost'")
	}

	smtpPort := 587
	if portVal, exists := emailConfig["smtpPort"]; exists {
		switch p := portVal.(type) {
		case float64:
			smtpPort = int(p)
		case int:
			smtpPort = p
		}
	}

	username, _ := emailConfig["username"].(string)
	password, _ := emailConfig["password"].(string)

	fromAddr, ok := emailConfig["fromAddr"].(string)
	if !ok || fromAddr == "" {
		return fmt.Errorf("email trigger config missing required field 'fromAddr'")
	}

	// Extract recipient (to) - required
	to, ok := emailConfig["to"].(string)
	if !ok || to == "" {
		return fmt.Errorf("email trigger config missing required field 'to'")
	}
	recipients := []string{to}

	// CC recipients (optional)
	if cc, exists := emailConfig["cc"]; exists {
		if ccStr, ok := cc.(string); ok && ccStr != "" {
			recipients = append(recipients, ccStr)
		}
	}

	// BCC recipients (optional)
	if bcc, exists := emailConfig["bcc"]; exists {
		if bccStr, ok := bcc.(string); ok && bccStr != "" {
			recipients = append(recipients, bccStr)
		}
	}

	// Extract subject - required
	subject, ok := emailConfig["subject"].(string)
	if !ok || subject == "" {
		return fmt.Errorf("email trigger config missing required field 'subject'")
	}

	// Extract body (prefer htmlBody, fallback to body)
	var body string
	var isHTML bool

	if htmlBody, exists := emailConfig["htmlBody"]; exists {
		if htmlBodyStr, ok := htmlBody.(string); ok && htmlBodyStr != "" {
			body = htmlBodyStr
			isHTML = true
		}
	}

	if body == "" {
		if plainBody, exists := emailConfig["body"]; exists {
			if bodyStr, ok := plainBody.(string); ok {
				body = bodyStr
				isHTML = false
			}
		}
	}

	// Build email message
	message := buildEmailMessage(fromAddr, recipients, subject, body, isHTML)

	// Setup SMTP authentication only when credentials are provided.
	// net/smtp.PlainAuth requires the server to advertise AUTH; passing
	// a non-nil auth against a server that doesn't advertise it fails
	// with "smtp: server doesn't support AUTH" even though anonymous
	// delivery would have succeeded. Treat empty username as "no auth"
	// so the executor works against both authenticated relays and
	// relays open on a private network.
	var auth smtp.Auth
	if username != "" {
		auth = smtp.PlainAuth("", username, password, smtpHost)
	}

	// Send email
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	logger.Info(fmt.Sprintf("[INFRA:EmailExecutor] Sending email to %s via %s", strings.Join(recipients, ", "), addr))

	err := smtp.SendMail(addr, auth, fromAddr, recipients, []byte(message))
	if err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:EmailExecutor] Failed to send email to %s", strings.Join(recipients, ", ")))
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.Info(fmt.Sprintf("[INFRA:EmailExecutor] Email sent successfully to %s", strings.Join(recipients, ", ")))
	return nil
}

// GetType returns the trigger type this executor handles.
func (e *EmailExecutor) GetType() string {
	return "email"
}

// buildEmailMessage constructs an RFC 5322 compliant email message.
func buildEmailMessage(from string, to []string, subject, body string, isHTML bool) string {
	var builder strings.Builder

	// Headers
	builder.WriteString(fmt.Sprintf("From: %s\r\n", from))
	builder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))
	builder.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	builder.WriteString("MIME-Version: 1.0\r\n")

	if isHTML {
		builder.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		builder.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	builder.WriteString("\r\n")

	// Body
	builder.WriteString(body)

	return builder.String()
}

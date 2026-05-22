package email

import (
	"triggers/src/modules/events/application/ports"
)

// EmailExecutor is an infrastructure adapter that implements the TriggerExecutor port.
//
// Following Hexagonal Architecture, this adapter:
// - Lives in the infrastructure layer
// - Implements the application port interface (ports.TriggerExecutor)
// - Contains concrete implementation details (SMTP client, network I/O)
// - Has framework dependencies (net/smtp)
//
// All configuration comes from the trigger document in MongoDB (config.email):
//
//	{
//	  "smtpHost": "smtp.gmail.com",
//	  "smtpPort": 587,
//	  "username": "user@example.com",    // Optional
//	  "password": "app-password",        // Optional
//	  "fromAddr": "alerts@example.com",
//	  "to": "admin@example.com",
//	  "cc": "ops@example.com",           // Optional
//	  "bcc": "audit@example.com",        // Optional
//	  "subject": "Alert: {{payload.alertType}}",
//	  "body": "Plain text body",         // Optional
//	  "htmlBody": "<html>...</html>"     // Optional
//	}
type EmailExecutor struct{}

// Compile-time check to ensure EmailExecutor implements ports.TriggerExecutor interface
var _ ports.TriggerExecutor = (*EmailExecutor)(nil)

package teams

import (
	"net/http"

	"triggers/src/modules/events/application/ports"
)

// TeamsExecutor handles Microsoft Teams webhook trigger execution.
//
// Config schema:
//   {
//     "webhookUrl": "https://outlook.office.com/webhook/...",
//     "title": "Alert: {{payload.alertType}}",
//     "text": "Sensor {{payload.sensor}} triggered: {{payload.message}}",
//     "themeColor": "FF0000",  // Optional: hex color (e.g., red for alerts)
//     "sections": [            // Optional: additional sections
//       {
//         "activityTitle": "Details",
//         "facts": [
//           {"name": "Sensor", "value": "{{payload.sensor}}"},
//           {"name": "Value", "value": "{{payload.value}}"}
//         ]
//       }
//     ]
//   }
//
// Microsoft Teams uses the MessageCard format:
// https://docs.microsoft.com/en-us/outlook/actionable-messages/message-card-reference
type TeamsExecutor struct {
	client *http.Client
}

// Compile-time check to ensure TeamsExecutor implements ports.TriggerExecutor interface
var _ ports.TriggerExecutor = (*TeamsExecutor)(nil)

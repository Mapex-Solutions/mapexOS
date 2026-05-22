package slack

import (
	"net/http"

	"triggers/src/modules/events/application/ports"
)

// SlackExecutor handles Slack webhook trigger execution.
//
// Config schema:
//   {
//     "webhookUrl": "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXX",
//     "text": "Alert: Sensor {{payload.sensor}} triggered",
//     "username": "Alert Bot",      // Optional: custom bot name
//     "icon_emoji": ":warning:",    // Optional: bot emoji icon
//     "channel": "#alerts",          // Optional: override default channel
//     "blocks": [                    // Optional: Slack Block Kit for rich formatting
//       {
//         "type": "section",
//         "text": {
//           "type": "mrkdwn",
//           "text": "*Alert*: {{payload.message}}"
//         }
//       }
//     ]
//   }
//
// Slack webhook format:
// https://api.slack.com/messaging/webhooks
type SlackExecutor struct {
	client *http.Client
}

// Compile-time check to ensure SlackExecutor implements ports.TriggerExecutor interface
var _ ports.TriggerExecutor = (*SlackExecutor)(nil)

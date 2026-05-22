package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/**
 * TECHNICAL TRIGGER CONFIGS
 * Domain layer configurations for technical triggers (HTTP, MQTT, RabbitMQ, NATS, WebSocket)
 * Mirrored from contracts but with bson tags for MongoDB persistence
 */

// HttpConfig defines the configuration for HTTP/HTTPS triggers
type HttpConfig struct {
	Endpoint string            `bson:"endpoint"`
	Method   string            `bson:"method"`
	Headers  map[string]string `bson:"headers,omitempty"`
	Body     map[string]any    `bson:"body,omitempty"`
	Timeout  *int              `bson:"timeout,omitempty"`
}

// MqttConfig defines the configuration for MQTT triggers
type MqttConfig struct {
	Broker   string         `bson:"broker"`
	Port     int            `bson:"port"`
	Topic    string         `bson:"topic"`
	Qos      int            `bson:"qos"`
	Username *string        `bson:"username,omitempty"`
	Password *string        `bson:"password,omitempty"`
	ClientId *string        `bson:"clientId,omitempty"`
	Message  map[string]any `bson:"message,omitempty"`
	UseTLS   *bool          `bson:"useTLS,omitempty"`
}

// RabbitmqConfig defines the configuration for RabbitMQ triggers
type RabbitmqConfig struct {
	Host         string         `bson:"host"`
	Port         int            `bson:"port"`
	Vhost        *string        `bson:"vhost,omitempty"`
	Username     string         `bson:"username"`
	Password     string         `bson:"password"`
	PublishMode  string         `bson:"publishMode"`
	Exchange     *string        `bson:"exchange,omitempty"`
	ExchangeType *string        `bson:"exchangeType,omitempty"`
	RoutingKey   *string        `bson:"routingKey,omitempty"`
	Queue        *string        `bson:"queue,omitempty"`
	Message      map[string]any `bson:"message,omitempty"`
	UseTLS       *bool          `bson:"useTLS,omitempty"`
}

// NatsConfig defines the configuration for NATS triggers
type NatsConfig struct {
	Server   string         `bson:"server"`
	Subject  string         `bson:"subject"`
	Username *string        `bson:"username,omitempty"`
	Password *string        `bson:"password,omitempty"`
	Token    *string        `bson:"token,omitempty"`
	Message  map[string]any `bson:"message,omitempty"`
	UseTLS   *bool          `bson:"useTLS,omitempty"`
}

// WebsocketConfig defines the configuration for WebSocket triggers
type WebsocketConfig struct {
	Url     string            `bson:"url"`
	Message map[string]any    `bson:"message,omitempty"`
	Headers map[string]string `bson:"headers,omitempty"`
}

/**
 * COMMUNICATION TRIGGER CONFIGS
 * Domain layer configurations for communication triggers (Email, Teams, Slack)
 */

// EmailConfig defines the configuration for Email triggers
type EmailConfig struct {
	// SMTP server configuration
	SmtpHost string  `bson:"smtpHost"`
	SmtpPort int     `bson:"smtpPort"`
	Username *string `bson:"username,omitempty"`
	Password *string `bson:"password,omitempty"`
	FromAddr string  `bson:"fromAddr"`

	// Email content
	To       string  `bson:"to"`
	Cc       *string `bson:"cc,omitempty"`
	Bcc      *string `bson:"bcc,omitempty"`
	Subject  string  `bson:"subject"`
	Body     *string `bson:"body,omitempty"`
	HtmlBody *string `bson:"htmlBody,omitempty"`
}

// TeamsConfig defines the configuration for Microsoft Teams triggers
type TeamsConfig struct {
	WebhookUrl string  `bson:"webhookUrl"`
	Title      string  `bson:"title"`
	Text       string  `bson:"text"`
	ThemeColor *string `bson:"themeColor,omitempty"`
}

// SlackConfig defines the configuration for Slack triggers
type SlackConfig struct {
	WebhookUrl string  `bson:"webhookUrl"`
	Channel    *string `bson:"channel,omitempty"`
	Username   *string `bson:"username,omitempty"`
	IconEmoji  *string `bson:"iconEmoji,omitempty"`
	Message    string  `bson:"message"`
}

/**
 * TRIGGER CONFIG (UNION TYPE)
 * Union type for all trigger configurations
 */

// TriggerConfig is a union type that holds the configuration for a specific trigger type.
// Only one config should be populated based on the TriggerType field.
type TriggerConfig struct {
	// Technical triggers
	Http      *HttpConfig      `bson:"http,omitempty"`
	Mqtt      *MqttConfig      `bson:"mqtt,omitempty"`
	Rabbitmq  *RabbitmqConfig  `bson:"rabbitmq,omitempty"`
	Nats      *NatsConfig      `bson:"nats,omitempty"`
	Websocket *WebsocketConfig `bson:"websocket,omitempty"`

	// Communication triggers
	Email *EmailConfig `bson:"email,omitempty"`
	Teams *TeamsConfig `bson:"teams,omitempty"`
	Slack *SlackConfig `bson:"slack,omitempty"`
}

/**
 * TRIGGER ENTITY
 * Main domain entity for triggers
 */

// Trigger represents a trigger entity in the system
// Triggers can be Technical (http, mqtt, rabbitmq, nats, websocket)
// or Communication (email, teams, slack)
//
// Template Resources Pattern:
// - isSystem=true: Global MAPEX templates (no orgId, no pathKey)
// - isTemplate=true: Vendor/Customer templates (inherited by descendants)
// - isSystem=false, isTemplate=false: Local triggers (org-specific)
type Trigger struct {
	ID model.ObjectId `bson:"_id,omitempty"`

	Name        string  `bson:"name"`
	Description *string `bson:"description,omitempty"`
	TriggerType string  `bson:"triggerType"`
	Category    string  `bson:"category"`
	Enabled     bool    `bson:"enabled"`

	// Config contains the trigger-specific configuration (defined locally)
	Config TriggerConfig `bson:"config"`

	// Template visibility flags
	IsSystem   bool `bson:"isSystem"`
	IsTemplate bool `bson:"isTemplate"`

	// Multi-tenant fields (populated by service from RequestContext)
	// For isSystem=true: orgId is nil, pathKey is empty
	// For isSystem=false: orgId and pathKey are populated
	OrgID   *model.ObjectId `bson:"orgId,omitempty"`
	PathKey string          `bson:"pathKey"`

	// Metadata (follows platform standard naming)
	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

// GetCreated implements the model interface for created timestamp
func (t *Trigger) GetCreated() time.Time {
	return t.Created
}

// GetUpdated implements the model interface for updated timestamp
func (t *Trigger) GetUpdated() time.Time {
	return t.Updated
}


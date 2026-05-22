package triggers

import (
	"fmt"

	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/* TECHNICAL TRIGGER CONFIGS */

// HttpConfig defines the configuration for HTTP/HTTPS triggers
type HttpConfig struct {
	Endpoint string            `json:"endpoint" validate:"required,url"`
	Method   string            `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Headers  map[string]string `json:"headers,omitempty"`
	Body     map[string]any    `json:"body,omitempty"`
	Timeout  *int              `json:"timeout,omitempty" validate:"omitempty,min=1000,max=300000"` // 1s to 5min in milliseconds
}

// MqttConfig defines the configuration for MQTT triggers
type MqttConfig struct {
	Broker   string         `json:"broker" validate:"required"`
	Port     int            `json:"port" validate:"required,min=1,max=65535"`
	Topic    string         `json:"topic" validate:"required,min=1"`
	Qos      int            `json:"qos" validate:"required,oneof=0 1 2"`
	Username *string        `json:"username,omitempty"`
	Password *string        `json:"password,omitempty"`
	ClientId *string        `json:"clientId,omitempty"`
	Message  map[string]any `json:"message,omitempty"`
	UseTLS   *bool          `json:"useTLS,omitempty"`
}

// RabbitmqConfig defines the configuration for RabbitMQ triggers
type RabbitmqConfig struct {
	Host         string         `json:"host" validate:"required"`
	Port         int            `json:"port" validate:"required,min=1,max=65535"`
	Vhost        *string        `json:"vhost,omitempty"`
	Username     string         `json:"username" validate:"required"`
	Password     string         `json:"password" validate:"required"`
	PublishMode  string         `json:"publishMode" validate:"required,oneof=exchange queue"`
	Exchange     *string        `json:"exchange,omitempty"`                                  // Required if publishMode=exchange
	ExchangeType *string        `json:"exchangeType,omitempty" validate:"omitempty,oneof=direct fanout topic headers"` // Required if publishMode=exchange
	RoutingKey   *string        `json:"routingKey,omitempty"`                                // Optional for fanout exchange
	Queue        *string        `json:"queue,omitempty"`                                     // Required if publishMode=queue
	Message      map[string]any `json:"message,omitempty"`
	UseTLS       *bool          `json:"useTLS,omitempty"`
}

// NatsConfig defines the configuration for NATS triggers
type NatsConfig struct {
	Server   string         `json:"server" validate:"required"`
	Subject  string         `json:"subject" validate:"required,min=1"`
	Username *string        `json:"username,omitempty"`
	Password *string        `json:"password,omitempty"`
	Token    *string        `json:"token,omitempty"`
	Message  map[string]any `json:"message,omitempty"`
	UseTLS   *bool          `json:"useTLS,omitempty"`
}

// WebsocketConfig defines the configuration for WebSocket triggers
type WebsocketConfig struct {
	Url     string            `json:"url" validate:"required,url"`
	Message map[string]any    `json:"message,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

/* COMMUNICATION TRIGGER CONFIGS */

// EmailConfig defines the configuration for Email triggers
type EmailConfig struct {
	// SMTP server configuration
	SmtpHost string  `json:"smtpHost" validate:"required,min=1"`
	SmtpPort int     `json:"smtpPort" validate:"required,min=1,max=65535"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	FromAddr string  `json:"fromAddr" validate:"required,min=3"`

	// Email content
	To       string  `json:"to" validate:"required,min=3"`
	Cc       *string `json:"cc,omitempty"`
	Bcc      *string `json:"bcc,omitempty"`
	Subject  string  `json:"subject" validate:"required,min=1"`
	Body     *string `json:"body,omitempty"`     // Plain text body
	HtmlBody *string `json:"htmlBody,omitempty"` // HTML body (takes precedence over Body)
}

// TeamsConfig defines the configuration for Microsoft Teams triggers
type TeamsConfig struct {
	WebhookUrl string  `json:"webhookUrl" validate:"required,url"`
	Title      string  `json:"title" validate:"required,min=1"`
	Text       string  `json:"text" validate:"required,min=1"`
	ThemeColor *string `json:"themeColor,omitempty" validate:"omitempty,hexcolor"` // Hex color without #
}

// SlackConfig defines the configuration for Slack triggers
type SlackConfig struct {
	WebhookUrl string  `json:"webhookUrl" validate:"required,url"`
	Channel    *string `json:"channel,omitempty"`                        // Override default channel
	Username   *string `json:"username,omitempty"`                       // Bot display name
	IconEmoji  *string `json:"iconEmoji,omitempty" validate:"omitempty"` // e.g., :robot_face:
	Message    string  `json:"message" validate:"required,min=1"`
}

/* TRIGGER CONFIG (UNION TYPE) */

// TriggerConfig is a union type that holds the configuration for a specific trigger type.
// Similar to ProtocolType in assets and Router in routegroups.
// Only one config should be populated based on the TriggerType field.
type TriggerConfig struct {
	// Technical triggers
	Http      *HttpConfig      `json:"http,omitempty"`
	Mqtt      *MqttConfig      `json:"mqtt,omitempty"`
	Rabbitmq  *RabbitmqConfig  `json:"rabbitmq,omitempty"`
	Nats      *NatsConfig      `json:"nats,omitempty"`
	Websocket *WebsocketConfig `json:"websocket,omitempty"`

	// Communication triggers
	Email *EmailConfig `json:"email,omitempty"`
	Teams *TeamsConfig `json:"teams,omitempty"`
	Slack *SlackConfig `json:"slack,omitempty"`
}

/* TRIGGER DTOs */

// TriggerCreate is the DTO for creating a new trigger
// ⚠️ IMPORTANT: orgId and pathKey can be provided in DTO but will be validated/overwritten by service from RequestContext
type TriggerCreate struct {
	Name        string        `json:"name" validate:"required,min=3,max=150"`
	Description *string       `json:"description,omitempty" validate:"omitempty,max=500"`
	TriggerType string        `json:"triggerType" validate:"required,oneof=http mqtt rabbitmq nats websocket email teams slack"`
	Category    string        `json:"category" validate:"required,oneof=technical communication"`
	Enabled     bool          `json:"enabled" validate:"required"`
	Config      TriggerConfig `json:"config" validate:"required"`

	// Template visibility flags
	IsSystem   bool `json:"isSystem"`
	IsTemplate bool `json:"isTemplate"`

	// Multi-tenant hierarchical fields (populated automatically by service from RequestContext)
	OrgID   *model.ObjectId `json:"orgId,omitempty" validate:"omitempty"`
	PathKey *string         `json:"pathKey,omitempty" validate:"omitempty"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

// TriggerUpdate is the DTO for updating a trigger
// All fields are optional (pointers) except Config which can be fully replaced
type TriggerUpdate struct {
	Name        *string        `json:"name,omitempty" validate:"omitempty,min=3,max=150"`
	Description *string        `json:"description,omitempty" validate:"omitempty,max=500"`
	TriggerType *string        `json:"triggerType,omitempty" validate:"omitempty,oneof=http mqtt rabbitmq nats websocket email teams slack"`
	Category    *string        `json:"category,omitempty" validate:"omitempty,oneof=technical communication"`
	Enabled     *bool          `json:"enabled,omitempty"`
	Config      *TriggerConfig `json:"config,omitempty"`

	// Template visibility flags (optional for updates)
	IsSystem   *bool `json:"isSystem,omitempty"`
	IsTemplate *bool `json:"isTemplate,omitempty"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

// TriggerResponse is the DTO for trigger responses
type TriggerResponse struct {
	ID          *common.ObjectID `json:"id,omitempty"`
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	TriggerType *string          `json:"triggerType,omitempty"`
	Category    *string          `json:"category,omitempty"`
	Enabled     *bool            `json:"enabled,omitempty"`
	Config      *TriggerConfig   `json:"config,omitempty"`

	// Template visibility flags
	IsSystem   *bool `json:"isSystem,omitempty"`
	IsTemplate *bool `json:"isTemplate,omitempty"`

	OrgID   *model.ObjectId  `json:"orgId,omitempty"`
	PathKey *string          `json:"pathKey,omitempty"`
	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

func (r *TriggerResponse) SetCreated(t *common.NullTime) { r.Created = t }
func (r *TriggerResponse) SetUpdated(t *common.NullTime) { r.Updated = t }

// TriggerQuery is the DTO for querying triggers with filters
type TriggerQuery struct {
	query.BaseQueryDTO // Embed BaseQueryDTO for includeChildren support

	ID          *string `query:"id"`
	Name        *string `query:"name"`
	TriggerType *string `query:"triggerType"`
	Category    *string `query:"category"`
	Enabled     *bool   `query:"enabled"`
	OrgID       *string `query:"orgId"`
	PathKey     *string `query:"pathKey"`

	// Template filters (for querying system/template triggers)
	IsSystem   *bool `query:"isSystem"`
	IsTemplate *bool `query:"isTemplate"`

	// Pagination (Note: IncludeChildren comes from BaseQueryDTO)
	Page    *int    `query:"page"`
	PerPage *int    `query:"perPage"`
	Sort    *string `query:"sort"`
}

// TriggerListResponse is the DTO for paginated trigger list responses
type TriggerListResponse struct {
	Data       []TriggerResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PerPage    int               `json:"perPage"`
	TotalPages int               `json:"totalPages"`
}

// TriggerExecuteEvent is the DTO for trigger execution events published to NATS.
// Subject: trigger.{triggerId}.execute
// Published by: Router (kind=trigger).
// The consumer receives this as a generic enriched event from the Router.
// Fields like OrgID arrive as string — consumer converts to ObjectId if needed.
type TriggerExecuteEvent struct {
	TriggerID      string                 `json:"triggerId"`
	ExecutionID    string                 `json:"executionId"`
	EventTrackerId string                 `json:"eventTrackerId"`
	Source         string                 `json:"source"`
	Payload        map[string]interface{} `json:"payload"`
	OrgID          string                 `json:"orgId"`
	PathKey        string                 `json:"pathKey"`
	Created        string                 `json:"created"`
}

/* TRANSFORMATIONS (Custom Validation) */

// Transform validates that the TriggerConfig matches the TriggerType.
// This ensures that if triggerType is "http", then config.http must be populated.
// Similar to Router.Transform() in routegroups contract.
func (c *TriggerConfig) Transform(triggerType string) error {
	// Count how many config fields are populated
	populatedCount := 0
	var populatedField string

	if c.Http != nil {
		populatedCount++
		populatedField = "http"
	}
	if c.Mqtt != nil {
		populatedCount++
		populatedField = "mqtt"
	}
	if c.Rabbitmq != nil {
		populatedCount++
		populatedField = "rabbitmq"
	}
	if c.Nats != nil {
		populatedCount++
		populatedField = "nats"
	}
	if c.Websocket != nil {
		populatedCount++
		populatedField = "websocket"
	}
	if c.Email != nil {
		populatedCount++
		populatedField = "email"
	}
	if c.Teams != nil {
		populatedCount++
		populatedField = "teams"
	}
	if c.Slack != nil {
		populatedCount++
		populatedField = "slack"
	}

	// Validate that exactly ONE config field is populated
	if populatedCount == 0 {
		return fmt.Errorf("config must have one field populated matching triggerType '%s'", triggerType)
	}
	if populatedCount > 1 {
		return fmt.Errorf("config must have only one field populated, found %d fields", populatedCount)
	}

	// Validate that the populated field matches the triggerType
	if populatedField != triggerType {
		return fmt.Errorf("config field '%s' does not match triggerType '%s'", populatedField, triggerType)
	}

	// Perform type-specific validations
	switch triggerType {
	case "http":
		return c.validateHttpConfig()
	case "mqtt":
		return c.validateMqttConfig()
	case "rabbitmq":
		return c.validateRabbitmqConfig()
	case "nats":
		return c.validateNatsConfig()
	case "websocket":
		return c.validateWebsocketConfig()
	case "email":
		return c.validateEmailConfig()
	case "teams":
		return c.validateTeamsConfig()
	case "slack":
		return c.validateSlackConfig()
	default:
		return fmt.Errorf("unknown trigger type: %s", triggerType)
	}
}

// validateHttpConfig performs additional validation for HTTP config
func (c *TriggerConfig) validateHttpConfig() error {
	if c.Http == nil {
		return fmt.Errorf("http config is required when triggerType is 'http'")
	}
	// Additional validations can be added here
	// e.g., validate that endpoint starts with http:// or https://
	return nil
}

// validateMqttConfig performs additional validation for MQTT config
func (c *TriggerConfig) validateMqttConfig() error {
	if c.Mqtt == nil {
		return fmt.Errorf("mqtt config is required when triggerType is 'mqtt'")
	}
	// Additional validations can be added here
	return nil
}

// validateRabbitmqConfig performs additional validation for RabbitMQ config
func (c *TriggerConfig) validateRabbitmqConfig() error {
	if c.Rabbitmq == nil {
		return fmt.Errorf("rabbitmq config is required when triggerType is 'rabbitmq'")
	}

	// Validate publish mode specific requirements
	switch c.Rabbitmq.PublishMode {
	case "exchange":
		if c.Rabbitmq.Exchange == nil || *c.Rabbitmq.Exchange == "" {
			return fmt.Errorf("exchange is required when publishMode is 'exchange'")
		}
		if c.Rabbitmq.ExchangeType == nil || *c.Rabbitmq.ExchangeType == "" {
			return fmt.Errorf("exchangeType is required when publishMode is 'exchange'")
		}
	case "queue":
		if c.Rabbitmq.Queue == nil || *c.Rabbitmq.Queue == "" {
			return fmt.Errorf("queue is required when publishMode is 'queue'")
		}
	default:
		return fmt.Errorf("invalid publishMode: %s (must be 'exchange' or 'queue')", c.Rabbitmq.PublishMode)
	}

	return nil
}

// validateNatsConfig performs additional validation for NATS config
func (c *TriggerConfig) validateNatsConfig() error {
	if c.Nats == nil {
		return fmt.Errorf("nats config is required when triggerType is 'nats'")
	}
	// Additional validations can be added here
	return nil
}

// validateWebsocketConfig performs additional validation for WebSocket config
func (c *TriggerConfig) validateWebsocketConfig() error {
	if c.Websocket == nil {
		return fmt.Errorf("websocket config is required when triggerType is 'websocket'")
	}
	// Additional validations can be added here
	return nil
}

// validateEmailConfig performs additional validation for Email config
func (c *TriggerConfig) validateEmailConfig() error {
	if c.Email == nil {
		return fmt.Errorf("email config is required when triggerType is 'email'")
	}
	// Additional validations can be added here
	return nil
}

// validateTeamsConfig performs additional validation for Teams config
func (c *TriggerConfig) validateTeamsConfig() error {
	if c.Teams == nil {
		return fmt.Errorf("teams config is required when triggerType is 'teams'")
	}
	// Additional validations can be added here
	return nil
}

// validateSlackConfig performs additional validation for Slack config
func (c *TriggerConfig) validateSlackConfig() error {
	if c.Slack == nil {
		return fmt.Errorf("slack config is required when triggerType is 'slack'")
	}
	// Additional validations can be added here
	return nil
}

// Transform validates that the TriggerCreate DTO is correctly structured.
// Calls TriggerConfig.Transform() to ensure config matches triggerType.
func (t *TriggerCreate) Transform() error {
	return t.Config.Transform(t.TriggerType)
}

// Transform validates that the TriggerUpdate DTO is correctly structured.
// Only validates if both TriggerType and Config are being updated.
func (t *TriggerUpdate) Transform() error {
	// If both triggerType and config are provided, validate they match
	if t.TriggerType != nil && t.Config != nil {
		return t.Config.Transform(*t.TriggerType)
	}
	// If only config is provided, we cannot validate without knowing the triggerType
	// This validation will be done at service layer where we have the existing trigger
	return nil
}

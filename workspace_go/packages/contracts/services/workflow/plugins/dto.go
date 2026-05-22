// Package plugins holds the cross-service / wire-format contracts for the
// workflow service plugins module.
//
// The PluginManifest aggregate is the JSON authority — HTTP handlers, the
// JetStream FANOUT cache invalidation flow, and the in-process TieredCache
// (L0 RAM / L1 Disk) all serialize through these structs. The bson tags
// alongside json tags reflect the current persistence model where the
// MongoDB document shape mirrors the wire shape one-to-one. A future
// refactor can split persistence (services-side bson-only entity) from
// wire (json-only DTO here) with an explicit mapper; until then this
// single struct serves both edges. Domain entities under
// services/workflow/src/modules/plugins/domain/entities/ alias these
// types so the entity files contain no `json:"..."` literals.
package plugins

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/*
 * HANDLE DEFINITIONS
 */

// HandleDef defines an input or output connection point on a node.
type HandleDef struct {
	ID             string `bson:"id" json:"id"`
	Label          string `bson:"label" json:"label"`
	Position       string `bson:"position" json:"position"`
	DataType       string `bson:"dataType,omitempty" json:"dataType,omitempty"`
	MaxConnections *int   `bson:"maxConnections,omitempty" json:"maxConnections,omitempty"`
	Color          string `bson:"color,omitempty" json:"color,omitempty"`
}

/*
 * AVAILABLE OUTPUTS
 */

// AvailableOutput describes a field available in the node's output for downstream use.
type AvailableOutput struct {
	Path        string `bson:"path" json:"path"`
	Description string `bson:"description" json:"description"`
}

/*
 * UNIFIED ACTION CONTRACT
 */

// HttpActionDef defines an HTTP request template.
type HttpActionDef struct {
	Method  string            `bson:"method" json:"method"`
	Path    string            `bson:"path" json:"path"`
	Headers map[string]string `bson:"headers,omitempty" json:"headers,omitempty"`
	Body    interface{}       `bson:"body,omitempty" json:"body,omitempty"`
	Timeout *int              `bson:"timeout,omitempty" json:"timeout,omitempty"`
}

// MqttActionDef defines an MQTT publish template.
type MqttActionDef struct {
	Topic   string      `bson:"topic" json:"topic"`
	QoS     *int        `bson:"qos,omitempty" json:"qos,omitempty"`
	Payload interface{} `bson:"payload,omitempty" json:"payload,omitempty"`
	Retain  *bool       `bson:"retain,omitempty" json:"retain,omitempty"`
}

// NatsActionDef defines a NATS publish template.
type NatsActionDef struct {
	Subject string      `bson:"subject" json:"subject"`
	Data    interface{} `bson:"data,omitempty" json:"data,omitempty"`
}

// ScriptActionDef defines a JavaScript execution action.
type ScriptActionDef struct {
	Code    string `bson:"code" json:"code"`
	Timeout *int   `bson:"timeout,omitempty" json:"timeout,omitempty"`
}

// ActionOutputDef defines how to extract/transform the response of an action.
type ActionOutputDef struct {
	DataPath  string `bson:"dataPath,omitempty" json:"dataPath,omitempty"`
	ValuePath string `bson:"valuePath,omitempty" json:"valuePath,omitempty"`
	LabelPath string `bson:"labelPath,omitempty" json:"labelPath,omitempty"`
	Transform string `bson:"transform,omitempty" json:"transform,omitempty"`
}

// ActionDef is the unified action contract.
type ActionDef struct {
	Type   string           `bson:"type" json:"type"`
	Http   *HttpActionDef   `bson:"http,omitempty" json:"http,omitempty"`
	Mqtt   *MqttActionDef   `bson:"mqtt,omitempty" json:"mqtt,omitempty"`
	Nats   *NatsActionDef   `bson:"nats,omitempty" json:"nats,omitempty"`
	Script *ScriptActionDef `bson:"script,omitempty" json:"script,omitempty"`
	Output *ActionOutputDef `bson:"output,omitempty" json:"output,omitempty"`
}

/*
 * FETCH OPTIONS
 */

// FetchOptionsPagination defines how to paginate fetchOptions results.
type FetchOptionsPagination struct {
	Mode           string `bson:"mode" json:"mode"`
	CursorParam    string `bson:"cursorParam,omitempty" json:"cursorParam,omitempty"`
	NextCursorPath string `bson:"nextCursorPath,omitempty" json:"nextCursorPath,omitempty"`
	PageParam      string `bson:"pageParam,omitempty" json:"pageParam,omitempty"`
	LimitParam     string `bson:"limitParam,omitempty" json:"limitParam,omitempty"`
	LimitDefault   *int   `bson:"limitDefault,omitempty" json:"limitDefault,omitempty"`
	TotalPath      string `bson:"totalPath,omitempty" json:"totalPath,omitempty"`
}

// FetchOptionsSearch defines server-side search support for fetchOptions.
type FetchOptionsSearch struct {
	Param     string `bson:"param" json:"param"`
	MinLength *int   `bson:"minLength,omitempty" json:"minLength,omitempty"`
}

// FetchOptionsDef configures a dynamic options loader at the manifest level.
type FetchOptionsDef struct {
	ActionDef  `bson:",inline" json:",inline"`
	Pagination *FetchOptionsPagination `bson:"pagination,omitempty" json:"pagination,omitempty"`
	Search     *FetchOptionsSearch     `bson:"search,omitempty" json:"search,omitempty"`
}

/*
 * NODE PROPERTY SYSTEM
 */

// PropertyRendering controls how a form field looks in the UI.
type PropertyRendering struct {
	Multiline      *bool    `bson:"multiline,omitempty" json:"multiline,omitempty"`
	Rows           *int     `bson:"rows,omitempty" json:"rows,omitempty"`
	Password       *bool    `bson:"password,omitempty" json:"password,omitempty"`
	Editor         string   `bson:"editor,omitempty" json:"editor,omitempty"`
	MultipleValues *bool    `bson:"multipleValues,omitempty" json:"multipleValues,omitempty"`
	Min            *float64 `bson:"min,omitempty" json:"min,omitempty"`
	Max            *float64 `bson:"max,omitempty" json:"max,omitempty"`
	Placeholder    string   `bson:"placeholder,omitempty" json:"placeholder,omitempty"`
	DateOnly       *bool    `bson:"dateOnly,omitempty" json:"dateOnly,omitempty"`
}

// FetchOptionsRule defines which fetchOptions entry to use based on form state.
type FetchOptionsRule struct {
	When  map[string][]interface{} `bson:"when" json:"when"`
	Key   string                   `bson:"key" json:"key"`
	Label string                   `bson:"label" json:"label"`
}

// PropertyFetchOptions configures dynamic dropdown fetching for a fieldSource property.
type PropertyFetchOptions struct {
	Rules     []FetchOptionsRule `bson:"rules" json:"rules"`
	DependsOn []string           `bson:"dependsOn,omitempty" json:"dependsOn,omitempty"`
}

// PropertyOption represents a single option in a dropdown list.
type PropertyOption struct {
	Label string      `bson:"label" json:"label"`
	Value interface{} `bson:"value" json:"value"`
}

// DisplayOptions controls conditional visibility of a property based on other fields.
type DisplayOptions struct {
	Show map[string][]interface{} `bson:"show,omitempty" json:"show,omitempty"`
}

// NodePropertyDef defines a single declarative form field for auto-generated config forms.
type NodePropertyDef struct {
	Name           string                `bson:"name" json:"name"`
	DisplayName    string                `bson:"displayName" json:"displayName"`
	Type           string                `bson:"type" json:"type"`
	Default        interface{}           `bson:"default" json:"default"`
	Hint           string                `bson:"hint,omitempty" json:"hint,omitempty"`
	Required       *bool                 `bson:"required,omitempty" json:"required,omitempty"`
	IsSecret       *bool                 `bson:"isSecret,omitempty" json:"isSecret,omitempty"`
	Options        []PropertyOption      `bson:"options,omitempty" json:"options,omitempty"`
	DisplayOptions *DisplayOptions       `bson:"displayOptions,omitempty" json:"displayOptions,omitempty"`
	AllowedSources []string              `bson:"allowedSources,omitempty" json:"allowedSources,omitempty"`
	Rendering      *PropertyRendering    `bson:"rendering,omitempty" json:"rendering,omitempty"`
	FetchOptions   *PropertyFetchOptions `bson:"fetchOptions,omitempty" json:"fetchOptions,omitempty"`
	Values         []NodePropertyDef     `bson:"values,omitempty" json:"values,omitempty"`
	NoticeType     string                `bson:"noticeType,omitempty" json:"noticeType,omitempty"`
}

/*
 * NODE HOOKS
 */

// NodeHooks defines lifecycle hooks for a node type.
type NodeHooks struct {
	Before  *ActionDef `bson:"before,omitempty" json:"before,omitempty"`
	After   *ActionDef `bson:"after,omitempty" json:"after,omitempty"`
	Destroy *ActionDef `bson:"destroy,omitempty" json:"destroy,omitempty"`
}

/*
 * NODE TYPE MANIFEST
 */

// NodeTimeoutDef specifies the default timeout for an async node type.
type NodeTimeoutDef struct {
	Duration     int    `bson:"duration" json:"duration"`
	Unit         string `bson:"unit" json:"unit"`
	EnableOutput bool   `bson:"enableOutput,omitempty" json:"enableOutput,omitempty"`
}

// NodeTypeManifest defines a single node type within a plugin manifest.
type NodeTypeManifest struct {
	Type             string                 `bson:"type" json:"type"`
	Label            string                 `bson:"label" json:"label"`
	Icon             string                 `bson:"icon" json:"icon"`
	Color            string                 `bson:"color" json:"color"`
	Description      string                 `bson:"description" json:"description"`
	Inputs           []HandleDef            `bson:"inputs" json:"inputs"`
	Outputs          []HandleDef            `bson:"outputs" json:"outputs"`
	ConfigSchema     map[string]interface{} `bson:"configSchema,omitempty" json:"configSchema,omitempty"`
	Properties       []NodePropertyDef      `bson:"properties,omitempty" json:"properties,omitempty"`
	Defaults         map[string]interface{} `bson:"defaults,omitempty" json:"defaults,omitempty"`
	Timeout          *NodeTimeoutDef        `bson:"timeout,omitempty" json:"timeout,omitempty"`
	AvailableOutputs []AvailableOutput      `bson:"availableOutputs,omitempty" json:"availableOutputs,omitempty"`
	Operations       map[string]ActionDef   `bson:"operations,omitempty" json:"operations,omitempty"`
	Hooks            *NodeHooks             `bson:"hooks,omitempty" json:"hooks,omitempty"`
}

/*
 * CREDENTIAL DEFINITIONS
 */

// CredentialFieldDef defines a single credential input field.
type CredentialFieldDef struct {
	Name        string           `bson:"name" json:"name"`
	DisplayName string           `bson:"displayName" json:"displayName"`
	Type        string           `bson:"type" json:"type"`
	Required    *bool            `bson:"required,omitempty" json:"required,omitempty"`
	IsSecret    *bool            `bson:"isSecret,omitempty" json:"isSecret,omitempty"`
	Hint        string           `bson:"hint,omitempty" json:"hint,omitempty"`
	Default     interface{}      `bson:"default,omitempty" json:"default,omitempty"`
	Options     []PropertyOption `bson:"options,omitempty" json:"options,omitempty"`
}

// CredentialDef describes one authentication method for a plugin.
type CredentialDef struct {
	ID     string               `bson:"id" json:"id"`
	Name   string               `bson:"name" json:"name"`
	Fields []CredentialFieldDef `bson:"fields" json:"fields"`
	Test   *ActionDef           `bson:"test,omitempty" json:"test,omitempty"`
}

/*
 * PLUGIN METADATA
 */

// PluginMetadata contains visual metadata for marketplace display.
type PluginMetadata struct {
	BrandIcon string `bson:"brandIcon,omitempty" json:"brandIcon,omitempty"`
	Color     string `bson:"color,omitempty" json:"color,omitempty"`
	Docs      string `bson:"docs,omitempty" json:"docs,omitempty"`
}

/*
 * PLUGIN DEFAULTS
 */

// PluginDefaults provides default values inherited by operations and fetchOptions.
type PluginDefaults struct {
	BaseUrl string `bson:"baseUrl,omitempty" json:"baseUrl,omitempty"`
	Timeout *int   `bson:"timeout,omitempty" json:"timeout,omitempty"`
}

/*
 * PLUGIN MANIFEST (Root Aggregate)
 */

// PluginManifest is the root aggregate for an integration plugin.
type PluginManifest struct {
	ID          model.ObjectId `bson:"_id,omitempty" json:"id"`
	PluginID    string         `bson:"pluginId" json:"pluginId"`
	Name        string         `bson:"name" json:"name"`
	Version     string         `bson:"version" json:"version"`
	Category    string         `bson:"category" json:"category"`
	Icon        string         `bson:"icon" json:"icon"`
	Color       string         `bson:"color" json:"color"`
	Description string         `bson:"description" json:"description"`

	Defaults PluginDefaults `bson:"defaults,omitempty" json:"defaults,omitempty"`

	Credentials []CredentialDef `bson:"credentials,omitempty" json:"credentials,omitempty"`

	FetchOptions map[string]FetchOptionsDef `bson:"fetchOptions,omitempty" json:"fetchOptions,omitempty"`

	NodeTypes []NodeTypeManifest `bson:"nodeTypes" json:"nodeTypes"`

	Metadata *PluginMetadata `bson:"metadata,omitempty" json:"metadata,omitempty"`

	Author string   `bson:"author,omitempty" json:"author,omitempty"`
	Tags   []string `bson:"tags,omitempty" json:"tags,omitempty"`

	IsTemplate bool            `bson:"isTemplate" json:"isTemplate"`
	OrgId      *model.ObjectId `bson:"orgId,omitempty" json:"orgId,omitempty"`
	PathKey    string          `bson:"pathKey" json:"pathKey"`

	Enabled bool `bson:"enabled" json:"enabled"`

	Created time.Time `bson:"created" json:"created"`
	Updated time.Time `bson:"updated" json:"updated"`
}

// GetCreated returns the creation timestamp.
func (p *PluginManifest) GetCreated() time.Time { return p.Created }

// GetUpdated returns the last update timestamp.
func (p *PluginManifest) GetUpdated() time.Time { return p.Updated }

/*
 * PLUGIN MANIFEST UPDATE (Partial update — bson-only on the Mongo edge,
 * field names match wire JSON case-insensitively for body decoding).
 */

// PluginManifestUpdate is used for PATCH operations.
type PluginManifestUpdate struct {
	Name        *string         `bson:"name,omitempty"`
	Version     *string         `bson:"version,omitempty"`
	Category    *string         `bson:"category,omitempty"`
	Icon        *string         `bson:"icon,omitempty"`
	Color       *string         `bson:"color,omitempty"`
	Description *string         `bson:"description,omitempty"`
	Enabled     *bool           `bson:"enabled,omitempty"`
	OrgId       *model.ObjectId `bson:"orgId,omitempty"`
	PathKey     *string         `bson:"pathKey,omitempty"`

	Defaults     *PluginDefaults             `bson:"defaults,omitempty"`
	Credentials  *[]CredentialDef            `bson:"credentials,omitempty"`
	FetchOptions *map[string]FetchOptionsDef `bson:"fetchOptions,omitempty"`
	NodeTypes    *[]NodeTypeManifest         `bson:"nodeTypes,omitempty"`
	Metadata     *PluginMetadata             `bson:"metadata,omitempty"`
	Author       *string                     `bson:"author,omitempty"`
	Tags         *[]string                   `bson:"tags,omitempty"`

	Updated time.Time `bson:"updated"`
}

// GetUpdated returns the update timestamp.
func (p *PluginManifestUpdate) GetUpdated() time.Time { return p.Updated }

/*
 * REQUEST-LEVEL DTOs
 */

// PluginId is the params DTO for plugin id extraction used by
// ValidationMiddleware to parse and validate the :id URL param.
type PluginId struct {
	PluginId string `json:"id" params:"id" validate:"required"`
}

// PluginQuery defines query parameters for listing plugins.
type PluginQuery struct {
	Name       *string   `json:"name,omitempty" query:"name"`
	Category   *string   `json:"category,omitempty" query:"category"`
	Enabled    *bool     `json:"enabled,omitempty" query:"enabled"`
	IsTemplate *bool     `json:"isTemplate,omitempty" query:"isTemplate"`
	Page       *int      `json:"page,omitempty" query:"page"`
	PerPage    *int      `json:"perPage,omitempty" query:"perPage"`
	Projection model.Map `json:"-"`
}

// GetPage returns the page number (defaults to 1).
func (q *PluginQuery) GetPage() int {
	if q.Page != nil && *q.Page > 0 {
		return *q.Page
	}
	return 1
}

// GetPerPage returns the per-page size (defaults to 20).
func (q *PluginQuery) GetPerPage() int {
	if q.PerPage != nil && *q.PerPage > 0 {
		return *q.PerPage
	}
	return 20
}

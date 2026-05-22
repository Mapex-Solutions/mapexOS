package datasources

import (
	"fmt"

	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

const (
	AuthTypeNone        = "none"
	AuthTypeAPIKey      = "apiKey"
	AuthTypeJWT         = "jwt"
	AuthTypeIPWhitelist = "ip_whitelist"
	AuthTypeOAuth2      = "oauth2"
)

type DataSourceId struct {
	DataSourceId string `params:"dataSourceId" validate:"required,mongoid"`
}

type AuthAPIKey struct {
	// Where to expect the key; you can enable one or both.
	Type      string `json:"type" validate:"required,oneof=header query"`
	FieldName string `json:"fieldName,omitempty" validate:"required,min=1"`
	Key       string `json:"key,omitempty" validate:"required,min=20"`
}

// JWT (ingress auth via bearer token validation)
type AuthJWT struct {
	Secret     string  `json:"secret,omitempty" validate:"required,min=6"`
	Algorithms string  `json:"algorithms,omitempty" validate:"required,oneof=HS256 HS512"`
	HeaderName *string `json:"headerName,omitempty" validate:"omitempty,min=1"`
}

// IP Whitelist (ingress by source IP/CIDR)
type AuthIPWhitelist struct {
	// Accepts IPv4/IPv6 CIDRs
	CIDRs []string `json:"cidrs" validate:"required,min=1"`
}

// OAuth2 (client credentials — best for server→server pulls/pushes we initiate)
type AuthOAuth2 struct {
	JWKSURL      string  `json:"jwksURL" validate:"required,url"`
	ClientID     *string `json:"clientId" validate:"omitempty,min=3"`
	ClientSecret *string `json:"clientSecret" validate:"omitempty,min=6"`
}

// NONE
type AuthNone struct{}

type DataSourceAuth struct {
	Type        string           `json:"type" validate:"required,oneof=apiKey jwt ip_whitelist oauth2 none"`
	APIKey      *AuthAPIKey      `json:"apiKey,omitempty" validate:"omitempty"`
	JWT         *AuthJWT         `json:"jwt,omitempty" validate:"omitempty"`
	IPWhitelist *AuthIPWhitelist `json:"ipWhitelist,omitempty" validate:"omitempty"`
	OAuth2      *AuthOAuth2      `json:"oauth2,omitempty" validate:"omitempty"`
	None        *AuthNone        `json:"none,omitempty" validate:"omitempty"`
}

type WorkingHours struct {
	Enabled  bool    `json:"enabled" validate:"required"`
	Days     []int   `json:"days,omitempty" validate:"omitempty,min=0,max=6"`
	StartAt  *string `json:"startAt,omitempty" validate:"omitempty"`
	EndAt    *string `json:"endAt,omitempty" validate:"omitempty"`
	TimeZone *string `json:"timeZone,omitempty" validate:"omitempty"`
}

type RateLimit struct {
	Type           string `json:"type" validate:"required,oneof=second minute hour"`
	Value          int    `json:"value" validate:"required,min=1"`
	BurstCapacity  int    `json:"burstCapacity" validate:"required,min=1"`
	ActionOnExceed string `json:"actionOnExceed" validate:"required,oneof=drop queue"`
}

type AssetBind struct {
	Type string `json:"type" validate:"required,oneof=fixedAssetId uuidField "`
	Data struct {
		UUIDField []string `json:"uuidField,omitempty" validate:"omitempty,min=1"`
		AssetId   *string  `json:"assetId,omitempty" validate:"omitempty,min=1"`
	} `json:"data" validate:"required"`
}

type DataSourceCreate struct {
	Name        string  `json:"name" validate:"required,min=1"`
	Enabled     bool    `json:"enabled"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Mode         string         `json:"mode" validate:"required,oneof=pull push X"`
	Protocol     string         `json:"protocol" validate:"omitempty,oneof=mqtt http"`
	WorkingHrs   *WorkingHours  `json:"workingHours,omitempty" validate:"omitempty"`
	RateLimit    *RateLimit     `json:"rateLimit,omitempty" validate:"omitempty"`
	Auth         DataSourceAuth `json:"auth" validate:"required"`
	AssetBind    AssetBind      `json:"assetBind" validate:"required"`

	// Multi-tenant fields (populated automatically by coverage middleware)
	OrgID   *model.ObjectId `json:"orgId,omitempty" validate:"omitempty"`
	PathKey *string         `json:"pathKey,omitempty" validate:"omitempty"`
}

type DataSourceUpdate struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1"`
	Enabled     *bool   `json:"enabled,omitempty" validate:"omitempty"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Mode         *string         `json:"mode,omitempty" validate:"omitempty,oneof=pull push X"`
	Protocol     *string         `json:"protocol,omitempty" validate:"omitempty,oneof=mqtt http"`
	WorkingHrs   *WorkingHours   `json:"workingHours,omitempty" validate:"omitempty"`
	RateLimit    *RateLimit      `json:"rateLimit,omitempty" validate:"omitempty"`
	Auth         *DataSourceAuth `json:"auth,omitempty" validate:"omitempty"`
	AssetBind    *AssetBind      `json:"assetBind,omitempty" validate:"omitempty"`
}

type DataSourceResponse struct {
	ID          *common.ObjectID `json:"id,omitempty"`
	Name        *string          `json:"name,omitempty"`
	Enabled     *bool            `json:"enabled,omitempty"`
	Description *string          `json:"description,omitempty"`
	Mode         *string          `json:"mode,omitempty"`
	Protocol     *string          `json:"protocol,omitempty"`

	// Multi-tenant fields
	OrgId      *common.ObjectID `json:"orgId,omitempty"`
	PathKey    *string          `json:"pathKey,omitempty"`
	CustomerID *common.ObjectID `json:"customerId,omitempty"`

	WorkingHrs *WorkingHours    `json:"workingHours,omitempty"`
	RateLimit  *RateLimit       `json:"rateLimit,omitempty"`
	Auth       *DataSourceAuth  `json:"auth,omitempty"`
	AssetBind  *AssetBind       `json:"assetBind,omitempty"`
	Created    *common.NullTime `json:"created,omitempty"`
	Updated    *common.NullTime `json:"updated,omitempty"`
}

func (d *DataSourceResponse) SetCreated(t *common.NullTime) { d.Created = t }
func (d *DataSourceResponse) SetUpdated(t *common.NullTime) { d.Updated = t }

func (d *DataSourceAuth) Transform() error {
	// Cross-field rules + normalization (no full revalidation here)
	switch d.Type {
	case "apiKey":
		if d.APIKey == nil {
			return fmt.Errorf("field 'apiKey' must be provided when type is 'apiKey'")
		}
		// normalize: clear others
		d.JWT, d.IPWhitelist, d.OAuth2, d.None = nil, nil, nil, nil

	case "jwt":
		if d.JWT == nil {
			return fmt.Errorf("field 'jwt' must be provided when type is 'jwt'")
		}
		d.APIKey, d.IPWhitelist, d.OAuth2, d.None = nil, nil, nil, nil

	case "ip_whitelist":
		if d.IPWhitelist == nil {
			return fmt.Errorf("field 'ipWhitelist' must be provided when type is 'ip_whitelist'")
		}
		d.APIKey, d.JWT, d.OAuth2, d.None = nil, nil, nil, nil

	case "oauth2":
		if d.OAuth2 == nil {
			return fmt.Errorf("field 'oauth2' must be provided when type is 'oauth2'")
		}
		d.APIKey, d.JWT, d.IPWhitelist, d.None = nil, nil, nil, nil

	case "none":
		if d.None == nil {
			return fmt.Errorf("field 'none' must be provided when type is 'none'")
		}
		d.APIKey, d.JWT, d.IPWhitelist, d.OAuth2 = nil, nil, nil, nil

	default:
		return fmt.Errorf("unsupported auth type: %s", d.Type)
	}
	return nil
}

// DataSourceQuery represents query parameters for listing data sources.
// Embeds BaseQueryDTO for standard pagination, sorting, and hierarchy support.
//
// Standard fields (from BaseQueryDTO):
//   - Projection: comma-separated fields to return
//   - Page: page number (default: 1)
//   - PerPage: items per page (default: 20)
//   - Sort: sort order (default: "created:desc")
//   - IncludeChildren: include child orgs hierarchically (default: false)
//
// Module-specific filters:
//   - Name: filter by data source name (partial match)
//   - Enabled: filter by enabled status (true/false)
//   - Mode: filter by mode (pull/push)
//   - Protocol: filter by protocol (http/mqtt)
//
// Organization filtering is handled automatically via RequestContext:
//   - No manual orgId needed
//   - Context-aware filtering via X-Org-Context header
//   - Hierarchical queries via includeChildren parameter
type DataSourceQuery struct {
	query.BaseQueryDTO

	// Module-specific filters
	Name      *string `query:"name" validate:"omitempty,max=100"`
	Enabled   *bool   `query:"enabled" validate:"omitempty"`
	Mode      *string `query:"mode" validate:"omitempty,oneof=pull push"`
	Protocol  *string `query:"protocol" validate:"omitempty,oneof=http"`
	Auth      *string `query:"auth" validate:"omitempty,oneof=apiKey jwt ip_whitelist oauth2 none"`
	AssetBind *string `query:"assetBind" validate:"omitempty,oneof=fixedAssetId uuidField"`
}

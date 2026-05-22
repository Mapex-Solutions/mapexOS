package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/* Value Objects (mirror DTOs) */

type WorkingHours struct {
	Enabled  bool    `bson:"enabled"`
	Days     []int   `bson:"days,omitempty"`
	StartAt  *string `bson:"startAt,omitempty"`
	EndAt    *string `bson:"endAt,omitempty"`
	TimeZone *string `bson:"timeZone,omitempty"`
}

type RateLimit struct {
	Type           string `bson:"type"` // second | minute | hour
	Value          int    `bson:"value"`
	BurstCapacity  int    `bson:"burstCapacity"`
	ActionOnExceed string `bson:"actionOnExceed"` // drop | queue
}

type AssetBindData struct {
	AssetId   *string  `bson:"assetId,omitempty"`
	UUIDField []string `bson:"uuidField,omitempty"`
}

type AssetBind struct {
	Type string        `bson:"type"`
	Data AssetBindData `bson:"data"`
}

/* Auth (strongly-typed like your DTOs) */

type AuthAPIKey struct {
	// Where to expect the key; you can enable one or both.
	Type      string `bson:"type"`                // header | query
	FieldName string `bson:"fieldName,omitempty"` // normalized to lower-case for BSON
	Key       string `bson:"key,omitempty"`
}

type AuthJWT struct {
	Secret     string  `bson:"secret,omitempty"`
	Algorithms string  `bson:"algorithms,omitempty"` // HS256 | HS512
	HeaderName *string `bson:"headerName,omitempty"` // Header name for JWT (default: Authorization)
}

type AuthIPWhitelist struct {
	CIDRs              []string `bson:"cidrs"`
	AllowPrivateRanges *bool    `bson:"allowPrivateRanges,omitempty"`
}

type AuthOAuth2 struct {
	JWKSURL      string  `bson:"jwksURL"`
	ClientID     *string `bson:"clientId,omitempty"`
	ClientSecret *string `bson:"clientSecret,omitempty"`
}

type AuthNone struct{}

type DataSourceAuth struct {
	Type        string           `bson:"type"` // apiKey | jwt | ip_whitelist | oauth2 | none
	APIKey      *AuthAPIKey      `bson:"apiKey,omitempty"`
	JWT         *AuthJWT         `bson:"jwt,omitempty"`
	IPWhitelist *AuthIPWhitelist `bson:"ipWhitelist,omitempty"`
	OAuth2      *AuthOAuth2      `bson:"oauth2,omitempty"`
	None        *AuthNone        `bson:"none,omitempty"`
}

/* Entity */

type DataSource struct {
	ID          model.ObjectId `bson:"_id,omitempty"`
	Name        string         `bson:"name"`
	Enabled     bool           `bson:"enabled"`
	Description *string        `bson:"description,omitempty"`
	Mode         string         `bson:"mode"`               // pull | push | X
	Protocol     string         `bson:"protocol,omitempty"` // mqtt | http

	// Multi-tenant fields
	OrgID      model.ObjectId  `bson:"orgId"`
	PathKey    string          `bson:"pathKey"`              // Hierarchical path for range queries
	CustomerID *model.ObjectId `bson:"customerId,omitempty"` // Tenant anchor (denormalized)

	WorkingHrs *WorkingHours  `bson:"workingHours,omitempty"`
	RateLimit  *RateLimit     `bson:"rateLimit,omitempty"`
	Auth       DataSourceAuth `bson:"auth"`
	AssetBind  AssetBind      `bson:"assetBind"`
	Created    time.Time      `bson:"created"`
	Updated    time.Time      `bson:"updated"`
}

// Keep these as requested
func (u *DataSource) GetCreated() time.Time { return u.Created }
func (u *DataSource) GetUpdated() time.Time { return u.Updated }

/* Update structs (all fields optional) */

type WorkingHoursUpdate struct {
	Enabled  *bool   `bson:"enabled,omitempty"`
	Days     *[]int  `bson:"days,omitempty"`
	StartAt  *string `bson:"startAt,omitempty"`
	EndAt    *string `bson:"endAt,omitempty"`
	TimeZone *string `bson:"timeZone,omitempty"`
}

type RateLimitUpdate struct {
	Type           *string `bson:"type,omitempty"`
	Value          *int    `bson:"value,omitempty"`
	BurstCapacity  *int    `bson:"burstCapacity,omitempty"`
	ActionOnExceed *string `bson:"actionOnExceed,omitempty"`
}

type AssetBindDataUpdate struct {
	AssetId   *string   `bson:"assetId,omitempty"`
	UUIDField *[]string `bson:"uuidField,omitempty"`
}

type AssetBindUpdate struct {
	Type *string              `bson:"type,omitempty"`
	Data *AssetBindDataUpdate `bson:"data,omitempty"`
}

// Auth updates (partial)
type AuthAPIKeyUpdate struct {
	Type      *string `bson:"type,omitempty"` // header | query
	FieldName *string `bson:"fieldName,omitempty"`
	Key       *string `bson:"key,omitempty"`
}

type AuthJWTUpdate struct {
	Secret     *string   `bson:"secret,omitempty"`
	Algorithms *[]string `bson:"algorithms,omitempty"`
	HeaderName *string   `bson:"headerName,omitempty"`
}

type AuthIPWhitelistUpdate struct {
	CIDRs              *[]string `bson:"cidrs,omitempty"`
	AllowPrivateRanges *bool     `bson:"allowPrivateRanges,omitempty"`
}

type AuthOAuth2Update struct {
	JWKSURL      *string `bson:"jwksURL,omitempty"`
	ClientID     *string `bson:"clientId,omitempty"`
	ClientSecret *string `bson:"clientSecret,omitempty"`
}

type AuthNoneUpdate struct{}

type DataSourceAuthUpdate struct {
	Type        *string                `bson:"type,omitempty"`
	APIKey      *AuthAPIKeyUpdate      `bson:"apiKey,omitempty"`
	JWT         *AuthJWTUpdate         `bson:"jwt,omitempty"`
	IPWhitelist *AuthIPWhitelistUpdate `bson:"ipWhitelist,omitempty"`
	OAuth2      *AuthOAuth2Update      `bson:"oauth2,omitempty"`
	None        *AuthNoneUpdate        `bson:"none,omitempty"`
}

// PATCH/UPDATE payload (every field optional)
type DataSourceUpdate struct {
	Name        *string `bson:"name,omitempty"`
	Enabled     *bool   `bson:"enabled,omitempty"`
	Description *string `bson:"description,omitempty"`
	Mode         *string `bson:"mode,omitempty"`
	Protocol     *string `bson:"protocol,omitempty"`

	// Multi-tenant fields (usually not updated, but available)
	OrgID      *model.ObjectId `bson:"orgId,omitempty"`
	PathKey    *string         `bson:"pathKey,omitempty"`
	CustomerID *model.ObjectId `bson:"customerId,omitempty"`

	WorkingHrs *WorkingHoursUpdate   `bson:"workingHours,omitempty"`
	RateLimit  *RateLimitUpdate      `bson:"rateLimit,omitempty"`
	Auth       *DataSourceAuthUpdate `bson:"auth,omitempty"`
	AssetBind  *AssetBindUpdate      `bson:"assetBind,omitempty"`

	// Auditing (usually only Updated is set by service/repo)
	Created *time.Time `bson:"created,omitempty"`
	Updated *time.Time `bson:"updated,omitempty"`
}

func (u *DataSourceUpdate) GetCreated() *time.Time { return u.Created }
func (u *DataSourceUpdate) GetUpdated() *time.Time { return u.Updated }

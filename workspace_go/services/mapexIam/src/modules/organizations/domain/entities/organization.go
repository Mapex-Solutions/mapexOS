package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type Organization struct {
	ID           model.ObjectId  `bson:"_id,omitempty"`
	Name         string          `bson:"name"`
	Type         string          `bson:"type"` // "vendor" | "customer" | "site" | "building" | "floor" | "zone"
	ParentOrgID  *model.ObjectId `bson:"parentOrgId,omitempty"`

	// Hierarchical PathKey for range queries
	Code         string  `bson:"code"`    // Local code (e.g., "0001")
	PathKey      string  `bson:"pathKey"` // Full hierarchical path (e.g., "000001/000001/0001")
	Depth        int     `bson:"depth"`   // Level in tree (0=vendor, 1=customer, 2=site, etc.)

	// Denormalized tenant anchor for fast queries
	CustomerID   *model.ObjectId `bson:"customerId,omitempty"`

	// Metadata
	ChildCount        int                `bson:"childCount"` // Number of direct children
	Logo              *string            `bson:"logo,omitempty"`
	Enabled           bool               `bson:"enabled"`
	Address           *Address           `bson:"address,omitempty"`
	Phone             *string            `bson:"phone,omitempty"`
	AuthConfig        AuthConfig         `bson:"authConfig"`
	Created           time.Time          `bson:"created"`
	Updated           time.Time          `bson:"updated"`
	AccessPolicy      AccessPolicy       `bson:"accessPolicy"`
}

type Address struct {
	City    string `bson:"city"`
	State   string `bson:"state"`
	Country string `bson:"country"`
	ZipCode string `bson:"zipCode"`
}

type AuthConfig struct {
	ProviderType     string                 `bson:"providerType"`
	IssuerURL        *string                `bson:"issuerUrl,omitempty"`
	ClientID         *string                `bson:"clientId,omitempty"`
	JWTClaimMappings map[string]string      `bson:"jwtClaimMappings,omitempty"`
	Metadata         map[string]interface{} `bson:"metadata,omitempty"`
}

type AccessPolicy struct {
	RolePolicy   string `bson:"rolePolicy"`   // "merge" | "strict"
	DefaultScope string `bson:"defaultScope"` // "local" | "recursive" (UX default only)

	// REMOVED: AllowDirectPermissions
	// V1 uses pure role-based architecture (no direct permissions)
	// If needed in future, permissions would be assigned directly to users via Membership
}

func (u *Organization) GetCreated() time.Time { return u.Created }
func (u *Organization) GetUpdated() time.Time { return u.Updated }
func (u *Organization) GetPathKey() string    { return u.PathKey }
func (u *Organization) GetCustomerID() *model.ObjectId { return u.CustomerID }

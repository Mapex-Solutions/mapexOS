package entities

import (
	"time"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexVault/credentials"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/**
 * Connection Status
 *
 * Domain re-exports of cross-service contract status values.
 */

type ConnectionStatus = contracts.ConnectionStatus

const (
	ConnectionStatusActive  = contracts.ConnectionStatusActive
	ConnectionStatusRevoked = contracts.ConnectionStatusRevoked
	ConnectionStatusExpired = contracts.ConnectionStatusExpired
)

/**
 * Connection Entity
 * Tracks which external account is connected to the platform.
 * Each connection is linked to a credential in the vault.
 * The connection exists independently of workflows — a user can
 * connect an account without having any automation running.
 */

type Connection struct {
	ID           model.ObjectId   `bson:"_id,omitempty" json:"id"`
	Provider     string           `bson:"provider" json:"provider"`       // "instagram", "tiktok", "twitter"
	AccountId    string           `bson:"accountId" json:"accountId"`     // provider-specific account ID
	AccountName  string           `bson:"accountName" json:"accountName"` // display name (@handle)
	Status       ConnectionStatus `bson:"status" json:"status"`
	CredentialId model.ObjectId   `bson:"credentialId" json:"credentialId"`
	UserId       model.ObjectId   `bson:"userId" json:"userId"`
	OrgId        *model.ObjectId  `bson:"orgId,omitempty" json:"orgId,omitempty"`
	PathKey      string           `bson:"pathKey" json:"pathKey"`
	Scopes       []string         `bson:"scopes" json:"scopes"`
	ConnectedAt  time.Time        `bson:"connectedAt" json:"connectedAt"`
	LastUsedAt   *time.Time       `bson:"lastUsedAt,omitempty" json:"lastUsedAt,omitempty"`
	Created      time.Time        `bson:"created" json:"created"`
	Updated      time.Time        `bson:"updated" json:"updated"`
}

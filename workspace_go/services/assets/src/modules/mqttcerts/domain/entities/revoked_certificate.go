package entities

import (
	"time"

	"assets/src/modules/mqttcerts/domain/constants"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// RevokedCertificate is one row in mqttRevokedCertificates. Status is
// implicit — every row in this collection is revoked. Mongo TTL on
// revokedAt auto-deletes rows after 30 days.
//
// Domain entity: bson-only tags, no json tags, no cross-service
// contract imports — keeps the persistence shape isolated from the
// wire shape and the contracts package.
type RevokedCertificate struct {
	ID          model.ObjectId             `bson:"_id,omitempty"`
	Serial      string                     `bson:"serial"`
	Fingerprint string                     `bson:"fingerprint"`
	AssetUUID   string                     `bson:"assetUUID"`
	OrgID       string                     `bson:"orgId"`
	SubjectCN   string                     `bson:"subjectCN"`
	IssuedAt    time.Time                  `bson:"issuedAt"`
	RevokedAt   time.Time                  `bson:"revokedAt"`
	Reason      constants.RevocationReason `bson:"reason"`
	Created     time.Time                  `bson:"created"`
	Updated     time.Time                  `bson:"updated"`
}

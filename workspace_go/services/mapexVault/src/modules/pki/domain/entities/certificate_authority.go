package entities

import (
	"time"

	"mapexVault/src/modules/pki/domain/constants"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// CertificateAuthority is the mapexVault-side persistence record for a
// PKI CA (root or intermediate). Private key bytes are encrypted via
// the existing envelope primitive (Master Key → DEK → payload). Cert
// PEM is plaintext — public material.
//
// JSON tags ABSENT on purpose: domain entities round-trip Mongo only.
// Wire format lives in packages/contracts; conversion happens at the
// application boundary via a mapper. Per /go-arch §6 strictly.
type CertificateAuthority struct {
	ID          model.ObjectId   `bson:"_id,omitempty"`
	Kind        constants.CAKind `bson:"kind"`
	IsSystem    bool             `bson:"isSystem"`
	SubjectCN   string           `bson:"subjectCN"`
	Fingerprint string           `bson:"fingerprint"`
	NotBefore   time.Time        `bson:"notBefore"`
	NotAfter    time.Time        `bson:"notAfter"`

	EncryptedDEK []byte `bson:"encryptedDEK"`
	DekNonce     []byte `bson:"dekNonce"`
	EncryptedKey []byte `bson:"encryptedKey"`
	KeyNonce     []byte `bson:"keyNonce"`
	CertPEM      []byte `bson:"certPEM"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

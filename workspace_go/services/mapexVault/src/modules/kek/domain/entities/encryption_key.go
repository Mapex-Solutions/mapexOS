package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// EncryptionKey is the mapexVault-side persistence record for a platform
// key-encryption key (KEK). The KEK is AES-256-GCM envelope-encrypted with the
// Master Key (Master Key -> DEK -> payload); only the envelope fields persist.
//
// JSON tags ABSENT on purpose: domain entities round-trip Mongo only. Wire
// format lives in packages/contracts; conversion happens at the application
// boundary. Per /go-arch §6.
type EncryptionKey struct {
	ID       model.ObjectId `bson:"_id,omitempty"`
	Context  string         `bson:"context"`
	IsSystem bool           `bson:"isSystem"`

	EncryptedDEK []byte `bson:"encryptedDEK"`
	DekNonce     []byte `bson:"dekNonce"`
	EncryptedKey []byte `bson:"encryptedKey"`
	KeyNonce     []byte `bson:"keyNonce"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

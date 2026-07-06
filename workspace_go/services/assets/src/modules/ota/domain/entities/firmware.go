package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// FirmwareStatus is the lifecycle state of a firmware artifact.
type FirmwareStatus string

const (
	FirmwarePendingUpload FirmwareStatus = "PENDING_UPLOAD"
	FirmwareReady         FirmwareStatus = "READY"
	FirmwareActive        FirmwareStatus = "ACTIVE"
	FirmwareDeprecated    FirmwareStatus = "DEPRECATED"
	FirmwareRevoked       FirmwareStatus = "REVOKED"
	// FirmwareAbandoned marks an upload that was started (:init) but never
	// finalized (:complete) before its abandon-check timer fired.
	FirmwareAbandoned FirmwareStatus = "ABANDONED"
	// FirmwarePurged marks an artifact whose .bin was deleted from object
	// storage after the referencing plan closed.
	FirmwarePurged FirmwareStatus = "PURGED"
)

// Firmware is a first-class firmware artifact: the bytes (in object storage)
// that migrate a device to the target template's version. The asset template
// stays the version of record; the firmware carries the binary + integrity
// metadata and references the target template it produces.
type Firmware struct {
	ID                model.ObjectId `bson:"_id,omitempty"`
	OrgID             model.ObjectId `bson:"orgId"`
	TargetTemplateID  model.ObjectId `bson:"targetTemplateId"`
	Version           string         `bson:"version"`
	Filename          string         `bson:"filename"`
	ObjectKey         string         `bson:"objectKey"`
	Size              int64          `bson:"size"`
	Checksum          string         `bson:"checksum"`
	ChecksumAlgorithm string         `bson:"checksumAlgorithm"`
	Signature         *string        `bson:"signature,omitempty"`
	Status            FirmwareStatus `bson:"status"`
	Created           time.Time      `bson:"created"`
	Updated           time.Time      `bson:"updated"`
}

func (f *Firmware) GetCreated() time.Time { return f.Created }

// CanReference reports whether a plan may reference this firmware. Only a READY
// or ACTIVE artifact has confirmed, intact bytes in object storage.
func (f *Firmware) CanReference() bool {
	return f.Status == FirmwareReady || f.Status == FirmwareActive
}

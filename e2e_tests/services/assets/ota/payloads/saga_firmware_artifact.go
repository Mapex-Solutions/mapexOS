// Package payloads holds the pure builders the OTA saga building blocks consume:
// a deterministic firmware artifact and an OTA plan-create request. They depend
// only on the run's runID (for isolation) — never on the bag or HTTP.
package payloads

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// FirmwareArtifact is a deterministic, self-describing firmware blob for OTA
// saga journeys. Built once from the runID so every derived value (bytes,
// sha256, size, version, filename) is stable within a run and unique across
// runs: the InitFirmware step declares size+sha256, UploadFirmwareBinary PUTs
// Bytes() to the presigned URL, and the device sim verifies the same sha256+size
// after downloading.
type FirmwareArtifact struct {
	runID  string
	bytes  []byte
	sha256 string
}

// NewFirmwareArtifact builds the deterministic artifact for the given runID.
func NewFirmwareArtifact(runID string) *FirmwareArtifact {
	// A few KB of deterministic content keyed by runID, so two runs never share
	// a checksum and the size is a realistic non-zero firmware size.
	block := fmt.Appendf(nil, "MAPEX-OTA-FIRMWARE;run=%s;", runID)
	buf := make([]byte, 0, 4096+len(block))
	for len(buf) < 4096 {
		buf = append(buf, block...)
	}
	sum := sha256.Sum256(buf)
	return &FirmwareArtifact{
		runID:  runID,
		bytes:  buf,
		sha256: base64.StdEncoding.EncodeToString(sum[:]),
	}
}

// Bytes returns the firmware payload uploaded to (and downloaded from) storage.
func (f *FirmwareArtifact) Bytes() []byte { return f.bytes }

// SHA256 returns the base64-encoded sha256 of Bytes() — the S3
// x-amz-checksum-sha256 encoding the assets service signs into the presigned PUT.
// Declared on init and verified by the device sim after download.
func (f *FirmwareArtifact) SHA256() string { return f.sha256 }

// Size returns the byte length of Bytes().
func (f *FirmwareArtifact) Size() int64 { return int64(len(f.bytes)) }

// Version returns a per-run firmware version stamped with the runID for isolation.
func (f *FirmwareArtifact) Version() string { return "1.0.0-" + f.runID }

// Filename returns the artifact filename stamped with the runID.
func (f *FirmwareArtifact) Filename() string { return fmt.Sprintf("firmware-%s.bin", f.runID) }

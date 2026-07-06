package ports

import (
	"context"
	"time"
)

// FirmwareStorePort is the object-storage boundary for firmware artifacts. The
// firmware bytes never traverse the Asset MS: the client uploads and the device
// downloads directly via presigned URLs; the service only mints URLs and
// confirms/deletes objects.
type FirmwareStorePort interface {
	// PresignPut returns a short-TTL presigned PUT URL bound to the given
	// base64-encoded SHA-256 (the x-amz-checksum-sha256 header) so the object
	// store validates the uploaded bytes' integrity on write. The client MUST
	// send that exact x-amz-checksum-sha256 header value when uploading.
	PresignPut(ctx context.Context, key, checksumSHA256B64 string, ttl time.Duration) (string, error)

	// PresignGet returns a short-TTL presigned GET URL for the device download.
	PresignGet(ctx context.Context, key string, ttl time.Duration) (string, error)

	// Stat confirms the object exists and returns its size and stored
	// base64-encoded SHA-256 checksum (empty if it was uploaded without one).
	Stat(ctx context.Context, key string) (size int64, checksumSHA256B64 string, err error)

	// Delete removes the object. Idempotent: deleting a missing object is a
	// no-op (returns nil).
	Delete(ctx context.Context, key string) error
}

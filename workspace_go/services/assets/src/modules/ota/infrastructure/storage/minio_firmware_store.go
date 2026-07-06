package storage

import (
	"context"
	"net/http"
	"time"

	"assets/src/modules/ota/application/ports"

	minioModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/minio"
	minio "github.com/minio/minio-go/v7"
)

// minioFirmwareStore implements ports.FirmwareStorePort against the firmware
// MinIO bucket. All operations use the raw minio-go client (with the client's
// configured bucket) so the key space is identical across presign/stat/delete.
type minioFirmwareStore struct {
	client *minioModel.MinIOClient
}

// NewMinIOFirmwareStore returns a FirmwareStorePort backed by the firmware
// bucket. The injected MinIOClient must be configured with that bucket.
func NewMinIOFirmwareStore(client *minioModel.MinIOClient) ports.FirmwareStorePort {
	return &minioFirmwareStore{client: client}
}

func (s *minioFirmwareStore) PresignPut(ctx context.Context, key, checksumSHA256B64 string, ttl time.Duration) (string, error) {
	extra := http.Header{}
	extra.Set("x-amz-checksum-sha256", checksumSHA256B64)

	u, err := s.client.GetRawClient().PresignHeader(
		ctx, http.MethodPut, s.client.GetBucketName(), key, ttl, nil, extra,
	)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (s *minioFirmwareStore) PresignGet(ctx context.Context, key string, ttl time.Duration) (string, error) {
	u, err := s.client.GetRawClient().PresignedGetObject(
		ctx, s.client.GetBucketName(), key, ttl, nil,
	)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (s *minioFirmwareStore) Stat(ctx context.Context, key string) (int64, string, error) {
	info, err := s.client.GetRawClient().StatObject(
		ctx, s.client.GetBucketName(), key, minio.StatObjectOptions{Checksum: true},
	)
	if err != nil {
		return 0, "", err
	}
	return info.Size, info.ChecksumSHA256, nil
}

func (s *minioFirmwareStore) Delete(ctx context.Context, key string) error {
	return s.client.GetRawClient().RemoveObject(
		ctx, s.client.GetBucketName(), key, minio.RemoveObjectOptions{},
	)
}

// Compile-time check that minioFirmwareStore implements the port.
var _ ports.FirmwareStorePort = (*minioFirmwareStore)(nil)

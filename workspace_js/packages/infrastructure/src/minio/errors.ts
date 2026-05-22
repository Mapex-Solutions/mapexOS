/**
 * MinIO Errors
 * Same structure as workspace_go/packages/infrastructure/minio/errors.go
 */

// ErrObjectNotFound is returned when the requested object does not exist.
export const ErrObjectNotFound = new Error('object not found');

// ErrBucketNotFound is returned when the bucket does not exist.
export const ErrBucketNotFound = new Error('bucket not found');

// ErrInvalidConfig is returned when the configuration is invalid.
export const ErrInvalidConfig = new Error('invalid minio configuration');

// ErrNilData is returned when attempting to upload nil data.
export const ErrNilData = new Error('data cannot be nil');

// ErrEmptyKey is returned when the object key is empty.
export const ErrEmptyKey = new Error('object key cannot be empty');

// ErrConnectionFailed is returned when MinIO connection fails.
export const ErrConnectionFailed = new Error('failed to connect to MinIO');

// ErrUploadFailed is returned when object upload fails.
export const ErrUploadFailed = new Error('failed to upload object');

// ErrDownloadFailed is returned when object download fails.
export const ErrDownloadFailed = new Error('failed to download object');

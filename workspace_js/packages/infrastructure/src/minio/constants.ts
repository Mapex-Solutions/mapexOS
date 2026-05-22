/**
 * MinIO Constants
 * Same structure as workspace_go/packages/infrastructure/minio/constants.go
 */

// Default content types for common file formats.
export const ContentTypeJSON = 'application/json';
export const ContentTypeBinary = 'application/octet-stream';
export const ContentTypeText = 'text/plain';
export const ContentTypeJavaScript = 'application/javascript';
export const ContentTypeMessagePack = 'application/msgpack';

// Default configuration values.
export const DefaultRegion = 'us-east-1';
export const DefaultMaxRetries = 3;

// Bucket prefixes for MapexOS.
export const BucketTemplates = 'mapex-templates';
export const BucketExports = 'mapex-exports';

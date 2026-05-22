/**
 * MinIO Package Exports
 * Same structure as workspace_go/packages/infrastructure/minio/
 */

// Main factory function
export { New } from './minio';

// Client class
export { MinIOClient } from './methods';

// Types
export type {
	Config,
	PutOptions,
	GetResult,
	ObjectInfo,
	ListOptions,
	StreamReader,
} from './types';

// Constants
export {
	ContentTypeJSON,
	ContentTypeBinary,
	ContentTypeText,
	ContentTypeJavaScript,
	ContentTypeMessagePack,
	DefaultRegion,
	DefaultMaxRetries,
	BucketTemplates,
	BucketExports,
} from './constants';

// Errors
export {
	ErrObjectNotFound,
	ErrBucketNotFound,
	ErrInvalidConfig,
	ErrNilData,
	ErrEmptyKey,
	ErrConnectionFailed,
	ErrUploadFailed,
	ErrDownloadFailed,
} from './errors';

// Internals (for advanced use cases)
export { prefixKey, buildPutOptions, convertObjectInfo, isNotFoundError } from './internals';

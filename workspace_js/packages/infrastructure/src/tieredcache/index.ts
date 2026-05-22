/**
 * TieredCache Package Exports
 * Same structure as workspace_go/packages/infrastructure/tieredcache/
 *
 * Note: Some exports use "Cache" prefix to avoid conflicts with MinIO exports
 */

// Main factory function (renamed to avoid conflict with MinIO.New)
export { New as NewTieredCache } from './tieredcache';

// Client class
export { TieredCacheClient } from './methods';

// Types (renamed to avoid conflict with MinIO.Config)
export type { Config as TieredCacheConfig, CacheStats, CacheEntry, LocalCacheLoader, L0Entry } from './types';

// Enums and helpers
export { CacheTier, tierToString } from './types';

// Constants
export {
	DefaultL0MaxSize,
	DefaultL0MaxItems,
	DefaultL0TTL,
	L0BufferItems,
	DefaultL1MaxSize,
	DefaultL1TTL,
	DefaultL1Dir,
	L1FileExtension,
	L1MetaExtension,
	L1SubdirDepth,
} from './constants';

// Errors (renamed to avoid conflicts with MinIO errors)
export {
	ErrCacheMiss,
	ErrL0InitFailed,
	ErrL1InitFailed,
	ErrL1DirNotWritable,
	ErrNoLoader,
	ErrEntryExpired,
	ErrInvalidConfig as ErrCacheInvalidConfig,
	ErrNilData as ErrCacheNilData,
	ErrEmptyKey as ErrCacheEmptyKey,
} from './errors';

// Internals (renamed prefixKey to avoid conflict with MinIO)
export {
	prefixKey as cacheKeyPrefix,
	hashKey,
	l1FilePath,
	l1MetaPath,
	writeL1,
	readL1,
	deleteL1Files,
	isL0Enabled,
	isL1Enabled,
} from './internals';

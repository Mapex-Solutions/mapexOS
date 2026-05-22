/**
 * TieredCache Errors
 * Same structure as workspace_go/packages/infrastructure/tieredcache/errors.go
 */

// ErrCacheMiss is returned when the key is not found in any cache tier.
export const ErrCacheMiss = new Error('cache miss: key not found');

// ErrL0InitFailed is returned when L0 (RAM) cache initialization fails.
export const ErrL0InitFailed = new Error('failed to initialize L0 RAM cache');

// ErrL1InitFailed is returned when L1 (Disk) cache initialization fails.
export const ErrL1InitFailed = new Error('failed to initialize L1 disk cache');

// ErrL1DirNotWritable is returned when L1 directory is not writable.
export const ErrL1DirNotWritable = new Error('L1 cache directory is not writable');

// ErrNoLoader is returned when trying to load from L2 without a loader.
export const ErrNoLoader = new Error('no L2 loader configured');

// ErrEntryExpired is returned when the cached entry has expired.
export const ErrEntryExpired = new Error('cache entry expired');

// ErrInvalidConfig is returned when the configuration is invalid.
export const ErrInvalidConfig = new Error('invalid cache configuration');

// ErrNilData is returned when attempting to cache nil data.
export const ErrNilData = new Error('cannot cache nil data');

// ErrEmptyKey is returned when the cache key is empty.
export const ErrEmptyKey = new Error('cache key cannot be empty');

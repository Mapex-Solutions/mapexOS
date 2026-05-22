/**
 * TieredCache Constants
 * Same structure as workspace_go/packages/infrastructure/tieredcache/constants.go
 */

// Default L0 (RAM) configuration.

// DefaultL0MaxSize is the default maximum RAM cache size (256 MB).
export const DefaultL0MaxSize = 256 * 1024 * 1024;

// DefaultL0MaxItems is the default maximum number of items in RAM cache.
export const DefaultL0MaxItems = 100_000;

// DefaultL0TTL is the default TTL for RAM cache entries (5 minutes in ms).
export const DefaultL0TTL = 5 * 60 * 1000;

// L0BufferItems is the number of items to buffer before eviction.
export const L0BufferItems = 64;

// Default L1 (Disk) configuration.

// DefaultL1MaxSize is the default maximum disk cache size (10 GB).
export const DefaultL1MaxSize = 10 * 1024 * 1024 * 1024;

// DefaultL1TTL is the default TTL for disk cache entries (1 hour in ms).
export const DefaultL1TTL = 60 * 60 * 1000;

// DefaultL1Dir is the default directory for disk cache.
export const DefaultL1Dir = '/tmp/mapexos-cache';

// L1FileExtension is the file extension for cached items.
export const L1FileExtension = '.cache';

// L1MetaExtension is the file extension for cache metadata.
export const L1MetaExtension = '.meta';

// L1SubdirDepth is the number of subdirectory levels for L1 cache.
// Depth 2 = 65,536 directories (256 × 256), ~3K files per dir for 200M assets.
export const L1SubdirDepth = 2;

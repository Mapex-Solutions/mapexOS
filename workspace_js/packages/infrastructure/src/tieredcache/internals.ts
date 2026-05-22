/**
 * TieredCache Internals
 * Same structure as workspace_go/packages/infrastructure/tieredcache/internals.go
 */

import * as fs from 'fs';
import * as path from 'path';
import * as crypto from 'crypto';
import type { Config, CacheEntry, CacheStats } from './types';
import { L1FileExtension, L1MetaExtension } from './constants';
import { ErrCacheMiss, ErrEntryExpired } from './errors';

/**
 * prefixKey adds the configured prefix to a cache key.
 */
export function prefixKey(keyPrefix: string, key: string): string {
	if (!keyPrefix) {
		return key;
	}
	return keyPrefix + ':' + key;
}

/**
 * hashKey creates a filesystem-safe hash of the cache key.
 * Used for L1 disk cache filenames.
 */
export function hashKey(key: string): string {
	const hash = crypto.createHash('sha256').update(key).digest('hex');
	return hash.substring(0, 32); // Use first 16 bytes (32 hex chars)
}

/**
 * l1FilePath returns the file path for a cached item in L1.
 * Uses 2-level subdirectory structure for better distribution with 200M+ assets.
 */
export function l1FilePath(l1Dir: string, key: string): string {
	const hash = hashKey(key);
	const level1 = hash.substring(0, 2);
	const level2 = hash.substring(2, 4);
	return path.join(l1Dir, level1, level2, hash + L1FileExtension);
}

/**
 * l1MetaPath returns the metadata file path for a cached item in L1.
 */
export function l1MetaPath(l1Dir: string, key: string): string {
	const hash = hashKey(key);
	const level1 = hash.substring(0, 2);
	const level2 = hash.substring(2, 4);
	return path.join(l1Dir, level1, level2, hash + L1MetaExtension);
}

/**
 * writeL1 writes data to L1 disk cache.
 */
export function writeL1(
	l1Dir: string,
	key: string,
	data: Buffer,
	ttlMs: number,
	stats: CacheStats,
): void {
	const filePath = l1FilePath(l1Dir, key);
	const metaPath = l1MetaPath(l1Dir, key);

	// Ensure directory exists
	const dir = path.dirname(filePath);
	fs.mkdirSync(dir, { recursive: true });

	// Write data file
	fs.writeFileSync(filePath, data);

	// Write metadata
	const entry: CacheEntry = {
		ExpiresAt: new Date(Date.now() + ttlMs),
		CreatedAt: new Date(),
		Size: data.length,
	};
	fs.writeFileSync(metaPath, JSON.stringify(entry));

	// Update L1 size stats
	stats.L1Size += data.length;
}

/**
 * readL1 reads data from L1 disk cache.
 */
export function readL1(l1Dir: string, key: string, stats: CacheStats): Buffer {
	const filePath = l1FilePath(l1Dir, key);
	const metaPath = l1MetaPath(l1Dir, key);

	// Check metadata for expiration
	let metaData: string;
	try {
		metaData = fs.readFileSync(metaPath, 'utf8');
	} catch (err) {
		if ((err as NodeJS.ErrnoException).code === 'ENOENT') {
			throw ErrCacheMiss;
		}
		throw err;
	}

	let entry: CacheEntry;
	try {
		entry = JSON.parse(metaData);
		entry.ExpiresAt = new Date(entry.ExpiresAt);
		entry.CreatedAt = new Date(entry.CreatedAt);
	} catch {
		// Corrupted metadata, delete files
		deleteL1Files(l1Dir, key, stats);
		throw ErrCacheMiss;
	}

	// Check expiration (lazy cleanup)
	if (Date.now() > entry.ExpiresAt.getTime()) {
		deleteL1Files(l1Dir, key, stats);
		stats.L1LazyExpired++;
		throw ErrEntryExpired;
	}

	// Read data file
	try {
		return fs.readFileSync(filePath);
	} catch (err) {
		if ((err as NodeJS.ErrnoException).code === 'ENOENT') {
			// Metadata exists but data doesn't - corrupted state
			deleteL1Files(l1Dir, key, stats);
			throw ErrCacheMiss;
		}
		throw err;
	}
}

/**
 * deleteL1Files removes both data and metadata files from L1.
 */
export function deleteL1Files(l1Dir: string, key: string, stats: CacheStats): void {
	const filePath = l1FilePath(l1Dir, key);
	const metaPath = l1MetaPath(l1Dir, key);

	// Get file size before deletion for stats
	try {
		const fileStats = fs.statSync(filePath);
		stats.L1Size -= fileStats.size;
	} catch {
		// File doesn't exist, ignore
	}

	try {
		fs.unlinkSync(filePath);
	} catch {
		// Ignore
	}
	try {
		fs.unlinkSync(metaPath);
	} catch {
		// Ignore
	}
}

/**
 * isL0Enabled returns true if L0 RAM cache is enabled.
 */
export function isL0Enabled(config: Config, l0Cache: Map<string, unknown> | null): boolean {
	return config.EnableL0 && l0Cache !== null;
}

/**
 * isL1Enabled returns true if L1 disk cache is enabled.
 */
export function isL1Enabled(config: Config, l1Dir: string): boolean {
	return config.EnableL1 && l1Dir !== '';
}

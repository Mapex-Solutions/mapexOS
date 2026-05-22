/**
 * TieredCache Factory
 * Same structure as workspace_go/packages/infrastructure/tieredcache/tieredcache.go
 */

import * as fs from 'fs';
import type { Config, LocalCacheLoader, L0Entry } from './types';
import {
	DefaultL0MaxSize,
	DefaultL0MaxItems,
	DefaultL0TTL,
	DefaultL1MaxSize,
	DefaultL1TTL,
	DefaultL1Dir,
} from './constants';
import { ErrInvalidConfig, ErrL0InitFailed, ErrL1InitFailed, ErrL1DirNotWritable } from './errors';
import { TieredCacheClient } from './methods';
import { HTTPClient } from '../httpclient';
import { getInfraLogger } from '../logger';

/**
 * New creates a new TieredCache client with optional L0 (RAM) and L1 (Disk).
 *
 * The cache follows a tiered architecture:
 *   - L0 (RAM): Ultra-fast in-memory cache using Map (~50µs)
 *   - L1 (Disk): Fast local NVMe/SSD storage (~500µs)
 *   - L2 (Remote): MinIO/S3 source of truth (configured via Config.L2Loader)
 *
 * Critical behavior:
 *   - L0 is enabled by default but can be disabled via config.EnableL0
 *   - L1 is enabled by default but can be disabled via config.EnableL1
 *   - L2 loader is configured via Config.EnableL2 and Config.L2Loader
 *   - At least one cache layer (L0, L1, or L2 loader) must be enabled
 *   - Cache follows LRU eviction with TTL support
 */
export function New(config: Partial<Config>): TieredCacheClient {
	const fullConfig = applyDefaults(config);

	const validationError = validateConfig(fullConfig);
	if (validationError) {
		throw new Error(`${ErrInvalidConfig.message}: ${validationError}`);
	}

	let l0Cache: Map<string, L0Entry> | null = null;
	let l1Dir = '';

	// Initialize L0 (RAM) cache if enabled
	if (fullConfig.EnableL0) {
		try {
			l0Cache = new Map<string, L0Entry>();
		} catch (err) {
			throw new Error(`${ErrL0InitFailed.message}: ${err}`);
		}
	}

	// Initialize L1 (Disk) cache if enabled
	if (fullConfig.EnableL1) {
		try {
			initL1Cache(fullConfig.L1Dir);
			l1Dir = fullConfig.L1Dir;
		} catch (err) {
			throw new Error(`${ErrL1InitFailed.message}: ${err}`);
		}
	}

	// Initialize L2 (Remote) loader if enabled
	let l2Loader: LocalCacheLoader | null = null;
	if (fullConfig.EnableL2) {
		if (!fullConfig.L2Loader) {
			throw new Error(`${ErrInvalidConfig.message}: L2Loader is required when EnableL2 is true`);
		}
		l2Loader = fullConfig.L2Loader;
	}

	// Initialize Fallback HTTP client if configured
	let fallbackClient: HTTPClient | null = null;
	let fallbackEndpoint = '';
	if (fullConfig.FallbackBaseURL) {
		fallbackClient = new HTTPClient({
			baseURL: fullConfig.FallbackBaseURL,
			apiKey: fullConfig.FallbackAPIKey,
			timeout: fullConfig.FallbackTimeout ?? 5000,
		});
		fallbackEndpoint = fullConfig.FallbackEndpoint ?? '';
	}
	const fallbackKeyTransformer = fullConfig.FallbackKeyTransformer ?? null;

	getInfraLogger().info(
		{ l0: fullConfig.EnableL0, l1: fullConfig.EnableL1, l2: !!fullConfig.EnableL2, fallback: !!fullConfig.FallbackBaseURL },
		'[INFRA:CACHE] Initialized',
	);

	return new TieredCacheClient(fullConfig, l0Cache, l1Dir, l2Loader, fallbackClient, fallbackEndpoint, fallbackKeyTransformer);
}

/**
 * validateConfig validates the cache configuration.
 */
function validateConfig(config: Config): string | null {
	if (config.L0MaxSize < 0) {
		return 'L0MaxSize cannot be negative';
	}
	if (config.L0MaxItems < 0) {
		return 'L0MaxItems cannot be negative';
	}
	if (config.EnableL1 && config.L1MaxSize < 0) {
		return 'L1MaxSize cannot be negative';
	}
	if (config.EnableL2 && !config.L2Loader) {
		return 'L2Loader is required when EnableL2 is true';
	}
	return null;
}

/**
 * applyDefaults applies default values to unset configuration fields.
 */
function applyDefaults(config: Partial<Config>): Config {
	return {
		EnableL0: config.EnableL0 ?? false,
		L0MaxSize: config.L0MaxSize ?? DefaultL0MaxSize,
		L0MaxItems: config.L0MaxItems ?? DefaultL0MaxItems,
		L0DefaultTTL: config.L0DefaultTTL ?? DefaultL0TTL,

		EnableL1: config.EnableL1 ?? false,
		L1Dir: config.L1Dir ?? DefaultL1Dir,
		L1MaxSize: config.L1MaxSize ?? DefaultL1MaxSize,
		L1DefaultTTL: config.L1DefaultTTL ?? DefaultL1TTL,

		// L2 (Remote) configuration
		EnableL2: config.EnableL2 ?? false,
		L2Loader: config.L2Loader,

		// Fallback configuration (optional)
		FallbackBaseURL: config.FallbackBaseURL,
		FallbackAPIKey: config.FallbackAPIKey,
		FallbackEndpoint: config.FallbackEndpoint,
		FallbackTimeout: config.FallbackTimeout,
		FallbackKeyTransformer: config.FallbackKeyTransformer,

		KeyPrefix: config.KeyPrefix ?? '',
		EnableMetrics: config.EnableMetrics ?? false,
	};
}

/**
 * initL1Cache initializes the L1 disk cache directory.
 */
function initL1Cache(dir: string): void {
	fs.mkdirSync(dir, { recursive: true });

	// Verify directory is writable
	const testFile = dir + '/.write_test';
	try {
		fs.writeFileSync(testFile, 'test');
		fs.unlinkSync(testFile);
	} catch {
		throw ErrL1DirNotWritable;
	}
}

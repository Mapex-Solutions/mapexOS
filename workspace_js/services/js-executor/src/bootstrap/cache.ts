import type { ConfigModule, Logger } from '@mapexos/microservices';

import { container } from 'tsyringe';

import { New as NewMinIO, NewTieredCache } from '@mapexos/infrastructure';

// InitCache registers MinIO clients and TieredCache instances (assets, templates, bytecode) in DI container.
/** Initializes TieredCache instances (asset, template, bytecode) and MinIO clients. */
export async function initCache(configModule: ConfigModule, logger: Logger) {

	// MinIO base configuration (shared across all buckets)
	const minioEndpoint = configModule.get('minio_endpoint') as string;
	const [minioHost, minioPortStr] = minioEndpoint.split(':');
	const minioPort = parseInt(minioPortStr) || 9000;

	const minioBaseConfig = {
		Endpoint: minioHost,
		Port: minioPort,
		AccessKeyID: configModule.get('minio_access_key') as string,
		SecretAccessKey: configModule.get('minio_secret_key') as string,
		UseSSL: configModule.get('minio_use_ssl') as boolean,
		Region: configModule.get('minio_region') as string,
	};

	// TieredCache shared configuration (L0 + L1)
	const serviceName = configModule.get('service_name') as string;
	const cacheL0MaxSize = configModule.get('cache_l0_max_size') as number;
	const cacheL0MaxItems = configModule.get('cache_l0_max_items') as number;
	const cacheL0TTL = (configModule.get('cache_l0_ttl_seconds') as number) * 1000; // Convert to ms
	const cacheL1Enabled = configModule.get('cache_l1_enabled') as boolean;
	const cacheL1BaseDir = configModule.get('cache_l1_dir') as string;
	const cacheL1Dir = `${cacheL1BaseDir}/${serviceName}`; // Include service name to avoid conflicts
	const cacheL1MaxSize = configModule.get('cache_l1_max_size') as number;
	const cacheL1TTL = (configModule.get('cache_l1_ttl_seconds') as number) * 1000; // Convert to ms

	// Fallback configuration - calls Assets Service when L2 misses
	const assetsServiceURL = configModule.get('assets_service_url') as string;
	const internalAPIKey = configModule.get('internal_api_key') as string;
	const fallbackTimeout = configModule.get('cache_fallback_timeout') as number;

	/** MinIO Client for Assets (L2 source of truth) */
	// No KeyPrefix - key already includes {orgId}/ from caller
	const minioAssetsClient = await NewMinIO({
		...minioBaseConfig,
		BucketName: configModule.get('minio_assets_bucket') as string,
	});

	container.register('MinIOAssetsClient', { useValue: minioAssetsClient });
	logger.info('[APP:BOOTSTRAP] MinIO Assets client initialized');

	/** MinIO Client for Templates (L2 source of truth) */
	// No KeyPrefix - key already includes {orgId}/ from caller
	const minioTemplatesClient = await NewMinIO({
		...minioBaseConfig,
		BucketName: configModule.get('minio_templates_bucket') as string,
	});

	container.register('MinIOTemplatesClient', { useValue: minioTemplatesClient });
	logger.info('[APP:BOOTSTRAP] MinIO Templates client initialized');

	/** TieredCache for Assets (L0 + L1 + L2 + Fallback) */
	// Key format: {orgId}/{assetUUID}
	const assetCache = NewTieredCache({
		EnableL0: true,
		L0MaxSize: cacheL0MaxSize,
		L0MaxItems: cacheL0MaxItems,
		L0DefaultTTL: cacheL0TTL,

		EnableL1: cacheL1Enabled,
		L1Dir: cacheL1Dir + '/assets',
		L1MaxSize: cacheL1MaxSize,
		L1DefaultTTL: cacheL1TTL,

		// L2 loader for assets (MinIO)
		// Key format: {orgId}/{assetUUID} → MinIO path: {orgId}/{assetUUID}.json
		EnableL2: true,
		L2Loader: async (key: string) => {
			try {
				const result = await minioAssetsClient.Get(key + '.json');
				return result.Data;
			} catch {
				return null;
			}
		},

		// No KeyPrefix - key already includes {orgId}/ from caller
		// Fallback configuration
		FallbackBaseURL: assetsServiceURL,
		FallbackAPIKey: internalAPIKey,
		FallbackEndpoint: '/internal/assets',
		FallbackTimeout: fallbackTimeout,
	});

	container.register('AssetCache', { useValue: assetCache });
	logger.info({ l1: cacheL1Enabled, fallback: assetsServiceURL }, '[APP:BOOTSTRAP] TieredCache (assets) initialized');

	/** TieredCache for Templates (L0 + L1 + L2 + Fallback) */
	// Key format: {orgId}/{templateId}
	const templateCache = NewTieredCache({
		// L0 RAM
		EnableL0: true,
		L0MaxSize: cacheL0MaxSize,
		L0MaxItems: cacheL0MaxItems,
		L0DefaultTTL: cacheL0TTL,
		
		// L1 Disk
		EnableL1: cacheL1Enabled,
		L1Dir: cacheL1Dir + '/templates',
		L1MaxSize: cacheL1MaxSize,
		L1DefaultTTL: cacheL1TTL,

		// L2 loader for templates (MinIO)
		// Key format: {orgId}/{templateId} → MinIO path: {orgId}/{templateId}.json
		EnableL2: true,
		L2Loader: async (key: string) => {
			try {
				const result = await minioTemplatesClient.Get(key + '.json');
				return result.Data;
			} catch {
				return null;
			}
		},

		// No KeyPrefix - key already includes {orgId}/ from caller
		// Fallback configuration
		FallbackBaseURL: assetsServiceURL,
		FallbackAPIKey: internalAPIKey,
		FallbackEndpoint: '/internal/templates',
		FallbackTimeout: fallbackTimeout,
	});

	container.register('TemplateCache', { useValue: templateCache });
	logger.info({ l1: cacheL1Enabled, fallback: assetsServiceURL }, '[APP:BOOTSTRAP] TieredCache (templates) initialized');

	/** MinIO Client for Bytecode (L2 source of truth for compiled scripts) */
	// Bytecode is shared across pods for horizontal scaling benefits
	const minioBytecodeClient = await NewMinIO({
		...minioBaseConfig,
		BucketName: configModule.get('minio_bytecode_bucket') as string,
		KeyPrefix: 'bytecode',
	});
	container.register('MinIOBytecodeClient', { useValue: minioBytecodeClient });
	logger.info('[APP:BOOTSTRAP] MinIO Bytecode client initialized');

	/** TieredCache for Bytecode (L1 + L2) */
	// L0 (RAM) is SKIPPED by default since ScriptRegistry handles in-memory caching
	// Scripts rarely change, so L2 (MinIO) benefits horizontal scaling:
	// - Pod A compiles and stores in MinIO
	// - Pod B, C, D can reuse without recompiling
	const skipBytecodeL0 = configModule.get('cache_bytecode_skip_l0') as boolean;
	const bytecodeCache = NewTieredCache({
		
		// L0 RAM
		EnableL0: !skipBytecodeL0, // Skip L0 since ScriptRegistry caches Script Objects in RAM
		L0MaxSize: cacheL0MaxSize,
		L0MaxItems: cacheL0MaxItems * 3, // 3 scripts per asset (decode, validation, transform)
		L0DefaultTTL: cacheL0TTL,
		
		// L1 Disk
		EnableL1: cacheL1Enabled,
		L1Dir: cacheL1Dir + '/bytecode',
		L1MaxSize: cacheL1MaxSize,
		L1DefaultTTL: cacheL1TTL,
		KeyPrefix: 'bytecode',

		// L2 loader for bytecode (MinIO)
		EnableL2: true,
		L2Loader: async (key: string) => {
			try {
				const result = await minioBytecodeClient.Get(key + '.bin');
				return result.Data;
			} catch {
				return null;
			}
		},
	});

	container.register('BytecodeCache', { useValue: bytecodeCache });
	logger.info({ l1: cacheL1Enabled }, '[APP:BOOTSTRAP] TieredCache (bytecode) initialized');
}

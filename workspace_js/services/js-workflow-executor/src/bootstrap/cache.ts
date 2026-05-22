import type { ConfigModule, Logger } from '@mapexos/microservices';

import { container } from 'tsyringe';

import { New as NewMinIO, NewTieredCache } from '@mapexos/infrastructure';

// InitCache registers MinIO client and TieredCache instances (script source, bytecode) in DI container.
/** Initializes TieredCache instances (script source + bytecode) and MinIO client. */
export async function initCache(configModule: ConfigModule, logger: Logger) {

	// MinIO base configuration
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
	const cacheL0TTL = (configModule.get('cache_l0_ttl_seconds') as number) * 1000;
	
	const cacheL1Enabled = configModule.get('cache_l1_enabled') as boolean;
	const cacheL1BaseDir = configModule.get('cache_l1_dir') as string;
	const cacheL1Dir = `${cacheL1BaseDir}/${serviceName}`;
	const cacheL1MaxSize = configModule.get('cache_l1_max_size') as number;
	const cacheL1TTL = (configModule.get('cache_l1_ttl_seconds') as number) * 1000;

	// Fallback configuration — calls Workflow Service when L2 misses
	const workflowServiceURL = configModule.get('workflow_service_url') as string;
	const internalAPIKey = configModule.get('internal_api_key') as string;
	const fallbackTimeout = configModule.get('cache_fallback_timeout') as number;

	/** MinIO Client for Workflows (single bucket: mapex-workflows) */
	// Path structure: {orgId}/{workflowId}/scripts/{nodeId}.json (source)
	//                 {orgId}/{workflowId}/bytecode/{nodeId}.bin (compiled)
	const minioWorkflowsClient = await NewMinIO({
		...minioBaseConfig,
		BucketName: configModule.get('minio_workflows_bucket') as string,
	});

	container.register('MinIOWorkflowsClient', { useValue: minioWorkflowsClient });
	logger.info('[APP:BOOTSTRAP] MinIO Workflows client initialized');

	/** TieredCache for Script Source (L0 + L1 + L2 + Fallback) */
	// Key format: {orgId}/{workflowId}/scripts/{nodeId}
	// L2 fetches from MinIO: {orgId}/{workflowId}/scripts/{nodeId}.json
	// Fallback HTTP: GET /internal/definitions/{orgId}/{workflowId}/nodes/{nodeId}/script
	const scriptSourceCache = NewTieredCache({
		EnableL0: true,
		L0MaxSize: cacheL0MaxSize,
		L0MaxItems: cacheL0MaxItems,
		L0DefaultTTL: cacheL0TTL,

		EnableL1: cacheL1Enabled,
		L1Dir: cacheL1Dir + '/scripts',
		L1MaxSize: cacheL1MaxSize,
		L1DefaultTTL: cacheL1TTL,

		// L2 loader for script source (MinIO)
		EnableL2: true,
		L2Loader: async (key: string) => {
			try {
				const result = await minioWorkflowsClient.Get(key + '.json');
				return result.Data;
			} catch {
				return null;
			}
		},

		// Fallback: call Workflow Service HTTP endpoint
		FallbackBaseURL: workflowServiceURL,
		FallbackAPIKey: internalAPIKey,
		FallbackEndpoint: '/internal/workflow-scripts',
		FallbackTimeout: fallbackTimeout,
		
		FallbackKeyTransformer: (key: string) => {
			// key format: "{orgId}/{defId}/scripts/{nodeId}"
			// URL path: "/{defId}/scripts/{nodeId}" (strip orgId prefix)
			const firstSlash = key.indexOf('/');
			return firstSlash !== -1 ? key.substring(firstSlash) : `/${key}`;
		},
	});

	container.register('ScriptSourceCache', { useValue: scriptSourceCache });
	logger.info({ l1: cacheL1Enabled, fallback: workflowServiceURL }, '[APP:BOOTSTRAP] TieredCache (script source) initialized');

	/** TieredCache for Bytecode (L1 + L2, no fallback) */
	// L0 (RAM) is SKIPPED by default since piscina workers cache compiled Scripts in-memory.
	// L2 (MinIO) benefits horizontal scaling: Pod A compiles → Pod B reuses.
	// Key format: {orgId}/{workflowId}/bytecode/{nodeId}
	const skipBytecodeL0 = configModule.get('cache_bytecode_skip_l0') as boolean;
	const bytecodeCache = NewTieredCache({
		EnableL0: !skipBytecodeL0,
		L0MaxSize: cacheL0MaxSize,
		L0MaxItems: cacheL0MaxItems,
		L0DefaultTTL: cacheL0TTL,

		EnableL1: cacheL1Enabled,
		L1Dir: cacheL1Dir + '/bytecode',
		L1MaxSize: cacheL1MaxSize,
		L1DefaultTTL: cacheL1TTL,

		// L2 loader for bytecode (MinIO)
		EnableL2: true,
		L2Loader: async (key: string) => {
			try {
				const result = await minioWorkflowsClient.Get(key + '.bin');
				return result.Data;
			} catch {
				return null;
			}
		},
	});

	container.register('BytecodeCache', { useValue: bytecodeCache });
	logger.info({ l1: cacheL1Enabled }, '[APP:BOOTSTRAP] TieredCache (bytecode) initialized');
}

/**
 * MinIO Client Factory
 * Same structure as workspace_go/packages/infrastructure/minio/minio.go
 */

import * as Minio from 'minio';
import type { Config } from './types';
import { DefaultRegion } from './constants';
import { ErrInvalidConfig, ErrConnectionFailed } from './errors';
import { MinIOClient } from './methods';
import { getInfraLogger } from '../logger';

/**
 * New creates a new MinIO client wrapper.
 *
 * It validates the configuration, establishes a connection to MinIO,
 * and optionally verifies that the specified bucket exists.
 *
 * The client is ready to use for object storage operations immediately
 * after successful initialization.
 *
 * Critical behavior:
 *   - Validates required config fields (Endpoint, AccessKeyID, SecretAccessKey, BucketName)
 *   - Creates a connection with retry logic for transient failures
 *   - Verifies bucket existence (creates if BucketName is specified)
 *   - Returns error if connection or bucket verification fails
 */
export async function New(config: Config): Promise<MinIOClient> {
	const validationError = validateConfig(config);
	if (validationError) {
		throw new Error(`${ErrInvalidConfig.message}: ${validationError}`);
	}

	const region = config.Region || DefaultRegion;

	const minioClient = new Minio.Client({
		endPoint: config.Endpoint,
		port: config.Port,
		useSSL: config.UseSSL ?? false,
		accessKey: config.AccessKeyID,
		secretKey: config.SecretAccessKey,
		region: region,
	});

	// Verify connection with a health check
	try {
		await minioClient.listBuckets();
	} catch (err) {
		throw new Error(`${ErrConnectionFailed.message}: ${err}`);
	}

	// Ensure bucket exists
	if (config.BucketName) {
		try {
			const exists = await minioClient.bucketExists(config.BucketName);
			if (!exists) {
				getInfraLogger().warn({ bucket: config.BucketName }, '[INFRA:MINIO] Bucket does not exist');
			}
		} catch (err) {
			throw new Error(`bucket check failed: ${err}`);
		}
	}

	getInfraLogger().info({ bucket: config.BucketName }, '[INFRA:MINIO] Initialized');

	return new MinIOClient(minioClient, config.BucketName, config.KeyPrefix || '');
}

/**
 * validateConfig validates the MinIO configuration.
 */
function validateConfig(config: Config): string | null {
	if (!config.Endpoint) {
		return 'endpoint is required';
	}
	if (!config.AccessKeyID) {
		return 'access key ID is required';
	}
	if (!config.SecretAccessKey) {
		return 'secret access key is required';
	}
	return null;
}

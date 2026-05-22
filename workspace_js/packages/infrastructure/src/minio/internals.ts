/**
 * MinIO Internals
 * Same structure as workspace_go/packages/infrastructure/minio/internals.go
 */

import type { BucketItemStat } from 'minio';
import type { PutOptions, ObjectInfo } from './types';
import { ContentTypeBinary } from './constants';

/**
 * prefixKey adds the configured prefix to an object key.
 */
export function prefixKey(keyPrefix: string, key: string): string {
	if (!keyPrefix) {
		return key;
	}
	return keyPrefix + '/' + key;
}

/**
 * buildPutOptions converts PutOptions to minio metadata object.
 */
export function buildPutOptions(opts?: PutOptions): Record<string, string> {
	const metadata: Record<string, string> = {};

	if (!opts) {
		metadata['Content-Type'] = ContentTypeBinary;
		return metadata;
	}

	metadata['Content-Type'] = opts.ContentType || ContentTypeBinary;

	if (opts.CacheControl) {
		metadata['Cache-Control'] = opts.CacheControl;
	}

	if (opts.UserMetadata) {
		for (const [key, value] of Object.entries(opts.UserMetadata)) {
			metadata[`X-Amz-Meta-${key}`] = value;
		}
	}

	return metadata;
}

/**
 * convertObjectInfo converts minio BucketItemStat to our ObjectInfo type.
 */
export function convertObjectInfo(info: BucketItemStat, key: string): ObjectInfo {
	return {
		Key: key,
		Size: info.size,
		LastModified: info.lastModified,
		ETag: info.etag,
		ContentType: info.metaData?.['content-type'],
	};
}

/**
 * isNotFoundError checks if the error indicates the object was not found.
 */
export function isNotFoundError(err: unknown): boolean {
	if (!err) {
		return false;
	}

	if (err instanceof Error) {
		const message = err.message.toLowerCase();
		return (
			message.includes('not found') ||
			message.includes('nosuchkey') ||
			message.includes('does not exist') ||
			message.includes('no such key')
		);
	}

	// Check for minio error response
	const errorObj = err as { code?: string };
	if (errorObj.code) {
		return errorObj.code === 'NoSuchKey' || errorObj.code === 'NoSuchBucket';
	}

	return false;
}

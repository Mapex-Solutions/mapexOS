/**
 * MinIO Methods
 * Same structure as workspace_go/packages/infrastructure/minio/methods.go
 */

import type { Readable } from 'stream';
import type { Client as MinioNativeClient } from 'minio';
import type { GetResult, ObjectInfo, PutOptions, ListOptions, StreamReader } from './types';
import { ContentTypeJSON } from './constants';
import {
	ErrEmptyKey,
	ErrNilData,
	ErrObjectNotFound,
	ErrUploadFailed,
	ErrDownloadFailed,
} from './errors';
import { prefixKey, buildPutOptions, convertObjectInfo, isNotFoundError } from './internals';

/**
 * MinIOClient wraps the minio-go client providing a simplified interface
 * for MapexOS services. Follows the hexagonal architecture pattern.
 */
export class MinIOClient {
	private readonly client: MinioNativeClient;
	private readonly bucketName: string;
	private readonly keyPrefix: string;

	constructor(client: MinioNativeClient, bucketName: string, keyPrefix: string = '') {
		this.client = client;
		this.bucketName = bucketName;
		this.keyPrefix = keyPrefix;
	}

	/**
	 * Put uploads an object to MinIO.
	 *
	 * It accepts a key (path), data as bytes, and optional PutOptions.
	 * The key is automatically prefixed using the MinIOClient's keyPrefix.
	 */
	async Put(key: string, data: Buffer, opts?: PutOptions): Promise<void> {
		if (!key) {
			throw ErrEmptyKey;
		}
		if (!data) {
			throw ErrNilData;
		}

		const prefixed = prefixKey(this.keyPrefix, key);
		const metadata = buildPutOptions(opts);

		try {
			await this.client.putObject(this.bucketName, prefixed, data, data.length, metadata);
		} catch (err) {
			throw new Error(`${ErrUploadFailed.message}: ${err}`);
		}
	}

	/**
	 * PutStream uploads an object from a stream (for large objects).
	 */
	async PutStream(
		key: string,
		stream: Readable,
		size: number,
		opts?: PutOptions,
	): Promise<void> {
		if (!key) {
			throw ErrEmptyKey;
		}

		const prefixed = prefixKey(this.keyPrefix, key);
		const metadata = buildPutOptions(opts);

		try {
			await this.client.putObject(this.bucketName, prefixed, stream, size, metadata);
		} catch (err) {
			throw new Error(`${ErrUploadFailed.message}: ${err}`);
		}
	}

	/**
	 * Get retrieves an object from MinIO and returns its content as bytes.
	 */
	async Get(key: string): Promise<GetResult> {
		if (!key) {
			throw ErrEmptyKey;
		}

		const prefixed = prefixKey(this.keyPrefix, key);

		try {
			// Get object info
			const stat = await this.client.statObject(this.bucketName, prefixed);

			// Get object data
			const stream = await this.client.getObject(this.bucketName, prefixed);
			const chunks: Buffer[] = [];

			for await (const chunk of stream) {
				chunks.push(chunk as Buffer);
			}

			const data = Buffer.concat(chunks);

			return {
				Data: data,
				ContentType: stat.metaData?.['content-type'] || 'application/octet-stream',
				Size: stat.size,
				LastModified: stat.lastModified,
				ETag: stat.etag,
				UserMetadata: stat.metaData,
			};
		} catch (err) {
			if (isNotFoundError(err)) {
				throw ErrObjectNotFound;
			}
			throw new Error(`${ErrDownloadFailed.message}: ${err}`);
		}
	}

	/**
	 * GetStream retrieves an object as a stream reader (for large objects).
	 */
	async GetStream(key: string): Promise<StreamReader> {
		if (!key) {
			throw ErrEmptyKey;
		}

		const prefixed = prefixKey(this.keyPrefix, key);

		try {
			const stat = await this.client.statObject(this.bucketName, prefixed);
			const stream = await this.client.getObject(this.bucketName, prefixed);

			return {
				Stream: stream,
				Info: convertObjectInfo(stat, prefixed),
			};
		} catch (err) {
			if (isNotFoundError(err)) {
				throw ErrObjectNotFound;
			}
			throw new Error(`${ErrDownloadFailed.message}: ${err}`);
		}
	}

	/**
	 * Delete removes an object from MinIO.
	 * Does not return an error if the object does not exist.
	 */
	async Delete(key: string): Promise<void> {
		if (!key) {
			throw ErrEmptyKey;
		}

		const prefixed = prefixKey(this.keyPrefix, key);

		try {
			await this.client.removeObject(this.bucketName, prefixed);
		} catch (err) {
			// Ignore not found errors on delete
			if (!isNotFoundError(err)) {
				throw err;
			}
		}
	}

	/**
	 * Exists checks if an object exists in MinIO.
	 */
	async Exists(key: string): Promise<boolean> {
		if (!key) {
			throw ErrEmptyKey;
		}

		const prefixed = prefixKey(this.keyPrefix, key);

		try {
			await this.client.statObject(this.bucketName, prefixed);
			return true;
		} catch (err) {
			if (isNotFoundError(err)) {
				return false;
			}
			throw err;
		}
	}

	/**
	 * Stat returns metadata about an object without downloading it.
	 */
	async Stat(key: string): Promise<ObjectInfo> {
		if (!key) {
			throw ErrEmptyKey;
		}

		const prefixed = prefixKey(this.keyPrefix, key);

		try {
			const stat = await this.client.statObject(this.bucketName, prefixed);
			return convertObjectInfo(stat, prefixed);
		} catch (err) {
			if (isNotFoundError(err)) {
				throw ErrObjectNotFound;
			}
			throw err;
		}
	}

	/**
	 * List returns a list of objects matching the given options.
	 */
	async List(opts?: ListOptions): Promise<ObjectInfo[]> {
		const prefix = opts?.Prefix ? prefixKey(this.keyPrefix, opts.Prefix) : this.keyPrefix;
		const recursive = opts?.Recursive ?? false;

		const objects: ObjectInfo[] = [];
		const stream = this.client.listObjectsV2(this.bucketName, prefix, recursive);

		for await (const obj of stream) {
			if (obj.name) {
				objects.push({
					Key: obj.name,
					Size: obj.size,
					LastModified: obj.lastModified,
					ETag: obj.etag,
				});

				if (opts?.MaxKeys && objects.length >= opts.MaxKeys) {
					break;
				}
			}
		}

		return objects;
	}

	/**
	 * Copy copies an object to a new key within the same bucket.
	 */
	async Copy(srcKey: string, dstKey: string): Promise<void> {
		if (!srcKey || !dstKey) {
			throw ErrEmptyKey;
		}

		const srcPrefixed = prefixKey(this.keyPrefix, srcKey);
		const dstPrefixed = prefixKey(this.keyPrefix, dstKey);

		try {
			await this.client.copyObject(
				this.bucketName,
				dstPrefixed,
				`/${this.bucketName}/${srcPrefixed}`,
			);
		} catch (err) {
			throw new Error(`failed to copy object: ${err}`);
		}
	}

	/**
	 * PutJSON is a convenience method for uploading JSON data.
	 */
	async PutJSON(key: string, data: Buffer): Promise<void> {
		return this.Put(key, data, { ContentType: ContentTypeJSON });
	}

	/**
	 * PutWithMetadata is a convenience method for uploading data with metadata.
	 */
	async PutWithMetadata(
		key: string,
		data: Buffer,
		contentType: string,
		metadata: Record<string, string>,
	): Promise<void> {
		return this.Put(key, data, {
			ContentType: contentType,
			UserMetadata: metadata,
		});
	}

	/**
	 * Ping checks the MinIO connection by listing buckets.
	 *
	 * @throws If the connection is down or unreachable
	 */
	async Ping(): Promise<void> {
		await this.client.listBuckets();
	}

	/**
	 * GetRawClient returns the underlying minio-go client.
	 * Use with caution - prefer using the wrapper methods.
	 */
	GetRawClient(): MinioNativeClient {
		return this.client;
	}

	/**
	 * GetBucketName returns the configured bucket name.
	 */
	GetBucketName(): string {
		return this.bucketName;
	}
}

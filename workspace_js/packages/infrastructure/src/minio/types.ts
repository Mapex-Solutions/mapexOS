/**
 * MinIO Types
 * Same structure as workspace_go/packages/infrastructure/minio/types.go
 */

import type { Readable } from 'stream';

/**
 * MinIOClient wraps the minio client providing a simplified interface
 * for MapexOS services. Follows the hexagonal architecture pattern.
 */
export interface MinIOClientInterface {
	bucketName: string;
	keyPrefix: string;
}

/**
 * Config holds the MinIO connection configuration.
 * Used for dependency injection via tsyringe.
 */
export interface Config {
	// Endpoint is the MinIO server address (e.g., "localhost:9000")
	Endpoint: string;

	// AccessKeyID is the MinIO access key (like AWS Access Key)
	AccessKeyID: string;

	// SecretAccessKey is the MinIO secret key (like AWS Secret Key)
	SecretAccessKey: string;

	// BucketName is the default bucket for this client instance
	BucketName: string;

	// KeyPrefix is prepended to all object keys (optional)
	KeyPrefix?: string;

	// UseSSL enables HTTPS connections
	UseSSL?: boolean;

	// Region is the S3 region (optional, defaults to "us-east-1")
	Region?: string;

	// Port is the MinIO server port (optional)
	Port?: number;
}

/**
 * PutOptions contains options for uploading objects.
 */
export interface PutOptions {
	// ContentType is the MIME type of the object (e.g., "application/json")
	ContentType?: string;

	// UserMetadata are custom key-value pairs stored with the object
	UserMetadata?: Record<string, string>;

	// CacheControl sets the Cache-Control header for the object
	CacheControl?: string;

	// Expires sets the expiration time for the object
	Expires?: Date;
}

/**
 * GetResult contains the result of a Get operation.
 */
export interface GetResult {
	// Data is the object content as bytes
	Data: Buffer;

	// ContentType is the MIME type of the object
	ContentType: string;

	// Size is the size of the object in bytes
	Size: number;

	// LastModified is the last modification time
	LastModified: Date;

	// ETag is the entity tag (hash) of the object
	ETag: string;

	// UserMetadata are custom key-value pairs stored with the object
	UserMetadata?: Record<string, string>;
}

/**
 * ObjectInfo contains metadata about an object.
 */
export interface ObjectInfo {
	// Key is the object key (path)
	Key: string;

	// Size is the size in bytes
	Size: number;

	// LastModified is the last modification time
	LastModified: Date;

	// ETag is the entity tag (hash)
	ETag: string;

	// ContentType is the MIME type
	ContentType?: string;
}

/**
 * ListOptions contains options for listing objects.
 */
export interface ListOptions {
	// Prefix filters objects by key prefix
	Prefix?: string;

	// Recursive lists objects recursively (default: false)
	Recursive?: boolean;

	// MaxKeys limits the number of objects returned (0 = no limit)
	MaxKeys?: number;
}

/**
 * StreamReader wraps ReadableStream for streaming large objects.
 */
export interface StreamReader {
	Stream: Readable;
	Info: ObjectInfo;
}

export interface LakeHouseCredentials {
	// AWS S3 / MinIO
	accessKey?: string;
	secretKey?: string;
	region?: string;
	bucket?: string;
	endpoint?: string;
	useSSL?: boolean;
	tags?: Record<string, string>;

	// Azure Blob Storage
	accountName?: string;
	accountKey?: string;
	containerName?: string;

	// Google Cloud Storage
	projectId?: string;
	keyFile?: string;
}
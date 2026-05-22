import type { Provider } from '../interfaces';

export const PROVIDERS: Provider[] = [
  {
    id: 'aws-s3',
    name: 'Amazon S3',
    description: 'AWS scalable storage',
    hasInfo: true,
    infoText: 'Native service of Amazon Web Services',
	  icon: 'mdi-aws',
    iconColor: 'orange',
  },
  {
    id: 'azure-blob',
    name: 'Azure Blob Storage',
    description: 'Microsoft storage service',
    icon: 'mdi-microsoft-azure',
    iconColor: 'blue',
  },
  {
    id: 'gcp-storage',
    name: 'Google Cloud Storage',
    description: 'Google scalable storage',
    icon: 'mdi-google-cloud',
    iconColor: 'red',
  },
  {
    id: 'minio',
    name: 'MinIO',
    description: 'High-performance S3-compatible storage',
    hasInfo: true,
    infoText: 'MinIO uses the S3-compatible API',
	  icon: 'mdi-database',
    iconColor: 'purple',
  },
];


export const DEFAULT_GCP_DATA = {
	projectId: '',
	keyFile: '',
	bucket: '',
	region: '',
};

export const DEFAULT_AWS_DATA = {
	accessKey: '',
	secretKey: '',
	bucket: '',
	region: 'us-east-1',
};

export const DEFAULT_AZURE_DATA = {
	accountName: '',
	accountKey: '',
	containerName: '',
	endpoint: '',      // blob uses a string endpoint
};

export const DEFAULT_MINIO_DATA = {
	accessKey: '',
	secretKey: '',
	bucket: '',
	region: '',        // optional, but if present must be string
	endpoint: '',      // now always a string
	useSSL: true,      // always a boolean
};


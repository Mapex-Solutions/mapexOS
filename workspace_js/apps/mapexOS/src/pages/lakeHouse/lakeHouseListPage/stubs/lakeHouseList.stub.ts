export const LAKE_HOUSE_LIST_STUB = [
  {
    type: 'aws-s3' as const,
    name: 'Amazon S3',
    description: 'Amazon S3 data lake configured for daily partitioned uploads.',
    status: true,
    pathConfig: {
      maxFileSize: 100,
      partitions: ['year', 'month', 'day', 'hour', 'asset_id', 'asset_type']
    },
    credentials: {
      bucket: 'my-amazon-s3-bucket',
      region: 'us-east-1'
    },
    frequency: {
      interval: 1,
      type: 'day',
      time: '00:00'
    },
    created: new Date().toDateString(),
  },
  {
    type: 'azure-blob' as const,
    name: 'Azure Blob Storage',
    description: 'Azure Blob Storage data lake with scalable object storage.',
    status: true,
    pathConfig: {
      maxFileSize: 100,
      partitions: ['year', 'month', 'day', 'hour', 'asset_id', 'asset_type']
    },
    credentials: {
      bucket: 'my-azure-blob-container',
      region: 'eastus'
    },
    frequency: {
      interval: 1,
      type: 'day',
      time: '00:00'
    },
    created: new Date().toDateString(),
  },
  {
    type: 'gcp-storage' as const,
    name: 'Google Cloud Storage',
    description: 'Google Cloud Storage data lake for high-throughput processing.',
    status: true,
    pathConfig: {
      maxFileSize: 100,
      partitions: ['year', 'month', 'day', 'hour', 'asset_id', 'asset_type']
    },
    credentials: {
      bucket: 'my-gcp-storage-bucket',
      region: 'us-central1'
    },
    frequency: {
      interval: 1,
      type: 'day',
      time: '00:00'
    },
    created: new Date().toDateString(),
  },
  {
    type: 'minio' as const,
    name: 'MinIO',
    description: 'MinIO instance for local S3-compatible object storage.',
    status: true,
    pathConfig: {
      maxFileSize: 100,
      partitions: ['year', 'month', 'day', 'hour', 'asset_id', 'asset_type']
    },
    credentials: {
      bucket: 'my-minio-bucket',
      region: 'minio-local'
    },
    frequency: {
      interval: 1,
      type: 'day',
      time: '00:00'
    },
    created: new Date().toDateString(),
  },
];

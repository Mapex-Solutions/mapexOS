import type { LakeHouseFrequency } from '../lakeHouseFrequency/interfaces';
import type { LakeHouseCredentials } from '../lakeHouseCredentials/interfaces';

interface LakeHousePathConfig {
	basePath: string;
	partitions: string[];
	compression?: 'gzip' | 'snappy' | 'lz4' | 'none';
	filePrefix?: string;
	maxFileSize?: number;
}


export interface LakeHouseConfigProps {
	name: string;
	description: string;
	type: 'aws-s3' | 'azure-blob' | 'gcp-storage' | 'minio';
	status: boolean;
	tenantId?: string;
	credentials: LakeHouseCredentials;
	frequency: LakeHouseFrequency;
	pathConfig: LakeHousePathConfig;
}
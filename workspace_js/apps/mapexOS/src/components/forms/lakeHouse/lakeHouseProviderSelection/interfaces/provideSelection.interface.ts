import type { LakeHouseConfigProps } from '@components/forms/lakeHouse';

export interface Provider {
	id: 'aws-s3' | 'azure-blob' | 'gcp-storage' | 'minio';
	name: string;
	description: string;
	hasInfo?: boolean;
	infoText?: string;
	icon: string;
	iconColor: string;
}

export interface Emits {
  (e: 'update:modelValue', value: LakeHouseConfigProps): void;
}
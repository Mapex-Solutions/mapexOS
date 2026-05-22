/**
 * LakeHouseSelectorDrawer Interfaces
 */

/**
 * Data Lake type
 */
export interface LakeHouseItem {
  /** Unique identifier */
  id?: string;

  /** Data lake type (aws-s3, azure-blob, gcp-storage, minio) */
  type: 'aws-s3' | 'azure-blob' | 'gcp-storage' | 'minio';

  /** Display name */
  name: string;

  /** Description */
  description?: string;

  /** Status (active/inactive) */
  status: boolean;

  /** Path configuration */
  pathConfig?: {
    maxFileSize: number;
    partitions: string[];
  };

  /** Credentials configuration */
  credentials?: {
    bucket: string;
    region: string;
  };

  /** Frequency configuration */
  frequency?: {
    interval: number;
    type: string;
    time: string;
  };

  /** Created date */
  created?: string;
}

/**
 * Props for LakeHouseSelectorDrawer component
 */
export interface LakeHouseSelectorDrawerProps {
  /** Whether the drawer is open */
  modelValue: boolean;

  /** Pre-selected data lake ID */
  selectedLakeHouseId?: string | null;
}

/**
 * Emits for LakeHouseSelectorDrawer component
 */
export interface LakeHouseSelectorDrawerEmits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'select', lakeHouse: LakeHouseItem): void;
  (e: 'cancel'): void;
}

import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import LakeHouseProviderSelection from './LakeHouseProviderSelection.vue';

vi.mock('@components/tooltips', () => ({
  AppTooltip: { name: 'AppTooltip', template: '<div />' },
}));

vi.mock('./constants', () => ({
  PROVIDERS: [
    { id: 'aws-s3', name: 'AWS S3', description: 'Amazon', icon: 'cloud', iconColor: 'orange', hasInfo: false },
    { id: 'minio', name: 'MinIO', description: 'Open Source', icon: 'storage', iconColor: 'red', hasInfo: true, infoText: 'S3 compatible' },
    { id: 'azure-blob', name: 'Azure Blob', description: 'Microsoft', icon: 'cloud', iconColor: 'blue', hasInfo: false },
    { id: 'gcp-storage', name: 'GCP Storage', description: 'Google', icon: 'cloud', iconColor: 'blue', hasInfo: false },
  ],
  DEFAULT_AWS_DATA: { accessKey: '', secretKey: '', region: '', bucket: '' },
  DEFAULT_MINIO_DATA: { accessKey: '', secretKey: '', region: '', bucket: '', endpoint: '', useSSL: true },
  DEFAULT_AZURE_DATA: { accountName: '', accountKey: '', containerName: '' },
  DEFAULT_GCP_DATA: { projectId: '', region: '', keyFile: '', bucket: '' },
}));

const BASE_MODEL = {
  name: '',
  status: true,
  description: '',
  type: '' as any,
  credentials: {} as any,
};

describe('LakeHouseProviderSelection', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(LakeHouseProviderSelection, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('selectProvider updates type to aws-s3', () => {
    const wrapper = mountWithPlugins(LakeHouseProviderSelection, {
      props: { modelValue: BASE_MODEL },
    });
    wrapper.vm.selectProvider({ id: 'aws-s3', name: 'AWS S3', description: '', icon: '', iconColor: '', hasInfo: false });
    expect(wrapper.vm.modelRef.type).toBe('aws-s3');
  });

  it('selectProvider resets credentials for the provider', () => {
    const wrapper = mountWithPlugins(LakeHouseProviderSelection, {
      props: { modelValue: { ...BASE_MODEL, type: 'aws-s3', credentials: { accessKey: 'old' } } },
    });
    wrapper.vm.selectProvider({ id: 'minio', name: 'MinIO', description: '', icon: '', iconColor: '', hasInfo: false });
    expect(wrapper.vm.modelRef.type).toBe('minio');
    expect(wrapper.vm.modelRef.credentials.accessKey).toBe('');
  });

  it('getDefaultCredentials returns correct defaults for azure', () => {
    const wrapper = mountWithPlugins(LakeHouseProviderSelection, {
      props: { modelValue: BASE_MODEL },
    });
    const creds = wrapper.vm.getDefaultCredentials('azure-blob');
    expect(creds).toHaveProperty('accountName');
  });

  it('getDefaultCredentials returns empty object for unknown provider', () => {
    const wrapper = mountWithPlugins(LakeHouseProviderSelection, {
      props: { modelValue: BASE_MODEL },
    });
    const creds = wrapper.vm.getDefaultCredentials('unknown');
    expect(Object.keys(creds)).toHaveLength(0);
  });
});

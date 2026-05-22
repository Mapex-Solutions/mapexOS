import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import LakeHouseCredentials from './LakeHouseCredentials.vue';

const BASE_MODEL = {
  type: 'aws-s3' as const,
  credentials: {
    accessKey: '',
    secretKey: '',
    region: '',
    bucket: '',
  },
};

describe('LakeHouseCredentials', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(LakeHouseCredentials, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes providerType from model', () => {
    const wrapper = mountWithPlugins(LakeHouseCredentials, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.providerType).toBe('aws-s3');
  });

  it('computes isS3Compatible as true for aws-s3', () => {
    const wrapper = mountWithPlugins(LakeHouseCredentials, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isS3Compatible).toBe(true);
  });

  it('computes isS3Compatible as true for minio', () => {
    const wrapper = mountWithPlugins(LakeHouseCredentials, {
      props: { modelValue: { ...BASE_MODEL, type: 'minio' } },
    });
    expect(wrapper.vm.isS3Compatible).toBe(true);
  });

  it('computes isS3Compatible as false for azure-blob', () => {
    const wrapper = mountWithPlugins(LakeHouseCredentials, {
      props: {
        modelValue: {
          type: 'azure-blob',
          credentials: { accountName: '', accountKey: '', containerName: '' },
        },
      },
    });
    expect(wrapper.vm.isS3Compatible).toBe(false);
  });

  it('computes canTestConnection as false when credentials empty', () => {
    const wrapper = mountWithPlugins(LakeHouseCredentials, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.canTestConnection).toBe(false);
  });

  it('computes canTestConnection as true when aws credentials filled', () => {
    const wrapper = mountWithPlugins(LakeHouseCredentials, {
      props: {
        modelValue: {
          type: 'aws-s3',
          credentials: {
            accessKey: 'AK',
            secretKey: 'SK',
            region: 'us-east-1',
            bucket: 'my-bucket',
          },
        },
      },
    });
    expect(wrapper.vm.canTestConnection).toBe(true);
  });

  it('computes regionHint for aws-s3', () => {
    const wrapper = mountWithPlugins(LakeHouseCredentials, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.regionHint).toContain('AWS region');
  });

  it('computes regionPlaceholder for aws-s3', () => {
    const wrapper = mountWithPlugins(LakeHouseCredentials, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.regionPlaceholder).toBe('us-east-1');
  });

  it('validates URL correctly', () => {
    const wrapper = mountWithPlugins(LakeHouseCredentials, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidUrl('https://example.com')).toBe(true);
    expect(wrapper.vm.isValidUrl('not-a-url')).toBe(false);
  });

  it('validates JSON correctly', () => {
    const wrapper = mountWithPlugins(LakeHouseCredentials, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidJson('{"key": "value"}')).toBe(true);
    expect(wrapper.vm.isValidJson('not json')).toBe(false);
  });
});

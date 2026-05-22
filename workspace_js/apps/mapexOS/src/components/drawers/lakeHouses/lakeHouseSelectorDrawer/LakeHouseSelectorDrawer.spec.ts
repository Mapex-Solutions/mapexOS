import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import LakeHouseSelectorDrawer from './LakeHouseSelectorDrawer.vue';

// Mock useLogger
vi.mock('@composables/useLogger', () => ({
  useLogger: vi.fn(() => ({
    warn: vi.fn(),
    error: vi.fn(),
    info: vi.fn(),
  })),
}));

describe('LakeHouseSelectorDrawer', () => {
  const defaultProps = {
    modelValue: false,
    selectedLakeHouseId: null,
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: defaultProps,
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes showDialog from modelValue', () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: { ...defaultProps, modelValue: true },
    });
    expect(wrapper.vm.showDialog).toBe(true);
  });

  it('emits update:modelValue when showDialog is set', () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: defaultProps,
    });
    wrapper.vm.showDialog = true;
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([true]);
  });

  it('computes statusOptions with 3 entries', () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: defaultProps,
    });
    expect(wrapper.vm.statusOptions).toHaveLength(3);
    expect(wrapper.vm.statusOptions[0].label).toBe('All');
  });

  it('computes typeOptions with 5 entries', () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: defaultProps,
    });
    expect(wrapper.vm.typeOptions).toHaveLength(5);
  });

  it('returns correct icon for getLakeHouseIcon', () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: defaultProps,
    });
    expect(wrapper.vm.getLakeHouseIcon('aws-s3')).toBe('mdi-aws');
    expect(wrapper.vm.getLakeHouseIcon('azure-blob')).toBe('mdi-microsoft-azure');
    expect(wrapper.vm.getLakeHouseIcon('gcp-storage')).toBe('mdi-google-cloud');
    expect(wrapper.vm.getLakeHouseIcon('minio')).toBe('mdi-database');
    expect(wrapper.vm.getLakeHouseIcon('unknown')).toBe('storage');
  });

  it('returns correct color for getLakeHouseIconColor', () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: defaultProps,
    });
    expect(wrapper.vm.getLakeHouseIconColor('aws-s3')).toBe('orange-6');
    expect(wrapper.vm.getLakeHouseIconColor('azure-blob')).toBe('blue-6');
    expect(wrapper.vm.getLakeHouseIconColor('gcp-storage')).toBe('red-6');
    expect(wrapper.vm.getLakeHouseIconColor('minio')).toBe('purple-6');
    expect(wrapper.vm.getLakeHouseIconColor('unknown')).toBe('purple-6');
  });

  it('returns correct label for getTypeLabel', () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: defaultProps,
    });
    expect(wrapper.vm.getTypeLabel('aws-s3')).toBe('AWS S3');
    expect(wrapper.vm.getTypeLabel('azure-blob')).toBe('Azure Blob');
    expect(wrapper.vm.getTypeLabel('gcp-storage')).toBe('GCP Storage');
    expect(wrapper.vm.getTypeLabel('minio')).toBe('MinIO');
    expect(wrapper.vm.getTypeLabel('other')).toBe('other');
  });

  it('isSelected returns true when id matches selectedLakeHouseId', () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: { ...defaultProps, selectedLakeHouseId: 'dl-1' },
    });
    expect(wrapper.vm.isSelected({ id: 'dl-1', name: 'Test', type: 'aws-s3', status: true })).toBe(true);
  });

  it('isSelected returns false when id does not match', () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: { ...defaultProps, selectedLakeHouseId: 'dl-1' },
    });
    expect(wrapper.vm.isSelected({ id: 'dl-2', name: 'Test', type: 'aws-s3', status: true })).toBe(false);
  });

  it('handleCancel emits cancel and closes drawer', () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: defaultProps,
    });
    wrapper.vm.handleCancel();
    expect(wrapper.emitted('cancel')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });

  it('selectLakeHouse emits select and closes drawer', () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: defaultProps,
    });
    const lakeHouse = { id: 'dl-1', name: 'Test', type: 'aws-s3' as const, status: true };
    wrapper.vm.selectLakeHouse(lakeHouse);
    expect(wrapper.emitted('select')).toBeTruthy();
    expect(wrapper.emitted('select')![0]).toEqual([lakeHouse]);
  });

  it('filteredLakeHouses filters by name', async () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: { ...defaultProps, modelValue: true },
    });
    // Wait for fetchLakeHouses to populate
    await new Promise(resolve => setTimeout(resolve, 400));
    wrapper.vm.filters.name = 'Amazon';
    expect(wrapper.vm.filteredLakeHouses.every(
      (dl: any) => dl.name.toLowerCase().includes('amazon') || dl.description?.toLowerCase().includes('amazon'),
    )).toBe(true);
  });

  it('filteredLakeHouses filters by status', async () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: { ...defaultProps, modelValue: true },
    });
    await new Promise(resolve => setTimeout(resolve, 400));
    wrapper.vm.filters.status = false;
    expect(wrapper.vm.filteredLakeHouses.every((dl: any) => dl.status === false)).toBe(true);
  });

  it('filteredLakeHouses filters by type', async () => {
    const wrapper = mountWithPlugins(LakeHouseSelectorDrawer, {
      props: { ...defaultProps, modelValue: true },
    });
    await new Promise(resolve => setTimeout(resolve, 400));
    wrapper.vm.filters.type = 'minio';
    expect(wrapper.vm.filteredLakeHouses.every((dl: any) => dl.type === 'minio')).toBe(true);
  });
});

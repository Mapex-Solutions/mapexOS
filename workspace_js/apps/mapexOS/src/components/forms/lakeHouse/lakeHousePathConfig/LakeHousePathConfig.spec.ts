import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import LakeHousePathConfig from './LakeHousePathConfig.vue';

vi.mock('@components/chips', () => ({
  DetailChip: { name: 'DetailChip', template: '<span />' },
}));

vi.mock('@components/tooltips', () => ({
  AppTooltip: { name: 'AppTooltip', template: '<div />' },
}));

vi.mock('vue-draggable-plus', () => ({
  VueDraggable: { name: 'VueDraggable', template: '<div><slot /></div>' },
}));

vi.mock('@utils/alert/notify', () => ({
  notifySuccess: vi.fn(),
  notifyFail: vi.fn(),
}));

vi.mock('./constants', () => ({
  DEFAULT_PATH_CONFIG: {
    basePath: 'datalake',
    partitions: [],
    compression: 'gzip',
    filePrefix: 'data_export_',
    maxFileSize: 100,
  },
}));

const BASE_MODEL = {
  pathConfig: {
    basePath: 'datalake',
    partitions: ['year', 'month'],
    compression: 'gzip',
    filePrefix: 'data_export_',
    maxFileSize: 100,
  },
};

describe('LakeHousePathConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(LakeHousePathConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes fullPathPreview with partitions', () => {
    const wrapper = mountWithPlugins(LakeHousePathConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.fullPathPreview).toContain('datalake/');
    expect(wrapper.vm.fullPathPreview).toContain('year=');
    expect(wrapper.vm.fullPathPreview).toContain('.gzip');
  });

  it('computes fullPathPreview without partitions', () => {
    const wrapper = mountWithPlugins(LakeHousePathConfig, {
      props: {
        modelValue: {
          pathConfig: { ...BASE_MODEL.pathConfig, partitions: [] },
        },
      },
    });
    expect(wrapper.vm.fullPathPreview).toBe('datalake/data_export__20251215_140000.json.gzip');
  });

  it('isPartitionSelected returns true for selected partition', () => {
    const wrapper = mountWithPlugins(LakeHousePathConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isPartitionSelected('year')).toBe(true);
    expect(wrapper.vm.isPartitionSelected('hour')).toBe(false);
  });

  it('addPartition adds a new partition', () => {
    const wrapper = mountWithPlugins(LakeHousePathConfig, {
      props: { modelValue: BASE_MODEL },
    });
    wrapper.vm.addPartition('day');
    expect(wrapper.vm.modelRef.pathConfig.partitions).toContain('day');
  });

  it('addPartition does not duplicate', () => {
    const wrapper = mountWithPlugins(LakeHousePathConfig, {
      props: { modelValue: BASE_MODEL },
    });
    const initial = wrapper.vm.modelRef.pathConfig.partitions.length;
    wrapper.vm.addPartition('year');
    expect(wrapper.vm.modelRef.pathConfig.partitions.length).toBe(initial);
  });

  it('removePartition removes by index', () => {
    const wrapper = mountWithPlugins(LakeHousePathConfig, {
      props: { modelValue: BASE_MODEL },
    });
    wrapper.vm.removePartition(0);
    expect(wrapper.vm.modelRef.pathConfig.partitions).not.toContain('year');
  });

  it('getPartitionLabel returns correct label', () => {
    const wrapper = mountWithPlugins(LakeHousePathConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.getPartitionLabel('year')).toBe('Year (YYYY)');
    expect(wrapper.vm.getPartitionLabel('unknown')).toBe('unknown');
  });

  it('getPartitionExample returns correct example', () => {
    const wrapper = mountWithPlugins(LakeHousePathConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.getPartitionExample('year')).toBe('year=2025');
  });
});

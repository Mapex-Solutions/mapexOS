import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import AssetSelector from './AssetSelector.vue';

/** Mock composables and services */
vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    debug: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  }),
}));

vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));

vi.mock('@services/mapex', () => ({
  apis: {
    assets: {
      asset: {
        list: vi.fn().mockResolvedValue({ items: [], pagination: { totalItems: 0, totalPages: 1 } }),
        getById: vi.fn().mockResolvedValue({ id: '1', name: 'Test' }),
      },
    },
    mapexOS: {
      lists: {
        list: vi.fn().mockResolvedValue({ items: [] }),
      },
    },
  },
}));

vi.mock('@components/tooltips', () => ({
  AppTooltip: { template: '<span />' },
}));

describe('AssetSelector', () => {
  const defaultProps = {
    modelValue: null as string | null,
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(AssetSelector, { props: { ...defaultProps } });
    expect(wrapper.exists()).toBe(true);
  });

  it('initializes filters with empty defaults', () => {
    const wrapper = mountWithPlugins(AssetSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.filters.search).toBe('');
    expect(wrapper.vm.filters.status).toBeNull();
    expect(wrapper.vm.filters.categoryId).toBeUndefined();
    expect(wrapper.vm.filters.manufacturerId).toBeUndefined();
    expect(wrapper.vm.filters.model).toBeUndefined();
  });

  it('statusOptions returns 3 options (All, Active, Inactive)', () => {
    const wrapper = mountWithPlugins(AssetSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.statusOptions).toHaveLength(3);
    expect(wrapper.vm.statusOptions[0].label).toBe('All');
    expect(wrapper.vm.statusOptions[1].label).toBe('Active');
    expect(wrapper.vm.statusOptions[2].label).toBe('Inactive');
  });

  it('statusOptions values are null, true, false', () => {
    const wrapper = mountWithPlugins(AssetSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.statusOptions[0].value).toBeNull();
    expect(wrapper.vm.statusOptions[1].value).toBe(true);
    expect(wrapper.vm.statusOptions[2].value).toBe(false);
  });

  it('showDialog is initially false', () => {
    const wrapper = mountWithPlugins(AssetSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.showDialog).toBe(false);
  });

  it('selectedAsset is initially null', () => {
    const wrapper = mountWithPlugins(AssetSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.selectedAsset).toBeNull();
  });

  it('assets is initially empty array', () => {
    const wrapper = mountWithPlugins(AssetSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.assets).toEqual([]);
  });

  it('clearSelection resets selectedAsset and emits null', () => {
    const wrapper = mountWithPlugins(AssetSelector, { props: { ...defaultProps } });
    wrapper.vm.clearSelection();
    expect(wrapper.vm.selectedAsset).toBeNull();
    expect(wrapper.emitted('update:modelValue')![0]![0]).toBeNull();
    expect(wrapper.emitted('update:selectedAsset')![0]![0]).toBeNull();
  });

  it('onFilterChange resets currentPage to 1', () => {
    const wrapper = mountWithPlugins(AssetSelector, { props: { ...defaultProps } });
    wrapper.vm.currentPage = 5;
    wrapper.vm.onFilterChange();
    expect(wrapper.vm.currentPage).toBe(1);
  });

  it('uses default label prop value', () => {
    const wrapper = mountWithPlugins(AssetSelector, { props: { ...defaultProps } });
    expect(wrapper.props('label')).toBe('Select Asset');
  });

  it('accepts custom label prop', () => {
    const wrapper = mountWithPlugins(AssetSelector, {
      props: { ...defaultProps, label: 'Choose Asset' },
    });
    expect(wrapper.props('label')).toBe('Choose Asset');
  });
});

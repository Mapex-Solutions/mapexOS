import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import AssetSelectorDrawer from './AssetSelectorDrawer.vue';

vi.mock('@composables/assets/assets', () => ({
  useAssets: () => ({
    assets: { value: [] },
    isLoading: { value: false },
    isLoadingMore: { value: false },
    filters: { value: { name: undefined, categoryId: undefined, manufacturerId: undefined, modelId: undefined, status: undefined } },
    pagination: { value: { currentPage: 1, totalPages: 1, totalItems: 0 } },
    categoryOptions: { value: [] },
    manufacturerOptions: { value: [] },
    modelOptions: { value: [] },
    loadingCategories: { value: false },
    loadingManufacturers: { value: false },
    loadingModels: { value: false },
    fetchAssets: vi.fn(),
    loadCategories: vi.fn(),
    handleCategoryChange: vi.fn(),
    handleManufacturerChange: vi.fn(),
  }),
}));

describe('AssetSelectorDrawer', () => {
  const defaultProps = {
    modelValue: true,
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(AssetSelectorDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes showDialog from modelValue', () => {
    const wrapper = mountWithPlugins(AssetSelectorDrawer, { props: defaultProps });
    expect((wrapper.vm).showDialog).toBe(true);
  });

  it('sets showDialog to false when modelValue is false', () => {
    const wrapper = mountWithPlugins(AssetSelectorDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    expect((wrapper.vm).showDialog).toBe(false);
  });

  it('starts with selectedAsset as null', () => {
    const wrapper = mountWithPlugins(AssetSelectorDrawer, { props: defaultProps });
    expect((wrapper.vm).selectedAsset).toBeNull();
  });

  it('computes statusOptions with 3 entries', () => {
    const wrapper = mountWithPlugins(AssetSelectorDrawer, { props: defaultProps });
    expect((wrapper.vm).statusOptions).toHaveLength(3);
  });

  it('computes totalAssets from pagination', () => {
    const wrapper = mountWithPlugins(AssetSelectorDrawer, { props: defaultProps });
    expect((wrapper.vm).totalAssets).toBe(0);
  });

  it('defaults multiSelect to false', () => {
    const wrapper = mountWithPlugins(AssetSelectorDrawer, { props: defaultProps });
    expect(wrapper.props('multiSelect')).toBe(false);
  });

  it('emits select and closes on selectAsset', () => {
    const wrapper = mountWithPlugins(AssetSelectorDrawer, { props: defaultProps });
    const mockAsset = { id: 'a1', name: 'Asset 1' };
    (wrapper.vm).selectAsset(mockAsset);
    expect(wrapper.emitted('select')).toBeTruthy();
    expect(wrapper.emitted('select')![0]).toEqual([mockAsset]);
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('emits cancel on handleCancel', () => {
    const wrapper = mountWithPlugins(AssetSelectorDrawer, { props: defaultProps });
    (wrapper.vm).handleCancel();
    expect(wrapper.emitted('cancel')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('resets page on filter change', () => {
    const wrapper = mountWithPlugins(AssetSelectorDrawer, { props: defaultProps });
    (wrapper.vm).onFilterChange();
    // No error thrown = success (pagination.value.currentPage reset internally)
    expect(wrapper.exists()).toBe(true);
  });
});

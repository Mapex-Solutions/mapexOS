import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import AssetTemplateSelectorDialog from './AssetTemplateSelectorDialog.vue';

vi.mock('@components/dialogs/common/genericSelectorDialog', () => ({
  GenericSelectorDialog: { name: 'GenericSelectorDialog', template: '<div />' },
}));

vi.mock('@composables/assets/assetTemplates', () => ({
  useAssetTemplates: () => ({
    templates: { value: [] },
    isLoading: { value: false },
    isLoadingMore: { value: false },
    filters: { value: { name: undefined, categoryId: undefined, status: undefined, manufacturerId: undefined, modelId: undefined } },
    pagination: { value: { currentPage: 1, totalPages: 1, totalItems: 0 } },
    categoryOptions: { value: [] },
    manufacturerOptions: { value: [] },
    modelOptions: { value: [] },
    loadingCategories: { value: false },
    loadingManufacturers: { value: false },
    loadingModels: { value: false },
    fetchTemplates: vi.fn(),
    loadCategories: vi.fn(),
    handleCategoryChange: vi.fn(),
    handleManufacturerChange: vi.fn(),
  }),
}));

const BASE_PROPS = {
  modelValue: false,
};

describe('AssetTemplateSelectorDialog', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes dialogTitle for multi-select mode', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDialog, {
      props: { ...BASE_PROPS, multiSelect: true },
    });
    expect(wrapper.vm.dialogTitle).toBe('Select Asset Templates');
  });

  it('computes dialogTitle for single-select mode', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDialog, {
      props: { ...BASE_PROPS, multiSelect: false },
    });
    expect(wrapper.vm.dialogTitle).toBe('Select Asset Template');
  });

  it('computes bannerText for multi-select', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDialog, {
      props: { ...BASE_PROPS, multiSelect: true },
    });
    expect(wrapper.vm.bannerText).toContain('Confirm');
  });

  it('computes bannerText for single-select', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDialog, {
      props: { ...BASE_PROPS, multiSelect: false },
    });
    expect(wrapper.vm.bannerText).toContain('Click');
  });

  it('computes totalTemplates from pagination', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.totalTemplates).toBe(0);
  });

  it('computes hasMorePages as false when on last page', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.hasMorePages).toBe(false);
  });

  it('computes statusOptions with three options', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.statusOptions).toHaveLength(3);
    expect(wrapper.vm.statusOptions[0].label).toBe('All');
  });

  it('emits select on handleSelect', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDialog, { props: BASE_PROPS });
    const items = [{ id: 'tpl-1', name: 'Template A' }];
    wrapper.vm.handleSelect(items);
    const emitted = wrapper.emitted('select')!;
    expect(emitted[0]![0]).toEqual(items);
  });

  it('emits cancel on handleCancel', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDialog, { props: BASE_PROPS });
    wrapper.vm.handleCancel();
    expect(wrapper.emitted('cancel')).toBeTruthy();
  });
});

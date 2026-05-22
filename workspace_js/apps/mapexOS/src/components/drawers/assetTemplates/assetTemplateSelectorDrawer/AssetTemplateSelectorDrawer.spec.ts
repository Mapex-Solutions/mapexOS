import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import AssetTemplateSelectorDrawer from './AssetTemplateSelectorDrawer.vue';

vi.mock('@composables/assets/assetTemplates', () => ({
  useAssetTemplates: () => ({
    templates: { value: [] },
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
    fetchTemplates: vi.fn(),
    loadCategories: vi.fn(),
    handleCategoryChange: vi.fn(),
    handleManufacturerChange: vi.fn(),
  }),
}));

vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    debug: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  }),
}));

describe('AssetTemplateSelectorDrawer', () => {
  const defaultProps = {
    modelValue: true,
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts with empty selectedTemplates', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDrawer, { props: defaultProps });
    expect((wrapper.vm).selectedTemplates).toEqual([]);
  });

  it('computes statusOptions with 3 entries', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDrawer, { props: defaultProps });
    expect((wrapper.vm).statusOptions).toHaveLength(3);
  });

  it('computes totalTemplates from pagination', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDrawer, { props: defaultProps });
    expect((wrapper.vm).totalTemplates).toBe(0);
  });

  it('defaults multiSelect to true', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDrawer, { props: defaultProps });
    expect(wrapper.props('multiSelect')).toBe(true);
  });

  it('computes drawerTitle for multiSelect', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDrawer, { props: defaultProps });
    expect((wrapper.vm).drawerTitle).toBe('Select Asset Templates');
  });

  it('computes drawerTitle for single select', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDrawer, {
      props: { ...defaultProps, multiSelect: false },
    });
    expect((wrapper.vm).drawerTitle).toBe('Select Asset Template');
  });

  it('toggleTemplate in single-select mode emits select and closes', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDrawer, {
      props: { ...defaultProps, multiSelect: false },
    });
    const mockTemplate = { id: 't1', name: 'Template 1' };
    (wrapper.vm).toggleTemplate(mockTemplate);
    expect(wrapper.emitted('select')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('toggleTemplate in multi-select mode adds to selection', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDrawer, { props: defaultProps });
    const mockTemplate = { id: 't1', name: 'Template 1' };
    (wrapper.vm).toggleTemplate(mockTemplate);
    expect((wrapper.vm).selectedTemplates).toHaveLength(1);
  });

  it('toggleTemplate in multi-select mode removes already-selected item', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDrawer, { props: defaultProps });
    const mockTemplate = { id: 't1', name: 'Template 1' };
    (wrapper.vm).toggleTemplate(mockTemplate);
    (wrapper.vm).toggleTemplate(mockTemplate);
    expect((wrapper.vm).selectedTemplates).toHaveLength(0);
  });

  it('isSelected returns true for selected template', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDrawer, { props: defaultProps });
    const mockTemplate = { id: 't1', name: 'Template 1' };
    (wrapper.vm).toggleTemplate(mockTemplate);
    expect((wrapper.vm).isSelected(mockTemplate)).toBe(true);
  });

  it('emits cancel on handleCancel', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDrawer, { props: defaultProps });
    (wrapper.vm).handleCancel();
    expect(wrapper.emitted('cancel')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('confirmSelection emits select and closes', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelectorDrawer, { props: defaultProps });
    const mockTemplate = { id: 't1', name: 'Template 1' };
    (wrapper.vm).toggleTemplate(mockTemplate);
    (wrapper.vm).confirmSelection();
    expect(wrapper.emitted('select')).toBeTruthy();
    expect(wrapper.emitted('select')![0]).toEqual([[mockTemplate]]);
  });
});

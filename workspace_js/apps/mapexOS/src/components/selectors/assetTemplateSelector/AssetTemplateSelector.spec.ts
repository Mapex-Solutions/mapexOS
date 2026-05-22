import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins, createMockTranslations } from '@src/test/helpers';
import AssetTemplateSelector from './AssetTemplateSelector.vue';

/** Mock composables and services */
vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    debug: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  }),
}));

vi.mock('@src/composables/i18n/components/selectors/useAssetTemplateSelectorTranslations', () => ({
  useAssetTemplateSelectorTranslations: () => createMockTranslations(),
}));

vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));

vi.mock('@services/mapex', () => ({
  apis: {
    assets: {
      assetTemplate: {
        list: vi.fn().mockResolvedValue({ items: [], pagination: { totalItems: 0, totalPages: 1 } }),
        getById: vi.fn().mockResolvedValue({ id: '1', name: 'Template1' }),
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

describe('AssetTemplateSelector', () => {
  const defaultProps = {
    modelValue: null as string | null,
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelector, { props: { ...defaultProps } });
    expect(wrapper.exists()).toBe(true);
  });

  it('initializes filters with empty defaults', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.filters.search).toBe('');
    expect(wrapper.vm.filters.isSystem).toBeNull();
    expect(wrapper.vm.filters.isTemplate).toBeNull();
    expect(wrapper.vm.filters.categoryId).toBeNull();
    expect(wrapper.vm.filters.manufacturerId).toBeNull();
    expect(wrapper.vm.filters.modelId).toBeNull();
  });

  it('templateTypeOptions has 3 options', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.templateTypeOptions).toHaveLength(3);
  });

  it('templateSourceOptions has 3 options', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.templateSourceOptions).toHaveLength(3);
  });

  it('templates is initially empty', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.templates).toEqual([]);
  });

  it('selectedTemplate is initially null', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.selectedTemplate).toBeNull();
  });

  it('totalTemplates is initially 0', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.totalTemplates).toBe(0);
  });

  it('onFilterChange resets currentPage to 1', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelector, { props: { ...defaultProps } });
    wrapper.vm.currentPage = 3;
    wrapper.vm.onFilterChange();
    expect(wrapper.vm.currentPage).toBe(1);
  });

  it('selectTemplate sets selectedTemplate and emits values', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelector, { props: { ...defaultProps } });
    const template = { id: 'tmpl-1', name: 'My Template', manufacturerName: 'Acme', modelName: 'X100' } as any;
    wrapper.vm.selectTemplate(template);
    expect(wrapper.vm.selectedTemplate).toEqual(template);
    expect(wrapper.emitted('update:modelValue')![0]![0]).toBe('tmpl-1');
    expect(wrapper.emitted('update:selectedTemplate')![0]![0]).toEqual(template);
  });

  it('onCategoryChange resets manufacturer and model filters', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelector, { props: { ...defaultProps } });
    wrapper.vm.filters.manufacturerId = 'mfr-1';
    wrapper.vm.filters.modelId = 'mdl-1';
    wrapper.vm.onCategoryChange('cat-1');
    expect(wrapper.vm.filters.manufacturerId).toBeNull();
    expect(wrapper.vm.filters.modelId).toBeNull();
  });

  it('onManufacturerChange resets model filter', () => {
    const wrapper = mountWithPlugins(AssetTemplateSelector, { props: { ...defaultProps } });
    wrapper.vm.filters.modelId = 'mdl-1';
    wrapper.vm.onManufacturerChange('mfr-1');
    expect(wrapper.vm.filters.modelId).toBeNull();
  });
});

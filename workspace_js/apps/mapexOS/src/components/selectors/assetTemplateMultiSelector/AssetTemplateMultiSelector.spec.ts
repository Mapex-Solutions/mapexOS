import { describe, it, expect, vi } from 'vitest';
import { ref } from 'vue';
import { mountWithPlugins } from '@src/test/helpers';
import AssetTemplateMultiSelector from './AssetTemplateMultiSelector.vue';

/** Mock composables and services */
vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    debug: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  }),
}));

vi.mock('@composables/assets/assetTemplates', () => ({
  useAssetTemplates: () => ({
    templates: ref([]),
    isLoading: ref(false),
    isLoadingMore: ref(false),
    filters: ref({ name: '', categoryId: null, manufacturerId: null, modelId: null, status: undefined }),
    pagination: ref({ currentPage: 1, totalPages: 1, totalItems: 0 }),
    categoryOptions: ref([]),
    manufacturerOptions: ref([]),
    modelOptions: ref([]),
    loadingCategories: ref(false),
    loadingManufacturers: ref(false),
    loadingModels: ref(false),
    fetchTemplates: vi.fn(),
    loadCategories: vi.fn(),
    handleCategoryChange: vi.fn(),
    handleManufacturerChange: vi.fn(),
  }),
}));

vi.mock('@services/mapex', () => ({
  apis: {
    assets: {
      assetTemplate: {
        getById: vi.fn().mockResolvedValue({ id: '1', name: 'Template1' }),
      },
    },
  },
}));

vi.mock('@components/tooltips', () => ({
  AppTooltip: { template: '<span />' },
}));

describe('AssetTemplateMultiSelector', () => {
  const defaultProps = {
    modelValue: [] as string[],
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(AssetTemplateMultiSelector, { props: { ...defaultProps } });
    expect(wrapper.exists()).toBe(true);
  });

  it('selectedTemplates is initially empty', () => {
    const wrapper = mountWithPlugins(AssetTemplateMultiSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.selectedTemplates).toEqual([]);
  });

  it('showDialog is initially false', () => {
    const wrapper = mountWithPlugins(AssetTemplateMultiSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.showDialog).toBe(false);
  });

  it('statusOptions has 3 options', () => {
    const wrapper = mountWithPlugins(AssetTemplateMultiSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.statusOptions).toHaveLength(3);
    expect(wrapper.vm.statusOptions[0].label).toBe('All');
  });

  it('selectedIds is empty when no templates selected', () => {
    const wrapper = mountWithPlugins(AssetTemplateMultiSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.selectedIds).toEqual([]);
  });

  it('toggleTemplate adds template to selection', () => {
    const wrapper = mountWithPlugins(AssetTemplateMultiSelector, { props: { ...defaultProps } });
    const template = { id: 'tmpl-1', name: 'T1', assetIdPath: '/path' } as any;
    wrapper.vm.toggleTemplate(template);
    expect(wrapper.vm.selectedTemplates).toHaveLength(1);
    expect(wrapper.vm.selectedTemplates[0].id).toBe('tmpl-1');
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });

  it('toggleTemplate removes template from selection when already selected', () => {
    const wrapper = mountWithPlugins(AssetTemplateMultiSelector, { props: { ...defaultProps } });
    const template = { id: 'tmpl-1', name: 'T1', assetIdPath: '/path' } as any;
    wrapper.vm.toggleTemplate(template);
    wrapper.vm.toggleTemplate(template);
    expect(wrapper.vm.selectedTemplates).toHaveLength(0);
  });

  it('isSelected returns false when template not in selection', () => {
    const wrapper = mountWithPlugins(AssetTemplateMultiSelector, { props: { ...defaultProps } });
    const template = { id: 'tmpl-1', name: 'T1' } as any;
    expect(wrapper.vm.isSelected(template)).toBe(false);
  });

  it('isSelected returns true when template is in selection', () => {
    const wrapper = mountWithPlugins(AssetTemplateMultiSelector, { props: { ...defaultProps } });
    const template = { id: 'tmpl-1', name: 'T1' } as any;
    wrapper.vm.toggleTemplate(template);
    expect(wrapper.vm.isSelected(template)).toBe(true);
  });

  it('removeTemplate removes by ID', () => {
    const wrapper = mountWithPlugins(AssetTemplateMultiSelector, { props: { ...defaultProps } });
    const t1 = { id: 'tmpl-1', name: 'T1' } as any;
    const t2 = { id: 'tmpl-2', name: 'T2' } as any;
    wrapper.vm.toggleTemplate(t1);
    wrapper.vm.toggleTemplate(t2);
    wrapper.vm.removeTemplate('tmpl-1');
    expect(wrapper.vm.selectedTemplates).toHaveLength(1);
    expect(wrapper.vm.selectedTemplates[0].id).toBe('tmpl-2');
  });

  it('clearAll empties selection and emits', () => {
    const wrapper = mountWithPlugins(AssetTemplateMultiSelector, { props: { ...defaultProps } });
    wrapper.vm.toggleTemplate({ id: 'tmpl-1', name: 'T1' } as any);
    wrapper.vm.clearAll();
    expect(wrapper.vm.selectedTemplates).toHaveLength(0);
  });

  it('extractedPaths returns paths from templates with assetIdPath', () => {
    const wrapper = mountWithPlugins(AssetTemplateMultiSelector, { props: { ...defaultProps } });
    wrapper.vm.toggleTemplate({ id: 'a', name: 'A', assetIdPath: '/a/path' } as any);
    wrapper.vm.toggleTemplate({ id: 'b', name: 'B' } as any);
    const paths = wrapper.vm.extractedPaths;
    expect(paths).toHaveLength(1);
    expect(paths[0].templateId).toBe('a');
    expect(paths[0].assetIdPath).toBe('/a/path');
  });
});

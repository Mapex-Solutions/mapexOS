import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import LoopNodeConfig from './LoopNodeConfig.vue';

vi.mock('@src/composables/workflow', () => ({
  useWorkflowContext: () => ({
    nodes: { value: [] },
    getNodeType: () => null,
  }),
  usePluginI18n: () => ({ t: (key: string) => key }),
}));

vi.mock('@components/forms/fieldSourceSelector', () => ({
  FieldSourceSelector: { name: 'FieldSourceSelector', template: '<div />' },
}));

vi.mock('@components/dialogs/common/assetTemplateSelectorDialog', () => ({
  AssetTemplateSelectorDialog: { name: 'AssetTemplateSelectorDialog', template: '<div />' },
}));

vi.mock('@components/dialogs/common/genericSelectorDialog', () => ({
  GenericSelectorDialog: { name: 'GenericSelectorDialog', template: '<div />' },
}));

vi.mock('@services/mapex', () => ({
  apis: { assets: { assetTemplate: { getAvailableFields: vi.fn(), getById: vi.fn() } } },
}));

const BASE_CONFIG: Record<string, unknown> = { _nodeId: 'loop-1' };

describe('LoopNodeConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(LoopNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes source with defaults when not set', () => {
    const wrapper = mountWithPlugins(LoopNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.source).toEqual({ type: 'state', value: '' });
  });

  it('computes source from config', () => {
    const wrapper = mountWithPlugins(LoopNodeConfig, {
      props: { config: { ...BASE_CONFIG, source: { type: 'event', value: 'data.items' } } },
    });
    expect(wrapper.vm.source).toEqual({ type: 'event', value: 'data.items' });
  });

  it('computes hasTemplates as false when no template IDs', () => {
    const wrapper = mountWithPlugins(LoopNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.hasTemplates).toBe(false);
  });

  it('computes hasTemplates as true when template IDs exist', () => {
    const wrapper = mountWithPlugins(LoopNodeConfig, {
      props: { config: { ...BASE_CONFIG, selectedTemplateIds: ['tpl-1'] } },
    });
    expect(wrapper.vm.hasTemplates).toBe(true);
  });

  it('computes nodeOutputOptions as empty when no nodes', () => {
    const wrapper = mountWithPlugins(LoopNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.nodeOutputOptions).toEqual([]);
  });

  it('emits update:config on handleSourceUpdate', () => {
    const wrapper = mountWithPlugins(LoopNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleSourceUpdate({ type: 'literal', value: '[1,2,3]' });
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ source: { type: 'literal', value: '[1,2,3]' } });
  });

  it('opens template dialog on handleOpenTemplateSelector', () => {
    const wrapper = mountWithPlugins(LoopNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleOpenTemplateSelector();
    expect(wrapper.vm.templateDialogOpen).toBe(true);
  });

  it('opens field selector on handleOpenEventSelector', () => {
    const wrapper = mountWithPlugins(LoopNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleOpenEventSelector();
    expect(wrapper.vm.fieldSelectorOpen).toBe(true);
    expect(wrapper.vm.fieldSearchQuery).toBe('');
  });

  it('handles field search query update', () => {
    const wrapper = mountWithPlugins(LoopNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleFieldSearch('temp');
    expect(wrapper.vm.fieldSearchQuery).toBe('temp');
  });
});

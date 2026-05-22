import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import SetStateNodeConfig from './SetStateNodeConfig.vue';

vi.mock('@src/composables/workflow', () => ({
  useWorkflowContext: () => ({
    states: { value: [{ field: 'counter', type: 'number' }] },
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

const BASE_CONFIG: Record<string, unknown> = { _nodeId: 'set-1' };

describe('SetStateNodeConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(SetStateNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes stateFields from workflow context', () => {
    const wrapper = mountWithPlugins(SetStateNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.stateFields).toEqual([{ name: 'counter', type: 'number' }]);
  });

  it('computes operation as set by default', () => {
    const wrapper = mountWithPlugins(SetStateNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.operation).toBe('set');
  });

  it('computes operation from config', () => {
    const wrapper = mountWithPlugins(SetStateNodeConfig, {
      props: { config: { ...BASE_CONFIG, operation: 'increment' } },
    });
    expect(wrapper.vm.operation).toBe('increment');
  });

  it('computes targetField from config', () => {
    const wrapper = mountWithPlugins(SetStateNodeConfig, {
      props: { config: { ...BASE_CONFIG, targetField: 'counter' } },
    });
    expect(wrapper.vm.targetField).toBe('counter');
  });

  it('computes valueSource with defaults', () => {
    const wrapper = mountWithPlugins(SetStateNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.valueSource).toEqual({ type: 'literal', value: '' });
  });

  it('computes needsValueInput as true for set operation', () => {
    const wrapper = mountWithPlugins(SetStateNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.needsValueInput).toBe(true);
  });

  it('computes needsValueInput as false for remove operation', () => {
    const wrapper = mountWithPlugins(SetStateNodeConfig, {
      props: { config: { ...BASE_CONFIG, operation: 'remove' } },
    });
    expect(wrapper.vm.needsValueInput).toBe(false);
  });

  it('emits update:config on updateOperation', () => {
    const wrapper = mountWithPlugins(SetStateNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateOperation('append');
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ operation: 'append' });
  });

  it('emits update:config on updateTargetField', () => {
    const wrapper = mountWithPlugins(SetStateNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateTargetField('counter');
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ targetField: 'counter' });
  });

  it('emits update:config on handleValueSourceUpdate', () => {
    const wrapper = mountWithPlugins(SetStateNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleValueSourceUpdate({ type: 'literal', value: '42' });
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ valueSource: { type: 'literal', value: '42' } });
  });
});

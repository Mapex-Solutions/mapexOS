import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import ConditionNodeConfig from './ConditionNodeConfig.vue';

vi.mock('@src/composables/workflow', () => ({
  useWorkflowContext: () => ({
    states: { value: [{ field: 'status', type: 'string' }] },
  }),
  usePluginI18n: () => ({ t: (key: string) => key }),
}));

vi.mock('@components/dialogs/common/assetTemplateSelectorDialog', () => ({
  AssetTemplateSelectorDialog: { name: 'AssetTemplateSelectorDialog', template: '<div />' },
}));

vi.mock('@components/dialogs/common/genericSelectorDialog', () => ({
  GenericSelectorDialog: { name: 'GenericSelectorDialog', template: '<div />' },
}));

vi.mock('./ConditionGroupCard.vue', () => ({
  default: { name: 'ConditionGroupCard', template: '<div />' },
}));

vi.mock('./ConditionItemCard.vue', () => ({
  default: { name: 'ConditionItemCard', template: '<div />' },
}));

vi.mock('@services/mapex', () => ({
  apis: { assets: { assetTemplate: { getAvailableFields: vi.fn(), getById: vi.fn() } } },
}));

vi.mock('@src/pages/automations/rules/createEditRulePage/constants', () => ({
  ComparisonOperator: {
    Equals: 'equals',
    NotEquals: 'notEquals',
    IsNull: 'isNull',
    IsNotNull: 'isNotNull',
  },
}));

vi.mock('../constants', () => ({
  GROUP_LOGIC_OPTIONS: [
    { value: 'AND', label: 'AND', icon: 'join_full', color: 'blue-6', description: 'All must match' },
    { value: 'OR', label: 'OR', icon: 'join_inner', color: 'orange-6', description: 'Any must match' },
  ],
}));

const BASE_CONFIG: Record<string, unknown> = {};

describe('ConditionNodeConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(ConditionNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes stateFields from workflow context', () => {
    const wrapper = mountWithPlugins(ConditionNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.stateFields).toEqual([{ name: 'status', type: 'string' }]);
  });

  it('computes rootLogic as AND by default', () => {
    const wrapper = mountWithPlugins(ConditionNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.rootLogic).toBe('AND');
  });

  it('computes rootLogic from config', () => {
    const wrapper = mountWithPlugins(ConditionNodeConfig, {
      props: { config: { logic: 'OR' } },
    });
    expect(wrapper.vm.rootLogic).toBe('OR');
  });

  it('computes items as empty array by default', () => {
    const wrapper = mountWithPlugins(ConditionNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.items).toEqual([]);
  });

  it('computes currentOperator from rootLogic', () => {
    const wrapper = mountWithPlugins(ConditionNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.currentOperator.value).toBe('AND');
  });

  it('emits update:config on updateRootLogic', () => {
    const wrapper = mountWithPlugins(ConditionNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateRootLogic('OR');
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ logic: 'OR' });
  });

  it('emits update:config on addCondition', () => {
    const wrapper = mountWithPlugins(ConditionNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.addCondition();
    const emitted = wrapper.emitted('update:config')!;
    const items = (emitted[0]![0] as any).items;
    expect(items).toHaveLength(1);
    expect(items[0].type).toBe('condition');
  });

  it('emits update:config on addGroup', () => {
    const wrapper = mountWithPlugins(ConditionNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.addGroup();
    const emitted = wrapper.emitted('update:config')!;
    const items = (emitted[0]![0] as any).items;
    expect(items).toHaveLength(1);
    expect(items[0].type).toBe('group');
  });

  it('removes item by index via removeItem', () => {
    const existingItems = [
      { type: 'condition', data: { id: 'c1', name: 'c1', field: { type: 'event', value: '' }, operator: 'equals', value: { type: 'input', value: '' } } },
      { type: 'condition', data: { id: 'c2', name: 'c2', field: { type: 'event', value: '' }, operator: 'equals', value: { type: 'input', value: '' } } },
    ];
    const wrapper = mountWithPlugins(ConditionNodeConfig, {
      props: { config: { items: existingItems } },
    });
    wrapper.vm.removeItem(0);
    const emitted = wrapper.emitted('update:config')!;
    const items = (emitted[0]![0] as any).items;
    expect(items).toHaveLength(1);
    expect(items[0].data.id).toBe('c2');
  });

  it('opens template dialog via handleEventFieldRequest when no templates', () => {
    const wrapper = mountWithPlugins(ConditionNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleEventFieldRequest({ side: 'field', conditionId: 'c1' });
    expect(wrapper.vm.templateDialogOpen).toBe(true);
  });

  it('opens field selector via handleEventFieldRequest when templates exist', () => {
    const wrapper = mountWithPlugins(ConditionNodeConfig, {
      props: { config: { selectedTemplateIds: ['tpl-1'] } },
    });
    wrapper.vm.handleEventFieldRequest({ side: 'field', conditionId: 'c1' });
    expect(wrapper.vm.fieldSelectorOpen).toBe(true);
  });
});

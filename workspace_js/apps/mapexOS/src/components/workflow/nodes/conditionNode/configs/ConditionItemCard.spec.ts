import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import ConditionItemCard from './ConditionItemCard.vue';

vi.mock('@src/composables/workflow', () => ({
  useWorkflowContext: () => ({
    nodes: { value: [] },
    getNodeType: () => null,
  }),
  usePluginI18n: () => ({ t: (key: string) => key }),
}));

vi.mock('@components/forms/fieldSourceSelector', () => ({
  FieldSourceSelector: { name: 'FieldSourceSelector', template: '<div />' },
  SOURCE_TYPE_OPTIONS: [
    { value: 'event', icon: 'bolt', color: 'blue-6' },
    { value: 'state', icon: 'storage', color: 'purple-6' },
    { value: 'literal', icon: 'edit', color: 'grey-6' },
  ],
}));

vi.mock('@components/tooltips/appTooltip', () => ({
  AppTooltip: { name: 'AppTooltip', template: '<div />' },
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
  CONDITION_OPERATOR_OPTIONS: [
    { value: 'equals', label: 'Equals', symbol: '=' },
    { value: 'notEquals', label: 'Not Equals', symbol: '!=' },
    { value: 'isNull', label: 'Is Null', symbol: 'null' },
  ],
}));

const BASE_CONDITION = {
  id: 'c1',
  name: 'condition',
  field: { type: 'event' as const, value: 'data.status' },
  operator: 'equals',
  value: { type: 'input' as const, value: 'ACTIVE' },
};

describe('ConditionItemCard', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: BASE_CONDITION, stateFields: [], canRemove: true },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts collapsed by default', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: BASE_CONDITION, stateFields: [], canRemove: true },
    });
    expect(wrapper.vm.isExpanded).toBe(false);
  });

  it('toggles expanded state', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: BASE_CONDITION, stateFields: [], canRemove: true },
    });
    wrapper.vm.toggleExpanded();
    expect(wrapper.vm.isExpanded).toBe(true);
  });

  it('computes isUnaryOperator as false for equals', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: BASE_CONDITION, stateFields: [], canRemove: true },
    });
    expect(wrapper.vm.isUnaryOperator).toBe(false);
  });

  it('computes isUnaryOperator as true for is_null', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: { ...BASE_CONDITION, operator: 'isNull' }, stateFields: [], canRemove: true },
    });
    expect(wrapper.vm.isUnaryOperator).toBe(true);
  });

  it('computes collapsedSummary with field, operator, and value', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: BASE_CONDITION, stateFields: [], canRemove: true },
    });
    expect(wrapper.vm.collapsedSummary).toBe('data.status = ACTIVE');
  });

  it('computes collapsedSummary for unary operator', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: { ...BASE_CONDITION, operator: 'isNull' }, stateFields: [], canRemove: true },
    });
    expect(wrapper.vm.collapsedSummary).toBe('data.status null');
  });

  it('computes fieldSource from condition', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: BASE_CONDITION, stateFields: [], canRemove: true },
    });
    expect(wrapper.vm.fieldSource).toEqual({ type: 'event', value: 'data.status' });
  });

  it('computes valueSource from condition', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: BASE_CONDITION, stateFields: [], canRemove: true },
    });
    expect(wrapper.vm.valueSource).toEqual({ type: 'input', value: 'ACTIVE' });
  });

  it('emits update:condition on handleFieldUpdate', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: BASE_CONDITION, stateFields: [], canRemove: true },
    });
    wrapper.vm.handleFieldUpdate({ type: 'state', value: 'counter' });
    const emitted = wrapper.emitted('update:condition')!;
    expect(emitted[0]![0]).toMatchObject({ field: { type: 'state', value: 'counter' } });
  });

  it('emits update:condition on handleValueUpdate', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: BASE_CONDITION, stateFields: [], canRemove: true },
    });
    wrapper.vm.handleValueUpdate({ type: 'literal', value: '42' });
    const emitted = wrapper.emitted('update:condition')!;
    expect(emitted[0]![0]).toMatchObject({ value: { type: 'literal', value: '42' } });
  });

  it('emits update:condition on updateOperator', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: BASE_CONDITION, stateFields: [], canRemove: true },
    });
    wrapper.vm.updateOperator('not_equals');
    const emitted = wrapper.emitted('update:condition')!;
    expect(emitted[0]![0]).toMatchObject({ operator: 'not_equals' });
  });

  it('emits select-event-field on handleEventFieldRequest', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: BASE_CONDITION, stateFields: [], canRemove: true },
    });
    wrapper.vm.handleEventFieldRequest('field');
    const emitted = wrapper.emitted('select-event-field')!;
    expect(emitted[0]![0]).toEqual({ side: 'field' });
  });

  it('starts name editing and saves', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: BASE_CONDITION, stateFields: [], canRemove: true },
    });
    wrapper.vm.startEditingName();
    expect(wrapper.vm.isEditingName).toBe(true);
    wrapper.vm.editableName = 'Renamed';
    wrapper.vm.saveName();
    expect(wrapper.vm.isEditingName).toBe(false);
    const emitted = wrapper.emitted('update:condition')!;
    expect(emitted[0]![0]).toMatchObject({ name: 'Renamed' });
  });

  it('cancels name editing restoring original', () => {
    const wrapper = mountWithPlugins(ConditionItemCard, {
      props: { condition: BASE_CONDITION, stateFields: [], canRemove: true },
    });
    wrapper.vm.startEditingName();
    wrapper.vm.editableName = 'Changed';
    wrapper.vm.cancelEditName();
    expect(wrapper.vm.isEditingName).toBe(false);
    expect(wrapper.vm.editableName).toBe('condition');
  });
});

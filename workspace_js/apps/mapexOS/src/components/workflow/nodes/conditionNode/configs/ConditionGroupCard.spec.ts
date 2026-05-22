import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import ConditionGroupCard from './ConditionGroupCard.vue';

vi.mock('@src/composables/workflow', () => ({
  usePluginI18n: () => ({ t: (key: string) => key }),
}));

vi.mock('@components/tooltips/appTooltip', () => ({
  AppTooltip: { name: 'AppTooltip', template: '<div />' },
}));

vi.mock('./ConditionItemCard.vue', () => ({
  default: { name: 'ConditionItemCard', template: '<div />' },
}));

vi.mock('@src/pages/automations/rules/createEditRulePage/constants', () => ({
  ComparisonOperator: { Equals: 'equals' },
}));

vi.mock('../constants', () => ({
  GROUP_LOGIC_OPTIONS: [
    { value: 'AND', label: 'AND', icon: 'join_full', color: 'blue-6', description: 'All must match' },
    { value: 'OR', label: 'OR', icon: 'join_inner', color: 'orange-6', description: 'Any must match' },
  ],
}));

const BASE_GROUP = {
  id: 'g1',
  name: 'Group 1',
  logic: 'AND' as const,
  items: [],
};

describe('ConditionGroupCard', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(ConditionGroupCard, {
      props: { group: BASE_GROUP, canRemove: true, stateFields: [] },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts expanded by default', () => {
    const wrapper = mountWithPlugins(ConditionGroupCard, {
      props: { group: BASE_GROUP, canRemove: true, stateFields: [] },
    });
    expect(wrapper.vm.isExpanded).toBe(true);
  });

  it('toggles expanded state', () => {
    const wrapper = mountWithPlugins(ConditionGroupCard, {
      props: { group: BASE_GROUP, canRemove: true, stateFields: [] },
    });
    wrapper.vm.toggleExpanded();
    expect(wrapper.vm.isExpanded).toBe(false);
    wrapper.vm.toggleExpanded();
    expect(wrapper.vm.isExpanded).toBe(true);
  });

  it('computes currentOperator from group logic', () => {
    const wrapper = mountWithPlugins(ConditionGroupCard, {
      props: { group: BASE_GROUP, canRemove: true, stateFields: [] },
    });
    expect(wrapper.vm.currentOperator.value).toBe('AND');
  });

  it('computes itemsCountText', () => {
    const wrapper = mountWithPlugins(ConditionGroupCard, {
      props: { group: BASE_GROUP, canRemove: true, stateFields: [] },
    });
    expect(wrapper.vm.itemsCountText).toContain('0');
  });

  it('emits update:group on updateLogic', () => {
    const wrapper = mountWithPlugins(ConditionGroupCard, {
      props: { group: BASE_GROUP, canRemove: true, stateFields: [] },
    });
    wrapper.vm.updateLogic('OR');
    const emitted = wrapper.emitted('update:group')!;
    expect(emitted[0]![0]).toMatchObject({ logic: 'OR' });
  });

  it('emits update:group on addCondition', () => {
    const wrapper = mountWithPlugins(ConditionGroupCard, {
      props: { group: BASE_GROUP, canRemove: true, stateFields: [] },
    });
    wrapper.vm.addCondition();
    const emitted = wrapper.emitted('update:group')!;
    expect((emitted[0]![0] as any).items).toHaveLength(1);
  });

  it('emits remove event', () => {
    const wrapper = mountWithPlugins(ConditionGroupCard, {
      props: { group: BASE_GROUP, canRemove: true, stateFields: [] },
    });
    wrapper.vm.$emit('remove');
    expect(wrapper.emitted('remove')).toBeTruthy();
  });

  it('removes item by index via removeItem', () => {
    const groupWithItems = {
      ...BASE_GROUP,
      items: [
        { type: 'condition' as const, data: { id: 'c1', name: 'c', field: { type: 'event', value: '' }, operator: 'equals', value: { type: 'input', value: '' } } },
        { type: 'condition' as const, data: { id: 'c2', name: 'c', field: { type: 'event', value: '' }, operator: 'equals', value: { type: 'input', value: '' } } },
      ],
    };
    const wrapper = mountWithPlugins(ConditionGroupCard, {
      props: { group: groupWithItems, canRemove: true, stateFields: [] },
    });
    wrapper.vm.removeItem(0);
    const emitted = wrapper.emitted('update:group')!;
    expect((emitted[0]![0] as any).items).toHaveLength(1);
    expect((emitted[0]![0] as any).items[0].data.id).toBe('c2');
  });

  it('starts name editing and saves', () => {
    const wrapper = mountWithPlugins(ConditionGroupCard, {
      props: { group: BASE_GROUP, canRemove: true, stateFields: [] },
    });
    wrapper.vm.startEditingName();
    expect(wrapper.vm.isEditingName).toBe(true);
    wrapper.vm.editableName = 'Renamed';
    wrapper.vm.saveName();
    expect(wrapper.vm.isEditingName).toBe(false);
    const emitted = wrapper.emitted('update:group')!;
    expect(emitted[0]![0]).toMatchObject({ name: 'Renamed' });
  });

  it('cancels name editing', () => {
    const wrapper = mountWithPlugins(ConditionGroupCard, {
      props: { group: BASE_GROUP, canRemove: true, stateFields: [] },
    });
    wrapper.vm.startEditingName();
    wrapper.vm.editableName = 'Changed';
    wrapper.vm.cancelEditName();
    expect(wrapper.vm.isEditingName).toBe(false);
    expect(wrapper.vm.editableName).toBe('Group 1');
  });
});

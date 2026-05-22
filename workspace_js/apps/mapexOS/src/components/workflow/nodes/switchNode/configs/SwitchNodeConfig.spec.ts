import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import SwitchNodeConfig from './SwitchNodeConfig.vue';

vi.mock('@src/composables/workflow', () => ({
  useWorkflowContext: () => ({
    states: { value: [{ field: 'status', type: 'string' }, { field: 'count', type: 'number' }] },
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

vi.mock('../../conditionNode/configs/ConditionGroupCard.vue', () => ({
  default: { name: 'ConditionGroupCard', template: '<div />' },
}));

vi.mock('../../conditionNode/configs/ConditionItemCard.vue', () => ({
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

vi.mock('../../conditionNode/constants', () => ({
  GROUP_LOGIC_OPTIONS: [
    { value: 'AND', label: 'AND', icon: 'join_full', color: 'blue-6', description: 'All must match' },
    { value: 'OR', label: 'OR', icon: 'join_inner', color: 'orange-6', description: 'Any must match' },
    { value: 'NAND', label: 'NAND', icon: 'join_full', color: 'red-6', description: 'Not all match' },
    { value: 'NOR', label: 'NOR', icon: 'join_inner', color: 'red-6', description: 'None match' },
  ],
}));

const BASE_CONFIG: Record<string, unknown> = {};

function makeSwitchCase(id: string, name: string, items: any[] = []) {
  return {
    id,
    name,
    condition: { id: `g_${id}`, name: 'Root', logic: 'AND', items },
  };
}

function makeConditionItem(id: string) {
  return {
    type: 'condition' as const,
    data: {
      id,
      name: 'condition',
      field: { type: 'event', value: '' },
      operator: 'equals',
      value: { type: 'input', value: '' },
    },
  };
}

function makeGroupItem(id: string, subItems: any[] = []) {
  return {
    type: 'group' as const,
    data: {
      id,
      name: 'Group 1',
      logic: 'AND',
      items: subItems,
    },
  };
}

describe('SwitchNodeConfig', () => {
  // ── Initial State ──────────────────────────────────────────────────

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes cases as empty array by default', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.cases).toEqual([]);
  });

  it('computes matchMode as "first" by default', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.matchMode).toBe('first');
  });

  it('computes matchMode from config', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { matchMode: 'all' } },
    });
    expect(wrapper.vm.matchMode).toBe('all');
  });

  it('computes stateFields from workflow context', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.stateFields).toEqual([
      { name: 'status', type: 'string' },
      { name: 'count', type: 'number' },
    ]);
  });

  it('computes activeCase as null when no cases', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.activeCase).toBeNull();
  });

  it('computes activeCaseLogic as AND by default', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.activeCaseLogic).toBe('AND');
  });

  it('computes caseOptions from cases array', () => {
    const cases = [makeSwitchCase('c1', 'Case 1', [makeConditionItem('cond1')]), makeSwitchCase('c2', 'Case 2')];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    expect(wrapper.vm.caseOptions).toEqual([
      { label: 'Case 1', value: 0, itemCount: 1 },
      { label: 'Case 2', value: 1, itemCount: 0 },
    ]);
  });

  it('computes matchModeOptions with two options', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.matchModeOptions).toHaveLength(2);
    expect(wrapper.vm.matchModeOptions[0].value).toBe('first');
    expect(wrapper.vm.matchModeOptions[1].value).toBe('all');
  });

  // ── addCase ────────────────────────────────────────────────────────

  it('emits update:config with new case on addCase', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.addCase();
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted).toHaveLength(1);
    const cases = (emitted[0]![0] as any).cases;
    expect(cases).toHaveLength(1);
    expect(cases[0]).toMatchObject({
      name: 'Case 1',
      condition: { logic: 'AND', items: [] },
    });
  });

  it('sets activeCaseIndex to new case on addCase', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.addCase();
    expect(wrapper.vm.activeCaseIndex).toBe(0);
  });

  it('appends to existing cases on addCase', () => {
    const cases = [makeSwitchCase('c1', 'Case 1')];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.addCase();
    const emitted = wrapper.emitted('update:config')!;
    const newCases = (emitted[0]![0] as any).cases;
    expect(newCases).toHaveLength(2);
    expect(newCases[0].id).toBe('c1');
    expect(newCases[1].name).toBe('Case 2');
  });

  // ── removeCase ─────────────────────────────────────────────────────

  it('does nothing on removeCase when activeCaseIndex is null', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.removeCase();
    expect(wrapper.emitted('update:config')).toBeUndefined();
  });

  it('removes the active case and adjusts index', () => {
    const cases = [makeSwitchCase('c1', 'Case 1'), makeSwitchCase('c2', 'Case 2')];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    wrapper.vm.removeCase();
    const emitted = wrapper.emitted('update:config')!;
    const newCases = (emitted[0]![0] as any).cases;
    expect(newCases).toHaveLength(1);
    expect(newCases[0].id).toBe('c2');
  });

  it('sets activeCaseIndex to null when last case is removed', () => {
    const cases = [makeSwitchCase('c1', 'Case 1')];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    wrapper.vm.removeCase();
    expect(wrapper.vm.activeCaseIndex).toBeNull();
  });

  it('clamps activeCaseIndex when removing last case in list', () => {
    const cases = [makeSwitchCase('c1', 'Case 1'), makeSwitchCase('c2', 'Case 2')];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 1;
    wrapper.vm.removeCase();
    expect(wrapper.vm.activeCaseIndex).toBe(0);
  });

  // ── updateMatchMode ────────────────────────────────────────────────

  it('emits update:config with new matchMode', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateMatchMode('all');
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).matchMode).toBe('all');
  });

  // ── updateActiveCase ───────────────────────────────────────────────

  it('does nothing on updateActiveCase when no active case', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateActiveCase({ name: 'Updated' });
    expect(wrapper.emitted('update:config')).toBeUndefined();
  });

  // ── updateRootLogic ────────────────────────────────────────────────

  it('updates root logic of the active case', () => {
    const cases = [makeSwitchCase('c1', 'Case 1')];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    wrapper.vm.updateRootLogic('OR');
    const emitted = wrapper.emitted('update:config')!;
    const updatedCases = (emitted[0]![0] as any).cases;
    expect(updatedCases[0].condition.logic).toBe('OR');
  });

  it('does nothing on updateRootLogic when no active case', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateRootLogic('OR');
    expect(wrapper.emitted('update:config')).toBeUndefined();
  });

  // ── addCondition ───────────────────────────────────────────────────

  it('adds a condition to the active case', () => {
    const cases = [makeSwitchCase('c1', 'Case 1')];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    wrapper.vm.addCondition();
    const emitted = wrapper.emitted('update:config')!;
    const items = (emitted[0]![0] as any).cases[0].condition.items;
    expect(items).toHaveLength(1);
    expect(items[0].type).toBe('condition');
    expect(items[0].data.operator).toBe('equals');
  });

  it('does nothing on addCondition when no active case', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.addCondition();
    expect(wrapper.emitted('update:config')).toBeUndefined();
  });

  // ── addGroup ───────────────────────────────────────────────────────

  it('adds a sub-group to the active case with one default condition', () => {
    const cases = [makeSwitchCase('c1', 'Case 1')];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    wrapper.vm.addGroup();
    const emitted = wrapper.emitted('update:config')!;
    const items = (emitted[0]![0] as any).cases[0].condition.items;
    expect(items).toHaveLength(1);
    expect(items[0].type).toBe('group');
    expect(items[0].data.items).toHaveLength(1);
    expect(items[0].data.items[0].type).toBe('condition');
  });

  it('does nothing on addGroup when no active case', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.addGroup();
    expect(wrapper.emitted('update:config')).toBeUndefined();
  });

  // ── updateItem ─────────────────────────────────────────────────────

  it('updates item at index in the active case', () => {
    const cond = makeConditionItem('cond1');
    const cases = [makeSwitchCase('c1', 'Case 1', [cond])];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    const updatedItem = {
      type: 'condition' as const,
      data: { ...cond.data, name: 'updated' },
    };
    wrapper.vm.updateItem(0, updatedItem);
    const emitted = wrapper.emitted('update:config')!;
    const items = (emitted[0]![0] as any).cases[0].condition.items;
    expect(items[0].data.name).toBe('updated');
  });

  // ── removeItem ─────────────────────────────────────────────────────

  it('removes item at index from active case', () => {
    const items = [makeConditionItem('c1'), makeConditionItem('c2')];
    const cases = [makeSwitchCase('case1', 'Case 1', items)];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    wrapper.vm.removeItem(0);
    const emitted = wrapper.emitted('update:config')!;
    const resultItems = (emitted[0]![0] as any).cases[0].condition.items;
    expect(resultItems).toHaveLength(1);
    expect(resultItems[0].data.id).toBe('c2');
  });

  // ── findConditionById ──────────────────────────────────────────────

  it('finds condition by ID in root items', () => {
    const cond = makeConditionItem('target');
    const cases = [makeSwitchCase('c1', 'Case 1', [cond])];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    const found = wrapper.vm.findConditionById('target');
    expect(found).toBeDefined();
    expect(found!.id).toBe('target');
  });

  it('finds condition by ID in sub-groups', () => {
    const subCond = makeConditionItem('nested');
    const group = makeGroupItem('g1', [subCond]);
    const cases = [makeSwitchCase('c1', 'Case 1', [group])];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    const found = wrapper.vm.findConditionById('nested');
    expect(found).toBeDefined();
    expect(found!.id).toBe('nested');
  });

  it('returns undefined for non-existent condition ID', () => {
    const cases = [makeSwitchCase('c1', 'Case 1', [makeConditionItem('c1')])];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    expect(wrapper.vm.findConditionById('nonexistent')).toBeUndefined();
  });

  // ── handleEventFieldRequest ────────────────────────────────────────

  it('opens template dialog when no templates are selected', () => {
    const cases = [makeSwitchCase('c1', 'Case 1')];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    wrapper.vm.handleEventFieldRequest({ side: 'field', conditionId: 'c1' });
    expect(wrapper.vm.templateDialogOpen).toBe(true);
    expect(wrapper.vm.fieldSelectorOpen).toBe(false);
  });

  it('opens field selector when templates are already selected', () => {
    const cases = [makeSwitchCase('c1', 'Case 1')];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases, selectedTemplateIds: ['tpl-1'] } },
    });
    wrapper.vm.activeCaseIndex = 0;
    wrapper.vm.handleEventFieldRequest({ side: 'field', conditionId: 'c1' });
    expect(wrapper.vm.fieldSelectorOpen).toBe(true);
    expect(wrapper.vm.templateDialogOpen).toBe(false);
  });

  it('sets activeFieldSide and activeConditionId on handleEventFieldRequest', () => {
    const cases = [makeSwitchCase('c1', 'Case 1')];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases, selectedTemplateIds: ['tpl-1'] } },
    });
    wrapper.vm.handleEventFieldRequest({ side: 'value', conditionId: 'cond-42' });
    expect(wrapper.vm.activeFieldSide).toBe('value');
    expect(wrapper.vm.activeConditionId).toBe('cond-42');
  });

  // ── fieldItems / filteredFieldItems ────────────────────────────────

  it('computes fieldItems from template fields cache', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.templateFieldsCache.set('tpl1', ['temperature', 'humidity']);
    wrapper.vm.templateNamesCache.set('tpl1', 'Weather Template');
    const items = wrapper.vm.fieldItems;
    expect(items).toHaveLength(2);
    expect(items[0]).toEqual({ id: 'tpl1:temperature', path: 'temperature', templateName: 'Weather Template' });
  });

  it('computes filteredFieldItems with search query', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.templateFieldsCache.set('tpl1', ['temperature', 'humidity']);
    wrapper.vm.templateNamesCache.set('tpl1', 'Template');
    wrapper.vm.fieldSearchQuery = 'temp';
    const items = wrapper.vm.filteredFieldItems;
    expect(items).toHaveLength(1);
    expect(items[0].path).toBe('temperature');
  });

  it('returns all fieldItems when search query is empty', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.templateFieldsCache.set('tpl1', ['a', 'b']);
    wrapper.vm.templateNamesCache.set('tpl1', 'T');
    wrapper.vm.fieldSearchQuery = '';
    expect(wrapper.vm.filteredFieldItems).toHaveLength(2);
  });

  // ── handleFieldSearch ──────────────────────────────────────────────

  it('updates fieldSearchQuery on handleFieldSearch', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleFieldSearch('test query');
    expect(wrapper.vm.fieldSearchQuery).toBe('test query');
  });

  // ── handleFieldSelect ──────────────────────────────────────────────

  it('updates condition field on handleFieldSelect for root condition', () => {
    const cond = makeConditionItem('target-cond');
    const cases = [makeSwitchCase('c1', 'Case 1', [cond])];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    wrapper.vm.activeConditionId = 'target-cond';
    wrapper.vm.activeFieldSide = 'field';
    wrapper.vm.handleFieldSelect([{ path: 'event.temperature' }]);
    const emitted = wrapper.emitted('update:config')!;
    const items = (emitted[0]![0] as any).cases[0].condition.items;
    expect(items[0].data.field).toEqual({ type: 'event', value: 'event.temperature', mode: 'dynamic' });
  });

  it('updates condition value on handleFieldSelect for value side', () => {
    const cond = makeConditionItem('target-cond');
    const cases = [makeSwitchCase('c1', 'Case 1', [cond])];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    wrapper.vm.activeConditionId = 'target-cond';
    wrapper.vm.activeFieldSide = 'value';
    wrapper.vm.handleFieldSelect([{ path: 'event.humidity' }]);
    const emitted = wrapper.emitted('update:config')!;
    const items = (emitted[0]![0] as any).cases[0].condition.items;
    expect(items[0].data.value).toEqual({ type: 'event', value: 'event.humidity', mode: 'dynamic' });
  });

  it('updates nested condition in sub-group on handleFieldSelect', () => {
    const subCond = makeConditionItem('nested-cond');
    const group = makeGroupItem('g1', [subCond]);
    const cases = [makeSwitchCase('c1', 'Case 1', [group])];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    wrapper.vm.activeConditionId = 'nested-cond';
    wrapper.vm.activeFieldSide = 'field';
    wrapper.vm.handleFieldSelect([{ path: 'event.pressure' }]);
    const emitted = wrapper.emitted('update:config')!;
    const items = (emitted[0]![0] as any).cases[0].condition.items;
    expect(items[0].data.items[0].data.field).toEqual({ type: 'event', value: 'event.pressure', mode: 'dynamic' });
  });

  it('does nothing on handleFieldSelect with empty array', () => {
    const cases = [makeSwitchCase('c1', 'Case 1')];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    wrapper.vm.handleFieldSelect([]);
    expect(wrapper.emitted('update:config')).toBeUndefined();
  });

  // ── currentOperator ────────────────────────────────────────────────

  it('computes currentOperator from activeCaseLogic', () => {
    const cases = [makeSwitchCase('c1', 'Case 1')];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    expect(wrapper.vm.currentOperator.value).toBe('AND');
  });

  // ── selectedFieldIds ───────────────────────────────────────────────

  it('computes selectedFieldIds as empty when no active condition', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.selectedFieldIds).toEqual([]);
  });

  // ── activeCaseItems ────────────────────────────────────────────────

  it('computes activeCaseItems as empty when no active case', () => {
    const wrapper = mountWithPlugins(SwitchNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.activeCaseItems).toEqual([]);
  });

  it('computes activeCaseItems from active case condition', () => {
    const cond = makeConditionItem('c1');
    const cases = [makeSwitchCase('case1', 'Case 1', [cond])];
    const wrapper = mountWithPlugins(SwitchNodeConfig, {
      props: { config: { cases } },
    });
    wrapper.vm.activeCaseIndex = 0;
    expect(wrapper.vm.activeCaseItems).toHaveLength(1);
    expect(wrapper.vm.activeCaseItems[0].data.id).toBe('c1');
  });
});

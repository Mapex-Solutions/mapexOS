import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import TriggerEventNodeConfig from './TriggerEventNodeConfig.vue';

vi.mock('@src/composables/workflow', () => ({
  usePluginI18n: () => ({ t: (key: string) => key }),
}));

vi.mock('@components/tooltips', () => ({
  AppTooltip: { name: 'AppTooltip', template: '<div />' },
}));

vi.mock('@components/dialogs/common/triggerSelectorDialog', () => ({
  TriggerSelectorDialog: { name: 'TriggerSelectorDialog', template: '<div />' },
}));

vi.mock('@components/dialogs/common/assetTemplateSelectorDialog', () => ({
  AssetTemplateSelectorDialog: { name: 'AssetTemplateSelectorDialog', template: '<div />' },
}));

vi.mock('@components/dialogs/common/genericSelectorDialog', () => ({
  GenericSelectorDialog: { name: 'GenericSelectorDialog', template: '<div />' },
}));

vi.mock('@components/forms/eventFieldInput', () => ({
  EventFieldInput: { name: 'EventFieldInput', template: '<div />' },
}));

vi.mock('@services/mapex', () => ({
  apis: {
    assets: { assetTemplate: { getAvailableFields: vi.fn(), getById: vi.fn() } },
    triggers: { trigger: { getById: vi.fn() } },
  },
}));

vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));

vi.mock('@src/pages/automations/rules/createEditRulePage/utils', () => ({
  extractFieldsFromConfig: vi.fn(() => ({})),
  parseAllFieldVariables: vi.fn(() => ({})),
}));

const BASE_CONFIG: Record<string, unknown> = {};

describe('TriggerEventNodeConfig', () => {
  // ── Initial State ──────────────────────────────────────────────────

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts with triggerDrawerOpen as false', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.triggerDrawerOpen).toBe(false);
  });

  it('starts with templateDrawerOpen as false', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.templateDrawerOpen).toBe(false);
  });

  it('starts with fieldSelectorOpen as false', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.fieldSelectorOpen).toBe(false);
  });

  // ── Computed: parsedVariables ──────────────────────────────────────

  it('computes parsedVariables as empty when no trigger is selected', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.parsedVariables).toEqual({});
  });

  // ── Computed: uniqueVariables ──────────────────────────────────────

  it('computes uniqueVariables as empty when no parsedVariables', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.uniqueVariables).toEqual([]);
  });

  // ── Computed: triggerVariables ──────────────────────────────────────

  it('computes triggerVariables as empty object by default', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.triggerVariables).toEqual({});
  });

  it('computes triggerVariables from config', () => {
    const variables = {
      'device.id': {
        path: 'device.id',
        placeholder: '{{device.id}}',
        fieldKey: 'url',
        value: { type: 'event', value: 'event.deviceId' },
      },
    };
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, {
      props: { config: { variables } },
    });
    expect(wrapper.vm.triggerVariables).toEqual(variables);
  });

  // ── Computed: configuredCount ──────────────────────────────────────

  it('computes configuredCount as 0 when no variables', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.configuredCount).toBe(0);
  });

  it('counts only variables with non-empty values', () => {
    const variables = {
      'a': { path: 'a', placeholder: '', fieldKey: '', value: { type: 'event', value: 'filled' } },
      'b': { path: 'b', placeholder: '', fieldKey: '', value: { type: 'event', value: '' } },
      'c': { path: 'c', placeholder: '', fieldKey: '', value: { type: 'literal', value: 'also filled' } },
    };
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, {
      props: { config: { variables } },
    });
    expect(wrapper.vm.configuredCount).toBe(2);
  });

  // ── Computed: hasTemplates ─────────────────────────────────────────

  it('computes hasTemplates as false when no templates selected', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.hasTemplates).toBe(false);
  });

  it('computes hasTemplates as true when templates are selected', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, {
      props: { config: { selectedTemplateIds: ['tpl1'] } },
    });
    expect(wrapper.vm.hasTemplates).toBe(true);
  });

  // ── Computed: fieldItems / filteredFieldItems ──────────────────────

  it('computes fieldItems from template fields cache', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.templateFieldsCache.set('tpl1', ['temp', 'humidity']);
    wrapper.vm.templateNamesCache.set('tpl1', 'Weather');
    expect(wrapper.vm.fieldItems).toHaveLength(2);
    expect(wrapper.vm.fieldItems[0]).toEqual({ id: 'tpl1:temp', path: 'temp', templateName: 'Weather' });
  });

  it('filters fieldItems by search query', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.templateFieldsCache.set('tpl1', ['temperature', 'humidity']);
    wrapper.vm.templateNamesCache.set('tpl1', 'T');
    wrapper.vm.fieldSearchQuery = 'hum';
    expect(wrapper.vm.filteredFieldItems).toHaveLength(1);
    expect(wrapper.vm.filteredFieldItems[0].path).toBe('humidity');
  });

  it('returns all fieldItems when search is empty', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.templateFieldsCache.set('tpl1', ['a', 'b']);
    wrapper.vm.templateNamesCache.set('tpl1', 'T');
    wrapper.vm.fieldSearchQuery = '  ';
    expect(wrapper.vm.filteredFieldItems).toHaveLength(2);
  });

  // ── Computed: selectedFieldIds ─────────────────────────────────────

  it('computes selectedFieldIds as empty when no active variable', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.selectedFieldIds).toEqual([]);
  });

  // ── getTriggerTypeIcon ─────────────────────────────────────────────

  it('returns correct icon for known trigger types', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.getTriggerTypeIcon('http')).toBe('http');
    expect(wrapper.vm.getTriggerTypeIcon('mqtt')).toBe('router');
    expect(wrapper.vm.getTriggerTypeIcon('rabbitmq')).toBe('cloud_queue');
    expect(wrapper.vm.getTriggerTypeIcon('nats')).toBe('cloud');
    expect(wrapper.vm.getTriggerTypeIcon('websocket')).toBe('cable');
    expect(wrapper.vm.getTriggerTypeIcon('email')).toBe('email');
    expect(wrapper.vm.getTriggerTypeIcon('teams')).toBe('groups');
    expect(wrapper.vm.getTriggerTypeIcon('slack')).toBe('chat');
  });

  it('returns default icon for unknown trigger type', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.getTriggerTypeIcon('unknown')).toBe('notifications_active');
    expect(wrapper.vm.getTriggerTypeIcon(undefined)).toBe('notifications_active');
  });

  // ── getCategoryColor ───────────────────────────────────────────────

  it('returns purple-6 for communication category', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.getCategoryColor('communication')).toBe('purple-6');
  });

  it('returns blue-6 for other categories', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.getCategoryColor('data')).toBe('blue-6');
    expect(wrapper.vm.getCategoryColor(undefined)).toBe('blue-6');
  });

  // ── getVariableValue ───────────────────────────────────────────────

  it('returns default FieldSourceValue when variable not in config', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.getVariableValue('some.path')).toEqual({ type: 'event', value: '' });
  });

  it('returns stored value for configured variable', () => {
    const variables = {
      'device.id': {
        path: 'device.id',
        placeholder: '{{device.id}}',
        fieldKey: 'url',
        value: { type: 'event', value: 'event.deviceId' },
      },
    };
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, {
      props: { config: { variables } },
    });
    expect(wrapper.vm.getVariableValue('device.id')).toEqual({ type: 'event', value: 'event.deviceId' });
  });

  // ── handleTriggerSelect ────────────────────────────────────────────

  it('emits update:config with trigger info on handleTriggerSelect', async () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    const trigger = { id: 'trig-1', name: 'My Trigger', triggerType: 'http', config: {} };
    await wrapper.vm.handleTriggerSelect(trigger as any);
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({
      triggerId: 'trig-1',
      triggerName: 'My Trigger',
      triggerType: 'http',
      variables: {},
    });
  });

  it('sets selectedTriggerFull on handleTriggerSelect', async () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    const trigger = { id: 'trig-1', name: 'T', triggerType: 'http', config: { url: 'http://test' } };
    await wrapper.vm.handleTriggerSelect(trigger as any);
    expect(wrapper.vm.selectedTriggerFull).toEqual(trigger);
  });

  // ── clearTrigger ───────────────────────────────────────────────────

  it('clears trigger and variables on clearTrigger', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, {
      props: { config: { triggerId: 'trig-1', triggerName: 'T', triggerType: 'http', variables: { a: {} } } },
    });
    wrapper.vm.clearTrigger();
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({
      triggerId: undefined,
      triggerName: undefined,
      triggerType: undefined,
      variables: {},
      selectedTemplateIds: [],
    });
  });

  it('resets internal state on clearTrigger', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, {
      props: { config: { selectedTemplateIds: ['tpl1'] } },
    });
    wrapper.vm.templateFieldsCache.set('tpl1', ['f1']);
    wrapper.vm.templateNamesCache.set('tpl1', 'T');
    wrapper.vm.clearTrigger();
    expect(wrapper.vm.selectedTriggerFull).toBeNull();
    expect(wrapper.vm.selectedTemplateIds).toEqual([]);
    expect(wrapper.vm.templateFieldsCache.size).toBe(0);
    expect(wrapper.vm.templateNamesCache.size).toBe(0);
  });

  // ── handleVariableUpdate ───────────────────────────────────────────

  it('emits update:config with updated variable on handleVariableUpdate', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleVariableUpdate('device.id', '{{device.id}}', 'url', { type: 'event', value: 'event.devId' });
    const emitted = wrapper.emitted('update:config')!;
    const vars = (emitted[0]![0] as any).variables;
    expect(vars['device.id']).toEqual({
      path: 'device.id',
      placeholder: '{{device.id}}',
      fieldKey: 'url',
      value: { type: 'event', value: 'event.devId' },
    });
  });

  // ── handleOpenEventSelector ────────────────────────────────────────

  it('opens field selector on handleOpenEventSelector', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleOpenEventSelector('device.id');
    expect(wrapper.vm.activeVariablePath).toBe('device.id');
    expect(wrapper.vm.fieldSelectorOpen).toBe(true);
    expect(wrapper.vm.fieldSearchQuery).toBe('');
  });

  // ── handleOpenTemplateSelector ─────────────────────────────────────

  it('opens template drawer on handleOpenTemplateSelector', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleOpenTemplateSelector('device.id');
    expect(wrapper.vm.activeVariablePath).toBe('device.id');
    expect(wrapper.vm.templateDrawerOpen).toBe(true);
  });

  // ── handleFieldSelect ──────────────────────────────────────────────

  it('does nothing on handleFieldSelect with empty items', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleFieldSelect([]);
    expect(wrapper.emitted('update:config')).toBeUndefined();
  });

  it('does nothing on handleFieldSelect when no active variable path', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.activeVariablePath = null;
    wrapper.vm.handleFieldSelect([{ path: 'event.temp' }]);
    expect(wrapper.emitted('update:config')).toBeUndefined();
  });

  // ── handleFieldSearch ──────────────────────────────────────────────

  it('updates fieldSearchQuery on handleFieldSearch', () => {
    const wrapper = mountWithPlugins(TriggerEventNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleFieldSearch('temperature');
    expect(wrapper.vm.fieldSearchQuery).toBe('temperature');
  });
});

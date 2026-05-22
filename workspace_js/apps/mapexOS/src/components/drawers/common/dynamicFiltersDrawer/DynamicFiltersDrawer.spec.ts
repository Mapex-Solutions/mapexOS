import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins, createMockTranslations } from '@src/test/helpers';
import DynamicFiltersDrawer from './DynamicFiltersDrawer.vue';

vi.mock('@composables/i18n/components/drawers/dynamicFiltersDrawer', () => ({
  useDynamicFiltersDrawerTranslations: () => createMockTranslations({
    operators: new Proxy({}, {
      get: (_t, prop) => ({ value: String(prop) }),
    }),
    fieldTypeHeaders: new Proxy({}, {
      get: (_t, prop) => ({ value: String(prop) }),
    }),
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

vi.mock('@services/mapex', () => ({
  apis: {
    assets: {
      asset: {
        list: vi.fn().mockResolvedValue({ items: [] }),
        getById: vi.fn().mockResolvedValue({ id: 'a-1', assetTemplateId: 'at-1' }),
      },
      assetTemplate: {
        getById: vi.fn().mockResolvedValue({ dynamicFields: [] }),
      },
    },
    businessRules: {
      list: vi.fn().mockResolvedValue({ items: [] }),
      getById: vi.fn().mockResolvedValue({ assetIds: [] }),
    },
  },
}));

describe('DynamicFiltersDrawer', () => {
  const defaultProps = {
    modelValue: true,
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('initializes sourceType as asset', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    expect(wrapper.vm.sourceType).toBe('asset');
  });

  it('initializes selectedSourceId as null', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    expect(wrapper.vm.selectedSourceId).toBeNull();
  });

  it('initializes activeFilters as empty', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    expect(wrapper.vm.activeFilters).toEqual([]);
  });

  it('initializes availableFields as empty', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    expect(wrapper.vm.availableFields).toEqual([]);
  });

  it('computes sourceOptions with 2 entries', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    expect(wrapper.vm.sourceOptions).toHaveLength(2);
  });

  it('computes allFieldsAdded as false when no fields', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    expect(wrapper.vm.allFieldsAdded).toBe(false);
  });

  it('returns correct field type icons', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.getFieldTypeIcon('number')).toBe('tag');
    expect(vm.getFieldTypeIcon('string')).toBe('text_fields');
    expect(vm.getFieldTypeIcon('boolean')).toBe('toggle_on');
    expect(vm.getFieldTypeIcon('date')).toBe('event');
    expect(vm.getFieldTypeIcon('unknown')).toBe('help_outline');
  });

  it('returns correct field type colors', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.getFieldTypeColor('number')).toBe('blue-7');
    expect(vm.getFieldTypeColor('string')).toBe('teal-7');
    expect(vm.getFieldTypeColor('boolean')).toBe('orange-7');
    expect(vm.getFieldTypeColor('date')).toBe('purple-7');
  });

  it('returns correct input types', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.getInputType('number')).toBe('number');
    expect(vm.getInputType('date')).toBe('datetime-local');
    expect(vm.getInputType('string')).toBe('text');
  });

  it('adds filter from available fields', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    const field = { key: 'temp', label: 'temperature', type: 'number', originalType: 'number', fieldId: 0 };
    vm.addFilter(field);
    expect(vm.activeFilters).toHaveLength(1);
    expect(vm.activeFilters[0].key).toBe('temp');
  });

  it('removes filter by key', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.activeFilters = [{ key: 'temp', label: 'temperature', type: 'number', operator: 'eq', value: '', fieldId: 0 }];
    vm.removeFilter('temp');
    expect(vm.activeFilters).toHaveLength(0);
  });

  it('resets all state on handleReset', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.selectedSourceId = 'test';
    vm.activeFilters = [{ key: 'temp' }];
    vm.handleReset();
    expect(vm.selectedSourceId).toBeNull();
    expect(vm.activeFilters).toHaveLength(0);
    expect(vm.sourceType).toBe('asset');
    expect(wrapper.emitted('reset')).toBeTruthy();
  });

  it('resets selection on handleSourceChange', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.selectedSourceId = 'test';
    vm.resolvedTemplateId = 'tpl-1';
    vm.handleSourceChange();
    expect(vm.selectedSourceId).toBeNull();
    expect(vm.resolvedTemplateId).toBeUndefined();
    expect(vm.availableFields).toHaveLength(0);
  });

  it('handles Enter key to apply when resolvedTemplateId exists', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.resolvedTemplateId = 'tpl-1';
    vm.selectedSourceId = 'src-1';
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter' }));
    expect(wrapper.emitted('apply')).toBeTruthy();
  });

  it('does not apply when resolvedTemplateId is undefined', () => {
    const wrapper = mountWithPlugins(DynamicFiltersDrawer, { props: defaultProps });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter' }));
    expect(wrapper.emitted('apply')).toBeFalsy();
  });
});

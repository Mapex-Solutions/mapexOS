import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins, createMockTranslations } from '@src/test/helpers';
import type { FilterListItem } from './interfaces';

// Mock @components/drawers barrel to avoid transitive monaco-editor resolution
vi.mock('@components/drawers', () => ({
  UserSelectorDrawer: { name: 'UserSelectorDrawer', template: '<div />' },
}));

// Mock the i18n composable
vi.mock('@composables/i18n', () => ({
  useListFilterTranslations: () => ({
    listFilter: createMockTranslations({
      title: { value: 'Filters' },
      active: { value: 'active' },
      clear: { value: 'Clear' },
      apply: { value: 'Apply' },
      hint: { value: 'Click to expand filters' },
      includeChildren: { value: 'Include Children' },
      includeChildrenTooltip: { value: 'Include sub-items' },
    }),
  }),
}));

import ListFilter from './ListFilter.vue';

const baseItems: FilterListItem[] = [
  { key: 'name', type: 'input', label: 'Name', icon: 'search' },
  {
    key: 'status',
    type: 'select',
    label: 'Status',
    icon: 'toggle_on',
    options: [
      { label: 'Active', value: 'active' },
      { label: 'Inactive', value: 'inactive' },
    ],
  },
  { key: 'tags', type: 'multiselect', label: 'Tags', icon: 'label', options: [] },
];

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(ListFilter, {
    props: {
      items: baseItems,
      ...overrides,
    },
  });
}

describe('ListFilter', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('isExpanded starts as false', () => {
    const wrapper = factory();
    expect(wrapper.vm.isExpanded).toBe(false);
  });

  it('toggle flips isExpanded', () => {
    const wrapper = factory();
    wrapper.vm.toggle();
    expect(wrapper.vm.isExpanded).toBe(true);
    wrapper.vm.toggle();
    expect(wrapper.vm.isExpanded).toBe(false);
  });

  it('hasActiveFilters is false when all filters are empty', () => {
    const wrapper = factory();
    expect(wrapper.vm.hasActiveFilters).toBe(false);
  });

  it('hasActiveFilters is true when an input filter has a value', () => {
    const wrapper = factory();
    wrapper.vm.filterValues.name = 'test';
    expect(wrapper.vm.hasActiveFilters).toBe(true);
  });

  it('activeFiltersCount counts active filters correctly', () => {
    const wrapper = factory();
    expect(wrapper.vm.activeFiltersCount).toBe(0);
    wrapper.vm.filterValues.name = 'test';
    expect(wrapper.vm.activeFiltersCount).toBe(1);
    wrapper.vm.filterValues.status = 'active';
    expect(wrapper.vm.activeFiltersCount).toBe(2);
  });

  it('clearFilters resets all filter values', () => {
    const wrapper = factory();
    wrapper.vm.filterValues.name = 'test';
    wrapper.vm.filterValues.status = 'active';
    wrapper.vm.clearFilters();
    expect(wrapper.vm.filterValues.name).toBe('');
    expect(wrapper.vm.filterValues.status).toBeNull();
    expect(wrapper.vm.filterValues.tags).toEqual([]);
  });

  it('clearFilters emits reset event', () => {
    const wrapper = factory();
    wrapper.vm.clearFilters();
    expect(wrapper.emitted('reset')).toBeTruthy();
  });

  it('applyFilters emits apply with current values', () => {
    const wrapper = factory();
    wrapper.vm.filterValues.name = 'sensor';
    wrapper.vm.applyFilters();
    const emitted = wrapper.emitted('apply');
    expect(emitted).toBeTruthy();
    expect(emitted![0]![0]).toMatchObject({ name: 'sensor' });
  });

  it('applyFilters includes includeChildren when showIncludeChildren is true', () => {
    const wrapper = factory({ showIncludeChildren: true, includeChildrenInitial: true });
    wrapper.vm.applyFilters();
    const emitted = wrapper.emitted('apply');
    expect(emitted![0]![0]).toHaveProperty('includeChildren', true);
  });

  it('removeFilter clears a specific filter', () => {
    const wrapper = factory();
    wrapper.vm.filterValues.name = 'test';
    wrapper.vm.removeFilter('name');
    expect(wrapper.vm.filterValues.name).toBe('');
  });

  it('removeFilter clears select filter to null', () => {
    const wrapper = factory();
    wrapper.vm.filterValues.status = 'active';
    wrapper.vm.removeFilter('status');
    expect(wrapper.vm.filterValues.status).toBeNull();
  });

  it('removeFilter clears multiselect filter to empty array', () => {
    const wrapper = factory();
    wrapper.vm.filterValues.tags = ['a', 'b'];
    wrapper.vm.removeFilter('tags');
    expect(wrapper.vm.filterValues.tags).toEqual([]);
  });

  it('activeFilterChips returns chips for active filters', () => {
    const wrapper = factory();
    wrapper.vm.filterValues.name = 'sensor';
    wrapper.vm.filterValues.status = 'active';
    const chips = wrapper.vm.activeFilterChips;
    expect(chips).toHaveLength(2);
    expect(chips[0]).toMatchObject({ key: 'name', value: 'sensor' });
    expect(chips[1]).toMatchObject({ key: 'status', value: 'Active' });
  });

  it('handleFieldChange emits fieldChange for watched fields', () => {
    const wrapper = factory({ watchFields: ['status'] });
    wrapper.vm.handleFieldChange('status', 'active');
    expect(wrapper.emitted('fieldChange')).toBeTruthy();
    expect(wrapper.emitted('fieldChange')![0]).toEqual(['status', 'active']);
  });

  it('handleFieldChange does not emit for non-watched fields', () => {
    const wrapper = factory({ watchFields: ['status'] });
    wrapper.vm.handleFieldChange('name', 'test');
    expect(wrapper.emitted('fieldChange')).toBeFalsy();
  });

  it('initializes filter values based on type', () => {
    const wrapper = factory();
    expect(wrapper.vm.filterValues.name).toBe('');
    expect(wrapper.vm.filterValues.status).toBeNull();
    expect(wrapper.vm.filterValues.tags).toEqual([]);
  });

  it('uses defaultValue when provided', () => {
    const itemsWithDefault: FilterListItem[] = [
      { key: 'name', type: 'input', label: 'Name', icon: 'search', defaultValue: 'default-name' },
    ];
    const wrapper = factory({ items: itemsWithDefault });
    expect(wrapper.vm.filterValues.name).toBe('default-name');
  });
});

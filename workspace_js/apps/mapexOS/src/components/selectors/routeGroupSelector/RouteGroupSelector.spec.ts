import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins, createMockTranslations } from '@src/test/helpers';
import RouteGroupSelector from './RouteGroupSelector.vue';

/** Mock composables and services */
vi.mock('@src/composables/i18n/components/selectors/useRouteGroupSelectorTranslations', () => ({
  useRouteGroupSelectorTranslations: () => createMockTranslations(),
}));

vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));

vi.mock('@services/mapex', () => ({
  apis: {
    router: {
      routegroup: {
        list: vi.fn().mockResolvedValue({ items: [], pagination: { totalItems: 0, totalPages: 1 } }),
      },
    },
  },
}));

vi.mock('@components/chips', () => ({
  SelectableChip: { template: '<span />' },
}));

describe('RouteGroupSelector', () => {
  const defaultProps = {
    modelValue: [] as string[],
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(RouteGroupSelector, { props: { ...defaultProps } });
    expect(wrapper.exists()).toBe(true);
  });

  it('initializes filters with empty defaults', () => {
    const wrapper = mountWithPlugins(RouteGroupSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.filters.search).toBe('');
    expect(wrapper.vm.filters.isSystem).toBeNull();
    expect(wrapper.vm.filters.isTemplate).toBeNull();
  });

  it('templateTypeOptions has 3 options', () => {
    const wrapper = mountWithPlugins(RouteGroupSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.templateTypeOptions).toHaveLength(3);
  });

  it('templateSourceOptions has 3 options', () => {
    const wrapper = mountWithPlugins(RouteGroupSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.templateSourceOptions).toHaveLength(3);
  });

  it('routeGroups is initially empty', () => {
    const wrapper = mountWithPlugins(RouteGroupSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.routeGroups).toEqual([]);
  });

  it('selectedRouteGroups is initially empty', () => {
    const wrapper = mountWithPlugins(RouteGroupSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.selectedRouteGroups).toEqual([]);
  });

  it('totalRouteGroups is initially 0', () => {
    const wrapper = mountWithPlugins(RouteGroupSelector, { props: { ...defaultProps } });
    expect(wrapper.vm.totalRouteGroups).toBe(0);
  });

  it('toggleRouteGroup adds route group to selection', () => {
    const wrapper = mountWithPlugins(RouteGroupSelector, { props: { ...defaultProps } });
    const rg = { id: 'rg-1', name: 'Route Group 1', enabled: true } as any;
    wrapper.vm.toggleRouteGroup(rg);
    expect(wrapper.vm.selectedRouteGroups).toHaveLength(1);
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:selectedRouteGroups')).toBeTruthy();
  });

  it('toggleRouteGroup removes route group when already selected', () => {
    const wrapper = mountWithPlugins(RouteGroupSelector, { props: { ...defaultProps } });
    const rg = { id: 'rg-1', name: 'Route Group 1', enabled: true } as any;
    wrapper.vm.toggleRouteGroup(rg);
    wrapper.vm.toggleRouteGroup(rg);
    expect(wrapper.vm.selectedRouteGroups).toHaveLength(0);
  });

  it('isSelected returns false for unselected group', () => {
    const wrapper = mountWithPlugins(RouteGroupSelector, { props: { ...defaultProps } });
    const rg = { id: 'rg-1', name: 'RG1' } as any;
    expect(wrapper.vm.isSelected(rg)).toBe(false);
  });

  it('isSelected returns true for selected group', () => {
    const wrapper = mountWithPlugins(RouteGroupSelector, { props: { ...defaultProps } });
    const rg = { id: 'rg-1', name: 'RG1' } as any;
    wrapper.vm.toggleRouteGroup(rg);
    expect(wrapper.vm.isSelected(rg)).toBe(true);
  });

  it('onFilterChange resets currentPage to 1', () => {
    const wrapper = mountWithPlugins(RouteGroupSelector, { props: { ...defaultProps } });
    wrapper.vm.currentPage = 5;
    wrapper.vm.onFilterChange();
    expect(wrapper.vm.currentPage).toBe(1);
  });

  it('multiple selections emit correct IDs array', () => {
    const wrapper = mountWithPlugins(RouteGroupSelector, { props: { ...defaultProps } });
    wrapper.vm.toggleRouteGroup({ id: 'rg-1', name: 'RG1' } as any);
    wrapper.vm.toggleRouteGroup({ id: 'rg-2', name: 'RG2' } as any);
    const lastEmit = wrapper.emitted('update:modelValue')!;
    const lastValue = lastEmit[lastEmit.length - 1]![0] as string[];
    expect(lastValue).toEqual(['rg-1', 'rg-2']);
  });
});

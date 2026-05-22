import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins, createMockTranslations } from '@src/test/helpers';
import RouteGroupSelectorDrawer from './RouteGroupSelectorDrawer.vue';

vi.mock('@composables/routing/routers', () => ({
  useRouters: () => ({
    routeGroups: { value: [] },
    isLoading: { value: false },
    isLoadingMore: { value: false },
    filters: { value: { name: '', isSystem: undefined, isTemplate: undefined } },
    pagination: { value: { currentPage: 1, totalPages: 1, totalItems: 0 } },
    fetchRouteGroups: vi.fn(),
  }),
}));

vi.mock('@src/composables/i18n/components/selectors/useRouteGroupSelectorTranslations', () => ({
  useRouteGroupSelectorTranslations: () => createMockTranslations(),
}));

describe('RouteGroupSelectorDrawer', () => {
  const defaultProps = {
    modelValue: true,
    selectedRouteGroupIds: [],
    multiSelect: true,
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(RouteGroupSelectorDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('syncs showDialog with modelValue prop', () => {
    const wrapper = mountWithPlugins(RouteGroupSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.showDialog).toBe(true);
  });

  it('emits update:modelValue when showDialog changes', async () => {
    const wrapper = mountWithPlugins(RouteGroupSelectorDrawer, { props: defaultProps });
    wrapper.vm.showDialog = false;
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('initializes selectedRouteGroups as empty', () => {
    const wrapper = mountWithPlugins(RouteGroupSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.selectedRouteGroups).toEqual([]);
  });

  it('identifies selected route group via isSelected', () => {
    const wrapper = mountWithPlugins(RouteGroupSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.selectedRouteGroups = [{ id: 'rg-1', name: 'Test' }];
    expect(vm.isSelected({ id: 'rg-1' })).toBe(true);
    expect(vm.isSelected({ id: 'rg-2' })).toBe(false);
  });

  it('toggles route group selection in multi-select mode', () => {
    const wrapper = mountWithPlugins(RouteGroupSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    const rg = { id: 'rg-1', name: 'Test' };
    vm.toggleRouteGroup(rg);
    expect(vm.selectedRouteGroups).toHaveLength(1);
    vm.toggleRouteGroup(rg);
    expect(vm.selectedRouteGroups).toHaveLength(0);
  });

  it('selects and closes in single-select mode', async () => {
    const wrapper = mountWithPlugins(RouteGroupSelectorDrawer, {
      props: { ...defaultProps, multiSelect: false },
    });
    const vm = wrapper.vm;
    const rg = { id: 'rg-1', name: 'Test' };
    vm.toggleRouteGroup(rg);
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('select')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });

  it('emits select on confirmSelection in multi-select', () => {
    const wrapper = mountWithPlugins(RouteGroupSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.selectedRouteGroups = [{ id: 'rg-1' }];
    vm.confirmSelection();
    expect(wrapper.emitted('select')).toBeTruthy();
  });

  it('emits cancel on handleCancel', () => {
    const wrapper = mountWithPlugins(RouteGroupSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.handleCancel();
    expect(wrapper.emitted('cancel')).toBeTruthy();
  });

  it('computes totalRouteGroups from pagination', () => {
    const wrapper = mountWithPlugins(RouteGroupSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.totalRouteGroups).toBe(0);
  });
});

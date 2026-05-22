import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import RoleMultiSelectorDrawer from './RoleMultiSelectorDrawer.vue';

vi.mock('@services/mapex', () => ({
  apis: {
    mapexOS: {
      roles: {
        list: vi.fn().mockResolvedValue({ items: [], pagination: { totalPages: 1, totalItems: 0 } }),
      },
    },
  },
}));

vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));

describe('RoleMultiSelectorDrawer', () => {
  const defaultProps = {
    modelValue: true,
    selectedRoleIds: [],
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(RoleMultiSelectorDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('syncs showDialog with modelValue prop', () => {
    const wrapper = mountWithPlugins(RoleMultiSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.showDialog).toBe(true);
  });

  it('emits update:modelValue when showDialog changes', async () => {
    const wrapper = mountWithPlugins(RoleMultiSelectorDrawer, { props: defaultProps });
    wrapper.vm.showDialog = false;
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('initializes filter state with undefined defaults', () => {
    const wrapper = mountWithPlugins(RoleMultiSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.filters.name).toBeUndefined();
    expect(wrapper.vm.filters.isSystem).toBeUndefined();
    expect(wrapper.vm.filters.isTemplate).toBeUndefined();
  });

  it('initializes pagination state', () => {
    const wrapper = mountWithPlugins(RoleMultiSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.currentPage).toBe(1);
    expect(wrapper.vm.perPage).toBe(15);
  });

  it('computes typeOptions', () => {
    const wrapper = mountWithPlugins(RoleMultiSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.typeOptions).toHaveLength(3);
  });

  it('computes templateOptions', () => {
    const wrapper = mountWithPlugins(RoleMultiSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.templateOptions).toHaveLength(3);
  });

  it('identifies selected role via isSelected', () => {
    const wrapper = mountWithPlugins(RoleMultiSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.selectedRoles = [{ id: 'r-1', name: 'Admin' }];
    expect(vm.isSelected({ id: 'r-1' })).toBe(true);
    expect(vm.isSelected({ id: 'r-2' })).toBe(false);
  });

  it('toggles role selection', () => {
    const wrapper = mountWithPlugins(RoleMultiSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    const role = { id: 'r-1', name: 'Admin' };
    vm.toggleRoleSelection(role);
    expect(vm.selectedRoles).toHaveLength(1);
    vm.toggleRoleSelection(role);
    expect(vm.selectedRoles).toHaveLength(0);
  });

  it('emits confirm with selected roles on handleConfirm', async () => {
    const wrapper = mountWithPlugins(RoleMultiSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.selectedRoles = [{ id: 'r-1', name: 'Admin' }];
    vm.handleConfirm();
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('confirm')).toBeTruthy();
    expect(wrapper.emitted('confirm')![0]![0]).toHaveLength(1);
  });

  it('computes canConfirm based on selection', () => {
    const wrapper = mountWithPlugins(RoleMultiSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.canConfirm).toBe(false);
    vm.selectedRoles = [{ id: 'r-1' }];
    expect(vm.canConfirm).toBe(true);
  });

  it('emits cancel on handleCancel', () => {
    const wrapper = mountWithPlugins(RoleMultiSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.handleCancel();
    expect(wrapper.emitted('cancel')).toBeTruthy();
  });

  it('handles ESC key when drawer is open', () => {
    const wrapper = mountWithPlugins(RoleMultiSelectorDrawer, { props: defaultProps });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('cancel')).toBeTruthy();
  });
});

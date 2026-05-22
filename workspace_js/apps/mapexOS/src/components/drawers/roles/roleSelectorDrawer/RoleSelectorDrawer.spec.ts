import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import RoleSelectorDrawer from './RoleSelectorDrawer.vue';

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

describe('RoleSelectorDrawer', () => {
  const defaultProps = {
    modelValue: true,
    selectedRoleId: null,
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(RoleSelectorDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('syncs showDialog with modelValue prop', () => {
    const wrapper = mountWithPlugins(RoleSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.showDialog).toBe(true);
  });

  it('emits update:modelValue when showDialog changes', async () => {
    const wrapper = mountWithPlugins(RoleSelectorDrawer, { props: defaultProps });
    wrapper.vm.showDialog = false;
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('initializes filter state with undefined defaults', () => {
    const wrapper = mountWithPlugins(RoleSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.filters.name).toBeUndefined();
    expect(wrapper.vm.filters.enabled).toBeUndefined();
    expect(wrapper.vm.filters.isTemplate).toBeUndefined();
  });

  it('computes statusOptions', () => {
    const wrapper = mountWithPlugins(RoleSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.statusOptions).toHaveLength(3);
  });

  it('computes templateOptions', () => {
    const wrapper = mountWithPlugins(RoleSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.templateOptions).toHaveLength(3);
  });

  it('identifies selected role via isSelected', () => {
    const wrapper = mountWithPlugins(RoleSelectorDrawer, {
      props: { ...defaultProps, selectedRoleId: 'r-1' },
    });
    const vm = wrapper.vm;
    expect(vm.isSelected({ id: 'r-1' })).toBe(true);
    expect(vm.isSelected({ id: 'r-2' })).toBe(false);
  });

  it('emits select and closes on selectRole', async () => {
    const wrapper = mountWithPlugins(RoleSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    const role = { id: 'r-1', name: 'Admin' };
    vm.selectRole(role);
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('select')).toBeTruthy();
    expect(wrapper.emitted('select')![0]).toEqual([role]);
  });

  it('emits cancel on handleCancel', () => {
    const wrapper = mountWithPlugins(RoleSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.handleCancel();
    expect(wrapper.emitted('cancel')).toBeTruthy();
  });

  it('handles ESC key when drawer is open', () => {
    const wrapper = mountWithPlugins(RoleSelectorDrawer, { props: defaultProps });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('cancel')).toBeTruthy();
  });

  it('does not handle ESC key when drawer is closed', () => {
    const wrapper = mountWithPlugins(RoleSelectorDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('cancel')).toBeFalsy();
  });

  it('initializes pagination state', () => {
    const wrapper = mountWithPlugins(RoleSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.currentPage).toBe(1);
    expect(wrapper.vm.totalPages).toBe(1);
    expect(wrapper.vm.perPage).toBe(15);
  });
});

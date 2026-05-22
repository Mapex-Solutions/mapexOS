import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import UserSelectorDrawer from './UserSelectorDrawer.vue';

vi.mock('@services/mapex', () => ({
  apis: {
    mapexOS: {
      users: {
        list: vi.fn().mockResolvedValue({ items: [], pagination: { totalPages: 1, totalItems: 0 } }),
      },
    },
  },
}));

vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));

describe('UserSelectorDrawer', () => {
  const defaultProps = {
    modelValue: true,
    selectedUserId: null,
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(UserSelectorDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('syncs showDialog with modelValue prop', () => {
    const wrapper = mountWithPlugins(UserSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.showDialog).toBe(true);
  });

  it('emits update:modelValue when showDialog changes', async () => {
    const wrapper = mountWithPlugins(UserSelectorDrawer, { props: defaultProps });
    wrapper.vm.showDialog = false;
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('initializes filter state with undefined defaults', () => {
    const wrapper = mountWithPlugins(UserSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.filters.email).toBeUndefined();
    expect(wrapper.vm.filters.firstName).toBeUndefined();
    expect(wrapper.vm.filters.lastName).toBeUndefined();
    expect(wrapper.vm.filters.enabled).toBeUndefined();
  });

  it('computes statusOptions', () => {
    const wrapper = mountWithPlugins(UserSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.statusOptions).toHaveLength(3);
  });

  it('identifies selected user via isSelected', () => {
    const wrapper = mountWithPlugins(UserSelectorDrawer, {
      props: { ...defaultProps, selectedUserId: 'u-1' },
    });
    const vm = wrapper.vm;
    expect(vm.isSelected({ id: 'u-1' })).toBe(true);
    expect(vm.isSelected({ id: 'u-2' })).toBe(false);
  });

  it('emits select and closes on selectUser', async () => {
    const wrapper = mountWithPlugins(UserSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    const user = { id: 'u-1', firstName: 'John', lastName: 'Doe' };
    vm.selectUser(user);
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('select')).toBeTruthy();
    expect(wrapper.emitted('select')![0]).toEqual([user]);
  });

  it('emits cancel on handleCancel', () => {
    const wrapper = mountWithPlugins(UserSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.handleCancel();
    expect(wrapper.emitted('cancel')).toBeTruthy();
  });

  it('handles ESC key when drawer is open', () => {
    const wrapper = mountWithPlugins(UserSelectorDrawer, { props: defaultProps });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('cancel')).toBeTruthy();
  });

  it('does not handle ESC key when drawer is closed', () => {
    const wrapper = mountWithPlugins(UserSelectorDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('cancel')).toBeFalsy();
  });

  it('computes display name correctly', () => {
    const wrapper = mountWithPlugins(UserSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.getUserDisplayName({ firstName: 'John', lastName: 'Doe' })).toBe('John Doe');
    expect(vm.getUserDisplayName({ email: 'john@test.com' })).toBe('john@test.com');
    expect(vm.getUserDisplayName({})).toBe('Unnamed User');
  });

  it('computes initials correctly', () => {
    const wrapper = mountWithPlugins(UserSelectorDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.getUserInitials({ firstName: 'John', lastName: 'Doe' })).toBe('JD');
    expect(vm.getUserInitials({ firstName: 'John' })).toBe('JO');
    expect(vm.getUserInitials({ email: 'john@test.com' })).toBe('JO');
    expect(vm.getUserInitials({})).toBe('??');
  });

  it('initializes pagination state', () => {
    const wrapper = mountWithPlugins(UserSelectorDrawer, { props: defaultProps });
    expect(wrapper.vm.currentPage).toBe(1);
    expect(wrapper.vm.totalPages).toBe(1);
    expect(wrapper.vm.perPage).toBe(15);
  });
});

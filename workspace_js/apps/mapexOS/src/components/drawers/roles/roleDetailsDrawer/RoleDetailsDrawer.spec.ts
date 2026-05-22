import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins, createMockTranslations } from '@src/test/helpers';
import RoleDetailsDrawer from './RoleDetailsDrawer.vue';

vi.mock('@composables/i18n', () => ({
  useRolesTranslations: () => createMockTranslations(),
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
    mapexOS: {
      roles: {
        getById: vi.fn().mockResolvedValue({
          id: 'role-1',
          name: 'Admin',
          isSystem: false,
          scope: 'global',
          permissions: ['users.list', 'users.read'],
          created: '2024-01-01',
          updated: '2024-06-01',
        }),
      },
    },
  },
}));

vi.mock('@utils/alert', () => ({
  notifyFail: vi.fn(),
}));

describe('RoleDetailsDrawer', () => {
  const defaultProps = {
    modelValue: true,
    roleId: 'role-1',
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(RoleDetailsDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts loading when opened with a roleId', () => {
    const wrapper = mountWithPlugins(RoleDetailsDrawer, { props: defaultProps });
    expect(wrapper.vm.loading).toBe(true);
  });

  it('initializes error as false', () => {
    const wrapper = mountWithPlugins(RoleDetailsDrawer, { props: defaultProps });
    expect(wrapper.vm.error).toBe(false);
  });

  it('computes isSystemRole as false when role is not system', async () => {
    const wrapper = mountWithPlugins(RoleDetailsDrawer, { props: defaultProps });
    await wrapper.vm.$nextTick();
    // Before fetch completes, role is null
    expect(wrapper.vm.isSystemRole).toBe(false);
  });

  it('computes groupedPermissions as empty when role is null', () => {
    const wrapper = mountWithPlugins(RoleDetailsDrawer, {
      props: { ...defaultProps, roleId: undefined },
    });
    expect(wrapper.vm.groupedPermissions).toEqual([]);
  });

  it('emits update:modelValue(false) on close', () => {
    const wrapper = mountWithPlugins(RoleDetailsDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.close();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('handles ESC key when drawer is open', () => {
    const wrapper = mountWithPlugins(RoleDetailsDrawer, { props: defaultProps });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });

  it('does not handle ESC key when drawer is closed', () => {
    const wrapper = mountWithPlugins(RoleDetailsDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('update:modelValue')).toBeFalsy();
  });

  it('formats resource names correctly', () => {
    const wrapper = mountWithPlugins(RoleDetailsDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.formatResourceName('auth')).toBe('Authentication');
    expect(vm.formatResourceName('assettemplates')).toBe('Asset Templates');
    expect(vm.formatResourceName('users')).toBe('Users');
  });

  it('returns correct resource icons', () => {
    const wrapper = mountWithPlugins(RoleDetailsDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.getResourceIcon('users')).toBe('person');
    expect(vm.getResourceIcon('roles')).toBe('admin_panel_settings');
    expect(vm.getResourceIcon('unknown')).toBe('vpn_key');
  });

  it('formats date correctly', () => {
    const wrapper = mountWithPlugins(RoleDetailsDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.formatDate(null)).toBe('-');
    expect(vm.formatDate(undefined)).toBe('-');
    expect(vm.formatDate('2024-01-15')).toMatch(/Jan 15, 2024/);
  });
});

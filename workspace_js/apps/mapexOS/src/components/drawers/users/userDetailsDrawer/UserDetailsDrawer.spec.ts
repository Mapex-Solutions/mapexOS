import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins, createMockTranslations } from '@src/test/helpers';
import UserDetailsDrawer from './UserDetailsDrawer.vue';

vi.mock('@composables/i18n', () => ({
  useUsersTranslations: () => createMockTranslations({
    drawer: new Proxy({}, {
      get: (_t, prop) => {
        if (prop === 'authProviders') {
          return new Proxy({}, {
            get: () => ({ value: 'INTERNAL' }),
          });
        }
        return new Proxy({}, {
          get: () => ({ value: String(prop) }),
        });
      },
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
    mapexOS: {
      users: {
        getById: vi.fn().mockResolvedValue({
          id: 'user-1',
          firstName: 'John',
          lastName: 'Doe',
          email: 'john@example.com',
          enabled: true,
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

describe('UserDetailsDrawer', () => {
  const defaultProps = {
    modelValue: true,
    userId: 'user-1',
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(UserDetailsDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts loading when opened with a userId', () => {
    const wrapper = mountWithPlugins(UserDetailsDrawer, { props: defaultProps });
    expect(wrapper.vm.loading).toBe(true);
  });

  it('initializes error as false', () => {
    const wrapper = mountWithPlugins(UserDetailsDrawer, { props: defaultProps });
    expect(wrapper.vm.error).toBe(false);
  });

  it('emits update:modelValue(false) on close', () => {
    const wrapper = mountWithPlugins(UserDetailsDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.close();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('handles ESC key when drawer is open', () => {
    const wrapper = mountWithPlugins(UserDetailsDrawer, { props: defaultProps });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });

  it('does not handle ESC key when drawer is closed', () => {
    const wrapper = mountWithPlugins(UserDetailsDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('update:modelValue')).toBeFalsy();
  });

  it('returns empty string from getUserFullName when user is null', () => {
    const wrapper = mountWithPlugins(UserDetailsDrawer, {
      props: { ...defaultProps, userId: undefined },
    });
    const vm = wrapper.vm;
    expect(vm.getUserFullName()).toBe('');
  });

  it('formats date correctly', () => {
    const wrapper = mountWithPlugins(UserDetailsDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.formatDate(null)).toBe('-');
    expect(vm.formatDate(undefined)).toBe('-');
    expect(vm.formatDate('2024-01-15')).toMatch(/Jan 15, 2024/);
  });

  it('returns correct auth provider icons', () => {
    const wrapper = mountWithPlugins(UserDetailsDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.getAuthProviderIcon()).toBe('lock');
  });

  it('returns organized access as empty when user has no memberships', () => {
    const wrapper = mountWithPlugins(UserDetailsDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.getOrganizedAccess()).toEqual([]);
  });

  it('emits edit and closes on handleEdit when user is set', () => {
    const wrapper = mountWithPlugins(UserDetailsDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    // Simulate user being loaded
    vm.user = { id: 'user-1', firstName: 'John' };
    vm.handleEdit();
    expect(wrapper.emitted('edit')).toBeTruthy();
    expect(wrapper.emitted('edit')![0]).toEqual(['user-1']);
  });
});

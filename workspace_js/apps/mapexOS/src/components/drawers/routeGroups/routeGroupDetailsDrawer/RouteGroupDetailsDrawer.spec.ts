import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins, createMockTranslations } from '@src/test/helpers';
import RouteGroupDetailsDrawer from './RouteGroupDetailsDrawer.vue';

vi.mock('@composables/i18n', () => ({
  useRouteGroupsTranslations: () => createMockTranslations({
    routerKinds: new Proxy({}, {
      get: () => ({ label: { value: 'mock-label' } }),
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
    router: {
      routegroup: {
        getById: vi.fn().mockResolvedValue({
          id: 'rg-1',
          name: 'Test Route Group',
          enabled: true,
          routers: [],
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

describe('RouteGroupDetailsDrawer', () => {
  const defaultProps = {
    modelValue: true,
    routeGroupId: 'rg-1',
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(RouteGroupDetailsDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts loading when opened with a routeGroupId', () => {
    const wrapper = mountWithPlugins(RouteGroupDetailsDrawer, { props: defaultProps });
    expect(wrapper.vm.loading).toBe(true);
  });

  it('initializes error as false', () => {
    const wrapper = mountWithPlugins(RouteGroupDetailsDrawer, { props: defaultProps });
    expect(wrapper.vm.error).toBe(false);
  });

  it('emits update:modelValue(false) on close', () => {
    const wrapper = mountWithPlugins(RouteGroupDetailsDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.close();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('handles ESC key when drawer is open', () => {
    const wrapper = mountWithPlugins(RouteGroupDetailsDrawer, { props: defaultProps });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });

  it('does not handle ESC key when drawer is closed', () => {
    const wrapper = mountWithPlugins(RouteGroupDetailsDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('update:modelValue')).toBeFalsy();
  });

  it('returns correct router icons', () => {
    const wrapper = mountWithPlugins(RouteGroupDetailsDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.getRouterIcon('lake_house')).toBe('storage');
    expect(vm.getRouterIcon('notification')).toBe('notifications');
    expect(vm.getRouterIcon('save_event')).toBe('save');
    expect(vm.getRouterIcon('unknown')).toBe('route');
  });

  it('returns correct router colors', () => {
    const wrapper = mountWithPlugins(RouteGroupDetailsDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.getRouterColor('lake_house')).toBe('purple-6');
    expect(vm.getRouterColor('unknown')).toBe('grey-6');
  });

  it('formats date correctly', () => {
    const wrapper = mountWithPlugins(RouteGroupDetailsDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.formatDate(undefined)).toBe('-');
    expect(vm.formatDate('2024-01-15')).toMatch(/Jan 15, 2024/);
  });

  it('returns router destination info', () => {
    const wrapper = mountWithPlugins(RouteGroupDetailsDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.getRouterDestination({ kind: 'save_event' })).toBe('Saves event to storage');
  });
});

import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import TriggerDetailDrawer from './TriggerDetailDrawer.vue';

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
    triggers: {
      trigger: {
        getById: vi.fn().mockResolvedValue({
          id: 'trigger-1',
          name: 'Test Trigger',
          isSystem: false,
          category: 'technical',
          triggerType: 'http',
          enabled: true,
          config: {},
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

describe('TriggerDetailDrawer', () => {
  const defaultProps = {
    modelValue: true,
    triggerId: 'trigger-1',
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(TriggerDetailDrawer, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts loading when opened with a triggerId', () => {
    const wrapper = mountWithPlugins(TriggerDetailDrawer, { props: defaultProps });
    expect(wrapper.vm.loading).toBe(true);
  });

  it('initializes error as false', () => {
    const wrapper = mountWithPlugins(TriggerDetailDrawer, { props: defaultProps });
    expect(wrapper.vm.error).toBe(false);
  });

  it('computes isSystemTrigger as false when trigger is null', () => {
    const wrapper = mountWithPlugins(TriggerDetailDrawer, { props: defaultProps });
    expect(wrapper.vm.isSystemTrigger).toBe(false);
  });

  it('computes configEntries as empty when trigger is null', () => {
    const wrapper = mountWithPlugins(TriggerDetailDrawer, { props: defaultProps });
    expect(wrapper.vm.configEntries).toEqual([]);
  });

  it('emits update:modelValue(false) on close', () => {
    const wrapper = mountWithPlugins(TriggerDetailDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    vm.close();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('handles ESC key when drawer is open', () => {
    const wrapper = mountWithPlugins(TriggerDetailDrawer, { props: defaultProps });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });

  it('does not handle ESC key when drawer is closed', () => {
    const wrapper = mountWithPlugins(TriggerDetailDrawer, {
      props: { ...defaultProps, modelValue: false },
    });
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(wrapper.emitted('update:modelValue')).toBeFalsy();
  });

  it('returns correct category icons', () => {
    const wrapper = mountWithPlugins(TriggerDetailDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.getCategoryIcon('technical')).toBe('dns');
    expect(vm.getCategoryIcon('communication')).toBe('chat');
    expect(vm.getCategoryIcon(undefined)).toBe('category');
  });

  it('returns correct category colors', () => {
    const wrapper = mountWithPlugins(TriggerDetailDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.getCategoryColor('technical')).toBe('purple');
    expect(vm.getCategoryColor('communication')).toBe('teal');
    expect(vm.getCategoryColor(undefined)).toBe('grey');
  });

  it('returns correct type icons', () => {
    const wrapper = mountWithPlugins(TriggerDetailDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.getTypeIcon('http')).toBe('http');
    expect(vm.getTypeIcon('mqtt')).toBe('wifi_tethering');
    expect(vm.getTypeIcon('email')).toBe('email');
  });

  it('formats config keys correctly', () => {
    const wrapper = mountWithPlugins(TriggerDetailDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.formatConfigKey('maxRetries')).toBe('Max Retries');
    expect(vm.formatConfigKey('timeout')).toBe('Timeout');
  });

  it('formats config values correctly', () => {
    const wrapper = mountWithPlugins(TriggerDetailDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.formatConfigValue(null)).toBe('-');
    expect(vm.formatConfigValue('test')).toBe('test');
    expect(vm.formatConfigValue(42)).toBe('42');
    expect(vm.formatConfigValue(true)).toBe('Yes');
    expect(vm.formatConfigValue(false)).toBe('No');
  });

  it('formats date correctly', () => {
    const wrapper = mountWithPlugins(TriggerDetailDrawer, { props: defaultProps });
    const vm = wrapper.vm;
    expect(vm.formatDate(null)).toBe('-');
    expect(vm.formatDate('2024-01-15')).toMatch(/Jan 15, 2024/);
  });
});

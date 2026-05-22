import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import PushConfig from './PushConfig.vue';
import type { ChannelPushProps } from './interfaces';

vi.mock('./constants', () => ({
  SERVICE_PROVIDERS: [
    { value: 'firebase', label: 'Firebase Cloud Messaging', description: "Google's service" },
    { value: 'apns', label: 'Apple Push Notification Service', description: "Apple's service" },
  ],
  PRIORITY_OPTIONS: [
    { value: 'high', label: 'High', description: 'Immediate delivery' },
    { value: 'normal', label: 'Normal', description: 'Standard delivery' },
    { value: 'low', label: 'Low', description: 'Battery-saving delivery' },
  ],
  SOUND_OPTIONS: ['default', 'none', 'custom1'],
}));

const BASE_MODEL: ChannelPushProps = {
  appName: 'MyApp',
  deviceCount: 42,
  apiKey: 'api-key-123',
  serviceProvider: 'firebase',
  priority: 'high',
  ttl: 3600,
  badge: true,
  sound: 'default',
  clickAction: 'OPEN_APP',
};

describe('PushConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(PushConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('initializes localData from modelValue', () => {
    const wrapper = mountWithPlugins(PushConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.localData.appName).toBe('MyApp');
    expect(wrapper.vm.localData.deviceCount).toBe(42);
    expect(wrapper.vm.localData.serviceProvider).toBe('firebase');
  });

  it('returns provider label from getProviderLabel', () => {
    const wrapper = mountWithPlugins(PushConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.getProviderLabel('firebase')).toBe('Firebase Cloud Messaging');
  });

  it('returns raw value from getProviderLabel for unknown provider', () => {
    const wrapper = mountWithPlugins(PushConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.getProviderLabel('unknown')).toBe('unknown');
  });

  it('returns priority label from getPriorityLabel', () => {
    const wrapper = mountWithPlugins(PushConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.getPriorityLabel('high')).toBe('High');
  });

  it('returns raw value from getPriorityLabel for unknown priority', () => {
    const wrapper = mountWithPlugins(PushConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.getPriorityLabel('unknown')).toBe('unknown');
  });

  it('initializes showApiKey as false', () => {
    const wrapper = mountWithPlugins(PushConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.showApiKey).toBe(false);
  });

  it('emits update:modelValue when localData changes', async () => {
    const wrapper = mountWithPlugins(PushConfig, {
      props: { modelValue: BASE_MODEL },
    });
    wrapper.vm.localData.appName = 'NewApp';
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });
});

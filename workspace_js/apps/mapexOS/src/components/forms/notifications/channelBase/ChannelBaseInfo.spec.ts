import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import ChannelBaseInfo from './ChannelBaseInfo.vue';

vi.mock('./constants', () => ({
  CHANNEL_TYPES: [
    { value: 'email', label: 'Email', icon: 'email', description: 'Email channel' },
    { value: 'slack', label: 'Slack', icon: 'chat', description: 'Slack channel' },
  ],
}));

const BASE_MODEL = {
  type: 'email' as const,
  name: 'My Channel',
  description: 'A channel',
  icon: 'email',
  status: 'Active' as const,
  created: '2024-01-01',
};

describe('ChannelBaseInfo', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(ChannelBaseInfo, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes isActive as true when status is Active', () => {
    const wrapper = mountWithPlugins(ChannelBaseInfo, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isActive).toBe(true);
  });

  it('computes isActive as false when status is Inactive', () => {
    const wrapper = mountWithPlugins(ChannelBaseInfo, {
      props: { modelValue: { ...BASE_MODEL, status: 'Inactive' } },
    });
    expect(wrapper.vm.isActive).toBe(false);
  });

  it('updates status to Inactive when isActive set to false', () => {
    const wrapper = mountWithPlugins(ChannelBaseInfo, {
      props: { modelValue: BASE_MODEL },
    });
    wrapper.vm.isActive = false;
    expect(wrapper.vm.localData.status).toBe('Inactive');
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });

  it('updates icon on handleTypeChange', () => {
    const wrapper = mountWithPlugins(ChannelBaseInfo, {
      props: { modelValue: BASE_MODEL },
    });
    wrapper.vm.handleTypeChange('slack');
    expect(wrapper.vm.localData.icon).toBe('chat');
  });

  it('emits update:modelValue on handleTypeChange', () => {
    const wrapper = mountWithPlugins(ChannelBaseInfo, {
      props: { modelValue: BASE_MODEL },
    });
    wrapper.vm.handleTypeChange('slack');
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });
});

import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import PushChannel from './PushChannel.vue';

function factory(channel: { type: 'push'; appName: string; deviceCount?: number }) {
  return mountWithPlugins(PushChannel, {
    props: { channel },
  });
}

describe('PushChannel', () => {
  it('renders without errors', () => {
    const wrapper = factory({ type: 'push', appName: 'MyApp' });
    expect(wrapper.exists()).toBe(true);
  });

  it('receives appName prop', () => {
    const wrapper = factory({ type: 'push', appName: 'MyApp', deviceCount: 50 });
    expect(wrapper.props('channel').appName).toBe('MyApp');
  });

  it('receives deviceCount prop', () => {
    const wrapper = factory({ type: 'push', appName: 'MyApp', deviceCount: 100 });
    expect(wrapper.props('channel').deviceCount).toBe(100);
  });

  it('renders without deviceCount (optional)', () => {
    const wrapper = factory({ type: 'push', appName: 'MyApp' });
    expect(wrapper.exists()).toBe(true);
    expect(wrapper.props('channel').deviceCount).toBeUndefined();
  });
});

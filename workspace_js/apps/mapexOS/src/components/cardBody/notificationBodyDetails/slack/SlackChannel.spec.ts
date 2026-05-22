import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import SlackChannel from './SlackChannel.vue';

function factory(channel: { type: 'slack'; workspace: string; channelName: string }) {
  return mountWithPlugins(SlackChannel, {
    props: { channel },
  });
}

describe('SlackChannel', () => {
  it('renders without errors', () => {
    const wrapper = factory({ type: 'slack', workspace: 'Acme Corp', channelName: '#alerts' });
    expect(wrapper.exists()).toBe(true);
  });

  it('receives workspace prop', () => {
    const wrapper = factory({ type: 'slack', workspace: 'Acme Corp', channelName: '#alerts' });
    expect(wrapper.props('channel').workspace).toBe('Acme Corp');
  });

  it('receives channelName prop', () => {
    const wrapper = factory({ type: 'slack', workspace: 'Acme Corp', channelName: '#general' });
    expect(wrapper.props('channel').channelName).toBe('#general');
  });
});

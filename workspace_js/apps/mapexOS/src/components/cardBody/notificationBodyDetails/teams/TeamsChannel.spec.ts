import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import TeamsChannel from './TeamsChannel.vue';

function factory(channel: { type: 'teams'; teamName: string; channelName: string }) {
  return mountWithPlugins(TeamsChannel, {
    props: { channel },
  });
}

describe('TeamsChannel', () => {
  it('renders without errors', () => {
    const wrapper = factory({ type: 'teams', teamName: 'Engineering', channelName: '#ops' });
    expect(wrapper.exists()).toBe(true);
  });

  it('receives teamName prop', () => {
    const wrapper = factory({ type: 'teams', teamName: 'Engineering', channelName: '#ops' });
    expect(wrapper.props('channel').teamName).toBe('Engineering');
  });

  it('receives channelName prop', () => {
    const wrapper = factory({ type: 'teams', teamName: 'Engineering', channelName: '#incidents' });
    expect(wrapper.props('channel').channelName).toBe('#incidents');
  });
});

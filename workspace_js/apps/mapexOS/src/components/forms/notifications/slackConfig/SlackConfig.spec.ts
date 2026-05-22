import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import SlackConfig from './SlackConfig.vue';
import type { ChannelSlackProps } from './interfaces';

const BASE_MODEL: ChannelSlackProps = {
  workspace: 'my-workspace',
  channelsName: ['#general'],
  webhookUrl: 'https://hooks.slack.com/services/T00/B00/xxxx',
  messageTemplate: '*{{title}}*\n{{message}}',
  botName: 'MapexBot',
};

describe('SlackConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(SlackConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('initializes localData from modelValue', () => {
    const wrapper = mountWithPlugins(SlackConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.localData.workspace).toBe('my-workspace');
    expect(wrapper.vm.localData.botName).toBe('MapexBot');
  });

  it('computes previewMessage replacing placeholders', () => {
    const wrapper = mountWithPlugins(SlackConfig, {
      props: { modelValue: BASE_MODEL },
    });
    const preview = wrapper.vm.previewMessage;
    expect(preview).toContain('Example Title');
    expect(preview).toContain('example message');
    expect(preview).not.toContain('{{title}}');
    expect(preview).not.toContain('{{message}}');
  });

  it('validates correct URL with isValidUrl', () => {
    const wrapper = mountWithPlugins(SlackConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidUrl('https://hooks.slack.com/test')).toBe(true);
  });

  it('rejects invalid URL with isValidUrl', () => {
    const wrapper = mountWithPlugins(SlackConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidUrl('not-valid')).toBe(false);
  });

  it('emits update:modelValue when localData changes', async () => {
    const wrapper = mountWithPlugins(SlackConfig, {
      props: { modelValue: BASE_MODEL },
    });
    wrapper.vm.localData.workspace = 'new-workspace';
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });
});

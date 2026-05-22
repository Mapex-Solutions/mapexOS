import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import TeamsConfig from './TeamsConfig.vue';
import type { ChannelTeamsProps } from './interfaces';

const BASE_MODEL: ChannelTeamsProps = {
  teamName: 'Engineering',
  channelsName: ['general'],
  webhookUrl: 'https://outlook.office.com/webhook/test',
  messageTemplate: '{"title":"{{title}}","message":"{{message}}"}',
  adaptiveCard: false,
};

describe('TeamsConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(TeamsConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('initializes localData from modelValue', () => {
    const wrapper = mountWithPlugins(TeamsConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.localData.teamName).toBe('Engineering');
    expect(wrapper.vm.localData.adaptiveCard).toBe(false);
  });

  it('computes previewMessage replacing placeholders', () => {
    const wrapper = mountWithPlugins(TeamsConfig, {
      props: { modelValue: BASE_MODEL },
    });
    const preview = wrapper.vm.previewMessage;
    expect(preview).toContain('Example Title');
    expect(preview).not.toContain('{{title}}');
  });

  it('computes formattedPreview as formatted JSON', () => {
    const wrapper = mountWithPlugins(TeamsConfig, {
      props: { modelValue: BASE_MODEL },
    });
    const formatted = wrapper.vm.formattedPreview;
    // Should be pretty-printed JSON
    expect(formatted).toContain('\n');
    expect(formatted).toContain('Example Title');
  });

  it('returns plain text from formattedPreview when template is not valid JSON', () => {
    const wrapper = mountWithPlugins(TeamsConfig, {
      props: { modelValue: { ...BASE_MODEL, messageTemplate: 'not json {{title}}' } },
    });
    const formatted = wrapper.vm.formattedPreview;
    expect(formatted).toContain('Example Title');
  });

  it('validates correct URL with isValidUrl', () => {
    const wrapper = mountWithPlugins(TeamsConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidUrl('https://example.com')).toBe(true);
  });

  it('rejects invalid URL with isValidUrl', () => {
    const wrapper = mountWithPlugins(TeamsConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidUrl('not-a-url')).toBe(false);
  });

  it('validates correct JSON with isValidJson', () => {
    const wrapper = mountWithPlugins(TeamsConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidJson('{"key":"value"}')).toBe(true);
  });

  it('rejects invalid JSON with isValidJson', () => {
    const wrapper = mountWithPlugins(TeamsConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidJson('{bad}')).toBe(false);
  });

  it('emits update:modelValue when localData changes', async () => {
    const wrapper = mountWithPlugins(TeamsConfig, {
      props: { modelValue: BASE_MODEL },
    });
    wrapper.vm.localData.teamName = 'New Team';
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });
});

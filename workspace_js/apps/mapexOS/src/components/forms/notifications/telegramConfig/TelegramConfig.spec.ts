import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import TelegramConfig from './TelegramConfig.vue';
import type { ChannelTelegramProps } from './interfaces';

vi.mock('./constants', () => ({
  PARSE_MODES: [
    { value: 'Markdown', label: 'Markdown', description: 'Basic formatting' },
    { value: 'HTML', label: 'HTML', description: 'HTML tags' },
    { value: 'MarkdownV2', label: 'Markdown V2', description: 'Latest version' },
    { value: 'None', label: 'None', description: 'No formatting' },
  ],
}));

const BASE_MODEL: ChannelTelegramProps = {
  botName: '@TestBot',
  chatNames: ['chat-1'],
  botToken: 'bot-token-123',
  parseMode: 'Markdown',
  disableNotification: false,
  messageTemplate: '**{{title}}**\n{{message}}',
};

describe('TelegramConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(TelegramConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('initializes localData from modelValue', () => {
    const wrapper = mountWithPlugins(TelegramConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.localData.botName).toBe('@TestBot');
    expect(wrapper.vm.localData.parseMode).toBe('Markdown');
  });

  it('computes previewMessage replacing placeholders', () => {
    const wrapper = mountWithPlugins(TelegramConfig, {
      props: { modelValue: BASE_MODEL },
    });
    const preview = wrapper.vm.previewMessage;
    expect(preview).toContain('Example Title');
    expect(preview).toContain('example message');
    expect(preview).not.toContain('{{title}}');
  });

  it('computes getParseClass for Markdown parse mode', () => {
    const wrapper = mountWithPlugins(TelegramConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.getParseClass).toBe('telegram-markdown');
  });

  it('computes getParseClass for MarkdownV2 parse mode', () => {
    const wrapper = mountWithPlugins(TelegramConfig, {
      props: { modelValue: { ...BASE_MODEL, parseMode: 'MarkdownV2' } },
    });
    expect(wrapper.vm.getParseClass).toBe('telegram-markdown');
  });

  it('computes getParseClass for HTML parse mode', () => {
    const wrapper = mountWithPlugins(TelegramConfig, {
      props: { modelValue: { ...BASE_MODEL, parseMode: 'HTML' } },
    });
    expect(wrapper.vm.getParseClass).toBe('telegram-html');
  });

  it('computes getParseClass as empty string for None', () => {
    const wrapper = mountWithPlugins(TelegramConfig, {
      props: { modelValue: { ...BASE_MODEL, parseMode: 'None' } },
    });
    expect(wrapper.vm.getParseClass).toBe('');
  });

  it('initializes showToken as false', () => {
    const wrapper = mountWithPlugins(TelegramConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.showToken).toBe(false);
  });

  it('emits update:modelValue when localData changes', async () => {
    const wrapper = mountWithPlugins(TelegramConfig, {
      props: { modelValue: BASE_MODEL },
    });
    wrapper.vm.localData.botName = '@NewBot';
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });
});

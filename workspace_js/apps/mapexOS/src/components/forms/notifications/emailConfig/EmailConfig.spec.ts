import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import EmailConfig from './EmailConfig.vue';
import type { ChannelEmailProps } from './interfaces';

const BASE_MODEL: ChannelEmailProps = {
  from: 'sender@example.com',
  subject: 'Test Subject',
  to: ['recipient@example.com'],
  cc: [],
  bcc: [],
  template: '<h1>{{title}}</h1><p>{{message}}</p>',
  smtp: {
    host: 'smtp.example.com',
    port: 587,
    secure: true,
    auth: {
      user: 'admin',
      password: 'secret',
    },
  },
  attachments: false,
};

describe('EmailConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(EmailConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('initializes localData from modelValue', () => {
    const wrapper = mountWithPlugins(EmailConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.localData.from).toBe('sender@example.com');
    expect(wrapper.vm.localData.subject).toBe('Test Subject');
  });

  it('computes previewEmail replacing placeholders', () => {
    const wrapper = mountWithPlugins(EmailConfig, {
      props: { modelValue: BASE_MODEL },
    });
    const preview = wrapper.vm.previewEmail;
    expect(preview).toContain('Example Title');
    expect(preview).toContain('This is an example message');
    expect(preview).not.toContain('{{title}}');
    expect(preview).not.toContain('{{message}}');
  });

  it('validates correct email with isValidEmail', () => {
    const wrapper = mountWithPlugins(EmailConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidEmail('user@example.com')).toBe(true);
  });

  it('rejects invalid email with isValidEmail', () => {
    const wrapper = mountWithPlugins(EmailConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidEmail('not-an-email')).toBe(false);
    expect(wrapper.vm.isValidEmail('')).toBe(false);
  });

  it('emits update:modelValue when localData changes', async () => {
    const wrapper = mountWithPlugins(EmailConfig, {
      props: { modelValue: BASE_MODEL },
    });
    wrapper.vm.localData.from = 'new@example.com';
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });
});

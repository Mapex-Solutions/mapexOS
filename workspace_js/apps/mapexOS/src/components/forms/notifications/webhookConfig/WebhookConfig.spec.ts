import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import WebhookConfig from './WebhookConfig.vue';
import type { ChannelWebhookProps } from './interfaces';

const BASE_MODEL: ChannelWebhookProps = {
  name: 'My Webhook',
  method: 'POST',
  url: 'https://example.com/webhook',
  headers: { 'Content-Type': 'application/json' },
  payload: '{"title":"{{title}}","message":"{{message}}","timestamp":"{{timestamp}}"}',
  timeout: 5000,
  retryCount: 3,
};

describe('WebhookConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(WebhookConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('initializes localData from modelValue', () => {
    const wrapper = mountWithPlugins(WebhookConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.localData.name).toBe('My Webhook');
    expect(wrapper.vm.localData.method).toBe('POST');
    expect(wrapper.vm.localData.url).toBe('https://example.com/webhook');
  });

  it('validates correct URL with isValidUrl', () => {
    const wrapper = mountWithPlugins(WebhookConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidUrl('https://example.com')).toBe(true);
  });

  it('rejects invalid URL with isValidUrl', () => {
    const wrapper = mountWithPlugins(WebhookConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidUrl('not-a-url')).toBe(false);
  });

  it('validates correct JSON with isValidJson', () => {
    const wrapper = mountWithPlugins(WebhookConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidJson('{"key":"value"}')).toBe(true);
  });

  it('rejects invalid JSON with isValidJson', () => {
    const wrapper = mountWithPlugins(WebhookConfig, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidJson('{invalid}')).toBe(false);
  });

  it('adds a header row with addHeader', () => {
    const wrapper = mountWithPlugins(WebhookConfig, {
      props: { modelValue: BASE_MODEL },
    });
    const initialCount = wrapper.vm.headerIndices.length;
    wrapper.vm.addHeader();
    expect(wrapper.vm.headerIndices.length).toBe(initialCount + 1);
  });

  it('removes a header row with removeHeader', () => {
    const wrapper = mountWithPlugins(WebhookConfig, {
      props: { modelValue: BASE_MODEL },
    });
    // After mount, initializeHeaders runs via onMounted
    // Add then remove
    wrapper.vm.addHeader();
    const countAfterAdd = wrapper.vm.headerIndices.length;
    const lastIndex = wrapper.vm.headerIndices[wrapper.vm.headerIndices.length - 1];
    wrapper.vm.removeHeader(lastIndex);
    expect(wrapper.vm.headerIndices.length).toBe(countAfterAdd - 1);
  });

  it('rebuilds headers object on updateHeaderKey', () => {
    const wrapper = mountWithPlugins(WebhookConfig, {
      props: { modelValue: BASE_MODEL },
    });
    // Set key and value manually then call updateHeaderKey
    const idx = wrapper.vm.headerIndices[0];
    wrapper.vm.headerKeys[idx] = 'Authorization';
    wrapper.vm.headerValues[idx] = 'Bearer token';
    wrapper.vm.updateHeaderKey();
    expect(wrapper.vm.localData.headers['Authorization']).toBe('Bearer token');
  });

  it('updates a specific header value with updateHeaderValue', () => {
    const wrapper = mountWithPlugins(WebhookConfig, {
      props: { modelValue: BASE_MODEL },
    });
    const idx = wrapper.vm.headerIndices[0];
    const key = wrapper.vm.headerKeys[idx];
    wrapper.vm.headerValues[idx] = 'text/plain';
    wrapper.vm.updateHeaderValue(idx);
    expect(wrapper.vm.localData.headers[key]).toBe('text/plain');
  });

  it('sets testing to true during testWebhook', async () => {
    const wrapper = mountWithPlugins(WebhookConfig, {
      props: { modelValue: BASE_MODEL },
    });
    const promise = wrapper.vm.testWebhook();
    expect(wrapper.vm.testing).toBe(true);
    await promise;
    expect(wrapper.vm.testing).toBe(false);
  });

  it('emits update:modelValue when localData changes', async () => {
    const wrapper = mountWithPlugins(WebhookConfig, {
      props: { modelValue: BASE_MODEL },
    });
    wrapper.vm.localData.name = 'Updated Webhook';
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
  });
});

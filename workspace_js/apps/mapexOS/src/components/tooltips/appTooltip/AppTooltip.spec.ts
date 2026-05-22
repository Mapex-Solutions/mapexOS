import { describe, it, expect, vi } from 'vitest';
import { shallowMount } from '@vue/test-utils';
import AppTooltip from './AppTooltip.vue';

vi.mock('quasar', () => ({
  useQuasar: () => ({
    platform: {
      is: { mobile: false },
      has: { touch: false },
    },
  }),
}));

describe('AppTooltip', () => {
  it('renders without errors', () => {
    const wrapper = shallowMount(AppTooltip, {
      props: { content: 'Hello' },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('shouldShow is true on desktop by default', () => {
    const wrapper = shallowMount(AppTooltip, {
      props: { content: 'Tip' },
    });
    expect((wrapper.vm as any).shouldShow).toBe(true);
  });

  it('shouldShow is false when disabled', () => {
    const wrapper = shallowMount(AppTooltip, {
      props: { content: 'Tip', disabled: true },
    });
    expect((wrapper.vm as any).shouldShow).toBe(false);
  });
});

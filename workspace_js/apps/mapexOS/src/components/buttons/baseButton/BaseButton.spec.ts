import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import BaseButton from './BaseButton.vue';

describe('BaseButton', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(BaseButton);
    expect(wrapper.exists()).toBe(true);
  });

  it('always includes rounded-borders class', () => {
    const wrapper = mountWithPlugins(BaseButton);
    expect(wrapper.vm.buttonClass).toContain('rounded-borders');
  });

  it('defaults to empty user class when none provided', () => {
    const wrapper = mountWithPlugins(BaseButton);
    expect(wrapper.vm.buttonClass).toBe('rounded-borders');
  });

  it('qBtnProps computed exists and is an object', () => {
    const wrapper = mountWithPlugins(BaseButton);
    expect(wrapper.vm.qBtnProps).toBeDefined();
    expect(typeof wrapper.vm.qBtnProps).toBe('object');
  });

  it('renders q-btn stub in shallow mode', () => {
    const wrapper = mountWithPlugins(BaseButton);
    expect(wrapper.find('q-btn-stub').exists()).toBe(true);
  });
});

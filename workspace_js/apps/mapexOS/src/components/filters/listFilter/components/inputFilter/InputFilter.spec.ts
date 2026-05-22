import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import InputFilter from './InputFilter.vue';

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(InputFilter, {
    props: {
      modelValue: '',
      label: 'Search',
      icon: 'search',
      ...overrides,
    },
  });
}

describe('InputFilter', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('value computed getter returns modelValue', () => {
    const wrapper = factory({ modelValue: 'hello' });
    expect(wrapper.vm.value).toBe('hello');
  });

  it('value computed setter emits update:modelValue', () => {
    const wrapper = factory();
    wrapper.vm.value = 'new value';
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual(['new value']);
  });

  it('handleEnter emits enter event', () => {
    const wrapper = factory();
    wrapper.vm.handleEnter();
    expect(wrapper.emitted('enter')).toBeTruthy();
    expect(wrapper.emitted('enter')).toHaveLength(1);
  });

  it('clearable defaults to true', () => {
    const wrapper = factory();
    expect(wrapper.props('clearable')).toBe(true);
  });

  it('disabled defaults to false', () => {
    const wrapper = factory();
    expect(wrapper.props('disabled')).toBe(false);
  });

  it('debounce defaults to 0', () => {
    const wrapper = factory();
    expect(wrapper.props('debounce')).toBe(0);
  });

  it('accepts custom debounce value', () => {
    const wrapper = factory({ debounce: 500 });
    expect(wrapper.props('debounce')).toBe(500);
  });

  it('accepts mask prop', () => {
    const wrapper = factory({ mask: '####' });
    expect(wrapper.props('mask')).toBe('####');
  });

  it('accepts type prop', () => {
    const wrapper = factory({ type: 'number' });
    expect(wrapper.props('type')).toBe('number');
  });

  it('accepts placeholder prop', () => {
    const wrapper = factory({ placeholder: 'Type here...' });
    expect(wrapper.props('placeholder')).toBe('Type here...');
  });
});

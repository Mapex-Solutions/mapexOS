import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import FieldInputSelector from './FieldInputSelector.vue';

describe('FieldInputSelector', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(FieldInputSelector);
    expect(wrapper.exists()).toBe(true);
  });

  it('selectedType defaults to literal when no modelValue', () => {
    const wrapper = mountWithPlugins(FieldInputSelector);
    expect(wrapper.vm.selectedType).toBe('literal');
  });

  it('selectedType reflects modelValue prop', () => {
    const wrapper = mountWithPlugins(FieldInputSelector, {
      props: { modelValue: 'event' },
    });
    expect(wrapper.vm.selectedType).toBe('event');
  });

  it('currentOption returns matching option for event type', () => {
    const wrapper = mountWithPlugins(FieldInputSelector, {
      props: { modelValue: 'event' },
    });
    expect(wrapper.vm.currentOption.value).toBe('event');
    expect(wrapper.vm.currentOption.icon).toBe('event');
    expect(wrapper.vm.currentOption.color).toBe('blue-6');
  });

  it('currentOption returns matching option for state type', () => {
    const wrapper = mountWithPlugins(FieldInputSelector, {
      props: { modelValue: 'state' },
    });
    expect(wrapper.vm.currentOption.value).toBe('state');
    expect(wrapper.vm.currentOption.icon).toBe('storage');
  });

  it('currentOption returns matching option for variable type', () => {
    const wrapper = mountWithPlugins(FieldInputSelector, {
      props: { modelValue: 'variable' },
    });
    expect(wrapper.vm.currentOption.value).toBe('variable');
    expect(wrapper.vm.currentOption.icon).toBe('code');
  });

  it('currentOption returns matching option for literal type', () => {
    const wrapper = mountWithPlugins(FieldInputSelector, {
      props: { modelValue: 'literal' },
    });
    expect(wrapper.vm.currentOption.value).toBe('literal');
    expect(wrapper.vm.currentOption.icon).toBe('format_quote');
  });

  it('currentOption falls back to literal (index 3) for unknown type', () => {
    const wrapper = mountWithPlugins(FieldInputSelector, {
      props: { modelValue: 'nonexistent' as any },
    });
    expect(wrapper.vm.currentOption.value).toBe('literal');
  });

  it('setting selectedType emits update:modelValue', () => {
    const wrapper = mountWithPlugins(FieldInputSelector, {
      props: { modelValue: 'literal' },
    });
    wrapper.vm.selectedType = 'event';
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]![0]).toBe('event');
  });
});

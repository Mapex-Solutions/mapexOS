import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import MultiSelectFilter from './MultiSelectFilter.vue';

const baseOptions = [
  { label: 'MQTT', value: 'mqtt', icon: 'mdi-wifi' },
  { label: 'HTTP', value: 'http', icon: 'mdi-web' },
  { label: 'TCP', value: 'tcp' },
];

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(MultiSelectFilter, {
    props: {
      modelValue: [],
      label: 'Protocol',
      icon: 'mdi-filter',
      options: baseOptions,
      ...overrides,
    },
  });
}

describe('MultiSelectFilter', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('value computed getter returns modelValue', () => {
    const wrapper = factory({ modelValue: ['mqtt'] });
    expect(wrapper.vm.value).toEqual(['mqtt']);
  });

  it('emits update:modelValue when value is set', async () => {
    const wrapper = factory();
    wrapper.vm.value = ['mqtt', 'http'];
    await wrapper.vm.$nextTick();
    const emitted = wrapper.emitted('update:modelValue');
    expect(emitted).toBeTruthy();
    expect(emitted![0]![0]).toEqual(['mqtt', 'http']);
  });

  it('hasIcons is true when at least one option has an icon', () => {
    const wrapper = factory();
    expect(wrapper.vm.hasIcons).toBe(true);
  });

  it('hasIcons is false when no options have icons', () => {
    const wrapper = factory({
      options: [{ label: 'A', value: 'a' }, { label: 'B', value: 'b' }],
    });
    expect(wrapper.vm.hasIcons).toBe(false);
  });

  it('defaults clearable to true', () => {
    const wrapper = factory();
    expect(wrapper.props('clearable')).toBe(true);
  });

  it('defaults disabled to false', () => {
    const wrapper = factory();
    expect(wrapper.props('disabled')).toBe(false);
  });
});

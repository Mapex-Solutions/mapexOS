import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import SelectFilter from './SelectFilter.vue';

const baseOptions = [
  { label: 'Active', value: 'active', icon: 'mdi-check-circle', color: 'green' },
  { label: 'Inactive', value: 'inactive' },
];

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(SelectFilter, {
    props: {
      modelValue: null,
      label: 'Status',
      icon: 'mdi-filter',
      options: baseOptions,
      ...overrides,
    },
  });
}

describe('SelectFilter', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('value computed getter returns modelValue', () => {
    const wrapper = factory({ modelValue: 'active' });
    expect(wrapper.vm.value).toBe('active');
  });

  it('emits update:modelValue when value is set', async () => {
    const wrapper = factory();
    wrapper.vm.value = 'inactive';
    await wrapper.vm.$nextTick();
    const emitted = wrapper.emitted('update:modelValue');
    expect(emitted).toBeTruthy();
    expect(emitted![0]![0]).toBe('inactive');
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

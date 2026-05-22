import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import DateRangeFilter from './DateRangeFilter.vue';

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(DateRangeFilter, {
    props: {
      modelValue: null,
      label: 'Created',
      icon: 'event',
      ...overrides,
    },
  });
}

describe('DateRangeFilter', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('value computed getter returns modelValue', () => {
    const wrapper = factory({ modelValue: { from: '2025-01-01', to: '2025-01-31' } });
    expect(wrapper.vm.value).toEqual({ from: '2025-01-01', to: '2025-01-31' });
  });

  it('fromDate computed getter returns from part of modelValue', () => {
    const wrapper = factory({ modelValue: { from: '2025-01-01 10:00:00' } });
    expect(wrapper.vm.fromDate).toBe('2025-01-01 10:00:00');
  });

  it('fromDate returns null when modelValue is null', () => {
    const wrapper = factory({ modelValue: null });
    expect(wrapper.vm.fromDate).toBeNull();
  });

  it('toDate computed getter returns to part of modelValue', () => {
    const wrapper = factory({ modelValue: { to: '2025-01-31 23:59:59' } });
    expect(wrapper.vm.toDate).toBe('2025-01-31 23:59:59');
  });

  it('toDate returns null when modelValue has no to field', () => {
    const wrapper = factory({ modelValue: { from: '2025-01-01' } });
    expect(wrapper.vm.toDate).toBeNull();
  });

  it('emits update:modelValue when fromDate is set', async () => {
    const wrapper = factory({ modelValue: { to: '2025-01-31' } });
    wrapper.vm.fromDate = '2025-01-01';
    await wrapper.vm.$nextTick();
    const emitted = wrapper.emitted('update:modelValue');
    expect(emitted).toBeTruthy();
    expect(emitted![0]![0]).toEqual({ to: '2025-01-31', from: '2025-01-01' });
  });

  it('emits update:modelValue with from removed when fromDate is set to null', async () => {
    const wrapper = factory({ modelValue: { from: '2025-01-01', to: '2025-01-31' } });
    wrapper.vm.fromDate = null;
    await wrapper.vm.$nextTick();
    const emitted = wrapper.emitted('update:modelValue');
    expect(emitted).toBeTruthy();
    expect(emitted![0]![0]).toEqual({ to: '2025-01-31' });
  });

  it('emits update:modelValue when toDate is set', async () => {
    const wrapper = factory({ modelValue: { from: '2025-01-01' } });
    wrapper.vm.toDate = '2025-01-31';
    await wrapper.vm.$nextTick();
    const emitted = wrapper.emitted('update:modelValue');
    expect(emitted).toBeTruthy();
    expect(emitted![0]![0]).toEqual({ from: '2025-01-01', to: '2025-01-31' });
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

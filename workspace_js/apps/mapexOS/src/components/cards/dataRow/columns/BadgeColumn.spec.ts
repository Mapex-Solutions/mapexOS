import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import BadgeColumn from './BadgeColumn.vue';
import type { DataRowColumn } from '../interfaces';

const baseColumn: DataRowColumn = {
  key: 'status',
  label: 'Status',
  type: 'badge',
  visible: 'always',
  color: 'green',
};

const baseRow = { id: '1', status: 'active' };

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(BadgeColumn, {
    props: {
      value: 'active',
      column: baseColumn,
      row: baseRow,
      ...overrides,
    },
  });
}

describe('BadgeColumn', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('displayValue returns raw value when no format function', () => {
    const wrapper = factory();
    expect(wrapper.vm.displayValue).toBe('active');
  });

  it('displayValue returns "N/A" when value is falsy and no format', () => {
    const wrapper = factory({ value: '' });
    expect(wrapper.vm.displayValue).toBe('N/A');
  });

  it('displayValue applies format function when provided', () => {
    const wrapper = factory({
      column: { ...baseColumn, format: (v: any) => v.toUpperCase() },
    });
    expect(wrapper.vm.displayValue).toBe('ACTIVE');
  });

  it('getColor returns static color from column', () => {
    const wrapper = factory();
    expect(wrapper.vm.getColor()).toBe('green');
  });

  it('getColor calls function when color is a function', () => {
    const wrapper = factory({
      column: { ...baseColumn, color: () => 'red-5' },
    });
    expect(wrapper.vm.getColor()).toBe('red-5');
  });

  it('getColor defaults to "grey" when no color set', () => {
    const wrapper = factory({
      column: { ...baseColumn, color: undefined },
    });
    expect(wrapper.vm.getColor()).toBe('grey');
  });
});

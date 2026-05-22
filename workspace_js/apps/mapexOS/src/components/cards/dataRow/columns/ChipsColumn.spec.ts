import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import ChipsColumn from './ChipsColumn.vue';
import type { DataRowColumn } from '../interfaces';

const baseColumn: DataRowColumn = {
  key: 'tags',
  label: 'Tags',
  type: 'chips',
  visible: 'always',
  color: 'teal',
};

const baseRow = { id: '1', tags: ['alpha', 'beta', 'gamma', 'delta', 'epsilon'] };

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(ChipsColumn, {
    props: {
      value: baseRow.tags,
      column: baseColumn,
      row: baseRow,
      ...overrides,
    },
  });
}

describe('ChipsColumn', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('allItems returns array value as-is', () => {
    const wrapper = factory();
    expect(wrapper.vm.allItems).toEqual(baseRow.tags);
  });

  it('allItems wraps string value into array', () => {
    const wrapper = factory({ value: 'solo' });
    expect(wrapper.vm.allItems).toEqual(['solo']);
  });

  it('allItems returns empty array for non-array/non-string', () => {
    const wrapper = factory({ value: 123 });
    expect(wrapper.vm.allItems).toEqual([]);
  });

  it('visibleItems shows at most 3 items', () => {
    const wrapper = factory();
    expect(wrapper.vm.visibleItems).toHaveLength(3);
    expect(wrapper.vm.visibleItems).toEqual(['alpha', 'beta', 'gamma']);
  });

  it('hiddenCount returns number of items beyond 3', () => {
    const wrapper = factory();
    expect(wrapper.vm.hiddenCount).toBe(2);
  });

  it('hiddenCount is 0 when items are 3 or fewer', () => {
    const wrapper = factory({ value: ['a', 'b'] });
    expect(wrapper.vm.hiddenCount).toBe(0);
  });

  it('containerClass includes justify-center when align is center', () => {
    const wrapper = factory({
      column: { ...baseColumn, align: 'center' },
    });
    expect(wrapper.vm.containerClass).toContain('justify-center');
  });

  it('containerClass includes justify-end when align is right', () => {
    const wrapper = factory({
      column: { ...baseColumn, align: 'right' },
    });
    expect(wrapper.vm.containerClass).toContain('justify-end');
  });

  it('containerClass includes justify-start by default', () => {
    const wrapper = factory();
    expect(wrapper.vm.containerClass).toContain('justify-start');
  });

  it('getColor returns static color', () => {
    const wrapper = factory();
    expect(wrapper.vm.getColor()).toBe('teal');
  });

  it('getColor calls function when color is a function', () => {
    const wrapper = factory({
      column: { ...baseColumn, color: () => 'orange' },
    });
    expect(wrapper.vm.getColor()).toBe('orange');
  });

  it('getColor defaults to "primary" when no color set', () => {
    const wrapper = factory({
      column: { ...baseColumn, color: undefined },
    });
    expect(wrapper.vm.getColor()).toBe('primary');
  });
});

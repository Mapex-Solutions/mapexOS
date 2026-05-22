import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import CodeColumn from './CodeColumn.vue';
import type { DataRowColumn } from '../interfaces';

const baseColumn: DataRowColumn = {
  key: 'uuid',
  label: 'UUID',
  type: 'code',
  visible: 'always',
};

const baseRow = { id: '1', uuid: 'abc-123-def' };

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(CodeColumn, {
    props: {
      value: 'abc-123-def',
      column: baseColumn,
      row: baseRow,
      ...overrides,
    },
  });
}

describe('CodeColumn', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('displayValue returns raw value when no format', () => {
    const wrapper = factory();
    expect(wrapper.vm.displayValue).toBe('abc-123-def');
  });

  it('displayValue returns "-" when value is falsy', () => {
    const wrapper = factory({ value: '' });
    expect(wrapper.vm.displayValue).toBe('-');
  });

  it('displayValue applies format function', () => {
    const wrapper = factory({
      column: { ...baseColumn, format: (v: any) => v.substring(0, 3) },
    });
    expect(wrapper.vm.displayValue).toBe('abc');
  });
});

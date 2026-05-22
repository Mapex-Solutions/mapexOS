import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import TextColumn from './TextColumn.vue';
import type { DataRowColumn } from '../interfaces';

const baseColumn: DataRowColumn = {
  key: 'name',
  label: 'Name',
  type: 'text',
  visible: 'always',
};

const baseRow = { id: '1', name: 'Device A', nested: { info: 'extra' } };

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(TextColumn, {
    props: {
      value: 'Device A',
      column: baseColumn,
      row: baseRow,
      ...overrides,
    },
  });
}

describe('TextColumn', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('displayValue returns raw value when no format', () => {
    const wrapper = factory();
    expect(wrapper.vm.displayValue).toBe('Device A');
  });

  it('displayValue returns "-" when value is falsy', () => {
    const wrapper = factory({ value: '' });
    expect(wrapper.vm.displayValue).toBe('-');
  });

  it('displayValue applies format function', () => {
    const wrapper = factory({
      column: { ...baseColumn, format: (v: any) => `[${v}]` },
    });
    expect(wrapper.vm.displayValue).toBe('[Device A]');
  });

  it('secondaryValue returns null when no secondary config', () => {
    const wrapper = factory();
    expect(wrapper.vm.secondaryValue).toBeNull();
  });

  it('secondaryValue calls secondary function when provided', () => {
    const wrapper = factory({
      column: { ...baseColumn, secondary: (_v: any, row: any) => row.id },
    });
    expect(wrapper.vm.secondaryValue).toBe('1');
  });

  it('secondaryValue resolves secondaryKey with dot notation', () => {
    const wrapper = factory({
      column: { ...baseColumn, secondaryKey: 'nested.info' },
    });
    expect(wrapper.vm.secondaryValue).toBe('extra');
  });

  it('secondaryValue returns null for missing secondaryKey path', () => {
    const wrapper = factory({
      column: { ...baseColumn, secondaryKey: 'nested.missing' },
    });
    expect(wrapper.vm.secondaryValue).toBeNull();
  });
});

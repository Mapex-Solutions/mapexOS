import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import ChipColumn from './ChipColumn.vue';
import type { DataRowColumn } from '../interfaces';

const baseColumn: DataRowColumn = {
  key: 'protocol',
  label: 'Protocol',
  type: 'chip',
  visible: 'always',
  color: 'blue',
};

const baseRow = { id: '1', protocol: 'MQTT' };

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(ChipColumn, {
    props: {
      value: 'MQTT',
      column: baseColumn,
      row: baseRow,
      ...overrides,
    },
  });
}

describe('ChipColumn', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('displayValue returns raw value when no format', () => {
    const wrapper = factory();
    expect(wrapper.vm.displayValue).toBe('MQTT');
  });

  it('displayValue returns "N/A" when value is falsy', () => {
    const wrapper = factory({ value: null });
    expect(wrapper.vm.displayValue).toBe('N/A');
  });

  it('displayValue applies format function', () => {
    const wrapper = factory({
      column: { ...baseColumn, format: (v: any) => `Proto: ${v}` },
    });
    expect(wrapper.vm.displayValue).toBe('Proto: MQTT');
  });

  it('getColor returns static color', () => {
    const wrapper = factory();
    expect(wrapper.vm.getColor()).toBe('blue');
  });

  it('getColor calls function color', () => {
    const wrapper = factory({
      column: { ...baseColumn, color: () => 'purple' },
    });
    expect(wrapper.vm.getColor()).toBe('purple');
  });

  it('getColor defaults to "primary" when no color set', () => {
    const wrapper = factory({
      column: { ...baseColumn, color: undefined },
    });
    expect(wrapper.vm.getColor()).toBe('primary');
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
});

import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import AvatarColumn from './AvatarColumn.vue';
import type { DataRowColumn } from '../interfaces';

const baseColumn: DataRowColumn = {
  key: 'avatar',
  label: 'Avatar',
  type: 'avatar',
  visible: 'always',
  icon: 'person',
  color: 'primary',
};

const baseRow = { id: '1', name: 'Test' };

function factory(overrides: Partial<{ value: any; column: DataRowColumn; row: any; mobile: boolean }> = {}) {
  return mountWithPlugins(AvatarColumn, {
    props: {
      value: 'test',
      column: baseColumn,
      row: baseRow,
      ...overrides,
    },
  });
}

describe('AvatarColumn', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('getIcon returns static icon string from column', () => {
    const wrapper = factory();
    expect(wrapper.vm.getIcon()).toBe('person');
  });

  it('getIcon calls function when icon is a function', () => {
    const wrapper = factory({
      column: { ...baseColumn, icon: (_v: any, row: any) => row.name },
    });
    expect(wrapper.vm.getIcon()).toBe('Test');
  });

  it('getIcon defaults to "person" when no icon is set', () => {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const { icon: _omitIcon, ...columnWithoutIcon } = baseColumn;
    const wrapper = factory({ column: columnWithoutIcon });
    expect(wrapper.vm.getIcon()).toBe('person');
  });

  it('getColor returns static color string from column', () => {
    const wrapper = factory();
    expect(wrapper.vm.getColor()).toBe('primary');
  });

  it('getColor calls function when color is a function', () => {
    const wrapper = factory({
      column: { ...baseColumn, color: () => 'red' },
    });
    expect(wrapper.vm.getColor()).toBe('red');
  });

  it('getColor defaults to "primary" when no color is set', () => {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const { color: _omitColor, ...columnWithoutColor } = baseColumn;
    const wrapper = factory({ column: columnWithoutColor });
    expect(wrapper.vm.getColor()).toBe('primary');
  });

  it('getTooltip returns null on mobile', () => {
    const wrapper = factory({
      mobile: true,
      column: { ...baseColumn, tooltip: 'Tip' },
    });
    expect(wrapper.vm.getTooltip()).toBeNull();
  });

  it('getTooltip returns null when no tooltip configured', () => {
    const wrapper = factory();
    expect(wrapper.vm.getTooltip()).toBeNull();
  });

  it('getTooltip returns static string', () => {
    const wrapper = factory({
      column: { ...baseColumn, tooltip: 'Hello' },
    });
    expect(wrapper.vm.getTooltip()).toBe('Hello');
  });

  it('getTooltip calls function when tooltip is a function', () => {
    const wrapper = factory({
      column: { ...baseColumn, tooltip: (_v: any, row: any) => row.name },
    });
    expect(wrapper.vm.getTooltip()).toBe('Test');
  });
});

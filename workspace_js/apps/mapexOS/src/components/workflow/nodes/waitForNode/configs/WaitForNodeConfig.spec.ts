import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import WaitForNodeConfig from './WaitForNodeConfig.vue';

vi.mock('@src/composables/workflow', () => ({
  useWorkflowContext: () => ({
    states: { value: [{ field: 'status', type: 'string' }, { field: 'count', type: 'number' }] },
  }),
  usePluginI18n: () => ({ t: (key: string) => key }),
}));

vi.mock('../constants', () => ({
  WAIT_FOR_OPERATORS: [
    { value: 'equals', label: 'Equals', symbol: '=' },
    { value: 'notEquals', label: 'Not Equals', symbol: '!=' },
    { value: 'isEmpty', label: 'Is Empty', symbol: 'empty' },
  ],
}));

const BASE_CONFIG: Record<string, unknown> = {};

describe('WaitForNodeConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(WaitForNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes stateFields from workflow context', () => {
    const wrapper = mountWithPlugins(WaitForNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.stateFields).toHaveLength(2);
    expect(wrapper.vm.stateFields[0].label).toBe('state.status');
  });

  it('computes fieldValue from config', () => {
    const wrapper = mountWithPlugins(WaitForNodeConfig, {
      props: { config: { field: 'status' } },
    });
    expect(wrapper.vm.fieldValue).toBe('status');
  });

  it('computes operatorValue with default equals', () => {
    const wrapper = mountWithPlugins(WaitForNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.operatorValue).toBe('equals');
  });

  it('computes compareSource with default literal', () => {
    const wrapper = mountWithPlugins(WaitForNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.compareSource).toBe('literal');
  });

  it('computes compareValue from config', () => {
    const wrapper = mountWithPlugins(WaitForNodeConfig, {
      props: { config: { compareTo: { source: 'literal', value: 'COMPLETE' } } },
    });
    expect(wrapper.vm.compareValue).toBe('COMPLETE');
  });

  it('computes intervalValue with default 30s', () => {
    const wrapper = mountWithPlugins(WaitForNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.intervalValue).toBe('30s');
  });

  it('computes timeoutValue with default 5m', () => {
    const wrapper = mountWithPlugins(WaitForNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.timeoutValue).toBe('5m');
  });

  it('computes maxCyclesValue with default 3', () => {
    const wrapper = mountWithPlugins(WaitForNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.maxCyclesValue).toBe(3);
  });

  it('computes hideCompare as false for binary operators', () => {
    const wrapper = mountWithPlugins(WaitForNodeConfig, {
      props: { config: { operator: 'equals' } },
    });
    expect(wrapper.vm.hideCompare).toBe(false);
  });

  it('computes hideCompare as true for unary operators', () => {
    const wrapper = mountWithPlugins(WaitForNodeConfig, {
      props: { config: { operator: 'isEmpty' } },
    });
    expect(wrapper.vm.hideCompare).toBe(true);
  });

  it('computes operatorOptions from constants', () => {
    const wrapper = mountWithPlugins(WaitForNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.operatorOptions).toHaveLength(3);
  });

  it('emits update:config via fieldValue setter', () => {
    const wrapper = mountWithPlugins(WaitForNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.fieldValue = 'count';
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ field: 'count' });
  });
});

import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import WaitSignalNodeConfig from './WaitSignalNodeConfig.vue';

vi.mock('@src/composables/workflow', () => ({
  useWorkflowContext: () => ({
    states: { value: [{ field: 'approved', type: 'boolean' }] },
  }),
  usePluginI18n: () => ({ t: (key: string) => key }),
}));

vi.mock('@components/tooltips/appTooltip', () => ({
  AppTooltip: { name: 'AppTooltip', template: '<div />' },
}));

vi.mock('@components/tooltips', () => ({
  AppTooltip: { name: 'AppTooltip', template: '<div />' },
}));

const BASE_CONFIG: Record<string, unknown> = {};

describe('WaitSignalNodeConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(WaitSignalNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes stateFields from workflow context', () => {
    const wrapper = mountWithPlugins(WaitSignalNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.stateFields).toHaveLength(1);
    expect(wrapper.vm.stateFields[0].label).toBe('state.approved');
  });

  it('computes signalName as empty string by default', () => {
    const wrapper = mountWithPlugins(WaitSignalNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.signalName).toBe('');
  });

  it('computes signalName from config', () => {
    const wrapper = mountWithPlugins(WaitSignalNodeConfig, {
      props: { config: { signalName: 'approval.granted' } },
    });
    expect(wrapper.vm.signalName).toBe('approval.granted');
  });

  it('computes timeoutValue with default 10m', () => {
    const wrapper = mountWithPlugins(WaitSignalNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.timeoutValue).toBe('10m');
  });

  it('computes maxCyclesValue with default 3', () => {
    const wrapper = mountWithPlugins(WaitSignalNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.maxCyclesValue).toBe(3);
  });

  it('computes mappings as empty array by default', () => {
    const wrapper = mountWithPlugins(WaitSignalNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.mappings).toEqual([]);
  });

  it('computes mappings from config', () => {
    const wrapper = mountWithPlugins(WaitSignalNodeConfig, {
      props: { config: { mappings: [{ from: 'result', to: 'approved' }] } },
    });
    expect(wrapper.vm.mappings).toHaveLength(1);
    expect(wrapper.vm.mappings[0]).toEqual({ from: 'result', to: 'approved' });
  });

  it('adds mapping via addMapping', () => {
    const wrapper = mountWithPlugins(WaitSignalNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.addMapping();
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).mappings).toEqual([{ from: '', to: '' }]);
  });

  it('removes mapping via removeMapping', () => {
    const wrapper = mountWithPlugins(WaitSignalNodeConfig, {
      props: { config: { mappings: [{ from: 'a', to: 'b' }, { from: 'c', to: 'd' }] } },
    });
    wrapper.vm.removeMapping(0);
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).mappings).toEqual([{ from: 'c', to: 'd' }]);
  });

  it('updates mapping from value via updateMappingFrom', () => {
    const wrapper = mountWithPlugins(WaitSignalNodeConfig, {
      props: { config: { mappings: [{ from: 'old', to: 'target' }] } },
    });
    wrapper.vm.updateMappingFrom(0, 'new');
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).mappings[0]).toEqual({ from: 'new', to: 'target' });
  });

  it('updates mapping to value via updateMappingTo', () => {
    const wrapper = mountWithPlugins(WaitSignalNodeConfig, {
      props: { config: { mappings: [{ from: 'src', to: 'old' }] } },
    });
    wrapper.vm.updateMappingTo(0, 'new');
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).mappings[0]).toEqual({ from: 'src', to: 'new' });
  });

  it('emits update:config via signalName setter', () => {
    const wrapper = mountWithPlugins(WaitSignalNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.signalName = 'my.signal';
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ signalName: 'my.signal' });
  });
});

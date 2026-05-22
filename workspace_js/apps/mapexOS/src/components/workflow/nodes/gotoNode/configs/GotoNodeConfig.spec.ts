import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import GotoNodeConfig from './GotoNodeConfig.vue';

vi.mock('@src/composables/workflow', () => ({
  useWorkflowContext: () => ({
    nodes: { value: [] },
  }),
  usePluginI18n: () => ({ t: (key: string) => key }),
}));

vi.mock('../constants', () => ({
  GOTO_COLOR_OPTIONS: [
    { value: 'deep-purple-6', hex: '#5e35b1' },
    { value: 'blue-6', hex: '#1e88e5' },
  ],
}));

const BASE_CONFIG: Record<string, unknown> = { _nodeId: 'goto-1' };

describe('GotoNodeConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes role as sender by default', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.role).toBe('sender');
  });

  it('computes role from config', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, {
      props: { config: { ...BASE_CONFIG, role: 'receiver' } },
    });
    expect(wrapper.vm.role).toBe('receiver');
  });

  it('computes pairLabel from config', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, {
      props: { config: { ...BASE_CONFIG, pairLabel: 'ErrorHandler' } },
    });
    expect(wrapper.vm.pairLabel).toBe('ErrorHandler');
  });

  it('computes pairColor with default', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.pairColor).toBe('deep-purple-6');
  });

  it('computes colorHex from preset', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.colorHex).toBe('#5e35b1');
  });

  it('computes isCustomColor as false for preset color', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.isCustomColor).toBe(false);
  });

  it('computes isCustomColor as true for custom hex', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, {
      props: { config: { ...BASE_CONFIG, pairColor: '#ff5722' } },
    });
    expect(wrapper.vm.isCustomColor).toBe(true);
  });

  it('emits update:config with role reset on updateRole', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateRole('receiver');
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ role: 'receiver', pairLabel: '' });
  });

  it('emits update:config on updatePairLabel', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updatePairLabel('MyLabel');
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ pairLabel: 'MyLabel' });
  });

  it('emits update:config on updatePairColor and hides hex input', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.showHexInput = true;
    wrapper.vm.updatePairColor('blue-6');
    expect(wrapper.vm.showHexInput).toBe(false);
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ pairColor: 'blue-6' });
  });

  it('applies valid hex color via applyHexColor', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.hexBuffer = '#ff5722';
    wrapper.vm.applyHexColor();
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ pairColor: '#ff5722' });
  });

  it('does not emit for invalid hex color', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.hexBuffer = 'invalid';
    wrapper.vm.applyHexColor();
    expect(wrapper.emitted('update:config')).toBeFalsy();
  });

  it('computes senderLabels from nodes', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.senderLabels).toEqual([]);
  });

  it('computes matchedPairs as empty when no pairLabel', () => {
    const wrapper = mountWithPlugins(GotoNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.matchedPairs).toEqual([]);
  });
});

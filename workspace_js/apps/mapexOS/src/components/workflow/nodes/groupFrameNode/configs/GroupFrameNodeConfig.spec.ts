import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import GroupFrameNodeConfig from './GroupFrameNodeConfig.vue';

vi.mock('@src/composables/workflow', () => ({
  usePluginI18n: () => ({ t: (key: string) => key }),
}));

vi.mock('../constants', () => ({
  FRAME_COLOR_OPTIONS: [
    { value: 'blue-grey', hex: '#78909c' },
    { value: 'teal', hex: '#009688' },
  ],
}));

const BASE_CONFIG: Record<string, unknown> = {};

describe('GroupFrameNodeConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes title as empty string by default', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.title).toBe('');
  });

  it('computes title from config', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, {
      props: { config: { title: 'My Frame' } },
    });
    expect(wrapper.vm.title).toBe('My Frame');
  });

  it('computes description from config', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, {
      props: { config: { description: 'A group' } },
    });
    expect(wrapper.vm.description).toBe('A group');
  });

  it('computes colorName with default blue-grey', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.colorName).toBe('blue-grey');
  });

  it('computes colorHex from preset', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.colorHex).toBe('#78909c');
  });

  it('computes width with default 300', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.width).toBe(300);
  });

  it('computes height with default 200', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.height).toBe(200);
  });

  it('emits update:config on updateTitle', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateTitle('New Title');
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ title: 'New Title' });
  });

  it('emits update:config on updateDescription', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateDescription('New desc');
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ description: 'New desc' });
  });

  it('emits update:config on updateColor', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateColor('teal');
    expect(wrapper.vm.showHexInput).toBe(false);
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ color: 'teal' });
  });

  it('does not emit updateWidth for values below 150', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateWidth(100);
    expect(wrapper.emitted('update:config')).toBeFalsy();
  });

  it('emits updateWidth for valid values', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateWidth(400);
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ width: 400 });
  });

  it('does not emit updateHeight for values below 100', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateHeight(50);
    expect(wrapper.emitted('update:config')).toBeFalsy();
  });

  it('applies valid hex color via applyHexColor', () => {
    const wrapper = mountWithPlugins(GroupFrameNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.hexBuffer = '#009688';
    wrapper.vm.applyHexColor();
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ color: '#009688' });
  });
});

import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import EndNodeConfig from './EndNodeConfig.vue';

vi.mock('@src/composables/workflow', () => ({
  usePluginI18n: () => ({ t: (key: string) => key }),
}));

vi.mock('@components/forms/fieldSourceSelector', () => ({
  FieldSourceSelector: { name: 'FieldSourceSelector', template: '<div />' },
}));

const BASE_CONFIG: Record<string, unknown> = {};

describe('EndNodeConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(EndNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes terminateWithError as false by default', () => {
    const wrapper = mountWithPlugins(EndNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.terminateWithError).toBe(false);
  });

  it('computes terminateWithError as true when set in config', () => {
    const wrapper = mountWithPlugins(EndNodeConfig, {
      props: { config: { terminateWithError: true } },
    });
    expect(wrapper.vm.terminateWithError).toBe(true);
  });

  it('computes errorCode from config', () => {
    const wrapper = mountWithPlugins(EndNodeConfig, {
      props: { config: { errorCode: 'ERR_001' } },
    });
    expect(wrapper.vm.errorCode).toBe('ERR_001');
  });

  it('computes errorCode as empty string by default', () => {
    const wrapper = mountWithPlugins(EndNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.errorCode).toBe('');
  });

  it('computes errorMessage with defaults when not set', () => {
    const wrapper = mountWithPlugins(EndNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.errorMessage).toEqual({ type: 'literal', value: '' });
  });

  it('computes errorMessage from config', () => {
    const wrapper = mountWithPlugins(EndNodeConfig, {
      props: { config: { errorMessage: { type: 'state', value: 'err.msg' } } },
    });
    expect(wrapper.vm.errorMessage).toEqual({ type: 'state', value: 'err.msg' });
  });

  it('emits update:config when updateTerminateWithError is called', () => {
    const wrapper = mountWithPlugins(EndNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateTerminateWithError(true);
    expect(wrapper.emitted('update:config')).toBeTruthy();
    expect(wrapper.emitted('update:config')![0]![0]).toEqual({ terminateWithError: true });
  });

  it('emits update:config when updateErrorCode is called', () => {
    const wrapper = mountWithPlugins(EndNodeConfig, {
      props: { config: { terminateWithError: true } },
    });
    wrapper.vm.updateErrorCode('ERR_100');
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ errorCode: 'ERR_100' });
  });

  it('emits update:config when handleErrorMessageUpdate is called', () => {
    const wrapper = mountWithPlugins(EndNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleErrorMessageUpdate({ type: 'input', value: 'hello' });
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ errorMessage: { type: 'input', value: 'hello' } });
  });
});

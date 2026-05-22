import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import CodeNodeConfig from './CodeNodeConfig.vue';

vi.mock('@src/composables/workflow', () => ({
  usePluginI18n: () => ({ t: (key: string) => key }),
}));

vi.mock('@components/dialogs/scriptEditor', () => ({
  ScriptEditorDialog: { name: 'ScriptEditorDialog', template: '<div />' },
}));

const BASE_CONFIG: Record<string, unknown> = {};

describe('CodeNodeConfig', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(CodeNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes script with default value when not set', () => {
    const wrapper = mountWithPlugins(CodeNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.script).toContain('return {};');
  });

  it('computes script from config', () => {
    const wrapper = mountWithPlugins(CodeNodeConfig, {
      props: { config: { script: 'return { x: 1 };' } },
    });
    expect(wrapper.vm.script).toBe('return { x: 1 };');
  });

  it('computes timeout with default 5000', () => {
    const wrapper = mountWithPlugins(CodeNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.timeout).toBe(5000);
  });

  it('computes timeout from config', () => {
    const wrapper = mountWithPlugins(CodeNodeConfig, {
      props: { config: { timeout: 10000 } },
    });
    expect(wrapper.vm.timeout).toBe(10000);
  });

  it('computes scriptPreview with first 8 lines', () => {
    const longScript = Array.from({ length: 12 }, (_, i) => `line ${i + 1}`).join('\n');
    const wrapper = mountWithPlugins(CodeNodeConfig, {
      props: { config: { script: longScript } },
    });
    expect(wrapper.vm.scriptPreview).toContain('line 8');
    expect(wrapper.vm.scriptPreview).toContain('...');
    expect(wrapper.vm.scriptPreview).not.toContain('line 9');
  });

  it('computes lineCount correctly', () => {
    const wrapper = mountWithPlugins(CodeNodeConfig, {
      props: { config: { script: 'a\nb\nc' } },
    });
    expect(wrapper.vm.lineCount).toBe(3);
  });

  it('computes editorGuidelines', () => {
    const wrapper = mountWithPlugins(CodeNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.editorGuidelines).toHaveLength(4);
    expect(wrapper.vm.editorGuidelines[0].code).toBe('state');
  });

  it('starts with editorOpen as false', () => {
    const wrapper = mountWithPlugins(CodeNodeConfig, { props: { config: BASE_CONFIG } });
    expect(wrapper.vm.editorOpen).toBe(false);
  });

  it('emits update:config on handleScriptChange', () => {
    const wrapper = mountWithPlugins(CodeNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.handleScriptChange('return { y: 2 };');
    const emitted = wrapper.emitted('update:config')!;
    expect(emitted[0]![0]).toMatchObject({ script: 'return { y: 2 };' });
  });

  it('emits update:config on updateTimeout with minimum 100', () => {
    const wrapper = mountWithPlugins(CodeNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateTimeout(50);
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).timeout).toBe(100);
  });

  it('emits update:config on updateTimeout with valid value', () => {
    const wrapper = mountWithPlugins(CodeNodeConfig, { props: { config: BASE_CONFIG } });
    wrapper.vm.updateTimeout(3000);
    const emitted = wrapper.emitted('update:config')!;
    expect((emitted[0]![0] as any).timeout).toBe(3000);
  });
});

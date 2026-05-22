import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import ScriptEditorDialog from './ScriptEditorDialog.vue';

// Mock Monaco editor
vi.mock('monaco-editor', () => ({
  default: {
    editor: {
      create: vi.fn(() => ({
        dispose: vi.fn(),
        getValue: vi.fn(() => ''),
        onDidChangeModelContent: vi.fn(),
      })),
    },
  },
  editor: {
    create: vi.fn(() => ({
      dispose: vi.fn(),
      getValue: vi.fn(() => ''),
      onDidChangeModelContent: vi.fn(),
    })),
  },
}));

// Mock monaco-theme utils
vi.mock('@utils/monaco-theme', () => ({
  registerMapexMonacoThemes: vi.fn(),
  getMapexMonacoTheme: vi.fn(() => 'mapex-light'),
  applyMapexMonacoTheme: vi.fn(),
}));

// Mock theme store
vi.mock('@stores/theme', () => ({
  useThemeStore: vi.fn(() => ({
    isDark: false,
  })),
}));

// Mock translation util
vi.mock('@utils/translation', () => ({
  useTS: vi.fn(() => (key: string) => key),
}));

describe('ScriptEditorDialog', () => {
  const defaultProps = {
    modelValue: true,
    title: 'Edit Script',
    scriptContent: 'console.log("hello");',
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(ScriptEditorDialog, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('accepts title prop', () => {
    const wrapper = mountWithPlugins(ScriptEditorDialog, { props: defaultProps });
    expect(wrapper.props('title')).toBe('Edit Script');
  });

  it('accepts scriptContent prop', () => {
    const wrapper = mountWithPlugins(ScriptEditorDialog, { props: defaultProps });
    expect(wrapper.props('scriptContent')).toBe('console.log("hello");');
  });

  it('defaults language to javascript', () => {
    const wrapper = mountWithPlugins(ScriptEditorDialog, { props: defaultProps });
    expect(wrapper.props('language')).toBe('javascript');
  });

  it('accepts custom language prop', () => {
    const wrapper = mountWithPlugins(ScriptEditorDialog, {
      props: { ...defaultProps, language: 'json' },
    });
    expect(wrapper.props('language')).toBe('json');
  });

  it('defaults guidelines to empty array', () => {
    const wrapper = mountWithPlugins(ScriptEditorDialog, { props: defaultProps });
    expect(wrapper.props('guidelines')).toEqual([]);
  });

  it('accepts guidelines prop', () => {
    const guidelines = [{ code: '{{event}}', description: 'The event payload' }];
    const wrapper = mountWithPlugins(ScriptEditorDialog, {
      props: { ...defaultProps, guidelines },
    });
    expect(wrapper.props('guidelines')).toEqual(guidelines);
  });

  it('has editorContainerRef as null initially', () => {
    const wrapper = mountWithPlugins(ScriptEditorDialog, { props: defaultProps });
    expect(wrapper.vm.editorContainerRef).toBeNull();
  });

  it('computes resolvedCloseTooltip from prop when provided', () => {
    const wrapper = mountWithPlugins(ScriptEditorDialog, {
      props: { ...defaultProps, closeTooltip: 'Custom Close' },
    });
    expect(wrapper.vm.resolvedCloseTooltip).toBe('Custom Close');
  });

  it('computes resolvedCloseTooltip from translation when not provided', () => {
    const wrapper = mountWithPlugins(ScriptEditorDialog, { props: defaultProps });
    // Falls through to the translation key since closeTooltip is ''
    expect(typeof wrapper.vm.resolvedCloseTooltip).toBe('string');
  });
});

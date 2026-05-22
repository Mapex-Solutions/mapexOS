import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import ScriptViewerDialog from './ScriptViewerDialog.vue';

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

// Mock notify utils
vi.mock('@utils/alert/notify', () => ({
  notifyInfo: vi.fn(),
  notifyFail: vi.fn(),
}));

// Mock useLogger
vi.mock('@composables/useLogger', () => ({
  useLogger: vi.fn(() => ({
    warn: vi.fn(),
    error: vi.fn(),
    info: vi.fn(),
  })),
}));

describe('ScriptViewerDialog', () => {
  const defaultProps = {
    modelValue: true,
    title: 'View Script',
    scriptContent: 'console.log("hello");',
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(ScriptViewerDialog, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('accepts title prop', () => {
    const wrapper = mountWithPlugins(ScriptViewerDialog, { props: defaultProps });
    expect(wrapper.props('title')).toBe('View Script');
  });

  it('accepts scriptContent prop', () => {
    const wrapper = mountWithPlugins(ScriptViewerDialog, { props: defaultProps });
    expect(wrapper.props('scriptContent')).toBe('console.log("hello");');
  });

  it('defaults language to javascript', () => {
    const wrapper = mountWithPlugins(ScriptViewerDialog, { props: defaultProps });
    expect(wrapper.props('language')).toBe('javascript');
  });

  it('defaults copyTooltip to "Copy script"', () => {
    const wrapper = mountWithPlugins(ScriptViewerDialog, { props: defaultProps });
    expect(wrapper.props('copyTooltip')).toBe('Copy script');
  });

  it('defaults closeTooltip to "Close"', () => {
    const wrapper = mountWithPlugins(ScriptViewerDialog, { props: defaultProps });
    expect(wrapper.props('closeTooltip')).toBe('Close');
  });

  it('accepts custom language prop', () => {
    const wrapper = mountWithPlugins(ScriptViewerDialog, {
      props: { ...defaultProps, language: 'json' },
    });
    expect(wrapper.props('language')).toBe('json');
  });

  it('has editorContainerRef as null initially', () => {
    const wrapper = mountWithPlugins(ScriptViewerDialog, { props: defaultProps });
    expect(wrapper.vm.editorContainerRef).toBeNull();
  });
});

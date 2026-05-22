import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ref } from 'vue';
import {
  MONACO_EDITOR_OPTIONS,
  createMonacoEditorManager,
  disposeMonacoEditors,
} from './monacoEditor.handler';
import type { MonacoEditorManager } from './monacoEditor.handler';
import * as monaco from 'monaco-editor';

/**
 * Mock useLogger to prevent side effects
 */
vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    debug: vi.fn(),
    error: vi.fn(),
    warn: vi.fn(),
    info: vi.fn(),
  }),
}));

/**
 * Mock useThemeStore
 */
vi.mock('@stores/theme', () => ({
  useThemeStore: () => ({ isDark: true }),
}));

describe('monacoEditor.handler', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('MONACO_EDITOR_OPTIONS', () => {
    it('has correct default language', () => {
      expect(MONACO_EDITOR_OPTIONS.language).toBe('javascript');
    });

    it('has minimap disabled', () => {
      expect(MONACO_EDITOR_OPTIONS.minimap?.enabled).toBe(false);
    });

    it('has automatic layout enabled', () => {
      expect(MONACO_EDITOR_OPTIONS.automaticLayout).toBe(true);
    });

    it('has word wrap enabled', () => {
      expect(MONACO_EDITOR_OPTIONS.wordWrap).toBe('on');
    });

    it('has font size 14', () => {
      expect(MONACO_EDITOR_OPTIONS.fontSize).toBe(14);
    });

    it('has format on paste and type enabled', () => {
      expect(MONACO_EDITOR_OPTIONS.formatOnPaste).toBe(true);
      expect(MONACO_EDITOR_OPTIONS.formatOnType).toBe(true);
    });
  });

  describe('createMonacoEditorManager', () => {
    it('returns manager with all required methods', () => {
      const containerRef = ref<HTMLElement | null>(null);
      const manager = createMonacoEditorManager(containerRef);

      expect(manager).toHaveProperty('setup');
      expect(manager).toHaveProperty('dispose');
      expect(manager).toHaveProperty('getValue');
      expect(manager).toHaveProperty('setValue');
      expect(typeof manager.setup).toBe('function');
      expect(typeof manager.dispose).toBe('function');
      expect(typeof manager.getValue).toBe('function');
      expect(typeof manager.setValue).toBe('function');
    });

    it('does not create editor when container ref is null', () => {
      const containerRef = ref<HTMLElement | null>(null);
      const createSpy = vi.spyOn(monaco.editor, 'create');

      const manager = createMonacoEditorManager(containerRef);
      manager.setup();

      expect(createSpy).not.toHaveBeenCalled();
    });

    it('creates editor when container ref has element', () => {
      const element = document.createElement('div');
      const containerRef = ref<HTMLElement | null>(element);
      const createSpy = vi.spyOn(monaco.editor, 'create');

      const manager = createMonacoEditorManager(containerRef, undefined, 'initial code');
      manager.setup();

      expect(createSpy).toHaveBeenCalledWith(
        element,
        expect.objectContaining({
          value: 'initial code',
          language: 'javascript',
        }),
      );
    });

    it('merges custom options with defaults', () => {
      const element = document.createElement('div');
      const containerRef = ref<HTMLElement | null>(element);
      const createSpy = vi.spyOn(monaco.editor, 'create');

      const manager = createMonacoEditorManager(
        containerRef,
        { language: 'typescript', fontSize: 16 },
      );
      manager.setup();

      expect(createSpy).toHaveBeenCalledWith(
        element,
        expect.objectContaining({
          language: 'typescript',
          fontSize: 16,
        }),
      );
    });

    it('disposes previous editor on re-setup', () => {
      const element = document.createElement('div');
      const containerRef = ref<HTMLElement | null>(element);
      const disposeFn = vi.fn();

      vi.spyOn(monaco.editor, 'create').mockReturnValue({
        dispose: disposeFn,
        getValue: () => '',
        setValue: () => {},
        onDidChangeModelContent: () => ({ dispose: () => {} }),
        getModel: () => null,
        updateOptions: () => {},
      } as any);

      const manager = createMonacoEditorManager(containerRef);
      manager.setup();
      manager.setup(); // second call should dispose first

      expect(disposeFn).toHaveBeenCalledTimes(1);
    });

    it('getValue returns empty string when no editor', () => {
      const containerRef = ref<HTMLElement | null>(null);
      const manager = createMonacoEditorManager(containerRef);

      expect(manager.getValue()).toBe('');
    });

    it('setValue does nothing when no editor', () => {
      const containerRef = ref<HTMLElement | null>(null);
      const manager = createMonacoEditorManager(containerRef);

      // Should not throw
      expect(() => manager.setValue('test')).not.toThrow();
    });

    it('dispose cleans up editor instance', () => {
      const element = document.createElement('div');
      const containerRef = ref<HTMLElement | null>(element);
      const disposeFn = vi.fn();

      vi.spyOn(monaco.editor, 'create').mockReturnValue({
        dispose: disposeFn,
        getValue: () => '',
        setValue: () => {},
        onDidChangeModelContent: () => ({ dispose: () => {} }),
        getModel: () => null,
        updateOptions: () => {},
      } as any);

      const manager = createMonacoEditorManager(containerRef);
      manager.setup();
      manager.dispose();

      expect(disposeFn).toHaveBeenCalledTimes(1);
    });

    it('dispose is safe to call multiple times', () => {
      const containerRef = ref<HTMLElement | null>(null);
      const manager = createMonacoEditorManager(containerRef);

      // Should not throw on double dispose
      expect(() => {
        manager.dispose();
        manager.dispose();
      }).not.toThrow();
    });

    it('calls onContentChange callback when content changes', () => {
      const element = document.createElement('div');
      const containerRef = ref<HTMLElement | null>(element);
      const onContentChange = vi.fn();
      let changeCallback: (() => void) | null = null;

      vi.spyOn(monaco.editor, 'create').mockReturnValue({
        dispose: vi.fn(),
        getValue: () => 'new content',
        setValue: vi.fn(),
        onDidChangeModelContent: (cb: () => void) => {
          changeCallback = cb;
          return { dispose: () => {} };
        },
        getModel: () => null,
        updateOptions: () => {},
      } as any);

      const manager = createMonacoEditorManager(containerRef, undefined, '', onContentChange);
      manager.setup();

      // Simulate content change
      expect(changeCallback).not.toBeNull();
      changeCallback!();

      expect(onContentChange).toHaveBeenCalledWith('new content');
    });
  });

  describe('disposeMonacoEditors', () => {
    it('calls dispose on all managers', () => {
      const disposeFn1 = vi.fn();
      const disposeFn2 = vi.fn();

      const managers: MonacoEditorManager[] = [
        { editor: null, setup: vi.fn(), dispose: disposeFn1, getValue: vi.fn(), setValue: vi.fn() },
        { editor: null, setup: vi.fn(), dispose: disposeFn2, getValue: vi.fn(), setValue: vi.fn() },
      ];

      disposeMonacoEditors(managers);

      expect(disposeFn1).toHaveBeenCalledTimes(1);
      expect(disposeFn2).toHaveBeenCalledTimes(1);
    });

    it('handles empty array', () => {
      expect(() => disposeMonacoEditors([])).not.toThrow();
    });
  });
});

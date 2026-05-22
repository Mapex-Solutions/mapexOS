/**
 * Monaco Editor Handler
 * Manages Monaco Editor instances creation, disposal, and configuration
 */

import * as monaco from 'monaco-editor';
import type { Ref } from 'vue';
import { useLogger } from '@composables/useLogger';
import { registerMapexMonacoThemes, getMapexMonacoTheme } from '@utils/monaco-theme';
import { useThemeStore } from '@stores/theme';

const logger = useLogger('MonacoEditorHandler');

/**
 * Default Monaco Editor options (theme is set dynamically via getMapexMonacoTheme)
 */
export const MONACO_EDITOR_OPTIONS: monaco.editor.IStandaloneEditorConstructionOptions = {
  language: 'javascript',
  automaticLayout: true,
  minimap: { enabled: false },
  scrollBeyondLastLine: false,
  wordWrap: 'on',
  fontSize: 14,
  lineNumbers: 'on',
  roundedSelection: false,
  cursorStyle: 'line',
  formatOnPaste: true,
  formatOnType: true,
};

/**
 * Monaco Editor Manager Interface
 */
export interface MonacoEditorManager {
  editor: monaco.editor.IStandaloneCodeEditor | null;
  setup: () => void;
  dispose: () => void;
  getValue: () => string;
  setValue: (value: string) => void;
}

/**
 * Create Monaco Editor Manager
 * Manages a single Monaco Editor instance lifecycle
 *
 * @param containerRef - Ref to the DOM element container
 * @param options - Monaco editor options (optional)
 * @param initialValue - Initial editor content
 * @param onContentChange - Callback when content changes
 * @returns MonacoEditorManager instance
 */
export function createMonacoEditorManager(
  containerRef: Ref<HTMLElement | null>,
  options?: Partial<monaco.editor.IStandaloneEditorConstructionOptions>,
  initialValue = '',
  onContentChange?: (value: string) => void
): MonacoEditorManager {
  let editorInstance: monaco.editor.IStandaloneCodeEditor | null = null;

  const setup = () => {
    if (!containerRef.value) {
      logger.warn('Container ref is null');
      return;
    }

    // Dispose previous instance if exists
    if (editorInstance) {
      editorInstance.dispose();
      editorInstance = null;
    }

    // Register themes and create new instance with correct theme
    registerMapexMonacoThemes();
    const themeStore = useThemeStore();
    editorInstance = monaco.editor.create(containerRef.value, {
      ...MONACO_EDITOR_OPTIONS,
      theme: getMapexMonacoTheme(themeStore.isDark),
      ...options,
      value: initialValue,
    });

    // Setup content change listener
    if (onContentChange) {
      editorInstance.onDidChangeModelContent(() => {
        const value = editorInstance?.getValue() || '';
        onContentChange(value);
      });
    }
  };

  const dispose = () => {
    if (editorInstance) {
      editorInstance.dispose();
      editorInstance = null;
    }
  };

  const getValue = (): string => {
    return editorInstance?.getValue() || '';
  };

  const setValue = (value: string) => {
    if (editorInstance) {
      editorInstance.setValue(value);
    }
  };

  return {
    editor: editorInstance,
    setup,
    dispose,
    getValue,
    setValue,
  };
}

/**
 * Dispose multiple Monaco Editor instances
 *
 * @param managers - Array of MonacoEditorManager instances
 */
export function disposeMonacoEditors(managers: MonacoEditorManager[]): void {
  managers.forEach(manager => manager.dispose());
}

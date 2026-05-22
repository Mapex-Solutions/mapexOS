/**
 * Centralized Monaco Editor theme definitions for MapexOS
 *
 * Provides both light and dark themes that match the MapexOS CSS variable system.
 * Includes token rules for JSON and JavaScript languages.
 *
 * Usage:
 *   import { registerMapexMonacoThemes, getMapexMonacoTheme, applyMapexMonacoTheme } from '@utils/monaco-theme';
 *   registerMapexMonacoThemes();
 *   monaco.editor.create(container, { theme: getMapexMonacoTheme(isDark) });
 */

import { monaco } from './monaco-setup';

/** Monaco theme name for the light variant */
export const MAPEX_THEME_LIGHT = 'mapex-light';

/** Monaco theme name for the dark variant */
export const MAPEX_THEME_DARK = 'mapex-dark';

/** Guard to prevent duplicate registration */
let themesRegistered = false;

/**
 * Register both MapexOS Monaco themes (light + dark).
 * Safe to call multiple times — only registers once.
 *
 * @returns {void}
 */
export function registerMapexMonacoThemes(): void {
  if (themesRegistered) return;

  monaco.editor.defineTheme(MAPEX_THEME_LIGHT, {
    base: 'vs',
    inherit: true,
    rules: [
      // JSON tokens
      { token: 'string.key.json', foreground: '0066CC' },
      { token: 'string.value.json', foreground: '008000' },
      { token: 'number.json', foreground: '0000FF' },
      { token: 'keyword.json', foreground: '800080' },
      { token: 'delimiter', foreground: '000000' },
      // JavaScript tokens
      { token: 'keyword', foreground: '0000FF' },
      { token: 'string', foreground: 'A31515' },
      { token: 'number', foreground: '098658' },
      { token: 'comment', foreground: '008000' },
      { token: 'function', foreground: '795E26' },
    ],
    colors: {
      'editor.background': '#FAFAFA',
      'editor.foreground': '#333333',
      'editorLineNumber.foreground': '#999999',
      'editor.selectionBackground': '#ADD6FF',
      'editorCursor.foreground': '#000000',
      'editor.lineHighlightBackground': '#F0F0F0',
    },
  });

  monaco.editor.defineTheme(MAPEX_THEME_DARK, {
    base: 'vs-dark',
    inherit: true,
    rules: [
      // JSON tokens
      { token: 'string.key.json', foreground: '9CDCFE' },
      { token: 'string.value.json', foreground: 'CE9178' },
      { token: 'number.json', foreground: 'B5CEA8' },
      { token: 'keyword.json', foreground: '569CD6' },
      { token: 'delimiter', foreground: 'D4D4D4' },
      // JavaScript tokens
      { token: 'keyword', foreground: '569CD6' },
      { token: 'string', foreground: 'CE9178' },
      { token: 'number', foreground: 'B5CEA8' },
      { token: 'comment', foreground: '6A9955' },
      { token: 'function', foreground: 'DCDCAA' },
    ],
    colors: {
      'editor.background': '#232333',
      'editor.foreground': '#D4D4D4',
      'editorLineNumber.foreground': '#6E6E7E',
      'editor.selectionBackground': '#264F78',
      'editorCursor.foreground': '#E0E0E0',
      'editor.lineHighlightBackground': '#2A2A3E',
    },
  });

  themesRegistered = true;
}

/**
 * Get the MapexOS Monaco theme name for the current mode
 *
 * @param {boolean} isDark - Whether dark mode is active
 * @returns {string} Theme name to use with Monaco editor
 */
export function getMapexMonacoTheme(isDark: boolean): string {
  return isDark ? MAPEX_THEME_DARK : MAPEX_THEME_LIGHT;
}

/**
 * Apply the MapexOS Monaco theme globally (affects ALL editor instances)
 *
 * @param {boolean} isDark - Whether dark mode is active
 * @returns {void}
 */
export function applyMapexMonacoTheme(isDark: boolean): void {
  registerMapexMonacoThemes();
  monaco.editor.setTheme(getMapexMonacoTheme(isDark));
}

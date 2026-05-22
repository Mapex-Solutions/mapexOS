import type { ThemeMode, ThemeState } from './types';

const STORAGE_KEY = 'mapex-theme';

export const actions = {
  /**
   * Set theme mode, resolve it, and persist to localStorage
   *
   * @param {ThemeMode} mode - 'light', 'dark', or 'system'
   */
  setTheme(mode: ThemeMode): void {
    const store = this as ThemeState & typeof actions;

    store.mode = mode;
    localStorage.setItem(STORAGE_KEY, mode);

    const resolved = mode === 'system'
      ? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
      : mode;

    store.resolvedTheme = resolved;
  },

  /**
   * Toggle between light and dark (skips system)
   */
  toggleTheme(): void {
    const store = this as ThemeState & typeof actions;

    const next = store.resolvedTheme === 'dark' ? 'light' : 'dark';
    store.setTheme(next);
  },

  /**
   * Initialize theme on app startup — reads localStorage + resolves system preference
   */
  initTheme(): void {
    const store = this as ThemeState & typeof actions;

    const saved = localStorage.getItem(STORAGE_KEY) as ThemeMode | null;
    store.setTheme(saved || 'system');
  },
};

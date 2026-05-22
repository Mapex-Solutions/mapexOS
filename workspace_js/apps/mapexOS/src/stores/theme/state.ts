import type { ThemeState } from './types';

/** localStorage key used to persist the user's theme preference */
const STORAGE_KEY = 'mapex-theme';

/**
 * Theme state factory.
 * Reads the persisted mode from localStorage (defaults to 'system').
 *
 * @returns {ThemeState} Initial theme state
 */
export const state = (): ThemeState => ({
  mode: (localStorage.getItem(STORAGE_KEY) as ThemeState['mode']) || 'system',
  resolvedTheme: 'light',
});

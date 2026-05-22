import type { ThemeState, ThemeMode } from './types';

export const getters = {
  /**
   * Whether the resolved theme is dark
   *
   * @param {ThemeState} state - Current theme state
   * @returns {boolean} True when resolved theme is 'dark'
   */
  isDark: (state: ThemeState): boolean => state.resolvedTheme === 'dark',

  /**
   * Current user-selected mode ('light', 'dark', or 'system')
   *
   * @param {ThemeState} state - Current theme state
   * @returns {ThemeMode} The persisted theme mode preference
   */
  currentMode: (state: ThemeState): ThemeMode => state.mode,
};

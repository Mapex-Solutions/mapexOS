/** Supported theme modes */
export type ThemeMode = 'light' | 'dark' | 'system';

/** Theme store state */
export interface ThemeState {
  /** User's preferred theme mode */
  mode: ThemeMode;
  /** Resolved active theme (after system detection) */
  resolvedTheme: 'light' | 'dark';
}

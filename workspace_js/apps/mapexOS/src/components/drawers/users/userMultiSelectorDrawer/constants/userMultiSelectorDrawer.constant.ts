/**
 * UserMultiSelectorDrawer Constants
 */

/** Default pagination settings */
export const USER_SELECTOR_DEFAULTS = {
  /** Items per page */
  PER_PAGE: 20,

  /** Debounce delay for search input (ms) */
  DEBOUNCE_MS: 300,
} as const;

/** Filter mode options for search toggle */
export const FILTER_MODE_OPTIONS = [
  { value: 'name' as const, icon: 'person', tooltip: 'Search by name' },
  { value: 'email' as const, icon: 'email', tooltip: 'Search by email' },
] as const;

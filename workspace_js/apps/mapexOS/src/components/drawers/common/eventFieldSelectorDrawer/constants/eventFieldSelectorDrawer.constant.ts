/**
 * Default debounce delay for search input (milliseconds)
 */
export const SEARCH_DEBOUNCE_DELAY = 300;

/**
 * Empty state messages
 */
export const EMPTY_STATE_MESSAGES = {
  NO_TEMPLATES: 'No templates selected. Please select templates first.',
  NO_FIELDS: 'No fields available from selected templates.',
  NO_RESULTS: 'No fields match your search.',
} as const;

/**
 * Loading state messages
 */
export const LOADING_MESSAGES = {
  FETCHING: 'Loading available fields...',
} as const;

/**
 * Maximum number of fields to display before scrolling
 */
export const MAX_VISIBLE_FIELDS = 50;

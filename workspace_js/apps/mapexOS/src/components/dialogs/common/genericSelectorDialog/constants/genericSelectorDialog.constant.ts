/** Default dialog card width in pixels */
export const DEFAULT_DIALOG_WIDTH = 600;

/** Default search input debounce in milliseconds */
export const DEFAULT_SEARCH_DEBOUNCE_MS = 500;

/** Default scroll threshold (0-1) to trigger load-more */
export const DEFAULT_SCROLL_THRESHOLD = 0.8;

/** Default property key used to identify items */
export const DEFAULT_ITEM_KEY = 'id';

/** Default icon for empty state */
export const DEFAULT_EMPTY_ICON = 'inbox';

/** Default text for empty state */
export const DEFAULT_EMPTY_TEXT = 'No results found';

/** Default text shown during initial loading */
export const DEFAULT_LOADING_TEXT = 'Loading...';

/** Default confirm button label */
export const DEFAULT_CONFIRM_LABEL = 'Confirm Selection';

/** Default cancel button label */
export const DEFAULT_CANCEL_LABEL = 'Cancel';

/** Default singular noun for footer count */
export const DEFAULT_ITEM_NOUN_SINGULAR = 'item';

/** Default plural noun for footer count */
export const DEFAULT_ITEM_NOUN_PLURAL = 'items';

/** Default info banner styling (uses CSS class for theme-aware colors) */
export const DEFAULT_INFO_BANNER = {
  icon: 'info',
  bgClass: 'selector-banner',
  textClass: '',
  iconColor: 'primary',
} as const;

/** Default active item highlight colors (uses MapexOS design tokens) */
export const DEFAULT_ACTIVE_ITEM_STYLE = {
  backgroundColor: 'var(--mapex-active-bg)',
  borderColor: 'var(--mapex-active-border)',
} as const;

/**
 * Scroll info from q-scroll-area
 */
export interface ScrollInfo {
  /** Current vertical scroll position */
  verticalPosition: number;

  /** Total vertical scrollable size */
  verticalSize: number;

  /** Visible container height */
  verticalContainerSize: number;
}

/**
 * Info banner configuration for the selector dialog
 */
export interface SelectorInfoBanner {
  /** Banner text content */
  text: string;

  /** Material icon name (default: 'info') */
  icon?: string;

  /** Quasar background class (default: 'bg-teal-1') */
  bgClass?: string;

  /** Quasar text class (default: 'text-teal-9') */
  textClass?: string;

  /** Icon color (default: 'teal-6') */
  iconColor?: string;
}

/**
 * Active item highlight style overrides
 */
export interface ActiveItemStyle {
  /** Background color for active items (default: 'rgba(0, 150, 136, 0.08)') */
  backgroundColor?: string;

  /** Left border color for active items (default: 'var(--q-primary)') */
  borderColor?: string;
}

/**
 * Scoped slot props for the #item slot
 */
export interface ItemSlotProps {
  /** The current item being rendered */
  item: any;

  /** Whether the item is currently selected */
  isSelected: boolean;

  /** Toggle function to select/deselect the item */
  toggle: () => void;
}

/**
 * Props for GenericSelectorDialog component
 * Encapsulates boilerplate for selector dialogs with search, filters, infinite scroll
 */
export interface GenericSelectorDialogProps {
  /** Controls dialog visibility (v-model) */
  modelValue: boolean;

  /** Dialog header title */
  title: string;

  /** Icon name displayed before title */
  icon?: string;

  /** Icon color (Quasar color name, default: 'primary') */
  iconColor?: string;

  /** Array of items to display in the list */
  items: any[];

  /** Key property to identify items (default: 'id') */
  itemKey?: string;

  /** Enable multi-select mode with checkboxes and confirm button (default: false) */
  multiSelect?: boolean;

  /** Array of pre-selected item IDs */
  selectedIds?: string[];

  /** Whether initial data is loading */
  loading?: boolean;

  /** Whether additional pages are loading */
  loadingMore?: boolean;

  /** Total count of items (for footer display) */
  totalItems?: number;

  /** Whether more pages are available for infinite scroll */
  hasMorePages?: boolean;

  /** Placeholder text for search input (default: 'Search...') */
  searchPlaceholder?: string;

  /** Info banner configuration (shown below header) */
  infoBanner?: SelectorInfoBanner;

  /** Dialog card width in pixels (default: 600) */
  width?: number;

  /** Search input debounce in milliseconds (default: 500) */
  searchDebounce?: number;

  /** Icon for empty state (default: 'inbox') */
  emptyIcon?: string;

  /** Text for empty state (default: 'No results found') */
  emptyText?: string;

  /** Text shown during initial loading (default: 'Loading...') */
  loadingText?: string;

  /** Confirm button label (default: 'Confirm Selection') */
  confirmLabel?: string;

  /** Cancel button label (default: 'Cancel') */
  cancelLabel?: string;

  /** Singular noun for item count in footer (default: 'item') */
  itemNounSingular?: string;

  /** Plural noun for item count in footer (default: 'items') */
  itemNounPlural?: string;

  /** Icon shown in footer before count (default: undefined) */
  footerIcon?: string;

  /** Icon shown in results header (default: undefined) */
  resultsIcon?: string;

  /** Custom active item highlight colors */
  activeItemStyle?: ActiveItemStyle;

  /** Show "Filters" overline header above search/filters (default: true) */
  showFiltersHeader?: boolean;

  /** Show built-in search input (default: true) */
  showSearch?: boolean;

  /** Scroll position threshold (0-1) to trigger load-more (default: 0.8) */
  scrollThreshold?: number;
}

/**
 * Emits for GenericSelectorDialog component
 */
export interface GenericSelectorDialogEmits {
  /** Update dialog visibility */
  (e: 'update:modelValue', value: boolean): void;

  /** Emit selected items (single-item array for single-select, multiple for multi-select) */
  (e: 'select', items: any[]): void;

  /** Emit cancel action */
  (e: 'cancel'): void;

  /** Emit search query after debounce */
  (e: 'search', query: string): void;

  /** Request next page of items */
  (e: 'load-more'): void;
}

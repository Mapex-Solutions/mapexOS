/**
 * List Header Menu Component Interfaces
 *
 * Generic menu component for list pages that provides:
 * - Items count display
 * - Items per page selection
 * - Column visibility toggles
 */

/**
 * Column visibility option
 */
export interface ListHeaderMenuColumn {
  /** Unique column key */
  key: string;
  /** Display label for the column */
  label: string;
  /** Whether this column is currently visible */
  visible: boolean;
}

/**
 * Props for ListHeaderMenu component
 */
export interface ListHeaderMenuProps {
  /** Total items count to display in button label */
  itemsCount: number;
  /** Singular label (e.g., "Asset", "Rule", "User") */
  itemLabel: string;
  /** Plural label (e.g., "Assets", "Rules", "Users") */
  itemLabelPlural?: string;
  /** Icon name for the button */
  icon?: string;
  /** Current items per page value */
  itemsPerPage: number;
  /** Available items per page options */
  itemsPerPageOptions?: number[];
  /** Column visibility configuration */
  columns?: ListHeaderMenuColumn[];
  /** Show items per page section */
  showItemsPerPage?: boolean;
  /** Show column visibility section */
  showColumnVisibility?: boolean;
  /** Indicates if results are filtered */
  filtered?: boolean;
  /** Show refresh button (default: true) */
  showRefresh?: boolean;
  /** Disables refresh button and spins the icon while the list is loading */
  refreshing?: boolean;
  /** Timestamp (Date or ms) of the last successful fetch — drives the "Updated Xs ago" caption */
  lastUpdatedAt?: Date | number | undefined;
}

/**
 * Emits for ListHeaderMenu component
 */
export interface ListHeaderMenuEmits {
  /** Emitted when items per page changes */
  (event: 'update:itemsPerPage', value: number): void;
  /** Emitted when column visibility changes */
  (event: 'update:columns', columns: ListHeaderMenuColumn[]): void;
  /** Emitted when the refresh button is clicked */
  (event: 'refresh'): void;
}

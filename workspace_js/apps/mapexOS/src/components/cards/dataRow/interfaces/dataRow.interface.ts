export type DataRowColumnType = 'text' | 'chip' | 'chips' | 'badge' | 'code' | 'date' | 'icon' | 'avatar';

/**
 * Column visibility using intuitive aliases and Quasar breakpoint names
 *
 * PREFERRED: Use intuitive device aliases:
 * - 'mobile': Show only on mobile (< 1024px) - maps to 'lt-md'
 * - 'tablet': Show on tablet and up (>= 1024px, < 1440px) - maps to 'gt-sm'
 * - 'laptop': Show on laptop and up (>= 1024px) - maps to 'gt-sm'
 * - 'desktop': Show only on desktop (>= 1440px) - maps to 'gt-md'
 *
 * ADVANCED: Quasar breakpoint names for fine-tuning:
 * - xs: < 600px
 * - sm: >= 600px && < 1024px
 * - md: >= 1024px && < 1440px
 * - lg: >= 1440px && < 1920px
 * - xl: >= 1920px
 *
 * Special values:
 * - 'always': Show on all screen sizes
 * - 'gt-xs': Show on >= 600px (greater than xs)
 * - 'gt-sm': Show on >= 1024px (greater than sm)
 * - 'gt-md': Show on >= 1440px (greater than md)
 * - 'lt-sm': Show on < 600px (less than sm)
 * - 'lt-md': Show on < 1024px (less than md)
 * - 'lt-lg': Show on < 1440px (less than lg)
 *
 * @example
 * // ✅ PREFERRED: Use intuitive aliases
 * visible: 'always'   // Show everywhere (essential columns)
 * visible: 'laptop'   // Show on laptop and up (>= 1024px) - MOST COMMON
 * visible: 'desktop'  // Show only on desktop (>= 1440px) - details
 * visible: 'mobile'   // Show only on mobile (< 1024px) - rare
 *
 * // ⚙️ ADVANCED: Fine-tune with Quasar classes
 * visible: 'gt-sm'    // Show on >= 1024px (md, lg, xl)
 * visible: 'gt-md'    // Show on >= 1440px (lg, xl)
 * visible: 'xs'       // Show only on < 600px
 */
export type DataRowColumnVisibility =
  | 'always'
  // Intuitive device aliases (PREFERRED)
  | 'mobile' | 'tablet' | 'laptop' | 'desktop'
  // Quasar exact sizes (advanced)
  | 'xs' | 'sm' | 'md' | 'lg' | 'xl'
  // Quasar comparison helpers (advanced)
  | 'gt-xs' | 'gt-sm' | 'gt-md' | 'gt-lg'
  | 'lt-sm' | 'lt-md' | 'lt-lg' | 'lt-xl';

/**
 * DataRow column definition following Quasar Table pattern
 *
 * Value Access Flow:
 * 1. DataRow extracts raw value using 'key' (supports dot notation: 'protocol.type')
 * 2. Column component receives raw value as prop
 * 3. Column component applies 'format' if defined
 *
 * @example
 * ```typescript
 * {
 *   key: 'protocol.type',           // Nested property access
 *   label: 'Protocol',
 *   type: 'chip',
 *   format: (val) => val?.toUpperCase(), // Applied by column component
 *   color: (val) => val === 'mqtt' ? 'purple' : 'blue'
 * }
 * ```
 */
export interface DataRowColumn {
  /**
   * Key to access value in data object
   * Supports dot notation for nested properties (e.g., 'protocol.type')
   */
  key: string;

  /** Column header label */
  label: string;

  /** Column type determines which component renders the value */
  type: DataRowColumnType;

  /** When to show this column (responsive behavior) */
  visible: DataRowColumnVisibility;

  /** Fixed column width in pixels */
  width?: number;

  /**
   * Color for the column component
   * Can be static string or function that receives (value, row)
   */
  color?: string | ((value: any, row: any) => string);

  /** Enable text ellipsis for overflow */
  ellipsis?: boolean;

  /**
   * Format function applied by column components
   * Receives raw value extracted by DataRow
   * Returns formatted string for display
   */
  format?: (value: any, row: any) => string;

  /**
   * Icon name or function to determine icon
   * Used by avatar and icon column types
   */
  icon?: string | ((value: any, row: any) => string);

  /**
   * Key for secondary text (shown below main text)
   * Used by text column type
   */
  secondaryKey?: string;

  /**
   * Format function for secondary text (shown below main text)
   * Similar to format but for secondary line
   * Receives raw value and row, returns formatted string
   * Used by text column type
   */
  secondary?: (value: any, row: any) => string;

  /**
   * Tooltip text or function to show on hover
   * Displayed only on laptop/desktop (>= 1024px)
   * Can be static string or function that receives (value, row)
   * Commonly used with avatar column to show type (site, zone, floor, etc.)
   *
   * @example
   * ```typescript
   * tooltip: (val, row) => row.type // Shows: "site", "zone", "floor"
   * tooltip: 'Click to view details' // Static text
   * ```
   */
  tooltip?: string | ((value: any, row: any) => string);

  /**
   * Text alignment for the column content
   * - 'left': Align content to the left (default)
   * - 'center': Center the content
   * - 'right': Align content to the right
   */
  align?: 'left' | 'center' | 'right';
}

/**
 * Configuration for a custom action in the DataRow menu
 */
export interface DataRowCustomAction {
  /** Unique key for the action */
  key: string;

  /** Display label for the action */
  label: string;

  /** Material icon name */
  icon: string;

  /** Icon color (Quasar color name) */
  color?: string;

  /** Description text shown below the label */
  description?: string;

  /**
   * Optional condition function to show/hide action for specific rows
   * @param row - The row data
   * @returns boolean indicating if action should be shown
   */
  condition?: (row: any) => boolean;
}

/**
 * Configuration for DataRow actions menu
 */
export interface DataRowActionConfig {
  showEdit?: boolean;
  showView?: boolean;
  showDelete?: boolean;

  /**
   * Custom actions to display in the menu
   * These are rendered before the standard actions (Edit, View, Delete)
   */
  customActions?: DataRowCustomAction[];
}

export interface DataRowProps {
  data: any;
  columns: DataRowColumn[];
  primaryKey?: string;
  showExpand?: boolean;
  showActions?: boolean;
  expandOnClick?: boolean;
  actions?: DataRowActionConfig;
}

export interface DataRowEmits {
  (e: 'click', data: any): void;
  (e: 'dblclick', data: any): void;
  (e: 'edit', data: any): void;
  (e: 'view', data: any): void;
  (e: 'delete', data: any): void;
  (e: 'expand', data: any, expanded: boolean): void;
  (e: 'action', key: string, data: any): void;
}

/**
 * DetailChip component props interface
 * Used for displaying information chips in detail drawers with proper spacing and styling
 */
export interface DetailChipProps {
  /** Chip label text (use either label or value) */
  label?: string | undefined;

  /** Chip value text (alias for label) */
  value?: string | number | undefined;

  /** Optional icon name (Material Icons) */
  icon?: string | undefined;

  /** Chip color variant */
  color?: DetailChipColor;

  /** Chip size variant */
  size?: DetailChipSize;

  /** Text color (defaults to white for colored chips, grey-9 for default) */
  textColor?: string;

  /** Make chip clickable */
  clickable?: boolean;

  /** Show outline style instead of filled */
  outline?: boolean;

  /** Make chip square instead of rounded */
  square?: boolean;

  /** Compact mode (uses Quasar's default dense sizing) */
  dense?: boolean;

  /** Rounded corners (default: true) */
  rounded?: boolean;
}

/**
 * Available color variants for DetailChip
 */
export type DetailChipColor =
  | 'primary'
  | 'secondary'
  | 'positive'
  | 'negative'
  | 'warning'
  | 'info'
  | 'teal'
  | 'cyan'
  | 'blue'
  | 'indigo'
  | 'purple'
  | 'deep-purple'
  | 'pink'
  | 'red'
  | 'orange'
  | 'deep-orange'
  | 'amber'
  | 'yellow'
  | 'lime'
  | 'green'
  | 'grey'
  | 'blue-grey'
  | 'default';

/**
 * Available size variants for DetailChip
 */
export type DetailChipSize = 'xs' | 'sm' | 'md' | 'lg';

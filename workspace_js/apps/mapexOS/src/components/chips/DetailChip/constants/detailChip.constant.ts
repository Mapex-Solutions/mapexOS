import type { DetailChipColor, DetailChipSize } from '../interfaces';

/**
 * Color mapping for DetailChip variants
 * Maps color names to Quasar color values with proper contrast
 */
export const DETAIL_CHIP_COLORS: Record<DetailChipColor, string> = {
  primary: 'primary',
  secondary: 'secondary',
  positive: 'positive',
  negative: 'negative',
  warning: 'warning',
  info: 'info',
  teal: 'teal-6',
  cyan: 'cyan-6',
  blue: 'blue-6',
  indigo: 'indigo-6',
  purple: 'purple-6',
  'deep-purple': 'deep-purple-6',
  pink: 'pink-6',
  red: 'red-6',
  orange: 'orange-6',
  'deep-orange': 'deep-orange-6',
  amber: 'amber-6',
  yellow: 'yellow-6',
  lime: 'lime-6',
  green: 'green-6',
  grey: 'grey-6',
  'blue-grey': 'blue-grey-6',
  default: 'grey-3',
};

/**
 * Text color mapping for DetailChip variants
 * Provides proper contrast for each background color
 */
export const DETAIL_CHIP_TEXT_COLORS: Record<DetailChipColor, string> = {
  primary: 'white',
  secondary: 'white',
  positive: 'white',
  negative: 'white',
  warning: 'white',
  info: 'white',
  teal: 'white',
  cyan: 'white',
  blue: 'white',
  indigo: 'white',
  purple: 'white',
  'deep-purple': 'white',
  pink: 'white',
  red: 'white',
  orange: 'white',
  'deep-orange': 'white',
  amber: 'grey-9',
  yellow: 'grey-9',
  lime: 'grey-9',
  green: 'white',
  grey: 'white',
  'blue-grey': 'white',
  default: 'grey-9',
};

/**
 * Size configuration for DetailChip
 * Defines font size and padding for each size variant
 */
export const DETAIL_CHIP_SIZES: Record<
  DetailChipSize,
  {
    fontSize: string;
    padding: string;
    iconSize: string;
  }
> = {
  xs: {
    fontSize: '0.65rem',
    padding: '2px 8px',
    iconSize: '14px',
  },
  sm: {
    fontSize: '0.75rem',
    padding: '4px 10px',
    iconSize: '16px',
  },
  md: {
    fontSize: '0.85rem',
    padding: '6px 12px',
    iconSize: '18px',
  },
  lg: {
    fontSize: '0.95rem',
    padding: '8px 16px',
    iconSize: '20px',
  },
};

/**
 * Default props for DetailChip component
 */
export const DETAIL_CHIP_DEFAULTS = {
  color: 'default' as DetailChipColor,
  size: 'md' as DetailChipSize,
  clickable: false,
  outline: false,
  square: false,
  dense: false,
  rounded: true,
};

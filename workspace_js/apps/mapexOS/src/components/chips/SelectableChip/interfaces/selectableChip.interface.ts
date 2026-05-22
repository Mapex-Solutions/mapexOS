import type { DetailChipColor, DetailChipSize } from '../../DetailChip';

/**
 * SelectableChip component props interface
 * Used for interactive/removable chips in drawers and filters
 */
export interface SelectableChipProps {
  /** Chip label text */
  label?: string | undefined;

  /** Optional icon name (Material Icons) */
  icon?: string | undefined;

  /** Chip color variant (same as DetailChip) */
  color?: DetailChipColor;

  /** Chip size variant */
  size?: DetailChipSize;

  /** Text color (defaults based on background color) */
  textColor?: string;

  /** Show remove button */
  removable?: boolean;

  /** Compact mode */
  dense?: boolean;

  /** Disable chip interactions */
  disable?: boolean;

  /** Make chip clickable */
  clickable?: boolean;

  /** Show outline style instead of filled */
  outline?: boolean;

  /** Make chip square instead of rounded */
  square?: boolean;
}

/**
 * SelectableChip component emits interface
 */
export interface SelectableChipEmits {
  /** Emitted when remove button is clicked */
  (e: 'remove'): void;

  /** Emitted when chip is clicked (if clickable) */
  (e: 'click'): void;
}

import type { DetailChipColor, DetailChipSize } from '../../DetailChip';

/**
 * Default props for SelectableChip component
 */
export const SELECTABLE_CHIP_DEFAULTS = {
  color: 'primary' as DetailChipColor,
  size: 'sm' as DetailChipSize,
  removable: true,
  dense: false,
  disable: false,
  clickable: false,
  outline: false,
  square: false,
};

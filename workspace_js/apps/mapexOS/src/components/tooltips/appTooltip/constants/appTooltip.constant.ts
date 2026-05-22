/** Default tooltip delay before showing in ms */
export const TOOLTIP_DEFAULT_DELAY = 300;

/** Default tooltip hide delay in ms */
export const TOOLTIP_DEFAULT_HIDE_DELAY = 0;

/** Default max width */
export const TOOLTIP_DEFAULT_MAX_WIDTH = '300px';

/** Default anchor position */
export const TOOLTIP_DEFAULT_ANCHOR = 'top middle' as const;

/** Default self position */
export const TOOLTIP_DEFAULT_SELF = 'bottom middle' as const;

/** Tooltip defaults object */
export const TOOLTIP_DEFAULTS = {
  DELAY: TOOLTIP_DEFAULT_DELAY,
  HIDE_DELAY: TOOLTIP_DEFAULT_HIDE_DELAY,
  MAX_WIDTH: TOOLTIP_DEFAULT_MAX_WIDTH,
  ANCHOR: TOOLTIP_DEFAULT_ANCHOR,
  SELF: TOOLTIP_DEFAULT_SELF,
} as const;

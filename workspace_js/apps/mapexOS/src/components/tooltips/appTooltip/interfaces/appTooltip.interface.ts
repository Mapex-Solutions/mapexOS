/**
 * Props for AppTooltip component
 */
export interface AppTooltipProps {
  /** Tooltip content (alternative to slot) */
  content?: string;

  /** Show tooltip on mobile/touch devices (default: false) */
  showOnMobile?: boolean;

  /** Delay before showing in ms (default: 300) */
  delay?: number;

  /** Hide delay in ms (default: 0) */
  hideDelay?: number;

  /** Anchor position */
  anchor?:
    | 'top left'
    | 'top middle'
    | 'top right'
    | 'center left'
    | 'center middle'
    | 'center right'
    | 'bottom left'
    | 'bottom middle'
    | 'bottom right';

  /** Self position */
  self?:
    | 'top left'
    | 'top middle'
    | 'top right'
    | 'center left'
    | 'center middle'
    | 'center right'
    | 'bottom left'
    | 'bottom middle'
    | 'bottom right';

  /** Offset [horizontal, vertical] */
  offset?: [number, number];

  /** Disable tooltip */
  disabled?: boolean;

  /** Max width */
  maxWidth?: string;
}

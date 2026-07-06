/** Visual intent of an information banner; drives its soft surface + icon color. */
export type InfoBannerVariant = 'info' | 'success' | 'warning' | 'danger' | 'neutral';

export interface InfoBannerProps {
  /** Visual intent. Defaults to `info`. */
  variant?: InfoBannerVariant;
  /** Leading icon name; overrides the variant default. */
  icon?: string;
  /** Optional bold heading rendered above the body slot. */
  title?: string;
  /** Tighter padding for compact contexts. */
  dense?: boolean;
}

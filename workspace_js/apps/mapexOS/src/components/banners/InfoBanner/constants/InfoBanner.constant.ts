import type { InfoBannerVariant } from '../interfaces';

/** Default leading icon per variant; overridable via the `icon` prop. */
export const INFO_BANNER_VARIANT_ICON: Record<InfoBannerVariant, string> = {
  info: 'info',
  success: 'check_circle',
  warning: 'warning',
  danger: 'error',
  neutral: 'info',
};

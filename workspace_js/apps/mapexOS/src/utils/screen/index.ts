import { Screen } from 'quasar'

/**
 * @description Verify if is a mobile resolution
 * @returns {boolean}
 */
export const isGreaterThanXS = () => Screen.gt.xs

/**
 * @description Verify if is a mobile resolution
 * @returns {boolean}
 */
export const isGreaterThanSM = () => Screen.gt.sm
/**
 * @description Verify if is a mobile resolution
 * @returns {boolean}
 */
export const isLessThanSM = () => Screen.lt.sm

/**
 * @description Verify if is a mobile resolution
 * @returns {boolean}
 */
export const isLessThanMD = () => Screen.lt.md

/**
 * @description Verify if is a mobile resolution
 * @returns {boolean}
 */
export const isMobile = () => Screen.lt.sm

/**
 * @description Verify if is a mobile or greater than mobile resolution
 * @returns {boolean}
 */
export const isMobileOrGreater = () => Screen.gt.sm

/**
 * @description Verify if is a tablet resolution
 * @returns {boolean}
 */
export const isTablet = () => Screen.gt.xs && Screen.lt.lg

/**
 * @description Verify if is a desktop resolution
 * @returns {boolean}
 */
export const isDesktop = () => Screen.gt.md

/**
 * @description Get a current screen name
 * @returns {string}
 */
export const currentScreenName = () => Screen.name

/**
 * @description Get a current screen width
 * @returns {string}
 */
export const width = () => Screen.width

/** Minimum recommended screen width (HD) */
export const MIN_RECOMMENDED_WIDTH = 1366;

/** Minimum recommended screen height (HD) */
export const MIN_RECOMMENDED_HEIGHT = 768;

/** LocalStorage key for resolution warning dismissal */
export const RESOLUTION_WARNING_KEY = 'resolution-warning-dismissed';

/**
 * Check if screen meets minimum HD resolution (1366x768)
 *
 * @returns {boolean} True if screen is below minimum resolution
 */
export function isBelowMinResolution(): boolean {
  return window.innerWidth < MIN_RECOMMENDED_WIDTH || window.innerHeight < MIN_RECOMMENDED_HEIGHT;
}

/**
 * Check if user has already dismissed the resolution warning
 *
 * @returns {boolean} True if warning was dismissed
 */
export function isResolutionWarningDismissed(): boolean {
  return localStorage.getItem(RESOLUTION_WARNING_KEY) === 'true';
}

/**
 * Persist resolution warning dismissal
 */
export function dismissResolutionWarning(): void {
  localStorage.setItem(RESOLUTION_WARNING_KEY, 'true');
}

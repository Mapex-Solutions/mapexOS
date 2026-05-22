/**
 * ChangePasswordPage Constants
 */

/** Minimum password length (matches ZodUserUpdateSchema min(8)) */
export const MIN_PASSWORD_LENGTH = 8;

/** Maximum password length (matches ZodUserUpdateSchema max(72)) */
export const MAX_PASSWORD_LENGTH = 72;

/** UX delay for spinner visibility (milliseconds) */
export const CHANGE_PASSWORD_DELAY_MS = 300 as const;

/** Default redirect path after successful password change */
export const DEFAULT_REDIRECT_PATH = '/home' as const;

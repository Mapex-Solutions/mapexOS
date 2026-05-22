/**
 * LoginPage Constants
 */

/**
 * Default language selection
 */
export const DEFAULT_LANGUAGE = 'en-US' as const;

/**
 * Available language options
 */
export const SUPPORTED_LANGUAGES = ['en-US', 'pt-BR'] as const;

/**
 * Language flag URLs
 */
export const LANGUAGE_FLAGS = {
  'en-US': 'https://flagcdn.com/w20/us.png',
  'pt-BR': 'https://flagcdn.com/w20/br.png',
} as const;

/**
 * Default redirect path after successful login
 */
export const DEFAULT_REDIRECT_PATH = '/home' as const;

/**
 * UX delay for spinner visibility (milliseconds)
 */
export const LOGIN_DELAY_MS = 300 as const;

/**
 * Email validation regex pattern
 */
export const EMAIL_REGEX = /.+@.+\..+/;

/**
 * Local storage key for user locale preference
 */
export const LOCALE_STORAGE_KEY = 'user-locale' as const;

/**
 * Image assets
 */
export const ASSETS = {
  ILLUSTRATION: 'iot-illustration.png',
  LOGO: 'mapex-logo.png',
} as const;

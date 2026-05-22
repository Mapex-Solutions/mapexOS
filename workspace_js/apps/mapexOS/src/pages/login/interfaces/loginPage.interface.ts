/**
 * LoginPage Interfaces
 */

/**
 * Supported language locale codes
 */
export type SupportedLocale = 'en-US' | 'pt-BR';

/**
 * Language option for dropdown selector
 */
export interface LanguageOption {
  /** Display label for the language */
  label: string;
  /** Locale code value */
  value: SupportedLocale;
  /** URL to flag icon */
  flag: string;
}

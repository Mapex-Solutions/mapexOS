import { useI18n } from 'vue-i18n';
import type { ComposerTranslation } from 'vue-i18n';

/**
 * Text transformation options for translations
 */
export interface TranslationOptions {
  /**
   * Capitalize first letter of the text
   * Example: "hello world" -> "Hello world"
   */
  capitalize?: boolean;

  /**
   * Transform all text to uppercase
   * Example: "hello world" -> "HELLO WORLD"
   */
  uppercase?: boolean;

  /**
   * Transform all text to lowercase
   * Example: "HELLO WORLD" -> "hello world"
   */
  lowercase?: boolean;

  /**
   * Capitalize first letter of each word (Title Case)
   * Example: "hello world" -> "Hello World"
   */
  titleCase?: boolean;
}

/**
 * Default translation options
 */
const defaultOptions: TranslationOptions = {
  capitalize: true,
  uppercase: false,
  lowercase: false,
  titleCase: false,
};

/**
 * Capitalizes the first letter of a string
 */
function capitalizeFirst(str: string): string {
  if (!str) return str;
  return str.charAt(0).toUpperCase() + str.slice(1);
}

/**
 * Converts string to title case (capitalizes first letter of each word)
 */
function toTitleCase(str: string): string {
  if (!str) return str;
  return str
    .split(' ')
    .map(word => capitalizeFirst(word.toLowerCase()))
    .join(' ');
}

/**
 * Applies text transformations to a string based on options
 */
function applyTransformations(text: string, options: TranslationOptions): string {
  let result = text;

  // Apply transformations in order of priority
  if (options.uppercase) {
    result = result.toUpperCase();
  } else if (options.lowercase) {
    result = result.toLowerCase();
  } else if (options.titleCase) {
    result = toTitleCase(result);
  } else if (options.capitalize) {
    result = capitalizeFirst(result);
  }

  return result;
}

/**
 * Translation utility composable with text formatting options
 *
 * @param options - Default transformation options to apply to all translations
 * @returns Translation function with applied transformations
 *
 * @example
 * ```ts
 * // Basic usage with capitalize (default)
 * const ts = useTS();
 * const text = ts('common.actions.save'); // "Save"
 *
 * // With interpolation
 * const msg = ts('common.messages.deletedSuccessfully', { item: 'User' }); // "User deleted successfully"
 *
 * // Override options per call
 * const upper = ts('common.actions.save', {}, { uppercase: true }); // "SAVE"
 *
 * // Title case
 * const title = useTS({ titleCase: true });
 * const text = title('pages.settings.title'); // "System Settings"
 *
 * // No transformation
 * const raw = useTS({ capitalize: false });
 * const text = raw('some.key'); // original text
 * ```
 */
export function useTS(defaultTransformOptions: TranslationOptions = defaultOptions) {
  const { t } = useI18n();

  /**
   * Translate a key with optional interpolation and formatting
   *
   * @param key - Translation key path
   * @param interpolation - Values to interpolate into the translation
   * @param transformOptions - Text transformation options (overrides default)
   * @returns Translated and formatted text
   */
  return (
    key: string,
    interpolation: Record<string, unknown> = {},
    transformOptions?: TranslationOptions
  ): string => {
    // Get translated text
    const translatedText = t(key, interpolation);

    // Merge options (per-call options override defaults)
    const finalOptions: TranslationOptions = {
      ...defaultTransformOptions,
      ...transformOptions,
    };

    // Apply transformations
    return applyTransformations(translatedText, finalOptions);
  };
}

/**
 * Raw translation without any transformations
 * Useful for cases where you need the exact translation from JSON
 *
 * @example
 * ```ts
 * const raw = useRawTS();
 * const text = raw('some.key'); // exact text from JSON
 * ```
 */
export function useRawTS(): ComposerTranslation {
  const { t } = useI18n();
  return t;
}

/**
 * Pre-configured translation functions for common use cases
 */
export const useTranslationPresets = () => {
  return {
    /** Standard capitalized translation (default) */
    ts: useTS({ capitalize: true }),

    /** Uppercase translation */
    tsUpper: useTS({ uppercase: true }),

    /** Lowercase translation */
    tsLower: useTS({ lowercase: true }),

    /** Title case translation */
    tsTitle: useTS({ titleCase: true }),

    /** Raw translation without formatting */
    tsRaw: useRawTS(),
  };
};

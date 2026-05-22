import { defineBoot } from '#q-app/wrappers';
import { createI18n } from 'vue-i18n';

import messages from 'src/i18n';

export type MessageLanguages = keyof typeof messages;
// Type-define 'en-US' as the master schema for the resource
export type MessageSchema = typeof messages['en-US'];

// See https://vue-i18n.intlify.dev/guide/advanced/typescript.html#global-resource-schema-type-definition
/* eslint-disable @typescript-eslint/no-empty-object-type */
declare module 'vue-i18n' {
  // define the locale messages schema
  export interface DefineLocaleMessage extends MessageSchema {}

  // define the datetime format schema
  export interface DefineDateTimeFormat {}

  // define the number format schema
  export interface DefineNumberFormat {}
}
/* eslint-enable @typescript-eslint/no-empty-object-type */

/**
 * Minimal interface for plugin translation registration and lookup.
 * Avoids complex vue-i18n generic inference issues.
 */
export interface I18nPluginBridge {
  global: {
    mergeLocaleMessage: (locale: string, messages: Record<string, unknown>) => void;
    t: (key: string) => string;
    locale: { value: string };
  };
}

/** Exported i18n instance — used by plugin registry to merge plugin translations at runtime */
export let i18nInstance: I18nPluginBridge | undefined;

export default defineBoot(({ app }) => {
  // Load saved locale from localStorage or use default
  const savedLocale = localStorage.getItem('user-locale') as MessageLanguages | null;
  const defaultLocale: MessageLanguages = 'en-US';
  const initialLocale = savedLocale && savedLocale in messages ? savedLocale : defaultLocale;

  const i18n = createI18n<{ message: MessageSchema }, MessageLanguages>({
    locale: initialLocale,
    fallbackLocale: defaultLocale,
    legacy: false,
    messages,
  });

  i18nInstance = i18n as unknown as I18nPluginBridge;

  // Set i18n instance on app
  app.use(i18n);
});

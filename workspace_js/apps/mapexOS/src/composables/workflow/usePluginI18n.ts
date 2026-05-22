import { useI18n } from 'vue-i18n';

/**
 * Scoped i18n composable for workflow plugins.
 * Wraps vue-i18n's `t()` with automatic namespace prefix `wf.{pluginId}.`.
 *
 * @param {string} pluginId - Plugin identifier (e.g., 'core-flow-control')
 * @returns {{ t: (key: string, params?: Record<string, unknown>) => string }} Scoped translation function
 */
export function usePluginI18n(pluginId: string) {
  const { t: globalT } = useI18n();

  /**
   * Translate a key scoped to this plugin's namespace
   *
   * @param {string} key - Translation key relative to plugin namespace (e.g., 'nodes.end.config.terminateWithError')
   * @param {Record<string, unknown>} [params] - Interpolation parameters
   * @returns {string} Translated string
   */
  function t(key: string, params?: Record<string, unknown>): string {
    const fullKey = `wf.${pluginId}.${key}`;
    return params ? globalT(fullKey, params) : globalT(fullKey);
  }

  return { t };
}

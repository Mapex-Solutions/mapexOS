import type { PluginRegistryState, PluginNodeType, PluginCategory } from './types';
import type { CatalogGroup } from '@src/pages/automations/workflows/createEditWorkflowPage/interfaces';
import { CATALOG_CATEGORIES } from '@src/pages/automations/workflows/createEditWorkflowPage/components/PluginCatalog/constants';
import { i18nInstance } from 'src/boot/i18n';

/**
 * Try to translate a key, returning the fallback if the key is not found.
 *
 * @param {string} key - i18n key
 * @param {string} fallback - Fallback string
 * @returns {string} Translated string or fallback
 */
function tryTranslate(key: string, fallback: string): string {
  if (!i18nInstance) return fallback;
  // Use te() to check existence first — avoids i18n console warnings for missing keys
  const global = i18nInstance.global as unknown as { te: (key: string) => boolean; t: (key: string) => string };
  if (!global.te(key)) return fallback;
  return String(global.t(key));
}

export const getters = {
  /**
   * Get all node types grouped by category for the plugin catalog.
   * Labels and descriptions are resolved from i18n with fallback to hardcoded English.
   *
   * @param {PluginRegistryState} state - Store state
   * @returns {CatalogGroup[]} Catalog groups
   */
  catalog(state: PluginRegistryState): CatalogGroup[] {
    // Touch reactive locale so the getter recomputes on locale change
    const _locale = i18nInstance?.global.locale.value;
    void _locale;

    const groupMap = new Map<PluginCategory, PluginNodeType[]>();

    for (const nodeType of state.nodeTypeMap.values()) {
      // Skip nodes hidden from the catalog (e.g., Start node)
      if (nodeType.catalogHidden) continue;

      // Find which plugin owns this node type
      let category: PluginCategory = 'custom';
      for (const plugin of state.plugins.values()) {
        if (plugin.nodeTypes.some((nt: PluginNodeType) => nt.type === nodeType.type)) {
          category = plugin.category;
          break;
        }
      }

      if (!groupMap.has(category)) {
        groupMap.set(category, []);
      }

      // Resolve translated label/description
      const pluginId = nodeType._pluginId;
      const shortName = nodeType.type.split('/').pop() || nodeType.type;

      const translatedNodeType: PluginNodeType = {
        ...nodeType,
        label: pluginId
          ? tryTranslate(`wf.${pluginId}.nodes.${shortName}.label`, nodeType.label)
          : nodeType.label,
        description: pluginId
          ? tryTranslate(`wf.${pluginId}.nodes.${shortName}.description`, nodeType.description)
          : nodeType.description,
      };

      groupMap.get(category)!.push(translatedNodeType);
    }

    // Known categories in fixed order
    const knownCategoryIds: Set<string> = new Set(CATALOG_CATEGORIES.map(c => c.category));
    const groups: CatalogGroup[] = CATALOG_CATEGORIES
      .filter(cat => groupMap.has(cat.category as PluginCategory))
      .map(cat => ({
        category: cat.category,
        label: tryTranslate(
          `pages.automations.createEditWorkflow.pluginCatalog.categories.${cat.category}`,
          cat.label,
        ),
        icon: cat.icon,
        nodeTypes: groupMap.get(cat.category as PluginCategory) || [],
      }));

    // Dynamic categories from marketplace plugins not in CATALOG_CATEGORIES
    for (const [category, nodeTypes] of groupMap) {
      if (knownCategoryIds.has(category)) continue;
      groups.push({
        category,
        label: tryTranslate(
          `pages.automations.createEditWorkflow.pluginCatalog.categories.${category}`,
          category.charAt(0).toUpperCase() + category.slice(1).replace(/_/g, ' '),
        ),
        icon: 'extension',
        nodeTypes,
      });
    }

    return groups;
  },

  /**
   * Get total number of registered node types
   *
   * @param {PluginRegistryState} state - Store state
   * @returns {number} Node type count
   */
  nodeTypeCount(state: PluginRegistryState): number {
    return state.nodeTypeMap.size;
  },

  /**
   * Get total number of registered plugins
   *
   * @param {PluginRegistryState} state - Store state
   * @returns {number} Plugin count
   */
  pluginCount(state: PluginRegistryState): number {
    return state.plugins.size;
  },
};

/** TYPE IMPORTS */
import type { PluginNodeType } from '@src/components/workflow/interfaces';

/**
 * Build default config from a node type's properties.
 * Each property's `default` value is used. Falls back to legacy `defaults` object
 * for backward compatibility with manifests that still have it.
 *
 * When `operation` is provided, only properties visible for that operation are included.
 * This is used to rebuild config from scratch when the user changes the operation.
 *
 * @param {PluginNodeType} nodeType - Node type definition
 * @param {string} [operation] - Selected operation value (filters properties by displayOptions)
 * @returns {Record<string, unknown>} Default config object
 */
export function buildDefaultConfig(nodeType: PluginNodeType, operation?: string): Record<string, unknown> {
  const config: Record<string, unknown> = {};

  if (nodeType.properties?.length) {
    for (const prop of nodeType.properties) {
      // Notices are display-only — never include in config
      if (prop.type === 'notice') continue;

      // When operation is specified, skip properties not visible for it
      if (operation && prop.displayOptions?.show?.operation) {
        const allowedOps = prop.displayOptions.show.operation as string[];
        if (!allowedOps.includes(operation)) continue;
      }

      if (prop.default !== undefined) {
        config[prop.name] = typeof prop.default === 'object' && prop.default !== null
          ? JSON.parse(JSON.stringify(prop.default))
          : prop.default;
      }
    }
  }

  // Legacy fallback: merge nodeType.defaults if properties didn't produce anything
  if (Object.keys(config).length === 0 && nodeType.defaults) {
    return { ...nodeType.defaults };
  }

  return config;
}

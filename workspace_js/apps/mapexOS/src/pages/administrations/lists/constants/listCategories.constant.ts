/**
 * List Categories Configuration
 *
 * This file defines the available categories and their corresponding list types.
 * Each category can have different list types that will be used to filter lists in the backend.
 */

export type ListCategory = 'iot';

export type IoTListType = 'manufacturers' | 'assets';

export type ListType = IoTListType;

/**
 * List types available for each category
 */
export const LIST_TYPES_BY_CATEGORY: Record<ListCategory, readonly ListType[]> = {
  iot: ['manufacturers', 'assets'] as const,
} as const;

/**
 * Get available list types for a specific category
 */
export function getListTypesByCategory(category: ListCategory | undefined): readonly ListType[] {
  if (!category) return [];
  return LIST_TYPES_BY_CATEGORY[category] || [];
}

/**
 * Check if a list type is valid for a given category
 */
export function isValidListType(category: ListCategory | undefined, type: ListType | undefined): boolean {
  if (!category || !type) return false;
  return LIST_TYPES_BY_CATEGORY[category]?.includes(type) ?? false;
}

/**
 * Get all available categories
 */
export function getAllCategories(): readonly ListCategory[] {
  return Object.keys(LIST_TYPES_BY_CATEGORY) as ListCategory[];
}

/**
 * Removes undefined, null, and empty string values from an object
 * This is useful for cleaning query parameters before sending to the API
 * to avoid sending "undefined" as a string in the URL
 *
 * @param params - The object to clean
 * @returns A new object with only defined values
 *
 * @example
 * ```ts
 * const params = {
 *   page: 1,
 *   name: undefined,
 *   enabled: false,
 *   mode: null,
 *   search: ''
 * };
 *
 * const cleaned = cleanQueryParams(params);
 * // Result: { page: 1, enabled: false }
 * ```
 */
export function cleanQueryParams<T extends Record<string, any>>(params: T): Partial<T> {
	const cleaned: Partial<T> = {};

	for (const key in params) {
		const value = params[key];

		// Keep values that are:
		// - Not undefined
		// - Not null
		// - Not empty strings
		// This means we keep: false, 0, and other falsy but valid values
		if (value !== undefined && value !== null && value !== '') {
			cleaned[key] = value;
		}
	}

	return cleaned;
}

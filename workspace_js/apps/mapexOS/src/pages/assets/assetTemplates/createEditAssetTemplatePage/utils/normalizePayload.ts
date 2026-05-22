/**
 * Normalizes a StandardizedPayload into an array of accessible field paths
 *
 * This function recursively traverses an object and generates dot-notation paths
 * for all primitive values and array elements. It follows specific rules:
 * - Arrays: Detects arrays, enters the first element (if object), does NOT use [0] notation
 * - Nested objects: Continues recursion indefinitely
 * - Primitives: Stops recursion and adds the path
 * - Returns sorted array alphabetically
 *
 * @param {any} obj - Object to normalize (typically a StandardizedPayload)
 * @param {string} prefix - Current path prefix (used internally for recursion)
 * @returns {string[]} Array of field paths (e.g., ["data.temperature", "data.location.lat"])
 *
 * @example
 * const payload = {
 *   eventType: "sensor.reading",
 *   data: {
 *     temperature: 23.5,
 *     sensors: [
 *       { id: 1, value: 10 },
 *       { id: 2, value: 20 }
 *     ]
 *   }
 * };
 *
 * const fields = normalizePayloadPaths(payload);
 * // Returns: ["data.sensors.id", "data.sensors.value", "data.temperature", "eventType"]
 */
export function normalizePayloadPaths(obj: any, prefix = ''): string[] {
	const paths: string[] = [];

	// If not an object or is null, return empty array
	if (typeof obj !== 'object' || obj === null) {
		return paths;
	}

	// Iterate over object entries
	for (const [key, value] of Object.entries(obj)) {
		const path = prefix ? `${prefix}.${key}` : key;

		if (Array.isArray(value)) {
			// ARRAYS: Enter first element (if exists and is object) and continue normalizing
			// DO NOT use [0] notation - Rule Engine will iterate the array
			if (value.length > 0 && typeof value[0] === 'object' && value[0] !== null) {
				// Array of objects - recurse into first element
				paths.push(...normalizePayloadPaths(value[0], path));
			} else {
				// Array of primitives or empty array - add path as-is
				paths.push(path);
			}
		} else if (typeof value === 'object' && value !== null) {
			// NESTED OBJECTS: Continue recursion
			paths.push(...normalizePayloadPaths(value, path));
		} else {
			// PRIMITIVES: Add the path (leaf node)
			paths.push(path);
		}
	}

	// Sort alphabetically for consistent output
	return paths.sort();
}

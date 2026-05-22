import { SchemaError } from '@mapexos/validations';

/**
 * Replaces path parameters in a URL path with their corresponding values.
 * Path parameters are identified by a colon followed by the parameter name (e.g., /:id/).
 *
 * @param path - The URL path containing parameters to be replaced (e.g., "/users/:userId/posts/:postId")
 * @param pathParams - An object containing parameter names as keys and their values to replace in the path
 * @returns The URL path with all parameters replaced by their corresponding values (URL-encoded)
 * @throws {SchemaError} If a parameter in the path doesn't have a corresponding value in pathParams
 */
export function replacePathParams(path: string, pathParams: Record<string, any>): string {
	return path.replace(/:([a-zA-Z0-9_-]+)/g, (_, key) => {
		if (!(key in pathParams)) {
			throw new SchemaError(`Missing path parameter: ${key}`);
		}
		const value = pathParams[key];
		return encodeURIComponent(value);
	});
}
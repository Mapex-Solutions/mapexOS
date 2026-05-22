/**
 * Enum to represent organization permission scope
 * - inherited: permissions apply only to this specific organization
 * - recursive: permissions apply to this organization and all descendants
 */
export enum OrganizationScopeEnum {
	INHERITED = 'inherited',
	RECURSIVE = 'recursive',
}

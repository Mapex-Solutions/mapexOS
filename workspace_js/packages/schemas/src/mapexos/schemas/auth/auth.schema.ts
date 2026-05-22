import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsBoolean, IsStringDateFormat, IsString } from '@mapexos/validations';
import { AuthTypeEnum, OrganizationTypeEnum, OrganizationScopeEnum } from '../../../common';

/**
 * Payload to log in a user
 */
export const ZodLoginSchema = z.object({
	email: StringAndNotBeEmpty,
	password: StringAndNotBeEmpty,
	keepConnected: IsBoolean.default(() => false),
});

/**
 * Response from the login endpoint
 */
export const ZodLoginResponseSchema = z.object({
	access_token: StringAndNotBeEmpty,
	refresh_token: StringAndNotBeEmpty,
	user: z.object({
		id: StringAndNotBeEmpty,
		email: StringAndNotBeEmpty,
		changePasswordNextLogin: IsBoolean,
		authProvider: z.object({
			type: z.nativeEnum(AuthTypeEnum),
			metadata: z.record(IsString, z.any()),
		}).optional(),
		firstName: StringAndNotBeEmpty,
		lastName: StringAndNotBeEmpty,
		phone: StringAndNotBeEmptyOrOptional,
		enabled: IsBoolean,
		startTour: IsBoolean,
		created: IsStringDateFormat,
		updated: IsStringDateFormat.optional(),
	}),
});

/**
 * Organization coverage item
 */
export const ZodOrganizationCoverageItemSchema = z.object({
	id: StringAndNotBeEmpty,
	name: StringAndNotBeEmpty,
	type: z.nativeEnum(OrganizationTypeEnum),
	pathKey: StringAndNotBeEmpty,
	scope: z.nativeEnum(OrganizationScopeEnum),
	membershipId: StringAndNotBeEmpty,
	roleIds: z.array(StringAndNotBeEmpty),
});

/**
 * Response from getUserCoverage endpoint
 */
export const ZodOrganizationCoverageResponseSchema = z.object({
	lastUpdated: IsStringDateFormat,
	organizations: z.array(ZodOrganizationCoverageItemSchema),
});

/**
 * Response from getMyPermissions endpoint
 */
export const ZodPermissionsResponseSchema = z.object({
	permissions: z.array(StringAndNotBeEmpty),
	version: z.number(),
});

/**
 * Build-authorization-cache request — mirror of Go
 * packages/contracts/services/mapexIam/auth/dtos.go::BuildAuthorizationCacheRequest.
 * orgId is optional; when empty, the IAM service builds root permissions.
 */
export const ZodBuildAuthorizationCacheRequestSchema = z.object({
	userId: StringAndNotBeEmpty,
	orgId: IsString.optional(),
});

/**
 * Build-coverage-cache request — mirror of Go
 * packages/contracts/services/mapexIam/auth/dtos.go::BuildCoverageCacheRequest.
 */
export const ZodBuildCoverageCacheRequestSchema = z.object({
	userId: StringAndNotBeEmpty,
});
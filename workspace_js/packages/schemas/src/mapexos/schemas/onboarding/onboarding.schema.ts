import { z, IsBoolean, IsMongoId, IsString } from '@mapexos/validations';
import { ZodAuthProviderSchema } from '../users/users.schema';
import { ZodUserResponseSchema } from '../users/users.schema';
import { ZodMembershipResponseSchema } from '../memberships/memberships.schema';

/**
 * MembershipData schema - Direct role assignment for a user
 * OrgID and Scope come from backend context:
 *   - OrgID: RequestContext.OrgContext (current org)
 *   - Scope: Organization.AccessPolicy.DefaultScope
 */
export const ZodMembershipDataSchema = z.object({
	roles: z.array(IsMongoId).min(1),
});

/**
 * ExistingGroupData schema - Join an existing group
 * User will be added to the group's members and inherit its roles.
 * The group ALREADY has a Membership (assigneeType='group') with roles defined.
 */
export const ZodExistingGroupDataSchema = z.object({
	groupId: IsMongoId,
});

/**
 * NewGroupData schema - Create a new group during onboarding
 * Creates the Group entity + Membership for the group, then adds user as member.
 */
export const ZodNewGroupDataSchema = z.object({
	name: IsString.min(3).max(150),
	description: IsString.max(500).optional(),
	roleIds: z.array(IsMongoId).min(1),
});

/**
 * GroupAccessData schema - Group-based access during user onboarding
 * EXACTLY ONE of existingGroup or newGroup must be provided.
 *
 * Use cases:
 *   - existingGroup: User joins an existing group (inherits its roles)
 *   - newGroup: Create a new group with specific roles, then user joins it
 */
export const ZodGroupAccessDataSchema = z.object({
	existingGroup: ZodExistingGroupDataSchema.optional(),
	newGroup: ZodNewGroupDataSchema.optional(),
}).refine(
	(data) => (data.existingGroup && !data.newGroup) || (!data.existingGroup && data.newGroup),
	{ message: 'Exactly one of existingGroup or newGroup must be provided' }
);

/**
 * CreateUserWithMemberships schema - Orchestration DTO for creating user + memberships
 * Requires at least one membership OR one group
 */
export const ZodCreateUserWithMembershipsSchema = z.object({
	// User data
	email: IsString.email().max(254),
	password: IsString.min(8).max(72).optional(),
	changePasswordNextLogin: IsBoolean,
	authProvider: ZodAuthProviderSchema.optional(), // V1: always internal, backend defaults
	firstName: IsString.min(2).max(100),
	lastName: IsString.min(2).max(100),
	phone: IsString.optional(),
	jobTitle: IsString.max(120).optional(),
	enabled: IsBoolean,
	avatar: IsString.url().optional(),

	// Direct role memberships (optional if groups provided)
	memberships: z.array(ZodMembershipDataSchema).optional(),

	// Group access (optional if memberships provided)
	// Each entry represents either joining an existing group OR creating a new group
	groups: z.array(ZodGroupAccessDataSchema).optional(),
}).refine(
	(data) => (data.memberships && data.memberships.length > 0) || (data.groups && data.groups.length > 0),
	{ message: 'At least one membership (direct roles) or group must be provided' }
);

/**
 * UpdateUserWithAccess path params schema
 * Used for PATCH /api/v1/onboarding/users/:userId
 */
export const ZodUpdateUserWithAccessParamsSchema = z.object({
	userId: IsMongoId,
});

/**
 * UpdateUserWithAccess schema - Orchestration DTO for updating user + access
 * Updates user data and replaces access configuration (memberships/groups)
 * OrgID and Scope come from backend context
 */
export const ZodUpdateUserWithAccessSchema = z.object({
	// User data (all optional for partial updates)
	firstName: IsString.min(2).max(100).optional(),
	lastName: IsString.min(2).max(100).optional(),
	phone: IsString.optional(),
	jobTitle: IsString.max(120).optional(),
	enabled: IsBoolean.optional(),
	avatar: IsString.url().optional(),
	password: IsString.min(8).max(72).optional(),
	changePasswordNextLogin: IsBoolean.optional(),

	// Direct role memberships (optional if groups provided)
	memberships: z.array(ZodMembershipDataSchema).optional(),

	// Group access (optional if memberships provided)
	// Each entry represents either joining an existing group OR creating a new group
	groups: z.array(ZodGroupAccessDataSchema).optional(),
}).refine(
	(data) => (data.memberships && data.memberships.length > 0) || (data.groups && data.groups.length > 0),
	{ message: 'At least one membership (direct roles) or group must be provided' }
);

/**
 * UserOnboardingResponse schema - Response from orchestration endpoint
 */
export const ZodUserOnboardingResponseSchema = z.object({
	user: ZodUserResponseSchema,
	memberships: z.array(ZodMembershipResponseSchema),
});

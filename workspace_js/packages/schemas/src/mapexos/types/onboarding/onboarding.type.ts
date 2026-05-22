import { z } from 'zod';
import {
	ZodMembershipDataSchema,
	ZodExistingGroupDataSchema,
	ZodNewGroupDataSchema,
	ZodGroupAccessDataSchema,
	ZodCreateUserWithMembershipsSchema,
	ZodUpdateUserWithAccessParamsSchema,
	ZodUpdateUserWithAccessSchema,
	ZodUserOnboardingResponseSchema,
} from '@/mapexos';

/**
 * Membership data for direct role assignment
 */
export type MembershipData = z.infer<typeof ZodMembershipDataSchema>;

/**
 * Existing group data - join an existing group (inherits its roles)
 */
export type ExistingGroupData = z.infer<typeof ZodExistingGroupDataSchema>;

/**
 * New group data - create a new group during onboarding
 */
export type NewGroupData = z.infer<typeof ZodNewGroupDataSchema>;

/**
 * Group access data - either existing or new group
 */
export type GroupAccessData = z.infer<typeof ZodGroupAccessDataSchema>;

/**
 * Create user with memberships request payload
 */
export type CreateUserWithMemberships = z.infer<typeof ZodCreateUserWithMembershipsSchema>;

/**
 * Update user with access path params
 */
export type UpdateUserWithAccessParams = z.infer<typeof ZodUpdateUserWithAccessParamsSchema>;

/**
 * Update user with access request payload
 */
export type UpdateUserWithAccess = z.infer<typeof ZodUpdateUserWithAccessSchema>;

/**
 * User onboarding response (user + created memberships)
 */
export type UserOnboardingResponse = z.infer<typeof ZodUserOnboardingResponseSchema>;

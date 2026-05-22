import { z } from 'zod';
import {
	ZodAuthProviderSchema,
	ZodAuthProviderUpdateSchema,
	ZodUserIdSchema,
	ZodUserCreateSchema,
	ZodUserUpdateSchema,
	ZodUserQuerySchema,
	ZodUserResponseSchema,
	ZodUserGroupInfoSchema,
	ZodUserMembershipInfoSchema,
} from '@/mapexos';

/**
 * Auth Provider types
 */
export type AuthProvider = z.infer<typeof ZodAuthProviderSchema>;
export type AuthProviderUpdate = z.infer<typeof ZodAuthProviderUpdateSchema>;

/**
 * User API types
 */
export type UserId = z.infer<typeof ZodUserIdSchema>;
export type UserCreate = z.infer<typeof ZodUserCreateSchema>;
export type UserUpdate = z.infer<typeof ZodUserUpdateSchema>;
export type UserQuery = z.infer<typeof ZodUserQuerySchema>;
export type UserResponse = z.infer<typeof ZodUserResponseSchema>;

/**
 * Enriched user data types (for detail view)
 */
export type UserGroupInfo = z.infer<typeof ZodUserGroupInfoSchema>;
export type UserMembershipInfo = z.infer<typeof ZodUserMembershipInfoSchema>;

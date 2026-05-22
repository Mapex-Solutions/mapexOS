import { z } from 'zod';
import {
	ZodMembershipIdSchema,
	ZodMembershipCreateSchema,
	ZodMembershipUpdateSchema,
	ZodMembershipQuerySchema,
	ZodMembershipResponseSchema,
} from '@/mapexos';

/**
 * Membership API types
 */
export type MembershipId = z.infer<typeof ZodMembershipIdSchema>;
export type MembershipCreate = z.infer<typeof ZodMembershipCreateSchema>;
export type MembershipUpdate = z.infer<typeof ZodMembershipUpdateSchema>;
export type MembershipQuery = z.infer<typeof ZodMembershipQuerySchema>;
export type MembershipResponse = z.infer<typeof ZodMembershipResponseSchema>;

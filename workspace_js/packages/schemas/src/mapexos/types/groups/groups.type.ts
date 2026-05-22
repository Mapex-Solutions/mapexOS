import { z } from 'zod';
import {
	ZodGroupIdSchema,
	ZodGroupCreateSchema,
	ZodGroupUpdateSchema,
	ZodGroupQuerySchema,
	ZodGroupResponseSchema,
	ZodGroupMembersQuerySchema,
	ZodGroupMemberResponseSchema,
	ZodGroupMemberAddSchema,
	ZodGroupMemberIdSchema,
} from '@/mapexos';

/**
 * Group API types
 */
export type GroupId = z.infer<typeof ZodGroupIdSchema>;
export type GroupCreate = z.infer<typeof ZodGroupCreateSchema>;
export type GroupUpdate = z.infer<typeof ZodGroupUpdateSchema>;
export type GroupQuery = z.infer<typeof ZodGroupQuerySchema>;
export type GroupResponse = z.infer<typeof ZodGroupResponseSchema>;
export type GroupMembersQuery = z.infer<typeof ZodGroupMembersQuerySchema>;
export type GroupMemberResponse = z.infer<typeof ZodGroupMemberResponseSchema>;
export type GroupMemberAdd = z.infer<typeof ZodGroupMemberAddSchema>;
export type GroupMemberId = z.infer<typeof ZodGroupMemberIdSchema>;

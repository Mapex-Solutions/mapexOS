import { z } from 'zod';
import {
	ZodOrganizationIdSchema,
	ZodAddressSchema,
	ZodAuthConfigSchema,
	ZodAccessPolicySchema,
	ZodOrganizationCreateSchema,
	ZodOrganizationUpdateSchema,
	ZodOrganizationQuerySchema,
	ZodOrganizationResponseSchema,
	ZodTreeQuerySchema,
	ZodTreeItemSchema,
	ZodTreeResponseSchema,
} from '@/mapexos';

/**
 * Organization nested types
 */
export type Address = z.infer<typeof ZodAddressSchema>;
export type AuthConfig = z.infer<typeof ZodAuthConfigSchema>;
export type AccessPolicy = z.infer<typeof ZodAccessPolicySchema>;

/**
 * Organization API types
 */
export type OrganizationId = z.infer<typeof ZodOrganizationIdSchema>;
export type OrganizationCreate = z.infer<typeof ZodOrganizationCreateSchema>;
export type OrganizationUpdate = z.infer<typeof ZodOrganizationUpdateSchema>;
export type OrganizationQuery = z.infer<typeof ZodOrganizationQuerySchema>;
export type OrganizationResponse = z.infer<typeof ZodOrganizationResponseSchema>;

/**
 * Tree navigation types
 */
export type TreeQuery = z.infer<typeof ZodTreeQuerySchema>;
export type TreeItem = z.infer<typeof ZodTreeItemSchema>;
export type TreeResponse = z.infer<typeof ZodTreeResponseSchema>;

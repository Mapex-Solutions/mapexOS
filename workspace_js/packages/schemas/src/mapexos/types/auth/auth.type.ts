import { z } from 'zod';
import {
	ZodLoginSchema,
	ZodLoginResponseSchema,
	ZodOrganizationCoverageItemSchema,
	ZodOrganizationCoverageResponseSchema,
	ZodPermissionsResponseSchema
} from '@/mapexos';

export type Login = z.infer<typeof ZodLoginSchema>
export type LoginResponse = z.infer<typeof ZodLoginResponseSchema>
export type OrganizationCoverageItem = z.infer<typeof ZodOrganizationCoverageItemSchema>
export type OrganizationCoverageResponse = z.infer<typeof ZodOrganizationCoverageResponseSchema>
export type PermissionsResponse = z.infer<typeof ZodPermissionsResponseSchema>

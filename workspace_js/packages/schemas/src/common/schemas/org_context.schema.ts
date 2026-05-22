import { z, StringAndNotBeEmpty } from '@mapexos/validations';

/**
 * Schema for organization context body parameter.
 * Used by API endpoints that require orgId to be passed via the X-Org-Context header.
 */
export const ZodOrgContextSchema = z.object({
	orgId: StringAndNotBeEmpty,
});

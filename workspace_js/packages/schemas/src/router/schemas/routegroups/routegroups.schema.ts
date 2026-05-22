import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsBoolean, IsString, IsNumber, IsStringDateFormat, IsMongoId } from '@mapexos/validations';

/**
 * Nested router type schemas
 */
const ZodLakeHouseDataSchema = z.object({
	lakeHouseId: IsMongoId,
	metadata: z.record(IsString, z.any()).optional(),
});

const ZodNotificationDataSchema = z.object({
	notificationId: IsMongoId,
	metadata: z.record(IsString, z.any()).optional(),
});

const ZodTriggerDataSchema = z.object({
	triggerId: IsMongoId,
	metadata: z.record(IsString, z.any()).optional(),
});

const ZodSaveEventDataSchema = z.object({
	metadata: z.record(IsString, z.any()).optional(),
});

const ZodWorkflowDataSchema = z.object({
	mode: z.enum(['newInstance', 'signal', 'signalOrStart']),
	data: z.record(IsString, z.any()),
	metadata: z.record(IsString, z.any()).optional(),
});

/**
 * Router schema - represents routing configuration
 */
const ZodRouterSchema = z.object({
	kind: z.enum(['lake_house', 'notification', 'trigger', 'save_event', 'workflow']),
	lakeHouse: ZodLakeHouseDataSchema.optional(),
	notification: ZodNotificationDataSchema.optional(),
	trigger: ZodTriggerDataSchema.optional(),
	saveEvent: ZodSaveEventDataSchema.optional(),
	workflow: ZodWorkflowDataSchema.optional(),
}).refine((data) => {
	// Validate that corresponding field is provided based on kind
	switch (data.kind) {
		case 'lake_house':
			return !!data.lakeHouse;
		case 'notification':
			return !!data.notification;
		case 'trigger':
			return !!data.trigger;
		case 'save_event':
			return true;
		case 'workflow':
			return !!data.workflow;
		default:
			return false;
	}
}, {
	message: "Router configuration must match the specified kind",
	path: ['kind'],
});

/**
 * RouteGroup ID parameter schema (for URL params - MongoDB ObjectID)
 */
export const ZodRouteGroupIdSchema = z.object({
	routeGroupId: IsMongoId,
});

/**
 * RouteGroup Create schema - Used for creating new route groups
 */
export const ZodRouteGroupCreateSchema = z.object({
	version: IsString.regex(/^\d+\.\d+\.\d+$/, 'Version must be in semver format (e.g., 1.0.0)'),
	name: StringAndNotBeEmpty,
	description: StringAndNotBeEmptyOrOptional,
	enabled: IsBoolean,
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),
	routers: z.array(ZodRouterSchema).min(1).optional(),
});

/**
 * RouteGroup Update schema - Used for updating existing route groups (all fields optional)
 */
export const ZodRouteGroupUpdateSchema = z.object({
	version: IsString.regex(/^\d+\.\d+\.\d+$/, 'Version must be in semver format (e.g., 1.0.0)').optional(),
	name: StringAndNotBeEmpty.optional(),
	description: StringAndNotBeEmptyOrOptional,
	enabled: IsBoolean.optional(),
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),
	routers: z.array(ZodRouterSchema).min(1).optional(),
});

/**
 * RouteGroup Response schema - Returned from API
 */
export const ZodRouteGroupResponseSchema = z.object({
	id: IsMongoId.optional(),
	version: IsString.optional(),
	name: IsString.optional(),
	description: IsString.optional(),
	enabled: IsBoolean.optional(),

	// Multi-tenant fields
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),
	orgId: IsString.optional(),
	pathKey: IsString.optional(),

	routers: z.array(ZodRouterSchema).optional(),

	created: IsStringDateFormat.optional(),
	updated: IsStringDateFormat.optional(),
});

/**
 * RouteGroup Query schema - Used for filtering/pagination
 */
export const ZodRouteGroupQuerySchema = z.object({
	// BaseQueryDTO fields
	projection: IsString.optional(),
	page: IsNumber.int().min(1).optional(),
	perPage: IsNumber.int().min(1).max(100).optional(),
	sort: IsString.optional(),
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	name: IsString.max(100).optional(),
	enabled: IsBoolean.optional(),
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),
	version: IsString.regex(/^\d+\.\d+\.\d+$/, 'Version must be in semver format (e.g., 1.0.0)').optional(),
	/**
	 * kinds optionally restricts the result set to RouteGroups whose every
	 * router.kind is contained in the provided set. Strict semantic.
	 * Mirrors Go: packages/contracts/services/router/routegroups/dto.go::RouteGroupQuery.Kinds
	 * Used by the asset wizard's Health step (HealthMonitoringSection) to surface
	 * only RouteGroups acceptable to validateHealthMonitorConfig.
	 */
	kinds: z.array(z.enum(['lake_house', 'notification', 'trigger', 'save_event', 'workflow'])).optional(),
});

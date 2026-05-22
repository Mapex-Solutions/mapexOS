import { z } from 'zod';
import {
	StringAndNotBeEmpty,
	IsString,
	IsBoolean,
	IsStringDateFormat,
	NumberIntAndPositive,
} from '@mapexos/validations';

/**
 * Zod schemas for the asset connectivity history API.
 *
 * Backed by the Go contract at
 * workspace_go/packages/contracts/services/events/asset_status/dto.go.
 * Zod types are inferred (never hand-rolled) and re-exported alongside
 * the schemas.
 */

/** Single row of asset_status_history (response shape). */
export const ZodAssetConnectivityEventSchema = z.object({
	created: IsStringDateFormat,
	orgId: StringAndNotBeEmpty,
	pathKey: IsString,
	assetUUID: StringAndNotBeEmpty,
	assetName: IsString.optional().default(''),
	eventId: StringAndNotBeEmpty,
	eventType: z.enum(['offline', 'online']),
	lastSeenAt: IsStringDateFormat.optional(),
	thresholdMinutes: z.number().int().min(0).optional(),
	missCount: z.number().int().min(0).optional(),
});
export type AssetConnectivityEvent = z.infer<typeof ZodAssetConnectivityEventSchema>;

/**
 * Query parameters for
 * GET /api/v1/events/assets/:assetUUID/connectivity_history.
 * Cursor fields mirror CursorQueryDTO; from/to/eventType are optional filters.
 */
export const ZodAssetConnectivityQuerySchema = z.object({
	cursor: IsString.optional(),
	limit: NumberIntAndPositive.optional(),
	direction: z.enum(['next', 'prev']).optional(),
	sortAsc: IsBoolean.optional(),
	includeChildren: IsBoolean.optional(),
	from: IsStringDateFormat.optional(),
	to: IsStringDateFormat.optional(),
	eventType: z.enum(['offline', 'online']).optional(),
	assetUUID: IsString.optional(),
});
export type AssetConnectivityQuery = z.infer<typeof ZodAssetConnectivityQuerySchema>;

/** Cursor-paginated response wrapper — mirrors EventsRawCursorResult shape. */
export const ZodAssetConnectivityCursorResultSchema = z.object({
	items: z.array(ZodAssetConnectivityEventSchema),
	nextCursor: IsStringDateFormat.optional(),
	prevCursor: IsStringDateFormat.optional(),
	hasNext: IsBoolean,
	hasPrevious: IsBoolean,
});
export type AssetConnectivityCursorResult = z.infer<typeof ZodAssetConnectivityCursorResultSchema>;

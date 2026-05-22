import { z, StringAndNotBeEmpty, IsBoolean } from '@mapexos/validations';

/**
 * Cross-service auth projection schema. Mirrors Go:
 * packages/contracts/services/assets/auth/dto.go::AuthProjection.
 *
 * Stored at MinIO bucket `mapex-asset-auth` under key `{assetUUID}.json`
 * (flat layout — assetUUID is globally unique via Mongo
 * idx_asset_uuid_unique). The broker plugin reads it on every CONNECT
 * lookup. Also returned by the assets-service internal endpoint
 * GET /internal/asset-auth/:assetUUID as the broker's L3 fallback.
 *
 * `type` is intentionally open so future auth surfaces (http_api,
 * lorawan_key, etc.) can reuse the same projection shape without
 * breaking the broker. Today only `mqtt` is valid.
 */
export const ZodAuthProjectionSchema = z.object({
	assetUUID: StringAndNotBeEmpty,
	orgId: StringAndNotBeEmpty,
	enabled: IsBoolean,
	type: z.enum(['mqtt']),
	authType: z.enum(['password', 'cert']),
	passwordHash: z.string().optional(),
	currentCertSerial: z.string().optional(),
});

import { z, StringAndNotBeEmpty, IsBoolean } from '@mapexos/validations';

/**
 * Cross-service auth projection schema. Mirrors Go:
 * packages/contracts/services/assets/auth/dto.go::AuthProjection.
 *
 * Stored at MinIO bucket `mapex-asset-auth` under key `{assetUUID}.json`
 * (flat layout — assetUUID is globally unique via Mongo
 * idx_asset_uuid_unique). The broker plugin reads it on every CONNECT
 * lookup. Also returned by the assets-service internal endpoint
 * GET /internal/asset_auth/:assetUUID as the broker's L3 fallback.
 *
 * `type` discriminates the auth surface and both blocks are nested and
 * symmetric: the `mqtt` block is set when type=mqtt; the `lorawan` block is set
 * when type=lorawan. The broker reads the mqtt block; the LNS reads the lorawan
 * block.
 */

// EncryptedKeys: the four envelope fields. Go []byte marshals to base64 JSON
// strings, so these are plain strings on the wire.
const ZodEncryptedKeysSchema = z.object({
	encryptedDek: z.string(),
	dekNonce: z.string(),
	encryptedKey: z.string(),
	keyNonce: z.string(),
});

// MqttAuth: the broker's CONNECT-decision fields. authType selects the mode,
// passwordHash backs password mode, currentCertSerial backs cert mode.
const ZodMqttAuthSchema = z.object({
	authType: z.enum(['password', 'cert']).optional(),
	passwordHash: z.string().optional(),
	currentCertSerial: z.string().optional(),
});

// LorawanGatewayAuth: per-gateway frequency plan + auth mode + cert serial or
// api-key hash (no device keys). The LNS Gateway Server builds the radio config +
// auth from it; apiKeyHash is set in key mode and bcrypt-compared at connect.
const ZodLorawanGatewayAuthSchema = z.object({
	authMode: z.enum(['eui', 'cert', 'key']),
	frequencyPlanId: StringAndNotBeEmpty,
	frequencyPlanIds: z.array(z.string()).optional(),
	latitude: z.number().optional(),
	longitude: z.number().optional(),
	altitude: z.number().optional(),
	currentCertSerial: z.string().optional(),
	apiKeyHash: z.string().optional(),
});

// LorawanAuth: kind discriminates a device (identity + profile + encrypted keys,
// read by the LNS NS/JS) from a gateway (frequency plan + auth, read by the LNS
// Gateway Server). Device fields are optional because a gateway omits them.
const ZodLorawanAuthSchema = z.object({
	kind: z.enum(['device', 'gateway']).optional(),
	devEui: z.string().optional(),
	joinEui: z.string().optional(),
	region: z.string().optional(),
	class: z.enum(['A', 'B', 'C']).optional(),
	macVersion: z.string().optional(),
	phyVersion: z.string().optional(),
	activation: z.enum(['otaa', 'abp']).optional(),
	keys: ZodEncryptedKeysSchema.optional(),
	gateway: ZodLorawanGatewayAuthSchema.optional(),
});

export const ZodAuthProjectionSchema = z.object({
	assetUUID: StringAndNotBeEmpty,
	orgId: StringAndNotBeEmpty,
	enabled: IsBoolean,
	type: z.enum(['mqtt', 'lorawan']),
	mqtt: ZodMqttAuthSchema.optional(),
	lorawan: ZodLorawanAuthSchema.optional(),
});

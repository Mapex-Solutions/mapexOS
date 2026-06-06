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

// LorawanAuth: identity + profile in clear, key material encrypted in keys.
const ZodLorawanAuthSchema = z.object({
	devEui: StringAndNotBeEmpty,
	joinEui: z.string().optional(),
	region: StringAndNotBeEmpty,
	class: z.enum(['A', 'B', 'C']),
	macVersion: StringAndNotBeEmpty,
	phyVersion: StringAndNotBeEmpty,
	activation: z.enum(['otaa', 'abp']),
	keys: ZodEncryptedKeysSchema,
});

export const ZodAuthProjectionSchema = z.object({
	assetUUID: StringAndNotBeEmpty,
	orgId: StringAndNotBeEmpty,
	enabled: IsBoolean,
	type: z.enum(['mqtt', 'lorawan']),
	mqtt: ZodMqttAuthSchema.optional(),
	lorawan: ZodLorawanAuthSchema.optional(),
});

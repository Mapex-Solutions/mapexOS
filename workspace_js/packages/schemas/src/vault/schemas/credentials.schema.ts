import { z, IsString, IsBoolean, IsNumber, IsMongoId, IsStringDateFormat, StringAndNotBeEmpty } from '@mapexos/validations';

/**
 * TokenRequestConfig — fully describes one HTTP request for token
 * acquisition or renewal. Mirror of Go
 * packages/contracts/services/mapexVault/credentials/dto.go::TokenRequestConfig.
 */
export const ZodTokenRequestConfigSchema = z.object({
  method: IsString.optional(),
  url: StringAndNotBeEmpty,
  contentType: IsString.optional(),
  headers: z.record(IsString, IsString).optional(),
  body: z.record(IsString, z.any()).optional(),
  queryParams: z.record(IsString, IsString).optional(),
  accessTokenPath: IsString.optional(),
  refreshTokenPath: IsString.optional(),
  expiresInPath: IsString.optional(),
});

/**
 * ProviderConfig — login + refresh request templates. Stored unencrypted
 * because the refresh consumer needs to read these fields to know
 * WHERE/HOW to refresh, without decrypting secrets.
 */
export const ZodProviderConfigSchema = z.object({
  loginConfig: ZodTokenRequestConfigSchema.optional(),
  refreshConfig: ZodTokenRequestConfigSchema.optional(),
});

/**
 * Credential Create schema — mirror of Go CreateCredentialDTO.
 */
export const ZodCredentialCreateSchema = z.object({
  name: StringAndNotBeEmpty,
  type: z.enum(['manual', 'oauth2', 'userAndPass']),
  pluginId: StringAndNotBeEmpty,
  credentialDefId: IsString.optional().default(''),
  data: z.record(IsString, z.any()),
  providerConfig: ZodProviderConfigSchema.optional(),
  isTemplate: IsBoolean.optional().default(false),
});

/**
 * Credential Update schema — mirror of Go UpdateCredentialDTO.
 */
export const ZodCredentialUpdateSchema = z.object({
  name: IsString.optional(),
  data: z.record(IsString, z.any()).optional(),
  providerConfig: ZodProviderConfigSchema.optional(),
  isTemplate: IsBoolean.optional(),
});

/**
 * Credential Query schema — mirror of Go CredentialQueryDTO.
 */
export const ZodCredentialQuerySchema = z.object({
  pluginId: IsString.optional(),
  type: z.enum(['manual', 'oauth2', 'userAndPass']).optional(),
  status: z.enum(['active', 'expired', 'revoked', 'error']).optional(),
  page: IsNumber.int().min(1).optional(),
  perPage: IsNumber.int().min(1).max(100).optional(),
});

/**
 * Credential Response schema — mirror of Go CredentialResponse.
 * Encrypted blobs are deliberately absent; this is the public projection.
 */
export const ZodCredentialResponseSchema = z.object({
  id: IsMongoId.optional(),
  name: IsString,
  type: z.enum(['manual', 'oauth2', 'userAndPass']),
  pluginId: IsString,
  credentialDefId: IsString.optional(),
  orgId: IsMongoId.optional(),
  pathKey: IsString.optional(),
  isTemplate: IsBoolean.optional(),
  status: z.enum(['active', 'expired', 'revoked', 'error']).optional(),
  tokenExpiresAt: IsStringDateFormat.optional(),
  lastRefreshedAt: IsStringDateFormat.optional(),
  refreshError: IsString.optional(),
  providerConfig: ZodProviderConfigSchema.optional(),
  created: IsStringDateFormat.optional(),
  updated: IsStringDateFormat.optional(),
});

/**
 * OAuth2 callback input — mirror of Go OAuthCallbackDTO.
 */
export const ZodOAuthCallbackSchema = z.object({
  provider: StringAndNotBeEmpty,
  code: StringAndNotBeEmpty,
  redirectUri: StringAndNotBeEmpty,
  pluginId: StringAndNotBeEmpty,
  credentialDefId: IsString.optional(),
});

/**
 * Connection Create schema — mirror of Go CreateConnectionDTO.
 */
export const ZodConnectionCreateSchema = z.object({
  provider: StringAndNotBeEmpty,
  accountId: StringAndNotBeEmpty,
  accountName: IsString.optional(),
  credentialId: IsMongoId,
  scopes: z.array(IsString).optional(),
});

/**
 * Connection Upsert schema — mirror of Go UpsertConnectionDTO. The upsert
 * key is (provider, accountId, orgId).
 */
export const ZodConnectionUpsertSchema = z.object({
  provider: StringAndNotBeEmpty,
  accountId: StringAndNotBeEmpty,
  accountName: IsString.optional(),
  credentialId: IsMongoId,
  scopes: z.array(IsString).optional(),
});

/**
 * Connection Query schema — mirror of Go ConnectionQueryDTO.
 */
export const ZodConnectionQuerySchema = z.object({
  provider: IsString.optional(),
  status: z.enum(['active', 'revoked', 'expired']).optional(),
  page: IsNumber.int().min(1).optional(),
  perPage: IsNumber.int().min(1).max(100).optional(),
});

/**
 * Connection Response schema — mirror of Go ConnectionResponse.
 */
export const ZodConnectionResponseSchema = z.object({
  id: IsMongoId.optional(),
  provider: IsString,
  accountId: IsString,
  accountName: IsString.optional(),
  status: z.enum(['active', 'revoked', 'expired']),
  credentialId: IsMongoId,
  userId: IsMongoId.optional(),
  orgId: IsMongoId.optional(),
  pathKey: IsString.optional(),
  scopes: z.array(IsString).optional(),
  connectedAt: IsStringDateFormat.optional(),
  lastUsedAt: IsStringDateFormat.optional(),
  created: IsStringDateFormat.optional(),
  updated: IsStringDateFormat.optional(),
});

import { z } from 'zod';
import {
  ZodTokenRequestConfigSchema,
  ZodProviderConfigSchema,
  ZodCredentialCreateSchema,
  ZodCredentialUpdateSchema,
  ZodCredentialQuerySchema,
  ZodCredentialResponseSchema,
  ZodOAuthCallbackSchema,
  ZodConnectionCreateSchema,
  ZodConnectionUpsertSchema,
  ZodConnectionQuerySchema,
  ZodConnectionResponseSchema,
} from '../schemas/credentials.schema';

export type TokenRequestConfig = z.infer<typeof ZodTokenRequestConfigSchema>;
export type ProviderConfig = z.infer<typeof ZodProviderConfigSchema>;
export type CredentialCreate = z.infer<typeof ZodCredentialCreateSchema>;
export type CredentialUpdate = z.infer<typeof ZodCredentialUpdateSchema>;
export type CredentialQuery = z.infer<typeof ZodCredentialQuerySchema>;
export type CredentialResponse = z.infer<typeof ZodCredentialResponseSchema>;

export type OAuthCallback = z.infer<typeof ZodOAuthCallbackSchema>;
export type ConnectionCreate = z.infer<typeof ZodConnectionCreateSchema>;
export type ConnectionUpsert = z.infer<typeof ZodConnectionUpsertSchema>;
export type ConnectionQuery = z.infer<typeof ZodConnectionQuerySchema>;
export type ConnectionResponse = z.infer<typeof ZodConnectionResponseSchema>;

export type CredentialType = 'manual' | 'oauth2' | 'userAndPass';
export type CredentialStatus = 'active' | 'expired' | 'revoked' | 'error';
export type ConnectionStatus = 'active' | 'revoked' | 'expired';

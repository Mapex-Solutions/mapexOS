import type { GetToken } from '@src/common';
import type { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';

/**
 * Vault credential types (inline — avoids circular dependency with schemas during build)
 */
export interface VaultCredentialCreate {
  name: string;
  type: 'manual' | 'oauth2' | 'userAndPass';
  pluginId: string;
  credentialDefId?: string;
  data: Record<string, unknown>;
  providerConfig?: Record<string, unknown>;
  isTemplate?: boolean;
}

export interface VaultCredentialUpdate {
  name?: string;
  data?: Record<string, unknown>;
  providerConfig?: Record<string, unknown>;
  isTemplate?: boolean;
}

export interface VaultCredentialQuery {
  pluginId?: string;
  type?: string;
  status?: string;
  page?: number;
  perPage?: number;
}

export interface VaultCredentialResponse {
  id: string;
  name: string;
  type: string;
  pluginId: string;
  status: string;
  [key: string]: unknown;
}

export interface PaginatedResponse<T> {
  items: T[];
  pagination: { page: number; perPage: number; totalPages: number; totalItems: number };
}

/**
 * Creates Vault Credential user API (JWT auth — called from frontend).
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
  const factory = createApiFactory(http);

  return factory({
    basePath: '/api/v1/credentials',
    useAuthJWT: true,
    getToken,
    methods: {
      list: {
        method: 'GET',
        path: '',
        queryParams: {} as VaultCredentialQuery,
        responseType: {} as PaginatedResponse<VaultCredentialResponse>,
      },
      create: {
        method: 'POST',
        path: '',
        bodyParams: {} as VaultCredentialCreate,
        responseType: {} as VaultCredentialResponse,
      },
      getById: {
        method: 'GET',
        path: '/:credentialId',
        pathParams: {} as { credentialId: string },
        responseType: {} as VaultCredentialResponse,
      },
      update: {
        method: 'PATCH',
        path: '/:credentialId',
        pathParams: {} as { credentialId: string },
        bodyParams: {} as VaultCredentialUpdate,
        responseType: {} as VaultCredentialResponse,
      },
      delete: {
        method: 'DELETE',
        path: '/:credentialId',
        pathParams: {} as { credentialId: string },
        responseType: {} as { deleted: boolean },
      },
      test: {
        method: 'POST',
        path: '/:credentialId/test',
        pathParams: {} as { credentialId: string },
        responseType: {} as { success: boolean },
      },
    },
  });
}

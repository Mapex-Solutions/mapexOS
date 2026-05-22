import { isEmpty } from 'lodash';
import type { ApiConfig, ApiInterceptors, SessionsConfig, GetToken } from '@src/common';
import { createHttp } from '@src/tools';
import { createVaultCredentialApi } from './credential';

/**
 * Creates Vault Service API with credential module.
 */
export function createVaultApi(config: ApiConfig, sessionsConfig?: SessionsConfig) {
  const interceptors = !isEmpty(sessionsConfig?.interceptors)
    ? sessionsConfig?.interceptors
    : config?.interceptors || {} as ApiInterceptors;

  const getToken: GetToken | undefined = sessionsConfig?.getToken;
  const http = createHttp(config, interceptors);

  const credential = createVaultCredentialApi(http, getToken);

  return {
    http,
    credential,
  };
}

export type VaultApiMethods = ReturnType<typeof createVaultApi>;

import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';
import { userApi } from './user.api';

export function createVaultCredentialApi(http: AxiosInstance, getToken: GetToken | undefined) {
  return {
    ...userApi(http, getToken),
  };
}

export type VaultCredentialApiMethods = ReturnType<typeof createVaultCredentialApi>;

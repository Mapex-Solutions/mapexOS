import type { ApiWrappers } from '@mapexos/apis';
import { initializeApis } from '@mapexos/apis';
import { useAuthStore } from '@stores/auth';
import { useOrganizationStore } from '@stores/organization';
import type { AxiosError, InternalAxiosRequestConfig } from 'axios';
import { useLogger } from '@composables/useLogger';

const logger = useLogger('MapexAPI');

/**
 * Token refresh queue mechanism
 *
 * Problem: When multiple requests fail with 401 simultaneously, we need to prevent
 * multiple refresh token calls to the backend.
 *
 * Solution: Queue pattern
 * - First 401: Sets isRefreshing = true, calls refresh endpoint
 * - Subsequent 401s: Add their resolve/reject to failedQueue and wait
 * - When refresh completes: processQueue() wakes all waiting requests with new token
 *
 * Example flow:
 * 1. Request A (401) → isRefreshing = true → Calls /auth/refresh
 * 2. Request B (401) → isRefreshing = true → Queued (waits)
 * 3. Request C (401) → isRefreshing = true → Queued (waits)
 * 4. Refresh succeeds → processQueue(null, newToken) → B and C wake up with token
 * 5. B and C retry with new token
 */
let isRefreshing = false;
let failedQueue: Array<{
  resolve: (token: string) => void;
  reject: (error: Error) => void;
}> = [];

/**
 * Process all queued requests after token refresh completes
 *
 * Wakes up all requests waiting in the queue by calling their resolve/reject
 *
 * @param error - If refresh failed, reject all queued requests with this error
 * @param token - If refresh succeeded, resolve all queued requests with new token
 */
function processQueue(error: Error | null, token: string | null = null) {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token!);
    }
  });

  failedQueue = [];
}

/**
 * Factory that reads env once and builds the client.
 */
function createApis(): ApiWrappers {
  const mapexBaseURL = process.env.MAPEXOS_API_BASE_URL as string;
  const httpGatewayBaseURL = process.env.HTTP_GATEWAY_API_BASE_URL as string;
  const assetsBaseURL = process.env.ASSETS_API_BASE_URL as string;
  const eventsBaseURL = process.env.EVENTS_API_BASE_URL as string;
  const routerBaseURL = process.env.ROUTER_API_BASE_URL as string;
  const jsExecutorBaseURL = process.env.JS_EXECUTOR_API_BASE_URL as string;
  const triggersBaseURL = process.env.TRIGGERS_API_BASE_URL as string;
  const workflowBaseURL = process.env.WORKFLOW_API_BASE_URL as string;
  const vaultBaseURL = process.env.VAULT_API_BASE_URL as string;

  if (!mapexBaseURL) logger.warn('MAPEXOS_API_BASE_URL is empty');
  if (!httpGatewayBaseURL) logger.warn('HTTP_GATEWAY_API_BASE_URL is empty');
  if (!assetsBaseURL) logger.warn('ASSETS_API_BASE_URL is empty');
  if (!eventsBaseURL) logger.warn('EVENTS_API_BASE_URL is empty');
  if (!routerBaseURL) logger.warn('ROUTER_API_BASE_URL is empty');
  if (!jsExecutorBaseURL) logger.warn('JS_EXECUTOR_API_BASE_URL is empty');
  if (!triggersBaseURL) logger.warn('TRIGGERS_API_BASE_URL is empty');
  if (!workflowBaseURL) logger.warn('WORKFLOW_API_BASE_URL is empty');
  if (!vaultBaseURL) logger.warn('VAULT_API_BASE_URL is empty');

  return initializeApis({
    mapexOS: { baseURL: mapexBaseURL },
    assets: { baseURL: assetsBaseURL },
    events: { baseURL: eventsBaseURL },
    router: { baseURL: routerBaseURL },
    httpGateway: { baseURL: httpGatewayBaseURL },
    jsExecutor: { baseURL: jsExecutorBaseURL },
    triggers: { baseURL: triggersBaseURL },
    workflows: { baseURL: workflowBaseURL },
    vault: { baseURL: vaultBaseURL },

    // All interceptors and session config
    interceptors: {
      onRequest: (config: any) => {
        // Inject Authorization header with Bearer token from auth store
        const authStore = useAuthStore();
        const token = authStore.accessToken;

        if (token) {
          config.headers = config.headers || {};
          config.headers['Authorization'] = `Bearer ${token}`;
        }

        // Inject X-Org-Context header with selected organization ID
        const organizationStore = useOrganizationStore();
        const orgId = organizationStore.selectedOrganizationId;

        if (orgId) {
          config.headers = config.headers || {};
          config.headers['X-Org-Context'] = orgId;
        }

        return config;
      },
      onResponse: (response: any) => {
        // Handle response interceptors if needed
        return response;
      },
      onError: async (error: AxiosError) => {
        const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

        // Check if this is a 401 error and not already retried
        if (error.response?.status === 401 && !originalRequest._retry) {

          // Skip refresh for login and refresh endpoints to prevent infinite loops
          if (originalRequest.url?.includes('/auth/login') || originalRequest.url?.includes('/auth/refresh')) {
            return Promise.reject(error);
          }

          // If already refreshing, queue this request
          if (isRefreshing) {
            try {
              // Create a Promise that will be resolved/rejected by processQueue()
              // when the refresh completes. This makes the request wait until
              // a new token is available.
              const token = await new Promise<string>((resolve, reject) => {

                // Store resolve/reject in queue - they will be called later
                failedQueue.push({ resolve, reject });
              });              

              originalRequest.headers['Authorization'] = `Bearer ${token}`;
              const axios = (await import('axios')).default;
              return axios(originalRequest);
            } catch (err) {
              // processQueue() rejected with error (refresh failed)
              const error = err instanceof Error ? err : new Error('Queued request failed');
              return Promise.reject(error);
            }
          }

          originalRequest._retry = true;
          isRefreshing = true;

          const authStore = useAuthStore();

          try {
            if (!authStore.refreshToken) {
              throw new Error('No refresh token available');
            }

            // Use axios directly to call refresh endpoint
            const axios = (await import('axios')).default;
            const headers = {'X-Refresh-Token': authStore.refreshToken};
            const response = await axios.post(`${mapexBaseURL}/auth/refresh`, {}, { headers });
            const { access_token, refresh_token } = response?.data?.data || {};

            // Update tokens in store and storage
            authStore.updateTokens(access_token, refresh_token);

            // Update the failed request's header
            originalRequest.headers['Authorization'] = `Bearer ${access_token}`;

            // Wake up all queued requests with the new token
            // This calls resolve(access_token) on all Promises waiting in the queue
            processQueue(null, access_token);

            // Retry original request
            return axios(originalRequest);
          } catch (err) {
            // Refresh failed - logout user            
            const error = err instanceof Error ? err : new Error('Refresh token failed');

            // Reject all queued requests with error
            // This calls reject(error) on all Promises waiting in the queue
            processQueue(error, null);
            authStore.logout();

            // Redirect to login
            if (typeof window !== 'undefined') {
              window.location.href = '/#/';
            }

            return Promise.reject(error);
          } finally {
            isRefreshing = false;
          }
        }

        return Promise.reject(error);
      },
    }
  });
}

/**
 * MapexOS-specific API type with all services guaranteed to be initialized
 */
export type MapexOSApis = Required<ApiWrappers>;

/**
 * HMR-safe global singleton cache.
 * Vite recreates modules on HMR, so we pin the instance on globalThis.
 */
const g = globalThis as unknown as { __MAPEX_APIS__?: MapexOSApis };
if (!g.__MAPEX_APIS__) {
  g.__MAPEX_APIS__ = createApis() as MapexOSApis;
}

/** The singleton instance to import everywhere */
export const apis: MapexOSApis = g.__MAPEX_APIS__;
export default apis;
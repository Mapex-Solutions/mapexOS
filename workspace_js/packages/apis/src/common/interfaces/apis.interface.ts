import { AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import type { Agent as HttpsAgent } from 'https'

/**
 * Get token function type.
 */
export type GetToken = () => string | Promise<string> | undefined

/**
 * API interceptors configuration interface.
 */
export interface ApiInterceptors {
  onRequest?: (config: InternalAxiosRequestConfig) => InternalAxiosRequestConfig | Promise<InternalAxiosRequestConfig>
  onResponse?: (response: AxiosResponse) => AxiosResponse | Promise<AxiosResponse>
  onError?: (error: any) => any
}

/**
 * Sessions configuration interface.
 */
export interface SessionsConfig {
  interceptors?: ApiInterceptors
  getToken?: GetToken
}

/**
 * API initialization configuration interface.
 */
export interface ApiConfig {
  baseURL: string
  headers?: { [key: string]: string }
  httpsAgent?: HttpsAgent

  /**
   * Local interceptors for all APIs
   * If local interceptors and global are provided, we will use local interceptors.
   */
  interceptors?: ApiInterceptors
}

/**
 * Interface for API initialization configuration.
 */
export interface ApiInitConfig extends SessionsConfig {

  /**
   * MapexOS platform services
   */
  mapexOS?: ApiConfig
  assets?: ApiConfig
  events?: ApiConfig
  router?: ApiConfig
  httpGateway?: ApiConfig
  jsExecutor?: ApiConfig
  triggers?: ApiConfig
  workflows?: ApiConfig
  vault?: ApiConfig
}

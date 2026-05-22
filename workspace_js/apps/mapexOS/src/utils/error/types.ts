import type { ComputedRef } from 'vue';

/**
 * HTTP Error Status Codes
 */
export enum HttpErrorCode {
  BAD_REQUEST = 400,
  UNAUTHORIZED = 401,
  FORBIDDEN = 403,
  NOT_FOUND = 404,
  CONFLICT = 409,
  UNPROCESSABLE_ENTITY = 422,
  INTERNAL_SERVER_ERROR = 500,
  SERVICE_UNAVAILABLE = 503,
}

/**
 * Error message map - Maps HTTP status codes or special keys to i18n message keys
 * Supports numeric HTTP status codes (400, 404, etc.) and special keys ('network', 'unknown')
 */
export interface ErrorMessageMap {
  [key: number]: ComputedRef<string> | string;
  network?: ComputedRef<string> | string;
  unknown?: ComputedRef<string> | string;
}

/**
 * API Error with response metadata
 */
export interface ApiError {
  response?: {
    status: number;
    data?: any;
  };
  message?: string;
}

/**
 * Options for handleApiError utility
 */
export interface HandleApiErrorOptions {
  /**
   * Custom error messages mapped to HTTP status codes
   * Example: { 409: t.notifications.alreadyExists, 422: t.notifications.validationFailed }
   */
  customMessages?: ErrorMessageMap;

  /**
   * Default fallback message if no custom message is found
   * If not provided, uses generic error message from global i18n
   */
  defaultMessage?: ComputedRef<string> | string;

  /**
   * Notification timeout in milliseconds
   * @default 5000
   */
  timeout?: number;

  /**
   * Whether to log the error to console
   * @default true
   */
  logError?: boolean;

  /**
   * Custom error logger function
   */
  onError?: (error: ApiError) => void;
}

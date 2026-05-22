/**
 * Timezone configuration
 */
export interface TimezoneConfig {
  /** How timezone is resolved */
  type: 'literal' | 'variable';

  /** Timezone value (IANA name or variable path) */
  value: string;
}

/**
 * Global retry policy for workflow node failures
 */
export interface RetryPolicy {
  /** Whether retry is enabled */
  enabled: boolean;

  /** Maximum number of retry attempts */
  maxAttempts: number;

  /** Initial interval between retries (e.g. '1s', '5s', '30s') */
  initialInterval: string;

  /** Backoff multiplier (e.g. 2.0 = exponential) */
  backoffMultiplier: number;

  /** Maximum interval cap (e.g. '5m') */
  maxInterval: string;

  /** Error codes that should NOT be retried (e.g. 'VALIDATION_FAILED') */
  nonRetryableErrors: string[];
}

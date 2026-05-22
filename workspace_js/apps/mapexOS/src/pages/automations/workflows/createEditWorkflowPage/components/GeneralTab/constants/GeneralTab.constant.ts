/**
 * Timezone type options
 */
export const TIMEZONE_TYPE_OPTIONS = [
  { label: 'Literal', value: 'literal' },
  { label: 'Inputs', value: 'variable' },
] as const;

/**
 * All IANA timezone options built from Intl API.
 * Format: { label: 'America/Sao_Paulo', value: 'America/Sao_Paulo', region: 'America' }
 */
export const IANA_TIMEZONE_OPTIONS = Intl.supportedValuesOf('timeZone').map((tz) => ({
  label: tz,
  value: tz,
  region: tz.split('/')[0] ?? '',
}));

/**
 * Standard non-retryable error codes for workflow retry policy
 */
export const NON_RETRYABLE_ERROR_OPTIONS = [
  { label: 'VALIDATION_FAILED', value: 'VALIDATION_FAILED', description: 'Input data failed validation' },
  { label: 'PERMISSION_DENIED', value: 'PERMISSION_DENIED', description: 'Insufficient permissions' },
  { label: 'NOT_FOUND', value: 'NOT_FOUND', description: 'Resource not found' },
  { label: 'INVALID_CONFIG', value: 'INVALID_CONFIG', description: 'Node configuration is invalid' },
  { label: 'SCRIPT_ERROR', value: 'SCRIPT_ERROR', description: 'Code node execution error' },
  { label: 'TIMEOUT', value: 'TIMEOUT', description: 'Execution timed out' },
  { label: 'CANCELLED', value: 'CANCELLED', description: 'Execution was cancelled' },
  { label: 'RATE_LIMITED', value: 'RATE_LIMITED', description: 'Rate limit exceeded' },
] as const;

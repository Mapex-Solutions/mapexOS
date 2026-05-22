import type { HttpDataSource } from '../interfaces/httpDataSource.interface';

/**
 * Initial step number for the HTTP data source form
 */
export const INITIAL_STEP = 1;

/**
 * Total number of steps in the HTTP data source form
 */
export const TOTAL_STEPS = 5;

/**
 * Step numbers enum for better readability
 */
export const STEP = {
  BASIC_INFO: 1,
  WORKING_HOURS: 2,
  AUTHENTICATION: 3,
  ASSET_BINDING: 4,
  REVIEW: 5,
} as const;

/**
 * Default initial values for HTTP data source form
 * Enabled defaults to true, Protocol/Mode always HTTP Push
 */
export const HTTP_DATASOURCE_DEFAULTS: HttpDataSource = {
  name: '',
  description: '',
  enabled: true, // Default enabled
  mode: 'Push', // Always Push for HTTP Gateway
  protocol: 'HTTP', // Always HTTP for HTTP Gateway
  enableWorkingHours: false,
  daysOfWeek: [],
  timeIntervals: [{ startTime: '09:00', endTime: '17:00' }],
  timezone: '',
  enableRateLimit: false,
  rateLimitType: null,
  rateLimitValue: 100,
  burstCapacity: 200,
  actionOnExceed: null,
  authType: null,
  apiKey: { headerApiKey: 'X-API-Key', valueApiKey: '' },
  jwt: { secretKey: '', headerName: 'Authorization' },
  ipWhitelist: { addresses: [] },
  oauth2: { jwksUrl: '' },
  bindingMode: null,
  directAssetId: null,
  directAssetIdPath: null,
  assetTemplateIds: [],
  customUuidPaths: [{ path: '' }],
  finalUuidPaths: [],
  payloadExample: '',
};

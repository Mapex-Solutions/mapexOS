// Time Interval Configuration
export interface TimeInterval {
  startTime: string;
  endTime: string;
}

// API Key Authentication
export interface ApiKeyConfig {
  headerApiKey: string;
  valueApiKey: string;
}

// JWT Authentication (HS256 Shared Key)
export interface JwtConfig {
  secretKey: string;
  headerName: string;
}

// IP Whitelist Authentication
export interface IpWhitelistConfig {
  addresses: string[];
}

// OAuth2 Authentication
export interface OAuth2Config {
  jwksUrl: string;
}

// Custom UUID Path for Asset Binding
export interface CustomUuidPath {
  path: string;
}

// Main Data Source Data
export interface HttpDataSource {
  name: string;
  description: string;
  enabled: boolean;
  mode: string | null;
  protocol: string | null;

  // Working Hours
  enableWorkingHours: boolean;
  daysOfWeek: string[];
  timeIntervals: TimeInterval[];
  timezone: string;

  // Rate Limit
  enableRateLimit: boolean;
  rateLimitType: string | null;
  rateLimitValue: number;
  burstCapacity: number;
  actionOnExceed: string | null;

  // Authentication
  authType: string | null;
  apiKey: ApiKeyConfig;
  jwt: JwtConfig;
  ipWhitelist: IpWhitelistConfig;
  oauth2: OAuth2Config;

  // Asset Binding - Enterprise Pattern
  bindingMode: string | null;
  directAssetId: string | null;              // For "Direct" mode - Asset MongoDB ID (for display)
  directAssetIdPath: string | null;          // For "Direct" mode - Path to extract UUID from payload (used by backend)
  assetTemplateIds: string[];                // For "FieldMapping" mode - template selection
  customUuidPaths: CustomUuidPath[];         // For "FieldMapping" mode - manual paths
  finalUuidPaths?: string[];                 // Computed - merged array for backend
  payloadExample: string;
}

// Step Props (for child components)
export interface StepProps {
  modelValue: Partial<HttpDataSource>;
}

export interface StepEmits {
  (e: 'update:modelValue', value: Partial<HttpDataSource>): void;
}

// Step Definition
export interface StepDefinition {
  label: string;
  icon: string;
  description: string;
}

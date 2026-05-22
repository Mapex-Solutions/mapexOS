import type { StepDefinition } from '../interfaces/httpDataSource.interface';

// Mode Options
// NOTE: Pull mode commented out - HTTP Gateway only supports Push mode
// Pull mode will be available in future dedicated gateway services
export const MODE_OPTIONS = [
  // { label: 'Pull', value: 'Pull', disable: true }, // Future: Polling Gateway
  { label: 'Push', value: 'Push' },
];

// Protocol Options
// NOTE: MQTT commented out - MQTT is a device protocol (configured in Asset)
// HTTP Gateway only handles HTTP/HTTPS communication
export const PROTOCOL_OPTIONS = [
  { label: 'HTTP', value: 'HTTP' },
  // { label: 'MQTT', value: 'MQTT', disable: true }, // Device protocol - see Asset configuration
];

// Rate Limit Type Options (Backend expects: 'second', 'minute', 'hour')
export const RATE_LIMIT_TYPE_OPTIONS = [
  { label: 'Requests per Second', value: 'second' },
  { label: 'Requests per Minute', value: 'minute' },
  { label: 'Requests per Hour', value: 'hour' },
];

// Rate Limit Action Options (Backend expects: 'drop', 'queue')
export const RATE_LIMIT_ACTION_OPTIONS = [
  { label: 'Drop Request', value: 'drop' },
  { label: 'Queue Request', value: 'queue' },
];

// Authentication Type Options (Backend expects: 'apiKey', 'jwt', 'ip_whitelist', 'oauth2', 'none')
export const AUTH_TYPE_OPTIONS = [
  { label: 'API Key', value: 'apiKey' },
  { label: 'JWT', value: 'jwt' },
  { label: 'IP Whitelist', value: 'ip_whitelist' },
  { label: 'OAuth2', value: 'oauth2' },
  { label: 'None', value: 'none' },
];

// Asset Binding Options (Backend expects: 'fixedAssetId', 'uuidField')
export const ASSET_BINDING_OPTIONS = [
  { label: 'Direct (Fixed Asset)', value: 'fixedAssetId' },
  { label: 'Field Mapping (Dynamic UUID)', value: 'uuidField' },
];

// Step Definitions (removed Protocol & Mode step - always HTTP Push)
export const STEPS: StepDefinition[] = [
  { label: 'Basic Information', icon: 'info', description: 'Name, description and status of the data source' },
  { label: 'Working Hours & Rate Limit', icon: 'schedule', description: 'Configure working hours and rate limit' },
  { label: 'Authentication', icon: 'lock', description: 'Configure authentication' },
  { label: 'Asset Binding', icon: 'device_unknown', description: 'Configure asset binding' },
  { label: 'Review', icon: 'check_circle', description: 'Review and confirm data source settings' },
];

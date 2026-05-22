export const ASSET_TEMPLATE_LIST_STUB = [
  {
    // Header card
    id: 1,
    name: 'Temperature Sensor Handler',
    description: 'Processes temperature data from Acme Corp sensors',
    icon: 'thermostat',
    status: 'Active',

    // Body card
    manufacturer: 'Acme Corp',
    deviceModel: 'TS-100',
    version: 'v1.0.0',
    hasPreprocessor: true,
    hasValidation: true,
    hasConversion: false,

    // Footer card
    created: '2024-03-15T10:30:00',
    updated: '2024-03-15T10:30:00',
  },
  {
    id: 2,
    name: 'Humidity Sensor Handler',
    description: 'Handles humidity readings from Beta Industries devices',
    icon: 'water_drop',
    status: 'Active',

    manufacturer: 'Beta Industries',
    deviceModel: 'HS-200',
    version: 'v1.1.0',
    hasPreprocessor: false,
    hasValidation: true,
    hasConversion: true,

    created: '2024-04-01T08:45:00',
    updated: '2024-04-02T09:00:00',
  },
  {
    id: 3,
    name: 'Pressure Sensor Handler',
    description: 'Converts pressure data for Gamma Solutions units',
    icon: 'compress',
    status: 'Inactive',

    manufacturer: 'Gamma Solutions',
    deviceModel: 'PS-300',
    version: 'v2.0.1',
    hasPreprocessor: true,
    hasValidation: false,
    hasConversion: true,

    created: '2024-05-10T14:20:00',
    updated: '2024-05-10T14:20:00',
  },
  {
    id: 4,
    name: 'Air Quality Handler',
    description: 'Validates and forwards air quality indices',
    icon: 'co2',
    status: 'Active',

    manufacturer: 'Delta Enviro',
    deviceModel: 'AQ-50X',
    version: 'v1.0.5',
    hasPreprocessor: true,
    hasValidation: true,
    hasConversion: true,

    created: '2024-06-12T11:15:00',
    updated: '2024-06-13T12:00:00',
  },
  {
    id: 5,
    name: 'Light Sensor Handler',
    description: 'Processes luminosity readings from Epsilon Tech',
    icon: 'brightness_auto',
    status: 'Active',

    manufacturer: 'Epsilon Tech',
    deviceModel: 'LS-900',
    version: 'v1.3.2',
    hasPreprocessor: false,
    hasValidation: false,
    hasConversion: false,

    created: '2024-07-05T09:00:00',
    updated: '2024-07-05T09:00:00',
  },
  {
    id: 6,
    name: 'Motion Detector Handler',
    description: 'Filters and alerts on motion events',
    icon: 'waving_hand',
    status: 'Inactive',

    manufacturer: 'Zeta Security',
    deviceModel: 'MD-Pro',
    version: 'v2.2.0',
    hasPreprocessor: true,
    hasValidation: true,
    hasConversion: false,

    created: '2024-08-20T16:30:00',
    updated: '2024-09-01T10:00:00',
  },
  {
    id: 7,
    name: 'GPS Tracker Handler',
    description: 'Normalizes GPS data for fleet tracking',
    icon: 'gps_fixed',
    status: 'Active',

    manufacturer: 'Eta Navigation',
    deviceModel: 'GPS-X1',
    version: 'v1.0.3',
    hasPreprocessor: true,
    hasValidation: false,
    hasConversion: true,

    created: '2024-09-15T13:45:00',
    updated: '2024-09-15T13:45:00',
  },
  {
    id: 8,
    name: 'Battery Monitor Handler',
    description: 'Monitors battery health metrics',
    icon: 'battery_charging_full',
    status: 'Active',

    manufacturer: 'Theta Power',
    deviceModel: 'BM-7',
    version: 'v3.0.0',
    hasPreprocessor: false,
    hasValidation: true,
    hasConversion: true,

    created: '2024-10-01T07:30:00',
    updated: '2024-10-02T08:00:00',
  },
  {
    id: 9,
    name: 'Proximity Sensor Handler',
    description: 'Handles proximity events for Kappa Devices',
    icon: 'gesture',
    status: 'Inactive',

    manufacturer: 'Kappa Devices',
    deviceModel: 'PRX-12',
    version: 'v1.5.4',
    hasPreprocessor: true,
    hasValidation: true,
    hasConversion: false,

    created: '2024-11-11T15:10:00',
    updated: '2024-11-11T15:10:00',
  },
];

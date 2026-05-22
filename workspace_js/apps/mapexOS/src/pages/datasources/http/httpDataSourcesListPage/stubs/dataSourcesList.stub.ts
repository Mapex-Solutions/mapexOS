export const DATA_SOURCES_LIST_STUB = [
  {
    id: 1,
    // Card Header
    name: 'Temperature Sensor',
    description: 'Fetches temperature data every 120s via HTTP',
    icon: 'thermostat',
    iconToolTip: 'PULL: This source will fetch data periodically',
    status: 'Active',

    // Card body
    protocol: 'HTTP',
    auth: 'ApiKey',
    rateLimitEnabled: true,
    workingHoursEnabled: false,

    // Card footer
    created: '2024-05-10T08:00:00',
    updated: '2024-05-12T09:15:00',
  },
  {
    id: 2,
    // Card Header
    name: 'Humidity Sensor',
    description: 'Retrieves humidity readings every 60s via MQTT',
    icon: 'water_drop',
    iconToolTip: 'PULL: Subscribes to humidity topic',
    status: 'Active',

    // Card body
    protocol: 'MQTT',
    auth: 'TLS',
    rateLimitEnabled: false,
    workingHoursEnabled: true,

    // Card footer
    created: '2024-05-11T10:20:00',
    updated: '2024-05-13T11:30:00',
  },
  {
    id: 3,
    // Card Header
    name: 'Pressure Sensor',
    description: 'Pushes pressure data on threshold breach',
    icon: 'speed',
    iconToolTip: 'PUSH: Sends data when value exceeds limits',
    status: 'Paused',

    // Card body
    protocol: 'WebSocket',
    auth: 'OAuth2',
    rateLimitEnabled: true,
    workingHoursEnabled: true,

    // Card footer
    created: '2024-05-12T09:00:00',
    updated: '2024-05-14T14:45:00',
  },
  {
    id: 4,
    // Card Header
    name: 'Motion Detector',
    description: 'Waits for motion events via HTTP callback',
    icon: 'motion_photos_on',
    iconToolTip: 'PUSH: HTTP callback on motion detected',
    status: 'Active',

    // Card body
    protocol: 'HTTP',
    auth: 'BearerToken',
    rateLimitEnabled: false,
    workingHoursEnabled: false,

    // Card footer
    created: '2024-05-13T13:15:00',
    updated: '2024-05-15T16:00:00',
  },
  {
    id: 5,
    // Card Header
    name: 'GPS Tracker',
    description: 'Polls GPS location every 300s via HTTP',
    icon: 'gps_fixed',
    iconToolTip: 'PULL: Periodic location polling',
    status: 'Active',

    // Card body
    protocol: 'HTTP',
    auth: 'ApiKey',
    rateLimitEnabled: true,
    workingHoursEnabled: false,

    // Card footer
    created: '2024-05-14T07:45:00',
    updated: '2024-05-16T08:20:00',
  },
  {
    id: 6,
    // Card Header
    name: 'Camera Feed',
    description: 'Streams images via RTSP protocol',
    icon: 'camera_alt',
    iconToolTip: 'PULL: Connects to RTSP stream',
    status: 'Inactive',

    // Card body
    protocol: 'RTSP',
    auth: 'None',
    rateLimitEnabled: false,
    workingHoursEnabled: false,

    // Card footer
    created: '2024-05-15T12:00:00',
    updated: '2024-05-17T12:30:00',
  },
  {
    id: 7,
    // Card Header
    name: 'Energy Meter',
    description: 'Receives consumption data via MQTT push',
    icon: 'flash_on',
    iconToolTip: 'PUSH: Real-time energy usage',
    status: 'Active',

    // Card body
    protocol: 'MQTT',
    auth: 'TLS',
    rateLimitEnabled: true,
    workingHoursEnabled: true,

    // Card footer
    created: '2024-05-16T11:10:00',
    updated: '2024-05-18T13:25:00',
  },
  {
    id: 8,
    // Card Header
    name: 'Air Quality Monitor',
    description: 'Fetches air quality index every 180s',
    icon: 'air',
    iconToolTip: 'PULL: Periodic AQI polling',
    status: 'Active',

    // Card body
    protocol: 'HTTP',
    auth: 'ApiKey',
    rateLimitEnabled: true,
    workingHoursEnabled: false,

    // Card footer
    created: '2024-05-17T08:30:00',
    updated: '2024-05-19T09:40:00',
  },
  {
    id: 9,
    // Card Header
    name: 'Soil Moisture Sensor',
    description: 'Pushes moisture readings when dry',
    icon: 'grass',
    iconToolTip: 'PUSH: Alerts when moisture below threshold',
    status: 'Active',

    // Card body
    protocol: 'LoRaWAN',
    auth: 'NetworkKey',
    rateLimitEnabled: false,
    workingHoursEnabled: true,

    // Card footer
    created: '2024-05-18T14:55:00',
    updated: '2024-05-20T15:05:00',
  },
];

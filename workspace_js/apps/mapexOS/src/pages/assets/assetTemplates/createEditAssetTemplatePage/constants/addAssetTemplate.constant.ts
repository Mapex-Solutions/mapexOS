import type { StepperVerticalItem } from '@components/steppers';
import type { TourStepDefinition } from '@composables/tour';
import type { AssetTemplateData, DynamicFieldType } from '../interfaces';

/**
 * Dynamic Fields Limits
 * Based on ClickHouse EVA storage optimization (MAP<UInt16, Type>)
 * UInt16 supports 0-65535, but we limit to practical usage for performance
 */
export const DYNAMIC_FIELDS_MAX = 100;
export const DYNAMIC_FIELDS_WARNING_THRESHOLD = 80;

export const STEPS: StepperVerticalItem[] = [
	{
		title: 'Basic Information',
		icon: 'mdi-information',
		description: 'Add template name, manufacturer, and device details',
	},
	{
		title: 'Asset ID Path',
		icon: 'mdi-routes',
		description: 'Define the path to identify your assets',
	},
	{
		title: 'Preprocessor Script',
		icon: 'mdi-code-braces',
		description: 'Prepare and normalize incoming data (optional)',
	},
	{
		title: 'Validation Script',
		icon: 'mdi-shield-check',
		description: 'Validate data integrity and format (optional)',
	},
	{
		title: 'Conversion Script',
		icon: 'mdi-swap-horizontal',
		description: 'Convert data to standard format (required)',
	},
	{
		title: 'Test Payload',
		icon: 'mdi-flask',
		description: 'Define test payload for validation',
	},
	{
		title: 'Testing',
		icon: 'mdi-test-tube',
		description: 'Execute and validate your scripts',
	},
	{
		title: 'Dynamic Fields',
		icon: 'mdi-database-cog',
		description: 'Map fields for ClickHouse storage',
	},
	{
		title: 'Review',
		icon: 'mdi-clipboard-check',
		description: 'Review all configuration before saving',
	},
];

/**
 * Dynamic Field Type Options
 * Used in Step 8 for selecting field types
 */
export interface DynamicFieldTypeOption {
	label: string;
	value: DynamicFieldType;
	icon: string;
	description: string;
}

export const DYNAMIC_FIELD_TYPE_OPTIONS: DynamicFieldTypeOption[] = [
	{
		label: 'String',
		value: 'string',
		icon: 'mdi-format-text',
		description: 'Text values',
	},
	{
		label: 'Number',
		value: 'number',
		icon: 'mdi-numeric',
		description: 'Integer or decimal values',
	},
	{
		label: 'Boolean',
		value: 'bool',
		icon: 'mdi-toggle-switch',
		description: 'True/False values',
	},
	{
		label: 'Date',
		value: 'date',
		icon: 'mdi-calendar-clock',
		description: 'Date and time values',
	},
	{
		label: 'Geo',
		value: 'geo',
		icon: 'mdi-map-marker',
		description: 'Geographic coordinates (lat/lng)',
	},
];

export const DEFAULT_ASSET_TEMPLATE_DATA: AssetTemplateData = {
	// Step 1: Basic Information
	name: '',
	enabled: true,
	description: '',

	// Asset Classification (FLAT)
	categoryId: undefined,
	categoryName: undefined,
	manufacturerId: undefined,
	manufacturerName: undefined,
	modelId: undefined,
	modelName: undefined,
	version: undefined,

	isSystem: false, // Always false for user-created templates
	isTemplate: false, // Template compartilhado (apenas Vendor/Customer)

	// Step 2: Asset ID Path
	assetIdPath: '',

	// Step 3: Preprocessor Script (optional) - Empty by default
	scriptProcessor: '',

	// Step 4: Validation Script (optional) - Empty by default
	scriptValidator: '',

	// Step 5: Conversion Script (required)
	scriptConversion: `// Conversion Script (Required)
// Convert data to MapexOS StandardizedPayload format
// Available: payload object
// MUST return: const result = StandardizedPayload
//
// StandardizedPayload Format (REQUIRED):
// {
//   eventType: string (required) - Type of event (e.g., "sensor.reading", "device.alert")
//   eventId: string (required) - Unique identifier for this event
//   data: object (required) - The actual data payload
//   metadata: object (optional) - Additional metadata
//   created: string (required) - ISO 8601 timestamp
// }

const result = {
  eventType: "sensor.reading",
  eventId: \`\${payload.deviceId}-\${Date.now()}\`,
  data: {
    temperature: payload.temperature,
    humidity: payload.humidity || null,
    battery: payload.battery || null,
    deviceId: payload.deviceId
  },
  metadata: {
    source: "asset-template",
    version: "1.0.0",
    processingTime: new Date().toISOString()
  },
  created: new Date(payload.timestamp).toISOString()
};`,

	// Step 6: Test Payload (MANDATORY - JSON payload for testing)
	scriptTest: `{
  "deviceId": "SENSOR001",
  "timestamp": "2024-01-15T10:30:00Z",
  "temp_c": 23.5,
  "humidity": 65.2,
  "battery": 85
}`,

	// Available Fields - populated automatically from Step 7 test results
	availableFields: [],

	// Step 8: Dynamic Fields - user-defined field mappings for ClickHouse
	dynamicFields: [],
};

/**
 * Tour step definitions for Asset Template wizard
 * Teaches users about each step of the wizard
 */
export const ASSET_TEMPLATE_TOUR_STEPS: TourStepDefinition[] = [
	{
		element: '#page-header-section',
		translationKey: 'pageOverview',
		side: 'bottom',
		align: 'start',
	},
	{
		element: '#stepper-section',
		translationKey: 'stepper',
		side: 'right',
		align: 'start',
	},
	{
		element: '#step-2',
		translationKey: 'assetIdPath',
		side: 'right',
		align: 'center',
	},
	{
		element: '#step-3',
		translationKey: 'preprocessorScript',
		side: 'right',
		align: 'center',
	},
	{
		element: '#step-4',
		translationKey: 'validationScript',
		side: 'right',
		align: 'center',
	},
	{
		element: '#step-5',
		translationKey: 'conversionScript',
		side: 'right',
		align: 'center',
	},
	{
		element: '#step-6',
		translationKey: 'testPayload',
		side: 'right',
		align: 'center',
	},
	{
		element: '#step-7',
		translationKey: 'testing',
		side: 'right',
		align: 'center',
	},
	{
		element: '#step-8',
		translationKey: 'dynamicFields',
		side: 'right',
		align: 'center',
	},
	{
		element: '#step-9',
		translationKey: 'review',
		side: 'right',
		align: 'center',
	},
];

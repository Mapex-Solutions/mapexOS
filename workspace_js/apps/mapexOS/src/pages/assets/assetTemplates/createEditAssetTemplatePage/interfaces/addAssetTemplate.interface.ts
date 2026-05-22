/**
 * Dynamic Field Type
 * Maps to ClickHouse typed Maps for storage
 */
export type DynamicFieldType = 'string' | 'number' | 'bool' | 'date' | 'geo';

/**
 * Dynamic Field Mapping Interface
 * Defines how to extract and type fields from payload for ClickHouse storage
 *
 * ClickHouse Maps:
 * - string → fieldsString Map(String, String)
 * - number → fieldsFloat64 Map(String, Float64)
 * - bool   → fieldsBool Map(String, UInt8)
 * - date   → fieldsDateTime64 Map(String, DateTime64)
 * - geo    → fieldsGeo Map(String, Tuple(Float64, Float64))
 */
export interface DynamicFieldMapping {
	/** Field name - will be the key in ClickHouse Map (e.g., "temperature") */
	field: string;

	/** Path to extract value from payload - autocomplete from availableFields (e.g., "data.item.value") */
	value: string;

	/** Data type - determines which ClickHouse Map to use */
	type: DynamicFieldType;

	/** Latitude path - only for 'geo' type */
	latitudePath?: string;

	/** Longitude path - only for 'geo' type */
	longitudePath?: string;
}

/**
 * Asset Template Data Interface
 * Maps to AssetTemplateCreate DTO from backend
 */
export interface AssetTemplateData {
	// Step 1: Basic Information
	name: string;
	enabled: boolean;
	description?: string | undefined;

	// Asset Classification (FLAT structure)
	categoryId?: string | undefined;
	categoryName?: string | undefined;
	manufacturerId?: string | undefined;
	manufacturerName?: string | undefined;
	modelId?: string | undefined;
	modelName?: string | undefined;
	version?: string | undefined;

	isSystem: boolean; // Always false for user-created templates
	isTemplate: boolean; // Template compartilhado (apenas Vendor/Customer podem criar)

	// Step 2: Asset ID Path
	assetIdPath: string; // Required: path to unique asset identifier in payload

	// Step 3: Preprocessor Script (optional)
	scriptProcessor?: string;

	// Step 4: Validation Script (REQUIRED per Go contract)
	scriptValidator: string;

	// Step 5: Conversion Script (required)
	scriptConversion: string;

	// Step 6: Test Payload (MANDATORY for CREATE, optional for UPDATE)
	// JSON payload example used to test the conversion script
	scriptTest?: string;

	// Available Fields (generated from Step 7 test)
	/** Array of normalized field paths from StandardizedPayload */
	availableFields?: string[];

	// Step 8: Dynamic Fields Mapping
	/** Dynamic field mappings for ClickHouse typed Maps */
	dynamicFields?: DynamicFieldMapping[];
}

/**
 * Test Step Result Interface
 */
export interface TestStep {
	name: string;
	success: boolean;
	error?: string;
	details?: any; // Additional error details from API
}

/**
 * Test Results Interface
 */
export interface TestResults {
	executed: boolean;
	success: boolean;
	steps: TestStep[];
	output: any;
	standardizedPayload?: any; // The final StandardizedPayload from conversion script
	responseData?: any; // Full data object from API response
	newPayload?: any; // NEW: The converted payload (response.data.data)
	logs: string[];
}

import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsBoolean, IsMongoId, IsString, IsNumber, NumberIntAndPositive } from '@mapexos/validations';

/**
 * ============================================================================
 * BUILDING BLOCK SCHEMAS — HANDLES & AVAILABLE OUTPUTS
 * ============================================================================
 */

/**
 * HandleDef schema - Input or output connection point on a node
 */
export const ZodHandleDefSchema = z.object({
	id: StringAndNotBeEmpty,
	label: StringAndNotBeEmpty,
	position: IsString,
	dataType: IsString.optional(),
	maxConnections: IsNumber.int().optional(),
	color: IsString.optional(),
});

/**
 * AvailableOutput schema - Describes a field available in the node's output
 */
export const ZodAvailableOutputSchema = z.object({
	path: StringAndNotBeEmpty,
	description: StringAndNotBeEmpty,
});

/**
 * @deprecated Use ZodAvailableOutputSchema instead.
 */
export const ZodOutputHintSchema = ZodAvailableOutputSchema;

/**
 * ============================================================================
 * UNIFIED ACTION CONTRACT
 * ============================================================================
 */

/**
 * HttpActionDef schema - HTTP request template
 */
export const ZodHttpActionDefSchema = z.object({
	method: StringAndNotBeEmpty,
	path: StringAndNotBeEmpty,
	headers: z.record(IsString, IsString).optional(),
	body: z.any().optional(),
	timeout: IsNumber.int().optional(),
});

/**
 * ActionOutputDef schema - How to extract/transform the response
 */
export const ZodActionOutputDefSchema = z.object({
	dataPath: IsString.optional(),
	valuePath: IsString.optional(),
	labelPath: IsString.optional(),
	transform: IsString.optional(),
});

/**
 * ActionDef schema - Unified action contract for operations, fetchOptions, credential test, hooks
 */
export const ZodActionDefSchema = z.object({
	type: StringAndNotBeEmpty, // "http", "mqtt", "nats", "script"
	http: ZodHttpActionDefSchema.optional(),
	output: ZodActionOutputDefSchema.optional(),
});

/**
 * ============================================================================
 * FETCH OPTIONS — DESIGN-TIME DYNAMIC DATA FETCHING
 * ============================================================================
 */

/**
 * FetchOptionsPagination schema
 */
export const ZodFetchOptionsPaginationSchema = z.object({
	mode: StringAndNotBeEmpty, // "cursor" | "page"
	cursorParam: IsString.optional(),
	nextCursorPath: IsString.optional(),
	pageParam: IsString.optional(),
	limitParam: IsString.optional(),
	limitDefault: IsNumber.int().optional(),
	totalPath: IsString.optional(),
});

/**
 * FetchOptionsSearch schema
 */
export const ZodFetchOptionsSearchSchema = z.object({
	param: StringAndNotBeEmpty,
	minLength: IsNumber.int().optional(),
});

/**
 * FetchOptionsDef schema - Manifest-level dynamic options loader.
 * Extends ActionDef with pagination and search support.
 */
export const ZodFetchOptionsDefSchema = ZodActionDefSchema.extend({
	pagination: ZodFetchOptionsPaginationSchema.optional(),
	search: ZodFetchOptionsSearchSchema.optional(),
});

/**
 * @deprecated Use ZodFetchOptionsDefSchema instead.
 */
export const ZodLoadOptionsDefSchema = ZodFetchOptionsDefSchema;

/**
 * ============================================================================
 * NODE PROPERTY SYSTEM — DECLARATIVE FORM DEFINITIONS
 * ============================================================================
 */

/**
 * PropertyRendering schema - Visual rendering options
 */
export const ZodPropertyRenderingSchema = z.object({
	multiline: IsBoolean.optional(),
	rows: IsNumber.int().optional(),
	password: IsBoolean.optional(),
	editor: IsString.optional(),
	multipleValues: IsBoolean.optional(),
	min: IsNumber.optional(),
	max: IsNumber.optional(),
	placeholder: IsString.optional(),
	dateOnly: IsBoolean.optional(),
});

/**
 * FetchOptionsRule schema - Determines which fetchOptions entry to use
 */
export const ZodFetchOptionsRuleSchema = z.object({
	when: z.record(IsString, z.array(z.any())),
	key: StringAndNotBeEmpty,
	label: StringAndNotBeEmpty,
});

/**
 * PropertyFetchOptions schema - Dynamic dropdown fetching config
 */
export const ZodPropertyFetchOptionsSchema = z.object({
	rules: z.array(ZodFetchOptionsRuleSchema),
	dependsOn: z.array(IsString).optional(),
});

/**
 * PropertyOption schema - Single option in a dropdown list
 */
export const ZodPropertyOptionSchema = z.object({
	label: StringAndNotBeEmpty,
	value: z.union([IsString, IsNumber, IsBoolean]),
});

/**
 * DisplayOptions schema - Conditional visibility based on other fields
 */
export const ZodDisplayOptionsSchema = z.object({
	show: z.record(IsString, z.array(z.any())).optional(),
});

/**
 * NodePropertyDef schema - Single declarative form field.
 * Supports both new format (rendering, fetchOptions) and legacy (typeOptions, placeholder).
 */
export const ZodNodePropertyDefSchema = z.object({
	name: StringAndNotBeEmpty,
	displayName: StringAndNotBeEmpty,
	type: StringAndNotBeEmpty,
	default: z.any(),
	hint: IsString.optional(),
	required: IsBoolean.optional(),
	options: z.array(ZodPropertyOptionSchema).optional(),
	displayOptions: ZodDisplayOptionsSchema.optional(),
	allowedSources: z.array(IsString).optional(),
	isSecret: IsBoolean.optional(),
	rendering: ZodPropertyRenderingSchema.optional(),
	fetchOptions: ZodPropertyFetchOptionsSchema.optional(),
	// Legacy compat
	placeholder: IsString.optional(),
	typeOptions: z.object({
		multiline: IsBoolean.optional(),
		rows: IsNumber.int().optional(),
		password: IsBoolean.optional(),
		editor: IsString.optional(),
		multipleValues: IsBoolean.optional(),
		minValue: IsNumber.optional(),
		maxValue: IsNumber.optional(),
		placeholder: IsString.optional(),
		dateOnly: IsBoolean.optional(),
		loadOptions: z.object({ key: StringAndNotBeEmpty, label: StringAndNotBeEmpty }).optional(),
		loadOptionsDependsOn: z.array(IsString).optional(),
	}).optional(),
	values: z.array(z.any()).optional(),
	noticeType: IsString.optional(),
});

/**
 * @deprecated Use ZodNodePropertyTypeOptsSchema for legacy compat only.
 */
export const ZodNodePropertyTypeOptsSchema = ZodNodePropertyDefSchema.shape.typeOptions!;
export const ZodLoadOptionsRefSchema = z.object({ key: StringAndNotBeEmpty, label: StringAndNotBeEmpty });

/**
 * ============================================================================
 * NODE HOOKS — LIFECYCLE
 * ============================================================================
 */

/**
 * NodeHooks schema - Lifecycle hooks for a node type
 */
export const ZodNodeHooksSchema = z.object({
	before: ZodActionDefSchema.optional(),
	after: ZodActionDefSchema.optional(),
	destroy: ZodActionDefSchema.optional(),
});

/**
 * ============================================================================
 * NODE TYPE MANIFEST
 * ============================================================================
 */

/**
 * NodeTypeManifest schema - Single node type within a plugin manifest
 */
/**
 * NodeTimeout schema - Default async timeout for a node type
 */
export const ZodNodeTimeoutDefSchema = z.object({
	duration: z.number().int().positive(),
	unit: IsString,
	enableOutput: z.boolean().optional().default(false),
});

export const ZodNodeTypeManifestSchema = z.object({
	type: StringAndNotBeEmpty,
	label: StringAndNotBeEmpty,
	icon: StringAndNotBeEmpty,
	color: StringAndNotBeEmpty,
	description: StringAndNotBeEmpty,
	inputs: z.array(ZodHandleDefSchema),
	outputs: z.array(ZodHandleDefSchema),
	configSchema: z.record(IsString, z.any()).optional(),
	properties: z.array(ZodNodePropertyDefSchema).optional(),
	defaults: z.record(IsString, z.any()).optional(),
	timeout: ZodNodeTimeoutDefSchema.optional(),
	availableOutputs: z.array(ZodAvailableOutputSchema).optional(),
	outputHints: z.array(ZodAvailableOutputSchema).optional(), // Legacy compat
	operations: z.record(IsString, ZodActionDefSchema).optional(),
	hooks: ZodNodeHooksSchema.optional(),
});

/**
 * ============================================================================
 * CREDENTIAL DEFINITIONS
 * ============================================================================
 */

/**
 * CredentialFieldDef schema - Single credential input field
 */
export const ZodCredentialFieldDefSchema = z.object({
	name: StringAndNotBeEmpty,
	displayName: StringAndNotBeEmpty,
	type: StringAndNotBeEmpty,
	required: IsBoolean.optional(),
	isSecret: IsBoolean.optional(),
	hint: IsString.optional(),
	default: z.any().optional(),
	options: z.array(ZodPropertyOptionSchema).optional(),
});

/**
 * CredentialDef schema - Single authentication method for a plugin
 */
export const ZodCredentialDefSchema = z.object({
	id: StringAndNotBeEmpty,
	name: StringAndNotBeEmpty,
	fields: z.array(ZodCredentialFieldDefSchema),
	test: ZodActionDefSchema.optional(),
});

/**
 * @deprecated Legacy schema — test was plain HTTP, not ActionDef
 */
export const ZodCredentialTestDefSchema = z.object({
	method: StringAndNotBeEmpty,
	path: StringAndNotBeEmpty,
	body: z.record(IsString, z.any()).optional(),
});

/**
 * ============================================================================
 * PLUGIN METADATA & DEFAULTS
 * ============================================================================
 */

/**
 * PluginMetadata schema - Visual metadata for marketplace display
 */
export const ZodPluginMetadataSchema = z.object({
	brandIcon: IsString.optional(),
	color: IsString.optional(),
	docs: IsString.optional(),
});

/**
 * PluginDefaults schema - Default values inherited by operations and fetchOptions
 */
export const ZodPluginDefaultsSchema = z.object({
	baseUrl: IsString.optional(),
	timeout: IsNumber.int().optional(),
});

/**
 * ============================================================================
 * PLUGIN DTOs
 * ============================================================================
 */

/**
 * Plugin ID parameter schema (for URL params — MongoDB ObjectID)
 */
export const ZodPluginIdSchema = z.object({
	id: IsMongoId,
});

/**
 * Plugin Query schema - Used for filtering and pagination
 */
export const ZodPluginQuerySchema = z.object({
	name: IsString.max(100).optional(),
	category: IsString.max(50).optional(),
	enabled: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
});

/**
 * Plugin Create schema - Used for creating new plugin manifests
 */
export const ZodPluginCreateSchema = z.object({
	pluginId: StringAndNotBeEmpty,
	name: IsString.min(1).max(255),
	version: StringAndNotBeEmpty,
	category: StringAndNotBeEmpty,
	icon: StringAndNotBeEmpty,
	color: StringAndNotBeEmpty,
	description: IsString.max(1000).optional().default(''),
	author: IsString.max(255).optional().default(''),
	tags: z.array(IsString.max(50)).optional().default([]),
	defaults: ZodPluginDefaultsSchema.optional(),
	credentials: z.array(ZodCredentialDefSchema).optional().default([]),
	fetchOptions: z.record(IsString, ZodFetchOptionsDefSchema).optional(),
	nodeTypes: z.array(ZodNodeTypeManifestSchema).min(1),
	metadata: ZodPluginMetadataSchema.optional(),
	isTemplate: IsBoolean.optional().default(false),
	enabled: IsBoolean.optional().default(true),
});

/**
 * Plugin Update schema - Used for updating existing manifests (all fields optional)
 */
export const ZodPluginUpdateSchema = z.object({
	name: IsString.min(1).max(255).optional(),
	version: StringAndNotBeEmptyOrOptional,
	category: StringAndNotBeEmptyOrOptional,
	icon: StringAndNotBeEmptyOrOptional,
	color: StringAndNotBeEmptyOrOptional,
	description: IsString.max(1000).optional(),
	author: IsString.max(255).optional(),
	tags: z.array(IsString.max(50)).optional(),
	defaults: ZodPluginDefaultsSchema.optional(),
	credentials: z.array(ZodCredentialDefSchema).optional(),
	fetchOptions: z.record(IsString, ZodFetchOptionsDefSchema).optional(),
	nodeTypes: z.array(ZodNodeTypeManifestSchema).optional(),
	metadata: ZodPluginMetadataSchema.optional(),
	enabled: IsBoolean.optional(),
});

/**
 * Plugin Response schema - Full manifest from backend API
 */
export const ZodPluginResponseSchema = z.object({
	id: IsMongoId.optional(),
	pluginId: StringAndNotBeEmpty,
	name: StringAndNotBeEmpty,
	version: StringAndNotBeEmpty,
	category: StringAndNotBeEmpty,
	icon: StringAndNotBeEmpty,
	color: StringAndNotBeEmpty,
	description: StringAndNotBeEmptyOrOptional,
	author: IsString.optional(),
	tags: z.array(IsString).optional(),
	defaults: ZodPluginDefaultsSchema.optional(),
	credentials: z.array(ZodCredentialDefSchema).optional(),
	fetchOptions: z.record(IsString, ZodFetchOptionsDefSchema).optional(),
	nodeTypes: z.array(ZodNodeTypeManifestSchema),
	metadata: ZodPluginMetadataSchema.optional(),
	isTemplate: IsBoolean.optional(),
	orgId: IsMongoId.optional(),
	pathKey: IsString.optional(),
	enabled: IsBoolean.optional(),
	created: StringAndNotBeEmptyOrOptional,
	updated: StringAndNotBeEmptyOrOptional,
});

import type { Message } from '@mapexos/infrastructure';
import type { AssetScripts } from '@modules/scripts/domain/types';
import type { ScriptExecutionResult, AssetReadModel, TemplateReadModel, ScriptServiceInternalDeps, BatchMessageResult, ScriptServiceMetrics, ScriptProcessorMessage, DataSource } from '@modules/scripts/application/types';
import type { ScriptSet, ScriptEngineServicePort, BytecodeCacheContext } from '@modules/engine/application/ports';
import type { ScriptServicePort } from '@modules/scripts/application/ports';
import type { AssetCachePort } from '@modules/scripts/application/ports/asset_cache_port';
import type { TemplateCachePort } from '@modules/scripts/application/ports/template_cache_port';
import type { EventPublisherPort } from '@modules/scripts/application/ports/event_publisher_port';

import type { Logger } from '@mapexos/microservices';

import { handleMqttBatch } from './script.handler_mqtt';
import { handleHttpBatch } from './script.handler_http';

import { isEmpty } from 'lodash';

import { getByPath } from '@mapexos/utils';
import { zodValidationError } from '@mapexos/validations';
import { ZodStandardizedPayloadSchema } from '@mapexos/schemas';

import { OOMError } from '@modules/engine/domain/errors';

import { PUBLIC_ORG_ID } from '@modules/scripts/application/constants';

/**
 * Application service for script execution orchestration.
 *
 * Pure orchestration — no NATS publishing, no cache invalidation, no message preprocessing.
 * Returns domain results only. Publishing and ACK/Nack handled by consumers (interface layer).
 *
 * Responsibilities:
 * - Resolve asset UUID from message
 * - Fetch scripts from cache (via ports)
 * - Execute script pipeline via ScriptEngineService
 * - Return execution result
 */
export class ScriptService implements ScriptServicePort {
	private readonly metrics?: ScriptServiceMetrics;

	constructor(
		private readonly logger: Logger,
		private readonly scriptEngine: ScriptEngineServicePort,
		private readonly assetCachePort: AssetCachePort,
		private readonly templateCachePort: TemplateCachePort,
		private readonly eventPublisher: EventPublisherPort,
		metrics?: ScriptServiceMetrics,
	) {
		this.metrics = metrics;
	}

	/** Internal deps passed to handler functions (TS has no partial classes) */
	private get internalDeps(): ScriptServiceInternalDeps {
		return {
			logger: this.logger,
			scriptEngine: this.scriptEngine,
			assetCachePort: this.assetCachePort,
			templateCachePort: this.templateCachePort,
			eventPublisher: this.eventPublisher,
			metrics: this.metrics,
		};
	}

	async handleMqttBatch(messages: Message[]): Promise<BatchMessageResult[]> {
		return handleMqttBatch(this.internalDeps, messages);
	}

	async handleHttpBatch(messages: Message[]): Promise<BatchMessageResult[]> {
		return handleHttpBatch(this.internalDeps, messages, this.resolveAssetUUID.bind(this));
	}

	/**
	 * Tests script execution with provided payload and scripts.
	 * Used by HTTP test endpoints.
	 *
	 * @param payload - Raw payload to process
	 * @param scripts - Scripts to execute (decode, validation, transform)
	 * @returns Script execution result
	 */
	async scripsTest(payload: any, scripts: AssetScripts): Promise<ScriptExecutionResult> {
		return this.runScriptPipeline(payload, scripts);
	}

	/**
	 * Fetches the scriptTest field from a template.
	 * Used by HTTP endpoints for template script testing.
	 *
	 * @param orgId - Organization ID (or "mapexos_public" for system templates)
	 * @param templateId - Template ID to fetch
	 * @returns Parsed scriptTest payload or null
	 */
	async getScriptTest(orgId: string, templateId: string): Promise<any | null> {
		const cacheKey = `${orgId}/${templateId}`;
		this.logger.info(`[SERVICE:Script] Fetching scriptTest from template: ${cacheKey}`);

		try {
			const template = await this.templateCachePort.get(cacheKey);
			if (!template) {
				this.logger.warn(`[SERVICE:Script] Template not found in cache: ${cacheKey}`);
				return null;
			}

			if (!template.scriptTest) {
				this.logger.info(`[SERVICE:Script] Template has no scriptTest: ${cacheKey}`);
				return null;
			}

			try {
				return JSON.parse(template.scriptTest);
			} catch {
				return template.scriptTest;
			}
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : String(error);
			this.logger.error(`[SERVICE:Script] Error fetching scriptTest: ${errorMessage}`);
			throw error;
		}
	}

	/**
	 * Generates a sample payload by executing template scripts against scriptTest.
	 * Used by Rule Test Runner UI.
	 *
	 * @param orgId - Organization ID (or "mapexos_public" for system templates)
	 * @param templateId - Template ID to fetch
	 * @returns Processed sample payload or null
	 */
	async getSamplePayload(orgId: string, templateId: string): Promise<any | null> {
		const cacheKey = `${orgId}/${templateId}`;
		this.logger.info(`[SERVICE:Script] Generating sample payload for template: ${cacheKey}`);

		try {
			const template = await this.templateCachePort.get(cacheKey);
			if (!template) {
				this.logger.warn(`[SERVICE:Script] Template not found in cache: ${cacheKey}`);
				return null;
			}

			if (!template.scriptTest) {
				this.logger.info(`[SERVICE:Script] Template has no scriptTest: ${cacheKey}`);
				return null;
			}

			let rawPayload: any;
			try {
				rawPayload = JSON.parse(template.scriptTest);
			} catch {
				rawPayload = template.scriptTest;
			}

			const scripts: AssetScripts = {
				decode: template.scriptProcessor || '',
				validation: template.scriptValidator,
				transform: template.scriptConversion,
			};

			if (isEmpty(scripts.transform)) {
				this.logger.warn(`[SERVICE:Script] Template has no transform script: ${cacheKey}`);
				return null;
			}

			const cacheContext: BytecodeCacheContext = {
				templateId: template.id,
				templateOrgId: orgId,
			};

			const result = await this.runScriptPipeline(rawPayload, scripts, cacheContext);

			if (!result.success) {
				throw new Error(`Script execution failed: ${result.error}`);
			}

			return result.standardizedPayload;
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : String(error);
			this.logger.error(`[SERVICE:Script] Error generating sample payload: ${errorMessage}`);
			throw error;
		}
	}

	/**
	 * Executes the full script pipeline for a message.
	 * Returns the execution result WITHOUT publishing — consumers handle publishing.
	 *
	 * Flow:
	 * 1. Resolve assetUUID
	 * 2. Fetch scripts from cache
	 * 3. Update dataSource with asset metadata
	 * 4. Execute decode → validate → transform
	 * 5. Validate standardized payload format
	 * 6. Return result (consumer publishes)
	 *
	 * @param message - Message with event payload and dataSource
	 * @returns Execution result with transformed payload or error details
	 * @throws OOMError - Propagated for NACK retry handling
	 */
	async executeScripts(message: ScriptProcessorMessage): Promise<ScriptExecutionResult> {
		const { event, dataSource, eventTrackerId } = message;
		this.logger.info(`[SERVICE:Script] Starting script execution, eventTrackerId: ${eventTrackerId}`);

		let assetUUID: string;
		let assetId: string;
		let cacheContext: BytecodeCacheContext | undefined;
		let executionResult: ScriptExecutionResult;

		try {
			assetUUID = await this.resolveAssetUUID(message);

			const orgId = dataSource.orgId || '';
			if (!orgId) {
				throw new Error('Missing required field: dataSource.orgId');
			}

			const {
				scripts,
				assetId: fetchedAssetId,
				debugEnabled: assetDebugEnabled,
				cacheContext: fetchedCacheContext,
				assetMetadata,
			} = await this.fetchAssetScripts(orgId, assetUUID);

			assetId = fetchedAssetId;
			cacheContext = fetchedCacheContext;

			// Update dataSource with asset metadata (source of truth)
			dataSource.pathKey = assetMetadata.pathKey;
			dataSource.name = assetMetadata.name;
			dataSource.description = assetMetadata.description;

			// Execute script pipeline
			executionResult = await this.runScriptPipeline(event, scripts, cacheContext);

			if (!executionResult.success) {
				throw new Error(executionResult.error);
			}

			// Validate standardized payload format
			void await ZodStandardizedPayloadSchema.parseAsync(executionResult.standardizedPayload);

			// Enrich result with metadata for consumer publishing
			executionResult.assetUUID = assetUUID;
			executionResult.assetId = assetId;
			executionResult.debugEnabled = assetDebugEnabled;

			this.logger.info(`[SERVICE:Script] Script execution completed: assetUUID=${assetUUID}, time=${executionResult.totalExecutionTime}`);
			return executionResult;

		} catch (error) {
			let errorMessage: any;

			if (error instanceof Error) errorMessage = error.message;
			else if ((error as any).name === 'ZodError') errorMessage = zodValidationError(error);
			else errorMessage = JSON.stringify(error);

			this.logger.error(`[SERVICE:Script] Script execution failed: ${errorMessage}`);

			const result: ScriptExecutionResult = {
				success: false,
				standardizedPayload: isEmpty(executionResult!) ? event : executionResult!.standardizedPayload,
				failedAt: executionResult?.failedAt ?? null,
				totalExecutionTime: executionResult?.totalExecutionTime ?? null,
				error: errorMessage,
				assetUUID: assetUUID! || 'unknown',
				assetId: assetId! || 'unknown',
				debugEnabled: false,
			};

			if (error instanceof OOMError) {
				throw error;
			}

			return result;
		}
	}

	/**
	 * Resolves asset UUID based on source type.
	 * MQTT/LoRaWAN: pre-provided in message. HTTP: extracted from event via assetBind.
	 *
	 * @param message - Script processor message
	 * @returns Resolved asset UUID
	 */
	async resolveAssetUUID(message: ScriptProcessorMessage): Promise<string> {
		const { sourceType, assetUUID, dataSource, event } = message;

		if (sourceType === 'mqtt' || sourceType === 'lorawan') {
			if (!assetUUID) {
				throw new Error(`assetUUID is required for sourceType '${sourceType}'`);
			}
			return assetUUID;
		}

		if (sourceType === 'http') {
			if (!dataSource.assetBind) {
				throw new Error('assetBind is required for HTTP source type');
			}
			return this.getAssetUUIDFromBind(dataSource.assetBind, event);
		}

		throw new Error(`Unknown sourceType: ${sourceType}`);
	}

	/**
	 * Retrieves asset UUID from assetBind configuration.
	 * fixedAssetId: returns fixed ID. uuidField: extracts from event payload.
	 *
	 * @param assetBind - Asset bind configuration
	 * @param event - Event payload to extract from
	 * @returns Asset UUID
	 */
	private async getAssetUUIDFromBind(assetBind: NonNullable<DataSource['assetBind']>, event: any): Promise<string> {
		const { type, data } = assetBind;

		if (type === 'fixedAssetId') {
			if (!data.assetId) {
				throw new Error('fixedAssetId is configured but assetId is not provided');
			}
			return getByPath(event, assetBind.data.uuidField[0]) as Promise<string>;
		}

		if (type === 'uuidField') {
			if (!data.uuidField || data.uuidField.length === 0) {
				throw new Error('uuidField is configured but no paths are provided');
			}

			const assetIds = [];
			for (const assetIdPath of data.uuidField) {
				const assetId = await getByPath(event, assetIdPath);
				if (assetId) assetIds.push(assetId);
			}

			if (assetIds.length === 0) {
				throw new Error(`No assetId found in event for paths: ${data.uuidField.join(', ')}`);
			}

			return assetIds[0];
		}

		throw new Error(`Unknown assetBind type: ${type}`);
	}

	/**
	 * Fetches scripts and metadata for an asset from TieredCache (via ports).
	 *
	 * Flow: asset cache → template cache → extract scripts + metadata.
	 *
	 * @param orgId - Organization ID
	 * @param assetUUID - Device UUID
	 * @returns Scripts, asset ID, debugEnabled, cache context, asset metadata
	 */
	async fetchAssetScripts(orgId: string, assetUUID: string): Promise<{
		scripts: AssetScripts;
		assetId: string;
		debugEnabled: boolean;
		cacheContext: BytecodeCacheContext;
		assetMetadata: { pathKey: string; name: string; description: string };
	}> {
		const assetCacheKey = `${orgId}/${assetUUID}`;
		this.logger.info(`[SERVICE:Script] Fetching asset from cache: ${assetCacheKey}`);

		const asset = await this.assetCachePort.get(assetCacheKey);
		if (!asset) {
			throw new Error(`Asset not found in cache: ${assetCacheKey}`);
		}

		this.logger.info(`[SERVICE:Script] Asset loaded: ${assetCacheKey}, debugEnabled: ${asset.debugEnabled}`);

		if (!asset.assetTemplateId) {
			throw new Error(`Asset ${assetUUID} has no template assigned`);
		}

		const templateOrgId = asset.assetTemplateOrgId || PUBLIC_ORG_ID;
		const templateCacheKey = `${templateOrgId}/${asset.assetTemplateId}`;

		const template = await this.templateCachePort.get(templateCacheKey);
		if (!template) {
			throw new Error(`Template not found in cache: ${templateCacheKey}`);
		}

		this.logger.info(`[SERVICE:Script] Template loaded: ${template.name}`);

		const scripts: AssetScripts = {
			decode: template.scriptProcessor || '',
			validation: template.scriptValidator,
			transform: template.scriptConversion,
		};

		if (isEmpty(scripts.transform)) {
			throw new Error(`Invalid or missing transform script for template: ${asset.assetTemplateId}`);
		}

		const cacheContext: BytecodeCacheContext = {
			templateId: template.id,
			templateOrgId,
		};

		const assetMetadata = {
			pathKey: asset.pathKey || '',
			name: asset.name || '',
			description: asset.description || '',
		};

		return { scripts, assetId: asset.id, debugEnabled: asset.debugEnabled, cacheContext, assetMetadata };
	}

	/**
	 * Executes the script pipeline via ScriptEngineService.
	 *
	 * @param payload - Input data to process
	 * @param scripts - Scripts for decode/validate/transform
	 * @param cacheContext - Template context for bytecode cache
	 * @returns Script execution result
	 */
	private async runScriptPipeline(payload: any, scripts: AssetScripts, cacheContext?: BytecodeCacheContext): Promise<ScriptExecutionResult> {
		this.logger.info('[SERVICE:Script] Starting script pipeline execution');

		const scriptSet: ScriptSet = {
			decode: scripts.decode,
			validation: scripts.validation,
			transform: scripts.transform,
		};

		const engineResult = await this.scriptEngine.runScriptPipeline(payload, scriptSet, cacheContext);

		const result: ScriptExecutionResult = {
			success: engineResult.success,
			standardizedPayload: engineResult.finalPayload,
			failedAt: engineResult.failedAt ?? null,
			totalExecutionTime: engineResult.totalPipelineTime ?? null,
			error: engineResult.error ?? null,
		};

		this.logger.info(`[SERVICE:Script] Pipeline ${result.success ? 'completed' : 'failed'}`);
		return result;
	}
}

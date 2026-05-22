import type { Message } from '@mapexos/infrastructure';
import type { AssetScripts } from '@modules/scripts/domain/types';
import type { ScriptExecutionResult, BatchMessageResult } from '@modules/scripts/application/types';
import type { ScriptProcessorMessage } from '@modules/scripts/application/types';
import type { BytecodeCacheContext } from '@modules/engine/application/ports';

/**
 * Port interface for Script Service (application layer contract).
 *
 * Pure orchestration — resolves assets, executes scripts, returns results.
 * No NATS publishing, no cache invalidation, no message lifecycle management.
 *
 * Consumers (interface layer) call these methods after preprocessing,
 * then handle publishing and ACK/Nack based on the returned results.
 */
export interface ScriptServicePort {
	/**
	 * Tests script execution with provided payload and scripts.
	 * Used by HTTP test endpoints. No publishing.
	 *
	 * @param payload - Input data to process
	 * @param scripts - Scripts to execute (decode, validation, transform)
	 * @returns Script execution result
	 */
	scripsTest(payload: any, scripts: AssetScripts): Promise<ScriptExecutionResult>;

	/**
	 * Executes the complete script pipeline for a message.
	 * Returns result WITHOUT publishing — consumers handle that.
	 *
	 * Result includes enriched metadata (assetUUID, assetId, debugEnabled)
	 * that consumers need for publishing decisions.
	 *
	 * @param message - Message with event payload and dataSource
	 * @returns Execution result with transformed payload or error
	 * @throws OOMError - Propagated for NACK retry handling
	 */
	executeScripts(message: ScriptProcessorMessage): Promise<ScriptExecutionResult>;

	/**
	 * Fetches scriptTest payload from a template via TieredCache.
	 * Used by HTTP endpoints for template testing.
	 *
	 * @param orgId - Organization ID (or "mapexos_public" for system templates)
	 * @param templateId - Template ID
	 * @returns Parsed scriptTest object or null
	 */
	getScriptTest(orgId: string, templateId: string): Promise<any | null>;

	/**
	 * Generates processed sample payload by executing scripts against scriptTest.
	 * Used by Rule Test Runner UI.
	 *
	 * @param orgId - Organization ID
	 * @param templateId - Template ID
	 * @returns Processed payload or null
	 */
	getSamplePayload(orgId: string, templateId: string): Promise<any | null>;

	/**
	 * Fetches scripts and metadata for an asset from cache.
	 * Used by consumers during batch preprocessing (Pass 2).
	 *
	 * @param orgId - Organization ID
	 * @param assetUUID - Device UUID
	 * @returns Scripts, asset ID, debugEnabled, cache context, asset metadata
	 */
	fetchAssetScripts(orgId: string, assetUUID: string): Promise<{
		scripts: AssetScripts;
		assetId: string;
		debugEnabled: boolean;
		cacheContext: BytecodeCacheContext;
		assetMetadata: { pathKey: string; name: string; description: string };
	}>;

	/**
	 * Resolves asset UUID based on source type and message content.
	 * Used by consumers during batch preprocessing (Pass 1).
	 *
	 * @param message - Script processor message
	 * @returns Resolved asset UUID
	 */
	resolveAssetUUID(message: ScriptProcessorMessage): Promise<string>;

	/**
	 * Processes a batch of MQTT telemetry messages.
	 * Parses MQTT subjects to extract orgId/assetUUID, normalizes into
	 * common batch contract, then runs the shared processing pipeline.
	 *
	 * @param messages - Raw NATS messages from MQTT-DATA stream
	 * @returns Per-message results (consumer uses for ACK/Nack decisions)
	 */
	handleMqttBatch(messages: Message[]): Promise<BatchMessageResult[]>;

	/**
	 * Processes a batch of HTTP datasource messages.
	 * Parses JSON payloads with dataSource/event fields, resolves assetUUID,
	 * normalizes into common batch contract, then runs the shared processing pipeline.
	 *
	 * @param messages - Raw NATS messages from PROCESSOR-JS-EXECUTE stream
	 * @returns Per-message results (consumer uses for ACK/Nack decisions)
	 */
	handleHttpBatch(messages: Message[]): Promise<BatchMessageResult[]>;
}

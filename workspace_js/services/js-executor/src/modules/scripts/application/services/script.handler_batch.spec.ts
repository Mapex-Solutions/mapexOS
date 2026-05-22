import 'reflect-metadata';
import { processBatch } from './script.handler_batch';

import type { Logger } from '@mapexos/microservices';
import type { ScriptEngineServicePort } from '@modules/engine/application/ports';
import type { EventPublisherPort } from '@modules/scripts/application/ports/event_publisher_port';
import type { AssetCachePort } from '@modules/scripts/application/ports/asset_cache_port';
import type { TemplateCachePort } from '@modules/scripts/application/ports/template_cache_port';
import type {
	AssetReadModel,
	BatchMessageInput,
	ScriptServiceInternalDeps,
	ScriptServiceMetrics,
	TemplateReadModel,
} from '@modules/scripts/application/types';

/**
 * Heartbeat gate truth-table (TKT-2026-0034):
 *   enabled=false                       → skip (reason='disabled')
 *   enabled=true  + heartbeatMode='implicit' (or undefined) → publish
 *   enabled=true  + heartbeatMode='explicit'                → skip (reason='explicit_mode')
 */
describe('script.handler_batch — publishHeartbeat gate', () => {
	const baseAsset = (): AssetReadModel => ({
		id: 'asset-mongo-id',
		uuid: 'sensor-001',
		orgId: 'org-1',
		pathKey: '000001',
		enabled: true,
		debugEnabled: false,
		name: 'Sensor 1',
		assetTemplateId: 'template-1',
		assetTemplateOrgId: 'org-1',
		created: new Date().toISOString(),
		updated: new Date().toISOString(),
	});

	const baseTemplate = (): TemplateReadModel => ({
		id: 'template-1',
		name: 'Template 1',
		enabled: true,
		scriptProcessor: '',
		scriptValidator: '',
		scriptConversion: 'const result = payload;',
		created: new Date().toISOString(),
		updated: new Date().toISOString(),
	});

	const buildDeps = (asset: AssetReadModel) => {
		const eventPublisher: jest.Mocked<EventPublisherPort> = {
			publishResult: jest.fn(),
			publishRawEvent: jest.fn(),
			publishExecutionLog: jest.fn(),
			publishHeartbeat: jest.fn(),
			flush: jest.fn().mockResolvedValue(undefined),
		} as any;

		const assetCachePort: jest.Mocked<AssetCachePort> = {
			get: jest.fn().mockResolvedValue(asset),
		} as any;

		const templateCachePort: jest.Mocked<TemplateCachePort> = {
			get: jest.fn().mockResolvedValue(baseTemplate()),
		} as any;

		const scriptEngine: jest.Mocked<ScriptEngineServicePort> = {
			runBatch: jest.fn().mockResolvedValue([
				{ success: true, finalPayload: { ok: true }, totalPipelineTime: 1 },
			]),
		} as any;

		const logger: jest.Mocked<Logger> = {
			info: jest.fn(),
			error: jest.fn(),
			warn: jest.fn(),
			debug: jest.fn(),
		} as any;

		const heartbeatsPublished = { inc: jest.fn() };
		const heartbeatsSkipped = { inc: jest.fn() };
		const metrics = {
			heartbeatsPublished,
			heartbeatsSkipped,
		} as unknown as ScriptServiceMetrics;

		const deps: ScriptServiceInternalDeps = {
			logger,
			scriptEngine,
			assetCachePort,
			templateCachePort,
			eventPublisher,
			metrics,
		};

		return { deps, eventPublisher, heartbeatsPublished, heartbeatsSkipped };
	};

	const buildInput = (): BatchMessageInput => ({
		index: 0,
		orgId: 'org-1',
		assetUUID: 'sensor-001',
		event: { temperature: 22 },
		sourceType: 'http',
		eventTrackerId: 'tracker-1',
	});

	it('enabled=true + heartbeatMode=implicit → publishHeartbeat called', async () => {
		const asset = baseAsset();
		asset.healthMonitor = {
			enabled: true,
			thresholdMinutes: 10,
			requiredMisses: 3,
			heartbeatMode: 'implicit',
		};
		const { deps, eventPublisher, heartbeatsPublished, heartbeatsSkipped } = buildDeps(asset);

		await processBatch(deps, [buildInput()]);

		expect(eventPublisher.publishHeartbeat).toHaveBeenCalledTimes(1);
		expect(eventPublisher.publishHeartbeat).toHaveBeenCalledWith({
			orgId: 'org-1',
			assetUUID: 'sensor-001',
			pathKey: '000001',
		});
		expect(heartbeatsPublished.inc).toHaveBeenCalledTimes(1);
		expect(heartbeatsSkipped.inc).not.toHaveBeenCalled();
	});

	it('enabled=true + heartbeatMode=explicit → publishHeartbeat NOT called', async () => {
		const asset = baseAsset();
		asset.healthMonitor = {
			enabled: true,
			thresholdMinutes: 10,
			requiredMisses: 3,
			heartbeatMode: 'explicit',
		};
		const { deps, eventPublisher, heartbeatsPublished, heartbeatsSkipped } = buildDeps(asset);

		await processBatch(deps, [buildInput()]);

		expect(eventPublisher.publishHeartbeat).not.toHaveBeenCalled();
		expect(heartbeatsPublished.inc).not.toHaveBeenCalled();
		expect(heartbeatsSkipped.inc).toHaveBeenCalledWith({ reason: 'explicit_mode' });
	});

	it('enabled=false (any mode) → publishHeartbeat NOT called', async () => {
		const asset = baseAsset();
		asset.healthMonitor = {
			enabled: false,
			thresholdMinutes: 10,
			requiredMisses: 3,
			heartbeatMode: 'explicit',
		};
		const { deps, eventPublisher, heartbeatsPublished, heartbeatsSkipped } = buildDeps(asset);

		await processBatch(deps, [buildInput()]);

		expect(eventPublisher.publishHeartbeat).not.toHaveBeenCalled();
		expect(heartbeatsPublished.inc).not.toHaveBeenCalled();
		expect(heartbeatsSkipped.inc).toHaveBeenCalledWith({ reason: 'disabled' });
	});
});

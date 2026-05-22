import 'reflect-metadata';
import { ScriptService } from './script.service';
import type { ScriptEngineServicePort } from '@modules/engine/application/ports';
import type { Logger } from '@mapexos/microservices';
import type { NatsBus } from '@mapexos/infrastructure';
import type { TieredCacheClient } from '@mapexos/infrastructure';
import type { AssetScripts } from '@modules/scripts/domain/types';
import type { ScriptExecutionResult, AssetReadModel, TemplateReadModel } from '@modules/scripts/application/types';
import type { ScriptProcessorMessage, DataSource } from '@modules/scripts/application/types';
import { CacheTier } from '@mapexos/infrastructure';

// Mock dependencies
jest.mock('@mapexos/validations', () => ({
  zodValidationError: jest.fn().mockReturnValue('Validation error'),
}));

jest.mock('@mapexos/schemas', () => ({
  ZodStandardizedPayloadSchema: {
    parseAsync: jest.fn(),
  },
}));

jest.mock('@mapexos/utils', () => ({
  getByPath: jest.fn(),
}));

describe('ScriptService', () => {
  let scriptService: ScriptService;
  let mockLogger: jest.Mocked<Logger>;
  let mockNatsBus: jest.Mocked<NatsBus>;
  let mockScriptEngine: jest.Mocked<ScriptEngineServicePort>;
  let mockAssetCache: jest.Mocked<TieredCacheClient>;
  let mockTemplateCache: jest.Mocked<TieredCacheClient>;

  const mockAssetScripts: AssetScripts = {
    decode: 'const result = payload;',
    validation: 'const result = payload;',
    transform: 'const result = payload;',
  };

  const mockAssetReadModel: AssetReadModel = {
    id: 'asset-mongo-id-123',
    uuid: 'asset-uuid-123',
    orgId: 'org-123',
    pathKey: '0001/0002/',
    enabled: true,
    debugEnabled: false,
    name: 'Test Asset',
    assetTemplateId: 'template-123',
    created: new Date().toISOString(),
    updated: new Date().toISOString(),
  };

  const mockTemplateReadModel: TemplateReadModel = {
    id: 'template-123',
    name: 'Test Template',
    enabled: true,
    scriptProcessor: 'const result = payload;',
    scriptValidator: 'const result = payload;',
    scriptConversion: 'const result = payload;',
    created: new Date().toISOString(),
    updated: new Date().toISOString(),
  };

  beforeEach(() => {
    // Mock Logger
    mockLogger = {
      info: jest.fn(),
      error: jest.fn(),
      warn: jest.fn(),
      debug: jest.fn(),
    } as any;

    // Mock NatsBus
    mockNatsBus = {
      publish: jest.fn().mockResolvedValue(undefined),
    } as any;

    // Mock ScriptEngineService
    mockScriptEngine = {
      runScriptPipeline: jest.fn(),
    } as any;

    // Mock TieredCacheClient for assets
    mockAssetCache = {
      Get: jest.fn().mockResolvedValue({
        data: Buffer.from(JSON.stringify(mockAssetReadModel)),
        tier: CacheTier.L0,
      }),
      Set: jest.fn(),
      Delete: jest.fn(),
      Invalidate: jest.fn(),

      Warmup: jest.fn(),
    } as any;

    // Mock TieredCacheClient for templates
    mockTemplateCache = {
      Get: jest.fn().mockResolvedValue({
        data: Buffer.from(JSON.stringify(mockTemplateReadModel)),
        tier: CacheTier.L0,
      }),
      Set: jest.fn(),
      Delete: jest.fn(),
      Invalidate: jest.fn(),

      Warmup: jest.fn(),
    } as any;

    const mockEventPublisher = {
      publishResult: jest.fn(),
      publishRawEvent: jest.fn(),
      publishExecutionLog: jest.fn(),
      publishHeartbeat: jest.fn(),
      flush: jest.fn().mockResolvedValue(undefined),
    } as any;

    scriptService = new ScriptService(
      mockLogger,
      mockScriptEngine,
      mockAssetCache as any,
      mockTemplateCache as any,
      mockEventPublisher,
    );

    jest.clearAllMocks();
  });

  describe('executeScripts', () => {
    const mockMessage: ScriptProcessorMessage = {
      sourceType: 'http',
      event: {
        eventId: 'test-event',
        data: { value: 123 }
      },
      dataSource: {
        id: 'test-source',
        orgId: 'test-org-id',
        pathKey: '0001/0002/',
        assetBind: {
          type: 'fixedAssetId',
          data: { assetId: 'asset-123', uuidField: ['metadata.assetId'] }
        }
      } as DataSource,
    };

    const mockExecutionResult: ScriptExecutionResult = {
      success: true,
      standardizedPayload: {
        eventType: 'test',
        eventId: 'test-event',
        data: { processed: true },
        metadata: {},
        created: new Date().toISOString()
      },
      failedAt: undefined,
      totalExecutionTime: 100,
      error: undefined,
    };

    beforeEach(() => {
      // Setup mocks for successful path
      const { getByPath } = require('@mapexos/utils');
      getByPath.mockResolvedValue('asset-uuid-123');

      mockScriptEngine.runScriptPipeline.mockResolvedValue({
        success: true,
        finalPayload: mockExecutionResult.standardizedPayload,
        totalPipelineTime: 100,
      });

      const { ZodStandardizedPayloadSchema } = require('@mapexos/schemas');
      ZodStandardizedPayloadSchema.parseAsync.mockResolvedValue(mockExecutionResult.standardizedPayload);
    });

    it('should execute scripts successfully using TieredCache', async () => {
      const result = await scriptService.executeScripts(mockMessage);

      expect(result.success).toBe(true);
      expect(mockAssetCache.Get).toHaveBeenCalled();
      // Template cache key format: {templateOrgId}/{templateId} - uses PUBLIC_ORG_ID when not specified
      expect(mockTemplateCache.Get).toHaveBeenCalledWith('mapexos_public/template-123');
      expect(mockLogger.info).toHaveBeenCalledWith(
        expect.stringContaining('Starting script execution')
      );
    });

    it('should handle missing asset in cache', async () => {
      mockAssetCache.Get.mockResolvedValue(null);

      const result = await scriptService.executeScripts(mockMessage);

      expect(result.success).toBe(false);
      expect(result.error).toContain('Asset not found');
    });

    it('should handle missing template in cache', async () => {
      mockTemplateCache.Get.mockResolvedValue(null);

      const result = await scriptService.executeScripts(mockMessage);

      expect(result.success).toBe(false);
      expect(result.error).toContain('Template not found');
    });

    it('should handle asset without template assigned', async () => {
      const assetWithoutTemplate = { ...mockAssetReadModel, assetTemplateId: undefined };
      mockAssetCache.Get.mockResolvedValue({
        data: Buffer.from(JSON.stringify(assetWithoutTemplate)),
        tier: CacheTier.L0,
      });

      const result = await scriptService.executeScripts(mockMessage);

      expect(result.success).toBe(false);
      expect(result.error).toContain('no template assigned');
    });

    it('should handle script execution failure', async () => {
      mockScriptEngine.runScriptPipeline.mockResolvedValue({
        success: false,
        finalPayload: null,
        failedAt: 'decode',
        error: 'Script error',
      });

      const result = await scriptService.executeScripts(mockMessage);

      expect(result.success).toBe(false);
      expect(result.error).toBe('Script error');
    });

    it('should publish results to NATS on success', async () => {
      await scriptService.executeScripts(mockMessage);

      expect(mockNatsBus.publish).toHaveBeenCalledWith(
        expect.stringContaining('route.'),
        expect.any(Object)
      );
    });
  });

  describe('scripsTest', () => {
    it('should execute test scripts successfully', async () => {
      const testPayload = { test: 'data' };
      const expectedResult = {
        success: true,
        finalPayload: {
          eventType: 'test',
          eventId: 'test-event',
          data: { processed: true },
          metadata: {},
          created: new Date().toISOString()
        },
        totalPipelineTime: 50,
      };

      mockScriptEngine.runScriptPipeline.mockResolvedValue(expectedResult);

      const result = await scriptService.scripsTest(testPayload, mockAssetScripts);

      expect(result.success).toBe(true);
      expect(mockScriptEngine.runScriptPipeline).toHaveBeenCalledWith(
        testPayload,
        expect.objectContaining({
          decode: mockAssetScripts.decode,
          validation: mockAssetScripts.validation,
          transform: mockAssetScripts.transform,
        }),
        undefined
      );
    });
  });

  // invalidateCache removed — now handled by CacheInvalidationAdapter

  describe('fetchAssetScripts', () => {
    it('should fetch asset and template from TieredCache', async () => {
      const result = await (scriptService as any).fetchAssetScripts('org-123', 'asset-uuid-123');

      // Cache key format: {orgId}/{assetUUID}
      expect(mockAssetCache.Get).toHaveBeenCalledWith('org-123/asset-uuid-123');
      // Template cache key format: {templateOrgId}/{templateId} - uses PUBLIC_ORG_ID when not specified
      expect(mockTemplateCache.Get).toHaveBeenCalledWith('mapexos_public/template-123');
      expect(result.scripts.transform).toBe(mockTemplateReadModel.scriptConversion);
      expect(result.assetId).toBe(mockAssetReadModel.id);
      expect(result.debugEnabled).toBe(false);
      // Check asset metadata is returned
      expect(result.assetMetadata).toEqual({
        pathKey: mockAssetReadModel.pathKey,
        name: mockAssetReadModel.name,
        description: '',
      });
    });

    it('should return debugEnabled from asset', async () => {
      const assetWithDebug = { ...mockAssetReadModel, debugEnabled: true };
      mockAssetCache.Get.mockResolvedValue({
        data: Buffer.from(JSON.stringify(assetWithDebug)),
        tier: CacheTier.L0,
      });

      const result = await (scriptService as any).fetchAssetScripts('org-123', 'asset-uuid-123');

      expect(result.debugEnabled).toBe(true);
    });

    it('should throw when asset not found', async () => {
      mockAssetCache.Get.mockResolvedValue(null);

      await expect((scriptService as any).fetchAssetScripts('org-123', 'missing-asset'))
        .rejects.toThrow('Asset not found');
    });

    it('should throw when template not found', async () => {
      mockTemplateCache.Get.mockResolvedValue(null);

      await expect((scriptService as any).fetchAssetScripts('org-123', 'asset-uuid-123'))
        .rejects.toThrow('Template not found');
    });

    it('should throw when template has no transform script', async () => {
      const templateWithoutScript = { ...mockTemplateReadModel, scriptConversion: '' };
      mockTemplateCache.Get.mockResolvedValue({
        data: Buffer.from(JSON.stringify(templateWithoutScript)),
        tier: CacheTier.L0,
      });

      await expect((scriptService as any).fetchAssetScripts('org-123', 'asset-uuid-123'))
        .rejects.toThrow('Invalid or missing transform script');
    });
  });

  describe('publishResult', () => {
    const mockDataSource: DataSource = {
      id: 'test-source',
      orgId: 'org-123',
      pathKey: '0001/0002/',
    };

    const mockResult: ScriptExecutionResult = {
      success: true,
      standardizedPayload: {
        eventType: 'test',
        eventId: 'test-event',
        data: { test: true },
        metadata: {},
        created: new Date().toISOString()
      },
      failedAt: undefined,
      totalExecutionTime: 100,
      error: undefined,
    };

    it('should publish successful result to NATS', async () => {
      await (scriptService as any).publishResult('asset-uuid', 'asset-id', mockResult, mockDataSource);

      expect(mockNatsBus.publish).toHaveBeenCalledWith(
        'route.execute',
        expect.objectContaining({
          assetUUID: 'asset-uuid',
          assetId: 'asset-id',
          orgId: 'org-123',
        })
      );
    });

    it('should not publish on failure', async () => {
      const failedResult = { ...mockResult, success: false };

      await (scriptService as any).publishResult('asset-uuid', 'asset-id', failedResult, mockDataSource);

      expect(mockNatsBus.publish).not.toHaveBeenCalled();
    });

    it('should handle NATS publishing errors', async () => {
      mockNatsBus.publish.mockRejectedValue(new Error('NATS error'));

      await expect((scriptService as any).publishResult('asset-uuid', 'asset-id', mockResult, mockDataSource))
        .rejects.toThrow('Failed to publish result: NATS error');
    });
  });

  describe('runScriptPipeline', () => {
    it('should delegate to ScriptEngineService', async () => {
      const mockPayload = { test: 'data' };

      mockScriptEngine.runScriptPipeline.mockResolvedValue({
        success: true,
        finalPayload: {
          eventType: 'test',
          eventId: 'test-event',
          data: { processed: true },
          metadata: {},
          created: new Date().toISOString()
        },
        totalPipelineTime: 75,
      });

      const result = await (scriptService as any).runScriptPipeline(
        mockPayload,
        mockAssetScripts,
        'asset-123'
      );

      expect(mockScriptEngine.runScriptPipeline).toHaveBeenCalledWith(
        mockPayload,
        {
          decode: mockAssetScripts.decode,
          validation: mockAssetScripts.validation,
          transform: mockAssetScripts.transform,
        },
        'asset-123'
      );

      expect(result.success).toBe(true);
    });
  });

  describe('debugEnabled functionality', () => {
    const mockMessage: ScriptProcessorMessage = {
      sourceType: 'http',
      event: {
        eventId: 'test-event',
        data: { value: 123 }
      },
      dataSource: {
        id: 'test-source',
        orgId: 'test-org-id',
        pathKey: '0001/0002/',
        assetBind: {
          type: 'fixedAssetId',
          data: { assetId: 'asset-123', uuidField: ['metadata.assetId'] }
        }
      } as DataSource,
    };

    beforeEach(() => {
      const { getByPath } = require('@mapexos/utils');
      getByPath.mockResolvedValue('asset-uuid-123');

      mockScriptEngine.runScriptPipeline.mockResolvedValue({
        success: true,
        finalPayload: { eventType: 'test', eventId: 'test', data: {}, metadata: {}, created: new Date().toISOString() },
        totalPipelineTime: 100,
      });

      const { ZodStandardizedPayloadSchema } = require('@mapexos/schemas');
      ZodStandardizedPayloadSchema.parseAsync.mockResolvedValue({});
    });

    it('should publish debug logs when debugEnabled is true', async () => {
      const assetWithDebug = { ...mockAssetReadModel, debugEnabled: true };
      mockAssetCache.Get.mockResolvedValue({
        data: Buffer.from(JSON.stringify(assetWithDebug)),
        tier: CacheTier.L0,
      });

      await scriptService.executeScripts(mockMessage);

      // Should publish raw event before execution
      expect(mockNatsBus.publish).toHaveBeenCalledWith(
        'events.raw',
        expect.objectContaining({
          threadId: 'asset-uuid-123',
          source: 'http_gateway',
          success: true,
        })
      );

      // Should publish JS execution log after execution
      expect(mockNatsBus.publish).toHaveBeenCalledWith(
        'events.logs.jsexecutor',
        expect.objectContaining({
          threadId: 'asset-uuid-123',
          execution: expect.objectContaining({
            success: true,
          }),
        })
      );

      // Should also publish to router
      expect(mockNatsBus.publish).toHaveBeenCalledWith(
        expect.stringContaining('route.'),
        expect.any(Object)
      );
    });

    it('should NOT publish debug logs when debugEnabled is false', async () => {
      // mockAssetReadModel has debugEnabled: false by default
      await scriptService.executeScripts(mockMessage);

      // Should NOT publish to events.raw (only when debugEnabled is true)
      expect(mockNatsBus.publish).not.toHaveBeenCalledWith(
        'events.raw',
        expect.any(Object)
      );

      // Should NOT publish to events.logs.jsexecutor on success (only when debugEnabled is true)
      expect(mockNatsBus.publish).not.toHaveBeenCalledWith(
        'events.logs.jsexecutor',
        expect.any(Object)
      );

      // Should still publish to router
      expect(mockNatsBus.publish).toHaveBeenCalledWith(
        expect.stringContaining('route.'),
        expect.any(Object)
      );
    });

    it('should handle raw event publish failure gracefully', async () => {
      const assetWithDebug = { ...mockAssetReadModel, debugEnabled: true };
      mockAssetCache.Get.mockResolvedValue({
        data: Buffer.from(JSON.stringify(assetWithDebug)),
        tier: CacheTier.L0,
      });

      // Make first publish fail (events.raw) but don't throw
      mockNatsBus.publish.mockRejectedValueOnce(new Error('NATS error'));
      mockNatsBus.publish.mockResolvedValue(undefined);

      // Should not throw - raw event publish failure is not critical
      const result = await scriptService.executeScripts(mockMessage);

      expect(result.success).toBe(true);
      expect(mockLogger.warn).toHaveBeenCalledWith(
        expect.stringContaining('Failed to publish raw event')
      );
    });
  });

  describe('publishRawEvent', () => {
    const mockDataSource: DataSource = {
      id: 'test-source',
      orgId: 'org-123',
      pathKey: '0001/0002/',
      name: 'Test Source',
      description: 'Test description',
    };

    it('should publish raw event to events.raw', async () => {
      await (scriptService as any).publishRawEvent('asset-uuid', mockDataSource, { test: 'data' }, 'http');

      expect(mockNatsBus.publish).toHaveBeenCalledWith(
        'events.raw',
        expect.objectContaining({
          threadId: 'asset-uuid',
          orgId: 'org-123',
          pathKey: '0001/0002/',
          event: { test: 'data' },
          source: 'http_gateway',
          success: true,
          error: '',
        })
      );
    });

    it('should map mqtt sourceType to mqtt_gateway', async () => {
      await (scriptService as any).publishRawEvent('asset-uuid', mockDataSource, { test: 'data' }, 'mqtt');

      expect(mockNatsBus.publish).toHaveBeenCalledWith(
        'events.raw',
        expect.objectContaining({
          source: 'mqtt_gateway',
        })
      );
    });

    it('should not throw on publish failure', async () => {
      mockNatsBus.publish.mockRejectedValue(new Error('NATS error'));

      // Should not throw
      await expect(
        (scriptService as any).publishRawEvent('asset-uuid', mockDataSource, {}, 'http')
      ).resolves.not.toThrow();

      expect(mockLogger.warn).toHaveBeenCalledWith(
        expect.stringContaining('Failed to publish raw event')
      );
    });
  });

  describe('publishJsExecutionLog', () => {
    const mockDataSource: DataSource = {
      id: 'test-source',
      orgId: 'org-123',
      pathKey: '0001/0002/',
      name: 'Test Source',
      description: 'Test description',
    };

    const mockResult: ScriptExecutionResult = {
      success: true,
      standardizedPayload: {
        eventType: 'test',
        eventId: 'test-event',
        data: { processed: true },
        metadata: {},
        created: new Date().toISOString(),
      },
      failedAt: null,
      totalExecutionTime: 150,
      error: null,
    };

    it('should publish execution log to events.logs.jsexecutor', async () => {
      await (scriptService as any).publishJsExecutionLog('asset-uuid', mockResult, mockDataSource);

      expect(mockNatsBus.publish).toHaveBeenCalledWith(
        'events.logs.jsexecutor',
        expect.objectContaining({
          threadId: 'asset-uuid',
          orgId: 'org-123',
          pathKey: '0001/0002/',
          name: 'Test Source',
          description: 'Test description',
          execution: expect.objectContaining({
            success: true,
            totalExecutionTime: 150,
          }),
          event: expect.objectContaining({
            eventType: 'test',
            data: { processed: true },
          }),
        })
      );
    });

    it('should not throw on publish failure', async () => {
      mockNatsBus.publish.mockRejectedValue(new Error('NATS error'));

      // Should not throw
      await expect(
        (scriptService as any).publishJsExecutionLog('asset-uuid', mockResult, mockDataSource)
      ).resolves.not.toThrow();

      expect(mockLogger.warn).toHaveBeenCalledWith(
        expect.stringContaining('Failed to publish JS execution log')
      );
    });
  });
});

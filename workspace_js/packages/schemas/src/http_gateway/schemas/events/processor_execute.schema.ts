import { z } from 'zod';
import { IsMongoId, IsString } from '@mapexos/validations';

import { ZodAssetBindSchema } from '../datasources/datasources.schema';

/**
 * Cross-service NATS payload for the SubjectProcessorJSExecute subject
 * (mapexos.processor.js.execute).
 *
 * Mirror of Go: packages/contracts/services/http_gateway/events/types.go::ProcessorExecutePayload
 *
 * Published by: http_gateway (Go).
 * Consumed by: js-executor (TS).
 *
 * Contract is intentionally minimal — js-executor reads pathKey, name, and
 * description from the Asset cache (source of truth) using {orgId}/{assetUUID}
 * as the key. Only orgId + assetBind (cache lookup + asset resolution inputs)
 * are wired here.
 */
export const ZodProcessorExecuteDataSourceSchema = z.object({
  orgId: IsMongoId,
  assetBind: ZodAssetBindSchema,
});

export const ZodProcessorExecutePayloadSchema = z.object({
  sourceType: z.literal('http'),
  dataSource: ZodProcessorExecuteDataSourceSchema,
  event: z.record(IsString, z.unknown()),
  eventTrackerId: IsString,
});

export type ProcessorExecuteDataSource = z.infer<typeof ZodProcessorExecuteDataSourceSchema>;
export type ProcessorExecutePayload = z.infer<typeof ZodProcessorExecutePayloadSchema>;

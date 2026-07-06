import { z } from 'zod';
import {
	ZodOTAPlanStatusSchema,
	ZodOTAExecutionStateSchema,
	ZodOTAFirmwareStatusSchema,
	ZodFirmwareInitRequestSchema,
	ZodFirmwareInitResponseSchema,
	ZodOTARolloutConfigSchema,
	ZodOTAPlanCreateRequestSchema,
	ZodOTAPlanUpdateRequestSchema,
	ZodOTAPlanCountersSchema,
	ZodOTAFirmwareInfoSchema,
	ZodOTAFirmwareDownloadResponseSchema,
	ZodOTAPlanResponseSchema,
	ZodOTAExecutionResponseSchema,
	ZodOTAPlanQuerySchema,
	ZodOTAExecutionQuerySchema,
	ZodOTAPlanIdSchema,
	ZodOTAFirmwareIdSchema,
} from '../../schemas/ota';

/**
 * TypeScript types inferred from Zod schemas for OTA remote firmware update.
 */

export type OTAPlanStatus = z.infer<typeof ZodOTAPlanStatusSchema>;
export type OTAExecutionState = z.infer<typeof ZodOTAExecutionStateSchema>;
export type OTAFirmwareStatus = z.infer<typeof ZodOTAFirmwareStatusSchema>;
export type FirmwareInitRequest = z.infer<typeof ZodFirmwareInitRequestSchema>;
export type FirmwareInitResponse = z.infer<typeof ZodFirmwareInitResponseSchema>;
export type OTARolloutConfig = z.infer<typeof ZodOTARolloutConfigSchema>;
export type OTAPlanCreateRequest = z.infer<typeof ZodOTAPlanCreateRequestSchema>;
export type OTAPlanUpdateRequest = z.infer<typeof ZodOTAPlanUpdateRequestSchema>;
export type OTAPlanCounters = z.infer<typeof ZodOTAPlanCountersSchema>;
export type OTAFirmwareInfo = z.infer<typeof ZodOTAFirmwareInfoSchema>;
export type OTAFirmwareDownloadResponse = z.infer<typeof ZodOTAFirmwareDownloadResponseSchema>;
export type OTAPlanResponse = z.infer<typeof ZodOTAPlanResponseSchema>;
export type OTAExecutionResponse = z.infer<typeof ZodOTAExecutionResponseSchema>;
export type OTAPlanQuery = z.infer<typeof ZodOTAPlanQuerySchema>;
export type OTAExecutionQuery = z.infer<typeof ZodOTAExecutionQuerySchema>;
export type OTAPlanId = z.infer<typeof ZodOTAPlanIdSchema>;
export type OTAFirmwareId = z.infer<typeof ZodOTAFirmwareIdSchema>;

import type {
	FirmwareInitRequest,
	FirmwareInitResponse,
	OTAPlanCreateRequest,
	OTAPlanUpdateRequest,
	OTAPlanQuery,
	OTAPlanResponse,
	OTAExecutionQuery,
	OTAExecutionResponse,
	OTAFirmwareDownloadResponse,
	OTAPlanId,
	OTAFirmwareId,
	PaginatedResponse,
} from '@mapexos/schemas';

import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodFirmwareInitRequestSchema,
	ZodOTAPlanCreateRequestSchema,
	ZodOTAPlanUpdateRequestSchema,
	ZodOTAPlanQuerySchema,
	ZodOTAExecutionQuerySchema,
	ZodOTAPlanIdSchema,
	ZodOTAFirmwareIdSchema,
} from '@mapexos/schemas';

/**
 * Creates the OTA Plans user API (external JWT endpoints on the assets service).
 *
 * Covers the operator surface of the OTA remote firmware update feature:
 * firmware upload lifecycle (init -> presigned PUT -> complete), plan CRUD,
 * and the per-device executions list.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing OTA external API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/ota',
		useAuthJWT: true,
		getToken,
		methods: {
			// FIRMWARE INIT - POST /firmware/init
			// Declares the artifact (target template, version, filename, size,
			// sha256 hex) and returns { firmwareId, uploadUrl } — the client PUTs
			// the .bin straight to object storage via the presigned URL.
			firmwareInit: {
				method: 'POST',
				path: '/firmware/init',
				bodyParams: {} as FirmwareInitRequest,
				bodySchema: ZodFirmwareInitRequestSchema,
				responseType: {} as FirmwareInitResponse,
			},

			// FIRMWARE COMPLETE - POST /firmware/:firmwareId/complete
			// Finalizes the upload: the platform confirms the stored object
			// (size + checksum) and flips the artifact to READY.
			firmwareComplete: {
				method: 'POST',
				path: '/firmware/:firmwareId/complete',
				pathParams: {} as OTAFirmwareId,
				paramSchema: ZodOTAFirmwareIdSchema,
				responseType: {} as { success: boolean },
			},

			// CREATE PLAN - POST /plans
			createPlan: {
				method: 'POST',
				path: '/plans',
				bodyParams: {} as OTAPlanCreateRequest,
				bodySchema: ZodOTAPlanCreateRequestSchema,
				responseType: {} as OTAPlanResponse,
			},

			// LIST PLANS - GET /plans
			listPlans: {
				method: 'GET',
				path: '/plans',
				queryParams: {} as OTAPlanQuery,
				querySchema: ZodOTAPlanQuerySchema,
				responseType: {} as PaginatedResponse<OTAPlanResponse>,
			},

			// GET PLAN - GET /plans/:planId
			getPlan: {
				method: 'GET',
				path: '/plans/:planId',
				pathParams: {} as OTAPlanId,
				paramSchema: ZodOTAPlanIdSchema,
				responseType: {} as OTAPlanResponse,
			},

			// UPDATE PLAN - PATCH /plans/:planId
			updatePlan: {
				method: 'PATCH',
				path: '/plans/:planId',
				pathParams: {} as OTAPlanId,
				paramSchema: ZodOTAPlanIdSchema,
				bodyParams: {} as OTAPlanUpdateRequest,
				bodySchema: ZodOTAPlanUpdateRequestSchema,
				responseType: {} as OTAPlanResponse,
			},

			// DELETE (CANCEL) PLAN - DELETE /plans/:planId
			// Marks the plan CANCELED; the platform's close routine finalizes
			// (timers, stragglers) asynchronously. The firmware artifact is retained.
			deletePlan: {
				method: 'DELETE',
				path: '/plans/:planId',
				pathParams: {} as OTAPlanId,
				paramSchema: ZodOTAPlanIdSchema,
				responseType: {} as { success: boolean },
			},

			// LIST EXECUTIONS - GET /plans/:planId/executions
			listExecutions: {
				method: 'GET',
				path: '/plans/:planId/executions',
				pathParams: {} as OTAPlanId,
				paramSchema: ZodOTAPlanIdSchema,
				queryParams: {} as OTAExecutionQuery,
				querySchema: ZodOTAExecutionQuerySchema,
				responseType: {} as PaginatedResponse<OTAExecutionResponse>,
			},

			// DOWNLOAD FIRMWARE - GET /plans/:planId/firmware/download
			// Mints a short-TTL presigned GET URL for the plan's firmware artifact.
			// Returns 404 when the object is no longer in storage.
			downloadFirmware: {
				method: 'GET',
				path: '/plans/:planId/firmware/download',
				pathParams: {} as OTAPlanId,
				paramSchema: ZodOTAPlanIdSchema,
				responseType: {} as OTAFirmwareDownloadResponse,
			},
		},
	});
}

import type {
	EventsRawQuery,
	EventsRawCursorResult,
	EventsJsExecQuery,
	EventsJsExecCursorResult,
	EventsRouterQuery,
	EventsRouterCursorResult,
	EventsBusinessRuleQuery,
	EventsBusinessRuleCursorResult,
	EventsTriggerQuery,
	EventsTriggerCursorResult,
	EventsStoreQuery,
	EventsStoreCursorResult,
	EventsWorkflowQuery,
	EventsWorkflowExecutionIdParam,
	EventsWorkflowResponse,
	EventsWorkflowCursorResult,
	EventsDLQQuery,
	EventsDLQCursorResult,
	EventsDLQCountsQuery,
	EventsDLQCountsResult,
	AssetConnectivityQuery,
	AssetConnectivityCursorResult,
} from '@mapexos/schemas';

import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodEventsRawQuerySchema,
	ZodEventsJsExecQuerySchema,
	ZodEventsRouterQuerySchema,
	ZodEventsBusinessRuleQuerySchema,
	ZodEventsTriggerQuerySchema,
	ZodEventsStoreQuerySchema,
	ZodEventsWorkflowQuerySchema,
	ZodEventsWorkflowExecutionIdParamSchema,
	ZodEventsDLQQuerySchema,
	ZodEventsDLQCountsQuerySchema,
	ZodAssetConnectivityQuerySchema,
} from '@mapexos/schemas';

/**
 * Creates Events user API with external endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing events external API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/events',
		useAuthJWT: true,
		getToken,
		methods: {
			/**
			 * LIST RAW EVENTS - GET /raw
			 *
			 * Uses cursor-based pagination for efficient querying of large datasets.
			 * Returns items with cursor metadata instead of traditional page/totalItems.
			 *
			 * Query params:
			 *   - cursor: timestamp to start from (optional)
			 *   - direction: "next" (older) or "prev" (newer)
			 *   - limit: max items (default: 20, max: 100)
			 *   - sortAsc: false = DESC (newest first)
			 *   - threadId, source, startTime, endTime: filters
			 */
			listRaw: {
				method: 'GET',
				path: '/raw',
				queryParams: {} as EventsRawQuery,
				querySchema: ZodEventsRawQuerySchema,
				responseType: {} as EventsRawCursorResult,
			},

			/**
			 * LIST JS EXECUTOR EVENTS - GET /jsexec
			 *
			 * Uses cursor-based pagination for efficient querying of large datasets.
			 * Returns items with cursor metadata instead of traditional page/totalItems.
			 *
			 * Query params:
			 *   - cursor: timestamp to start from (optional)
			 *   - direction: "next" (older) or "prev" (newer)
			 *   - limit: max items (default: 20, max: 100)
			 *   - sortAsc: false = DESC (newest first)
			 *   - threadId, success, startTime, endTime: filters
			 */
			listJsExec: {
				method: 'GET',
				path: '/jsexec',
				queryParams: {} as EventsJsExecQuery,
				querySchema: ZodEventsJsExecQuerySchema,
				responseType: {} as EventsJsExecCursorResult,
			},

			/**
			 * LIST ROUTER EVENTS - GET /router
			 *
			 * Uses cursor-based pagination for efficient querying of large datasets.
			 * Returns items with cursor metadata instead of traditional page/totalItems.
			 *
			 * Query params:
			 *   - cursor: timestamp to start from (optional)
			 *   - direction: "next" (older) or "prev" (newer)
			 *   - limit: max items (default: 20, max: 100)
			 *   - sortAsc: false = DESC (newest first)
			 *   - threadId, assetId, routerId, success, publishedCount, startTime, endTime: filters
			 */
			listRouter: {
				method: 'GET',
				path: '/router',
				queryParams: {} as EventsRouterQuery,
				querySchema: ZodEventsRouterQuerySchema,
				responseType: {} as EventsRouterCursorResult,
			},

			/**
			 * LIST BUSINESS RULE EVENTS - GET /businessrule
			 *
			 * Uses cursor-based pagination for efficient querying of large datasets.
			 * Returns items with cursor metadata instead of traditional page/totalItems.
			 *
			 * Query params:
			 *   - cursor: timestamp to start from (optional)
			 *   - direction: "next" (older) or "prev" (newer)
			 *   - limit: max items (default: 20, max: 100)
			 *   - sortAsc: false = DESC (newest first)
			 *   - threadId, ruleId, businessRuleId, matched, startTime, endTime: filters
			 */
			listBusinessRule: {
				method: 'GET',
				path: '/businessrule',
				queryParams: {} as EventsBusinessRuleQuery,
				querySchema: ZodEventsBusinessRuleQuerySchema,
				responseType: {} as EventsBusinessRuleCursorResult,
			},

			/**
			 * LIST TRIGGER EVENTS - GET /trigger
			 *
			 * Uses cursor-based pagination for efficient querying of large datasets.
			 * Returns items with cursor metadata instead of traditional page/totalItems.
			 *
			 * Query params:
			 *   - cursor: timestamp to start from (optional)
			 *   - direction: "next" (older) or "prev" (newer)
			 *   - limit: max items (default: 20, max: 100)
			 *   - sortAsc: false = DESC (newest first)
			 *   - triggerId, triggerType, category, source, success, startTime, endTime: filters
			 */
			listTrigger: {
				method: 'GET',
				path: '/trigger',
				queryParams: {} as EventsTriggerQuery,
				querySchema: ZodEventsTriggerQuerySchema,
				responseType: {} as EventsTriggerCursorResult,
			},

			/**
			 * QUERY STORE EVENTS - POST /store/query
			 *
			 * Uses POST to support EVA dynamic field filters (evaFilters array in body).
			 * Uses cursor-based pagination for efficient querying of processed events.
			 *
			 * Body params (JSON):
			 *   - cursor: timestamp to start from (optional)
			 *   - direction: "next" (older) or "prev" (newer)
			 *   - limit: max items (default: 20, max: 100)
			 *   - sortAsc: false = DESC (newest first)
			 *   - threadId, assetId, templateId, eventType, source, startTime, endTime: filters
			 *   - evaFilters: array of EVA dynamic field filters with operators (optional)
			 */
			queryStore: {
				method: 'POST',
				path: '/store/query',
				bodyParams: {} as EventsStoreQuery,
				bodySchema: ZodEventsStoreQuerySchema,
				responseType: {} as EventsStoreCursorResult,
			},

			/**
			 * LIST WORKFLOW EVENTS - GET /workflow
			 *
			 * Uses cursor-based pagination for efficient querying of workflow execution history.
			 * Returns terminal execution events (completed, failed, cancelled) from ClickHouse cold storage.
			 *
			 * Query params:
			 *   - cursor: timestamp to start from (optional)
			 *   - direction: "next" (older) or "prev" (newer)
			 *   - limit: max items (default: 20, max: 100)
			 *   - sortAsc: false = DESC (newest first)
			 *   - eventTrackerId, workflowUUID, instanceId, definitionId, status, success, startTime, endTime: filters
			 */
			listWorkflow: {
				method: 'GET',
				path: '/workflow',
				queryParams: {} as EventsWorkflowQuery,
				querySchema: ZodEventsWorkflowQuerySchema,
				responseType: {} as EventsWorkflowCursorResult,
			},

			/**
			 * GET WORKFLOW EVENT BY EXECUTION ID - GET /workflow/execution/:executionId
			 *
			 * Retrieves a single workflow execution event by its executionId (MongoDB _id hex).
			 */
			getWorkflowByExecutionId: {
				method: 'GET',
				path: '/workflow/execution/:executionId',
				pathParams: {} as EventsWorkflowExecutionIdParam,
				paramSchema: ZodEventsWorkflowExecutionIdParamSchema,
				responseType: {} as EventsWorkflowResponse,
			},

			// DLQ COUNTS - GET /dlq/counts
			getDLQCounts: {
				method: 'GET',
				path: '/dlq/counts',
				queryParams: {} as EventsDLQCountsQuery,
				querySchema: ZodEventsDLQCountsQuerySchema,
				responseType: {} as EventsDLQCountsResult,
			},

			// LIST DLQ EVENTS - GET /dlq
			listDLQ: {
				method: 'GET',
				path: '/dlq',
				queryParams: {} as EventsDLQQuery,
				querySchema: ZodEventsDLQQuerySchema,
				responseType: {} as EventsDLQCursorResult,
			},

			/**
			 * LIST CONNECTIVITY HISTORY (asset-scoped) - GET /assets/:assetUUID/connectivity_history
			 *
			 * Cursor-paginated offline/online transitions for a single asset.
			 * Backed by the asset_status_history ClickHouse table.
			 */
			listAssetConnectivity: {
				method: 'GET',
				path: '/assets/:assetUUID/connectivity_history',
				pathParams: {} as { assetUUID: string },
				queryParams: {} as AssetConnectivityQuery,
				querySchema: ZodAssetConnectivityQuerySchema,
				responseType: {} as AssetConnectivityCursorResult,
			},

			/**
			 * LIST CONNECTIVITY HISTORY (org-wide) - GET /connectivity_history
			 *
			 * Cursor-paginated offline/online transitions across ALL assets in
			 * the current organization context. Same filters as the asset-scoped
			 * variant plus optional `assetUUID` query filter.
			 */
			listConnectivityHistory: {
				method: 'GET',
				path: '/connectivity_history',
				queryParams: {} as AssetConnectivityQuery,
				querySchema: ZodAssetConnectivityQuerySchema,
				responseType: {} as AssetConnectivityCursorResult,
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

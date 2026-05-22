import type { Request, Response, RequestHandler } from '@mapexos/microservices';
import type { ScriptServicePort } from '@modules/scripts/application/ports';

import { success, notFound, internalError, badRequest } from '@mapexos/microservices';

/**
 * Creates a request handler for fetching scriptTest from a template.
 *
 * This is an INTERNAL endpoint for service-to-service communication.
 * Protected by API Key authentication.
 *
 * Flow:
 * 1. Receives orgId and templateId from path params
 * 2. Fetches template from TieredCache using cache key: {orgId}/{templateId}
 * 3. Returns only the scriptTest field
 *
 * The orgId parameter determines where to look for the template:
 * - "mapexos_public" for system templates (isSystem=true)
 * - Organization's actual ID for private templates
 *
 *
 * @param service - The script service port interface
 * @returns A RequestHandler function that fetches and returns the scriptTest
 */
export function getTemplateScriptTestInternal(service: ScriptServicePort): RequestHandler {
	return async function (req: Request, res: Response) {
		const { orgId, templateId } = req.params;

		if (!orgId) {
			return badRequest(res, 'orgId is required');
		}

		if (!templateId) {
			return badRequest(res, 'templateId is required');
		}

		try {
			// Use getScriptTest which leverages TieredCache (L0 → L1 → L2 → Fallback)
			const scriptTest = await service.getScriptTest(orgId as string, templateId as string);

			if (scriptTest === null) {
				return notFound(res, 'Template or scriptTest not found');
			}

			return success(res, { scriptTest });
		} catch (error: any) {
			return internalError(res, error.message);
		}
	};
}

/**
 * Creates a request handler for generating sample payload from a template.
 *
 * This is an INTERNAL endpoint for Rule Test Runner UI.
 * Protected by API Key authentication.
 *
 * Flow:
 * 1. Receives orgId and templateId from path params
 * 2. Fetches template from TieredCache
 * 3. Extracts scriptTest (raw sample payload)
 * 4. Executes script pipeline: parseScript → validationScript → transformScript
 * 5. Returns the processed/standardized payload
 *
 * The orgId parameter determines where to look for the template:
 * - "mapexos_public" for system templates (isSystem=true)
 * - Organization's actual ID for private templates
 *
 * Used by Rule Test Runner UI to get an editable sample payload.
 * The user can then modify this payload before testing the rule.
 *
 * @param service - The script service port interface
 * @returns A RequestHandler function that generates and returns the sample payload
 */
export function getSamplePayloadInternal(service: ScriptServicePort): RequestHandler {
	return async function (req: Request, res: Response) {
		const { orgId, templateId } = req.params;

		if (!orgId) {
			return badRequest(res, 'orgId is required');
		}

		if (!templateId) {
			return badRequest(res, 'templateId is required');
		}

		try {
			// Execute scripts and get processed sample payload
			const samplePayload = await service.getSamplePayload(orgId as string, templateId as string);

			if (samplePayload === null) {
				return notFound(res, 'Template, scriptTest, or scripts not found');
			}

			return success(res, { samplePayload });
		} catch (error: any) {
			return internalError(res, error.message);
		}
	};
}

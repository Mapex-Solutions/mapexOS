import type { ScriptTest } from '@mapexos/schemas';
import type { AssetScripts } from '@modules/scripts/domain/types';
import type { Request, Response, RequestHandler } from '@mapexos/microservices';
import type { ScriptServicePort } from '@modules/scripts/application/ports';

import { success, internalError, notFound, badRequest, getDTO } from '@mapexos/microservices';

import { PUBLIC_ORG_ID } from '@modules/scripts/application/constants';

/**
 * Creates a request handler for testing script execution.
 *
 * @param service - The script service port interface
 * @returns A RequestHandler function that executes the test and returns the result
 */
export function scriptTest(service: ScriptServicePort): RequestHandler {
	return async function (_: Request, res: Response) {
		const bodyData = getDTO<ScriptTest>(res, 'bodyDTO');

		try {
			const { event, ...scripts } = bodyData;
			const data = await service.scripsTest(event, scripts as AssetScripts);
			return success(res, data);
		} catch (error) {
			return internalError(res, error.message);
		}
	};
}

/**
 * Creates a request handler for fetching scriptTest from a template.
 *
 * Retrieves the scriptTest payload from a template via TieredCache.
 * Used by the Rule Test Runner UI to populate the test payload based on selected template.
 *
 * @param service - The script service port interface
 * @returns A RequestHandler function that fetches and returns the scriptTest
 */
export function getTemplateScriptTest(service: ScriptServicePort): RequestHandler {
	return async function (req: Request, res: Response) {
		const { templateId } = req.params;
		const { orgId, isSystem } = req.query;

		if (!templateId) {
			return badRequest(res, 'templateId is required');
		}

		try {
			// For system templates, use PUBLIC_ORG_ID
			// For org templates, use the provided orgId
			const templateOrgId = isSystem === 'true' ? PUBLIC_ORG_ID : (orgId as string);

			if (!templateOrgId) {
				return badRequest(res, 'orgId is required for non-system templates');
			}

			const scriptTest = await service.getScriptTest(templateOrgId, templateId as string);

			if (scriptTest === null) {
				return notFound(res, 'Template or scriptTest not found');
			}

			return success(res, { scriptTest });
		} catch (error) {
			return internalError(res, error.message);
		}
	};
}

/**
 * Creates a request handler for generating sample payload from a template.
 *
 * This endpoint:
 * 1. Fetches the template from TieredCache
 * 2. Extracts scriptTest (raw sample payload)
 * 3. Executes script pipeline: parseScript → validationScript → transformScript
 * 4. Returns the processed/standardized payload
 *
 * Used by Rule Test Runner UI to get an editable sample payload.
 * The user can then modify this payload before testing the rule.
 *
 * Query params:
 * - orgId: Organization ID (required for non-system templates)
 * - isSystem: "true" if system template (uses PUBLIC_ORG_ID)
 *
 * @param service - The script service port interface
 * @returns A RequestHandler function that generates and returns the sample payload
 */
export function getSamplePayload(service: ScriptServicePort): RequestHandler {
	return async function (req: Request, res: Response) {
		const { templateId } = req.params;
		const { orgId, isSystem } = req.query;

		if (!templateId) {
			return badRequest(res, 'templateId is required');
		}

		try {
			// For system templates, use PUBLIC_ORG_ID
			// For org templates, use the provided orgId
			const templateOrgId = isSystem === 'true' ? PUBLIC_ORG_ID : (orgId as string);

			if (!templateOrgId) {
				return badRequest(res, 'orgId is required for non-system templates');
			}

			// Execute scripts and get processed sample payload
			const samplePayload = await service.getSamplePayload(templateOrgId, templateId as string);

			if (samplePayload === null) {
				return notFound(res, 'Template, scriptTest, or scripts not found');
			}

			return success(res, { samplePayload });
		} catch (error: any) {
			return internalError(res, error.message);
		}
	};
}
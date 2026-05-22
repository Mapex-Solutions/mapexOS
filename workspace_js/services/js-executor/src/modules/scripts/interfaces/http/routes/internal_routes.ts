import type { ScriptServicePort } from '@modules/scripts/application/ports';

import { Router } from '@mapexos/microservices';
import {
	getTemplateScriptTestInternal,
	getSamplePayloadInternal,
} from '@modules/scripts/interfaces/http/handlers/internal_handler';

/**
 * Registers internal API routes for service-to-service communication.
 *
 * These routes are protected by API Key authentication and are used internally
 *
 * Routes:
 * - GET /:orgId/:templateId/script_test - Fetch raw scriptTest from template
 *   - orgId: "mapexos_public" for system templates, or org's ID for private templates
 *   - templateId: The MongoDB ObjectID of the template
 *   - Returns: { scriptTest: object } - The raw scriptTest payload
 *
 * - GET /:orgId/:templateId/sample_payload - Generate processed sample payload
 *   - orgId: "mapexos_public" for system templates, or org's ID for private templates
 *   - templateId: The MongoDB ObjectID of the template
 *   - Returns: { samplePayload: object } - The processed/standardized payload
 *   - Used by Rule Test Runner UI to get an editable sample payload
 *
 * @param service - The script service port interface
 * @returns {Router} A Router instance containing the registered internal routes.
 */
export function registerInternalRoutes(service: ScriptServicePort): Router {
	const group = Router();

	/** Get raw scriptTest from template using orgId + templateId (TieredCache key format) */
	group.get('/:orgId/:templateId/script_test', getTemplateScriptTestInternal(service));

	/** Get processed sample payload by executing template scripts against scriptTest */
	group.get('/:orgId/:templateId/sample_payload', getSamplePayloadInternal(service));

	return group;
}

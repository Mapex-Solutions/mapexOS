import type { Validation } from '@mapexos/microservices';
import type { ScriptServicePort } from '@modules/scripts/application/ports';

import { Router } from '@mapexos/microservices';
import { validationMiddleware, newValidation } from '@mapexos/microservices';
import { ScriptTestDto } from '@modules/scripts/application/dtos';
import { scriptTest, getTemplateScriptTest, getSamplePayload } from '@modules/scripts/interfaces/http';

/**
 * Registers and returns a group of HTTP routes for the application.
 *
 * Routes:
 * - POST /test - Test script execution with provided payload and scripts
 * - GET /templates/:templateId/script_test - Get raw scriptTest from template
 * - GET /templates/:templateId/sample_payload - Get processed sample payload (scripts executed)
 *
 * @param service - The script service port interface
 * @returns {Router} A Router instance containing the registered routes.
 */
export function registerRoutes(service: ScriptServicePort): Router {
	/** Router instance */
	const group = Router();

	/** Test script execution with provided payload and scripts */
	const validationTestsDto = newValidation(ScriptTestDto, null, null);
	group.post('/test', validationMiddleware(validationTestsDto), scriptTest(service));

	/** Get raw scriptTest from template via TieredCache */
	group.get('/templates/:templateId/script_test', getTemplateScriptTest(service));

	/** Get processed sample payload by executing template scripts against scriptTest */
	group.get('/templates/:templateId/sample_payload', getSamplePayload(service));

	/** Return the router group */
	return group;
}

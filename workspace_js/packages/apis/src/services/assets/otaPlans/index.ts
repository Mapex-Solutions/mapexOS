import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './user.api';
import { sha256OfFirmware, uploadFirmwareBinary } from './upload.api';

/**
 * Creates the OTA Plans module API (assets service, /api/v1/ota).
 *
 * User methods (JWT auth) are spread to the root level; the presigned-upload
 * helpers (direct browser -> object storage, no platform auth) ride along so
 * the whole firmware flow lives behind one module.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with OTA plan/firmware/execution API methods + upload helpers
 *
 * @example
 * const otaPlans = createOtaPlansApi(http, getToken);
 *
 * // Firmware flow: declare -> PUT to storage -> finalize
 * const { hex, base64 } = await otaPlans.sha256OfFirmware(file);
 * const { firmwareId, uploadUrl } = await otaPlans.firmwareInit({ targetTemplateId, version, filename, size, sha256: hex });
 * await otaPlans.uploadFirmwareBinary(uploadUrl, file, base64);
 * await otaPlans.firmwareComplete({ firmwareId });
 *
 * // Plans
 * await otaPlans.createPlan({ firmwareId, sourceTemplateId, assetIds, startAt, maxTime, name, rolloutConfig });
 * await otaPlans.listPlans({ page: 1, perPage: 20 });
 * await otaPlans.listExecutions({ planId }, { state: 'DOWNLOADING' });
 */
export function createOtaPlansApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),
		sha256OfFirmware,
		uploadFirmwareBinary,
	};
}

export type OtaPlansApiMethods = ReturnType<typeof createOtaPlansApi>;

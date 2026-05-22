import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { membershipApi } from './membership.api';

/**
 * Creates Membership module API.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with membership API methods
 *
 * @example
 * const membershipApi = createMembershipApi(http, getToken);
 *
 * // Membership methods (JWT)
 * await membershipApi.list({ page: 1, perPage: 20 });
 * await membershipApi.create({ assigneeType: 'user', assigneeId: '...', ... });
 * await membershipApi.getById({ membershipId: '...' });
 * await membershipApi.update({ membershipId: '...' }, { roleIds: [...] });
 * await membershipApi.delete({ membershipId: '...' });
 */
export function createMembershipApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...membershipApi(http, getToken),
	};
}

export type MembershipApiMethods = ReturnType<typeof createMembershipApi>;

import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './user.api';

/**
 * Creates Organization module API.
 *
 * Organization methods (JWT auth) are spread to the root level.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with organization API methods
 *
 * @example
 * const organizationApi = createOrganizationApi(http, getToken);
 *
 * // Organization methods (JWT)
 * await organizationApi.list({ page: 1, perPage: 20 });
 * await organizationApi.tree({ limit: 50, sortAsc: true });
 * await organizationApi.create({ name: 'New Org', type: 'customer', ... });
 * await organizationApi.getById({ organizationId: '...' });
 * await organizationApi.update({ organizationId: '...' }, { name: 'Updated Org' });
 * await organizationApi.delete({ organizationId: '...' });
 */
export function createOrganizationApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),
	};
}

export type OrganizationApiMethods = ReturnType<typeof createOrganizationApi>;

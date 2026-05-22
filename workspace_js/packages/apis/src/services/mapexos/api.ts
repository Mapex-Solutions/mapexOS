import { isEmpty } from 'lodash';
import { ApiConfig, ApiInterceptors, SessionsConfig, GetToken } from '@src/common';
import { createHttp } from '@src/tools';

import { createAuthApi } from './auth';
import { createUserApi } from './users';
import { createGroupApi } from './groups';
import { createOrganizationApi } from './organizations';
import { createRoleApi } from './roles';
import { createListApi } from './lists';
import { createMembershipApi } from './memberships';
import { createOnboardingApi } from './onboarding';

export function createMapexOSApi(config: ApiConfig, sessionsConfig?: SessionsConfig) {

	/** Global /Local interceptors for the mapexos API */
	const interceptors = !isEmpty(sessionsConfig?.interceptors) ? sessionsConfig?.interceptors : config?.interceptors || {} as ApiInterceptors;
	const getToken: GetToken | undefined = sessionsConfig?.getToken;

	/** HTTP instance */
	const http = createHttp(config, interceptors);

	/** Modules within the service */
	const auth = createAuthApi(http, getToken, sessionsConfig);
	const users = createUserApi(http, getToken);
	const groups = createGroupApi(http, getToken);
	const organizations = createOrganizationApi(http, getToken);
	const roles = createRoleApi(http, getToken);
	const lists = createListApi(http, getToken);
	const memberships = createMembershipApi(http, getToken);
	const onboarding = createOnboardingApi(http, getToken);

	return {
		http,
		auth,
		users,
		groups,
		organizations,
		roles,
		lists,
		memberships,
		onboarding,
	};
}

/**
 * Interface for API initialization configuration.
 */
export type MapexOSApi = ReturnType<typeof createMapexOSApi>;
import type {
	CreateUserWithMemberships,
	UpdateUserWithAccess,
	UpdateUserWithAccessParams,
	UserOnboardingResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodCreateUserWithMembershipsSchema,
	ZodUpdateUserWithAccessParamsSchema,
	ZodUpdateUserWithAccessSchema,
} from '@mapexos/schemas';

/**
 * Creates Onboarding API for orchestrated user creation with memberships.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing onboarding API methods
 */
export function onboardingApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/onboarding',
		useAuthJWT: true,
		getToken,
		methods: {
			/**
			 * Create user with memberships in a single atomic operation
			 * POST /api/v1/onboarding/users
			 *
			 * This endpoint creates a user AND their memberships (direct roles or groups)
			 * ensuring the user is immediately accessible in the organization.
			 */
			createUserWithMemberships: {
				method: 'POST',
				path: '/users',
				bodyParams: {} as CreateUserWithMemberships,
				bodySchema: ZodCreateUserWithMembershipsSchema,
				responseType: {} as UserOnboardingResponse,
			},

			/**
			 * Update user with access configuration in a single atomic operation
			 * PATCH /api/v1/onboarding/users/:userId
			 *
			 * This endpoint updates user data AND replaces their access configuration
			 * (memberships/groups) ensuring atomic updates within the organization context.
			 *
			 * Backend logic:
			 * 1. Updates user data
			 * 2. Removes old memberships/groups in current org
			 * 3. Creates new memberships/groups as specified
			 * 4. All within a MongoDB transaction
			 */
			updateUserWithAccess: {
				method: 'PATCH',
				path: '/users/:userId',
				pathParams: {} as UpdateUserWithAccessParams,
				paramSchema: ZodUpdateUserWithAccessParamsSchema,
				bodyParams: {} as UpdateUserWithAccess,
				bodySchema: ZodUpdateUserWithAccessSchema,
				responseType: {} as UserOnboardingResponse,
			},
		},
	});
}

export type OnboardingApiMethods = ReturnType<typeof onboardingApi>;

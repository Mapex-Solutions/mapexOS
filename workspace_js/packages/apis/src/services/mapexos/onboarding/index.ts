import type { GetToken } from '@src/common';
import type { AxiosInstance } from 'axios';
import { onboardingApi } from './onboarding.api';

export function createOnboardingApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return onboardingApi(http, getToken);
}

export * from './onboarding.api';

import type { AuthState } from './types';
import { apis } from '@src/services/mapex';
import { useOrganizationStore } from '@stores/organization';
import { usePermissionStore } from '@stores/permission';
import storage from '@utils/storages';
import { useLogger } from '@composables/useLogger';

const logger = useLogger('AuthStore');

export const actions = {

	/**
	 * Authenticates a user with the provided email and password, and updates the store with the authentication tokens and user information.
	 * After successful authentication, initializes the organization store with user's organization coverage.
	 *
	 * @param email - The email address of the user attempting to log in.
	 * @param password - The password associated with the user's email.
	 * @param keepConnected - A boolean indicating whether the user should remain logged in across sessions.
	 * @returns A promise that resolves when the login process is complete and the store is updated.
	 * @throws Will throw an error if the login process fails.
	 */
	async login(email: string, password: string, keepConnected: boolean) {
		const store = this as AuthState & typeof actions;

		try {
			// 1. Authenticate user
			const data = await apis.mapexOS?.auth.login({
				email,
				password,
				keepConnected,
			});

			store.accessToken = data.access_token;
			store.refreshToken = data.refresh_token;
			store.user = data.user;
			store.keepConnected = keepConnected;

			// 2. Persist tokens based on keepConnected preference
			store.persistTokens();

			// 3. Initialize organization coverage after successful login
			const orgStore = useOrganizationStore();
			await orgStore.initializeAfterLogin();

		} catch (error) {
			logger.error('Error logging in', error);
			throw error;
		}
	},

	/**
	 * Logs out the current user by clearing authentication tokens and user information from both the store and storage.
	 * Also clears the organization store coverage data.
	 */
	logout() {
		const store = this as AuthState & typeof actions;

		// Clear store state
		store.accessToken = '';
		store.refreshToken = '';
		store.user = null;
		store.keepConnected = false;

		// Clear from both storages to be safe
		storage.local.remove('auth_tokens');
		storage.session.remove('auth_tokens');

		// Clear organization store
		const orgStore = useOrganizationStore();
		orgStore.clearCoverage();

		// Clear permission store
		const permStore = usePermissionStore();
		permStore.clearPermissions();
	},

	/**
	 * Persists authentication tokens to storage based on keepConnected preference.
	 * If keepConnected is true, tokens are stored in localStorage for persistence across sessions.
	 * If false, tokens are stored in sessionStorage for the current session only.
	 */
	persistTokens() {
		const store = this as AuthState & typeof actions;

		const authData = {
			accessToken: store.accessToken,
			refreshToken: store.refreshToken,
			user: store.user,
			keepConnected: store.keepConnected,
		};

		if (store.keepConnected) {
			// Persistent storage - survives browser restart
			storage.local.set('auth_tokens', authData);
		} else {
			// Session storage - cleared on browser close
			storage.session.set('auth_tokens', authData);
		}
	},

	/**
	 * Hydrates the auth store from storage on app load or refresh.
	 * Checks both localStorage and sessionStorage for stored authentication data.
	 * Returns true if valid auth data was found and loaded.
	 */
	hydrateFromStorage(): boolean {
		const store = this as AuthState & typeof actions;

		// Try localStorage first (keepConnected users)
		let authData = storage.local.get('auth_tokens');

		// If not in local, try session storage
		if (!authData) {
			authData = storage.session.get('auth_tokens');
		}

		// If found valid auth data, restore to store
		if (authData?.accessToken && authData?.user) {
			store.accessToken = authData.accessToken;
			store.refreshToken = authData.refreshToken;
			store.user = authData.user;
			store.keepConnected = authData.keepConnected || false;
			return true;
		}

		return false;
	},

	/**
	 * Updates the access and refresh tokens in the store and persists them.
	 * Used after token refresh.
	 */
	updateTokens(accessToken: string, refreshToken: string) {
		const store = this as AuthState & typeof actions;

		store.accessToken = accessToken;
		store.refreshToken = refreshToken;

		// Re-persist with same storage strategy
		store.persistTokens();
	},

	setEmail() {
	},

	setPassword() {
	},

	setKeepConnected() {
	},
};

/**
 * Authentication Handler Service
 *
 * Decoupled service for handling authentication-related operations.
 * This service is designed to work independently of any authentication
 * implementation, making it easy to integrate with any auth system later.
 *
 * Future integration points:
 * - JWT token management
 * - Session management
 * - OAuth/SSO integration
 * - Refresh token handling
 * - Multi-factor authentication
 */

import { useRouter } from 'vue-router';
import { useAuthStore } from '@stores/auth';

/**
 * Clears all authentication data from stores and browser storage
 * Delegates to authStore.logout() which handles Pinia state, token storage,
 * and organization store cleanup. Also clears remaining browser storage.
 */
export function clearAuthData(): void {
  // Clear auth store (Pinia state + token storage + org store)
  const authStore = useAuthStore();
  authStore.logout();

  // Clear any remaining browser storage
  localStorage.clear();
  sessionStorage.clear();
}

/**
 * Handles user logout
 * - Clears all auth data (Pinia stores + browser storage)
 * - Redirects to login page
 *
 * @param router - Vue Router instance for navigation
 */
export function handleLogout(router: ReturnType<typeof useRouter>): void {
  // Clear all authentication data
  clearAuthData();

  // Redirect to login page
  void router.push('/');
}

/**
 * Composable hook for authentication operations
 * Provides a clean interface for components to use auth functionality
 */
export function useAuth() {
  const router = useRouter();

  return {
    /**
     * Logout user and redirect
     */
    logout: () => handleLogout(router),

    /**
     * Clear auth data without redirect
     */
    clearData: () => clearAuthData(),

    // Future methods to be implemented:
    // login: (credentials) => {},
    // refreshToken: () => {},
    // checkAuth: () => {},
    // getUser: () => {},
  };
}

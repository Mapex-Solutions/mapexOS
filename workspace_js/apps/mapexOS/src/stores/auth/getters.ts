import type { AuthState } from './types'

export const getters = {
  isAuthenticated: (s: AuthState) => Boolean(s.accessToken && s.user),
}
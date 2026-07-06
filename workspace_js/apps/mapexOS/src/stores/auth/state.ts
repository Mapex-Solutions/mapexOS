// src/stores/auth/state.ts
import type { AuthState } from './types'

export const state = (): AuthState => ({
  keepConnected: true,

  accessToken: '',
  refreshToken: '',
  user: null,

  loading: false,
  error: null,
})

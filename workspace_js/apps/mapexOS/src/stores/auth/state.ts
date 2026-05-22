// src/stores/auth/state.ts
import type { AuthState } from './types'

export const state = (): AuthState => ({
  email: "admin@mapex.global",    
  password: "mapex123",
  keepConnected: true,

  accessToken: '',
  refreshToken: '',
  user: null,

  loading: false,
  error: null,
})

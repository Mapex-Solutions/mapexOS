export type User = {
	authProvider: {
		internal: string,
		metadata: any,
	};
  changePasswordNextLogin: boolean;
  created: string; // ISO date string
  email: string;
  enabled: boolean;
  startTour: boolean;
  firstName: string;
  id: string;
  lastName: string;
  phone: string;
  updated: string; // ISO date string
}

export interface AuthState {
	// session persistence preference
	keepConnected: boolean;

	// session
	accessToken: string;
	refreshToken: string;
	user: any;

	// ui
	loading: boolean;
	error: string | null;
}
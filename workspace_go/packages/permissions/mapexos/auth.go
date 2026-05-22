package permissions

// Auth Permissions
const (
	// AuthLogin - Permission to login
	AuthLogin = "auth.login"

	// AuthLogout - Permission to logout
	AuthLogout = "auth.logout"

	// AuthRefreshToken - Permission to refresh authentication token
	AuthRefreshToken = "auth.refresh"

	// AuthChangePassword - Permission to change password
	AuthChangePassword = "auth.changepassword"

	// AuthResetPassword - Permission to reset password
	AuthResetPassword = "auth.resetpassword"

	// AuthAll - Wildcard permission for all auth operations
	AuthAll = "auth.*"
)

package auth

type LoginDTO struct {
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=8"`
	KeepConnected bool   `json:"keepConnected"`
}

type BuildAuthorizationCacheRequest struct {
	UserID string `json:"userId" validate:"required"`
	OrgID  string `json:"orgId"` // Optional - when empty, builds root permissions
}

type BuildCoverageCacheRequest struct {
	UserID string `json:"userId" validate:"required"`
}

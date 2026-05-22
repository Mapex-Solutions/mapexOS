package constants

import "time"

// Token TTLs centralized so Login and RefreshToken share the same values
// without drifting between the two orchestrations.
const (
	AccessTokenTTL  = 30 * time.Minute
	RefreshTokenTTL = 7 * 24 * 60 * time.Minute
)

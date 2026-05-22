package constants

// RevocationReason is the domain enum for why a cert was revoked.
// String values mirror the cross-service contract for storage consistency.
type RevocationReason string

const (
	ReasonReplaced   RevocationReason = "replaced"
	ReasonUserAction RevocationReason = "user_action"
)

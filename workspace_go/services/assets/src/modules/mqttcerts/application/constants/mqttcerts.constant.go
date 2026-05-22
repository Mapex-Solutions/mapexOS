package constants

import "time"

const (
	DefaultDeviceCertTTLDays  = 365
	CABootstrapInitialTimeout = 5 * time.Second
	CABootstrapBackoffMin     = 1 * time.Second
	CABootstrapBackoffMax     = 30 * time.Second

	// CertTTL bounds applied when the operator supplies a per-asset
	// `certTTL` override. Min = 1 day so an accidental zero never
	// produces an instantly-expired cert. Max = 10 years; longer than
	// that the CA's own validity becomes the limiting factor anyway,
	// so refusing here is the safer path than emitting a cert that
	// outlives its issuer.
	MinDeviceCertTTLDays = 1
	MaxDeviceCertTTLDays = 3650
)

// CertTTL unit identifiers. Mirrors the platform contract's enum
// (`oneof=day week month year`) — kept here so the signer and the
// validator import a single source of truth.
const (
	CertTTLUnitDay   = "day"
	CertTTLUnitWeek  = "week"
	CertTTLUnitMonth = "month"
	CertTTLUnitYear  = "year"
)

// CertTTLUnitToDays returns the day-count multiplier for the supplied
// unit token. Unknown / empty unit returns 0 so callers can detect
// validation failure without an out-of-band error code. Month + year
// are deliberate approximations (30 + 365) — cert resolution is day-
// scale and calendar-accurate math would buy nothing.
func CertTTLUnitToDays(unit string) int {
	switch unit {
	case CertTTLUnitDay:
		return 1
	case CertTTLUnitWeek:
		return 7
	case CertTTLUnitMonth:
		return 30
	case CertTTLUnitYear:
		return 365
	default:
		return 0
	}
}

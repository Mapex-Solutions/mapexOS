package constants

// CAKind discriminates which CA in the hierarchy this record represents.
// Domain-owned vocabulary — NOT imported from packages/contracts.
type CAKind string

const (
	CAKindRoot         CAKind = "root-ca"
	CAKindIntermediate CAKind = "intermediate-ca"
)

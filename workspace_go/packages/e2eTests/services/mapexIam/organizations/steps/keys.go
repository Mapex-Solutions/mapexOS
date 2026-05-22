package steps

// Bag keys this package writes. Other packages reading these keys import
// the constants from here.
const (
	// BagKeyOrgID is the organization id every saga step downstream
	// reads to bind resources (asset templates, route groups) to the
	// active org context. The runner also reads it to set
	// X-Org-Context on every per-service client.
	BagKeyOrgID = "iam.organizationID"

	// BagKeyOrgPathKey is the materialised dot-path of the org tree
	// (e.g. "0000000000000000000aa001.0000000000000000000aa007") —
	// useful for assertions that need to verify ancestry without
	// walking the org graph by hand.
	BagKeyOrgPathKey = "iam.organizationPathKey"
)

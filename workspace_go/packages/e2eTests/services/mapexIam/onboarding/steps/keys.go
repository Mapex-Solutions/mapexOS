package steps

// Bag keys this package writes. Other packages reading these keys import
// the constants from here.
const (
	// BagKeyUserID is the freshly-created user id returned by the
	// onboarding orchestrator endpoint.
	BagKeyUserID = "iam.userID"

	// BagKeyUserEmail is the email used during creation. AuthenticateUser
	// reads it to know which credentials to send when the saga later
	// switches from the seed admin to the saga test user.
	BagKeyUserEmail = "iam.userEmail"

	// BagKeyGroupID is the first group id returned by the onboarding
	// response. Saga journeys that exercise group-scoped behaviour read
	// it; cleanup is covered by the org cascade.
	BagKeyGroupID = "iam.groupID"
)

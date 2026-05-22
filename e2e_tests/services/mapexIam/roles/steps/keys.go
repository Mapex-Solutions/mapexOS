package steps

// Bag keys this package writes. Other packages reading these keys import
// the constants from here.
const (
	// BagKeyRoleID is the role id returned by CreateRole. Onboarding
	// reads it when binding a role to a freshly-created user.
	BagKeyRoleID = "iam.roleID"
)

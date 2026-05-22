package steps

// Bag keys this package writes. Other packages that read from these keys
// import the constants instead of using string literals so renames are
// caught at compile time and grep finds every producer + consumer.
const (
	// BagKeyUserJWT is the bearer token published by SeedAdminLogin or
	// AuthenticateUser. Asserts and downstream steps use it to inspect
	// claims; the runner already attaches it to every per-service HTTP
	// client via ClientSet.SetBearer.
	BagKeyUserJWT = "iam.userJWT"
)

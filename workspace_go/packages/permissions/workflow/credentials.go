package permissions

// Credential Permissions
const (
	// CredentialList - Permission to list all credentials
	CredentialList = "credentials.list"

	// CredentialCreate - Permission to create a new credential
	CredentialCreate = "credentials.create"

	// CredentialRead - Permission to read a specific credential
	CredentialRead = "credentials.read"

	// CredentialUpdate - Permission to update a credential
	CredentialUpdate = "credentials.update"

	// CredentialDelete - Permission to delete a credential
	CredentialDelete = "credentials.delete"

	// CredentialTest - Permission to test a credential
	CredentialTest = "credentials.test"

	// CredentialAll - Wildcard permission for all credential operations
	CredentialAll = "credentials.*"
)

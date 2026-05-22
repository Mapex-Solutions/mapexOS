package ports

import "context"

// VaultPort provides credential decryption for plugin execution.
// Implemented by the VaultClient (HTTP call to mapexVault MS).
type VaultPort interface {
	// DecryptCredential returns plaintext credential data by ID.
	DecryptCredential(ctx context.Context, id string) (map[string]interface{}, error)
}

package ports

import "context"

// LorawanKEKClientPort wraps the HTTP call to GET /internal/kek/:context on
// mapexVault. Returns the decrypted KEK as a 64-hex string. Used by the asset
// service OnMount to load the KEK once at boot.
type LorawanKEKClientPort interface {
	FetchKEK(ctx context.Context, keyContext string) (string, error)
}

package services

import (
	"crypto/sha256"
	"encoding/hex"

	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// verifyBundleChecksum hard-verifies that rawBytes hashes to declaredSha256 (a
// hex-encoded SHA-256 over the artifact's exact bytes as served — never a
// re-serialized form, matching the content-hashing pattern npm, pip, Docker,
// Homebrew and Terraform all use). On a mismatch it returns the structured,
// UI-facing rejection error (HTTP 422) so the install creates nothing — the
// standing integrity standard for every marketplace install flow here.
func (s *AssetTemplateService) verifyBundleChecksum(rawBytes []byte, declaredSha256 string) error {
	sum := sha256.Sum256(rawBytes)
	computed := hex.EncodeToString(sum[:])
	if computed != declaredSha256 {
		return &customErrors.ServerCustomError{
			Code:   httpStatus.UNPROCESSABLE_ENTITY,
			Errors: []string{"CHECKSUM_MISMATCH: The template downloaded from the marketplace does not match its published checksum and was not installed."},
		}
	}
	return nil
}

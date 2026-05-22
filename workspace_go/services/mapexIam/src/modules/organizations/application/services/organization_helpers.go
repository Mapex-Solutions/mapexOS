package services

import (
	"strings"

	"mapexIam/src/modules/organizations/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// toBase36 converts a decimal number to Base36 (0-9, A-Z) with zero-padding.
//
// Parameters:
//   - num: The decimal number to convert.
//   - width: The desired string width with zero-padding.
//
// Returns:
//   - A Base36 string representation (e.g., "00001A", "0001", "01").
func toBase36(num int, width int) string {
	const base36Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if num == 0 {
		return strings.Repeat("0", width)
	}

	var result []byte
	for num > 0 {
		remainder := num % 36
		result = append([]byte{base36Chars[remainder]}, result...)
		num = num / 36
	}

	// Add leading zeros if needed
	for len(result) < width {
		result = append([]byte{'0'}, result...)
	}

	return string(result)
}

// buildOrganizationCode generates a unique code for the organization based on the parent's childCount.
// The code uses Base36 encoding with variable width based on organization type:
//   - Vendor/Customer: 6 chars (36^6 = 2.1 billion)
//   - Site/Building: 4 chars (36^4 = 1.6 million)
//   - Floor/Zone: 3 chars (36^3 = 46K)
//
// Parameters:
//   - childCount: The current number of children the parent has (will be incremented by 1).
//   - orgType: The type of organization being created.
//
// Returns:
//   - A Base36 string with appropriate width (e.g., "000001", "0001", "01").
func buildOrganizationCode(childCount int, orgType string) string {
	newChildNumber := childCount + 1

	var width int
	switch orgType {
	case "vendor", "customer":
		width = 6
	case "site", "building":
		width = 4
	case "floor", "zone":
		width = 3
	default:
		width = 6 // Default to 6 for unknown types
	}

	return toBase36(newChildNumber, width)
}

// buildPathKey constructs the full hierarchical path for the organization.
// If parent exists, it appends the new code to the parent's pathKey.
// If no parent exists (vendor), it returns just the code.
//
// Parameters:
//   - parentPathKey: The parent organization's pathKey (empty string if no parent).
//   - code: The organization's local code.
//
// Returns:
//   - The full pathKey string (e.g., "000001/000002/000003").
func buildPathKey(parentPathKey, code string) string {
	if parentPathKey == "" {
		return code
	}
	return parentPathKey + "/" + code
}

// calculateDepth determines the organization's depth in the hierarchy tree.
// Depth indicates the level: 0=vendor, 1=customer, 2=site, 3=building, etc.
//
// Parameters:
//   - parent: Pointer to the parent organization entity (nil if no parent).
//
// Returns:
//   - The depth as an integer (0 if no parent, parent.Depth + 1 otherwise).
func calculateDepth(parent *entities.Organization) int {
	if parent == nil {
		return 0 // Root level (vendor)
	}
	return parent.Depth + 1
}

// determineCustomerID determines the customerID for the new organization.
// - If type is "customer", the customerID is the organization's own ID.
// - Otherwise, it inherits the customerID from the parent.
//
// Parameters:
//   - orgType: The type of the organization being created.
//   - orgID: The ID of the organization being created.
//   - parent: Pointer to the parent organization entity (nil if no parent).
//
// Returns:
//   - A pointer to the customerID (ObjectId), or nil if no customer context exists.
func determineCustomerID(orgType string, orgID model.ObjectId, parent *entities.Organization) *model.ObjectId {
	if orgType == "customer" {
		return &orgID
	}
	if parent != nil {
		return parent.CustomerID
	}
	return nil
}

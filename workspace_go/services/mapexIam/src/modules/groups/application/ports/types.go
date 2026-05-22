package ports

/* GroupBasicInfo */

// GroupBasicInfo is a lightweight DTO for cross-domain queries.
// Contains only the fields other domains need, avoiding entity leakage.
type GroupBasicInfo struct {
	ID          string
	Name        string
	Description *string
}

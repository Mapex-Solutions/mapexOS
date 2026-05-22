package retention

import "fmt"

// RetentionPolicyLimits defines the minimum, maximum, and default retention days for a table
type RetentionPolicyLimits struct {
	DefaultDays uint16
	MinDays     uint16
	MaxDays     uint16
	Name        string // Human-readable name for the table
}

// RetentionPoliciesLimits defines retention limits for all ClickHouse event tables
// All values are in DAYS (uint16)
//
// Rationale:
//   - eventsRaw, eventsRouter: High volume debug logs (1-3 days max) - exported to Parquet for cold storage
//   - eventsJsExecutor: Script execution stats (3 days) - operational debugging
//   - events: Core IoT analytics (7-365 days) - business-critical data
//   - eventsBusinessRule: Audit trail for business rules (1-15 days)
//   - eventsBusinessRuleCondition: Detailed condition logs (1-365 days)
//   - eventsAudit: Compliance requirement (minimum 365 days, up to 7 years)
//   - eventsNotifications: User-facing notification history (30-365 days)
var RetentionPoliciesLimits = map[string]RetentionPolicyLimits{
	"eventsRaw": {
		DefaultDays: 1,
		MinDays:     1,
		MaxDays:     3,
		Name:        "Events Raw",
	},
	"eventsJsExecutor": {
		DefaultDays: 3,
		MinDays:     1,
		MaxDays:     3,
		Name:        "Events JS Executor",
	},
	"eventsRouter": {
		DefaultDays: 1,
		MinDays:     1,
		MaxDays:     3,
		Name:        "Events Router",
	},
	"events": {
		DefaultDays: 90,
		MinDays:     7,
		MaxDays:     365,
		Name:        "Events",
	},
	"eventsBusinessRule": {
		DefaultDays: 1,
		MinDays:     1,
		MaxDays:     15,
		Name:        "Events Business Rule",
	},
	"eventsBusinessRuleCondition": {
		DefaultDays: 1,
		MinDays:     1,
		MaxDays:     365,
		Name:        "Events Business Rule Condition",
	},
	"eventsAudit": {
		DefaultDays: 365,
		MinDays:     365,  // Compliance requirement: minimum 1 year
		MaxDays:     2555, // 7 years (legal/fiscal retention)
		Name:        "Events Audit",
	},
	"eventsNotifications": {
		DefaultDays: 30,
		MinDays:     1,
		MaxDays:     365,
		Name:        "Events Notifications",
	},
	"eventsWorkflow": {
		DefaultDays: 7,
		MinDays:     1,
		MaxDays:     365,
		Name:        "Events Workflow",
	},
	"asset_status_history": {
		DefaultDays: 7,
		MinDays:     1,
		MaxDays:     90,
		Name:        "Asset Connectivity History",
	},
}

// ValidateRetentionPolicy validates if the retention days are within allowed limits for a specific table
func ValidateRetentionPolicy(table string, days uint16) error {
	limits, exists := RetentionPoliciesLimits[table]
	if !exists {
		return fmt.Errorf("invalid table: %s", table)
	}

	if days < limits.MinDays || days > limits.MaxDays {
		return fmt.Errorf(
			"retention for %s must be between %d and %d days (got %d)",
			table, limits.MinDays, limits.MaxDays, days,
		)
	}

	return nil
}

// GetDefaultRetentionPolicies returns a map of default retention days for all tables
func GetDefaultRetentionPolicies() map[string]uint16 {
	defaults := make(map[string]uint16)
	for table, limits := range RetentionPoliciesLimits {
		defaults[table] = limits.DefaultDays
	}
	return defaults
}

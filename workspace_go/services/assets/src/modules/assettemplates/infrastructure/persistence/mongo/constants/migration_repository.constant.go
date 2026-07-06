package constants

import model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"

// MigrationPlansCollection is the MongoDB collection name for template migration plans.
const MigrationPlansCollection = "template_migration_plans"

// MigrationExecutionsCollection is the MongoDB collection name for per-asset migration executions.
const MigrationExecutionsCollection = "template_migration_executions"

// MigrationPlansIndexes defines the indexes for the template_migration_plans collection.
var MigrationPlansIndexes = []model.IndexDefinition{
	{Name: "idx_org", Keys: map[string]int{"orgId": 1}},
	{Name: "idx_status", Keys: map[string]int{"status": 1}},
	{Name: "idx_from_template", Keys: map[string]int{"fromTemplateId": 1}},
}

// MigrationExecutionsIndexes defines the indexes for the template_migration_executions
// collection (plus a compound planId+status for fast per-status counts).
var MigrationExecutionsIndexes = []model.IndexDefinition{
	{Name: "idx_plan", Keys: map[string]int{"planId": 1}},
	{Name: "idx_status", Keys: map[string]int{"status": 1}},
	{Name: "idx_plan_status", Keys: map[string]int{"planId": 1, "status": 1}},
}

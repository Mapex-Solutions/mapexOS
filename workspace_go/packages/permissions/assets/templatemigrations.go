package permissions

// TemplateMigration Permissions
const (
	// TemplateMigrationList - Permission to list all template migrations
	TemplateMigrationList = "templatemigrations.list"

	// TemplateMigrationCreate - Permission to create a new template migration
	TemplateMigrationCreate = "templatemigrations.create"

	// TemplateMigrationRead - Permission to read a specific template migration
	TemplateMigrationRead = "templatemigrations.read"

	// TemplateMigrationUpdate - Permission to update a template migration
	TemplateMigrationUpdate = "templatemigrations.update"

	// TemplateMigrationDelete - Permission to delete a template migration
	TemplateMigrationDelete = "templatemigrations.delete"

	// TemplateMigrationAll - Wildcard permission for all template migration operations
	TemplateMigrationAll = "templatemigrations.*"
)

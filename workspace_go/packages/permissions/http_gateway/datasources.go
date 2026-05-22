package permissions

// DataSource Permissions
const (
	// DatasourceList - Permission to list all datasources
	DatasourceList = "datasources.list"

	// DatasourceCreate - Permission to create a new datasources
	DatasourceCreate = "datasources.create"

	// DatasourceRead - Permission to read a specific datasources
	DatasourceRead = "datasources.read"

	// DatasourceUpdate - Permission to update a datasources
	DatasourceUpdate = "datasources.update"

	// DatasourceDelete - Permission to delete a datasources
	DatasourceDelete = "datasources.delete"

	// DatasourceAll - Wildcard permission for all datasources operations
	DatasourceAll = "datasources.*"
)

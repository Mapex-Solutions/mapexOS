package permissions

// Job Permissions for the Scheduler service
const (
	// JobList - Permission to list all scheduled jobs
	JobList = "jobs.list"

	// JobCreate - Permission to create a new scheduled job
	JobCreate = "jobs.create"

	// JobRead - Permission to read a specific scheduled job
	JobRead = "jobs.read"

	// JobUpdate - Permission to update a scheduled job
	JobUpdate = "jobs.update"

	// JobDelete - Permission to delete a scheduled job
	JobDelete = "jobs.delete"

	// JobAll - Wildcard permission for all job operations
	JobAll = "jobs.*"
)

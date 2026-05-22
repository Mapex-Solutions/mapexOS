package permissions

// Workflow Permissions
const (
	// WorkflowList - Permission to list all workflows
	WorkflowList = "workflows.list"

	// WorkflowCreate - Permission to create a new workflow
	WorkflowCreate = "workflows.create"

	// WorkflowRead - Permission to read a specific workflow
	WorkflowRead = "workflows.read"

	// WorkflowUpdate - Permission to update a workflow
	WorkflowUpdate = "workflows.update"

	// WorkflowDelete - Permission to delete a workflow
	WorkflowDelete = "workflows.delete"

	// WorkflowAll - Wildcard permission for all workflow operations
	WorkflowAll = "workflows.*"

	// WorkflowInstanceList - Permission to list workflow instance configs
	WorkflowInstanceList = "workflows.instances.list"

	// WorkflowInstanceRead - Permission to read a specific workflow instance config
	WorkflowInstanceRead = "workflows.instances.read"

	// WorkflowInstanceCreate - Permission to create a workflow instance config
	WorkflowInstanceCreate = "workflows.instances.create"

	// WorkflowInstanceUpdate - Permission to update a workflow instance config
	WorkflowInstanceUpdate = "workflows.instances.update"

	// WorkflowInstanceDelete - Permission to delete a workflow instance config
	WorkflowInstanceDelete = "workflows.instances.delete"

	// WorkflowExecutionList - Permission to list workflow executions
	WorkflowExecutionList = "workflows.executions.list"

	// WorkflowExecutionRead - Permission to read a specific workflow execution
	WorkflowExecutionRead = "workflows.executions.read"

	// WorkflowInstanceCancel - Permission to cancel a running workflow execution
	WorkflowInstanceCancel = "workflows.instances.cancel"

	// WorkflowInstanceSignal - Permission to send a signal to a waiting workflow execution
	WorkflowInstanceSignal = "workflows.instances.signal"

	// WorkflowInstanceExecute - Permission to execute a workflow instance (start a new execution)
	WorkflowInstanceExecute = "workflows.instances.execute"
)

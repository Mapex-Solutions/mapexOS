package ports

import (
	"workflow/src/modules/instances/domain/entities"
)

// Public type aliases — other modules import from here, NEVER from domain/entities
type WorkflowInstance = entities.WorkflowInstance

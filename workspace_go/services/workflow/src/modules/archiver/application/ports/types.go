package ports

import (
	sharedTypes "workflow/src/shared/types"
)

// Port-level type aliases — expose shared/cross-module types through the port boundary.
// Application services import these aliases instead of reaching into interfaces/message.

// StateEvent is the WORKFLOW-STATE lifecycle event consumed by the archiver.
// Mirrors the intra-service alias in interfaces/message/types.go but is exposed here
// so application code does not have to import interfaces/message (see Hexagonal layering).
type StateEvent = sharedTypes.StateEvent

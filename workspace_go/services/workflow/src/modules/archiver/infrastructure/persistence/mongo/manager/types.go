package manager

import (
	"workflow/src/modules/archiver/application/ports"

	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
)

/*
 * MONGO MANAGER ADAPTER TYPES
 * Adapter wraps the concrete mongoManager.MongoManager so the archiver
 * application layer depends only on the MongoManagerPort interface.
 */

// MongoManagerAdapter implements ports.MongoManagerPort by delegating to
// the concrete MongoDB manager from mapexGoKit.
type MongoManagerAdapter struct {
	mgr *mongoManager.MongoManager
}

// Compile-time check.
var _ ports.MongoManagerPort = (*MongoManagerAdapter)(nil)

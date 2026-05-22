package mongo

import (
	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
)

// MongoManagerAdapter is the infrastructure adapter that implements
// onboarding_orchestrator/application/ports.MongoManagerPort by delegating
// to the mongoManager.MongoManager driver from mapexGoKit.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: ports.MongoManagerPort
//   - Adapter: MongoManagerAdapter (this struct)
//
// The adapter keeps the concrete *mongoManager.MongoManager driver out of
// the application layer's DI struct.
type MongoManagerAdapter struct {
	mm *mongoManager.MongoManager
}

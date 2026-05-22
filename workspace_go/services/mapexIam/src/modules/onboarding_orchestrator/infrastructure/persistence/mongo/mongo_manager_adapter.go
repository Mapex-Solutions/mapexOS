package mongo

import (
	"context"

	"mapexIam/src/modules/onboarding_orchestrator/application/ports"

	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
)

// NewMongoManagerAdapter builds a MongoManagerAdapter around the given
// mongoManager.MongoManager driver. Returns the adapter as the application
// port interface so callers do not depend on the concrete driver type.
//
// Parameters:
//   - mm: The underlying mongoManager.MongoManager driver (provided by bootstrap).
//
// Returns:
//   - ports.MongoManagerPort: The adapter as an interface.
func NewMongoManagerAdapter(mm *mongoManager.MongoManager) ports.MongoManagerPort {
	return &MongoManagerAdapter{mm: mm}
}

// RunTransaction executes the provided function within a MongoDB transaction.
// Delegates to the underlying driver while preserving the application port's
// TransactionFunc signature.
func (a *MongoManagerAdapter) RunTransaction(ctx context.Context, txnFunc ports.TransactionFunc) (interface{}, error) {
	return a.mm.RunTransaction(ctx, mongoManager.TransactionFunc(txnFunc))
}

// Compile-time check to ensure MongoManagerAdapter implements MongoManagerPort.
var _ ports.MongoManagerPort = (*MongoManagerAdapter)(nil)

package ports

import (
	"context"
)

// TransactionFunc defines a function to be executed within a MongoDB transaction.
// It receives a context with the session bound — all MongoDB operations using
// this context will automatically participate in the transaction.
//
// Return an error to abort the transaction, or nil to commit.
type TransactionFunc func(ctx context.Context) (interface{}, error)

// MongoManagerPort is the driven port exposing MongoDB transaction support to
// the onboarding orchestrator without leaking the concrete
// *mongoManager.MongoManager driver into the application DI layer.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: MongoManagerPort (this interface)
//   - Adapter: infrastructure/persistence/mongo.MongoManagerAdapter
//
// Only the method actually consumed by UserOnboardingService (RunTransaction)
// is exposed here.
type MongoManagerPort interface {
	// RunTransaction executes the provided function within a MongoDB transaction.
	// All MongoDB operations using the provided context participate in the
	// transaction. Returns the function's result on success, or the error that
	// caused the transaction to abort.
	RunTransaction(ctx context.Context, txnFunc TransactionFunc) (interface{}, error)
}

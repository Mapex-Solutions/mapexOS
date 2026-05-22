package di

import (
	"workflow/src/modules/archiver/application/ports"
	"workflow/src/modules/archiver/domain/repositories"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

/*
 * ARCHIVER SERVICE DI
 * Aggregates all dependencies for the ArchiverService.
 */

// ArchiverServiceDependenciesInjection aggregates all dependencies required by ArchiverService.
type ArchiverServiceDependenciesInjection struct {
	dig.In

	// ArchiveRepo provides BulkWrite operations to the instances collection
	ArchiveRepo repositories.ArchiveRepository

	// KVStore provides hot state reads (for terminal events) and key deletion
	KVStore natsModel.KeyValueStore

	// Publisher publishes NATS messages (workflow events to ClickHouse)
	Publisher natsModel.Publisher `name:"core"`

	// MongoManager provides backpressure mode and write latency recording
	MongoManager ports.MongoManagerPort
}

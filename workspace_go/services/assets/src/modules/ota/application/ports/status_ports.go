package ports

import (
	"context"

	"assets/src/modules/ota/domain/entities"

	otaEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/events"
)

// LiveStatePort persists the live execution state (Redis read model) for fast
// status-screen reads.
type LiveStatePort interface {
	SetExecutionLive(ctx context.Context, exec *entities.OTAExecution) error
}

// StatusHistoryPublisherPort publishes a status advisory to the events stream
// for ClickHouse history.
type StatusHistoryPublisherPort interface {
	Publish(ctx context.Context, advisory otaEvents.OTAStatusAdvisory) error
}

// AssetTemplateSwitcherPort switches an asset to a new template — called when a
// device reports UPDATED. Implemented via the assets module's port (never a
// direct write to the assets collection).
type AssetTemplateSwitcherPort interface {
	SwitchTemplate(ctx context.Context, assetID, targetTemplateID string) error
}

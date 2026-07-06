package adapters

import (
	"context"

	"assets/src/modules/ota/application/ports"

	hmPorts "assets/src/modules/healthmonitor/application/ports"
)

// presenceReaderAdapter implements ports.PresenceReaderPort by READING the
// healthmonitor presence read model (it does not modify healthmonitor).
type presenceReaderAdapter struct {
	repo hmPorts.HealthRepository
}

// NewPresenceReaderAdapter returns a PresenceReaderPort over the healthmonitor
// health repository.
func NewPresenceReaderAdapter(repo hmPorts.HealthRepository) ports.PresenceReaderPort {
	return &presenceReaderAdapter{repo: repo}
}

func (a *presenceReaderAdapter) IsOnline(ctx context.Context, orgID, assetUUID string) (bool, error) {
	return a.repo.IsKnownOnline(ctx, orgID, assetUUID)
}

var _ ports.PresenceReaderPort = (*presenceReaderAdapter)(nil)

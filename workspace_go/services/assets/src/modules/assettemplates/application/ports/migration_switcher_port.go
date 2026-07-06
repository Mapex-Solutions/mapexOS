package ports

import "context"

// TemplateSwitcherPort rebinds one asset to a target template through the assets
// module, which owns the read-model re-projection and fanout invalidation side
// effects.
type TemplateSwitcherPort interface {
	SwitchTemplate(ctx context.Context, assetID, targetTemplateID string) error
}

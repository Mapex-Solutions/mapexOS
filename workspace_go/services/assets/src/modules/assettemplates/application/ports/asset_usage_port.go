package ports

import (
	"context"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

// AssetUsagePort counts, through the assets module, how many assets reference a
// template. The marketplace uninstall guard uses it to refuse removing a
// template that is still in use. Org scoping travels via the request context.
type AssetUsagePort interface {
	CountAssetsUsingTemplate(ctx context.Context, requestContext *reqCtx.RequestContext, templateID string) (int64, error)
}

package ports

import (
	ctx "context"

	"assets/src/modules/assets/application/dtos"

	assetsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
	assetsAuthContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/auth"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

// AssetServicePort defines the contract for asset-related business operations.
// This interface follows Hexagonal Architecture principles by defining the
// application's core business logic contract that can be implemented by
// different service implementations.
//
// All methods operate on DTOs for input/output to maintain clean separation
// between the domain layer and external consumers.
type AssetServicePort interface {
	// CreateAsset creates a new asset entity from the provided DTO.
	// Uses RequestContext to populate orgId and pathKey for multi-tenant support.
	// For MQTT-protocol assets, the operator-supplied password (or
	// platform-generated random) is bcrypt-hashed before persistence;
	// the plaintext is never stored and never returned on read endpoints.
	CreateAsset(ctx ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.AssetCreateDTO) (*dtos.AssetResponse, error)

	// GetAssetById retrieves an asset by its unique identifier.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - assetId: Unique identifier of the asset
	//
	// Returns:
	//   - AssetResponse: The found asset
	//   - error: If asset not found or operation fails
	GetAssetById(ctx ctx.Context, assetId *string) (*dtos.AssetResponse, error)

	// GetByUUID returns the asset entity (not the response DTO) for the
	// given device UUID. Used by cross-module callers (devices module's
	// refresh-token flow) that need entity-level fields not exposed on
	// the response shape — most notably Protocol.Mqtt.TokenJti, which
	// drives the jti-match invariant for refresh. Returns NOT_FOUND when
	// the UUID is unknown. The returned type is `ports.Asset`, an alias
	// for `domain/entities.Asset` (entity-boundary alias — cross-module
	// callers depend on the ports package, not the domain entity).
	GetByUUID(ctx ctx.Context, assetUUID string) (*Asset, error)

	// UpdateAssetById updates an existing asset identified by its unique ID.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - assetId: Unique identifier of the asset to update
	//   - dto: Fields to update
	//
	// Returns:
	//   - AssetResponse: The updated asset
	//   - error: If update fails or asset not found
	UpdateAssetById(ctx ctx.Context, assetId *string, dto *dtos.AssetUpdateDTO) (*dtos.AssetResponse, error)

	// GenerateMqttPassword returns a fresh random alphanumeric password
	// (24 chars) for the operator to drop into the asset's MQTT config
	// before submit. Stateless — does not touch any asset record.
	GenerateMqttPassword(ctx ctx.Context) (*dtos.GenerateMqttPasswordResponseDTO, error)

	// DeleteAssetById removes an asset identified by its unique ID.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - assetId: Unique identifier of the asset to delete
	//
	// Returns:
	//   - map[string]bool: Success indicator
	//   - error: If deletion fails or asset not found
	DeleteAssetById(ctx ctx.Context, assetId *string) (map[string]bool, error)

	// GetAssets retrieves a paginated and filtered list of assets.
	// Uses RequestContext from coverage middleware for context-aware org filtering with hierarchical support.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Contains ScopedOrgIds, OrgContext, OrgContextData with PathKey
	//   - query: Filters, pagination, projection, and includeChildren flag
	//
	// Returns:
	//   - PaginatedResult: Matching assets as DTOs and pagination metadata
	//   - error: If query fails
	//
	// Security: Uses orgfilter.BuildOrgFilter() for automatic org filtering based on RequestContext
	GetAssets(ctx ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.AssetQueryDTO) (*model.PaginatedResult[dtos.AssetResponse], error)

	// GetAssetByMqttUsername retrieves the public asset response by
	// MQTT username. Used by callers that need the read-friendly
	// shape; PasswordHash is NOT included. The broker plugin does NOT
	// use this — it consumes the full AssetReadModel via the L2 MinIO
	// bucket or the GetAssetReadModelByUUID internal endpoint, which
	// carry PasswordHash + CurrentCert for local auth decisions.
	GetAssetByMqttUsername(ctx ctx.Context, username string) (*dtos.AssetResponse, error)

	// CountAssets returns the total count of assets for the given org context.
	// Uses Redis cache with 6h TTL, invalidated on create/delete.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Contains org access data from coverage middleware
	//
	// Returns:
	//   - int64: Total count of matching assets
	//   - error: If query fails
	CountAssets(ctx ctx.Context, requestContext *reqCtx.RequestContext) (int64, error)

	// GetAssetReadModelByUUID retrieves the asset read model by UUID.
	// This method is used by internal endpoints for cache fallback.
	//
	// Called by: TieredCache fallback when L2 (MinIO) cache miss
	// Purpose: Fetch asset from MongoDB and repopulate L2 cache
	//
	// The method also writes the read model to MinIO (L2) to ensure
	// subsequent requests can find it in cache.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - assetUUID: The device UUID (not MongoDB ID)
	//
	// Returns:
	//   - AssetReadModel: The asset read model (same format as L2 cache)
	//   - error: If asset not found
	GetAssetReadModelByUUID(ctx ctx.Context, assetUUID string) (*assetsContract.AssetReadModel, error)

	// GetAuthProjectionByUUID returns the slim auth-only projection
	// consumed by the broker plugin as the L3 fallback when L1+L2 both
	// miss. Spawns an async goroutine to warm the mapex-asset-auth
	// bucket back up after returning the projection — the HTTP response
	// does not wait for the warm-up.
	GetAuthProjectionByUUID(ctx ctx.Context, assetUUID string) (*assetsAuthContract.AuthProjection, error)

	// ProcessL2WriteRetry is invoked by the L2 sync fallback consumer
	// when a previous synchronous L2 write failed. The method re-fetches
	// the asset from Mongo by ID (NOT trusting the event payload — Mongo
	// is the source of truth) and re-runs the L2 sync. On success it
	// emits the existing fanout.asset.invalidate so caches downstream
	// refresh. Returns an error if Mongo lookup fails or if the write
	// continues to fail after the helper's internal handling — the
	// consumer NAKs on error to let NATS retry with backoff.
	ProcessL2WriteRetry(ctx ctx.Context, assetId string) error

	// SetCurrentCert mirrors a freshly-issued device cert's metadata
	// onto the asset entity and runs the standard fanout side effects
	// (L2 rewrite + FANOUT cache invalidation). Called by the
	// mqttcerts module after IssueCert returns the PEM bundle so the
	// broker plugin sees the new active serial on the next CONNECT
	// (or sooner, via the FANOUT-driven L1 eviction). Returns
	// NOT_FOUND when assetUUID is unknown.
	SetCurrentCert(ctx ctx.Context, assetUUID string, cert AssetCertificateInput) error

	// ClearCurrentCertBySerial finds the asset whose `currentCert.serial`
	// matches the supplied serial, clears its `currentCert` subdoc, and
	// runs the standard fanout side effects. Returns the cleared asset's
	// UUID so the caller (mqttcerts revoke flow) can stamp the audit row
	// with the asset link. Returns NOT_FOUND when no asset carries that
	// serial (already-revoked or never-issued).
	ClearCurrentCertBySerial(ctx ctx.Context, serial string) (assetUUID string, err error)
}

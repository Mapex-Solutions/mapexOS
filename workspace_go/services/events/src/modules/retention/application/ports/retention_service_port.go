package ports

import (
	ctx "context"

	"events/src/modules/retention/application/dtos"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

// RetentionServicePort defines the contract for retention policy business operations.
//
// This port interface enables Hexagonal Architecture by decoupling the business
// logic from its implementation. The Events service depends on this interface
// to get retention days for each organization.
type RetentionServicePort interface {
	// GetRetentionPolicies retrieves a paginated and filtered list of retention policies.
	// Uses RequestContext from coverage middleware for context-aware org filtering.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Request context with org access data (from coverage middleware)
	//   - query: Filters, pagination, and projection options
	//
	// Returns:
	//   - PaginatedResult: Matching retention policies and pagination metadata
	//   - error: If query fails
	GetRetentionPolicies(ctx ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.RetentionPolicyQueryDTO) (*model.PaginatedResult[dtos.RetentionPolicyResponse], error)

	// GetRetentionPolicyById retrieves a retention policy by its unique identifier.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - policyId: The unique identifier of the retention policy
	//
	// Returns:
	//   - RetentionPolicyResponse: The retention policy if found
	//   - error: If not found or retrieval fails
	GetRetentionPolicyById(ctx ctx.Context, policyId *string) (*dtos.RetentionPolicyResponse, error)

	// UpsertRetentionPolicy creates or updates a retention policy by orgId + type.
	// Uses RequestContext to populate orgId and pathKey for multi-tenant support.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Contains OrgContext (selected orgId) and OrgContextData (pathKey)
	//   - dto: Retention policy upsert data (type, name, retentionDays, enabled)
	//
	// Returns:
	//   - RetentionPolicyResponse: The created/updated retention policy
	//   - error: If upsert fails or validation error
	UpsertRetentionPolicy(ctx ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.RetentionPolicyUpsertDTO) (*dtos.RetentionPolicyResponse, error)

	// DeleteRetentionPolicyById removes a retention policy by its unique identifier.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - policyId: The unique identifier of the retention policy to delete
	//
	// Returns:
	//   - map[string]bool: Success indicator ({"success": true})
	//   - error: If deletion fails or retention policy not found
	DeleteRetentionPolicyById(ctx ctx.Context, policyId *string) (map[string]bool, error)

	// GetRetentionDays retrieves the retention days for a specific table and organization.
	// This is the PRIMARY method used by EventService to determine TTL per record.
	// Uses cache-aside pattern with automatic fallback to defaults.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - orgId: The organization ID
	//   - tableName: The table name (e.g., "events", "eventsRaw")
	//
	// Returns:
	//   - uint16: Retention days for the table
	//   - error: If any critical error occurs
	GetRetentionDays(ctx ctx.Context, orgId string, tableName string) (uint16, error)

	// CreateDefaultPolicies creates 8 default retention policy documents for a new organization.
	// Called by the NATS consumer when an organization is created.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - orgId: The organization ID (hex string)
	//   - pathKey: The organization's hierarchical path key
	//
	// Returns:
	//   - error: If creation fails
	CreateDefaultPolicies(ctx ctx.Context, orgId string, pathKey string) error

	// SeedPlatformPolicies inserts the platform-level retention rows (no
	// orgId). Currently a single row — asset_status_history with 7 days
	// default. Idempotent: re-running is safe.
	SeedPlatformPolicies(ctx ctx.Context) error

	// ApplyAssetStatusHistoryTTL applies the given retention days to the
	// asset_status_history ClickHouse table via ALTER TABLE ... MODIFY TTL.
	// Called from the upsert path whenever the policy for that type changes.
	ApplyAssetStatusHistoryTTL(ctx ctx.Context, days uint16) error

	// HandleOrgCreatedEvent processes an organization.created NATS message.
	// The service owns the full lifecycle: parse, validate, act.
	// Returns nil if the message should be acked (success OR unrecoverable
	// parse/validation issue that must not be redelivered) and a non-nil
	// error when the message should be nacked for retry.
	HandleOrgCreatedEvent(msg *natsModel.Message) error
}

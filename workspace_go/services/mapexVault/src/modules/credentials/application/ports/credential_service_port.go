package ports

import (
	"context"

	"mapexVault/src/modules/credentials/application/dtos"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

// CredentialServicePort defines the business operations for vault credential management.
type CredentialServicePort interface {
	// Credential CRUD
	CreateCredential(ctx context.Context, requestContext *reqCtx.RequestContext, dto *dtos.CreateCredentialDTO) (*dtos.CredentialResponse, error)
	GetCredentialById(ctx context.Context, id string) (*dtos.CredentialResponse, error)
	UpdateCredentialById(ctx context.Context, id string, dto *dtos.UpdateCredentialDTO) (*dtos.CredentialResponse, error)
	DeleteCredentialById(ctx context.Context, id string) (map[string]bool, error)
	GetCredentials(ctx context.Context, requestContext *reqCtx.RequestContext, query *dtos.CredentialQueryDTO) (*model.PaginatedResult[dtos.CredentialResponse], error)

	// Internal (for service-to-service)
	DecryptCredential(ctx context.Context, id string) (map[string]interface{}, error)

	// Test
	TestCredential(ctx context.Context, id string) (map[string]bool, error)

	// OAuth2
	HandleOAuthCallback(ctx context.Context, requestContext *reqCtx.RequestContext, dto *dtos.OAuthCallbackDTO) (*dtos.CredentialResponse, error)

	// userAndPass
	HandleLogin(ctx context.Context, credentialId string) error

	// Connection CRUD
	CreateConnection(ctx context.Context, requestContext *reqCtx.RequestContext, dto *dtos.CreateConnectionDTO) (*dtos.ConnectionResponse, error)
	GetConnections(ctx context.Context, requestContext *reqCtx.RequestContext, query *dtos.ConnectionQueryDTO) (*model.PaginatedResult[dtos.ConnectionResponse], error)
	UpsertConnection(ctx context.Context, requestContext *reqCtx.RequestContext, dto *dtos.UpsertConnectionDTO) (*dtos.ConnectionResponse, error)

	// Reconciler — safety-net loop invoked by the VAULT-RECONCILER consumer.
	// Reseeds any VAULT-SCHEDULE timer missing for active credentials.
	RunReconcile(ctx context.Context)

	// HandleRefreshMessage processes a scheduled refresh message from
	// vault.schedule.fired. Invoked by the refresh consumer registered in module.go.
	HandleRefreshMessage(msg *natsModel.Message)
}

package services

import (
	"context"
	"fmt"

	"mapexVault/src/modules/credentials/application/constants"
	"mapexVault/src/modules/credentials/application/di"
	"mapexVault/src/modules/credentials/application/dtos"
	"mapexVault/src/modules/credentials/application/ports"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	common "github.com/Mapex-Solutions/mapexGoKit/microservices/common"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// Compile-time checks
var _ ports.CredentialServicePort = (*CredentialService)(nil)
var _ common.Mountable = (*CredentialService)(nil)

// New creates a new CredentialService with all dependencies injected.
func New(deps di.CredentialServiceDependenciesInjection) ports.CredentialServicePort {
	return &CredentialService{deps: deps}
}

// OnMount is the lifecycle hook called by common.RunLifecycleHooks after
// DI is fully wired. Seeds initial refresh schedules for existing
// credentials and arms the first reconcile timer on VAULT-RECONCILER.
func (s *CredentialService) OnMount() {
	logger.Info("[SERVICE:Credential] OnMount: seeding refresh schedules and arming reconciler")
	s.bootstrapSeed()
	s.scheduleNextReconcile()
}

// CreateCredential orchestrates credential creation: envelope-encrypt the
// secret payload -> build the entity with org/path scoping -> persist ->
// return the public response DTO (no encrypted blobs).
func (s *CredentialService) CreateCredential(ctx context.Context, requestContext *reqCtx.RequestContext, dto *dtos.CreateCredentialDTO) (*dtos.CredentialResponse, error) {
	env, err := encryptData(s.deps.Encryption, dto.Data)
	if err != nil {
		return nil, err
	}
	entity := s.buildCredentialEntity(requestContext, dto, env)
	result, err := s.deps.CredentialRepo.Create(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Failed to create: %w", err)
	}
	logger.Info(fmt.Sprintf("[SERVICE:Credential] Created credential %s (type=%s, plugin=%s)", result.ID.Hex(), result.Type, result.PluginId))
	return toCredentialResponse(result), nil
}

// UpdateCredentialById applies a partial update to a credential, re-running
// envelope encryption only when the patch carries new secret data.
// Returns the public response DTO.
func (s *CredentialService) UpdateCredentialById(ctx context.Context, id string, dto *dtos.UpdateCredentialDTO) (*dtos.CredentialResponse, error) {
	update, err := s.buildCredentialUpdateMap(dto)
	if err != nil {
		return nil, err
	}
	result, err := s.deps.CredentialRepo.FindByIdAndUpdate(ctx, &id, update)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Failed to update: %w", err)
	}
	logger.Info(fmt.Sprintf("[SERVICE:Credential] Updated credential %s", id))
	return toCredentialResponse(result), nil
}

// GetCredentialById fetches a credential by id and shapes the public
// response DTO (without the encrypted blobs). Returns a wrapped error
// when the repository layer fails or the id is unknown.
func (s *CredentialService) GetCredentialById(ctx context.Context, id string) (*dtos.CredentialResponse, error) {
	cred, err := s.deps.CredentialRepo.FindById(ctx, &id)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Not found: %w", err)
	}
	if cred == nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Credential %s not found", id)
	}
	return toCredentialResponse(cred), nil
}

// GetCredentials returns the paginated, org-scoped credential list.
// Orchestration: filter -> paginate -> query -> map.
func (s *CredentialService) GetCredentials(ctx context.Context, requestContext *reqCtx.RequestContext, query *dtos.CredentialQueryDTO) (*model.PaginatedResult[dtos.CredentialResponse], error) {
	filters := s.buildCredentialListFilters(requestContext, query)
	pagination := s.buildCredentialListPagination(query)
	result, err := s.deps.CredentialRepo.FindWithFilters(ctx, filters, pagination, model.Map{"created": -1})
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Failed to list: %w", err)
	}
	return s.mapCredentialList(result), nil
}

// DeleteCredentialById orchestrates removal: purge any pending refresh
// schedule first (best-effort) -> delete the row in Mongo. The schedule
// purge runs first so a still-pending timer cannot fire after deletion.
func (s *CredentialService) DeleteCredentialById(ctx context.Context, id string) (map[string]bool, error) {
	subject := fmt.Sprintf("%s.%s", constants.VaultScheduleSubjectPrefix, id)
	if err := s.deps.ScheduleManager.PurgeStreamSubject(constants.VaultScheduleStreamName, subject); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Credential] Failed to purge schedule for %s: %v", id, err))
	}
	if err := s.deps.CredentialRepo.DeleteById(ctx, &id); err != nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Failed to delete: %w", err)
	}
	logger.Info(fmt.Sprintf("[SERVICE:Credential] Deleted credential %s", id))
	return map[string]bool{"deleted": true}, nil
}

// DecryptCredential returns the plaintext credential data + metadata.
// Internal API only — never exposed via the public HTTP surface.
func (s *CredentialService) DecryptCredential(ctx context.Context, id string) (map[string]interface{}, error) {
	cred, err := s.fetchCredentialOrError(ctx, id, "decrypt")
	if err != nil {
		return nil, err
	}
	data, err := decryptData(s.deps.Encryption, cred)
	if err != nil {
		return nil, err
	}
	data["__pluginId"] = cred.PluginId
	data["__credentialDefId"] = cred.CredentialDefId
	return data, nil
}

// TestCredential decrypts the credential and reports a boolean success
// envelope. The decrypt itself is the implicit test — if envelope
// decryption fails the credential is considered broken regardless of the
// underlying reason.
func (s *CredentialService) TestCredential(ctx context.Context, id string) (map[string]bool, error) {
	if _, err := s.DecryptCredential(ctx, id); err != nil {
		return map[string]bool{"success": false}, err
	}
	return map[string]bool{"success": true}, nil
}

// HandleLogin performs token acquisition using the credential's LoginConfig.
// Steps: load + decrypt -> execute the configured token request -> persist
// the resulting tokens (encrypted at rest).
func (s *CredentialService) HandleLogin(ctx context.Context, credentialId string) error {
	cred, data, err := s.loadAndDecryptForLogin(ctx, credentialId)
	if err != nil {
		return err
	}
	resp, err := s.executeTokenRequest(cred, cred.ProviderConfig.LoginConfig, data)
	if err != nil {
		return fmt.Errorf("[SERVICE:Credential] Login failed for %s: %w", credentialId, err)
	}
	if err := s.updateCredentialTokens(cred, data, resp); err != nil {
		return fmt.Errorf("[SERVICE:Credential] Failed to update tokens after login for %s: %w", credentialId, err)
	}
	logger.Info(fmt.Sprintf("[SERVICE:Credential] Login completed for credential %s", credentialId))
	return nil
}

// HandleOAuthCallback is the OAuth2 authorization-code callback entry
// point. The actual provider-specific token exchange (Instagram, TikTok,
// etc.) is not yet implemented — this method enforces input validation
// and surfaces a clear "not implemented" error until each provider lands.
func (s *CredentialService) HandleOAuthCallback(ctx context.Context, requestContext *reqCtx.RequestContext, dto *dtos.OAuthCallbackDTO) (*dtos.CredentialResponse, error) {
	if dto.PluginId == "" {
		return nil, fmt.Errorf("[SERVICE:Credential] pluginId is required for OAuth2 callback")
	}
	return nil, fmt.Errorf("[SERVICE:Credential] OAuth2 provider-specific exchange not yet implemented — requires provider ticket (Instagram, TikTok, etc.)")
}

// HandleRefreshMessage processes a scheduled refresh message from
// vault.schedule.fired. Steps: parse payload -> load + filter inactive
// credentials -> select login or refresh config -> run token refresh ->
// persist new tokens. Always Acks (failures self-mark the credential).
func (s *CredentialService) HandleRefreshMessage(msg *natsModel.Message) {
	credentialId, ok := s.parseRefreshPayload(msg)
	if !ok {
		return
	}
	cred, ok := s.loadCredentialForRefresh(credentialId, msg)
	if !ok {
		return
	}
	cfg, ok := s.selectTokenRefreshConfig(cred, credentialId, msg)
	if !ok {
		return
	}
	s.runTokenRefresh(cred, cfg, credentialId, msg)
}

// CreateConnection creates a new connection entity scoped to the caller's
// org. Steps: build entity with org/path -> persist.
func (s *CredentialService) CreateConnection(ctx context.Context, requestContext *reqCtx.RequestContext, dto *dtos.CreateConnectionDTO) (*dtos.ConnectionResponse, error) {
	entity, err := s.buildConnectionEntity(requestContext, dto)
	if err != nil {
		return nil, err
	}
	result, err := s.deps.ConnectionRepo.Create(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Failed to create connection: %w", err)
	}
	return toConnectionResponse(result), nil
}

// GetConnections returns the paginated, org-scoped connection list.
// Orchestration: filter -> paginate -> query -> map.
func (s *CredentialService) GetConnections(ctx context.Context, requestContext *reqCtx.RequestContext, query *dtos.ConnectionQueryDTO) (*model.PaginatedResult[dtos.ConnectionResponse], error) {
	filters := s.buildConnectionListFilters(requestContext, query)
	pagination := s.buildConnectionListPagination(query)
	result, err := s.deps.ConnectionRepo.FindWithFilters(ctx, filters, pagination, model.Map{"connectedAt": -1})
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Failed to list connections: %w", err)
	}
	return s.mapConnectionList(result), nil
}

// UpsertConnection creates or updates a connection by (provider, accountId,
// orgId). Used for OAuth flows where the same external account reconnects.
func (s *CredentialService) UpsertConnection(ctx context.Context, requestContext *reqCtx.RequestContext, dto *dtos.UpsertConnectionDTO) (*dtos.ConnectionResponse, error) {
	entity, orgId, err := s.buildUpsertConnectionEntity(requestContext, dto)
	if err != nil {
		return nil, err
	}
	result, err := s.deps.ConnectionRepo.UpsertByAccount(ctx, dto.Provider, dto.AccountId, orgId, entity)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:Credential] Failed to upsert connection: %w", err)
	}
	return toConnectionResponse(result), nil
}

// RunReconcile is invoked by the reconcile consumer on
// vault.reconcile.fired. Iterates active credentials, reseeds any missing
// refresh schedules, and re-arms the next reconcile timer.
func (s *CredentialService) RunReconcile(ctx context.Context) {
	credentials, err := s.deps.CredentialRepo.FindActiveWithTokenExpiry(ctx)
	if err != nil {
		logger.Error(err, "[SERVICE:Credential] Reconciler failed to query credentials")
		s.scheduleNextReconcile()
		return
	}
	checked, reseeded := s.reseedMissingSchedules(credentials)
	logger.Info(fmt.Sprintf("[SERVICE:Credential] Reconciler completed: checked=%d reseeded=%d", checked, reseeded))
	s.scheduleNextReconcile()
}

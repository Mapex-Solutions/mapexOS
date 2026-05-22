package collection

import (
	"context"

	"mapexVault/src/modules/credentials/domain/entities"
	"mapexVault/src/modules/credentials/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// NewCredentialRepository creates a repository for Credential entities.
//
// Ensures two indexes on first initialization (idempotent):
//   - idx_reconciler_active_expiry: optimizes the reconciler/bootstrap query
//     {status, type, tokenExpiresAt} used to find active credentials eligible
//     for scheduled refresh.
//   - idx_org_status: optimizes per-organization credential listings.
func NewCredentialRepository(m *manager.MongoManager) repositories.CredentialRepository {
	mdl := model.New[entities.Credential](m.GetDatabase(), credentialCollectionName, model.Config{
		Indexes: []model.IndexDefinition{
			{
				Name: "idx_reconciler_active_expiry",
				Keys: map[string]int{
					"status":         1,
					"type":           1,
					"tokenExpiresAt": 1,
				},
			},
			{
				Name: "idx_org_status",
				Keys: map[string]int{
					"orgId":  1,
					"status": 1,
				},
			},
		},
	})
	return &credentialRepository{model: mdl}
}

// Create inserts a new Credential entity into MongoDB.
func (r *credentialRepository) Create(ctx context.Context, entity *entities.Credential) (*entities.Credential, error) {
	return r.model.CreateOne(ctx, entity)
}

// FindById retrieves a Credential by its MongoDB ObjectId.
func (r *credentialRepository) FindById(ctx context.Context, id *string) (*entities.Credential, error) {
	return r.model.FindByID(ctx, *id)
}

// FindByIdAndUpdate updates a Credential by ID using $set operator.
func (r *credentialRepository) FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.Credential, error) {
	_id, err := model.ToObjectID(*id)
	if err != nil {
		return nil, err
	}

	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc}

	return r.model.FindOneAndUpdate(ctx, &query, &update, &options)
}

// DeleteById removes a Credential by its MongoDB ObjectId.
func (r *credentialRepository) DeleteById(ctx context.Context, id *string) error {
	_id, err := model.ToObjectID(*id)
	if err != nil {
		return err
	}

	query := model.Map{"_id": _id}
	return r.model.DeleteOne(ctx, &query)
}

// FindWithFilters retrieves a paginated list of credentials.
func (r *credentialRepository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	sort model.Map,
) (*model.PaginatedResult[entities.Credential], error) {
	opts := &model.CommonOpts{}
	if len(sort) > 0 {
		opts.Sort = sort
	}
	// Exclude encrypted fields from list queries
	opts.Projection = model.Map{
		"encryptedDEK":  0,
		"dekNonce":      0,
		"encryptedData": 0,
		"dataNonce":     0,
	}
	return r.model.FindByOffset(ctx, filters, pagination, opts)
}

// FindActiveWithTokenExpiry returns all active oauth2/userAndPass credentials
// with non-nil tokenExpiresAt. Used by bootstrap seed to publish initial schedules.
func (r *credentialRepository) FindActiveWithTokenExpiry(ctx context.Context) ([]entities.Credential, error) {
	query := model.Map{
		"tokenExpiresAt": model.Map{"$ne": nil},
		"status":         string(entities.CredentialStatusActive),
		"type": model.Map{"$in": []string{
			string(entities.CredentialOAuth2),
			string(entities.CredentialUserAndPass),
		}},
	}

	result, err := r.model.FindByOffset(ctx, query, &model.PaginationOpts{
		Page:    1,
		PerPage: 500,
	}, &model.CommonOpts{})
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// CountDocuments returns the count of documents matching the filter.
func (r *credentialRepository) CountDocuments(ctx context.Context, filters model.Map) (int64, error) {
	count, err := r.model.DIRECT().CountDocuments(ctx, filters)
	return count, err
}

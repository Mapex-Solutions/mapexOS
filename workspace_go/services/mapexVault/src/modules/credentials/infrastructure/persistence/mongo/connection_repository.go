package collection

import (
	"context"

	"mapexVault/src/modules/credentials/domain/entities"
	"mapexVault/src/modules/credentials/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// NewConnectionRepository creates a repository for Connection entities.
func NewConnectionRepository(m *manager.MongoManager) repositories.ConnectionRepository {
	mdl := model.New[entities.Connection](m.GetDatabase(), connectionCollectionName, model.Config{})
	return &connectionRepository{model: mdl}
}

// Create inserts a new Connection entity into MongoDB.
func (r *connectionRepository) Create(ctx context.Context, entity *entities.Connection) (*entities.Connection, error) {
	return r.model.CreateOne(ctx, entity)
}

// FindById retrieves a Connection by its MongoDB ObjectId.
func (r *connectionRepository) FindById(ctx context.Context, id *string) (*entities.Connection, error) {
	return r.model.FindByID(ctx, *id)
}

// FindByIdAndUpdate updates a Connection by ID using $set operator.
func (r *connectionRepository) FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.Connection, error) {
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

// DeleteById removes a Connection by its MongoDB ObjectId.
func (r *connectionRepository) DeleteById(ctx context.Context, id *string) error {
	_id, err := model.ToObjectID(*id)
	if err != nil {
		return err
	}

	query := model.Map{"_id": _id}
	return r.model.DeleteOne(ctx, &query)
}

// FindWithFilters retrieves a paginated list of connections.
func (r *connectionRepository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	sort model.Map,
) (*model.PaginatedResult[entities.Connection], error) {
	opts := &model.CommonOpts{}
	if len(sort) > 0 {
		opts.Sort = sort
	}
	return r.model.FindByOffset(ctx, filters, pagination, opts)
}

// UpsertByAccount creates or updates a connection by provider + accountId + orgId.
// Used when a user reconnects the same external account — overwrites instead of duplicating.
func (r *connectionRepository) UpsertByAccount(
	ctx context.Context,
	provider string,
	accountId string,
	orgId *model.ObjectId,
	connection *entities.Connection,
) (*entities.Connection, error) {
	query := model.Map{
		"provider":  provider,
		"accountId": accountId,
	}
	if orgId != nil {
		query["orgId"] = orgId
	}

	update := model.Map{
		"$set": model.Map{
			"accountName":  connection.AccountName,
			"status":       connection.Status,
			"credentialId": connection.CredentialId,
			"userId":       connection.UserId,
			"pathKey":      connection.PathKey,
			"scopes":       connection.Scopes,
			"connectedAt":  connection.ConnectedAt,
			"updated":      connection.Updated,
		},
		"$setOnInsert": model.Map{
			"provider":  provider,
			"accountId": accountId,
			"orgId":     orgId,
			"created":   connection.Created,
		},
	}

	upsert := true
	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{
		Upsert:         &upsert,
		ReturnDocument: &returnDoc,
	}

	return r.model.FindOneAndUpdate(ctx, &query, &update, &options)
}

// CountDocuments returns the count of documents matching the filter.
func (r *connectionRepository) CountDocuments(ctx context.Context, filters model.Map) (int64, error) {
	count, err := r.model.DIRECT().CountDocuments(ctx, filters)
	return count, err
}

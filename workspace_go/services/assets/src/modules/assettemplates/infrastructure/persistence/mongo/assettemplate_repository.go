package collection

import (
	"context"
	"fmt"

	"assets/src/modules/assettemplates/domain/entities"
	"assets/src/modules/assettemplates/domain/repositories"
	"assets/src/modules/assettemplates/infrastructure/persistence/mongo/constants"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// New creates and returns a generic repository for the XXX entity.
// It accepts a *MongoManager to obtain the database connection,
//
// Then calls model.New to initialize a Model[XXX],
// targeting the configured collection name with default settings.
func New(m *manager.MongoManager) repositories.AssetTemplateRepository {
	mdl := model.New[entities.Assettemplate](m.GetDatabase(), constants.CollectionName, model.Config{
		Indexes: constants.Indexes,
	})
	return &repository{model: mdl}
}

/* REPOSITORY METHODS */

// Create inserts a new Assettemplate entity into the repository.
// It accepts a context for cancellation and timeouts, and a pointer to a Assettemplate entity to be persisted.
// Internally, it calls the underlying model’s CreateOne method to store the document in MongoDB.
// It returns the created Assettemplate (populated with any database-assigned fields) and an error if something goes wrong.
func (r *repository) Create(ctx context.Context, u *entities.Assettemplate) (*entities.Assettemplate, error) {
	event, err := r.model.CreateOne(ctx, u)
	return event, err
}

// FindById retrieves a Assettemplate entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a string representing the event ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - dataSourceId: A string representing the unique identifier of the Assettemplate to be retrieved.
//
// Returns:
//   - A pointer to the Assettemplate entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindById(ctx context.Context, dataSourceId *string) (*entities.Assettemplate, error) {
	retData, _ := r.model.FindByID(ctx, *dataSourceId)
	return retData, nil
}

// FindByMarketplaceGuidAndOrg returns the caller-org's link record for a
// marketplace-installed template, or nil when that org has not installed it.
func (r *repository) FindByMarketplaceGuidAndOrg(ctx context.Context, marketplaceGuid string, orgId model.ObjectId) (*entities.Assettemplate, error) {
	query := model.Map{"marketplaceGuid": marketplaceGuid, "orgId": orgId}
	retData, _ := r.model.FindOne(ctx, &query, nil)
	return retData, nil
}

// FindMarketplaceContentByGuid returns the shared content document for a
// marketplace template — the org-less record that holds the heavy body once per
// guid — or nil when it has not been downloaded yet.
func (r *repository) FindMarketplaceContentByGuid(ctx context.Context, marketplaceGuid string) (*entities.Assettemplate, error) {
	query := model.Map{"marketplaceGuid": marketplaceGuid, "orgId": model.Map{"$exists": false}, "isMarketplace": true}
	retData, _ := r.model.FindOne(ctx, &query, nil)
	return retData, nil
}

// UpsertMarketplaceContent creates or updates in place the single shared content
// document for a marketplaceGuid and returns it with its assigned _id. The
// org-less + isMarketplace filter targets the shared document only, never a
// per-org link, so a version bump replaces the shared body for every tenant.
func (r *repository) UpsertMarketplaceContent(ctx context.Context, content *entities.Assettemplate) (*entities.Assettemplate, error) {
	if content == nil || content.MarketplaceGuid == nil {
		return nil, fmt.Errorf("marketplace content requires a marketplaceGuid")
	}

	set, err := marketplaceContentSet(content)
	if err != nil {
		return nil, err
	}

	query := model.Map{"marketplaceGuid": *content.MarketplaceGuid, "orgId": model.Map{"$exists": false}, "isMarketplace": true}
	update := model.Map{"$set": set}

	upsert := true
	returnDoc := model.ReturnDoc(1)
	opts := &model.CommonOpts{Upsert: &upsert, ReturnDocument: &returnDoc}

	return r.model.FindOneAndUpdate(ctx, &query, &update, opts)
}

// marketplaceContentSet turns the shared content entity into a $set document,
// dropping _id so the upsert never tries to overwrite the immutable key. orgId /
// pathKey are omitempty on the entity and stay absent for the org-less doc.
func marketplaceContentSet(content *entities.Assettemplate) (model.Map, error) {
	data, err := bson.Marshal(content)
	if err != nil {
		return nil, err
	}
	var set model.Map
	if err := bson.Unmarshal(data, &set); err != nil {
		return nil, err
	}
	delete(set, "_id")
	return set, nil
}

// FindInstalledGuids returns the subset of the given marketplace guids the org
// has installed, in one query. It reads only the marketplaceGuid field and caps
// the page at the input size (each guid links at most once per org).
func (r *repository) FindInstalledGuids(ctx context.Context, guids []string, orgId model.ObjectId) ([]string, error) {
	if len(guids) == 0 {
		return []string{}, nil
	}

	filters := model.Map{"marketplaceGuid": model.Map{"$in": guids}, "orgId": orgId}
	pagination := &model.PaginationOpts{Page: 1, PerPage: int64(len(guids))}
	opts := &model.CommonOpts{Projection: model.Map{"marketplaceGuid": 1}}

	result, err := r.model.FindByOffset(ctx, filters, pagination, opts)
	if err != nil {
		return nil, err
	}

	installed := make([]string, 0, len(result.Items))
	for i := range result.Items {
		if g := result.Items[i].MarketplaceGuid; g != nil {
			installed = append(installed, *g)
		}
	}
	return installed, nil
}

// FindByIdAndUpdate updates a Assettemplate entity in the repository by its ID.
// It accepts a context for cancellation and timeouts, a pointer to the event ID,
// and a map containing the fields to be updated.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - dataSourceId: A pointer to the string representing the unique identifier of the Assettemplate to be updated.
//   - payload: A map containing the fields and their new values to update in the Assettemplate entity.
//
// Returns:
//   - A pointer to the updated Assettemplate entity, populated with the new values from the database.
//   - An error if the update operation fails or if the Assettemplate is not found.
func (r *repository) FindByIdAndUpdate(ctx context.Context, dataSourceId *string, payload map[string]any) (*entities.Assettemplate, error) {

	_id, _ := model.ToObjectID(*dataSourceId)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc} // Return the updated document

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

// DeleteById removes a Assettemplate entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a pointer to the event ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - dataSourceId: A pointer to the string representing the unique identifier of the Assettemplate to be deleted.
//
// Returns:
//   - An error if the deletion operation fails or if the Assettemplate is not found.
//   - nil if the Assettemplate is successfully deleted.
func (r *repository) DeleteById(ctx context.Context, dataSourceId *string) error {

	// Convert the string ID to an ObjectID for MongoDB
	_id, _ := model.ToObjectID(*dataSourceId)
	query := model.Map{"_id": _id}

	// Delete one document from the collection
	err := r.model.DeleteOne(ctx, &query)
	return err
}

// FindWithFilters retrieves a paginated list of Assettemplate entities from the repository,
// applying filters, pagination, and projection options.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - filters: A map of filters to apply to the query (e.g., {"isSystem": true, "status": true}).
//   - pagination: Pagination options including page number and items per page.
//   - projection: A map specifying which fields to include in the results (e.g., {"name": 1, "isSystem": 1}).
//
// Returns:
//   - A pointer to PaginatedResult containing the matching Assettemplate entities and pagination metadata.
//   - An error if the query fails.
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.Assettemplate], error) {
	// Build options with projection
	opts := &model.CommonOpts{}
	if len(projection) > 0 {
		opts.Projection = projection
	}

	// Execute paginated query
	result, err := r.model.FindByOffset(ctx, filters, pagination, opts)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateMany updates multiple Assettemplate entities matching the filter criteria.
// This method is used for bulk updates, such as synchronizing denormalized names.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - filter: A map of filters to match documents (e.g., {"manufacturerId": ObjectId(...)})
//   - update: A map containing the update operations (e.g., {"$set": {"manufacturerName": "Milesight"}})
//
// Returns:
//   - int64: The number of documents matched by the filter
//   - error: If the update operation fails
func (r *repository) UpdateMany(ctx context.Context, filter model.Map, update model.Map) (int64, error) {
	result, err := r.model.FindAndUpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.MatchedCount, nil
}

// CountDocuments counts documents matching the provided filters.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals
//   - filters: A map of filters to apply to the count query
//
// Returns:
//   - int64: The number of matching documents
//   - error: If the count operation fails
func (r *repository) CountDocuments(ctx context.Context, filters model.Map) (int64, error) {
	count, err := r.model.DIRECT().CountDocuments(ctx, filters)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Compile-time check to ensure repository implements AssetTemplateRepository
var _ repositories.AssetTemplateRepository = (*repository)(nil)

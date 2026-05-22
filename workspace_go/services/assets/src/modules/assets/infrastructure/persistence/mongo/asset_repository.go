package collection

import (
	"context"
	"time"

	"assets/src/modules/assets/domain/entities"
	"assets/src/modules/assets/domain/repositories"
	"assets/src/modules/assets/infrastructure/persistence/mongo/constants"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

var _ repositories.AssetRepository = (*repository)(nil)

// New creates and returns a generic repository for the XXX entity.
// It accepts a *MongoManager to obtain the database connection,
//
// Then calls model.New to initialize a Model[XXX],
// targeting the configured collection name with default settings.
func New(m *manager.MongoManager) repositories.AssetRepository {
	mdl := model.New[entities.Asset](m.GetDatabase(), constants.CollectionName, model.Config{
		Indexes: constants.Indexes,
	})
	return &repository{model: mdl}
}

/* REPOSITORY METHODS */

// Create inserts a new Asset entity into the repository.
// It accepts a context for cancellation and timeouts, and a pointer to a Asset entity to be persisted.
// Internally, it calls the underlying model’s CreateOne method to store the document in MongoDB.
// It returns the created Asset (populated with any database-assigned fields) and an error if something goes wrong.
func (r *repository) Create(ctx context.Context, u *entities.Asset) (*entities.Asset, error) {
	event, err := r.model.CreateOne(ctx, u)
	return event, err
}

// FindById retrieves a Asset entity from the repository by its ID.
// It accepts a context for cancellation and timeouts, and a string representing the event ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - assetId: A string representing the unique identifier of the Asset to be retrieved.
//
// Returns:
//   - A pointer to the Asset entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindById(ctx context.Context, assetId *string) (*entities.Asset, error) {
	retData, err := r.model.FindByID(ctx, *assetId)
	return retData, err
}

// FindByAssetUUID retrieves a Asset entity from the repository by its assetUUID field.
// This method is used to find assets by their device identifier (devEUI, deviceId, etc).
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - assetUUID: The device identifier (assetUUID field) to search for.
//
// Returns:
//   - A pointer to the Asset entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindByAssetUUID(ctx context.Context, assetUUID *string) (*entities.Asset, error) {
	query := model.Map{"assetUUID": *assetUUID}
	retData, _ := r.model.FindOne(ctx, &query, nil)
	return retData, nil
}

// FindByMqttUsername retrieves a Asset entity from the repository by its MQTT username.
// This method is used by Auth Callout to find assets for authentication.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - username: The MQTT username (protocol.mqtt.username field) to search for.
//
// Returns:
//   - A pointer to the Asset entity if found, or nil if not found.
//   - An error if the retrieval operation fails.
func (r *repository) FindByMqttUsername(ctx context.Context, username string) (*entities.Asset, error) {
	query := model.Map{"protocol.mqtt.username": username}
	retData, err := r.model.FindOne(ctx, &query, nil)
	return retData, err
}

// FindByIdAndUpdate updates a Asset entity in the repository by its ID.
// It accepts a context for cancellation and timeouts, a pointer to the event ID,
// and a map containing the fields to be updated.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - assetId: A pointer to the string representing the unique identifier of the Asset to be updated.
//   - payload: A map containing the fields and their new values to update in the Asset entity.
//
// Returns:
//   - A pointer to the updated Asset entity, populated with the new values from the database.
//   - An error if the update operation fails or if the Asset is not found.
func (r *repository) FindByIdAndUpdate(ctx context.Context, assetId *string, payload map[string]any) (*entities.Asset, error) {

	_id, _ := model.ToObjectID(*assetId)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc} // Return the updated document

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

// DeleteById removes a Asset entity from the repository by its ID and orgId.
// It accepts a context for cancellation and timeouts, orgId for multi-tenant isolation,
// and a pointer to the asset ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - orgId: A pointer to the organization ID to ensure multi-tenant data isolation.
//   - assetId: A pointer to the string representing the unique identifier of the Asset to be deleted.
//
// Returns:
//   - An error if the deletion operation fails or if the Asset is not found.
//   - nil if the Asset is successfully deleted.
func (r *repository) DeleteById(ctx context.Context, assetId *string) error {

	// Convert the string ID to an ObjectID for MongoDB
	_id, _ := model.ToObjectID(*assetId)

	// Query using only _id (coverage middleware already validated access)
	query := model.Map{"_id": _id}

	// Delete one document from the collection
	err := r.model.DeleteOne(ctx, &query)
	return err
}

// FindWithFilters retrieves a paginated list of Asset entities from the repository,
// applying filters, pagination, and projection options.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - filters: A map of filters to apply to the query (e.g., {"status": true, "assetType": "sensor"}).
//   - pagination: Pagination options including page number and items per page.
//   - projection: A map specifying which fields to include in the results (e.g., {"name": 1, "status": 1}).
//
// Returns:
//   - A pointer to PaginatedResult containing the matching Asset entities and pagination metadata.
//   - An error if the query fails.
func (r *repository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.Asset], error) {
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

// UpdateHealthStatusWithChangedAt atomically updates healthStatus and healthStatusChangedAt
// for an asset identified by assetUUID. Called ONLY on state transitions (online→offline,
// offline→online, unknown→online). Both fields are written in a single $set so a caller
// that receives nil error is guaranteed the pair flipped together.
func (r *repository) UpdateHealthStatusWithChangedAt(ctx context.Context, assetUUID *string, status string, changedAt time.Time) error {
	query := model.Map{"assetUUID": *assetUUID}
	update := model.Map{"$set": model.Map{
		"healthStatus":          status,
		"healthStatusChangedAt": changedAt,
	}}
	_, err := r.model.FindOneAndUpdate(ctx, &query, &update, nil)
	return err
}

// FindWithFiltersAndTemplate retrieves assets with template data joined via $lookup aggregation.
// Routes to an optimized fast path when no template filters are present (90%+ of requests),
// or falls back to the full pipeline when template-level filtering is needed.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals
//   - assetFilters: Filters for assets collection (e.g., {"orgId": ObjectId, "status": true})
//   - templateFilters: Filters for template fields (e.g., {"categoryId": ObjectId, "manufacturerId": ObjectId})
//   - pagination: Page and PerPage for pagination
//   - sort: Sort specification (e.g., {"created": -1})
//
// Returns:
//   - PaginatedResult with AssetWithTemplate entities
//   - Error if aggregation fails
func (r *repository) FindWithFiltersAndTemplate(
	ctx context.Context,
	assetFilters model.Map,
	templateFilters model.Map,
	pagination *model.PaginationOpts,
	sort model.Map,
) (*model.PaginatedResult[entities.AssetWithTemplate], error) {

	if len(templateFilters) == 0 {
		return r.findAssetsOptimized(ctx, assetFilters, pagination, sort)
	}
	return r.findAssetsWithTemplateFilters(ctx, assetFilters, templateFilters, pagination, sort)
}

// findAssetsOptimized executes the fast-path aggregation when no template filters are present.
// Pipeline: $match → $sort → $facet { metadata: [$count], data: [$skip, $limit, $lookup, $unwind, $project] }
// The $lookup moves inside $facet.data AFTER $skip/$limit, so only the page-sized slice (e.g. 20 docs)
// is joined — instead of joining ALL matching documents before paginating.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals
//   - assetFilters: Filters for assets collection (e.g., {"orgId": ObjectId, "enabled": true})
//   - pagination: Page and PerPage for pagination
//   - sort: Sort specification (e.g., {"created": -1})
//
// Returns:
//   - PaginatedResult with AssetWithTemplate entities
//   - Error if aggregation fails
func (r *repository) findAssetsOptimized(
	ctx context.Context,
	assetFilters model.Map,
	pagination *model.PaginationOpts,
	sort model.Map,
) (*model.PaginatedResult[entities.AssetWithTemplate], error) {

	pipeline := []model.Map{}

	// Stage 1: Match assets by filters (uses idx_org_created index)
	if len(assetFilters) > 0 {
		pipeline = append(pipeline, model.Map{"$match": assetFilters})
	}

	// Stage 2: Sort (default: created descending — covered by idx_org_created)
	if len(sort) > 0 {
		pipeline = append(pipeline, model.Map{"$sort": sort})
	} else {
		pipeline = append(pipeline, model.Map{"$sort": model.Map{"created": -1}})
	}

	// Stage 3: Facet — count runs on filtered set (no JOIN), data paginates then joins
	dataPipeline := []model.Map{}
	if pagination != nil {
		skip := (pagination.Page - 1) * pagination.PerPage
		dataPipeline = append(dataPipeline, model.Map{"$skip": skip})
		dataPipeline = append(dataPipeline, model.Map{"$limit": pagination.PerPage})
	}

	// $lookup + $unwind + $project run ONLY on the page-sized slice
	dataPipeline = append(dataPipeline, model.Map{
		"$lookup": model.Map{
			"from":         "assets_templates",
			"localField":   "assetTemplateId",
			"foreignField": "_id",
			"as":           "template",
		},
	})
	dataPipeline = append(dataPipeline, model.Map{
		"$unwind": model.Map{
			"path":                       "$template",
			"preserveNullAndEmptyArrays": false,
		},
	})
	dataPipeline = append(dataPipeline, templateProjectStage())

	pipeline = append(pipeline, model.Map{
		"$facet": model.Map{
			"metadata": []model.Map{
				{"$count": "total"},
			},
			"data": dataPipeline,
		},
	})

	return r.decodeFacetResults(ctx, pipeline, pagination)
}

// findAssetsWithTemplateFilters executes the full pipeline when template-level filtering is needed.
// Pipeline: $match → $lookup → $unwind → $match(template.*) → $project → $sort → $facet
// Template filters require $lookup before $match, so the JOIN runs on all matched assets.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals
//   - assetFilters: Filters for assets collection
//   - templateFilters: Filters for template fields (e.g., {"categoryId": ObjectId})
//   - pagination: Page and PerPage for pagination
//   - sort: Sort specification
//
// Returns:
//   - PaginatedResult with AssetWithTemplate entities
//   - Error if aggregation fails
func (r *repository) findAssetsWithTemplateFilters(
	ctx context.Context,
	assetFilters model.Map,
	templateFilters model.Map,
	pagination *model.PaginationOpts,
	sort model.Map,
) (*model.PaginatedResult[entities.AssetWithTemplate], error) {

	pipeline := []model.Map{}

	// Stage 1: Match assets by filters
	if len(assetFilters) > 0 {
		pipeline = append(pipeline, model.Map{"$match": assetFilters})
	}

	// Stage 2: Lookup asset_templates collection (JOIN)
	pipeline = append(pipeline, model.Map{
		"$lookup": model.Map{
			"from":         "assets_templates",
			"localField":   "assetTemplateId",
			"foreignField": "_id",
			"as":           "template",
		},
	})

	// Stage 3: Unwind template array
	pipeline = append(pipeline, model.Map{
		"$unwind": model.Map{
			"path":                       "$template",
			"preserveNullAndEmptyArrays": false,
		},
	})

	// Stage 4: Match template filters (categoryId, manufacturerId, modelId)
	templateMatch := model.Map{}
	for key, value := range templateFilters {
		templateMatch["template."+key] = value
	}
	pipeline = append(pipeline, model.Map{"$match": templateMatch})

	// Stage 5: Project fields
	pipeline = append(pipeline, templateProjectStage())

	// Stage 6: Sort
	if len(sort) > 0 {
		pipeline = append(pipeline, model.Map{"$sort": sort})
	} else {
		pipeline = append(pipeline, model.Map{"$sort": model.Map{"created": -1}})
	}

	// Stage 7: Facet — count + pagination
	dataPipeline := []model.Map{}
	if pagination != nil {
		skip := (pagination.Page - 1) * pagination.PerPage
		dataPipeline = append(dataPipeline, model.Map{"$skip": skip})
		dataPipeline = append(dataPipeline, model.Map{"$limit": pagination.PerPage})
	}

	pipeline = append(pipeline, model.Map{
		"$facet": model.Map{
			"metadata": []model.Map{
				{"$count": "total"},
			},
			"data": dataPipeline,
		},
	})

	return r.decodeFacetResults(ctx, pipeline, pagination)
}

// decodeFacetResults executes the aggregation pipeline and decodes the $facet output
// into a PaginatedResult. Shared by both fast and full pipeline paths.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals
//   - pipeline: The complete aggregation pipeline to execute
//   - pagination: Pagination options for calculating total pages
//
// Returns:
//   - PaginatedResult with AssetWithTemplate entities
//   - Error if aggregation or decoding fails
func (r *repository) decodeFacetResults(
	ctx context.Context,
	pipeline []model.Map,
	pagination *model.PaginationOpts,
) (*model.PaginatedResult[entities.AssetWithTemplate], error) {

	cursor, err := r.model.DIRECT().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var facetResults []struct {
		Metadata []struct {
			Total int64 `bson:"total"`
		} `bson:"metadata"`
		Data []entities.AssetWithTemplate `bson:"data"`
	}

	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, err
	}

	totalItems := int64(0)
	items := []entities.AssetWithTemplate{}

	if len(facetResults) > 0 {
		if len(facetResults[0].Metadata) > 0 {
			totalItems = facetResults[0].Metadata[0].Total
		}
		items = facetResults[0].Data
	}

	totalPages := int64(0)
	if pagination != nil && pagination.PerPage > 0 {
		totalPages = (totalItems + pagination.PerPage - 1) / pagination.PerPage
	}

	return &model.PaginatedResult[entities.AssetWithTemplate]{
		Items: items,
		Pagination: model.Pagination{
			Page:       pagination.Page,
			PerPage:    pagination.PerPage,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}, nil
}

// templateProjectStage returns the $project stage that flattens template classification
// data into the root document. Shared by both pipeline paths.
func templateProjectStage() model.Map {
	return model.Map{
		"$project": model.Map{
			"_id":             1,
			"name":            1,
			"enabled":         1,
			"debugEnabled":    1,
			"description":     1,
			"assetUUID":       1,
			"assetTemplateId": 1,
			"orgId":           1,
			"pathKey":         1,
			"customerId":      1,
			"routeGroupIds":   1,
			"protocol":        1,
			"latitude":        1,
			"longitude":       1,
			"created":         1,
			"updated":         1,

			// Active MQTT device cert metadata — needed so the list view
			// can render the "no certificate" warning chip only when a
			// cert-mode asset is actually missing its current cert.
			// Without this the aggregation drops it and every cert-mode
			// row looks unconfigured even after IssueCert succeeded.
			"currentCert": 1,

			// Health monitoring fields — needed for list enrichment. Without these
			// the aggregation drops them, forcing the API to surface empty status
			// and disabling the Redis enrichment (which is guarded by healthMonitor).
			"healthMonitor":         1,
			"healthStatus":          1,
			"healthStatusChangedAt": 1,

			// Flatten template classification data
			"categoryId":       "$template.categoryId",
			"categoryName":     "$template.categoryName",
			"manufacturerId":   "$template.manufacturerId",
			"manufacturerName": "$template.manufacturerName",
			"modelId":          "$template.modelId",
			"modelName":        "$template.modelName",
			"version":          "$template.version",
		},
	}
}

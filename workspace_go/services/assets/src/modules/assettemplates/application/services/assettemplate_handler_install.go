package services

import (
	ctx "context"
	"errors"
	"strings"

	"assets/src/modules/assettemplates/application/dtos"
	"assets/src/modules/assettemplates/application/ports"
	"assets/src/modules/assettemplates/domain/entities"

	v1lists "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/lists"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// resolveInstallOrg validates that the request carries an org context and returns
// the caller's org id (ObjectId + string form) and pathKey. Install is always
// org-scoped, so a missing org context is a 400.
func (s *AssetTemplateService) resolveInstallOrg(rc *reqCtx.RequestContext) (model.ObjectId, string, *string, error) {
	var zeroOrg model.ObjectId
	if rc == nil || rc.OrgContext == nil || *rc.OrgContext == "" {
		return zeroOrg, "", nil, &customErrors.ServerCustomError{Code: httpStatus.BAD_REQUEST, Errors: []string{"an organization context is required to install a template"}}
	}
	orgID, err := model.ToObjectID(*rc.OrgContext)
	if err != nil {
		return zeroOrg, "", nil, &customErrors.ServerCustomError{Code: httpStatus.BAD_REQUEST, Errors: []string{"invalid organization context"}}
	}
	var pathKey *string
	if rc.OrgContextData != nil && rc.OrgContextData.PathKey != "" {
		pk := rc.OrgContextData.PathKey
		pathKey = &pk
	}
	return orgID, *rc.OrgContext, pathKey, nil
}

// fetchAndVerifyBundle fetches the bundle and hard-verifies its sha256 BEFORE any
// resolve/cache/create — a tampered bundle is rejected regardless of whether the
// template already exists locally. A marketplace 404 maps to a 404.
func (s *AssetTemplateService) fetchAndVerifyBundle(c ctx.Context, vendor, slug string) (*ports.MarketplaceBundleFetch, error) {
	fetch, err := s.deps.MarketplaceClient.FetchBundle(c, vendor, slug)
	if err != nil {
		if errors.Is(err, ports.ErrTemplateNotFound) {
			return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"marketplace template not found"}}
		}
		return nil, err
	}
	if err := s.verifyBundleChecksum(fetch.RawBytes, fetch.DeclaredSha256); err != nil {
		return nil, err
	}
	return fetch, nil
}

// resolveClassification turns the bundle's stable slugs into org-scoped list ids
// (manufacturer, category, model — the model parented on the resolved
// manufacturer), creating each list when absent. An empty slug yields a nil id.
func (s *AssetTemplateService) resolveClassification(c ctx.Context, fetch *ports.MarketplaceBundleFetch, orgId string, pathKey *string, shareWithChildren bool) (categoryId, manufacturerId, modelId *model.ObjectId, err error) {
	b := fetch.Bundle

	manufacturerId, err = s.resolveClassificationOne(c, "asset_manufacturer", b.ManufacturerSlug, b.ManufacturerName, nil, orgId, pathKey, shareWithChildren)
	if err != nil {
		return nil, nil, nil, err
	}
	categoryId, err = s.resolveClassificationOne(c, "asset_category", b.CategorySlug, b.CategoryName, nil, orgId, pathKey, shareWithChildren)
	if err != nil {
		return nil, nil, nil, err
	}
	var parentId *string
	if manufacturerId != nil {
		p := manufacturerId.Hex()
		parentId = &p
	}
	modelId, err = s.resolveClassificationOne(c, "asset_model", b.ModelSlug, b.ModelName, parentId, orgId, pathKey, shareWithChildren)
	if err != nil {
		return nil, nil, nil, err
	}
	return categoryId, manufacturerId, modelId, nil
}

// resolveClassificationOne resolves a single classification slug into an
// org-scoped list id, or nil when the slug is empty.
func (s *AssetTemplateService) resolveClassificationOne(c ctx.Context, listType, slug, name string, parentId *string, orgId string, pathKey *string, shareWithChildren bool) (*model.ObjectId, error) {
	if slug == "" {
		return nil, nil
	}
	id, err := s.deps.ListsClient.Resolve(c, v1lists.ListResolveRequest{
		Type:              listType,
		Slug:              slug,
		Name:              name,
		ParentId:          parentId,
		OrgId:             orgId,
		PathKey:           pathKey,
		ShareWithChildren: shareWithChildren,
	})
	if err != nil {
		return nil, err
	}
	oid, err := model.ToObjectID(id)
	if err != nil {
		return nil, err
	}
	return &oid, nil
}

// buildInstalledLink builds the caller-org link record: only the light metadata
// (name/description/version/classification names) + resolved classification ids +
// marketplace identity + org scope. The heavy body (scripts, dynamic/available
// fields) is NOT stored here — it lives once in the shared content cache and is
// hydrated on read, so N orgs never duplicate it.
func (s *AssetTemplateService) buildInstalledLink(fetch *ports.MarketplaceBundleFetch, orgID model.ObjectId, pathKey *string, shareWithChildren bool, categoryId, manufacturerId, modelId *model.ObjectId) *entities.Assettemplate {
	b := fetch.Bundle
	guid := fetch.MarketplaceGuid
	sha := fetch.DeclaredSha256

	entity := &entities.Assettemplate{
		Name:            localizedValue(b.Name),
		Enabled:         true,
		IsSystem:        false,
		IsTemplate:      shareWithChildren,
		OrgID:           &orgID,
		PathKey:         pathKey,
		MarketplaceGuid: &guid,
		Sha256:          &sha,
		CategoryId:      categoryId,
		ManufacturerId:  manufacturerId,
		ModelId:         modelId,
		AssetIDPath:     b.AssetIDPath,
	}
	if desc := localizedValue(b.Description); desc != "" {
		entity.Description = &desc
	}
	if b.CategoryName != "" {
		cn := b.CategoryName
		entity.CategoryName = &cn
	}
	if b.ManufacturerName != "" {
		mn := b.ManufacturerName
		entity.ManufacturerName = &mn
	}
	if b.ModelName != "" {
		mdn := b.ModelName
		entity.ModelName = &mdn
	}
	if b.Version != "" {
		v := b.Version
		entity.Version = &v
	}
	return entity
}

// upsertInstalledLink reuses the caller-org's existing link (idempotent install)
// or creates it. A duplicate-key race (a concurrent install created the same
// org+guid first) is treated as success by re-fetching.
func (s *AssetTemplateService) upsertInstalledLink(c ctx.Context, orgID model.ObjectId, marketplaceGuid string, link *entities.Assettemplate) (*entities.Assettemplate, error) {
	existing, err := s.deps.AssetTemplateRepo.FindByMarketplaceGuidAndOrg(c, marketplaceGuid, orgID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	created, err := s.deps.AssetTemplateRepo.Create(c, link)
	if err != nil {
		if isDuplicateKey(err) {
			return s.deps.AssetTemplateRepo.FindByMarketplaceGuidAndOrg(c, marketplaceGuid, orgID)
		}
		return nil, err
	}
	return created, nil
}

// installResponse maps the link record to the response DTO and stamps its origin.
func (s *AssetTemplateService) installResponse(record *entities.Assettemplate) *dtos.AssetTemplateResponse {
	resp, _ := mapper.EntityToDto[entities.Assettemplate, dtos.AssetTemplateResponse](record)
	source := "marketplace"
	resp.Source = &source
	return resp
}

// localizedValue picks the en-US string from a locale map, falling back to any
// available value, then empty.
func localizedValue(m map[string]string) string {
	if m == nil {
		return ""
	}
	if v, ok := m["en-US"]; ok && v != "" {
		return v
	}
	for _, v := range m {
		if v != "" {
			return v
		}
	}
	return ""
}

// isDuplicateKey reports a MongoDB unique-index violation via the E11000 marker.
func isDuplicateKey(err error) bool {
	return err != nil && strings.Contains(err.Error(), "E11000")
}

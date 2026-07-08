package services

import (
	ctx "context"
	"encoding/json"

	"assets/src/modules/assettemplates/application/constants"
	"assets/src/modules/assettemplates/application/dtos"
)

// contentCacheKey builds the tiered-cache key for a marketplace template's shared
// content, keyed by marketplaceGuid.
func contentCacheKey(marketplaceGuid string) string {
	return constants.AssetTemplateContentKeyPrefix + marketplaceGuid
}

// cacheSharedContent stores the raw bundle bytes (the shared, immutable content)
// in the tiered cache under marketplaceGuid, so every org that installs the same
// template hydrates one shared copy instead of duplicating the heavy body per
// organization.
func (s *AssetTemplateService) cacheSharedContent(c ctx.Context, marketplaceGuid string, rawBytes []byte) error {
	return s.deps.TieredCache.Set(c, contentCacheKey(marketplaceGuid), rawBytes, constants.AssetTemplateContentTTL)
}

// getCachedContent reads and decodes the shared bundle content for a
// marketplaceGuid, returning (nil, nil) on a cache miss so the caller can decide
// to re-fetch from the marketplace.
func (s *AssetTemplateService) getCachedContent(c ctx.Context, marketplaceGuid string) (*dtos.MarketplaceBundle, error) {
	data, tier, err := s.deps.TieredCache.Get(c, contentCacheKey(marketplaceGuid))
	if err != nil {
		return nil, err
	}
	if tier < 0 || len(data) == 0 {
		return nil, nil
	}

	var bundle dtos.MarketplaceBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return nil, err
	}
	return &bundle, nil
}

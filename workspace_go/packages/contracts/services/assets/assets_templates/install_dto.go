package assetstemplate

// InstallParams are the path params for the marketplace install/uninstall
// routes: POST/DELETE /:marketplaceVendor/:marketplaceSlug/install.
type InstallParams struct {
	MarketplaceVendor string `params:"marketplaceVendor" validate:"required"`
	MarketplaceSlug   string `params:"marketplaceSlug" validate:"required"`
}

// InstallBody is the install request body: the caller's only choice is whether
// descendant organizations inherit the installed template.
type InstallBody struct {
	ShareWithChildren bool `json:"shareWithChildren"`
}

// MarketplaceBundle is the JSON body the marketplace bundle endpoint serves and
// the install client decodes. Classification travels as stable slugs + display
// names — the install resolves the slugs into org-scoped list ids; the original
// placeholder ObjectIds are never published. marketplaceGuid/sha256 arrive as
// response headers, read separately by the client, not fields of this body.
type MarketplaceBundle struct {
	Name             map[string]string         `json:"name"`
	Description      map[string]string         `json:"description"`
	CategorySlug     string                    `json:"categorySlug"`
	CategoryName     string                    `json:"categoryName"`
	ManufacturerSlug string                    `json:"manufacturerSlug"`
	ManufacturerName string                    `json:"manufacturerName"`
	ModelSlug        string                    `json:"modelSlug"`
	ModelName        string                    `json:"modelName"`
	Version          string                    `json:"version"`
	AssetIDPath      string                    `json:"assetIdPath"`
	ScriptTest       *string                   `json:"scriptTest,omitempty"`
	ScriptProcessor  *string                   `json:"scriptProcessor,omitempty"`
	ScriptValidator  string                    `json:"scriptValidator"`
	ScriptConversion string                    `json:"scriptConversion"`
	AvailableFields  []string                  `json:"availableFields"`
	DynamicFields    []MarketplaceDynamicField `json:"dynamicFields"`
	NextFieldId      uint16                    `json:"nextFieldId"`
}

// MarketplaceDynamicField is one entry of MarketplaceBundle.DynamicFields — the
// wire shape decoded from the external bundle, mapped to the persisted entity's
// dynamic field by the install service.
type MarketplaceDynamicField struct {
	FieldId uint16 `json:"fieldId"`
	Field   string `json:"field"`
	Value   string `json:"value"`
	Type    string `json:"type"`
	Status  uint8  `json:"status"`
}

// MarketplaceInstalledCheckRequest is the batch installed-check body: the
// marketplace GUIDs to test against the caller org's installs in one call.
type MarketplaceInstalledCheckRequest struct {
	MarketplaceGuids []string `json:"marketplaceGuids" validate:"required"`
}

// MarketplaceInstalledCheckResponse returns the subset of the requested GUIDs
// the caller org has installed.
type MarketplaceInstalledCheckResponse struct {
	Installed []string `json:"installed"`
}

// ErrCodeTemplateInUse is the machine-readable error code returned with HTTP 403
// when an uninstall is refused because assets still reference the template.
const ErrCodeTemplateInUse = "TEMPLATE_IN_USE"

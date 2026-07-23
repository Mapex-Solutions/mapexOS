package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets_templates"
)

type (
	AssetTemplateCreateDTO = v1.AssetTemplateCreate
	AssetTemplateUpdateDTO = v1.AssetTemplateUpdate
	AssetTemplateIdDto     = v1.AssetTemplateId
	AssetTemplateQueryDto  = v1.AssetTemplateQuery
	AssetTemplateResponse  = v1.AssetTemplateResponse
	DynamicField           = v1.DynamicField

	FieldVocabularyQuery    = v1.FieldVocabularyQuery
	FieldVocabularyField    = v1.FieldVocabularyField
	FieldVocabularyGroup    = v1.FieldVocabularyGroup
	FieldVocabularyResponse = v1.FieldVocabularyResponse

	InstallParams           = v1.InstallParams
	InstallBody             = v1.InstallBody
	CloneBody               = v1.CloneBody
	MarketplaceBundle       = v1.MarketplaceBundle
	MarketplaceDynamicField = v1.MarketplaceDynamicField

	MarketplaceInstalledCheckRequest  = v1.MarketplaceInstalledCheckRequest
	MarketplaceInstalledCheckResponse = v1.MarketplaceInstalledCheckResponse
)

// ErrCodeTemplateInUse is surfaced (with HTTP 403) when an uninstall is refused
// because assets still reference the template.
const ErrCodeTemplateInUse = v1.ErrCodeTemplateInUse

// ErrCodeTemplateReadonly is surfaced (with HTTP 403) when an edit is refused
// because the template came from the marketplace; the client clones it to edit.
const ErrCodeTemplateReadonly = v1.ErrCodeTemplateReadonly

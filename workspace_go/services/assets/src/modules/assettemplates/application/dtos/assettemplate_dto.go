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
	MarketplaceBundle       = v1.MarketplaceBundle
	MarketplaceDynamicField = v1.MarketplaceDynamicField
)

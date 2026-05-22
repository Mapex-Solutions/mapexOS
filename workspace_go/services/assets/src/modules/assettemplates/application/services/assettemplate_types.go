package services

import (
	"assets/src/modules/assettemplates/application/di"
	"assets/src/modules/assettemplates/domain/entities"
)

// AssetTemplateService provides methods for managing assettemplate-related operations.
// It serves as an application service layer that interacts with the
// AssetTemplateRepository to perform domain-level actions on Assettemplate entities.
//
// This service implements the AssetTemplateServicePort interface, following
// Hexagonal Architecture principles by depending on interfaces rather
// than concrete implementations.
type AssetTemplateService struct {
	deps di.AssetTemplateServiceDependenciesInjection
}

// DynamicFieldsResult holds processed dynamic fields and updated NextFieldId
type DynamicFieldsResult struct {
	Fields      []entities.DynamicField
	NextFieldId uint16
}

package collection

import (
	"mapexVault/src/modules/credentials/domain/entities"
	"mapexVault/src/modules/credentials/domain/repositories"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// Compile-time checks
var _ repositories.CredentialRepository = (*credentialRepository)(nil)
var _ repositories.ConnectionRepository = (*connectionRepository)(nil)

type credentialRepository struct {
	model *model.Model[entities.Credential]
}

type connectionRepository struct {
	model *model.Model[entities.Connection]
}

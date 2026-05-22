package routes

import (
	"mapexVault/src/modules/credentials/application/dtos"
	"mapexVault/src/modules/credentials/application/ports"
	"mapexVault/src/modules/credentials/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/vault"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
)

// RegisterRoutes registers external credential API routes (JWT auth).
func RegisterRoutes(router fiber.Router, service ports.CredentialServicePort) {
	createDto := validation.NewValidation(&dtos.CreateCredentialDTO{}, nil, nil)
	router.Post("/",
		validation.ValidationMiddleware(createDto),
		permissionMw.RequirePermission(perms.CredentialCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateCredential(service),
	)

	queryDto := validation.NewValidation(nil, &dtos.CredentialQueryDTO{}, nil)
	router.Get("/",
		validation.ValidationMiddleware(queryDto),
		permissionMw.RequirePermission(perms.CredentialRead),
		coverageMw.InjectRequestContext(),
		handlers.GetCredentials(service),
	)

	router.Get("/:credentialId",
		permissionMw.RequirePermission(perms.CredentialRead),
		coverageMw.InjectRequestContext(),
		handlers.GetCredentialById(service),
	)

	updateDto := validation.NewValidation(&dtos.UpdateCredentialDTO{}, nil, nil)
	router.Patch("/:credentialId",
		validation.ValidationMiddleware(updateDto),
		permissionMw.RequirePermission(perms.CredentialUpdate),
		coverageMw.InjectRequestContext(),
		handlers.UpdateCredentialById(service),
	)

	router.Delete("/:credentialId",
		permissionMw.RequirePermission(perms.CredentialDelete),
		coverageMw.InjectRequestContext(),
		handlers.DeleteCredentialById(service),
	)

	router.Post("/:credentialId/test",
		permissionMw.RequirePermission(perms.CredentialRead),
		coverageMw.InjectRequestContext(),
		handlers.TestCredential(service),
	)
}

// RegisterInternalRoutes registers internal credential API routes (API key auth).
func RegisterInternalRoutes(router fiber.Router, service ports.CredentialServicePort) {
	router.Get("/:credentialId/decrypt", handlers.DecryptCredential(service))
}

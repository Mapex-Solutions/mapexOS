package routes

import (
	"mapexVault/src/modules/credentials/application/dtos"
	"mapexVault/src/modules/credentials/application/ports"
	"mapexVault/src/modules/credentials/interfaces/http/handlers"

	perms "github.com/Mapex-Solutions/MapexOS/permissions/vault"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers external credential API routes (JWT auth).
//
// Routes are registered through the swagger wrapper: each NewValidation declares
// the input contract once (used to both validate and document), the module tag is
// declared once on Wrap, and each route's summary, description, and response type
// are attached via the fluent builder.
func RegisterRoutes(router web.Router, service ports.CredentialServicePort) {

	r := swagger.Wrap(router).Tag("Credentials")

	// Create a new credential.
	createDto := validation.NewValidation(&dtos.CreateCredentialDTO{}, nil, nil)
	r.Post("/", createDto, swagger.Expose,
		permissionMw.RequirePermission(perms.CredentialCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreateCredential(service),
	).
		Summary("Create credential").
		Description("Creates and envelope-encrypts a new credential. Org scoping is applied from the request context.").
		Returns(&dtos.CredentialResponse{})

	// List credentials with filters and pagination.
	queryDto := validation.NewValidation(nil, &dtos.CredentialQueryDTO{}, nil)
	r.Get("/", queryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.CredentialRead),
		coverageMw.InjectRequestContext(),
		handlers.GetCredentials(service),
	).
		Summary("List credentials").
		Description("Returns a paginated, filterable list of credentials scoped to the caller's organization. Secret material is never returned.").
		Returns(&model.PaginatedResult[dtos.CredentialResponse]{})

	// Get credential by ID.
	getByIdDto := validation.NewValidation(nil, nil, nil)
	r.Get("/:credentialId", getByIdDto, swagger.Expose,
		permissionMw.RequirePermission(perms.CredentialRead),
		coverageMw.InjectRequestContext(),
		handlers.GetCredentialById(service),
	).
		Summary("Get credential by ID").
		Description("Retrieves a single credential's metadata by its ID. Secret material is never returned.").
		Returns(&dtos.CredentialResponse{})

	// Update credential by ID.
	updateDto := validation.NewValidation(&dtos.UpdateCredentialDTO{}, nil, nil)
	r.Patch("/:credentialId", updateDto, swagger.Expose,
		permissionMw.RequirePermission(perms.CredentialUpdate),
		coverageMw.InjectRequestContext(),
		handlers.UpdateCredentialById(service),
	).
		Summary("Update credential").
		Description("Partially updates an existing credential. Only provided fields are changed.").
		Returns(&dtos.CredentialResponse{})

	// Delete credential by ID.
	deleteDto := validation.NewValidation(nil, nil, nil)
	r.Delete("/:credentialId", deleteDto, swagger.Expose,
		permissionMw.RequirePermission(perms.CredentialDelete),
		coverageMw.InjectRequestContext(),
		handlers.DeleteCredentialById(service),
	).
		Summary("Delete credential").
		Description("Deletes a credential by its ID.").
		Returns(map[string]bool{})

	// Test a credential's connectivity.
	testDto := validation.NewValidation(nil, nil, nil)
	r.Post("/:credentialId/test", testDto, swagger.Expose,
		permissionMw.RequirePermission(perms.CredentialRead),
		coverageMw.InjectRequestContext(),
		handlers.TestCredential(service),
	).
		Summary("Test credential").
		Description("Decrypts the credential and verifies it can authenticate against its provider.").
		Returns(map[string]bool{})
}

// RegisterInternalRoutes registers internal credential API routes (API key auth).
// These are MS-to-MS endpoints and are kept out of the public OpenAPI document.
func RegisterInternalRoutes(router web.Router, service ports.CredentialServicePort) {

	r := swagger.Wrap(router)

	// Decrypt a credential for an internal caller (never exposed publicly).
	decryptDto := validation.NewValidation(nil, nil, nil)
	r.Get("/:credentialId/decrypt", decryptDto, swagger.Hidden,
		handlers.DecryptCredential(service),
	)
}

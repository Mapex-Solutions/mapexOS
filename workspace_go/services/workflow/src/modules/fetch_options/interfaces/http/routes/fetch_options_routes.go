package routes

import (
	"workflow/src/modules/fetch_options/application/ports"
	"workflow/src/modules/fetch_options/application/types"
	"workflow/src/modules/fetch_options/interfaces/http/handlers"

	perms "github.com/Mapex-Solutions/MapexOS/permissions/workflow"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers fetch-options HTTP routes.
//
// Routes:
//
//	POST / - Resolve dynamic options for a plugin/node configuration field
func RegisterRoutes(router web.Router, service ports.FetchOptionsServicePort) {

	r := swagger.Wrap(router).Tag("Fetch Options")

	// Resolve dynamic options. The body (credentialId, pluginId, resourceKey,
	// dependsOn) is parsed in the handler, so it is not declared as a contract here.
	fetchOptionsDto := validation.NewValidation(nil, nil, nil)
	r.Post("/", fetchOptionsDto, swagger.Expose,
		permissionMw.RequirePermission(perms.CredentialRead),
		coverageMw.InjectRequestContext(),
		handlers.FetchOptions(service),
	).
		Summary("Fetch options").
		Description("Decrypts the referenced credential, loads the plugin manifest, and proxies a call to the provider to resolve dynamic options for a configuration field. Body: credentialId, pluginId, resourceKey, dependsOn.").
		Returns(&[]types.FetchOptionsItem{})
}

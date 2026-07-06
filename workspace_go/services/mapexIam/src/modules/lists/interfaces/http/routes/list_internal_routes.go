package routes

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"mapexIam/src/modules/lists/application/dtos"
	"mapexIam/src/modules/lists/application/ports"
	"mapexIam/src/modules/lists/interfaces/http/handlers"

	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
)

// RegisterInternalListRoutes registers the service-to-service list routes under
// /internal/lists. The group is API-key protected in module.go; these routes are
// only ever called by other services, never end users.
//
// Endpoints:
//   - POST /internal/lists/resolve - resolve-or-create an org-scoped list by slug
func RegisterInternalListRoutes(group web.Router, service ports.ListServicePort) {
	resolveDto := validation.NewValidation(&dtos.ListResolveRequest{}, nil, nil)
	group.Post("/resolve", validation.ValidationMiddleware(resolveDto), handlers.ResolveList(service))
}

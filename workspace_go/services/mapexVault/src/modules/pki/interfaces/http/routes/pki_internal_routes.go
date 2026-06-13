package routes

import (
	"mapexVault/src/modules/pki/interfaces/http/handlers"

	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterInternalRoutes mounts the pki internal endpoints under /internal/pki.
// The router passed in MUST already be API-key gated by the caller (module.go
// wires apikeymw on the group). These are MS-to-MS endpoints, kept out of the
// public OpenAPI document via swagger.Hidden.
func RegisterInternalRoutes(router web.Router, h *handlers.PkiInternalHandler) {
	r := swagger.Wrap(router)

	dto := validation.NewValidation(nil, nil, nil)
	r.Get("/intermediate_ca_bundle", dto, swagger.Hidden, h.GetIntermediateCABundle)
	r.Get("/ca_chain", dto, swagger.Hidden, h.GetCAChain)
	r.Post("/sign_server", dto, swagger.Hidden, h.SignServer)
}

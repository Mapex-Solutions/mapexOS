package routes

import (
	"mapexVault/src/modules/kek/interfaces/http/handlers"

	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterInternalRoutes mounts the kek internal endpoints under /internal/kek.
// The router passed in MUST already be API-key gated by the caller (module.go
// wires apikeymw on the group). These are MS-to-MS endpoints, kept out of the
// public OpenAPI document via swagger.Hidden.
func RegisterInternalRoutes(router web.Router, h *handlers.KekInternalHandler) {
	r := swagger.Wrap(router)

	dto := validation.NewValidation(nil, nil, nil)
	r.Get("/:context", dto, swagger.Hidden, h.GetKEK)
}

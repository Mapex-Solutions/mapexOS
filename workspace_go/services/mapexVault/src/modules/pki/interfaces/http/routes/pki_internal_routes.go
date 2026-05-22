package routes

import (
	"mapexVault/src/modules/pki/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

// RegisterInternalRoutes mounts the pki internal endpoints under
// /internal/pki. The router passed in MUST already be API-key gated by
// the caller (module.go wires apikeymw on the group).
func RegisterInternalRoutes(router fiber.Router, h *handlers.PkiInternalHandler) {
	router.Get("/intermediate_ca_bundle", h.GetIntermediateCABundle)
	router.Get("/ca_chain", h.GetCAChain)
	router.Post("/sign_server", h.SignServer)
}

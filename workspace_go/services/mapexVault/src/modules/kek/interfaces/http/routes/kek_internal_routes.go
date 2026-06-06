package routes

import (
	"mapexVault/src/modules/kek/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

// RegisterInternalRoutes mounts the kek internal endpoints under /internal/kek.
// The router passed in MUST already be API-key gated by the caller (module.go
// wires apikeymw on the group).
func RegisterInternalRoutes(router fiber.Router, h *handlers.KekInternalHandler) {
	router.Get("/:context", h.GetKEK)
}

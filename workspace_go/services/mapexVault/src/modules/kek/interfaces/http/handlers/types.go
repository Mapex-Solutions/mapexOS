package handlers

import (
	"errors"
	"strings"

	kekPorts "mapexVault/src/modules/kek/application/ports"
	service "mapexVault/src/modules/kek/application/services"

	"github.com/gofiber/fiber/v2"
)

// KekInternalHandler bundles the service port for the internal endpoints.
type KekInternalHandler struct {
	service kekPorts.KekServicePort
}

// NewKekInternalHandler constructs the handler.
func NewKekInternalHandler(s kekPorts.KekServicePort) *KekInternalHandler {
	return &KekInternalHandler{service: s}
}

// mapServiceErrToStatus translates service-layer errors into HTTP status codes.
// The caller's boot retry loop relies on this mapping: 503 = "KEK not seeded
// yet, keep retrying"; 500 = "real failure, surface it".
func mapServiceErrToStatus(err error) int {
	if err == nil {
		return fiber.StatusOK
	}
	if errors.Is(err, service.ErrKEKNotSeeded) {
		return fiber.StatusServiceUnavailable
	}
	// The repository may wrap mongo's "no documents" as a plain error chain that
	// does not implement Is/As for the sentinel; sniff the message so the boot
	// race maps to 503 instead of leaking as 500.
	msg := err.Error()
	if strings.Contains(msg, "document not found") || strings.Contains(msg, "no documents") {
		return fiber.StatusServiceUnavailable
	}
	return fiber.StatusInternalServerError
}

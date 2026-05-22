package handlers

import (
	"errors"
	"strings"

	pkiPorts "mapexVault/src/modules/pki/application/ports"
	service "mapexVault/src/modules/pki/application/services"

	"github.com/gofiber/fiber/v2"
)

// PkiInternalHandler bundles the service port for the internal endpoints.
type PkiInternalHandler struct {
	service pkiPorts.PkiServicePort
}

// NewPkiInternalHandler constructs the handler.
func NewPkiInternalHandler(s pkiPorts.PkiServicePort) *PkiInternalHandler {
	return &PkiInternalHandler{service: s}
}

// mapServiceErrToStatus translates service-layer errors into HTTP
// status codes. The Assets MS bootstrap retry loop relies on this
// mapping: 503 = "CA not yet bootstrapped, keep retrying"; 500 =
// "real failure, surface it". Any 502 here would leak the wrong
// semantics (BadGateway is for proxy/upstream failures, not for "the
// thing you asked for doesn't exist yet").
func mapServiceErrToStatus(err error) int {
	if err == nil {
		return fiber.StatusOK
	}
	if errors.Is(err, service.ErrCANotBootstrapped) {
		return fiber.StatusServiceUnavailable
	}
	// The repository wraps mongo's "no documents" as a plain error
	// chain that may not implement Is/As for the original sentinel.
	// Fall back to substring sniffing so the OnMount race
	// (handler hits the route before bootstrap completes) maps to
	// 503 instead of leaking as 500.
	msg := err.Error()
	if strings.Contains(msg, "document not found") || strings.Contains(msg, "no documents") {
		return fiber.StatusServiceUnavailable
	}
	return fiber.StatusInternalServerError
}

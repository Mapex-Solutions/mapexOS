package handlers

import (
	"mapexVault/src/modules/credentials/application/ports"

	"github.com/gofiber/fiber/v2"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

func DecryptCredential(service ports.CredentialServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("credentialId")
		data, err := service.DecryptCredential(c.UserContext(), id)
		if err != nil {
			return err
		}
		return response.Success(c, data)
	}
}

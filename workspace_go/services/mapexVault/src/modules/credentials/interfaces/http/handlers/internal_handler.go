package handlers

import (
	"mapexVault/src/modules/credentials/application/ports"

	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

func DecryptCredential(service ports.CredentialServicePort) web.Handler {
	return func(c *web.Ctx) error {
		id := c.Params("credentialId")
		data, err := service.DecryptCredential(c.UserContext(), id)
		if err != nil {
			return err
		}
		return response.Success(c, data)
	}
}

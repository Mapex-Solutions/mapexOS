package handlers

import (
	"mapexVault/src/modules/credentials/application/dtos"
	"mapexVault/src/modules/credentials/application/ports"

	"github.com/gofiber/fiber/v2"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

func CreateCredential(service ports.CredentialServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		rc, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found", nil)
		}
		dto, _ := requestValidation.GetDTO[*dtos.CreateCredentialDTO](c, "bodyDTO")
		result, err := service.CreateCredential(ctx, rc, dto)
		if err != nil {
			return err
		}
		return response.Created(c, result)
	}
}

func GetCredentialById(service ports.CredentialServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("credentialId")
		result, err := service.GetCredentialById(c.UserContext(), id)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}
}

func UpdateCredentialById(service ports.CredentialServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("credentialId")
		dto, _ := requestValidation.GetDTO[*dtos.UpdateCredentialDTO](c, "bodyDTO")
		result, err := service.UpdateCredentialById(c.UserContext(), id, dto)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}
}

func DeleteCredentialById(service ports.CredentialServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("credentialId")
		result, err := service.DeleteCredentialById(c.UserContext(), id)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}
}

func GetCredentials(service ports.CredentialServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rc, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found", nil)
		}
		query, _ := requestValidation.GetDTO[*dtos.CredentialQueryDTO](c, "queryDTO")
		result, err := service.GetCredentials(c.UserContext(), rc, query)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}
}

func TestCredential(service ports.CredentialServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("credentialId")
		result, err := service.TestCredential(c.UserContext(), id)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}
}

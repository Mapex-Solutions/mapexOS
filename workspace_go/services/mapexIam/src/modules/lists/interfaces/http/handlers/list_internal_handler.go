package handlers

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"mapexIam/src/modules/lists/application/dtos"
	"mapexIam/src/modules/lists/application/ports"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// ResolveList returns a service-to-service handler that resolves a classification
// slug into an org-scoped list id, creating the list when absent. The org is
// carried in the request body — this is an internal call with no user context.
func ResolveList(service ports.ListServicePort) web.Handler {
	return func(c *web.Ctx) error {
		req, _ := requestValidation.GetDTO[*dtos.ListResolveRequest](c, "bodyDTO")
		resp, err := service.ResolveOrCreateList(c.UserContext(), *req)
		if err != nil {
			return err
		}
		return response.Success(c, resp)
	}
}

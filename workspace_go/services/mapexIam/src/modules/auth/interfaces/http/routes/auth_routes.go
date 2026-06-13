package routes

import (
	"mapexIam/src/modules/auth/application/dtos"
	"mapexIam/src/modules/auth/application/ports"
	"mapexIam/src/modules/auth/interfaces/http/handlers"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	refremw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/refreshTokenExtractor"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers the auth HTTP routes. Base path: /auth.
//
// Following Hexagonal Architecture, this function accepts the service port interface
// rather than a concrete service implementation. Routes are registered through the
// swagger wrapper.
func RegisterRoutes(group web.Router, service ports.AuthServicePort) {

	r := swagger.Wrap(group).Tag("Auth")

	loginDtos := validation.NewValidation(&dtos.LoginDTO{}, nil, nil)
	r.Post("/login", loginDtos, swagger.Expose, handlers.Login(service)).
		Summary("Log in").
		Description("Authenticates a user with email and password and returns an access token bundle.").
		Returns(map[string]interface{}{})

	noContract := validation.NewValidation(nil, nil, nil)
	r.Post("/logout", noContract, swagger.Expose,
		authmw.AuthMiddleware(config.GetAuthConfig()),
		handlers.Logout(service),
	).
		Summary("Log out").
		Description("Invalidates the caller's current session/token.").
		Returns(map[string]interface{}{})

	r.Post("/refresh", noContract, swagger.Expose,
		refremw.RefreshTokenExtractor(),
		handlers.RefreshToken(service),
	).
		Summary("Refresh access token").
		Description("Exchanges a valid refresh token for a new access token bundle.").
		Returns(map[string]interface{}{})

	r.Get("/users/me/coverage", noContract, swagger.Expose,
		authmw.AuthMiddleware(config.GetAuthConfig()),
		handlers.GetMyCoverage(service),
	).
		Summary("Get my authorization coverage").
		Description("Returns the authenticated caller's authorization coverage (organizations and scopes the caller can act on).").
		Returns(map[string]interface{}{})

	r.Get("/me/permissions", noContract, swagger.Expose,
		authmw.AuthMiddleware(config.GetAuthConfig()),
		handlers.GetMyPermissions(service),
	).
		Summary("Get my permissions").
		Description("Returns the authenticated caller's effective permission set for the current organization context.").
		Returns(map[string]interface{}{})
}

package routes

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/auth/application/dtos"
	"mapexIam/src/modules/auth/application/ports"
	"mapexIam/src/modules/auth/interfaces/http/handlers"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	authmw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	refremw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/refreshTokenExtractor"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
)

func RegisterRoutes(group fiber.Router, service ports.AuthServicePort) {

	// Log in
	loginDtos := validation.NewValidation(&dtos.LoginDTO{}, nil, nil)
	group.Post("/login", validation.ValidationMiddleware(loginDtos), handlers.Login(service))

	// Log out
	group.Post("/logout", authmw.AuthMiddleware(config.GetAuthConfig()), handlers.Logout(service))

	// Refresh token
	group.Post(
		"/refresh",
		refremw.RefreshTokenExtractor(),
		handlers.RefreshToken(service),
	)

	// Get my coverage (organizations accessible to current user)
	group.Get(
		"/users/me/coverage",
		authmw.AuthMiddleware(config.GetAuthConfig()),
		handlers.GetMyCoverage(service),
	)

	// Get my permissions (resolved permissions for current user in current org)
	group.Get(
		"/me/permissions",
		authmw.AuthMiddleware(config.GetAuthConfig()),
		handlers.GetMyPermissions(service),
	)
}

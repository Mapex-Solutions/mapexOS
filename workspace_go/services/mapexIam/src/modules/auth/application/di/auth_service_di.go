package di

import (
	"go.uber.org/dig"

	"mapexIam/src/modules/auth/domain/repositories"

	membershipPorts "mapexIam/src/modules/memberships/application/ports"
	rolePorts "mapexIam/src/modules/roles/application/ports"
	userPorts "mapexIam/src/modules/users/application/ports"

	middlewaresAuth "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
)

type AuthServiceDI struct {
	dig.In

	// Same domain repositories
	Repo                  repositories.AuthRepository
	SessionRepo           repositories.SessionRepository
	CoverageCacheRepo     repositories.CoverageCacheRepository
	AuthCacheRepo         repositories.AuthorizationCacheRepository

	// Auth Config
	AuthConfig middlewaresAuth.AuthConfig

	// Other domains (using ports for testability and decoupling)
	UserService       userPorts.UserServicePort
	MembershipService membershipPorts.MembershipServicePort
	RoleService       rolePorts.RoleServicePort
}

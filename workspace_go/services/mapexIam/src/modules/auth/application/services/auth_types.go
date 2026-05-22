package services

import (
	"mapexIam/src/modules/auth/application/di"
)

// AuthService coordinates authentication use cases such as login and
// refresh-token issuance. It depends on other application services and
// infrastructure via the AuthServiceDI container.
type AuthService struct {
	di di.AuthServiceDI
}

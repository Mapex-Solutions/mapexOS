package repositories

import (
	"context"
	"mapexIam/src/modules/auth/domain/entities"
)

type AuthRepository interface {
	Login(ctx context.Context, u *entities.Auth) (*entities.Auth, error)
}

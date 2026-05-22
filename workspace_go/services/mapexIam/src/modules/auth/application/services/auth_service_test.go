package services_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"mapexIam/src/modules/auth/application/di"
	"mapexIam/src/modules/auth/application/dtos"
	authServices "mapexIam/src/modules/auth/application/services"
	authRepos "mapexIam/src/modules/auth/domain/repositories"
	membershipPorts "mapexIam/src/modules/memberships/application/ports"
	rolePorts "mapexIam/src/modules/roles/application/ports"
	userDtos "mapexIam/src/modules/users/application/dtos"
	userPorts "mapexIam/src/modules/users/application/ports"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	middlewaresAuth "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	"github.com/Mapex-Solutions/mapexGoKit/utils/typeconv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
 * Inline port mocks — per /go-arch §14, tests use stdlib + inline mocks of
 * port interfaces. Each fake exposes function fields for the methods the
 * tests configure; unused ports are stubbed with no-op methods so the DI
 * struct stays fully wired without dragging in testify.
 */

type fakeUserService struct {
	getUserByEmailFn func(ctx context.Context, email *string) (*userPorts.User, error)
}

func (f *fakeUserService) CreateUser(_ context.Context, _ *userDtos.UserCreateDTO) (*userDtos.UserResponse, error) {
	return nil, nil
}
func (f *fakeUserService) GetUserById(_ context.Context, _ *string) (*userDtos.UserResponse, error) {
	return nil, nil
}
func (f *fakeUserService) UpdateUserById(_ context.Context, _ *string, _ *userDtos.UserUpdateDTO) (*userDtos.UserResponse, error) {
	return nil, nil
}
func (f *fakeUserService) DeleteUserById(_ context.Context, _ *string) (map[string]bool, error) {
	return nil, nil
}
func (f *fakeUserService) GetUserByEmail(ctx context.Context, email *string) (*userPorts.User, error) {
	if f.getUserByEmailFn != nil {
		return f.getUserByEmailFn(ctx, email)
	}
	return nil, nil
}
func (f *fakeUserService) GetUsers(_ context.Context, _ *reqCtx.RequestContext, _ *userDtos.UserQueryDto) (*model.PaginatedResult[userDtos.UserResponse], error) {
	return nil, nil
}
func (f *fakeUserService) CountUsers(_ context.Context, _ *reqCtx.RequestContext) (int64, error) {
	return 0, nil
}

type fakeSessionRepo struct {
	storeRefreshTokenFn func(ctx context.Context, userId, sessionId, refreshToken string, ttl time.Duration) error
}

func (f *fakeSessionRepo) StoreRefreshToken(ctx context.Context, userId, sessionId, refreshToken string, ttl time.Duration) error {
	if f.storeRefreshTokenFn != nil {
		return f.storeRefreshTokenFn(ctx, userId, sessionId, refreshToken, ttl)
	}
	return nil
}
func (f *fakeSessionRepo) GetRefreshToken(_ context.Context, _, _ string) (string, error) {
	return "", nil
}
func (f *fakeSessionRepo) InvalidateRefreshToken(_ context.Context, _, _ string) error { return nil }

/*
 * Test fixtures
 */

func createTestUser(email string, enabled bool) *userPorts.User {
	objectId := primitive.NewObjectID()
	// Real bcrypt hash of "validPassword".
	hashedPassword := "$2a$10$A8mC2Zwvq5Rv/TcWcQeif.oDYAHab/MQw8ppn5jv7roA82yG8EY9a"
	return &userPorts.User{
		ID:        model.ObjectId(objectId),
		Email:     email,
		Password:  &hashedPassword,
		FirstName: "Test",
		LastName:  "User",
		Enabled:   enabled,
	}
}

// silence unused-import warnings on optional helper.
var _ = typeconv.PtrString

func buildAuthDI(userService userPorts.UserServicePort, sessionRepo authRepos.SessionRepository) di.AuthServiceDI {
	return di.AuthServiceDI{
		UserService:       userService,
		SessionRepo:       sessionRepo,
		CoverageCacheRepo: nil,
		AuthCacheRepo:     nil,
		Repo:              nil,
		MembershipService: membershipPorts.MembershipServicePort(nil),
		RoleService:       rolePorts.RoleServicePort(nil),
		AuthConfig: middlewaresAuth.AuthConfig{
			Secret: "test-secret-key-for-jwt-signing",
		},
	}
}

/*
 * New (constructor)
 */

func TestAuthService_New(t *testing.T) {
	authDI := di.AuthServiceDI{
		AuthConfig: middlewaresAuth.AuthConfig{Secret: "test-secret"},
	}

	service := authServices.New(authDI)
	if service == nil {
		t.Fatal("expected non-nil service")
	}
}

/*
 * Login
 */

func TestAuthService_Login_Success(t *testing.T) {
	ctx := context.Background()
	email := "test@example.com"
	password := "validPassword"
	testUser := createTestUser(email, true)

	userService := &fakeUserService{
		getUserByEmailFn: func(_ context.Context, _ *string) (*userPorts.User, error) {
			return testUser, nil
		},
	}
	storeCalls := 0
	sessionRepo := &fakeSessionRepo{
		storeRefreshTokenFn: func(_ context.Context, _, _, _ string, _ time.Duration) error {
			storeCalls++
			return nil
		},
	}

	authService := authServices.New(buildAuthDI(userService, sessionRepo))
	loginDto := &dtos.LoginDTO{Email: email, Password: password}

	result, err := authService.Login(ctx, loginDto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map[string]interface{}, got %T", result)
	}
	if resultMap["access_token"] == "" || resultMap["access_token"] == nil {
		t.Fatal("expected non-empty access_token")
	}
	if resultMap["refresh_token"] == "" || resultMap["refresh_token"] == nil {
		t.Fatal("expected non-empty refresh_token")
	}
	if resultMap["user"] == nil {
		t.Fatal("expected non-nil user")
	}
	if storeCalls != 1 {
		t.Fatalf("expected StoreRefreshToken called exactly once, got %d", storeCalls)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	ctx := context.Background()
	email := "nonexistent@example.com"

	userService := &fakeUserService{
		getUserByEmailFn: func(_ context.Context, _ *string) (*userPorts.User, error) {
			return nil, nil
		},
	}
	authService := authServices.New(buildAuthDI(userService, &fakeSessionRepo{}))

	result, err := authService.Login(ctx, &dtos.LoginDTO{Email: email, Password: "password"})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatalf("expected nil result, got %#v", result)
	}
	if !strings.Contains(err.Error(), "invalid credentials") {
		t.Fatalf("expected error to contain 'invalid credentials', got: %v", err)
	}
}

func TestAuthService_Login_UserDisabled(t *testing.T) {
	ctx := context.Background()
	email := "disabled@example.com"
	testUser := createTestUser(email, false)

	userService := &fakeUserService{
		getUserByEmailFn: func(_ context.Context, _ *string) (*userPorts.User, error) {
			return testUser, nil
		},
	}
	authService := authServices.New(buildAuthDI(userService, &fakeSessionRepo{}))

	result, err := authService.Login(ctx, &dtos.LoginDTO{Email: email, Password: "password"})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatalf("expected nil result, got %#v", result)
	}
	if !strings.Contains(err.Error(), "user blocked") {
		t.Fatalf("expected error to contain 'user blocked', got: %v", err)
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	ctx := context.Background()
	email := "test@example.com"
	testUser := createTestUser(email, true) // hash is for "validPassword"

	userService := &fakeUserService{
		getUserByEmailFn: func(_ context.Context, _ *string) (*userPorts.User, error) {
			return testUser, nil
		},
	}
	authService := authServices.New(buildAuthDI(userService, &fakeSessionRepo{}))

	result, err := authService.Login(ctx, &dtos.LoginDTO{Email: email, Password: "wrongPassword"})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatalf("expected nil result, got %#v", result)
	}
	if !strings.Contains(err.Error(), "invalid credentials") {
		t.Fatalf("expected error to contain 'invalid credentials', got: %v", err)
	}
}

func TestAuthService_Login_CacheStoreSuccess(t *testing.T) {
	ctx := context.Background()
	email := "test@example.com"
	testUser := createTestUser(email, true)

	userService := &fakeUserService{
		getUserByEmailFn: func(_ context.Context, _ *string) (*userPorts.User, error) {
			return testUser, nil
		},
	}
	authService := authServices.New(buildAuthDI(userService, &fakeSessionRepo{}))

	result, err := authService.Login(ctx, &dtos.LoginDTO{Email: email, Password: "validPassword"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

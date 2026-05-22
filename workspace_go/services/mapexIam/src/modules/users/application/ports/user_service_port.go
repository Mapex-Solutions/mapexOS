package ports

import (
	"context"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"mapexIam/src/modules/users/application/dtos"
	"mapexIam/src/modules/users/domain/entities" // Still needed for GetUserByEmail
)

// UserServicePort defines the inbound port (Hexagonal Architecture) for user-related operations.
// This port allows other layers (handlers, other services) to depend on user operations
// without coupling to the concrete UserService implementation.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: This interface (defines the contract)
//   - Adapter: UserService (implements the contract)
//
// Used by:
//   - HTTP Handlers (for API endpoints)
//   - AuthService (for login and token refresh operations)
//   - Other application services requiring user operations
type UserServicePort interface {
	// CreateUser creates a new user with the provided data.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - dto: User creation data
	//
	// Returns:
	//   - *dtos.UserResponse: Created user data DTO
	//   - error: Error if creation fails
	CreateUser(ctx context.Context, dto *dtos.UserCreateDTO) (*dtos.UserResponse, error)

	// GetUserById retrieves a user by their unique identifier.
	// Returns a DTO (without sensitive fields like password).
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - userId: User's unique identifier
	//
	// Returns:
	//   - *dtos.UserResponse: User data DTO (sanitized)
	//   - error: Error if user not found or database error
	GetUserById(ctx context.Context, userId *string) (*dtos.UserResponse, error)

	// UpdateUserById updates an existing user's information.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - userId: User's unique identifier
	//   - dto: Fields to update
	//
	// Returns:
	//   - *dtos.UserResponse: Updated user data DTO
	//   - error: Error if update fails or user not found
	UpdateUserById(ctx context.Context, userId *string, dto *dtos.UserUpdateDTO) (*dtos.UserResponse, error)

	// DeleteUserById removes a user by their unique identifier.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - userId: User's unique identifier
	//
	// Returns:
	//   - map[string]bool: Success indicator
	//   - error: Error if deletion fails or user not found
	DeleteUserById(ctx context.Context, userId *string) (map[string]bool, error)

	// GetUserByEmail retrieves a user entity by email address.
	// Used primarily for authentication flows.
	// Returns full entity (includes password hash).
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - email: User's email address
	//
	// Returns:
	//   - *entities.User: Full user entity (includes password hash)
	//   - error: Error if user not found or database error
	GetUserByEmail(ctx context.Context, email *string) (*entities.User, error)

	// GetUsers retrieves a paginated and filtered list of users.
	// Uses RequestContext for context-aware organization filtering.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - requestContext: RequestContext from InjectRequestContext middleware
	//   - query: Query parameters (filters, pagination, projection)
	//
	// Returns:
	//   - *model.PaginatedResult[dtos.UserResponse]: Paginated results with DTOs
	//   - error: Error if query fails
	GetUsers(ctx context.Context, requestContext *reqCtx.RequestContext, query *dtos.UserQueryDto) (*model.PaginatedResult[dtos.UserResponse], error)

	// CountUsers returns the total count of users for the given org context.
	// Implements cache-aside: check Redis first, fallback to MongoDB CountDocuments.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - requestContext: Contains org access data from coverage middleware
	//
	// Returns:
	//   - int64: Total count of matching users
	//   - error: If query fails
	CountUsers(ctx context.Context, requestContext *reqCtx.RequestContext) (int64, error)
}

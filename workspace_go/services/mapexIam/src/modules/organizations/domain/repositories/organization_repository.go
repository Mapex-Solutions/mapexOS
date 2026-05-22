package repositories

import (
	"context"
	"mapexIam/src/modules/organizations/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type OrganizationRepository interface {
	Create(ctx context.Context, u *entities.Organization) (*entities.Organization, error)
	FindById(ctx context.Context, organizationId *string) (*entities.Organization, error)
	FindByIds(ctx context.Context, organizationIds []string) ([]*entities.Organization, error)
	FindByIdAndUpdate(ctx context.Context, organizationId *string, payload map[string]any) (*entities.Organization, error)
	DeleteById(ctx context.Context, organizationId *string) error
	FindWithFilters(
		ctx context.Context,
		filters model.Map,
		pagination *model.PaginationOpts,
		projection model.Map,
	) (*model.PaginatedResult[entities.Organization], error)
	FindWithCursor(
		ctx context.Context,
		filters model.Map,
		cursorOpts *model.CursorOpts,
		projection model.Map,
	) (*model.CursorResult[entities.Organization], error)
}

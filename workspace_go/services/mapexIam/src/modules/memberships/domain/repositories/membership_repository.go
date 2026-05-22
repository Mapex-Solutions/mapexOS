package repositories

import (
	"context"
	"mapexIam/src/modules/memberships/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type MembershipRepository interface {
	Create(ctx context.Context, u *entities.Membership) (*entities.Membership, error)
	FindById(ctx context.Context, membershipId *string) (*entities.Membership, error)
	FindByIdAndUpdate(ctx context.Context, membershipId *string, payload map[string]any) (*entities.Membership, error)
	DeleteById(ctx context.Context, membershipId *string) error
	FindByUserId(ctx context.Context, userId *string) ([]*entities.Membership, error)
	FindByGroupIds(ctx context.Context, groupIds []string) ([]*entities.Membership, error)
	FindByAssigneeAndOrg(ctx context.Context, assigneeType string, assigneeId string, orgId string) (*entities.Membership, error)
	FindWithFilters(
		ctx context.Context,
		filters model.Map,
		pagination *model.PaginationOpts,
		projection model.Map,
	) (*model.PaginatedResult[entities.Membership], error)
}

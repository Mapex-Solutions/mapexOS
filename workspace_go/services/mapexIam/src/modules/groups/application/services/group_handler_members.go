package services

import (
	ctx "context"
	"fmt"
	"time"

	"mapexIam/src/modules/groups/application/dtos"
	"mapexIam/src/modules/groups/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// resolveAddMemberInputs combines the group fetch + id parsing required by
// AddMemberToGroup so the public orchestration stays a single named-step call.
func (s *GroupService) resolveAddMemberInputs(c ctx.Context, groupId, userId string) (*entities.Group, model.ObjectId, model.ObjectId, error) {
	group, err := s.fetchGroupOr404(c, &groupId)
	if err != nil {
		return nil, model.ObjectId{}, model.ObjectId{}, err
	}
	groupObjectID, userObjectID, err := s.parseMemberIds(groupId, userId)
	if err != nil {
		return nil, model.ObjectId{}, model.ObjectId{}, err
	}
	return group, groupObjectID, userObjectID, nil
}

// parseMemberIds converts the raw groupId/userId strings into ObjectId
// values, returning BAD_REQUEST when either id is malformed.
func (s *GroupService) parseMemberIds(groupId, userId string) (model.ObjectId, model.ObjectId, error) {
	groupObjectID, err := model.ToObjectID(groupId)
	if err != nil {
		return model.ObjectId{}, model.ObjectId{}, &customErrors.ServerCustomError{
			Code:   status.BAD_REQUEST,
			Errors: []string{"Invalid group ID format"},
		}
	}
	userObjectID, err := model.ToObjectID(userId)
	if err != nil {
		return model.ObjectId{}, model.ObjectId{}, &customErrors.ServerCustomError{
			Code:   status.BAD_REQUEST,
			Errors: []string{"Invalid user ID format"},
		}
	}
	return groupObjectID, userObjectID, nil
}

// persistGroupMember inserts a row into the GroupMember junction table
// carrying the denormalized OrgID + PathKey for hierarchical range queries.
func (s *GroupService) persistGroupMember(c ctx.Context, group *entities.Group, groupObjectID, userObjectID model.ObjectId) error {
	var orgObjectID model.ObjectId
	if group.OrgID != nil {
		orgObjectID = *group.OrgID
	}
	now := time.Now()
	groupMember := &entities.GroupMember{
		GroupID: groupObjectID,
		UserID:  userObjectID,
		OrgID:   orgObjectID,
		PathKey: group.PathKey,
		AddedAt: now,
		Created: now,
		Updated: now,
	}
	if _, err := s.deps.GroupMemberRepo.Create(c, groupMember); err != nil {
		return &customErrors.ServerCustomError{
			Code:   status.INTERNAL_SERVER_ERROR,
			Errors: []string{fmt.Sprintf("Failed to add user to group %s: %v", groupObjectID.Hex(), err)},
		}
	}
	return nil
}

// cappedPaginationOpts returns the page/perPage from the query DTO with
// the per-page value capped at 100 so callers cannot blow up the response.
func (s *GroupService) cappedPaginationOpts(query *dtos.GroupMembersQueryDto) (int64, int64) {
	perPage := int64(query.GetPerPage())
	if perPage > 100 {
		perPage = 100
	}
	return int64(query.GetPage()), perPage
}

// enrichGroupMemberDetails projects each junction row into a response DTO
// and attaches the user's email + name fields fetched from the user service.
// User-fetch failures are tolerated (the member row still ships).
func (s *GroupService) enrichGroupMemberDetails(c ctx.Context, members []*entities.GroupMember) []dtos.GroupMemberResponse {
	dtoItems := make([]dtos.GroupMemberResponse, len(members))
	for i, member := range members {
		userId := member.UserID.Hex()
		dtoItems[i] = dtos.GroupMemberResponse{
			ID:      &member.ID,
			UserID:  &member.UserID,
			GroupID: &member.GroupID,
			OrgID:   &member.OrgID,
			AddedAt: &member.AddedAt,
			Created: &member.Created,
		}
		if member.AddedBy != nil {
			dtoItems[i].AddedBy = member.AddedBy
		}
		user, err := s.deps.UserService.GetUserById(c, &userId)
		if err == nil && user != nil {
			dtoItems[i].UserEmail = user.Email
			dtoItems[i].UserFirstName = user.FirstName
			dtoItems[i].UserLastName = user.LastName
		}
	}
	return dtoItems
}

// buildGroupMembersResult wraps the enriched member rows with pagination
// metadata for the response DTO.
func (s *GroupService) buildGroupMembersResult(items []dtos.GroupMemberResponse, page, perPage, totalItems int64) *model.PaginatedResult[dtos.GroupMemberResponse] {
	totalPages := (totalItems + perPage - 1) / perPage
	return &model.PaginatedResult[dtos.GroupMemberResponse]{
		Items: items,
		Pagination: model.Pagination{
			Page:       page,
			PerPage:    perPage,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}
}

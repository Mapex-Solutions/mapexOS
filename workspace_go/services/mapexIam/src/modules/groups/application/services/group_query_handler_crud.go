package services

import (
	"context"
	"fmt"

	"mapexIam/src/modules/groups/application/ports"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// getGroupBasicInfo fetches a group and projects the basic info DTO.
func (s *GroupQueryService) getGroupBasicInfo(ctx context.Context, groupId string) (*ports.GroupBasicInfo, error) {
	group, err := s.deps.Repo.FindById(ctx, &groupId)
	if err != nil {
		return nil, fmt.Errorf("failed to get group %s: %w", groupId, err)
	}
	if group == nil {
		return nil, nil
	}

	return &ports.GroupBasicInfo{
		ID:          group.ID.Hex(),
		Name:        group.Name,
		Description: group.Description,
	}, nil
}

// getAllUserGroupIds queries the GroupMember junction table by userId and
// returns the list of group IDs the user belongs to across all organizations.
func (s *GroupQueryService) getAllUserGroupIds(ctx context.Context, userId string) ([]string, error) {
	groupMembers, err := s.deps.GroupMemberRepo.FindByUserId(ctx, userId)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:GroupQuery] Failed to get group memberships for user=%s: %v", userId, err))
		return []string{}, nil
	}

	groupIds := make([]string, len(groupMembers))
	for i, gm := range groupMembers {
		groupIds[i] = gm.GroupID.Hex()
	}

	return groupIds, nil
}

// getUserIdsByGroupIds fetches members by group IDs and returns a deduplicated
// list of user IDs.
func (s *GroupQueryService) getUserIdsByGroupIds(ctx context.Context, groupIds []string) ([]string, error) {
	if len(groupIds) == 0 {
		return []string{}, nil
	}

	groupMembers, err := s.deps.GroupMemberRepo.FindByGroupIds(ctx, groupIds)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by group IDs: %w", err)
	}

	userIdMap := make(map[string]bool, len(groupMembers))
	for _, gm := range groupMembers {
		userIdMap[gm.UserID.Hex()] = true
	}

	userIds := make([]string, 0, len(userIdMap))
	for userId := range userIdMap {
		userIds = append(userIds, userId)
	}

	return userIds, nil
}

// countGroupsByUserIds returns a userId -> groupCount mapping for the given
// user IDs via individual GroupMember lookups.
func (s *GroupQueryService) countGroupsByUserIds(ctx context.Context, userIds []string) (map[string]int, error) {
	result := make(map[string]int, len(userIds))

	for _, userId := range userIds {
		result[userId] = 0
	}

	for _, userId := range userIds {
		groupMembers, err := s.deps.GroupMemberRepo.FindByUserId(ctx, userId)
		if err == nil {
			result[userId] = len(groupMembers)
		}
	}

	return result, nil
}

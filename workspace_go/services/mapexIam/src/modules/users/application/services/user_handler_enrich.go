package services

import (
	ctx "context"

	"mapexIam/src/modules/users/application/dtos"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// getUserGroups returns the groups a user belongs to and the count.
// Uses GroupQueryService (cross-domain port) instead of direct repository access.
func (s *UserService) getUserGroups(context ctx.Context, userId string) ([]dtos.UserGroupInfo, *int) {
	groups := []dtos.UserGroupInfo{}

	// Get group IDs via GroupQueryService (cross-domain port)
	groupIds, err := s.deps.GroupQueryService.GetAllUserGroupIds(context, userId)
	if err != nil || len(groupIds) == 0 {
		count := 0
		return groups, &count
	}

	// Get group details for each group via GroupQueryService
	for _, groupId := range groupIds {
		groupInfo, err := s.deps.GroupQueryService.GetGroupBasicInfo(context, groupId)
		if err != nil || groupInfo == nil {
			continue
		}

		info := dtos.UserGroupInfo{
			ID:          groupInfo.ID,
			Name:        groupInfo.Name,
			Description: groupInfo.Description,
		}
		groups = append(groups, info)
	}

	count := len(groups)
	return groups, &count
}

// getUserMemberships returns the organization memberships for a user (direct + via groups).
// Uses MembershipService and GroupQueryService (cross-domain ports) instead of direct repo access.
func (s *UserService) getUserMemberships(context ctx.Context, userId string) []dtos.UserMembershipInfo {
	memberships := []dtos.UserMembershipInfo{}

	// Track unique org memberships to avoid duplicates
	orgMembershipMap := make(map[string]dtos.UserMembershipInfo)

	// 1. Get direct user memberships via MembershipService port
	directMemberships, err := s.deps.MembershipService.GetDirectUserMemberships(context, userId)
	if err == nil {
		for _, membership := range directMemberships {
			orgId := membership.OrgID.Hex()

			// Get org details
			org, err := s.deps.OrgService.GetOrganizationById(context, &orgId)
			if err != nil || org == nil {
				continue
			}

			// Get role names
			roleNames := s.getRoleNames(context, membership.RoleIds)

			info := dtos.UserMembershipInfo{
				OrgID:     orgId,
				OrgName:   *org.Name,
				OrgType:   *org.Type,
				Scope:     membership.Scope,
				RoleNames: roleNames,
				Via:       "direct",
			}

			// Direct membership takes precedence
			orgMembershipMap[orgId] = info
		}
	}

	// 2. Get memberships via groups using GroupQueryService + MembershipService ports
	groupIds, err := s.deps.GroupQueryService.GetAllUserGroupIds(context, userId)
	if err == nil && len(groupIds) > 0 {
		// Build group name map for "via" field
		groupNameMap := make(map[string]string)
		for _, groupId := range groupIds {
			groupInfo, err := s.deps.GroupQueryService.GetGroupBasicInfo(context, groupId)
			if err == nil && groupInfo != nil {
				groupNameMap[groupId] = groupInfo.Name
			}
		}

		// Get group memberships via MembershipService port
		groupMemberships, err := s.deps.MembershipService.GetMembershipsByGroupIds(context, groupIds)
		if err == nil {
			for _, membership := range groupMemberships {
				orgId := membership.OrgID.Hex()

				// Skip if direct membership already exists for this org
				if _, exists := orgMembershipMap[orgId]; exists {
					continue
				}

				// Get org details
				org, err := s.deps.OrgService.GetOrganizationById(context, &orgId)
				if err != nil || org == nil {
					continue
				}

				// Get role names
				roleNames := s.getRoleNames(context, membership.RoleIds)

				// Get group name for "via" field
				groupName := groupNameMap[membership.AssigneeID.Hex()]
				via := "Group: " + groupName

				info := dtos.UserMembershipInfo{
					OrgID:     orgId,
					OrgName:   *org.Name,
					OrgType:   *org.Type,
					Scope:     membership.Scope,
					RoleNames: roleNames,
					Via:       via,
				}

				orgMembershipMap[orgId] = info
			}
		}
	}

	// Convert map to slice
	for _, info := range orgMembershipMap {
		memberships = append(memberships, info)
	}

	return memberships
}

// getRoleNames returns the names of roles given their IDs
func (s *UserService) getRoleNames(context ctx.Context, roleIds []model.ObjectId) []string {
	names := make([]string, 0, len(roleIds))

	for _, roleId := range roleIds {
		roleIdStr := roleId.Hex()
		role, err := s.deps.RoleService.GetRoleById(context, &roleIdStr)
		if err != nil || role == nil {
			continue
		}
		if role.Name != nil {
			names = append(names, *role.Name)
		}
	}

	return names
}

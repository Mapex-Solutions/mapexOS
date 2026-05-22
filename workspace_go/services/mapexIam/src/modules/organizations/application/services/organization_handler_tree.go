package services

import (
	ctx "context"
	"fmt"

	events "mapexIam/src/modules/cache_invalidation/application/events"
	"mapexIam/src/modules/organizations/domain/entities"

	contractsDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/organizations"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// collectAncestorOrgIds walks the parent chain from the given org upwards,
// collecting all ancestor organization IDs until it reaches the root.
func (s *OrganizationService) collectAncestorOrgIds(start *entities.Organization) []string {
	ancestorOrgIds := []string{}
	current := start
	for current != nil {
		ancestorOrgIds = append(ancestorOrgIds, current.ID.Hex())
		if current.ParentOrgID != nil {
			parentIdHex := current.ParentOrgID.Hex()
			current, _ = s.deps.Repo.FindById(ctx.Background(), &parentIdHex)
		} else {
			current = nil
		}
	}
	return ancestorOrgIds
}

// publishOrgHierarchyChangedEvent publishes a NATS event to invalidate
// coverage cache for users with recursive memberships on ancestor orgs.
// Runs asynchronously — walks the parent chain, marshals, and publishes.
func (s *OrganizationService) publishOrgHierarchyChangedEvent(orgID string, parentOrg *entities.Organization, action string) {
	go func() {
		ancestorOrgIds := s.collectAncestorOrgIds(parentOrg)
		if len(ancestorOrgIds) == 0 {
			return
		}

		event := events.NewOrgHierarchyChangedEvent(orgID, ancestorOrgIds, action, "")

		subject := fmt.Sprintf("mapexos.cache.invalidation.organization.%s.hierarchy.changed", orgID)
		if err := s.deps.NatsBus.Publish(natsModel.PublishConfig{Subject: subject, Data: event}); err != nil {
			logger.Error(err, fmt.Sprintf("[SERVICE:Organization] Failed to publish OrgHierarchyChangedEvent for org=%s", orgID))
			return
		}
		logger.Info(fmt.Sprintf("[SERVICE:Organization] Published OrgHierarchyChangedEvent (action=%s) for org=%s with %d ancestors", action, orgID, len(ancestorOrgIds)))
	}()
}

// publishOrgCreatedEvent publishes a NATS event when an organization is
// created. Consumed by the Events service to auto-create default
// retention policies. Runs asynchronously.
func (s *OrganizationService) publishOrgCreatedEvent(org *entities.Organization) {
	go func() {
		event := contractsDtos.OrganizationCreatedEvent{
			OrgId:   org.ID.Hex(),
			PathKey: org.PathKey,
			Name:    org.Name,
			Type:    org.Type,
		}

		subject := "mapexos.events.organization.created"
		if err := s.deps.NatsBus.Publish(natsModel.PublishConfig{Subject: subject, Data: event}); err != nil {
			logger.Error(err, fmt.Sprintf("[SERVICE:Organization] Failed to publish OrganizationCreatedEvent for org=%s", org.ID.Hex()))
			return
		}
		logger.Info(fmt.Sprintf("[SERVICE:Organization] Published OrganizationCreatedEvent for org=%s pathKey=%s", org.ID.Hex(), org.PathKey))
	}()
}

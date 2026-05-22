package ports

import "router/src/modules/routegroups/domain/entities"

// Port-level type aliases — expose domain entities through the port boundary.
// Other modules import these types from ports, NEVER from domain/entities directly.

type RouteGroup = entities.RouteGroup
type Router = entities.Router
type MatchConfig = entities.MatchConfig
type MatchRule = entities.MatchRule
type LakeHouseData = entities.LakeHouseData
type NotificationData = entities.NotificationData
type TriggerData = entities.TriggerData
type SaveEventData = entities.SaveEventData
type WorkflowData = entities.WorkflowData

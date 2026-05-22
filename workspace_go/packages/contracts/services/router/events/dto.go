package events

import (
	"fmt"

	"github.com/Mapex-Solutions/MapexOS/contracts/common"
)

type RouteGroupId struct {
	RouteGroupId string `params:"routeGroupId" validate:"required,mongoid"`
}

type LakeHouseData struct {
	LakeHouseId string                 `json:"lakeHouseId" validate:"required,mongoid"`
}

type NotificationData struct {
	NotificationId string                 `json:"notificationId" validate:"required,mongoid"`
}

type SaveEventData struct {
}

type Router struct {
	Kind         string            `json:"kind" validate:"required,oneof=lake_house notification save_event"`
	LakeHouse    *LakeHouseData    `json:"lakeHouse,omitempty"`
	Notification *NotificationData `json:"notification,omitempty"`
	SaveEvent    *SaveEventData    `json:"saveEvent,omitempty"`
}

type RouteGroupCreate struct {
	Version     string   `json:"version" validate:"required,semver"`
	Name        string   `json:"name" validate:"required,min=1"`
	Description string   `json:"description" validate:"required,min=1"`
	Enabled     bool     `json:"enabled"`
	OrgId       string   `json:"orgId" validate:"required,min=3"`
	Routers     []Router `json:"routers" validate:"required,min=1,dive"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

type RouteGroupUpdate struct {
	Version     *string   `json:"version,omitempty" validate:"omitempty,semver"`
	Name        *string   `json:"name,omitempty" validate:"omitempty,min=1"`
	Description *string   `json:"description,omitempty" validate:"omitempty,min=1"`
	Enabled     *bool     `json:"enabled,omitempty"`
	OrgId       *string   `json:"orgId,omitempty" validate:"omitempty,min=3"`
	Routers     *[]Router `json:"routers,omitempty" validate:"omitempty,min=1,dive"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

type RouteGroupResponse struct {
	ID          *common.ObjectID `json:"id,omitempty"`
	Version     *string          `json:"version,omitempty"`
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	Enabled     *bool            `json:"enabled,omitempty"`
	OrgId       *string          `json:"orgId,omitempty"`
	Routers     *[]Router        `json:"routers,omitempty"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

func (r *RouteGroupResponse) SetCreated(t *common.NullTime) { r.Created = t }
func (r *RouteGroupResponse) SetUpdated(t *common.NullTime) { r.Updated = t }

/** TRANSFORMATIONS **/

func (r *Router) Transform() error {
	switch r.Kind {
	case "save_event":
		return nil


	case "lake_house":
		if r.LakeHouse == nil {
			return fmt.Errorf("field 'lakeHouse' must be provided when kind is 'lake_house'")
		}

	case "notification":
		if r.Notification == nil {
			return fmt.Errorf("field 'notification' must be provided when kind is 'notification'")
		}

	default:
		return fmt.Errorf("invalid router kind: %s", r.Kind)
	}

	return nil
}

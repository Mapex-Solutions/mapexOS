package dtos

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// OrganizationCreatedEvent is the NATS event published when a new organization is created.
// Consumed by the Events service to auto-create default retention policies.
type OrganizationCreatedEvent struct {
	// OrgId is the hex string of the new organization's ObjectID
	OrgId string `json:"orgId"`

	// PathKey is the hierarchical path key of the organization
	PathKey string `json:"pathKey"`

	// Name is the organization name
	Name string `json:"name"`

	// Type is the organization type (vendor, customer, site, building, floor, zone)
	Type string `json:"type"`
}

// SubjectOrganizationCreated is the NATS subject on which mapexIam
// publishes OrganizationCreatedEvent. Consumed by the events service
// retention module to auto-create default retention policies for the
// newly created organization. Resolved at package init —
// e.g. "dev.mapexos.events.organization.created".
var SubjectOrganizationCreated = config.Subject("events", "organization.created")

// StreamOrganizationEvents is the NATS JetStream stream that carries
// organization lifecycle events emitted by mapexIam. Consumed by the
// events service retention module (org_created consumer). Resolved at
// package init — e.g. "DEV-MAPEXOS-MAPEXIAM-ORGANIZATIONS".
var StreamOrganizationEvents = config.StreamName("MAPEXIAM", "ORGANIZATIONS")

// EventTypeOrganizationCreated tags DLQ messages produced by the
// retention/org_created consumer in the events service.
const EventTypeOrganizationCreated = "organization.created"

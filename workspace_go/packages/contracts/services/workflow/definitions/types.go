package definitions

/*
 * CROSS-SERVICE NATS MESSAGE TYPES
 * Payloads published by the workflow service and consumed by other services
 * (Go or TypeScript). Keep field names and JSON tags in sync with
 * workspace_js/packages/schemas/src/services/workflow/definitions/.
 */

// DefinitionInvalidatePayload is the FANOUT message payload for workflow
// definition cache invalidation. Published when code nodes are created,
// updated, or deleted in a definition. Consumers (e.g. js-workflow-executor)
// invalidate L0 (RAM) + L1 (Disk) for the specified nodeIds.
type DefinitionInvalidatePayload struct {
	OrgId        string   `json:"orgId"`
	DefinitionId string   `json:"definitionId"`
	NodeIds      []string `json:"nodeIds"`
}

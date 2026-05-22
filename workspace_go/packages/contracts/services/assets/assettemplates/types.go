package assettemplates

// ListNameUpdatedEvent is the payload published by the mapexos (lists) service
// on the ListNameUpdatedSubject whenever a list's name changes (manufacturer,
// model, or category). It is consumed by the assets service to propagate the
// renamed value into denormalized AssetTemplate documents.
//
// Ownership: published by the mapexos service.
// Consumers: assets service (assettemplates module).
//
// Reciprocity: mirrored by workspace_js/packages/schemas/src/services/assets/assettemplates.
type ListNameUpdatedEvent struct {
	ListId   string `json:"listId"`
	ListType string `json:"listType"`
	NewName  string `json:"newName"`
	OrgId    string `json:"orgId"`
}

// TemplateInvalidatePayload is the cross-service wire contract for the FANOUT subject mapexos.fanout.template.invalidate. Published by the assets service (assettemplates module); consumed by router, events, and js-executor.
type TemplateInvalidatePayload struct {
	OrgId      string `json:"orgId"`
	TemplateId string `json:"templateId"`
}

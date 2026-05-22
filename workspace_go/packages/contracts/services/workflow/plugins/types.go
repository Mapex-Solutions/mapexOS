package plugins

// PluginInvalidatePayload is the cross-service FANOUT wire contract for mapexos.fanout.workflow.plugin.invalidate. Published by the workflow service plugins module; consumed by workflow self-fanout (and any future js-workflow-executor subscription).
type PluginInvalidatePayload struct {
	PluginID string `json:"pluginId"`
	Action   string `json:"action"` // "create", "update", "delete"
}

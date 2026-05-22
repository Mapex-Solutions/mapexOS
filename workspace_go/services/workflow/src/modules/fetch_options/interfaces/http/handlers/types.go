package handlers

/*
 * FETCH OPTIONS HANDLER TYPES
 * HTTP request/response structs for the load_options endpoint.
 */

// FetchOptionsRequest is the POST body for the load_options endpoint.
type FetchOptionsRequest struct {
	CredentialId string            `json:"credentialId" validate:"required"`
	PluginId     string            `json:"pluginId,omitempty"`
	ResourceKey  string            `json:"resourceKey" validate:"required"`
	DependsOn    map[string]string `json:"dependsOn,omitempty"`
}

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	pluginPorts "workflow/src/modules/plugins/application/ports"
	"workflow/src/modules/fetch_options/application/types"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/jsrunner"
)

// resolveCredentialAndPlugin decrypts the credential via vault and pulls the
// pluginId from the metadata when the caller didn't supply one. Strips
// internal __pluginId/__credentialDefId keys so the rest of the pipeline
// only sees real credential data for template substitution.
func (s *FetchOptionsService) resolveCredentialAndPlugin(ctx context.Context, credentialId, pluginId string) (map[string]interface{}, string, error) {
	credData, err := s.deps.Vault.DecryptCredential(ctx, credentialId)
	if err != nil {
		logger.Error(err, "[SERVICE:FetchOptions] Failed to decrypt credential")
		return nil, "", &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"Credential not found or decrypt failed"}}
	}
	if pluginId == "" {
		if pid, ok := credData["__pluginId"].(string); ok {
			pluginId = pid
		}
	}
	delete(credData, "__pluginId")
	delete(credData, "__credentialDefId")
	return credData, pluginId, nil
}

// resolveFetchOptionsDef looks up the plugin manifest and its FetchOptions
// entry for the given resource key. Surfaces 404 when the manifest or key
// is missing, 400 when the entry is not http-typed.
func (s *FetchOptionsService) resolveFetchOptionsDef(ctx context.Context, pluginId, resourceKey string) (*pluginPorts.PluginManifest, *pluginPorts.FetchOptionsDef, error) {
	manifest, err := s.deps.PluginRepo.FindByPluginId(ctx, pluginId)
	if err != nil || manifest == nil {
		logger.Error(fmt.Errorf("manifest not found: %s", pluginId), "[SERVICE:FetchOptions] Manifest not found")
		return nil, nil, &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"Plugin manifest not found for " + pluginId}}
	}
	fetchOptDef, ok := manifest.FetchOptions[resourceKey]
	if !ok {
		logger.Error(fmt.Errorf("key not found: %s", resourceKey), "[SERVICE:FetchOptions] FetchOptions key not in manifest")
		return nil, nil, &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"FetchOptions key '" + resourceKey + "' not found in plugin " + pluginId}}
	}
	if fetchOptDef.Type != "http" || fetchOptDef.Http == nil {
		return nil, nil, &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{"FetchOptions must be of type http"}}
	}
	return manifest, &fetchOptDef, nil
}

// buildFetchRequest assembles the HTTP request: resolve {{credentials.*}}
// and {{dependsOn.*}} placeholders in path + headers, default to GET when
// the manifest doesn't pin a method, and stamp the headers on the request.
func (s *FetchOptionsService) buildFetchRequest(ctx context.Context, manifest *pluginPorts.PluginManifest, def *pluginPorts.FetchOptionsDef, credData map[string]interface{}, dependsOn map[string]string) (*http.Request, error) {
	path := resolveFetchTemplates(def.Http.Path, credData, dependsOn)
	method := def.Http.Method
	if method == "" {
		method = "GET"
	}
	url := manifest.Defaults.BaseUrl + path
	logger.Debug(fmt.Sprintf("[SERVICE:FetchOptions] %s %s", method, url))
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:FetchOptions] Failed to build request: %w", err)
	}
	for k, v := range def.Http.Headers {
		req.Header.Set(k, resolveFetchTemplates(v, credData, dependsOn))
	}
	return req, nil
}

// executeFetchRequest performs the outbound HTTP call with a 15s timeout
// and decodes the JSON body. Non-2xx responses are surfaced as 400 contract
// errors carrying the upstream body so the caller can debug.
func (s *FetchOptionsService) executeFetchRequest(req *http.Request) (interface{}, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:FetchOptions] HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:FetchOptions] Failed to read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &customErrors.ServerCustomError{
			Code:   status.BAD_REQUEST,
			Errors: []string{fmt.Sprintf("External API returned %d: %s", resp.StatusCode, string(body))},
		}
	}
	var responseData interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, fmt.Errorf("[SERVICE:FetchOptions] Failed to parse response: %w", err)
	}
	return responseData, nil
}

// extractFetchItems projects the raw response into the {label,value} item
// list. Two modes: a JS transform script (when present) runs the response
// through ES5; otherwise dataPath/valuePath/labelPath drive a generic
// extraction over the JSON tree.
func (s *FetchOptionsService) extractFetchItems(ctx context.Context, responseData interface{}, def *pluginPorts.FetchOptionsDef) ([]types.FetchOptionsItem, error) {
	transformScript, dataPath, valuePath, labelPath := "", "", "", ""
	if def.Output != nil {
		transformScript = def.Output.Transform
		dataPath = def.Output.DataPath
		valuePath = def.Output.ValuePath
		labelPath = def.Output.LabelPath
	}
	if transformScript != "" {
		jsResults, err := jsrunner.RunTransform(ctx, transformScript, responseData)
		if err != nil {
			return nil, fmt.Errorf("[SERVICE:FetchOptions] Transform script failed: %w", err)
		}
		items := make([]types.FetchOptionsItem, 0, len(jsResults))
		for _, r := range jsResults {
			items = append(items, types.FetchOptionsItem{Label: r.Label, Value: r.Value})
		}
		logger.Debug(fmt.Sprintf("[SERVICE:FetchOptions] Returning %d items", len(items)))
		return items, nil
	}
	items, err := extractOptions(responseData, dataPath, valuePath, labelPath)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE:FetchOptions] Failed to extract options: %w", err)
	}
	logger.Debug(fmt.Sprintf("[SERVICE:FetchOptions] Returning %d items", len(items)))
	return items, nil
}

// resolveFetchTemplates substitutes {{credentials.*}} and {{dependsOn.*}}
// placeholders in a string. Non-string credential values are skipped.
func resolveFetchTemplates(input string, credData map[string]interface{}, dependsOn map[string]string) string {
	out := input
	for key, val := range credData {
		if strVal, ok := val.(string); ok {
			out = strings.ReplaceAll(out, "{{credentials."+key+"}}", strVal)
		}
	}
	for key, val := range dependsOn {
		out = strings.ReplaceAll(out, "{{dependsOn."+key+"}}", val)
	}
	return out
}

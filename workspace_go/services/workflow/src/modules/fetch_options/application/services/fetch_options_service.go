package services

import (
	"context"
	"fmt"

	"workflow/src/modules/fetch_options/application/di"
	"workflow/src/modules/fetch_options/application/ports"
	"workflow/src/modules/fetch_options/application/types"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// Compile-time check to ensure FetchOptionsService implements FetchOptionsServicePort.
var _ ports.FetchOptionsServicePort = (*FetchOptionsService)(nil)

// New creates a new FetchOptionsService.
func New(deps di.FetchOptionsServiceDI) ports.FetchOptionsServicePort {
	return &FetchOptionsService{
		deps: deps,
	}
}

// FetchOptions executes the fetchOptions proxy flow:
// decrypt credential -> resolve manifest entry -> build templated HTTP
// request -> call the external API -> transform/extract the response into
// {label,value} items. Returns the items slice or a contract error on any
// failure (vault miss, manifest miss, HTTP non-2xx, transform error).
func (s *FetchOptionsService) FetchOptions(ctx context.Context, credentialId string, pluginId string, resourceKey string, dependsOn map[string]string) ([]types.FetchOptionsItem, error) {
	logger.Debug(fmt.Sprintf("[SERVICE:FetchOptions] Called: credentialId=%s, pluginId=%s, resourceKey=%s", credentialId, pluginId, resourceKey))
	credData, pluginId, err := s.resolveCredentialAndPlugin(ctx, credentialId, pluginId)
	if err != nil {
		return nil, err
	}
	manifest, fetchOptDef, err := s.resolveFetchOptionsDef(ctx, pluginId, resourceKey)
	if err != nil {
		return nil, err
	}
	req, err := s.buildFetchRequest(ctx, manifest, fetchOptDef, credData, dependsOn)
	if err != nil {
		return nil, err
	}
	responseData, err := s.executeFetchRequest(req)
	if err != nil {
		return nil, err
	}
	return s.extractFetchItems(ctx, responseData, fetchOptDef)
}

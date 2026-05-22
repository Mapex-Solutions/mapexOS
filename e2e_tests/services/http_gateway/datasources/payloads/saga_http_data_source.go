// Package payloads holds canonical DataSourceCreate fixtures for the
// http_gateway datasources module.
//
// The HTTP-protocol connectivity-action journey provisions a push-mode
// HTTP DataSource so the asset can issue explicit heartbeats and
// telemetry events against /api/v1/heartbeat and /api/v1/events with
// the same auth surface a production device would use.
package payloads

import (
	"fmt"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/http_gateway/datasources"
)

// SagaApiKeyHeaderName is the request header the saga's HTTP datasource
// expects the apiKey on. Kept as a constant so the saga's heartbeat /
// publish steps reference the exact same string without re-declaring it.
const SagaApiKeyHeaderName = "X-API-Key"

// DataSourceCreateBuilder wraps contracts.DataSourceCreate so journeys
// can override individual fields without redeclaring the canonical
// baseline.
type DataSourceCreateBuilder struct {
	spec   contracts.DataSourceCreate
	apiKey string
}

// Build returns the contracts payload ready for POST /api/v1/data_sources.
func (b *DataSourceCreateBuilder) Build() contracts.DataSourceCreate { return b.spec }

// ApiKey returns the plaintext apiKey the builder embedded so the
// CreateDataSource step can publish it on the bag for the heartbeat /
// publish steps to present in the request header.
func (b *DataSourceCreateBuilder) ApiKey() string { return b.apiKey }

// SagaHttpDataSource returns the canonical push-mode HTTP datasource
// the connectivity-action HTTP journey uses. The DataSource exposes an
// apiKey auth surface (header X-API-Key) and binds events to the
// assetUUID field on the request body — matching the heartbeat
// contract { "assetUUID": "..." }.
//
// Inputs:
//   - runID  saga run identifier embedded in the name so multiple
//            concurrent runs on the same Mongo do not collide on the
//            name unique-by-(orgId, name) index.
//
// Defaults:
//   - Mode:     push (device pushes to /api/v1/events and /api/v1/heartbeat)
//   - Protocol: http
//   - Auth:     apiKey via header X-API-Key, key derived from runID and
//               long enough to pass the contract's min=20 validator.
//   - Bind:     uuidField on path "assetUUID" — the heartbeat body
//               key the http_gateway resolves against the assets read
//               model.
func SagaHttpDataSource(runID string) *DataSourceCreateBuilder {
	apiKey := fmt.Sprintf("saga-apikey-%s-padded-1234567890", runID)
	spec := contracts.DataSourceCreate{
		Name:    fmt.Sprintf("saga-http-ds-%s", runID),
		Enabled: true,
		Mode:    "push",
		Protocol: "http",
		Auth: contracts.DataSourceAuth{
			Type: contracts.AuthTypeAPIKey,
			APIKey: &contracts.AuthAPIKey{
				Type:      "header",
				FieldName: SagaApiKeyHeaderName,
				Key:       apiKey,
			},
		},
	}
	spec.AssetBind.Type = "uuidField"
	spec.AssetBind.Data.UUIDField = []string{"assetUUID"}
	return &DataSourceCreateBuilder{spec: spec, apiKey: apiKey}
}

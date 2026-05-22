package gateways

/**
 * GATEWAY MESSAGE CONTRACT
 *
 * Common contract for all gateways (HTTP, MQTT, LoRaWAN) to send messages
 * to the JS-Executor service for script processing.
 *
 * This ensures consistency across all gateway types and simplifies
 * the JS-Executor's message handling.
 *
 * PUBLISHERS:
 *   - HTTP Gateway: Converts DataSourceResponse → GatewayMessage
 *   - MQTT Gateway: Converts AssetReadModel → GatewayMessage
 *   - LoRaWAN Gateway: Converts AssetReadModel → GatewayMessage
 *
 * CONSUMER:
 *   - JS-Executor: Receives GatewayMessage on "processor.js.execute" subject
 */

// SourceType identifies the origin gateway
type SourceType string

const (
	SourceTypeHTTP    SourceType = "http"
	SourceTypeMQTT    SourceType = "mqtt"
	SourceTypeLoRaWAN SourceType = "lorawan"
)

// GatewayMessage is the common contract for all gateways
// sending messages to the JS-Executor for script processing.
//
// Fields:
//   - SourceType: Origin gateway (http, mqtt, lorawan)
//   - DataSource: Normalized source configuration with tenant info
//   - Event: Raw payload from device/API
//   - AssetUUID: Pre-resolved for MQTT/LoRaWAN, resolved by JS-Executor for HTTP
type GatewayMessage struct {
	SourceType SourceType `json:"sourceType"`
	DataSource DataSource `json:"dataSource"`
	Event      any        `json:"event"`
	AssetUUID  string     `json:"assetUUID,omitempty"`
}

// DataSource contains the normalized source configuration.
//
// For HTTP Gateway: Populated from HTTP DataSource entity
// For MQTT/LoRaWAN: Populated from Asset Read Model
//
// Fields:
//   - ID: Source identifier (DataSource ID or Asset ID)
//   - Name: Human-readable name
//   - Enabled: Whether processing is enabled
//   - DebugEnabled: If true, raw events are logged to ClickHouse
//   - Description: Optional description
//   - OrgId: Organization ID for tenant isolation
//   - PathKey: Hierarchical path for range queries
type DataSource struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Enabled      bool   `json:"enabled"`
	DebugEnabled bool   `json:"debugEnabled"`
	Description  string `json:"description,omitempty"`
	OrgId        string `json:"orgId"`
	PathKey      string `json:"pathKey"`
}

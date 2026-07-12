package constants

import (
	"os"
	"strconv"
	"time"
)

// TimeoutMultiplier scales the saga poll timeouts for slow environments (CI). Read
// from SAGA_TIMEOUT_MULTIPLIER (a float), clamped to a 1.0 minimum so it can only
// EXTEND a wait, never shorten it; default 1.0 leaves the tuned timeouts untouched.
var TimeoutMultiplier = getTimeoutMultiplier()

func getTimeoutMultiplier() float64 {
	if v := os.Getenv("SAGA_TIMEOUT_MULTIPLIER"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 1.0 {
			return f
		}
	}
	return 1.0
}

// ScaleTimeout multiplies a base poll timeout by TimeoutMultiplier so a single env
// var stretches every wait for a slow stack without editing each assert.
func ScaleTimeout(d time.Duration) time.Duration {
	return time.Duration(float64(d) * TimeoutMultiplier)
}

// Service URLs
var (
	MapexosURL  = getEnv("MAPEXOS_URL", "http://localhost:5000")
	RouterURL   = getEnv("ROUTER_URL", "http://localhost:5003")
	AssetsURL   = getEnv("ASSETS_URL", "http://localhost:5002")
	GatewayURL  = getEnv("GATEWAY_URL", "http://localhost:5001")
	EventsURL   = getEnv("EVENTS_URL", "http://localhost:5004")
	TriggersURL = getEnv("TRIGGERS_URL", "http://localhost:5006")
	WorkflowURL = getEnv("WORKFLOW_URL", "http://localhost:5007")
	VaultURL    = getEnv("VAULT_URL", "http://localhost:5010")
)

// Internal API key for /internal/* routes (asset L3 fallback, asset-auth
// L3 fallback, healthmonitor force-offline). Default mirrors the value
// the standalone compose ships in services/envs/global.env so saga
// runs against `make standalone` work out of the box.
var (
	InternalApiKey = getEnv("INTERNAL_API_KEY", "5230c2e2-e245-468d-89e8-94154cf520d0")
)

// SinkHost is the host a trigger/broker config advertises so the service
// under test can reach the in-process server the saga starts. Default
// "localhost" for a fully local run (go run — everything on the host). Set
// SAGA_SINK_HOST=host.docker.internal when the services run in Docker and
// must reach the host-side server over the bridge. Only the host is
// configurable; the port is the OS-assigned ephemeral port the server
// publishes on the bag.
var (
	SinkHost = getEnv("SAGA_SINK_HOST", "127.0.0.1")
)

// MQTT broker URLs. Plaintext on :1883 (password mode), mTLS on :8883
// (cert mode). Saga connectivity / telemetry phases hit these as a
// real device would.
var (
	MqttBrokerURL    = getEnv("MQTT_BROKER_URL", "tcp://localhost:1883")
	MqttBrokerTLSURL = getEnv("MQTT_BROKER_TLS_URL", "ssl://localhost:8883")
)

// MongoDB
var (
	MongoURI      = getEnv("MONGO_URI", "mongodb://localhost:27017")
	MongoDatabase = getEnv("MONGO_DATABASE", "mapexos_test")
)

// Redis
var (
	RedisHost = getEnv("REDIS_HOST", "localhost")
	RedisPort = getEnv("REDIS_PORT", "6379")
	RedisDB   = 1 // Use DB 1 for tests
)

// ClickHouse (events raw-event store). The infra readiness probe TCP-dials the HTTP
// plane (:8123), host-reachable; the events service uses the native plane internally.
var (
	ClickHouseURL = getEnv("CLICKHOUSE_URL", "tcp://localhost:8123")
)

// mapexLNS (LoRaWAN Network Server) ingress endpoints the LoRaWAN e2e drives real
// gateway traffic at. UDP is the Semtech packet-forwarder ingress; Basics Station
// is the WebSocket ingress. Always env-overridable for flexibility; localhost by
// default so a local stack works out of the box.
var (
	LNSUDPHost     = getEnv("LNS_UDP_HOST", "127.0.0.1")
	LNSUDPPort     = getEnvInt("LNS_UDP_PORT", 1700)
	LNSBStationURI = getEnv("LNS_BSTATION_URI", "ws://127.0.0.1:1887")
	// LNSURL is the LNS HTTP plane (only /health; the gateway + uplink planes are
	// UDP/WebSocket). The infra readiness probe GETs LNSURL+"/health".
	LNSURL = getEnv("LNS_URL", "http://localhost:5011")
)

// js-executor (V8 script execution for IoT events, :8000) and js-workflow-executor
// (V8 execution of workflow code nodes, :8001) HTTP planes — only /health; both are
// otherwise NATS consumers. The infra readiness probe GETs them so the event / workflow
// pipelines have a running executor before a journey fires.
var (
	JsExecutorURL         = getEnv("JS_EXECUTOR_URL", "http://localhost:8000")
	JsWorkflowExecutorURL = getEnv("JS_WORKFLOW_EXECUTOR_URL", "http://localhost:8001")
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if n, err := strconv.Atoi(value); err == nil {
			return n
		}
	}
	return defaultValue
}

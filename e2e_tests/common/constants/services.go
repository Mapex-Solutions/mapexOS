package constants

import "os"

// Service URLs
var (
	MapexosURL  = getEnv("MAPEXOS_URL", "http://localhost:5000")
	RouterURL   = getEnv("ROUTER_URL", "http://localhost:5003")
	AssetsURL   = getEnv("ASSETS_URL", "http://localhost:5002")
	GatewayURL  = getEnv("GATEWAY_URL", "http://localhost:5001")
	EventsURL   = getEnv("EVENTS_URL", "http://localhost:5004")
	TriggersURL = getEnv("TRIGGERS_URL", "http://localhost:5006")
	WorkflowURL = getEnv("WORKFLOW_URL", "http://localhost:5007")
)

// Internal API key for /internal/* routes (asset L3 fallback, asset-auth
// L3 fallback, healthmonitor force-offline). Default mirrors the value
// the standalone compose ships in services/envs/global.env so saga
// runs against `make standalone` work out of the box.
var (
	InternalApiKey = getEnv("INTERNAL_API_KEY", "5230c2e2-e245-468d-89e8-94154cf520d0")
)

// Trigger HTTP sink URL — the connectivity_actions_*/phase2_trigger
// journey starts a local HTTP server bound to TriggerSinkBindAddr and
// configures the trigger to POST against this URL. Because the
// triggers service typically runs inside Docker, the URL the trigger
// sees must resolve to the host network (host.docker.internal on
// Linux/Mac compose with extra_hosts: host-gateway). Override per
// environment if needed.
var (
	// TriggerSinkURL is the URL written into the trigger HTTP config
	// — what the triggers service POSTs against. Defaults to
	// localhost because the Mapex services run on the host in dev;
	// override to host.docker.internal:11010 when running the
	// triggers service inside Docker.
	TriggerSinkURL = getEnv("SAGA_TRIGGER_SINK_URL", "http://localhost:11010")

	// TriggerSinkBindAddr is the host:port the sink listens on inside
	// the test process. Bind to 0.0.0.0 so Docker bridge / WSL2 can
	// route to it when needed; localhost binding would silently drop
	// external connections.
	TriggerSinkBindAddr = getEnv("SAGA_TRIGGER_SINK_BIND_ADDR", "0.0.0.0:11010")
)

// MQTT broker URLs. Plaintext on :1883 (password mode), mTLS on :8883
// (cert mode). Saga connectivity / telemetry phases hit these as a
// real device would.
var (
	MqttBrokerURL    = getEnv("MQTT_BROKER_URL", "tcp://localhost:1883")
	MqttBrokerTLSURL = getEnv("MQTT_BROKER_TLS_URL", "ssl://localhost:8883")
)

// WebSocket sink — used by the WebSocket trigger smoke journey.
var (
	// WsSinkBindAddr is the host:port the in-process WS sink listens
	// on. /ws is the upgrade endpoint, / is a health probe.
	WsSinkBindAddr = getEnv("SAGA_TRIGGER_WS_BIND_ADDR", "0.0.0.0:11026")

	// WsSinkURL is what gets written into the trigger config's
	// websocket.url field. ws:// because the sink has no TLS.
	WsSinkURL = getEnv("SAGA_TRIGGER_WS_URL", "ws://localhost:11026/ws")
)

// SMTP sink — used by the Email trigger smoke journey. The journey
// starts a local SMTP server bound to SmtpSinkBindAddr and configures
// the trigger with SmtpSinkHost / SmtpSinkPort so the triggers
// service delivers mail to it. host-bound just like the HTTP sink, so
// the same Docker-host caveat applies (override SAGA_TRIGGER_SMTP_HOST
// to host.docker.internal when the triggers service runs in a
// container).
var (
	// SmtpSinkBindAddr is the host:port the in-process SMTP server
	// listens on. 0.0.0.0 so Docker bridge / WSL2 can reach it.
	SmtpSinkBindAddr = getEnv("SAGA_TRIGGER_SMTP_BIND_ADDR", "0.0.0.0:11025")

	// SmtpSinkHost is what gets written into the trigger config's
	// smtpHost field — must resolve from the triggers service POV.
	SmtpSinkHost = getEnv("SAGA_TRIGGER_SMTP_HOST", "localhost")

	// SmtpSinkPort is what gets written into the trigger config's
	// smtpPort field. Kept separate from BindAddr so we can override
	// host alone (Docker scenarios) without retyping the port.
	SmtpSinkPort = getEnv("SAGA_TRIGGER_SMTP_PORT", "11025")
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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

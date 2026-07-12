package infra

import (
	"net"
	"net/http"
	"strings"
	"time"

	constants "github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
)

// probeTimeout bounds a single readiness probe so a down service fails fast
// instead of hanging the run.
const probeTimeout = 2 * time.Second

// serviceMap is the registry of provisionable services: logical name → compose
// service name (as it appears in that service's compose file), readiness probe, and
// how to run it from source (local mode). HTTP services probe GET /health; the
// non-HTTP ones (broker, nats) dial their TCP port; minio prefers its health endpoint.
// The compose NAMES must match the deploy compose exactly (mapexos→mapex-iam,
// http_gateway→http-gateway). local is nil for infra services (docker-only).
var serviceMap = map[string]serviceEntry{
	// mapex services stack (services/docker-compose.yml) — go run / npm run in local mode
	"mapexos":              {compose: "mapex-iam", ready: httpHealthy(constants.MapexosURL), local: goRun("mapexIam")},
	"http_gateway":         {compose: "http-gateway", ready: httpHealthy(constants.GatewayURL), local: goRun("http_gateway")},
	"assets":               {compose: "assets", ready: httpHealthy(constants.AssetsURL), local: goRun("assets")},
	"router":               {compose: "router", ready: httpHealthy(constants.RouterURL), local: goRun("router")},
	"events":               {compose: "events", ready: httpHealthy(constants.EventsURL), local: goRun("events")},
	"triggers":             {compose: "triggers", ready: httpHealthy(constants.TriggersURL), local: goRun("triggers")},
	"workflow":             {compose: "workflow", ready: httpHealthy(constants.WorkflowURL), local: goRun("workflow")},
	"vault":                {compose: "mapex-vault", ready: httpHealthy(constants.VaultURL), local: goRun("mapexVault")},
	"lns":                  {compose: "lns", ready: httpHealthy(constants.LNSURL), local: lnsRun()},
	"js_executor":          {compose: "js-executor", ready: httpHealthy(constants.JsExecutorURL), local: npmDev("js-executor")},
	"js_workflow_executor": {compose: "js-workflow-executor", ready: httpHealthy(constants.JsWorkflowExecutorURL), local: npmDev("js-workflow-executor")},

	// infra stack (infra/docker-compose.yml) — docker-only (no local run)
	"mongodb":    {compose: "mongodb", ready: tcpOpen(constants.MongoURI)},
	"redis":      {compose: "redis", ready: tcpOpen("tcp://" + constants.RedisHost + ":" + constants.RedisPort)},
	"clickhouse": {compose: "clickhouse", ready: tcpOpen(constants.ClickHouseURL)},
	"nats":       {compose: "nats-core", ready: tcpOpen("tcp://localhost:4222")},
	"broker":     {compose: "mapex-broker-mqtt", ready: tcpOpen(constants.MqttBrokerURL)},
	"minio":      {compose: "minio", ready: httpHealthyOrTCP("http://localhost:9000/minio/health/live", "localhost:9000")},
}

// The two provisioning stacks. Each is brought up (or reused) independently, with its
// own compose file and its own mode (docker | local) — see ensure_all.go.
var (
	infraServices = []string{"mongodb", "redis", "clickhouse", "minio", "nats", "broker"}
	mapexServices = []string{
		"mapexos", "http_gateway", "assets", "router", "events",
		"triggers", "workflow", "vault", "lns", "js_executor", "js_workflow_executor",
	}
)

// probe reports whether the named service is ready and whether it is known. It is
// a package var so tests replace it with a fake — no real network in unit tests.
var probe = func(service string) (ready, known bool) {
	e, ok := serviceMap[service]
	if !ok {
		return false, false
	}
	return e.ready(), true
}

// httpOK reports whether a GET of url returns 200 within probeTimeout.
func httpOK(url string) bool {
	client := &http.Client{Timeout: probeTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// httpHealthy returns a probe that GETs base+"/health".
func httpHealthy(base string) func() bool {
	return func() bool { return httpOK(base + "/health") }
}

// tcpOpen returns a probe that dials the host:port parsed from url.
func tcpOpen(url string) func() bool {
	return func() bool {
		conn, err := net.DialTimeout("tcp", hostPort(url), probeTimeout)
		if err != nil {
			return false
		}
		_ = conn.Close()
		return true
	}
}

// httpHealthyOrTCP returns a probe that tries the HTTP health endpoint, falling
// back to a TCP dial (used for minio, which exposes both).
func httpHealthyOrTCP(healthURL, hostport string) func() bool {
	tcp := tcpOpen("tcp://" + hostport)
	return func() bool { return httpOK(healthURL) || tcp() }
}

// hostPort strips a scheme:// prefix, returning host:port for net.Dial.
func hostPort(url string) string {
	if _, after, found := strings.Cut(url, "://"); found {
		return after
	}
	return url
}

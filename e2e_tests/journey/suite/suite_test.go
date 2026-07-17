//go:build saga

// Package suite is the single execution entry point for every saga journey. It
// replaces the per-journey journey_test.go / TestJourney wrappers (which would
// double-run as parallel packages and reintroduce port collisions): journeys are
// listed once in the registry and driven by TestSuite.
//
// Lifecycle (§10, §11):
//
//	TestMain: infra.EnsureAll()  → stack UP once (reuse what runs via CLI)
//	          m.Run()            → TestSuite
//	          teardown()         → stack DOWN (only what EnsureAll started)
//
// TestSuite runs each registered journey as a parallel subtest; every entity is
// runID-stamped so parallel journeys never collide (single-tenant by design, §13).
//
// Selection: go test -tags=saga -run 'TestSuite/iot/lorawan_journey_otaa/phase2_sensor_uplink'
package suite

import (
	"os"
	"testing"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/infra"
	leakcheck "github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/leakcheck"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	iam "github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/iam_bootstrap"

	// automations — each trigger's connectivity (phase1) + event pipeline (phase2)
	emailConn "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_email/phase1_connectivity"
	emailEvents "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_email/phase2_event_pipeline"
	httpConn "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_http/phase1_connectivity"
	httpEvents "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_http/phase2_event_pipeline"
	mqttConn "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_mqtt/phase1_connectivity"
	mqttEvents "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_mqtt/phase2_event_pipeline"
	natsConn "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_nats/phase1_connectivity"
	natsEvents "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_nats/phase2_event_pipeline"
	rabbitConn "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_rabbitmq/phase1_connectivity"
	rabbitEvents "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_rabbitmq/phase2_event_pipeline"
	slackConn "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_slack/phase1_connectivity"
	slackEvents "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_slack/phase2_event_pipeline"
	teamsConn "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_teams/phase1_connectivity"
	teamsEvents "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_teams/phase2_event_pipeline"
	wsConn "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_websocket/phase1_connectivity"
	wsEvents "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/automations/trigger_websocket/phase2_event_pipeline"

	// iot — connectivity actions (workflow phase1 → trigger phase2)
	connHTTPWorkflow "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/connectivity_actions_http/phase1_workflow"
	connHTTPTrigger "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/connectivity_actions_http/phase2_trigger"
	connMQTTWorkflow "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/connectivity_actions_mqtt/phase1_workflow"
	connMQTTTrigger "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/connectivity_actions_mqtt/phase2_trigger"

	// iot — gateway provisioning
	gwCert "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/gateway_cert_provisioning"
	gwKey "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/gateway_key_provisioning"
	gwProv "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/gateway_provisioning"

	// iot — mqtt broker auth
	brokerPwd "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/mqtt_broker_auth/phase1_password_user"
	brokerCert "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/mqtt_broker_auth/phase2_cert_user"

	// iot — lorawan (phased sensor journeys)
	bspUplink "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/lorawan_journey_bsp/phase1_sensor_uplink"
	bspPresence "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/lorawan_journey_bsp/phase2_sensor_presence"
	otaaConnectivity "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/lorawan_journey_otaa/phase1_gateway_connectivity"
	otaaUplink "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/lorawan_journey_otaa/phase2_sensor_uplink"
	otaaPresence "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/lorawan_journey_otaa/phase3_sensor_presence"

	// iot — lorawan device codec (uplink bytes → decoded semantic fields)
	lorawanCodec "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/lorawan_journey_codec"

	// iot — ota firmware update (http poll + mqtt push)
	otaHTTP "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/ota_http"
	otaMQTT "github.com/Mapex-Solutions/MapexOS/e2eTests/journey/iot/ota_mqtt"
)

// TestMain is the ONLY entry point: it owns bringing the stack up once and tearing
// down only what it started, wrapping the whole suite run.
func TestMain(m *testing.M) {
	// After each journey's rollback, sweep the public API for entities that still
	// carry that journey's RunID (a Compensate gap). Off unless wired here.
	saga.AfterRun = leakcheck.Hook

	// Start assets with a short OTA reconciler scan so the OTA journeys' dispatch
	// and early-close land within the device-wait and assert budgets; the default
	// (60s) is too slow — the MQTT device only waits 45s for the pushed command.
	// Suite-wide by design (one shared stack); only the OTA journeys read it.
	teardown := infra.EnsureAll(infra.WithServiceEnv("assets", "OTA_SCAN_INTERVAL=15"))
	// Defer teardown inside a func so it still runs if a journey panics out of m.Run
	// (a plain teardown() after m.Run would be skipped on panic). SIGKILL still leaks.
	code := func() int { defer teardown(); return m.Run() }()
	os.Exit(code)
}

// journey is one registered flow: a display name and its Run function.
type journey struct {
	name string
	run  func(t *testing.T)
}

// registry lists every journey the suite runs. This is REGISTRATION (start) order
// only — with t.Parallel there is NO completion ordering, so every journey must be
// fully self-contained and never rely on another finishing first (phases included;
// each phase provisions its own inputs). The grouping below is for readers, not
// sequencing. Adding a journey is one aliased import + one line here.
var registry = []journey{
	{"iam_bootstrap", iam.Run},

	{"automations/trigger_http/connectivity", httpConn.Run},
	{"automations/trigger_http/event_pipeline", httpEvents.Run},
	{"automations/trigger_mqtt/connectivity", mqttConn.Run},
	{"automations/trigger_mqtt/event_pipeline", mqttEvents.Run},
	{"automations/trigger_nats/connectivity", natsConn.Run},
	{"automations/trigger_nats/event_pipeline", natsEvents.Run},
	{"automations/trigger_rabbitmq/connectivity", rabbitConn.Run},
	{"automations/trigger_rabbitmq/event_pipeline", rabbitEvents.Run},
	{"automations/trigger_slack/connectivity", slackConn.Run},
	{"automations/trigger_slack/event_pipeline", slackEvents.Run},
	{"automations/trigger_teams/connectivity", teamsConn.Run},
	{"automations/trigger_teams/event_pipeline", teamsEvents.Run},
	{"automations/trigger_websocket/connectivity", wsConn.Run},
	{"automations/trigger_websocket/event_pipeline", wsEvents.Run},
	{"automations/trigger_email/connectivity", emailConn.Run},
	{"automations/trigger_email/event_pipeline", emailEvents.Run},

	{"iot/connectivity_actions_http/workflow", connHTTPWorkflow.Run},
	{"iot/connectivity_actions_http/trigger", connHTTPTrigger.Run},
	{"iot/connectivity_actions_mqtt/workflow", connMQTTWorkflow.Run},
	{"iot/connectivity_actions_mqtt/trigger", connMQTTTrigger.Run},

	{"iot/gateway_provisioning", gwProv.Run},
	{"iot/gateway_key_provisioning", gwKey.Run},
	{"iot/gateway_cert_provisioning", gwCert.Run},

	{"iot/mqtt_broker_auth/password", brokerPwd.Run},
	{"iot/mqtt_broker_auth/cert", brokerCert.Run},

	{"iot/lorawan_journey_bsp/phase1_sensor_uplink", bspUplink.Run},
	{"iot/lorawan_journey_bsp/phase2_sensor_presence", bspPresence.Run},
	{"iot/lorawan_journey_otaa/phase1_gateway_connectivity", otaaConnectivity.Run},
	{"iot/lorawan_journey_otaa/phase2_sensor_uplink", otaaUplink.Run},
	{"iot/lorawan_journey_otaa/phase3_sensor_presence", otaaPresence.Run},

	{"iot/lorawan_journey_codec", lorawanCodec.Run},

	{"iot/ota_http", otaHTTP.Run},
	{"iot/ota_mqtt", otaMQTT.Run},
}

// TestSuite runs every registered journey as a parallel subtest.
func TestSuite(t *testing.T) {
	for _, j := range registry {
		j := j
		t.Run(j.name, func(t *testing.T) {
			t.Parallel()
			j.run(t)
		})
	}
}

package message

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// NATS streams, subjects, and message identifiers used by the healthmonitor
// module. These values are intra-service: producers and consumers all live
// inside the assets service and share this single source of truth. Stream,
// subject, and durable names resolve at package init from GO_ENV so the
// same binary serves multiple environments on a shared NATS cluster.
//
// External producers (js-executor, http_gateway, MQTT republish) publish to
// HeartbeatSubject — the subject is part of the cross-service contract and
// duplicated here only because the consumer needs to subscribe to it.

// HeartbeatStreamName is the JetStream stream that captures all
// heartbeat publishes (any producer, any origin). Resolved at package
// init — e.g. "DEV-MAPEXOS-ASSETS-HEARTBEAT".
var HeartbeatStreamName = config.StreamName("ASSETS", "HEARTBEAT")

// HeartbeatSubject is the wildcard subject the heartbeat consumer
// subscribes to. The trailing token is the orgId. Resolved at package
// init — e.g. "dev.mapexos.asset.heartbeat.>".
var HeartbeatSubject = config.Subject("asset", "heartbeat") + ".>"

// HeartbeatEventType is the event-type label used by the DLQ policy.
const HeartbeatEventType = "asset.heartbeat"

// ScanStreamName is the JetStream stream that drives periodic scan
// scheduling. WorkQueue retention with AllowMsgSchedules. Resolved at
// package init — e.g. "DEV-MAPEXOS-ASSETS-HEALTH-MONITOR".
var ScanStreamName = config.StreamName("ASSETS", "HEALTH-MONITOR")

// ScanSubject is the subject the scan consumer subscribes to. The
// service publishes the next scan to ScanScheduleSubject, JetStream
// then redirects it to this subject when the schedule fires. Resolved
// at package init — e.g. "dev.mapexos.healthmonitor.scan".
var ScanSubject = config.Subject("healthmonitor", "scan")

// ScanScheduleSubject is the subject used by the service to publish
// the next scheduled scan with a TargetSubject of ScanSubject. Resolved
// at package init — e.g. "dev.mapexos.healthmonitor.scan.schedule".
var ScanScheduleSubject = config.Subject("healthmonitor", "scan.schedule")

// ScanMsgId is a fixed message id used together with stream
// Duplicates to ensure exactly-one pending schedule across pods.
const ScanMsgId = "hm-scan"

// MqttPresenceStreamName is the JetStream stream that captures broker
// presence advisories published by the Mosquitto broker plugin
// (mapex-broker-plugin). Single subject — both connect and disconnect
// share it, discriminated by the Event field on the payload. The stream
// feeds two separate consumers (one filter per Event) inside the
// healthmonitor module — see interfaces/message/consumers/presence/.
// Cluster-ready: replicas parameterized via STREAM_REPLICAS at stream
// creation time. Resolved at package init —
// e.g. "DEV-MAPEXOS-ASSETS-MQTT-PRESENCE".
var MqttPresenceStreamName = config.StreamName("ASSETS", "MQTT-PRESENCE")

// MqttPresenceAdvisorySubject is the single subject the broker plugin
// publishes to on every device CONNECT and DISCONNECT. Payload schema:
// healthmonitor.PresenceAdvisory (Event field = "connect"|"disconnect").
// The two healthmonitor consumers subscribe to the same subject and
// gate by Event in their handlers.
var MqttPresenceAdvisorySubject = config.Subject("mqtt", "presence.advisory")

// MqttPresenceEventType is the DLQ event-type label for the presence
// consumer family.
const MqttPresenceEventType = "mqtt-presence"

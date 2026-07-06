package message

import (
	otaEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/events"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// OTA timer scheduling subjects. The subjects are STATIC — the plan id travels
// in the message payload (and the MsgId), NEVER in the subject.
var (
	// SubjectOTASchedule is the single static subject where scheduled OTA timer
	// messages wait until their @at time.
	SubjectOTASchedule = config.Subject("ota", "schedule")

	// Target subjects the scheduled messages are delivered to when they fire.
	SubjectOTATimerStart = config.Subject("ota", "timer.start")
	SubjectOTATimerClose = config.Subject("ota", "timer.close")
	SubjectOTATimerScan  = config.Subject("ota", "timer.scan")

	// OTAScheduleStream carries the scheduled + delivered timer messages AND the
	// inbound device status advisories (AllowMsgSchedules, file storage — created
	// in bootstrap). Subjects: ota.> (schedule, timer.*, status.advisory).
	OTAScheduleStream = config.StreamName("OTA", "SCHEDULE")

	// SubjectOTAStatusAdvisory re-exports the cross-service inbound advisory
	// subject (published by the HTTP gateway / mapexMQTTBroker; consumed here).
	SubjectOTAStatusAdvisory = otaEvents.SubjectOTAStatusAdvisory
)

// Downlink subjects are cross-service (consumed by mapexMQTTBroker / mapexLNS)
// and live in the shared contract: contracts/services/assets/downlink.
var (
	// SubjectOTATimerAbandon — the firmware abandon-check fires here (routed by
	// the timers consumer to the firmware service).
	SubjectOTATimerAbandon = config.Subject("ota", "timer.abandon")
)

// ScanMsgId is the FIXED message id for the single global pacing scan (only one
// scan is pending network-wide, like the healthmonitor scan).
const ScanMsgId = "ota-scan"

// AbandonMsgId dedups a firmware's abandon-check timer (firmwareId in the
// MsgId + payload, never the subject).
func AbandonMsgId(firmwareID string) string { return "ota-abandon-" + firmwareID }

// StartMsgId / CloseMsgId dedup a plan's start/close timer. The plan id lives in
// the MsgId metadata + payload — NOT the subject.
func StartMsgId(planID string) string { return "ota-start-" + planID }
func CloseMsgId(planID string) string { return "ota-close-" + planID }

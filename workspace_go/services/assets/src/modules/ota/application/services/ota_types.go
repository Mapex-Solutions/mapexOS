package services

import (
	"time"

	"assets/src/modules/ota/application/ports"
	"assets/src/modules/ota/domain/repositories"
)

// checksumAlgorithmSHA256 is the algorithm recorded on firmware artifacts.
const checksumAlgorithmSHA256 = "SHA-256"

// protocolHTTP devices poll for their job rather than being pushed.
const protocolHTTP = "http"

// mqttStatusEventType is the event type an MQTT device reports OTA progress
// under, following the broker's device topic contract
// events/{assetUUID}/{eventType} (the assetUUID is pinned to the authenticated
// session by the broker ACL).
const mqttStatusEventType = "ota_status"

// mqttStatusTopicFor builds the device's own OTA status topic.
func mqttStatusTopicFor(assetUUID string) string {
	return "events/" + assetUUID + "/" + mqttStatusEventType
}

// httpStatusReportPath is the gateway endpoint an HTTP device reports its OTA
// status to (authenticated per data source).
const httpStatusReportPath = "/api/v1/ota/status"

// Close reasons passed to ClosePlan.
const (
	CloseReasonMaxTime     = "maxTime"
	CloseReasonAllTerminal = "all_terminal"
	CloseReasonCancel      = "cancel"
)

// FirmwareServiceDeps aggregates the firmware service dependencies.
type FirmwareServiceDeps struct {
	FirmwareRepo repositories.FirmwareRepository
	Store        ports.FirmwareStorePort
	Scheduler    ports.FirmwareSchedulerPort
	PresignTTL   time.Duration
	AbandonTTL   time.Duration
}

// FirmwareService implements OTAFirmwareServicePort.
type FirmwareService struct {
	deps FirmwareServiceDeps
}

// PlanServiceDeps aggregates the plan service dependencies.
type PlanServiceDeps struct {
	PlanRepo      repositories.OTAPlanRepository
	ExecutionRepo repositories.OTAExecutionRepository
	FirmwareRepo  repositories.FirmwareRepository
	Scheduler     ports.OTASchedulerPort
	Store         ports.FirmwareStorePort
	PresignTTL    time.Duration
}

// PlanService implements OTAPlanServicePort.
type PlanService struct {
	deps PlanServiceDeps
}

// ReconcilerDeps aggregates the reconciler dependencies.
type ReconcilerDeps struct {
	PlanRepo      repositories.OTAPlanRepository
	ExecutionRepo repositories.OTAExecutionRepository
	FirmwareRepo  repositories.FirmwareRepository
	Store         ports.FirmwareStorePort
	Presence      ports.PresenceReaderPort
	Edge          ports.EdgeDispatchPort
	AssetReader   ports.AssetReaderPort
	PresignTTL    time.Duration
	MaxAttempts   int
	DefaultRate   int
}

// Reconciler drives a plan's executions toward the target template.
type Reconciler struct {
	deps ReconcilerDeps
}

// ReconcilerTimersDeps aggregates the timer-service dependencies.
type ReconcilerTimersDeps struct {
	PlanRepo      repositories.OTAPlanRepository
	ExecutionRepo repositories.OTAExecutionRepository
	FirmwareRepo  repositories.FirmwareRepository
	Store         ports.FirmwareStorePort
	Scheduler     ports.OTASchedulerPort
	Reconciler    *Reconciler
	ScanInterval  time.Duration
}

// ReconcilerTimers drives the plan timers and the global pacing scan.
type ReconcilerTimers struct {
	deps ReconcilerTimersDeps
}

// StatusHandlerDeps aggregates the status-handler dependencies.
type StatusHandlerDeps struct {
	ExecutionRepo    repositories.OTAExecutionRepository
	PlanRepo         repositories.OTAPlanRepository
	LiveState        ports.LiveStatePort
	History          ports.StatusHistoryPublisherPort
	TemplateSwitcher ports.AssetTemplateSwitcherPort
}

// StatusHandler consumes a device status advisory and updates the execution,
// counters, live state, history, and (on UPDATED) the asset template.
type StatusHandler struct {
	deps StatusHandlerDeps
}

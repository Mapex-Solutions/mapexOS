package di

import (
	"go.uber.org/dig"

	"assets/src/modules/ota/application/ports"
	"assets/src/modules/ota/domain/repositories"
)

// FirmwareServiceDI aggregates the firmware service dependencies. Config-derived
// scalars (TTLs) are resolved in the module provider, not injected.
type FirmwareServiceDI struct {
	dig.In
	FirmwareRepo repositories.FirmwareRepository
	Store        ports.FirmwareStorePort
	Scheduler    ports.FirmwareSchedulerPort
}

// PlanServiceDI aggregates the plan service dependencies.
type PlanServiceDI struct {
	dig.In
	PlanRepo      repositories.OTAPlanRepository
	ExecutionRepo repositories.OTAExecutionRepository
	FirmwareRepo  repositories.FirmwareRepository
	Scheduler     ports.OTASchedulerPort
	Store         ports.FirmwareStorePort
}

// ReconcilerDI aggregates the reconciler-engine dependencies. Pacing/retry
// scalars are resolved from config in the module provider.
type ReconcilerDI struct {
	dig.In
	PlanRepo      repositories.OTAPlanRepository
	ExecutionRepo repositories.OTAExecutionRepository
	FirmwareRepo  repositories.FirmwareRepository
	Store         ports.FirmwareStorePort
	Presence      ports.PresenceReaderPort
	Edge          ports.EdgeDispatchPort
	AssetReader   ports.AssetReaderPort
}

// ReconcilerTimersDI aggregates the timer-service dependencies. The concrete
// *Reconciler is injected as a separate provider argument to keep this struct
// free of the services package.
type ReconcilerTimersDI struct {
	dig.In
	PlanRepo      repositories.OTAPlanRepository
	ExecutionRepo repositories.OTAExecutionRepository
	FirmwareRepo  repositories.FirmwareRepository
	Store         ports.FirmwareStorePort
	Scheduler     ports.OTASchedulerPort
}

// StatusHandlerDI aggregates the status-handler dependencies.
type StatusHandlerDI struct {
	dig.In
	ExecutionRepo    repositories.OTAExecutionRepository
	PlanRepo         repositories.OTAPlanRepository
	LiveState        ports.LiveStatePort
	History          ports.StatusHistoryPublisherPort
	TemplateSwitcher ports.AssetTemplateSwitcherPort
}

package services

import (
	"context"
	"strings"
	"time"

	"assets/src/modules/ota/domain/entities"

	otaEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/events"
)

// NewStatusHandler returns a status handler over the given dependencies.
func NewStatusHandler(deps StatusHandlerDeps) *StatusHandler {
	return &StatusHandler{deps: deps}
}

// HandleStatus processes one device status advisory. It is a no-op for unknown,
// already-terminal, or non-owned executions (stale / duplicate / forged reports).
func (h *StatusHandler) HandleStatus(c context.Context, adv otaEvents.OTAStatusAdvisory) error {
	if adv.OTAExecutionID == "" {
		return nil
	}
	exec, err := h.deps.ExecutionRepo.FindById(c, &adv.OTAExecutionID)
	if err != nil {
		return err
	}
	if exec == nil || exec.ID.IsZero() || exec.IsTerminal() {
		return nil
	}
	// Ownership: the reporting device must own this execution.
	if exec.AssetUUID != "" && adv.AssetUUID != "" && exec.AssetUUID != adv.AssetUUID {
		return nil
	}

	newState, ok := mapDeviceStatus(adv.Status)
	if !ok {
		return nil
	}

	fields := map[string]any{
		"state":      string(newState),
		"percentage": adv.Progress,
		"updated":    time.Now(),
	}
	if adv.Error != "" {
		fields["error"] = adv.Error
	}
	updated, err := h.deps.ExecutionRepo.FindByIdAndUpdate(c, &adv.OTAExecutionID, fields)
	if err != nil {
		return err
	}

	h.bumpCounters(c, exec.PlanID.Hex(), newState)
	if updated != nil {
		_ = h.deps.LiveState.SetExecutionLive(c, updated)
	}
	_ = h.deps.History.Publish(c, adv)

	if newState == entities.ExecUpdated {
		h.switchTemplate(c, exec)
	}
	return nil
}

func (h *StatusHandler) bumpCounters(c context.Context, planID string, state entities.ExecutionState) {
	switch state {
	case entities.ExecUpdated:
		_ = h.deps.PlanRepo.IncrementCounter(c, planID, "counters.succeeded", 1)
	case entities.ExecFailed:
		_ = h.deps.PlanRepo.IncrementCounter(c, planID, "counters.failed", 1)
	}
}

func (h *StatusHandler) switchTemplate(c context.Context, exec *entities.OTAExecution) {
	planID := exec.PlanID.Hex()
	plan, err := h.deps.PlanRepo.FindById(c, &planID)
	if err != nil || plan == nil || plan.ID.IsZero() {
		return
	}
	_ = h.deps.TemplateSwitcher.SwitchTemplate(c, exec.AssetID.Hex(), plan.TargetTemplateID.Hex())
}

// mapDeviceStatus maps a device-reported status string to an execution state.
func mapDeviceStatus(s string) (entities.ExecutionState, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "downloading":
		return entities.ExecDownloading, true
	case "downloaded":
		return entities.ExecDownloaded, true
	case "verified":
		return entities.ExecVerified, true
	case "updating":
		return entities.ExecUpdating, true
	case "updated":
		return entities.ExecUpdated, true
	case "failed":
		return entities.ExecFailed, true
	}
	return "", false
}

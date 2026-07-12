package payloads

import "fmt"

// SagaOtaSourceTemplate is the template an OTA asset starts on (the plan's
// sourceTemplateId). It reuses the canonical temperature template — a minimal
// valid template — with a distinct, runID-stamped name so it never collides with
// the target template created in the same journey.
func SagaOtaSourceTemplate(runID string) *AssetTemplateCreateBuilder {
	return SagaTemperatureTemplate(runID).WithName(fmt.Sprintf("saga-ota-source-%s", runID))
}

// SagaOtaTargetTemplate is the template the firmware migrates the asset TO (the
// firmware's targetTemplateId; the switch destination). Same minimal valid shape
// as the source, distinguished only by its runID-stamped name.
func SagaOtaTargetTemplate(runID string) *AssetTemplateCreateBuilder {
	return SagaTemperatureTemplate(runID).WithName(fmt.Sprintf("saga-ota-target-%s", runID))
}

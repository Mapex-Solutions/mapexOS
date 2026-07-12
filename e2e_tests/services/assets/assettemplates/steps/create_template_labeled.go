package steps

import (
	"fmt"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/payloads"
)

// CreateTemplateWithLabel is the label-scoped variant of CreateTemplateWith: it
// POSTs whatever template the injected builder produces and writes the returned
// id under TemplateIDKey(label) instead of the shared BagKeyTemplateID. A journey
// that provisions several templates (e.g. an OTA source + target) uses distinct
// labels so their ids and idempotent Compensate never collide.
//
// Reads (bag):
//   - none — auth headers live on the HTTP client; runID lives on Context.
//
// Writes (bag):
//   - TemplateIDKey(label)  string  Mongo ObjectID hex of the new template
//
// Compensate: DELETE /api/v1/asset_templates/{id}, reading the id back from
// TemplateIDKey(label) (404-tolerant, bag-driven) — shared with createTemplateStep.
func CreateTemplateWithLabel(build func(runID string) *payloads.AssetTemplateCreateBuilder, label string) saga.Step {
	name := fmt.Sprintf("assets/assettemplates.CreateTemplate[%s]", label)
	return createTemplateStep(name, TemplateIDKey(label), build)
}

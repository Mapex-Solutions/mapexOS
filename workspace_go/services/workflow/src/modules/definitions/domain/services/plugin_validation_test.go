package services

import (
	"testing"

	"workflow/src/modules/definitions/domain/entities"
)

func TestExtractRequiredPlugins_CoreOnly(t *testing.T) {
	nodes := []entities.WorkflowNode{
		{ID: "1", Type: "core/start"},
		{ID: "2", Type: "core/condition"},
		{ID: "3", Type: "core/end"},
	}

	result := ExtractRequiredPlugins(nodes)
	if len(result) != 0 {
		t.Errorf("expected 0 required plugins, got %d: %v", len(result), result)
	}
}

func TestExtractRequiredPlugins_WithMarketplace(t *testing.T) {
	nodes := []entities.WorkflowNode{
		{ID: "1", Type: "core/start"},
		{ID: "2", Type: "telegram/sendMessage"},
		{ID: "3", Type: "telegram/getUpdates"},
		{ID: "4", Type: "slack/postMessage"},
		{ID: "5", Type: "core/end"},
	}

	result := ExtractRequiredPlugins(nodes)
	if len(result) != 2 {
		t.Fatalf("expected 2 required plugins, got %d: %v", len(result), result)
	}
	if result[0] != "slack" || result[1] != "telegram" {
		t.Errorf("expected [slack, telegram], got %v", result)
	}
}

func TestExtractRequiredPlugins_Empty(t *testing.T) {
	result := ExtractRequiredPlugins(nil)
	if len(result) != 0 {
		t.Errorf("expected 0 required plugins for nil nodes, got %d", len(result))
	}
}

func TestExtractRequiredPlugins_NoSlash(t *testing.T) {
	nodes := []entities.WorkflowNode{
		{ID: "1", Type: "unknown"},
	}

	result := ExtractRequiredPlugins(nodes)
	if len(result) != 0 {
		t.Errorf("expected 0 required plugins for type without slash, got %d", len(result))
	}
}

func TestComputeDefinitionStatus_AllPresent(t *testing.T) {
	status, missing := ComputeDefinitionStatus(
		[]string{"slack", "telegram"},
		[]string{"openai", "slack", "telegram"},
	)

	if status != entities.StatusValid {
		t.Errorf("expected StatusValid, got %s", status)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing plugins, got %v", missing)
	}
}

func TestComputeDefinitionStatus_SomeMissing(t *testing.T) {
	status, missing := ComputeDefinitionStatus(
		[]string{"openai", "slack", "telegram"},
		[]string{"slack"},
	)

	if status != entities.StatusPluginMissing {
		t.Errorf("expected StatusPluginMissing, got %s", status)
	}
	if len(missing) != 2 {
		t.Fatalf("expected 2 missing plugins, got %d: %v", len(missing), missing)
	}
	if missing[0] != "openai" || missing[1] != "telegram" {
		t.Errorf("expected [openai, telegram], got %v", missing)
	}
}

func TestComputeDefinitionStatus_NoRequired(t *testing.T) {
	status, missing := ComputeDefinitionStatus(nil, []string{"slack"})

	if status != entities.StatusValid {
		t.Errorf("expected StatusValid for no required plugins, got %s", status)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing plugins, got %v", missing)
	}
}

func TestComputeDefinitionStatus_NoEnabled(t *testing.T) {
	status, missing := ComputeDefinitionStatus(
		[]string{"telegram"},
		nil,
	)

	if status != entities.StatusPluginMissing {
		t.Errorf("expected StatusPluginMissing, got %s", status)
	}
	if len(missing) != 1 || missing[0] != "telegram" {
		t.Errorf("expected [telegram], got %v", missing)
	}
}

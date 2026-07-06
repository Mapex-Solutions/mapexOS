package entities

import "testing"

func TestCanTransition_AllowedAndRejectedMoves(t *testing.T) {
	tests := []struct {
		name string
		from PlanStatus
		to   PlanStatus
		want bool
	}{
		// allowed
		{"pending to running", PlanPending, PlanRunning, true},
		{"scheduled to running", PlanScheduled, PlanRunning, true},
		{"running to complete", PlanRunning, PlanComplete, true},
		{"running to completed_with_errors", PlanRunning, PlanCompletedWithErrors, true},
		{"pending to cancelled", PlanPending, PlanCancelled, true},
		{"scheduled to cancelled", PlanScheduled, PlanCancelled, true},
		// rejected
		{"running to pending", PlanRunning, PlanPending, false},
		{"complete to running", PlanComplete, PlanRunning, false},
		{"cancelled to running", PlanCancelled, PlanRunning, false},
		{"complete to cancelled", PlanComplete, PlanCancelled, false},
		{"running to cancelled", PlanRunning, PlanCancelled, false},
		{"scheduled to complete", PlanScheduled, PlanComplete, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CanTransition(tt.from, tt.to); got != tt.want {
				t.Fatalf("CanTransition(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestMigrationPlan_IsEditable(t *testing.T) {
	tests := []struct {
		name   string
		status PlanStatus
		want   bool
	}{
		{"pending is editable", PlanPending, true},
		{"scheduled is editable", PlanScheduled, true},
		{"running is not editable", PlanRunning, false},
		{"complete is not editable", PlanComplete, false},
		{"completed_with_errors is not editable", PlanCompletedWithErrors, false},
		{"cancelled is not editable", PlanCancelled, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := &MigrationPlan{Status: tt.status}
			if got := plan.IsEditable(); got != tt.want {
				t.Fatalf("IsEditable() with status %q = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestMigrationPlan_IsTerminal(t *testing.T) {
	tests := []struct {
		name   string
		status PlanStatus
		want   bool
	}{
		{"pending is not terminal", PlanPending, false},
		{"scheduled is not terminal", PlanScheduled, false},
		{"running is not terminal", PlanRunning, false},
		{"complete is terminal", PlanComplete, true},
		{"completed_with_errors is terminal", PlanCompletedWithErrors, true},
		{"cancelled is terminal", PlanCancelled, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := &MigrationPlan{Status: tt.status}
			if got := plan.IsTerminal(); got != tt.want {
				t.Fatalf("IsTerminal() with status %q = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

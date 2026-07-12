package saga

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

// recordStep is a fake Step that appends "do:<name>" / "comp:<name>" to a shared log
// so a test can assert execution + rollback order. failDo makes Do return an error
// (after recording the attempt) to drive the item-failure path.
func recordStep(log *[]string, name string, failDo bool) Step {
	return Step{
		Name: name,
		Do: func(*Context) error {
			*log = append(*log, "do:"+name)
			if failDo {
				return errors.New("boom " + name)
			}
			return nil
		},
		Compensate: func(*Context) error {
			*log = append(*log, "comp:"+name)
			return nil
		},
	}
}

// TestRun_OrderAndRollback proves the guaranteed order — items in order, then every
// executed Step's Compensate in reverse — on success AND on an item failure (only the
// executed items roll back). It drives the internal run (which returns the failure
// instead of calling t.Errorf) so the failure case stays green.
func TestRun_OrderAndRollback(t *testing.T) {
	tests := []struct {
		name      string
		failStepB bool
		wantLog   []string
		wantFail  bool
	}{
		{"success", false, []string{"do:A", "do:B", "comp:B", "comp:A"}, false},
		{"item failure", true, []string{"do:A", "do:B", "comp:A"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var log []string
			a := recordStep(&log, "A", false)
			b := recordStep(&log, "B", tt.failStepB)

			sctx := newContext(t, context.Background(), "test", ClientSet{})
			name, err := run(sctx, a, b)

			if (err != nil) != tt.wantFail {
				t.Fatalf("err=%v (name=%q), wantFail=%v", err, name, tt.wantFail)
			}
			if !reflect.DeepEqual(log, tt.wantLog) {
				t.Fatalf("order = %v, want %v", log, tt.wantLog)
			}
		})
	}
}

// TestRun_RollsBackFullList proves the public Run rolls back the FULL executed list in
// reverse — guarding the closure-captured executed slice (a past bug captured it at
// len 0, making rollback a silent no-op).
func TestRun_RollsBackFullList(t *testing.T) {
	var log []string
	a := recordStep(&log, "A", false)
	b := recordStep(&log, "B", false)

	Run(t, context.Background(), "test", ClientSet{}, a, b)

	want := []string{"do:A", "do:B", "comp:B", "comp:A"}
	if !reflect.DeepEqual(log, want) {
		t.Fatalf("order = %v, want %v (rollback must run the full executed list in reverse)", log, want)
	}
}

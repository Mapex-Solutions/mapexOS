package saga

import "testing"

// TestNewTestContext_Usable proves the exported constructor yields a Context an
// external step/assert test can drive: the passed handle and a stable RunID, plus a
// working bag round-trip.
func TestNewTestContext_Usable(t *testing.T) {
	c := NewTestContext(t, ClientSet{})

	if c.T != t {
		t.Error("Context.T is not the passed *testing.T")
	}
	if c.RunID != "test" {
		t.Errorf("RunID = %q, want %q", c.RunID, "test")
	}

	c.Set("k", "v")
	got, ok := c.Get("k")
	if !ok || got != "v" {
		t.Errorf("bag round-trip: got %v ok=%v, want v true", got, ok)
	}
}

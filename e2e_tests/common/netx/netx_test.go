package netx

import "testing"

// TestFreeListener_ReturnsUsableFreePort proves a single call yields a listener on
// a real, non-zero port.
func TestFreeListener_ReturnsUsableFreePort(t *testing.T) {
	ln, port, err := FreeListener()
	if err != nil {
		t.Fatalf("FreeListener: %v", err)
	}
	defer ln.Close()

	if port <= 0 {
		t.Errorf("port = %d, want > 0", port)
	}
}

// TestFreeListener_DistinctPortsWhileHeld proves two listeners held open at the same
// time never share a port — the property that lets parallel journeys each grab their
// own sink port without colliding.
func TestFreeListener_DistinctPortsWhileHeld(t *testing.T) {
	ln1, port1, err := FreeListener()
	if err != nil {
		t.Fatalf("FreeListener #1: %v", err)
	}
	defer ln1.Close()

	ln2, port2, err := FreeListener()
	if err != nil {
		t.Fatalf("FreeListener #2: %v", err)
	}
	defer ln2.Close()

	if port1 == port2 {
		t.Errorf("two concurrently-held listeners share port %d, want distinct", port1)
	}
}

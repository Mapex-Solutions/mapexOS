package services

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestWaitForActiveWalkers_BlocksUntilDone(t *testing.T) {
	svc := &RuntimeService{}

	svc.activeWalks.Add(2)

	var waited atomic.Bool

	go func() {
		svc.WaitForActiveWalkers()
		waited.Store(true)
	}()

	// Should NOT have returned yet
	time.Sleep(50 * time.Millisecond)
	if waited.Load() {
		t.Fatal("WaitForActiveWalkers returned before walkers finished")
	}

	// Finish one walker
	svc.activeWalks.Done()
	time.Sleep(10 * time.Millisecond)
	if waited.Load() {
		t.Fatal("WaitForActiveWalkers returned with 1 walker still active")
	}

	// Finish second walker
	svc.activeWalks.Done()
	time.Sleep(10 * time.Millisecond)
	if !waited.Load() {
		t.Fatal("WaitForActiveWalkers did not return after all walkers finished")
	}
}

func TestWaitForActiveWalkers_ReturnsImmediatelyWhenNoWalkers(t *testing.T) {
	svc := &RuntimeService{}

	done := make(chan struct{})
	go func() {
		svc.WaitForActiveWalkers()
		close(done)
	}()

	select {
	case <-done:
		// returned immediately — correct
	case <-time.After(1 * time.Second):
		t.Fatal("WaitForActiveWalkers blocked with no active walkers")
	}
}

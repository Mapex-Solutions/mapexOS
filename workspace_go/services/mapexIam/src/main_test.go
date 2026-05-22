package main

import (
	"context"
	"errors"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/shutdown"
)

// sendSIGINT sends SIGINT to the current process after a delay.
func sendSIGINT(delay time.Duration) {
	go func() {
		time.Sleep(delay)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGINT)
	}()
}

// hookTracker records hook execution order and timing.
type hookTracker struct {
	mu        sync.Mutex
	order     []string
	startedAt map[string]time.Time
	endedAt   map[string]time.Time
}

func newHookTracker() *hookTracker {
	return &hookTracker{
		startedAt: make(map[string]time.Time),
		endedAt:   make(map[string]time.Time),
	}
}

func (h *hookTracker) record(name string, duration time.Duration) {
	h.mu.Lock()
	h.startedAt[name] = time.Now()
	h.mu.Unlock()

	time.Sleep(duration)

	h.mu.Lock()
	h.order = append(h.order, name)
	h.endedAt[name] = time.Now()
	h.mu.Unlock()
}

func (h *hookTracker) getOrder() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	cp := make([]string, len(h.order))
	copy(cp, h.order)
	return cp
}

// TestShutdownIntegration_MirrorsMainConfig verifies the exact hook
// configuration from main.go: Fiber at P0, MongoDB/Redis App/Redis Shared/NATS at P5.
// Fiber must complete before any P5 hook starts.
func TestShutdownIntegration_MirrorsMainConfig(t *testing.T) {
	sm := shutdown.New()
	tracker := newHookTracker()

	// P0: Fiber — simulate drain (20ms)
	sm.RegisterFunc("fiber", 0, func(_ context.Context) error {
		tracker.record("fiber", 20*time.Millisecond)
		return nil
	})

	// P5: MongoDB — simulate disconnect (15ms)
	sm.RegisterFunc("mongodb", 5, func(_ context.Context) error {
		tracker.record("mongodb", 15*time.Millisecond)
		return nil
	})

	// P5: Redis App — simulate close (10ms)
	sm.RegisterFunc("redis-app", 5, func(_ context.Context) error {
		tracker.record("redis-app", 10*time.Millisecond)
		return nil
	})

	// P5: Redis Shared — simulate close (10ms)
	sm.RegisterFunc("redis-shared", 5, func(_ context.Context) error {
		tracker.record("redis-shared", 10*time.Millisecond)
		return nil
	})

	// P5: NATS — simulate close (10ms)
	sm.RegisterFunc("nats", 5, func(_ context.Context) error {
		tracker.record("nats", 10*time.Millisecond)
		return nil
	})

	sendSIGINT(50 * time.Millisecond)
	sm.WaitForSignal(15 * time.Second)

	order := tracker.getOrder()

	// All 5 hooks must have executed
	if len(order) != 5 {
		t.Fatalf("expected 5 hooks executed, got %d: %v", len(order), order)
	}

	// Fiber (P0) must finish first
	if order[0] != "fiber" {
		t.Fatalf("expected fiber first, got %v", order)
	}

	// P5 hooks must all be present (order among them is non-deterministic)
	p5 := make(map[string]bool)
	for _, name := range order[1:] {
		p5[name] = true
	}
	for _, expected := range []string{"mongodb", "redis-app", "redis-shared", "nats"} {
		if !p5[expected] {
			t.Fatalf("missing %s in P5 group, got %v", expected, order[1:])
		}
	}

	// Fiber must end before any P5 hook starts
	tracker.mu.Lock()
	fiberEnd := tracker.endedAt["fiber"]
	for _, name := range []string{"mongodb", "redis-app", "redis-shared", "nats"} {
		if tracker.startedAt[name].Before(fiberEnd) {
			t.Fatalf("%s started before fiber finished", name)
		}
	}
	tracker.mu.Unlock()
}

// TestShutdownIntegration_P5ConnectionsConcurrent verifies that MongoDB,
// Redis App, Redis Shared and NATS close concurrently (same P5 priority).
func TestShutdownIntegration_P5ConnectionsConcurrent(t *testing.T) {
	sm := shutdown.New()

	var running atomic.Int32
	var maxConcurrent atomic.Int32

	makeP5Hook := func(name string) func(context.Context) error {
		return func(_ context.Context) error {
			cur := running.Add(1)
			for {
				old := maxConcurrent.Load()
				if cur <= old || maxConcurrent.CompareAndSwap(old, cur) {
					break
				}
			}
			time.Sleep(40 * time.Millisecond)
			running.Add(-1)
			return nil
		}
	}

	// P0: Fiber (fast)
	sm.RegisterFunc("fiber", 0, func(_ context.Context) error { return nil })

	// P5: All connections
	sm.RegisterFunc("mongodb", 5, makeP5Hook("mongodb"))
	sm.RegisterFunc("redis-app", 5, makeP5Hook("redis-app"))
	sm.RegisterFunc("redis-shared", 5, makeP5Hook("redis-shared"))
	sm.RegisterFunc("nats", 5, makeP5Hook("nats"))

	sendSIGINT(50 * time.Millisecond)
	sm.WaitForSignal(15 * time.Second)

	if maxConcurrent.Load() < 4 {
		t.Fatalf("expected all 4 P5 hooks running concurrently, max concurrent was %d", maxConcurrent.Load())
	}
}

// TestShutdownIntegration_MongoErrorDoesNotBlockOthers verifies that
// if MongoDB Close fails, Redis and NATS still close successfully.
func TestShutdownIntegration_MongoErrorDoesNotBlockOthers(t *testing.T) {
	sm := shutdown.New()

	var redisClosed atomic.Bool
	var natsClosed atomic.Bool

	sm.RegisterFunc("fiber", 0, func(_ context.Context) error { return nil })

	sm.RegisterFunc("mongodb", 5, func(_ context.Context) error {
		return errors.New("mongo: connection reset by peer")
	})
	sm.RegisterFunc("redis-app", 5, func(_ context.Context) error {
		redisClosed.Store(true)
		return nil
	})
	sm.RegisterFunc("nats", 5, func(_ context.Context) error {
		natsClosed.Store(true)
		return nil
	})

	sendSIGINT(50 * time.Millisecond)
	sm.WaitForSignal(15 * time.Second)

	if !redisClosed.Load() {
		t.Fatal("redis should close despite mongo error")
	}
	if !natsClosed.Load() {
		t.Fatal("nats should close despite mongo error")
	}
}

// TestShutdownIntegration_CompletesWithinTimeout verifies the full
// shutdown sequence finishes well under the 15s production timeout.
func TestShutdownIntegration_CompletesWithinTimeout(t *testing.T) {
	sm := shutdown.New()

	sm.RegisterFunc("fiber", 0, func(_ context.Context) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	sm.RegisterFunc("mongodb", 5, func(_ context.Context) error {
		time.Sleep(20 * time.Millisecond)
		return nil
	})
	sm.RegisterFunc("redis-app", 5, func(_ context.Context) error {
		time.Sleep(15 * time.Millisecond)
		return nil
	})
	sm.RegisterFunc("redis-shared", 5, func(_ context.Context) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	sm.RegisterFunc("nats", 5, func(_ context.Context) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	sendSIGINT(50 * time.Millisecond)
	start := time.Now()
	sm.WaitForSignal(15 * time.Second)
	elapsed := time.Since(start)

	if elapsed > 1*time.Second {
		t.Fatalf("shutdown took too long: %s (expected <1s)", elapsed)
	}
}

// TestShutdownIntegration_SlowFiberDrainsBeforeConnectionClose simulates
// Fiber taking time to drain in-flight requests. Connections must wait.
func TestShutdownIntegration_SlowFiberDrainsBeforeConnectionClose(t *testing.T) {
	sm := shutdown.New()
	tracker := newHookTracker()

	// P0: Fiber drains slowly (100ms)
	sm.RegisterFunc("fiber", 0, func(_ context.Context) error {
		tracker.record("fiber", 100*time.Millisecond)
		return nil
	})

	// P5: Quick connection closes
	sm.RegisterFunc("mongodb", 5, func(_ context.Context) error {
		tracker.record("mongodb", 5*time.Millisecond)
		return nil
	})
	sm.RegisterFunc("redis-app", 5, func(_ context.Context) error {
		tracker.record("redis-app", 5*time.Millisecond)
		return nil
	})

	sendSIGINT(50 * time.Millisecond)
	sm.WaitForSignal(15 * time.Second)

	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	fiberEnd := tracker.endedAt["fiber"]
	for _, name := range []string{"mongodb", "redis-app"} {
		if tracker.startedAt[name].Before(fiberEnd) {
			t.Fatalf("%s started at %v, before fiber ended at %v", name, tracker.startedAt[name], fiberEnd)
		}
	}
}

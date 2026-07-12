package infra

import (
	"context"
	"fmt"
	"time"

	constants "github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
)

const (
	bringUpTimeout = 60 * time.Second
	bringUpTick    = time.Second
)

// waitReady polls the service's probe until ready or a bounded deadline (scaled by
// SAGA_TIMEOUT_MULTIPLIER for slow environments). Called after a docker-compose
// bring-up to confirm the service actually came up.
func waitReady(ctx context.Context, svc string) error {
	timeout := constants.ScaleTimeout(bringUpTimeout)
	deadline := time.Now().Add(timeout)
	for {
		if ready, _ := probe(svc); ready {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("infra: %q not ready within %v after bring-up", svc, timeout)
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("infra: %q readiness wait cancelled: %w", svc, ctx.Err())
		case <-time.After(bringUpTick):
		}
	}
}

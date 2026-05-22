package steps

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// StartTestSink starts a local HTTP server bound to
// constants.TriggerSinkBindAddr that responds 200 OK to every POST.
// The connectivity-action phase2_trigger journey adds this as its
// first step so the triggers service has a real responder to POST
// against — the resulting events_trigger row carries success=true,
// which is what the saga's downstream assert expects.
//
// Writes (bag):
//   - BagKeyTriggerSinkServer  *http.Server   for Compensate to stop.
//   - BagKeyTriggerSinkHits    *atomic.Int64  POST counter, optional
//     fast path for asserts that want to confirm the sink saw the
//     traffic without polling /api/v1/events/trigger.
//
// Compensate: graceful shutdown of the sink with a 2s deadline. Safe
// to call even when Do failed before publishing the server (the bag
// lookup short-circuits).
//
// Linux note: when triggers MS runs in Docker, it must reach the
// host-bound sink. Override SAGA_TRIGGER_SINK_URL to a host the
// container can resolve (host.docker.internal — already exposed by
// the standalone compose's `extra_hosts: host-gateway`).
func StartTestSink() saga.Step {
	return saga.Step{
		Name: "triggers/triggers.StartTestSink",
		Do: func(c *saga.Context) error {
			hits := &atomic.Int64{}

			mux := http.NewServeMux()
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					hits.Add(1)
				}
				w.WriteHeader(http.StatusOK)
				// "ok" literal matches Slack's incoming-webhook contract;
				// HTTP / Teams executors don't validate the body so the
				// shared response is harmless for them.
				_, _ = w.Write([]byte("ok"))
			})

			ln, err := net.Listen("tcp", constants.TriggerSinkBindAddr)
			if err != nil {
				return fmt.Errorf("listen %s: %w", constants.TriggerSinkBindAddr, err)
			}

			srv := &http.Server{
				Handler:           mux,
				ReadHeaderTimeout: 5 * time.Second,
			}
			go func() { _ = srv.Serve(ln) }()

			c.Set(BagKeyTriggerSinkServer, srv)
			c.Set(BagKeyTriggerSinkHits, hits)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			v, ok := c.Get(BagKeyTriggerSinkServer)
			if !ok {
				return nil
			}
			srv, ok := v.(*http.Server)
			if !ok {
				return nil
			}
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_ = srv.Shutdown(shutdownCtx)
			return nil
		},
	}
}

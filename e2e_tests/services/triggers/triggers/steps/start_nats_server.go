package steps

import (
	"fmt"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// StartNatsServer boots an embedded NATS server (nats-server/v2) on a
// free ephemeral port without auth. The trigger payload reads the
// resulting URL from the bag so the trigger can publish to the saga
// own server, isolated from the platform NATS that requires auth.
//
// Writes (bag):
//   - BagKeyNatsServer  *natsserver.Server  for Compensate to Shutdown.
//   - BagKeyNatsURL     string              "nats://127.0.0.1:<port>".
//
// Compensate: server.Shutdown() + WaitForShutdown. Idempotent.
func StartNatsServer() saga.Step {
	return saga.Step{
		Name: "triggers/triggers.StartNatsServer",
		Do: func(c *saga.Context) error {
			opts := &natsserver.Options{
				Host:           "127.0.0.1",
				Port:           -1, // ask the runtime for a free port
				NoLog:          true,
				NoSigs:         true,
				MaxControlLine: 4096,
			}
			srv, err := natsserver.NewServer(opts)
			if err != nil {
				return fmt.Errorf("nats new server: %w", err)
			}
			go srv.Start()
			if !srv.ReadyForConnections(5 * time.Second) {
				srv.Shutdown()
				return fmt.Errorf("nats not ready within 5s")
			}

			c.Set(BagKeyNatsServer, srv)
			c.Set(BagKeyNatsURL, srv.ClientURL())
			return nil
		},
		Compensate: func(c *saga.Context) error {
			v, ok := c.Get(BagKeyNatsServer)
			if !ok {
				return nil
			}
			srv, ok := v.(*natsserver.Server)
			if !ok {
				return nil
			}
			srv.Shutdown()
			return nil
		},
	}
}

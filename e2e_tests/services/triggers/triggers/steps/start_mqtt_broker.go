package steps

import (
	"fmt"
	"io"
	"log/slog"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/netx"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// StartMqttBroker boots an in-process MQTT broker on a free ephemeral
// port using mochi-mqtt's AllowHook (no auth). Publishes the host and
// port on the bag so the trigger payload can dial the just-started
// broker without depending on the platform's authenticated broker.
//
// Writes (bag):
//   - BagKeyMqttBroker      *mqtt.Server  for Compensate to stop.
//   - BagKeyMqttBrokerHost  string        advertise host (constants.SinkHost).
//   - BagKeyMqttBrokerPort  int           OS-assigned port.
//
// Compensate: server.Close(). Idempotent.
func StartMqttBroker() saga.Step {
	return saga.Step{
		Name: "triggers/triggers.StartMqttBroker",
		Do: func(c *saga.Context) error {
			// Bind an ephemeral port and hand mochi the LIVE listener — never
			// close it and let mochi re-bind the same address (that check-then-bind
			// gap is a TOCTOU race another parallel journey could win).
			ln, port, err := netx.FreeListener()
			if err != nil {
				return fmt.Errorf("listen ephemeral mqtt port: %w", err)
			}

			server := mqtt.New(&mqtt.Options{
				// Silence the broker; saga test stdout already carries the run log.
				Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
			})
			if err := server.AddHook(new(auth.AllowHook), nil); err != nil {
				return fmt.Errorf("mqtt allow hook: %w", err)
			}
			if err := server.AddListener(listeners.NewNet("saga-mqtt", ln)); err != nil {
				return fmt.Errorf("mqtt add listener: %w", err)
			}
			go func() { _ = server.Serve() }()

			c.Set(BagKeyMqttBroker, server)
			c.Set(BagKeyMqttBrokerHost, constants.SinkHost)
			c.Set(BagKeyMqttBrokerPort, port)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			v, ok := c.Get(BagKeyMqttBroker)
			if !ok {
				return nil
			}
			srv, ok := v.(*mqtt.Server)
			if !ok {
				return nil
			}
			return srv.Close()
		},
	}
}

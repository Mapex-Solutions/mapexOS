package steps

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// StartMqttBroker boots an in-process MQTT broker on a free ephemeral
// port using mochi-mqtt's AllowHook (no auth). Publishes the host and
// port on the bag so the trigger payload can dial the just-started
// broker without depending on the platform's authenticated broker.
//
// Writes (bag):
//   - BagKeyMqttBroker      *mqtt.Server  for Compensate to stop.
//   - BagKeyMqttBrokerHost  string        bind host ("127.0.0.1").
//   - BagKeyMqttBrokerPort  int           OS-assigned port.
//
// Compensate: server.Close(). Idempotent.
func StartMqttBroker() saga.Step {
	return saga.Step{
		Name: "triggers/triggers.StartMqttBroker",
		Do: func(c *saga.Context) error {
			// Pre-bind to find a free port so we can publish the
			// host:port on the bag before the broker starts serving.
			// mochi listeners.Init opens the bind so we can read
			// Address() afterwards anyway, but resolving the port up
			// front keeps the bag write atomic with the bind.
			ln, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				return fmt.Errorf("listen ephemeral mqtt port: %w", err)
			}
			addr := ln.Addr().String()
			_ = ln.Close() // mochi opens its own listener on the same addr.

			host, portStr, err := net.SplitHostPort(addr)
			if err != nil {
				return fmt.Errorf("parse mqtt addr %q: %w", addr, err)
			}
			port, err := strconv.Atoi(portStr)
			if err != nil {
				return fmt.Errorf("parse mqtt port %q: %w", portStr, err)
			}

			server := mqtt.New(&mqtt.Options{
				// Silence the broker; saga test stdout already carries the run log.
				Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
			})
			if err := server.AddHook(new(auth.AllowHook), nil); err != nil {
				return fmt.Errorf("mqtt allow hook: %w", err)
			}
			tcp := listeners.NewTCP(listeners.Config{
				ID:      "saga-mqtt",
				Address: addr,
			})
			if err := server.AddListener(tcp); err != nil {
				return fmt.Errorf("mqtt add listener: %w", err)
			}
			go func() { _ = server.Serve() }()

			c.Set(BagKeyMqttBroker, server)
			c.Set(BagKeyMqttBrokerHost, host)
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

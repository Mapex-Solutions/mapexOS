package steps

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// StartRabbitmqContainer spins up an ephemeral RabbitMQ container via
// testcontainers-go on a random port (guest/guest, no extra plugins).
// The trigger payload reads the parsed host/port/user/pass from the
// bag so the saga's trigger publishes to a container the journey owns,
// isolated from any RabbitMQ the developer happens to have running.
//
// Writes (bag):
//   - BagKeyRabbitmqContainer  *rabbitmq.RabbitMQContainer  for Compensate.
//   - BagKeyRabbitmqHost       string                       container host.
//   - BagKeyRabbitmqPort       int                          mapped AMQP port.
//   - BagKeyRabbitmqUser       string                       "guest".
//   - BagKeyRabbitmqPass       string                       "guest".
//
// Compensate: container.Terminate(). First start can take 10–30 s while
// docker pulls the image; subsequent runs reuse the cached image.
func StartRabbitmqContainer() saga.Step {
	return saga.Step{
		Name: "triggers/triggers.StartRabbitmqContainer",
		Do: func(c *saga.Context) error {
			ctnr, err := rabbitmq.Run(c.Stdctx, "rabbitmq:3.12.11-management-alpine")
			if err != nil {
				return fmt.Errorf("start rabbitmq container: %w", err)
			}
			amqpURL, err := ctnr.AmqpURL(c.Stdctx)
			if err != nil {
				_ = ctnr.Terminate(context.Background())
				return fmt.Errorf("get rabbitmq amqpURL: %w", err)
			}
			// amqpURL shape: "amqp://guest:guest@host:port"
			user, pass, hostPort, err := splitAmqpURL(amqpURL)
			if err != nil {
				_ = ctnr.Terminate(context.Background())
				return fmt.Errorf("parse rabbitmq amqpURL %q: %w", amqpURL, err)
			}
			host, portStr, err := net.SplitHostPort(hostPort)
			if err != nil {
				_ = ctnr.Terminate(context.Background())
				return fmt.Errorf("parse rabbitmq host:port %q: %w", hostPort, err)
			}
			port, err := strconv.Atoi(portStr)
			if err != nil {
				_ = ctnr.Terminate(context.Background())
				return fmt.Errorf("parse rabbitmq port %q: %w", portStr, err)
			}

			c.Set(BagKeyRabbitmqContainer, ctnr)
			c.Set(BagKeyRabbitmqHost, host)
			c.Set(BagKeyRabbitmqPort, port)
			c.Set(BagKeyRabbitmqUser, user)
			c.Set(BagKeyRabbitmqPass, pass)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			v, ok := c.Get(BagKeyRabbitmqContainer)
			if !ok {
				return nil
			}
			ctnr, ok := v.(*rabbitmq.RabbitMQContainer)
			if !ok {
				return nil
			}
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			return ctnr.Terminate(ctx)
		},
	}
}

// splitAmqpURL parses "amqp://user:pass@host:port[/vhost]" into its
// pieces. The testcontainers helper only returns the canonical
// `amqp://guest:guest@host:port` shape so we lean on string slicing
// rather than pulling in a URL parser.
func splitAmqpURL(amqpURL string) (user, pass, hostPort string, err error) {
	const prefix = "amqp://"
	if !strings.HasPrefix(amqpURL, prefix) {
		return "", "", "", fmt.Errorf("missing %q prefix", prefix)
	}
	rest := strings.TrimPrefix(amqpURL, prefix)

	atIdx := strings.LastIndex(rest, "@")
	if atIdx < 0 {
		return "", "", "", fmt.Errorf("missing '@' separator")
	}
	creds := rest[:atIdx]
	hostPort = rest[atIdx+1:]

	colonIdx := strings.Index(creds, ":")
	if colonIdx < 0 {
		return "", "", "", fmt.Errorf("missing ':' between user and pass")
	}
	user = creds[:colonIdx]
	pass = creds[colonIdx+1:]

	// Trim optional /vhost suffix on host:port.
	if slashIdx := strings.Index(hostPort, "/"); slashIdx >= 0 {
		hostPort = hostPort[:slashIdx]
	}
	return user, pass, hostPort, nil
}

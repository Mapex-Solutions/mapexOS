package steps

import (
	"fmt"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// bagPort reads an int port previously published to the bag by a sink step,
// returning a clear error when the key is missing or holds a non-int (which means
// the sink step did not run before the create step).
func bagPort(c *saga.Context, key string) (int, error) {
	v, ok := c.Get(key)
	if !ok {
		return 0, fmt.Errorf("sink port: bag key %q missing (sink step not run?)", key)
	}
	port, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("sink port: bag key %q is not int (%T)", key, v)
	}
	return port, nil
}

// httpSinkURL builds the http://host:port callback URL for the in-process HTTP sink
// from the ephemeral address StartTestSink published on the bag. The HTTP, Slack, and
// Teams triggers all POST to this same sink, so their create steps share this helper.
//
// Reads (bag):
//   - BagKeyTriggerSinkHost  string  advertise host (constants.SinkHost).
//   - BagKeyTriggerSinkPort  int     OS-assigned ephemeral port.
func httpSinkURL(c *saga.Context) (string, error) {
	host := c.MustGetString(BagKeyTriggerSinkHost)
	port, err := bagPort(c, BagKeyTriggerSinkPort)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s:%d", host, port), nil
}

// wsSinkURL builds the ws://host:port/ws callback URL for the in-process WebSocket
// sink from the ephemeral address StartWebsocketSink published on the bag.
//
// Reads (bag):
//   - BagKeyWsHost  string  advertise host (constants.SinkHost).
//   - BagKeyWsPort  int     OS-assigned ephemeral port.
func wsSinkURL(c *saga.Context) (string, error) {
	host := c.MustGetString(BagKeyWsHost)
	port, err := bagPort(c, BagKeyWsPort)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("ws://%s:%d/ws", host, port), nil
}

// smtpSinkAddr returns the advertise host and ephemeral port StartSmtpSink published
// on the bag. The email trigger config carries host and port as separate fields, so
// this returns them split rather than as a URL.
//
// Reads (bag):
//   - BagKeySmtpHost  string  advertise host (constants.SinkHost).
//   - BagKeySmtpPort  int     OS-assigned ephemeral port.
func smtpSinkAddr(c *saga.Context) (host string, port int, err error) {
	host = c.MustGetString(BagKeySmtpHost)
	port, err = bagPort(c, BagKeySmtpPort)
	return host, port, err
}

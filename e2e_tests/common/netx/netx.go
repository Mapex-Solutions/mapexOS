// Package netx holds small, generic networking helpers shared across the e2e
// building blocks. It depends on nothing in the suite (no saga engine, no
// service), so any step that starts an in-process server can reuse it.
package netx

import (
	"fmt"
	"net"
	"strconv"
)

// FreeListener binds an OS-assigned free TCP port on ALL interfaces (0.0.0.0) and
// returns the open listener plus its port. Binding "0.0.0.0:0" lets the OS pick a
// guaranteed-free port atomically (no check-then-bind race) AND makes the server
// reachable both from the host and from a service running in Docker (which reaches
// the host over the bridge). The advertise host a caller puts in a callback URL is a
// separate policy decision (loopback for local, host.docker.internal for Docker) —
// not this listener's bind address — so FreeListener returns only the port. The
// caller owns the listener and must serve or Close it.
func FreeListener() (ln net.Listener, port int, err error) {
	ln, err = net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, 0, fmt.Errorf("netx: free listener: %w", err)
	}

	addr := ln.Addr().String()
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		_ = ln.Close()
		return nil, 0, fmt.Errorf("netx: parse listener addr %q: %w", addr, err)
	}

	port, err = strconv.Atoi(portStr)
	if err != nil {
		_ = ln.Close()
		return nil, 0, fmt.Errorf("netx: parse listener port %q: %w", portStr, err)
	}

	return ln, port, nil
}

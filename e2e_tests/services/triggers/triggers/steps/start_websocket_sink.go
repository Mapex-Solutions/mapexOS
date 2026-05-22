package steps

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

// StartWebsocketSink boots a minimal WebSocket server on
// constants.WsSinkBindAddr. /ws is the upgrade endpoint; on each
// successful handshake the sink increments the hit counter, reads
// frames into the bag's last-message slot, and waits for the client
// (the WebSocket trigger executor) to close. The smoke only needs
// the handshake to succeed for events_trigger to mark success=true,
// but capturing the first frame lets a content-key assert validate
// the message bytes when desired.
//
// Writes (bag):
//   - BagKeyWsServer       *http.Server  for Compensate to stop.
//   - BagKeyWsHits         *atomic.Int64 incremented on each handshake.
//   - BagKeyWsLastMessage  **string      pointer slot for the most
//     recently received frame payload (text or binary, encoded as the
//     literal UTF-8 bytes paho-style).
func StartWebsocketSink() saga.Step {
	return saga.Step{
		Name: "triggers/triggers.StartWebsocketSink",
		Do: func(c *saga.Context) error {
			hits := &atomic.Int64{}
			lastPtr := new(*string)

			mux := http.NewServeMux()
			mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
				conn, err := wsUpgrader.Upgrade(w, r, nil)
				if err != nil {
					return
				}
				hits.Add(1)
				// Read one frame so the executor's write completes
				// cleanly; ignore subsequent frames and shut down.
				_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
				_, msg, rerr := conn.ReadMessage()
				if rerr == nil {
					s := string(msg)
					*lastPtr = &s
				}
				_ = conn.Close()
			})
			// Health endpoint so a journey can sanity-check the sink is
			// up before letting the trigger fire.
			mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ws-sink ok"))
			})

			ln, err := net.Listen("tcp", constants.WsSinkBindAddr)
			if err != nil {
				return fmt.Errorf("listen %s: %w", constants.WsSinkBindAddr, err)
			}
			srv := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
			go func() { _ = srv.Serve(ln) }()

			c.Set(BagKeyWsServer, srv)
			c.Set(BagKeyWsHits, hits)
			c.Set(BagKeyWsLastMessage, lastPtr)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			v, ok := c.Get(BagKeyWsServer)
			if !ok {
				return nil
			}
			srv, ok := v.(*http.Server)
			if !ok {
				return nil
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_ = srv.Shutdown(ctx)
			return nil
		},
	}
}

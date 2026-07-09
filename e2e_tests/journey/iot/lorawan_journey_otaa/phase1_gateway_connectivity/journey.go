// Package phase1_gateway_connectivity exercises a LoRaWAN gateway asset's
// connectivity lifecycle end to end against the live stack, over BOTH radio
// transports (Semtech UDP and Basics Station key mode).
//
// The gateway provision + connect + presence sequences are the shared
// common/journey/lorawan_gateway building block — this phase just composes them.
// A LoRaWAN gateway asset needs neither an asset template nor a route group.
//
// Outcome on PASS (per transport block):
//   - A gateway asset is provisioned and served in the LNS-facing L3 shape.
//   - The simulated gateway connects to mapexLNS; its keepalive/stats make the
//     LNS publish a "connect" presence advisory and the asset goes online.
//   - ForceOfflineByAdmin (internal ops endpoint) flips it back to offline.
//   - Every gateway (labelled per transport) is torn down in compensation.
//
// Outcome on FAIL:
//   - The failing step / assert name surfaces in the saga log.
package phase1_gateway_connectivity

import (
	"context"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	bootstrap "github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/iam_bootstrap"
	gateway "github.com/Mapex-Solutions/MapexOS/e2eTests/common/journey/lorawan_gateway"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/utils"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// Items composes the shared gateway connect + presence sequences for a UDP
// gateway and a Basics Station gateway.
func Items() []saga.Item {
	var items []saga.Item
	// UDP transport (eui-mode gateway): connect → online → force offline → offline.
	items = append(items, gateway.ConnectUDPItems("gwUdp")...)
	items = append(items, gateway.PresenceItems("gwUdp")...)
	// Basics Station transport (key-mode gateway, bearer token): same lifecycle.
	items = append(items, gateway.ConnectBasicStationItems("gwBs")...)
	items = append(items, gateway.PresenceItems("gwBs")...)
	return items
}

// Run fronts the shared IAM bootstrap (common/journey/iam_bootstrap) and runs
// this phase's items under one rollback chain.
func Run(t *testing.T) {
	t.Helper()
	if err := utils.SetupE2EEnvironment(); err != nil {
		t.Fatalf("setup e2e environment: %v", err)
	}
	runID := random.NewRunID()
	clients := bootstrap.NewClients()
	items := append(bootstrap.BootstrapItems(), Items()...)
	saga.Run(t, context.Background(), runID, clients, items...)
}

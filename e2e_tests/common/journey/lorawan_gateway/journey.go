// Package lorawan_gateway is the shared LoRaWAN gateway building block. It lives
// under common/journey/ because more than one journey (lorawan_journey_otaa,
// lorawan_journey_bsp) reuses the same "provision + connect a gateway" and
// "gateway presence" item sequences. A journey must never import another
// journey's phase, so the shared sequence lives here — the sensor journeys carry
// the uplink over a gateway these items stand up.
//
// The functions return item slices a journey composes into its own saga.Run.
// They connect a gateway under gwLabel and leave it on the bag (via the assets
// step package's label-scoped keys) for the sensor steps to ride.
package lorawan_gateway

import (
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"

	assetAsserts "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/asserts"
	assetPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
)

// ConnectUDPItems provisions an eui-mode gateway asset under gwLabel, asserts its
// LNS-facing L3 auth projection, and connects the simulated gateway over the
// Semtech UDP ingress — leaving a live gateway on the bag for sensor steps to ride.
func ConnectUDPItems(gwLabel string) []saga.Item {
	return []saga.Item{
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanGatewayFor(gwLabel), gwLabel),
		assetAsserts.AssertAuthProjectionGatewayByLabel(gwLabel, assetPayloads.SagaGatewayFrequencyPlan),
		assetSteps.ConnectLorawanGatewayUDP(gwLabel),
		// Gate on the gateway going online before any sensor rides it. The Semtech-UDP
		// frontend drops packets that arrive while the gateway connection is still being
		// registered (no retry, no queue), so an uplink fired in that window is lost. The
		// gateway flips online only after the LNS has registered the connection.
		assetAsserts.AssertHealthStatusByLabel(gwLabel, "online"),
	}
}

// ConnectBasicStationItems provisions a key-mode gateway asset under gwLabel,
// asserts its L3 auth projection, and connects the simulated gateway over the
// Basics Station WebSocket ingress with the provisioned bearer token.
func ConnectBasicStationItems(gwLabel string) []saga.Item {
	return []saga.Item{
		assetSteps.CreateAssetWithLabel(assetPayloads.SagaLorawanGatewayKeyFor(gwLabel), gwLabel),
		assetAsserts.AssertAuthProjectionGatewayByLabel(gwLabel, assetPayloads.SagaGatewayFrequencyPlan),
		assetSteps.ConnectLorawanGatewayBasicStation(gwLabel),
		// Same readiness gate as the UDP path: wait until the LNS has registered the
		// gateway (online) before a sensor rides it, so the first uplink is not raced
		// against connection setup.
		assetAsserts.AssertHealthStatusByLabel(gwLabel, "online"),
	}
}

// PresenceItems asserts the connected gateway went online, forces it offline via
// the internal ops endpoint, and asserts it flipped to offline. Compose it after
// a Connect*Items sequence for the same gwLabel to exercise the presence lifecycle.
func PresenceItems(gwLabel string) []saga.Item {
	return []saga.Item{
		assetAsserts.AssertHealthStatusByLabel(gwLabel, "online"),
		assetSteps.ForceOfflineByAdminByLabel(gwLabel, "saga-lorawan-gw-offline"),
		assetAsserts.AssertHealthStatusByLabel(gwLabel, "offline"),
	}
}

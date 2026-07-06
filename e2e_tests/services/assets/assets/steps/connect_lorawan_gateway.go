package steps

import (
	"fmt"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/lorawansim"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
)

// ConnectLorawanGatewayUDP dials the labelled gateway asset at the mapexLNS Semtech
// UDP ingress and keeps the link open — its keepalive + stats loops make the
// gateway show online on the LNS. The gateway EUI is the asset's AssetUUID
// (SagaLorawanGatewayFor sets AssetUUID = the EUI).
//
// Reads (bag):
//   - AssetUUIDKey(gwLabel)  string  the gateway EUI, set by CreateAssetWithLabel
//
// Writes (bag):
//   - LorawanGatewayKey(gwLabel)  *lorawansim.Gateway
//
// Compensate: closes the gateway link.
func ConnectLorawanGatewayUDP(gwLabel string) saga.Step {
	return saga.Step{
		Name: "assets/assets.ConnectLorawanGatewayUDP[" + gwLabel + "]",
		Do: func(c *saga.Context) error {
			eui := c.MustGetString(AssetUUIDKey(gwLabel))
			gw, err := lorawansim.NewGateway(c.Stdctx, lorawansim.GatewayConfig{
				EUI:  eui,
				Link: "udp",
				Host: constants.LNSUDPHost,
				Port: constants.LNSUDPPort,
			})
			if err != nil {
				return fmt.Errorf("connect lorawan gateway (udp) %q: %w", gwLabel, err)
			}
			c.Set(LorawanGatewayKey(gwLabel), gw)
			return nil
		},
		Compensate: closeLorawanGateway(gwLabel),
	}
}

// ConnectLorawanGatewayBasicStation dials the labelled key-mode gateway at the
// mapexLNS Basics Station WebSocket ingress, presenting the bearer token the
// gateway asset was provisioned with (SagaGatewayAPIKey). The LNS bcrypt-compares
// it against the gateway's stored apiKeyHash.
func ConnectLorawanGatewayBasicStation(gwLabel string) saga.Step {
	return saga.Step{
		Name: "assets/assets.ConnectLorawanGatewayBasicStation[" + gwLabel + "]",
		Do: func(c *saga.Context) error {
			eui := c.MustGetString(AssetUUIDKey(gwLabel))
			gw, err := lorawansim.NewGateway(c.Stdctx, lorawansim.GatewayConfig{
				EUI:    eui,
				Link:   "basicstation",
				LNSURI: constants.LNSBStationURI,
				Token:  payloads.SagaGatewayAPIKey,
			})
			if err != nil {
				return fmt.Errorf("connect lorawan gateway (basicstation) %q: %w", gwLabel, err)
			}
			c.Set(LorawanGatewayKey(gwLabel), gw)
			return nil
		},
		Compensate: closeLorawanGateway(gwLabel),
	}
}

// closeLorawanGateway returns the shared Compensate that closes the labelled
// gateway link if it was opened.
func closeLorawanGateway(gwLabel string) func(*saga.Context) error {
	return func(c *saga.Context) error {
		if v, ok := c.Get(LorawanGatewayKey(gwLabel)); ok {
			_ = v.(*lorawansim.Gateway).Close()
		}
		return nil
	}
}

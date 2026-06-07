package services

import (
	"testing"

	"assets/src/modules/assets/domain/entities"

	assetsAuthContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/auth"
)

func TestBuildAuthProjection(t *testing.T) {
	lat, lon := 1.5, 2.5

	tests := []struct {
		name   string
		asset  *entities.Asset
		verify func(t *testing.T, proj assetsAuthContract.AuthProjection)
	}{
		{
			name: "lorawan gateway projects the gateway block and no keys",
			asset: &entities.Asset{
				AssetUUID: "0102030405060708",
				Enabled:   true,
				Protocol: entities.ProtocolType{
					Type: "lorawan",
					Lorawan: &entities.LorawanConfig{
						Kind: "gateway",
						Gateway: &entities.LorawanGatewayConfig{
							AuthMode:        "cert",
							FrequencyPlanID: "AU_915_928_FSB_2_BR",
							Latitude:        &lat,
							Longitude:       &lon,
						},
					},
				},
				CurrentCert: &entities.AssetCertificate{Serial: "ABCD1234"},
			},
			verify: func(t *testing.T, proj assetsAuthContract.AuthProjection) {
				if proj.Type != "lorawan" {
					t.Fatalf("Type = %q, want lorawan", proj.Type)
				}
				if proj.Mqtt != nil {
					t.Errorf("Mqtt must be nil for a gateway")
				}
				if proj.Lorawan == nil || proj.Lorawan.Kind != "gateway" {
					t.Fatalf("Lorawan block missing or wrong kind: %+v", proj.Lorawan)
				}
				if proj.Lorawan.Keys.EncryptedDEK != nil || proj.Lorawan.Keys.EncryptedKey != nil {
					t.Errorf("gateway must carry no key material, got %+v", proj.Lorawan.Keys)
				}
				gw := proj.Lorawan.Gateway
				if gw == nil {
					t.Fatalf("gateway sub-block missing")
				}
				if gw.AuthMode != "cert" || gw.FrequencyPlanID != "AU_915_928_FSB_2_BR" {
					t.Errorf("gateway block = %+v", gw)
				}
				if gw.CurrentCertSerial != "ABCD1234" {
					t.Errorf("CurrentCertSerial = %q, want ABCD1234 (from asset.CurrentCert)", gw.CurrentCertSerial)
				}
				if gw.Latitude == nil || *gw.Latitude != 1.5 {
					t.Errorf("Latitude not propagated: %v", gw.Latitude)
				}
			},
		},
		{
			name: "lorawan device projects identity + keys and no gateway block",
			asset: &entities.Asset{
				AssetUUID: "1122334455667788",
				Enabled:   true,
				Protocol: entities.ProtocolType{
					Type: "lorawan",
					Lorawan: &entities.LorawanConfig{
						Kind:       "device",
						DevEUI:     "1122334455667788",
						Region:     "EU_863_870",
						Activation: "otaa",
						Keys: entities.EncryptedKeys{
							EncryptedDEK: []byte("dek"),
							EncryptedKey: []byte("key"),
						},
					},
				},
			},
			verify: func(t *testing.T, proj assetsAuthContract.AuthProjection) {
				if proj.Type != "lorawan" {
					t.Fatalf("Type = %q, want lorawan", proj.Type)
				}
				if proj.Lorawan == nil || proj.Lorawan.Kind != "device" {
					t.Fatalf("device Lorawan block missing or wrong kind: %+v", proj.Lorawan)
				}
				if proj.Lorawan.Gateway != nil {
					t.Errorf("device must not carry a gateway block")
				}
				if proj.Lorawan.DevEUI != "1122334455667788" {
					t.Errorf("DevEUI = %q", proj.Lorawan.DevEUI)
				}
				if string(proj.Lorawan.Keys.EncryptedDEK) != "dek" {
					t.Errorf("device must carry the encrypted key material")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.verify(t, buildAuthProjection(tt.asset))
		})
	}
}

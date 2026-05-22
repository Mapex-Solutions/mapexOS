package steps

import (
	"fmt"
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/mqttclient"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// ConnectMqttPassword opens an MQTT CONNECT against the platform
// broker using username=assetUUID + the plaintext password the saga
// supplied on create. Stores the live *mqttclient.Client on the bag
// so Publish/Disconnect steps re-use the same socket.
//
// Reads (bag):
//   - BagKeyAssetUUID         string  set by CreateAsset
//   - BagKeyAssetMqttPassword string  set by CreateAsset
//
// Writes (bag):
//   - BagKeyMqttClient        *mqttclient.Client
//   - BagKeyMqttConnectedAt   time.Time
//
// Compensate: best-effort Disconnect so a panicking journey leaves no
// dangling socket against the broker.
func ConnectMqttPassword() saga.Step {
	return saga.Step{
		Name: "assets/assets.ConnectMqttPassword",
		Do: func(c *saga.Context) error {
			uuid := c.MustGetString(BagKeyAssetUUID)
			pwd := c.MustGetString(BagKeyAssetMqttPassword)

			cli, err := mqttclient.New(mqttclient.Config{
				BrokerURL: constants.MqttBrokerURL,
				ClientID:  uuid,
				Username:  uuid,
				Password:  pwd,
			})
			if err != nil {
				return fmt.Errorf("build mqtt client: %w", err)
			}
			if err := cli.Connect(c.Stdctx); err != nil {
				return fmt.Errorf("mqtt connect (password) for %s: %w", uuid, err)
			}
			c.Set(BagKeyMqttClient, cli)
			c.Set(BagKeyMqttConnectedAt, time.Now().UTC())
			return nil
		},
		Compensate: func(c *saga.Context) error {
			v, ok := c.Get(BagKeyMqttClient)
			if !ok {
				return nil
			}
			cli, ok := v.(*mqttclient.Client)
			if !ok {
				return nil
			}
			cli.Disconnect(0)
			return nil
		},
	}
}

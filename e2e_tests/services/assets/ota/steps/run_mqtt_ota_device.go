package steps

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/mqttclient"

	downlink "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/downlink"
	otaDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/dtos"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
)

// otaStatusProgression is the ordered report sequence both device sims publish,
// mirroring a real device's download → verify → apply lifecycle.
var otaStatusProgression = []struct {
	status   string
	progress int32
}{
	{"downloading", 20},
	{"downloaded", 40},
	{"verified", 60},
	{"updating", 80},
	{"updated", 100},
}

// SubscribeOtaCommand subscribes the connected MQTT sim to its OTA command topic
// (commands/{assetUUID}/ota_update) and buffers the received OTAUpdateCommand on
// the bag. It MUST run BEFORE CreatePlan: the reconciler pushes the command once,
// to an online device, so the subscription has to be live before dispatch.
//
// Reads (bag):
//   - assetSteps.BagKeyMqttClient  *mqttclient.Client  set by ConnectMqttPassword
//   - assetSteps.BagKeyAssetUUID   string              set by CreateAsset
//
// Writes (bag):
//   - bagKeyOtaCmdChan  chan downlink.OTAUpdateCommand (buffered, size 1)
//
// Compensate: no-op — the socket is torn down by ConnectMqttPassword's Compensate.
func SubscribeOtaCommand() saga.Step {
	return saga.Step{
		Name: "assets/ota.SubscribeOtaCommand",
		Do: func(c *saga.Context) error {
			cli, err := mqttClientFromBag(c)
			if err != nil {
				return err
			}
			uuid := c.MustGetString(assetSteps.BagKeyAssetUUID)
			// The broker delivers device commands on commands/{assetUUID}/{commandType}.
			topic := "commands/" + uuid + "/" + downlink.CommandTypeOTAUpdate

			cmdCh := make(chan downlink.OTAUpdateCommand, 1)
			c.Set(bagKeyOtaCmdChan, cmdCh)

			handler := func(_ string, payload []byte) {
				var cmd downlink.OTAUpdateCommand
				if err := json.Unmarshal(payload, &cmd); err != nil {
					return
				}
				select {
				case cmdCh <- cmd:
				default: // already buffered; ignore duplicate deliveries
				}
			}
			if err := cli.Subscribe(c.Stdctx, topic, 1, handler); err != nil {
				return fmt.Errorf("subscribe ota command on %q: %w", topic, err)
			}
			return nil
		},
		Compensate: func(_ *saga.Context) error { return nil },
	}
}

// RunMqttOtaDevice waits for the pushed OTA command, downloads + verifies the
// firmware, records the executionId, and publishes the status progression on
// events/{assetUUID}/ota_status. It MUST run AFTER CreatePlan (which triggers the
// dispatch) and after SubscribeOtaCommand (which armed the buffer).
//
// Reads (bag):
//   - assetSteps.BagKeyMqttClient / assetSteps.BagKeyAssetUUID / bagKeyOtaCmdChan
//
// Writes (bag):
//   - BagKeyExecutionID  string  the execution the device reports against
//
// Compensate: no-op.
func RunMqttOtaDevice() saga.Step {
	return saga.Step{
		Name: "assets/ota.RunMqttOtaDevice",
		Do: func(c *saga.Context) error {
			cli, err := mqttClientFromBag(c)
			if err != nil {
				return err
			}
			uuid := c.MustGetString(assetSteps.BagKeyAssetUUID)

			cmdCh, err := otaCmdChanFromBag(c)
			if err != nil {
				return err
			}

			var cmd downlink.OTAUpdateCommand
			select {
			case cmd = <-cmdCh:
			case <-time.After(45 * time.Second):
				return fmt.Errorf("run mqtt ota device: no ota_update command received within 45s for %s", uuid)
			case <-c.Stdctx.Done():
				return fmt.Errorf("run mqtt ota device: cancelled: %w", c.Stdctx.Err())
			}

			if err := downloadAndVerifyFirmware(c.Stdctx, cmd); err != nil {
				return err
			}
			c.Set(BagKeyExecutionID, cmd.ExecutionID)

			topic := "events/" + uuid + "/ota_status"
			for _, s := range otaStatusProgression {
				report := otaDtos.OTAStatusRequestDTO{
					AssetUUID:   uuid,
					ExecutionID: cmd.ExecutionID,
					Status:      s.status,
					Progress:    s.progress,
				}
				body, err := json.Marshal(report)
				if err != nil {
					return fmt.Errorf("marshal ota status %q: %w", s.status, err)
				}
				if err := cli.Publish(c.Stdctx, topic, 1, false, body); err != nil {
					return fmt.Errorf("publish ota status %q on %q: %w", s.status, topic, err)
				}
				time.Sleep(150 * time.Millisecond)
			}
			return nil
		},
		Compensate: func(_ *saga.Context) error { return nil },
	}
}

// mqttClientFromBag returns the live MQTT client ConnectMqttPassword put on the bag.
func mqttClientFromBag(c *saga.Context) (*mqttclient.Client, error) {
	v, ok := c.Get(assetSteps.BagKeyMqttClient)
	if !ok {
		return nil, fmt.Errorf("mqtt ota device: mqtt client missing on bag (ConnectMqttPassword must run first)")
	}
	cli, ok := v.(*mqttclient.Client)
	if !ok {
		return nil, fmt.Errorf("mqtt ota device: bag[%s] is not *mqttclient.Client (%T)", assetSteps.BagKeyMqttClient, v)
	}
	return cli, nil
}

// otaCmdChanFromBag returns the command channel SubscribeOtaCommand armed.
func otaCmdChanFromBag(c *saga.Context) (chan downlink.OTAUpdateCommand, error) {
	v, ok := c.Get(bagKeyOtaCmdChan)
	if !ok {
		return nil, fmt.Errorf("run mqtt ota device: command channel missing (SubscribeOtaCommand must run before CreatePlan)")
	}
	ch, ok := v.(chan downlink.OTAUpdateCommand)
	if !ok {
		return nil, fmt.Errorf("run mqtt ota device: bag[%s] is not a command channel (%T)", bagKeyOtaCmdChan, v)
	}
	return ch, nil
}

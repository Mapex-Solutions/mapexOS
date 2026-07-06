package payloads

import (
	"fmt"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets_templates"

	"github.com/Mapex-Solutions/mapexGoKit/utils/zerovalue"
)

// SagaLorawanSensorTemplate returns the canonical template for a LoRaWAN sensor
// saga. Unlike the MQTT temperature template (whose script reads a JSON body),
// js-executor hands a LoRaWAN uplink to the codec as
// `{ bytes, fPort, fCnt, rxInfo }` — the decrypted FRMPayload bytes plus frame
// metadata. The template's first script IS the device codec: it turns those bytes
// into the platform's canonical StandardizedPayload shape
// (`{ eventType, eventId, data, created }`) so the engine validation passes and the
// raw event + downstream consumers (router, events) see a record they can persist.
//
// The e2e asserts the RAW ingested event (bytes + fPort/fCnt/rxInfo), not the
// decoded output, so this codec stays a deterministic pass-through of the frame.
func SagaLorawanSensorTemplate(runID string) *AssetTemplateCreateBuilder {
	const scriptConversion = `const result = {
  eventType: 'lorawan',
  eventId: String(payload.fCnt),
  data: { bytes: payload.bytes, fPort: payload.fPort, fCnt: payload.fCnt },
  created: new Date().toISOString(),
};`
	return &AssetTemplateCreateBuilder{
		spec: contracts.AssetTemplateCreate{
			Name:             fmt.Sprintf("saga-lorawan-%s", runID),
			Description:      zerovalue.Ptr("Saga LoRaWAN sensor template (codec: uplink frame → StandardizedPayload)"),
			Enabled:          true,
			AssetIDPath:      "assetUUID",
			ScriptConversion: scriptConversion,
		},
	}
}

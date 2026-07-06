package nats

import (
	"context"

	"assets/src/modules/ota/application/ports"

	downlinkContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/downlink"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// edgeDispatchAdapter implements ports.EdgeDispatchPort by publishing the
// shared DownlinkEnvelope to the STATIC downlink subjects (mqtt.downlink /
// lorawan.downlink); org + assetUUID + command travel in the payload.
type edgeDispatchAdapter struct {
	pub natsModel.Publisher
}

// NewEdgeDispatchAdapter returns an EdgeDispatchPort over the NATS publisher.
func NewEdgeDispatchAdapter(pub natsModel.Publisher) ports.EdgeDispatchPort {
	return &edgeDispatchAdapter{pub: pub}
}

func (a *edgeDispatchAdapter) Dispatch(ctx context.Context, protocol, orgID, assetUUID string, payload ports.OTACommandPayload) error {
	subject := downlinkContract.SubjectMQTTDownlink
	if protocol == "lorawan" {
		subject = downlinkContract.SubjectLoRaWANDownlink
	}
	envelope, err := downlinkContract.NewOTAUpdateEnvelope(orgID, assetUUID, payload)
	if err != nil {
		return err
	}
	return a.pub.Publish(natsModel.PublishConfig{
		Ctx:     ctx,
		Subject: subject,
		Data:    envelope,
	})
}

var _ ports.EdgeDispatchPort = (*edgeDispatchAdapter)(nil)

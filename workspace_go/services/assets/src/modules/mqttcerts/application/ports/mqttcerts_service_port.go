package ports

import (
	"context"

	dtos "assets/src/modules/mqttcerts/application/dtos"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

// MqttCertsServicePort is the public surface of the device-cert
// lifecycle. JWT-gated external handlers depend on this interface.
type MqttCertsServicePort interface {
	OnMount()
	IsCAReady() bool
	IssueCert(ctx context.Context, rc *reqCtx.RequestContext, req *dtos.IssueCertRequest) (*dtos.IssueCertResponse, error)
	RevokeCert(ctx context.Context, rc *reqCtx.RequestContext, serial string, reason string) error
	ListRevokedByAsset(ctx context.Context, rc *reqCtx.RequestContext, assetUUID string) ([]*dtos.RevokedCertResponse, error)
	HardDeleteByAssetUUID(ctx context.Context, assetUUID string) error
}

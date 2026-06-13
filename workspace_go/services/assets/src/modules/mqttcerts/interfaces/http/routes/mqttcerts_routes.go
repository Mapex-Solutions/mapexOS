package routes

import (
	dtos "assets/src/modules/mqttcerts/application/dtos"
	mqttPorts "assets/src/modules/mqttcerts/application/ports"
	"assets/src/modules/mqttcerts/interfaces/http/handlers"
	localMw "assets/src/modules/mqttcerts/interfaces/http/middlewares"

	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes mounts /api/v1/mqtt_certs/*. The router passed in MUST already
// be JWT-gated + coverage-injected by the caller (module.go composes that group).
// The CA-ready gate (RequireCAReady) and per-route permissions are applied here.
//
// Routes are registered through the swagger wrapper: each NewValidation declares
// the input contract once (used to both validate and document), the module tag is
// declared once on Wrap, and each route's summary, description, and response type
// are attached via the fluent builder.
func RegisterRoutes(router web.Router, h *handlers.MqttCertsHandler, svc mqttPorts.MqttCertsServicePort) {

	r := swagger.Wrap(router).Tag("MQTT Certificates")

	// Issue a new device cert. Returns the cert + key + CA chain once (the key is
	// never persisted server-side; the operator downloads it immediately).
	issueDto := validation.NewValidation(&dtos.IssueCertRequest{}, nil, nil)
	r.Post("/", issueDto, swagger.Expose,
		localMw.RequireCAReady(svc),
		h.IssueCert,
	).
		Summary("Issue device certificate").
		Description("Issues a new X.509 client certificate for an asset's MQTT mutual-TLS auth. Returns the cert, private key, and CA chain once — the key is never persisted. Pass force=true to replace an asset's existing current cert. Returns 409 when a current cert exists and force is not set.").
		Returns(&dtos.IssueCertResponse{})

	// Revoke a cert by serial.
	revokeDto := validation.NewValidation(nil, nil, nil)
	r.Delete("/:serial", revokeDto, swagger.Expose,
		localMw.RequireCAReady(svc),
		h.RevokeCert,
	).
		Summary("Revoke certificate").
		Description("Revokes a certificate by its serial, moving it to the revoked-certificates collection.").
		Returns(map[string]bool{})

	// List revoked certs for a single asset.
	listDto := validation.NewValidation(nil, &dtos.ListRevokedQuery{}, nil)
	r.Get("/", listDto, swagger.Expose,
		localMw.RequireCAReady(svc),
		h.ListByAsset,
	).
		Summary("List revoked certificates").
		Description("Returns the revoked certificates for a single asset, identified by the assetUUID query parameter.").
		Returns([]dtos.RevokedCertResponse{})
}

// RegisterGatewayRoutes mounts /api/v1/gateway_certs/* (issue only). Same CA-ready
// gate; reuses the shared signer/CA machinery, scoped to lorawan gateways by the
// service eligibility check.
func RegisterGatewayRoutes(router web.Router, h *handlers.MqttCertsHandler, svc mqttPorts.MqttCertsServicePort) {

	r := swagger.Wrap(router).Tag("Gateway Certificates")

	// Issue a LoRaWAN gateway cert (eligibility enforced server-side).
	issueDto := validation.NewValidation(&dtos.IssueCertRequest{}, nil, nil)
	r.Post("/", issueDto, swagger.Expose,
		localMw.RequireCAReady(svc),
		h.IssueGatewayCert,
	).
		Summary("Issue gateway certificate").
		Description("Issues an X.509 certificate for a LoRaWAN gateway asset, reusing the shared signer/CA machinery. The service rejects assets that are not eligible lorawan gateways.").
		Returns(&dtos.IssueCertResponse{})
}

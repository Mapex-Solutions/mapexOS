package di

import (
	assetPorts "assets/src/modules/assets/application/ports"
	mqttPorts "assets/src/modules/mqttcerts/application/ports"
	domainRepos "assets/src/modules/mqttcerts/domain/repositories"

	"go.uber.org/dig"
)

// MqttCertsServiceDI aggregates dependencies for the MqttCertsService.
// Every field is a port interface — never a concrete driver.
//
// AssetService is the cross-module link used to reflect cert lifecycle
// events back onto the asset entity (mqttcerts owns the cert; the
// assets module owns `asset.currentCert` + the L2 / FANOUT side
// effects). The reverse path — assets module cleaning up revoked rows
// on asset deletion — flows through MqttCertsServicePort's
// HardDeleteByAssetUUID, so the dependency stays one-way per call site.
type MqttCertsServiceDI struct {
	dig.In
	RevokedRepo      domainRepos.RevokedRepository
	MapexVaultClient mqttPorts.MapexVaultClientPort
	Signer           mqttPorts.X509SignerPort
	CAStore          mqttPorts.CAStorePort
	AssetService     assetPorts.AssetServicePort
}

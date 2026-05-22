package steps

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/mqttclient"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
)

// ConnectMqttCert opens an MQTT CONNECT against the platform's TLS
// listener (port 8883) presenting the device cert + private key
// captured by IssueCert. The broker plugin verifies the cert serial
// against asset.currentCert.serial; the username on the wire is the
// bare assetUUID so the broker plugin's cache lookup uses the same
// key as the password path.
//
// Reads (bag):
//   - BagKeyAssetUUID         string  set by CreateAsset
//   - BagKeyAssetCertPEM      []byte  set by IssueCert
//   - BagKeyAssetKeyPEM       []byte  set by IssueCert
//   - BagKeyAssetCAChainPEM   []byte  set by IssueCert
//
// Writes (bag):
//   - BagKeyMqttClient        *mqttclient.Client
//   - BagKeyMqttConnectedAt   time.Time
//
// Compensate: best-effort Disconnect.
func ConnectMqttCert() saga.Step {
	return saga.Step{
		Name: "assets/assets.ConnectMqttCert",
		Do: func(c *saga.Context) error {
			uuid := c.MustGetString(BagKeyAssetUUID)
			certPEM := mustGetBytes(c, BagKeyAssetCertPEM)
			keyPEM := mustGetBytes(c, BagKeyAssetKeyPEM)
			caPEM := mustGetBytes(c, BagKeyAssetCAChainPEM)

			tlsCfg, err := buildDeviceTLSConfig(certPEM, keyPEM, caPEM)
			if err != nil {
				return fmt.Errorf("build mTLS config for %s: %w", uuid, err)
			}

			cli, err := mqttclient.New(mqttclient.Config{
				BrokerURL: constants.MqttBrokerTLSURL,
				ClientID:  uuid,
				Username:  uuid,
				TLSConfig: tlsCfg,
			})
			if err != nil {
				return fmt.Errorf("build mqtt client: %w", err)
			}
			if err := cli.Connect(c.Stdctx); err != nil {
				return fmt.Errorf("mqtt connect (cert) for %s: %w", uuid, err)
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

// buildDeviceTLSConfig produces the *tls.Config a paho client needs
// to present an X.509 client cert. The CA chain is loaded into
// RootCAs so the client verifies the broker's server cert as well —
// both ends use the same platform PKI, so a single pool suffices.
func buildDeviceTLSConfig(certPEM, keyPEM, caPEM []byte) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("parse device cert/key: %w", err)
	}
	pool := x509.NewCertPool()
	if len(caPEM) > 0 && !pool.AppendCertsFromPEM(caPEM) {
		return nil, fmt.Errorf("append CA chain to pool: no PEM blocks parsed")
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// mustGetBytes pulls a []byte from the bag or panics — the saga's
// MustGetString equivalent, kept here so the cert step has the same
// fail-loud signature as the rest of the file's MustGet* calls.
func mustGetBytes(c *saga.Context, key string) []byte {
	v, ok := c.Get(key)
	if !ok {
		panic(fmt.Sprintf("saga bag missing %q", key))
	}
	b, ok := v.([]byte)
	if !ok {
		panic(fmt.Sprintf("saga bag[%q] is not []byte (%T)", key, v))
	}
	return b
}

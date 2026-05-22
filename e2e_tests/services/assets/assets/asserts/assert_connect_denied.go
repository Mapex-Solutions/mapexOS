package asserts

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/mqttclient"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
)

// AssertConnectDeniedPassword opens a fresh MQTT CONNECT with the
// same credentials the saga used earlier and expects the broker to
// reject it. Used right after DeleteAsset to prove the FANOUT-driven
// invalidation reached the broker plugin's L1 (Pebble) cache — if
// the entry hadn't been evicted the deleted asset would still
// authenticate from cache. A short connect deadline keeps the assert
// snappy when the broker correctly denies on the first round trip.
//
// Reads (bag):
//   - assetSteps.BagKeyAssetUUID         string  set by CreateAsset
//   - assetSteps.BagKeyAssetMqttPassword string  set by CreateAsset
func AssertConnectDeniedPassword() saga.Assert {
	return saga.Assert{
		Name: "assets/assets.AssertConnectDeniedPassword",
		Check: func(c *saga.Context) error {
			uuid := c.MustGetString(assetSteps.BagKeyAssetUUID)
			pwd := c.MustGetString(assetSteps.BagKeyAssetMqttPassword)
			return attemptConnectAndExpectDeny(c, uuid, mqttclient.Config{
				BrokerURL:      constants.MqttBrokerURL,
				ClientID:       uuid + "-deny-probe",
				Username:       uuid,
				Password:       pwd,
				ConnectTimeout: 5 * time.Second,
			})
		},
	}
}

// AssertConnectDeniedCert opens a fresh mTLS CONNECT with the cert
// captured by IssueCert and expects the broker to reject it. Used
// after DeleteAsset to prove the FANOUT invalidation reached the
// broker's L1 — without it the cached entry would still allow the
// CONNECT even though asset.currentCert is gone from Mongo.
//
// Reads (bag):
//   - assetSteps.BagKeyAssetUUID        string  set by CreateAsset
//   - assetSteps.BagKeyAssetCertPEM     []byte  set by IssueCert
//   - assetSteps.BagKeyAssetKeyPEM      []byte  set by IssueCert
//   - assetSteps.BagKeyAssetCAChainPEM  []byte  set by IssueCert
func AssertConnectDeniedCert() saga.Assert {
	return saga.Assert{
		Name: "assets/assets.AssertConnectDeniedCert",
		Check: func(c *saga.Context) error {
			uuid := c.MustGetString(assetSteps.BagKeyAssetUUID)
			certPEM, keyPEM, caPEM, err := readCertBundle(c)
			if err != nil {
				return fmt.Errorf("deny-probe (cert): %w", err)
			}
			tlsCfg, err := buildClientTLSConfig(certPEM, keyPEM, caPEM)
			if err != nil {
				return fmt.Errorf("deny-probe (cert): build TLS: %w", err)
			}
			return attemptConnectAndExpectDeny(c, uuid, mqttclient.Config{
				BrokerURL:      constants.MqttBrokerTLSURL,
				ClientID:       uuid + "-deny-probe",
				Username:       uuid,
				TLSConfig:      tlsCfg,
				ConnectTimeout: 5 * time.Second,
			})
		},
	}
}

// readCertBundle pulls the PEM triple from the bag with typed errors
// — separate from the bagBytes path so the assert can surface "no
// bundle in bag" without conflating with TLS parsing failures.
func readCertBundle(c *saga.Context) ([]byte, []byte, []byte, error) {
	cert, ok := c.Get(assetSteps.BagKeyAssetCertPEM)
	if !ok {
		return nil, nil, nil, fmt.Errorf("cert PEM missing on bag")
	}
	key, ok := c.Get(assetSteps.BagKeyAssetKeyPEM)
	if !ok {
		return nil, nil, nil, fmt.Errorf("key PEM missing on bag")
	}
	ca, ok := c.Get(assetSteps.BagKeyAssetCAChainPEM)
	if !ok {
		return nil, nil, nil, fmt.Errorf("CA chain PEM missing on bag")
	}
	return cert.([]byte), key.([]byte), ca.([]byte), nil
}

// buildClientTLSConfig mirrors ConnectMqttCert's helper but lives in
// the asserts package so the deny-probe stays self-contained.
func buildClientTLSConfig(certPEM, keyPEM, caPEM []byte) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("parse cert/key: %w", err)
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

// attemptConnectAndExpectDeny opens a CONNECT with the supplied
// config and returns nil only when the broker refuses. A successful
// CONNECT means the FANOUT invalidation or asset-delete path failed
// upstream — the saga reports the asset uuid so the operator can
// chase it in the broker plugin logs.
func attemptConnectAndExpectDeny(c *saga.Context, uuid string, cfg mqttclient.Config) error {
	cli, err := mqttclient.New(cfg)
	if err != nil {
		return fmt.Errorf("deny-probe: build mqtt client: %w", err)
	}
	connectCtx, cancel := context.WithTimeout(c.Stdctx, cfg.ConnectTimeout)
	defer cancel()

	err = cli.Connect(connectCtx)
	if err == nil {
		cli.Disconnect(0)
		return fmt.Errorf("deny-probe: expected broker to reject CONNECT for deleted asset %s, but it succeeded", uuid)
	}
	return nil
}

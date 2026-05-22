package services

import (
	"context"
	"fmt"
	"time"

	appConsts "assets/src/modules/mqttcerts/application/constants"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// tryBootstrapSync makes ONE attempt with a short timeout. Returns true
// on success (CA loaded into the RAM store), false otherwise.
func (s *MqttCertsService) tryBootstrapSync() bool {
	ctx, cancel := context.WithTimeout(context.Background(), appConsts.CABootstrapInitialTimeout)
	defer cancel()
	ca, err := s.deps.MapexVaultClient.FetchIntermediateCABundle(ctx)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:MqttCerts] bootstrap sync failed: %v", err))
		return false
	}
	s.deps.CAStore.Set(ca)
	return true
}

// markCAReady is the logging point; the flag flip is implicit via CAStore.Set.
func (s *MqttCertsService) markCAReady() {
	logger.Info("[SERVICE:MqttCerts] CA loaded from mapexVault — ready (caReady=true)")
}

// spawnBootstrapGoroutine retries until the CA loads. Exponential
// backoff 1s -> 30s, no max attempts. Exits on parent ctx cancel.
func (s *MqttCertsService) spawnBootstrapGoroutine(parentCtx context.Context) {
	backoff := appConsts.CABootstrapBackoffMin
	for {
		select {
		case <-parentCtx.Done():
			return
		case <-time.After(backoff):
		}
		attemptCtx, cancel := context.WithTimeout(parentCtx, appConsts.CABootstrapInitialTimeout)
		ca, err := s.deps.MapexVaultClient.FetchIntermediateCABundle(attemptCtx)
		cancel()
		if err == nil {
			s.deps.CAStore.Set(ca)
			s.markCAReady()
			return
		}
		logger.Warn(fmt.Sprintf("[SERVICE:MqttCerts] bootstrap retry failed: %v (next in %s)", err, backoff))
		backoff *= 2
		if backoff > appConsts.CABootstrapBackoffMax {
			backoff = appConsts.CABootstrapBackoffMax
		}
	}
}

package services

import (
	"context"
	"fmt"
	"time"

	appConstants "assets/src/modules/assets/application/constants"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// tryLoadLorawanKEK makes ONE attempt with a short timeout to fetch the KEK and
// load it into the cipher. Returns true on success.
func (s *AssetService) tryLoadLorawanKEK() bool {
	ctx, cancel := context.WithTimeout(context.Background(), appConstants.LorawanKEKBootstrapInitialTimeout)
	defer cancel()
	return s.fetchAndStoreLorawanKEK(ctx)
}

// retryLoadLorawanKEK retries until the KEK loads. Exponential backoff
// 1s -> 30s, no max attempts. Exits on parent ctx cancel.
func (s *AssetService) retryLoadLorawanKEK(parentCtx context.Context) {
	logger.Warn("[SERVICE:Assets] LoRaWAN KEK: initial fetch failed, retrying in background")
	backoff := appConstants.LorawanKEKBootstrapBackoffMin
	for {
		select {
		case <-parentCtx.Done():
			return
		case <-time.After(backoff):
		}
		attemptCtx, cancel := context.WithTimeout(parentCtx, appConstants.LorawanKEKBootstrapInitialTimeout)
		ok := s.fetchAndStoreLorawanKEK(attemptCtx)
		cancel()
		if ok {
			return
		}
		backoff *= 2
		if backoff > appConstants.LorawanKEKBootstrapBackoffMax {
			backoff = appConstants.LorawanKEKBootstrapBackoffMax
		}
	}
}

// fetchAndStoreLorawanKEK fetches the KEK from mapexVault and loads it into the
// cipher, logging the outcome. Returns true only when the cipher is ready.
func (s *AssetService) fetchAndStoreLorawanKEK(ctx context.Context) bool {
	kekHex, err := s.deps.LorawanKEKClient.FetchKEK(ctx, appConstants.LorawanKEKContext)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Assets] LoRaWAN KEK fetch failed err=%v", err))
		return false
	}
	if err := s.deps.LorawanKEKCipher.SetKey(kekHex); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Assets] LoRaWAN KEK load failed err=%v", err))
		return false
	}
	logger.Info("[SERVICE:Assets] LoRaWAN KEK loaded from mapexVault, ready")
	return true
}

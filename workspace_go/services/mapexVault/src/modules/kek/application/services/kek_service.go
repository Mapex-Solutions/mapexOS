package services

import (
	"context"
	"errors"
	"fmt"

	kekDi "mapexVault/src/modules/kek/application/di"
	kekDtos "mapexVault/src/modules/kek/application/dtos"
	kekPorts "mapexVault/src/modules/kek/application/ports"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// Compile-time port check.
var _ kekPorts.KekServicePort = (*KekService)(nil)

// ErrKEKNotSeeded is returned when a context is requested before the Mongo seed
// has populated encryptionKeys. The HTTP layer maps it to 503 so the caller's
// boot retry loop keeps trying until the mongodb-init container loads the seed.
var ErrKEKNotSeeded = errors.New("kek not seeded")

// New constructs the KEK service. Returns the port interface; the concrete
// struct stays unexposed.
func New(deps kekDi.KekServiceDI) kekPorts.KekServicePort {
	return &KekService{deps: deps}
}

// GetKEK loads the KEK for a context and returns it decrypted. Used by the
// assets MS and the LNS at boot to seed their in-RAM envelope service.
func (s *KekService) GetKEK(ctx context.Context, keyContext string) (*kekDtos.KEKResponse, error) {
	logger.Info("[SERVICE:Kek] GetKEK: loading entity from Mongo context=" + keyContext)
	entity, err := s.loadByContext(ctx, keyContext)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Kek] GetKEK: load failed context=%s err=%v", keyContext, err))
		return nil, fmt.Errorf("load kek: %w", err)
	}
	plaintext, err := s.decryptKEK(entity)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Kek] GetKEK: decrypt failed context=%s err=%v", keyContext, err))
		return nil, fmt.Errorf("decrypt kek: %w", err)
	}
	logger.Info("[SERVICE:Kek] GetKEK: decrypt ok context=" + keyContext)
	return &kekDtos.KEKResponse{Context: keyContext, Kek: string(plaintext)}, nil
}

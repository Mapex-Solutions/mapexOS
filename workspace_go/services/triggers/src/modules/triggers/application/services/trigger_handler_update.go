package services

import (
	"time"

	"triggers/src/modules/triggers/application/dtos"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// buildTriggerUpdatePayload composes the $set map for an update, stamping
// the updated timestamp at the same time. Surfaces the canonical 400
// contract error when the DTO cannot be mapped.
func (s *TriggerService) buildTriggerUpdatePayload(dto *dtos.UpdateTriggerDto) (map[string]interface{}, error) {
	payload, err := mapper.DtoToMap(dto)
	if err != nil {
		return nil, &customErrors.ServerCustomError{Code: status.BAD_REQUEST, Errors: []string{err.Error()}}
	}
	payload["updated"] = time.Now()
	return payload, nil
}

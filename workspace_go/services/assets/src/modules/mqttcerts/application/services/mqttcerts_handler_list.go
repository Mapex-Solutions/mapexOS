package services

import (
	dtos "assets/src/modules/mqttcerts/application/dtos"
	"assets/src/modules/mqttcerts/domain/entities"

	mapper "github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

func (s *MqttCertsService) mapToResponseDTOs(rows []*entities.RevokedCertificate) []*dtos.RevokedCertResponse {
	out := make([]*dtos.RevokedCertResponse, 0, len(rows))
	for _, r := range rows {
		dto, _ := mapper.EntityToDto[entities.RevokedCertificate, dtos.RevokedCertResponse](r)
		out = append(out, dto)
	}
	return out
}

package services

import (
	ctx "context"

	"assets/src/modules/assets/application/dtos"
	"assets/src/modules/assets/domain/entities"

	assetsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// buildSimpleResponse maps an entity to the API response shape with no
// enrichment. Used by paths that don't need template/route-group data
// (e.g. the auth lookup by MQTT username).
func (s *AssetService) buildSimpleResponse(entity *entities.Asset) *dtos.AssetResponse {
	resp, _ := mapper.EntityToDto[entities.Asset, dtos.AssetResponse](entity)
	if !entity.ID.IsZero() {
		resp.ID = &entity.ID
	}
	enabled := entity.Enabled
	resp.Enabled = &enabled
	debugEnabled := entity.DebugEnabled
	resp.DebugEnabled = &debugEnabled
	if entity.CurrentCert != nil {
		resp.CurrentCert = &assetsContract.AssetCertificate{
			Serial:      entity.CurrentCert.Serial,
			Fingerprint: entity.CurrentCert.Fingerprint,
			SubjectCN:   entity.CurrentCert.SubjectCN,
			IssuedAt:    entity.CurrentCert.IssuedAt,
			ExpiresAt:   entity.CurrentCert.ExpiresAt,
		}
	}
	return resp
}

// buildCreateResponse shapes the response after a successful insert.
// The MQTT password the operator supplied is intentionally NOT echoed —
// the response carries only the persisted public identity (clientId +
// username); the plaintext password lives only in the request body and
// is the operator's copy of record.
func (s *AssetService) buildCreateResponse(entity *entities.Asset) *dtos.AssetResponse {
	base := s.buildSimpleResponse(entity)
	if !entity.AssetTemplateID.IsZero() {
		templateId := entity.AssetTemplateID.Hex()
		base.AssetTemplateID = &templateId
	}
	return base
}

// buildUpdateResponse shapes the response after a successful patch.
func (s *AssetService) buildUpdateResponse(entity *entities.Asset) *dtos.AssetResponse {
	return s.buildSimpleResponse(entity)
}

// buildGetByIdResponse shapes the response for the single-asset read.
// Includes the Mongo-persisted health flip; Redis enrichment runs after.
func (s *AssetService) buildGetByIdResponse(entity *entities.Asset) *dtos.AssetResponse {
	resp := s.buildSimpleResponse(entity)
	if !entity.AssetTemplateID.IsZero() {
		templateId := entity.AssetTemplateID.Hex()
		resp.AssetTemplateID = &templateId
	}
	resp.HealthStatusChangedAt = entity.HealthStatusChangedAt
	resp.HealthMonitor = convertHealthMonitor(entity.HealthMonitor)
	return resp
}

// enrichWithTemplateClassification overlays manufacturer/model/category
// metadata from the template on the response. Non-critical: a missing
// template (template was deleted) leaves the fields nil and the response
// is still returned.
func (s *AssetService) enrichWithTemplateClassification(c ctx.Context, entity *entities.Asset, resp *dtos.AssetResponse) {
	if entity.AssetTemplateID.IsZero() {
		return
	}
	templateIdStr := entity.AssetTemplateID.Hex()
	template, err := s.deps.AssetTemplateRepo.FindById(c, &templateIdStr)
	if err != nil || template == nil || template.ID.IsZero() {
		return
	}

	resp.ManufacturerName = template.ManufacturerName
	resp.ModelName = template.ModelName
	resp.CategoryName = template.CategoryName
	resp.Version = template.Version

	templateName := template.Name
	resp.AssetTemplateName = &templateName
	if template.AssetIDPath != "" {
		resp.AssetIdPath = &template.AssetIDPath
	}
	if template.ManufacturerId != nil {
		v := template.ManufacturerId.Hex()
		resp.ManufacturerId = &v
	}
	if template.ModelId != nil {
		v := template.ModelId.Hex()
		resp.ModelId = &v
	}
	if template.CategoryId != nil {
		v := template.CategoryId.Hex()
		resp.CategoryId = &v
	}
}

// enrichWithRouteGroupNames resolves human-readable group names by
// calling the Router service. Non-critical — empty/error leaves the field
// nil.
func (s *AssetService) enrichWithRouteGroupNames(c ctx.Context, entity *entities.Asset, resp *dtos.AssetResponse) {
	if len(entity.RouteGroupIds) == 0 {
		return
	}
	names, err := s.deps.RouteGroupPort.GetNamesByIds(c, entity.RouteGroupIds)
	if err == nil && len(names) > 0 {
		resp.RouteGroupNames = &names
	}
}

// mapListEntitiesToDtos converts each AssetWithTemplate to a response
// DTO, fixes the bool->*bool conversions copier can't handle, and
// normalizes HealthStatus + HealthStatusChangedAt onto each item.
func (s *AssetService) mapListEntitiesToDtos(items []entities.AssetWithTemplate) []dtos.AssetResponse {
	out := make([]dtos.AssetResponse, len(items))
	for i := range items {
		entity := items[i]
		dto, _ := mapper.EntityToDto[entities.AssetWithTemplate, dtos.AssetResponse](&entity)

		enabled := entity.Enabled
		dto.Enabled = &enabled
		debugEnabled := entity.DebugEnabled
		dto.DebugEnabled = &debugEnabled

		dto.HealthStatusChangedAt = entity.HealthStatusChangedAt
		dto.HealthMonitor = convertHealthMonitor(entity.HealthMonitor)

		if entity.HealthStatus != "" {
			status := entity.HealthStatus
			dto.HealthStatus = &status
		} else {
			dto.HealthStatus = nil
		}

		out[i] = *dto
	}
	return out
}

// buildReadModel produces the denormalized read-model for the
// internal cache-fallback endpoint. Does the manual ID/Description/
// Protocol/HealthMonitor conversions copier can't handle reliably.
func (s *AssetService) buildReadModel(asset *entities.Asset, templateOrgId string) *assetsContract.AssetReadModel {
	rm, _ := mapper.EntityToDto[entities.Asset, assetsContract.AssetReadModel](asset)
	rm.ID = asset.ID.Hex()
	rm.UUID = asset.AssetUUID
	rm.OrgId = asset.OrgID.Hex()
	rm.AssetTemplateID = asset.AssetTemplateID.Hex()
	rm.AssetTemplateOrgID = templateOrgId
	if asset.Description != nil {
		rm.Description = *asset.Description
	}
	if asset.Protocol.Type != "" {
		rm.Protocol = &assetsContract.ProtocolType{Type: asset.Protocol.Type}
		if asset.Protocol.Mqtt != nil {
			rm.Protocol.Mqtt = &assetsContract.MqttConfig{
				ClientId:     asset.Protocol.Mqtt.ClientId,
				Username:     asset.Protocol.Mqtt.Username,
				AuthType:     asset.Protocol.Mqtt.AuthType,
				PasswordHash: asset.Protocol.Mqtt.PasswordHash,
			}
		}
	}
	if asset.CurrentCert != nil {
		rm.CurrentCert = &assetsContract.AssetCertificate{
			Serial:      asset.CurrentCert.Serial,
			Fingerprint: asset.CurrentCert.Fingerprint,
			SubjectCN:   asset.CurrentCert.SubjectCN,
			IssuedAt:    asset.CurrentCert.IssuedAt,
			ExpiresAt:   asset.CurrentCert.ExpiresAt,
		}
	}
	rm.HealthMonitor = convertHealthMonitor(asset.HealthMonitor)
	return rm
}

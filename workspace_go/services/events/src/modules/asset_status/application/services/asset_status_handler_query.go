package services

import (
	"events/src/modules/asset_status/application/dtos"
	"events/src/modules/asset_status/domain/entities"

	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
)

// buildAssetStatusFilters translates the org-filter map and optional query
// filters into the []chModel.Filter slice the repository expects.
func buildAssetStatusFilters(orgFilter map[string]interface{}, query *dtos.AssetConnectivityHistoryQuery) []chModel.Filter {
	filters := []chModel.Filter{}

	if orgId, ok := orgFilter["orgId"]; ok && orgId != nil {
		filters = append(filters, chModel.Filter{
			Field:    "org_id",
			Operator: chModel.OpEqual,
			Value:    orgId,
		})
	}

	if pathKey, ok := orgFilter["pathKey"]; ok && pathKey != nil {
		if s, isStr := pathKey.(string); isStr {
			filters = append(filters, chModel.Filter{
				Field:    "path_key",
				Operator: chModel.OpEqual,
				Value:    s,
			})
		} else if m, isMap := pathKey.(map[string]interface{}); isMap {
			if gte, ok := m["$gte"]; ok {
				if lt, ok := m["$lt"]; ok {
					filters = append(filters, chModel.Filter{
						Field:    "path_key",
						Operator: chModel.OpBetween,
						Value:    gte,
						EndValue: lt,
					})
				}
			}
		}
	}

	if query.EventType != nil && *query.EventType != "" {
		filters = append(filters, chModel.Filter{
			Field:    "event_type",
			Operator: chModel.OpEqual,
			Value:    *query.EventType,
		})
	}

	if query.AssetUUID != nil && *query.AssetUUID != "" {
		filters = append(filters, chModel.Filter{
			Field:    "asset_uuid",
			Operator: chModel.OpEqual,
			Value:    *query.AssetUUID,
		})
	}

	if query.From != nil && query.To != nil {
		filters = append(filters, chModel.Filter{
			Field:    "created",
			Operator: chModel.OpBetween,
			Value:    query.From.UTC(),
			EndValue: query.To.UTC(),
		})
	} else {
		if query.From != nil {
			filters = append(filters, chModel.Filter{
				Field:    "created",
				Operator: chModel.OpGreaterEqual,
				Value:    query.From.UTC(),
			})
		}
		if query.To != nil {
			filters = append(filters, chModel.Filter{
				Field:    "created",
				Operator: chModel.OpLessEqual,
				Value:    query.To.UTC(),
			})
		}
	}

	return filters
}

// buildAssetStatusCursorOpts assembles a TimeCursorOpts from the embedded
// CursorQueryDTO getters on the request payload.
func buildAssetStatusCursorOpts(query *dtos.AssetConnectivityHistoryQuery) *chModel.TimeCursorOpts {
	cursorOpts := &chModel.TimeCursorOpts{
		Limit:     int64(query.GetLimit()),
		Direction: query.GetDirection(),
		SortAsc:   query.GetSortAsc(),
	}
	if query.GetCursor() != nil {
		cursorOpts.Cursor = *query.GetCursor()
	}
	return cursorOpts
}

// buildAssetStatusResponse maps the cursor result into the public DTO,
// surfacing next/previous cursors only when the repository populated them.
func buildAssetStatusResponse(result *chModel.TimeCursorResult[entities.AssetStatusEvent]) *dtos.AssetConnectivityCursorResult {
	items := make([]dtos.AssetConnectivityEvent, 0, len(result.Items))
	for _, entity := range result.Items {
		items = append(items, dtos.AssetConnectivityEvent{
			Created:          entity.Created,
			OrgId:            entity.OrgId,
			PathKey:          entity.PathKey,
			AssetUUID:        entity.AssetUUID,
			AssetName:        entity.AssetName,
			EventId:          entity.EventId,
			EventType:        entity.EventType,
			LastSeenAt:       entity.LastSeenAt,
			ThresholdMinutes: entity.ThresholdMinutes,
			MissCount:        entity.MissCount,
		})
	}

	response := &dtos.AssetConnectivityCursorResult{
		Items:       items,
		HasNext:     result.HasNext,
		HasPrevious: result.HasPrevious,
	}
	if !result.NextCursor.IsZero() {
		response.NextCursor = &result.NextCursor
	}
	if !result.PrevCursor.IsZero() {
		response.PrevCursor = &result.PrevCursor
	}
	return response
}

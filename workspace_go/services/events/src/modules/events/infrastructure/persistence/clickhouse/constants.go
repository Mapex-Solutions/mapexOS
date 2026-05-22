package clickhouseRepo

import (
	dtos "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
)

// evaOperatorMap maps EvaFilterOperator to ClickHouse FilterOperator.
var evaOperatorMap = map[dtos.EvaFilterOperator]chModel.FilterOperator{
	dtos.EvaOpEqual:        chModel.OpEqual,
	dtos.EvaOpNotEqual:     chModel.OpNotEqual,
	dtos.EvaOpGreater:      chModel.OpGreater,
	dtos.EvaOpGreaterEqual: chModel.OpGreaterEqual,
	dtos.EvaOpLess:         chModel.OpLess,
	dtos.EvaOpLessEqual:    chModel.OpLessEqual,
	dtos.EvaOpBetween:      chModel.OpBetween,
	dtos.EvaOpLike:         chModel.OpLike,
}

// evaBucketColumn maps EVA bucket name to ClickHouse MAP column name.
var evaBucketColumn = map[string]string{
	"number": "eva_number",
	"string": "eva_string",
	"bool":   "eva_bool",
	"date":   "eva_date",
}

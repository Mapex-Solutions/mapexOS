package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/http_gateway/datasources"
)

type (
	DataSourceCreateDTO = v1.DataSourceCreate
	DataSourceUpdateDTO = v1.DataSourceUpdate
	DataSourceIdDto     = v1.DataSourceId
	DataSourceResponse  = v1.DataSourceResponse
	DataSourceQueryDTO  = v1.DataSourceQuery
)

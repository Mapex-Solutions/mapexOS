package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// connAdapter wraps a raw clickhouse-go driver.Conn and exposes only the
// subset of operations declared by ports.ClickHouseConnPort. This keeps the
// concrete driver type out of the application layer.
type connAdapter struct {
	conn driver.Conn
}

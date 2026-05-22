package clickhouse

import (
	ctx "context"

	"events/src/modules/retention/application/ports"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// Compile-time check to ensure connAdapter implements ClickHouseConnPort.
var _ ports.ClickHouseConnPort = (*connAdapter)(nil)

// NewConnAdapter builds a ClickHouseConnPort backed by a raw driver.Conn.
//
// Parameters:
//   - conn: The raw ClickHouse driver connection.
//
// Returns:
//   - ports.ClickHouseConnPort: Port-scoped wrapper around the connection.
func NewConnAdapter(conn driver.Conn) ports.ClickHouseConnPort {
	return &connAdapter{conn: conn}
}

// Exec forwards the statement to the underlying driver.Conn.
//
// Parameters:
//   - c: Context controlling cancellation and timeouts.
//   - query: The SQL statement to execute.
//   - args: Optional parameter bindings.
//
// Returns:
//   - error: Propagated verbatim from the underlying driver.
func (a *connAdapter) Exec(c ctx.Context, query string, args ...any) error {
	return a.conn.Exec(c, query, args...)
}

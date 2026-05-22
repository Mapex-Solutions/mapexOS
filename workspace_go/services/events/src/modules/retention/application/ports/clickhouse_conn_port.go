package ports

import (
	ctx "context"
)

// ClickHouseConnPort is the minimal surface the retention service needs from
// the raw ClickHouse connection. Only methods actually exercised by the
// application layer are declared here so the concrete driver type never
// leaks past the infrastructure boundary.
//
// Currently limited to `Exec` because retention applies TTL changes via
// `ALTER TABLE ... MODIFY TTL` — it does not Query or stream rows.
type ClickHouseConnPort interface {
	// Exec runs a statement without returning rows.
	//
	// Parameters:
	//   - ctx: Context controlling cancellation and timeouts.
	//   - query: The SQL statement to execute (e.g., ALTER TABLE ... MODIFY TTL).
	//   - args: Optional parameter bindings for the statement.
	//
	// Returns:
	//   - error: Propagated verbatim from the underlying ClickHouse driver.
	Exec(ctx ctx.Context, query string, args ...any) error
}

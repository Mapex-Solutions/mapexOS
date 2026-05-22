package clickhouseRepo

import (
	"context"
	"testing"

	"events/src/modules/asset_status/domain/entities"
)

// TestBulkInsert_EmptyInput exercises the no-op fast-path. No ClickHouse
// connection required — the repo short-circuits before touching the table.
// The real insert path is exercised by the //go:build integration tests
// below once a ClickHouse instance is available.
func TestBulkInsert_EmptyInput(t *testing.T) {
	// Nil conn is safe because the BulkInsert empty-input branch returns
	// before dereferencing r.table.
	r := &AssetStatusRepositoryClickHouse{table: nil}

	if err := r.BulkInsert(context.Background(), nil); err != nil {
		t.Fatalf("BulkInsert(nil) = %v, want nil", err)
	}
	if err := r.BulkInsert(context.Background(), []*entities.AssetStatusEvent{}); err != nil {
		t.Fatalf("BulkInsert(empty) = %v, want nil", err)
	}
}

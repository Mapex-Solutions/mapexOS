package clickhouseRepo

import (
	"context"
	"fmt"

	"events/src/modules/events/domain/entities"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/**
 * Legacy Event Repository Methods (events table)
 */

// Save stores a single event in ClickHouse.
//
// This method uses ClickHouse's async insert optimization for better performance.
// ClickHouse will automatically batch inserts internally.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - event: The event to store
//
// Returns:
//   - error: If the insert fails
func (r *EventRepositoryClickHouse) Save(ctx context.Context, event *entities.Event) error {
	query := `
		INSERT INTO events (
			created,
			thread_id,
			asset_id,
			asset_template_id,
			org_id,
			path_key,
			event_type,
			source,
			payload,
			metadata,
			eva_number,
			eva_string,
			eva_bool,
			eva_date,
			retention_days
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	if err := r.conn.Exec(ctx, query,
		event.Created,
		event.ThreadId,
		event.AssetId,
		event.AssetTemplateId,
		event.OrgId,
		event.PathKey,
		event.EventType,
		event.Source,
		event.Payload,
		event.Metadata,
		event.EvaNumber,
		event.EvaString,
		event.EvaBool,
		event.EvaDate,
		event.RetentionDays,
	); err != nil {
		logger.Error(err, "[REPO:Event] Failed to insert event")
		return fmt.Errorf("failed to insert event: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:Event] Event saved: assetId=%s, type=%s", event.AssetId, event.EventType))
	return nil
}

// SaveBatch stores multiple events in ClickHouse efficiently.
//
// This method uses ClickHouse's batch insert API for optimal performance
// when processing multiple events at once (e.g., from NATS batch consumer).
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - events: Slice of events to store
//
// Returns:
//   - error: If the batch insert fails
func (r *EventRepositoryClickHouse) SaveBatch(ctx context.Context, events []*entities.Event) error {
	if len(events) == 0 {
		return nil
	}

	batch, err := r.conn.PrepareBatch(ctx, `
		INSERT INTO events (
			created,
			thread_id,
			asset_id,
			asset_template_id,
			org_id,
			path_key,
			event_type,
			source,
			payload,
			metadata,
			eva_number,
			eva_string,
			eva_bool,
			eva_date,
			retention_days
		)
	`)

	if err != nil {
		logger.Error(err, "[REPO:Event] Failed to prepare batch")
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	// Append all events to the batch
	for _, event := range events {
		if err := batch.Append(
			event.Created,
			event.ThreadId,
			event.AssetId,
			event.AssetTemplateId,
			event.OrgId,
			event.PathKey,
			event.EventType,
			event.Source,
			event.Payload,
			event.Metadata,
			event.EvaNumber,
			event.EvaString,
			event.EvaBool,
			event.EvaDate,
			event.RetentionDays,
		); err != nil {
			logger.Error(err, "[REPO:Event] Failed to append event to batch")
			return fmt.Errorf("failed to append event to batch: %w", err)
		}
	}

	// Send the batch to ClickHouse
	if err := batch.Send(); err != nil {
		logger.Error(err, "[REPO:Event] Failed to send batch")
		return fmt.Errorf("failed to send batch: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:Event] Batch saved: %d events", len(events)))
	return nil
}

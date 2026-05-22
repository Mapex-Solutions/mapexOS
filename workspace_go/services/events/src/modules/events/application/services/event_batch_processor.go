package services

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// processBatchParallel executes Phase 1 (parallel parse/validate/map) using a bounded worker pool.
//
// Worker pool size: min(runtime.NumCPU()*2, len(messages)).
// Each goroutine writes to its own results[idx] slot — no contention.
//
// The processFunc receives the message index and Message pointer, and must return
// a messageResult[T] with either action="reject" (invalid) or action="pending" (valid entity).
func processBatchParallel[T any](
	messages []*natsModel.Message,
	processFunc func(idx int, msg *natsModel.Message) messageResult[T],
) []messageResult[T] {
	workers := runtime.NumCPU() * 2
	if workers > len(messages) {
		workers = len(messages)
	}

	results := make([]messageResult[T], len(messages))
	work := make(chan int, len(messages))
	for i := range messages {
		work <- i
	}
	close(work)

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range work {
				results[idx] = processFunc(idx, messages[idx])
			}
		}()
	}
	wg.Wait()

	return results
}

// collectValidEntities extracts valid entities and their corresponding messages
// from Phase 1 results for the Phase 2 bulk insert.
func collectValidEntities[T any](results []messageResult[T]) ([]*T, []*natsModel.Message) {
	entities := make([]*T, 0, len(results))
	validMsgs := make([]*natsModel.Message, 0, len(results))
	for i := range results {
		if results[i].action == "pending" && results[i].entity != nil {
			entities = append(entities, results[i].entity)
			validMsgs = append(validMsgs, results[i].msg)
		}
	}
	return entities, validMsgs
}

// orchestrateBatch executes the full three-phase pipeline used by every
// retry-aware events consumer (raw, jsexec, router, businessrule, trigger,
// workflow, store):
//   - Phase 1: parallel parse + validate; reject failures into the DLQ.
//   - Phase 2: bulk insert via the caller-provided closure.
//   - Phase 3: Ack on success or Nack with backoff on insert failure.
//
// Centralizing the loop keeps every consumer's metrics + lifecycle calls
// identical — the only per-consumer difference is which entity is parsed
// and which insert method is dispatched.
func orchestrateBatch[T any](
	s *EventService,
	consumer string,
	tableName string,
	messages []*natsModel.Message,
	parseFn func(idx int, msg *natsModel.Message) messageResult[T],
	insertFn func(ctx context.Context, entities []*T) error,
) error {
	if len(messages) == 0 {
		return nil
	}

	batchStart := time.Now()
	logger.Info(fmt.Sprintf("[SERVICE:Event] Processing %s event batch: %d messages", consumer, len(messages)))
	s.deps.Metrics.EventsBatchSize.WithLabelValues(consumer).Observe(float64(len(messages)))

	results := processBatchParallel(messages, parseFn)
	rejectInvalidResults(s, consumer, results)

	entities, validMessages := collectValidEntities(results)
	if len(entities) == 0 {
		logger.Warn(fmt.Sprintf("[SERVICE:Event] No valid %s events in batch after validation", consumer))
		s.deps.Metrics.EventProcessingDuration.WithLabelValues(consumer).Observe(time.Since(batchStart).Seconds())
		return nil
	}

	insertErr := runBulkInsert(s, tableName, entities, insertFn)
	finalizeBatchAck(s, consumer, validMessages, insertErr)
	s.deps.Metrics.EventProcessingDuration.WithLabelValues(consumer).Observe(time.Since(batchStart).Seconds())
	if insertErr == nil {
		logger.Info(fmt.Sprintf("[SERVICE:Event] Successfully saved %s event batch: %d events", consumer, len(entities)))
	}
	return nil
}

// orchestrateDLQBatch is the DLQ variant: never Nacks (no retry loop) and
// emits "ack_skip" results immediately so they never bubble into the
// bulk insert. Mirrors the layout of orchestrateBatch.
func orchestrateDLQBatch[T any](
	s *EventService,
	consumer string,
	tableName string,
	messages []*natsModel.Message,
	parseFn func(idx int, msg *natsModel.Message) messageResult[T],
	insertFn func(ctx context.Context, entities []*T) error,
) error {
	if len(messages) == 0 {
		return nil
	}

	batchStart := time.Now()
	logger.Info(fmt.Sprintf("[SERVICE:Event] Processing %s event batch: %d messages", consumer, len(messages)))
	s.deps.Metrics.EventsBatchSize.WithLabelValues(consumer).Observe(float64(len(messages)))

	results := processBatchParallel(messages, parseFn)
	ackSkipParseFailures(s, consumer, results)

	entities, _ := collectValidEntities(results)
	if len(entities) == 0 {
		logger.Warn(fmt.Sprintf("[SERVICE:Event] No valid %s events in batch after parsing", consumer))
		s.deps.Metrics.EventProcessingDuration.WithLabelValues(consumer).Observe(time.Since(batchStart).Seconds())
		return nil
	}

	insertErr := runBulkInsert(s, tableName, entities, insertFn)
	finalizeDLQBatchAck(s, consumer, results, insertErr)
	s.deps.Metrics.EventProcessingDuration.WithLabelValues(consumer).Observe(time.Since(batchStart).Seconds())
	if insertErr == nil {
		logger.Info(fmt.Sprintf("[SERVICE:Event] Successfully saved %s event batch: %d events", consumer, len(entities)))
	}
	return nil
}

// rejectInvalidResults Rejects every Phase-1 reject directly into the
// configured DLQ and records the per-message failure metric.
func rejectInvalidResults[T any](s *EventService, consumer string, results []messageResult[T]) {
	for i := range results {
		if results[i].action == "reject" {
			results[i].msg.Reject(results[i].rejectReason)
			s.deps.Metrics.MessagesTotal.WithLabelValues(consumer, "reject").Inc()
			s.deps.Metrics.EventsProcessed.WithLabelValues(consumer, "error").Inc()
		}
	}
}

// ackSkipParseFailures handles the DLQ-specific path: parse failures are
// Acked (not Rejected) so they never re-enter the loop.
func ackSkipParseFailures[T any](s *EventService, consumer string, results []messageResult[T]) {
	for i := range results {
		if results[i].action == "ack_skip" {
			results[i].msg.Ack()
			s.deps.Metrics.MessagesTotal.WithLabelValues(consumer, "ack").Inc()
			s.deps.Metrics.EventsProcessed.WithLabelValues(consumer, "error").Inc()
		}
	}
}

// runBulkInsert times the insert call and emits the table-scoped batch
// metric. Returns the insert error verbatim so the caller can drive
// Ack/Nack.
func runBulkInsert[T any](
	s *EventService,
	tableName string,
	entities []*T,
	insertFn func(ctx context.Context, entities []*T) error,
) error {
	s.deps.Metrics.ClickHouseInsertBatchSize.WithLabelValues(tableName).Observe(float64(len(entities)))
	insertStart := time.Now()

	err := insertFn(context.Background(), entities)

	s.deps.Metrics.ClickHouseInsertDuration.WithLabelValues(tableName).Observe(time.Since(insertStart).Seconds())
	if err != nil {
		s.deps.Metrics.ClickHouseInsertTotal.WithLabelValues(tableName, "error").Inc()
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] Failed to save %s batch: %d events", tableName, len(entities)))
	} else {
		s.deps.Metrics.ClickHouseInsertTotal.WithLabelValues(tableName, "ok").Inc()
	}
	return err
}

// finalizeBatchAck Acks every valid message on insert success or Nacks
// each one on failure (NATS retries them per the consumer's retry
// policy). Records dlq vs nack based on retry attempt headroom.
func finalizeBatchAck(s *EventService, consumer string, validMessages []*natsModel.Message, insertErr error) {
	for _, msg := range validMessages {
		if insertErr != nil {
			msg.Nack(insertErr)
			if attempt, max, _ := msg.GetRetryInfo(); max > 0 && attempt > max {
				s.deps.Metrics.MessagesTotal.WithLabelValues(consumer, "dlq").Inc()
			} else {
				s.deps.Metrics.MessagesTotal.WithLabelValues(consumer, "nack").Inc()
			}
			s.deps.Metrics.EventsProcessed.WithLabelValues(consumer, "error").Inc()
		} else {
			msg.Ack()
			s.deps.Metrics.MessagesTotal.WithLabelValues(consumer, "ack").Inc()
			s.deps.Metrics.EventsProcessed.WithLabelValues(consumer, "success").Inc()
		}
	}
}

// finalizeDLQBatchAck Acks every "pending" result in the DLQ flow,
// regardless of insert outcome — the DLQ stream never Nacks because
// redelivery would feed the failure back into the same DLQ.
func finalizeDLQBatchAck[T any](s *EventService, consumer string, results []messageResult[T], insertErr error) {
	for i := range results {
		if results[i].action != "pending" {
			continue
		}
		results[i].msg.Ack()
		s.deps.Metrics.MessagesTotal.WithLabelValues(consumer, "ack").Inc()
		if insertErr != nil {
			s.deps.Metrics.EventsProcessed.WithLabelValues(consumer, "error").Inc()
		} else {
			s.deps.Metrics.EventsProcessed.WithLabelValues(consumer, "success").Inc()
		}
	}
}

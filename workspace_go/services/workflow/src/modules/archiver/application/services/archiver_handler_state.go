package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	appConstants "workflow/src/modules/archiver/application/constants"
	"workflow/src/modules/archiver/application/ports"
	archiverTypes "workflow/src/modules/archiver/application/types"
	"workflow/src/modules/archiver/domain/repositories"
	runtimePorts "workflow/src/modules/runtime/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// classifiedState carries the per-bucket payload that ProcessStateBatch
// produces during classification and consumes during the write phase.
type classifiedState struct {
	createdStubs       []repositories.LightweightExecution
	waitingUpdates     []repositories.WaitingUpdate
	resumedIDs         []string
	terminalExecutions []*runtimePorts.WorkflowExecution
	terminalKVKeys     []string
	refs               []archiverTypes.MsgRef
}

// applyArchiverBackpressure pauses or warns based on MongoDB write P99 so
// the consumer respects the BackpressureMode set by the manager. Backoff
// mode sleeps the configured pause; throttled mode logs but proceeds.
func (s *ArchiverService) applyArchiverBackpressure() {
	bpMode := s.deps.MongoManager.GetBackpressureMode()
	if bpMode == ports.BackpressureBackoff {
		logger.Warn(fmt.Sprintf("[SERVICE:Archiver] Backoff mode active (P99=%dms), pausing %s before batch",
			s.deps.MongoManager.WriteP99(), appConstants.BackpressureBackoffPause))
		time.Sleep(appConstants.BackpressureBackoffPause)
	} else if bpMode == ports.BackpressureThrottled {
		logger.Warn(fmt.Sprintf("[SERVICE:Archiver] Throttled mode active (P99=%dms)",
			s.deps.MongoManager.WriteP99()))
	}
}

// classifyStateEvents walks the inbound batch, decodes each message, and
// routes it into the matching bucket (created/waiting/resumed/terminal).
// Malformed messages are rejected in place; terminal events that can't be
// loaded from KV are nacked. The returned struct drives the write phase.
func (s *ArchiverService) classifyStateEvents(messages []*natsModel.Message) classifiedState {
	classified := classifiedState{}
	for _, msg := range messages {
		var event ports.StateEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			msg.Reject(fmt.Sprintf("invalid state event: %s", err))
			continue
		}
		msg.OrgId = event.OrgID
		switch {
		case isCreatedEvent(msg.Subject):
			stub, err := buildLightweightStub(event)
			if err != nil {
				msg.Reject(fmt.Sprintf("invalid created event: %s", err))
				continue
			}
			classified.createdStubs = append(classified.createdStubs, stub)
			classified.refs = append(classified.refs, archiverTypes.MsgRef{Msg: msg, Batch: appConstants.BatchTagCreated})
		case isWaitingEvent(msg.Subject):
			update, err := buildWaitingUpdate(event)
			if err != nil {
				msg.Reject(fmt.Sprintf("invalid waiting event: %s", err))
				continue
			}
			classified.waitingUpdates = append(classified.waitingUpdates, update)
			classified.refs = append(classified.refs, archiverTypes.MsgRef{Msg: msg, Batch: appConstants.BatchTagWaiting})
		case isResumedEvent(msg.Subject):
			classified.resumedIDs = append(classified.resumedIDs, event.InstanceID)
			classified.refs = append(classified.refs, archiverTypes.MsgRef{Msg: msg, Batch: appConstants.BatchTagResumed})
		case isTerminalEvent(msg.Subject):
			instance, kvKey, err := s.fetchFullState(event.InstanceID)
			if err != nil {
				logger.Warn(fmt.Sprintf("[SERVICE:Archiver] KV Get failed for %s, NACKing: %s", event.InstanceID, err))
				msg.Nack(err)
				continue
			}
			classified.terminalExecutions = append(classified.terminalExecutions, instance)
			classified.terminalKVKeys = append(classified.terminalKVKeys, kvKey)
			classified.refs = append(classified.refs, archiverTypes.MsgRef{Msg: msg, Batch: appConstants.BatchTagTerminal})
		default:
			msg.Reject(fmt.Sprintf("unknown state event subject: %s", msg.Subject))
		}
	}
	return classified
}

// runArchiverWriteBatches runs the four MongoDB bulk-write batches against
// the classified buckets, recording write latency on each batch and acking
// (or nacking) the matching message refs based on the per-batch outcome.
func (s *ArchiverService) runArchiverWriteBatches(ctx context.Context, c classifiedState) {
	if len(c.createdStubs) > 0 {
		s.runCreatedBatch(ctx, c.createdStubs, c.refs)
	}
	if len(c.waitingUpdates) > 0 {
		s.runWaitingBatch(ctx, c.waitingUpdates, c.refs)
	}
	if len(c.resumedIDs) > 0 {
		s.runResumedBatch(ctx, c.resumedIDs, c.refs)
	}
	if len(c.terminalExecutions) > 0 {
		s.runTerminalBatch(ctx, c.terminalExecutions, c.terminalKVKeys, c.refs)
	}
}

// runCreatedBatch persists the lightweight insert batch for "created" events.
func (s *ArchiverService) runCreatedBatch(ctx context.Context, stubs []repositories.LightweightExecution, refs []archiverTypes.MsgRef) {
	start := time.Now()
	if err := s.deps.ArchiveRepo.BulkInsertLightweight(ctx, stubs); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Archiver] BulkInsertLightweight failed for %d stubs", len(stubs)))
		nackBatch(refs, appConstants.BatchTagCreated)
	} else {
		ackBatch(refs, appConstants.BatchTagCreated)
		logger.Info(fmt.Sprintf("[SERVICE:Archiver] Archived %d created events", len(stubs)))
	}
	s.deps.MongoManager.RecordWriteLatency(time.Since(start))
}

// runWaitingBatch persists the waiting-status update batch.
func (s *ArchiverService) runWaitingBatch(ctx context.Context, updates []repositories.WaitingUpdate, refs []archiverTypes.MsgRef) {
	start := time.Now()
	if err := s.deps.ArchiveRepo.BulkUpdateWaiting(ctx, updates); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Archiver] BulkUpdateWaiting failed for %d updates", len(updates)))
		nackBatch(refs, appConstants.BatchTagWaiting)
	} else {
		ackBatch(refs, appConstants.BatchTagWaiting)
		logger.Info(fmt.Sprintf("[SERVICE:Archiver] Archived %d waiting events", len(updates)))
	}
	s.deps.MongoManager.RecordWriteLatency(time.Since(start))
}

// runResumedBatch clears waiting timers for resumed instances.
func (s *ArchiverService) runResumedBatch(ctx context.Context, ids []string, refs []archiverTypes.MsgRef) {
	start := time.Now()
	if err := s.deps.ArchiveRepo.BulkUpdateResumed(ctx, ids); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Archiver] BulkUpdateResumed failed for %d instances", len(ids)))
		nackBatch(refs, appConstants.BatchTagResumed)
	} else {
		ackBatch(refs, appConstants.BatchTagResumed)
		logger.Info(fmt.Sprintf("[SERVICE:Archiver] Archived %d resumed events", len(ids)))
	}
	s.deps.MongoManager.RecordWriteLatency(time.Since(start))
}

// runTerminalBatch upserts the full-state batch for terminal events,
// publishes the matching ClickHouse cold-storage event, and cleans up the
// KV entries the runtime no longer needs.
func (s *ArchiverService) runTerminalBatch(ctx context.Context, executions []*runtimePorts.WorkflowExecution, kvKeys []string, refs []archiverTypes.MsgRef) {
	expireAt := time.Now().Add(appConstants.TerminalExecutionTTL)
	for _, exec := range executions {
		exec.ExpireAt = &expireAt
		if exec.CompletedAt == nil {
			now := time.Now()
			exec.CompletedAt = &now
		}
	}
	start := time.Now()
	if err := s.deps.ArchiveRepo.BulkUpsertFull(ctx, executions); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Archiver] BulkUpsertFull failed for %d executions", len(executions)))
		nackBatch(refs, appConstants.BatchTagTerminal)
		return
	}
	s.deps.MongoManager.RecordWriteLatency(time.Since(start))
	for _, exec := range executions {
		s.publishWorkflowEvent(exec)
	}
	for _, kvKey := range kvKeys {
		if err := s.deps.KVStore.Delete(kvKey); err != nil {
			logger.Warn(fmt.Sprintf("[SERVICE:Archiver] KV Delete failed for %s: %s", kvKey, err))
		}
	}
	ackBatch(refs, appConstants.BatchTagTerminal)
	logger.Info(fmt.Sprintf("[SERVICE:Archiver] Archived %d terminal events, cleaned %d KV keys, published to EVENTS-WORKFLOW",
		len(executions), len(kvKeys)))
}

package nats

import (
	"time"

	"assets/src/modules/assettemplates/application/ports"
	migrationMsg "assets/src/modules/assettemplates/interfaces/message"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// migrationScheduler implements ports.MigrationSchedulerPort using native NATS
// message scheduling. The start timer is published to the SINGLE static
// SubjectMigrationSchedule; the plan id travels in the payload + MsgId, never in
// the subject.
type migrationScheduler struct {
	sm natsModel.ScheduleManager
}

// NewMigrationScheduler returns a MigrationSchedulerPort backed by the NATS
// ScheduleManager.
func NewMigrationScheduler(sm natsModel.ScheduleManager) ports.MigrationSchedulerPort {
	return &migrationScheduler{sm: sm}
}

func (s *migrationScheduler) ScheduleStart(planID string, at time.Time) error {
	return s.sm.PublishScheduled(natsModel.ScheduledPublishConfig{
		Subject:       migrationMsg.SubjectMigrationSchedule,
		TargetSubject: migrationMsg.SubjectMigrationTimerStart,
		ScheduleAt:    at,
		Data:          map[string]string{"planId": planID},
		MsgId:         migrationMsg.MigrationStartMsgId(planID),
	})
}

// CancelStart is a no-op. The ScheduleManager exposes no per-MsgId cancel
// primitive, and the shared static subject rules out a subject-scoped purge of a
// single plan. Cancellation is enforced by the runner's status guard: a
// fired-but-cancelled plan is a no-op.
func (s *migrationScheduler) CancelStart(planID string) error {
	return nil
}

// Compile-time check that migrationScheduler implements the port.
var _ ports.MigrationSchedulerPort = (*migrationScheduler)(nil)

package nats

import (
	"time"

	"assets/src/modules/assettemplates/application/ports"
	migrationMsg "assets/src/modules/assettemplates/interfaces/message"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// migrationScheduler implements ports.MigrationSchedulerPort using native NATS
// message scheduling (ADR-51). In ADR-51 the SUBJECT is the schedule's identity, so
// each plan's start timer is published to its OWN schedule subject
// (SubjectMigrationSchedule + "." + planID). A shared static subject would make two
// concurrent plans overwrite each other's schedule slot — only the last would fire.
// The delivery TARGET stays static (SubjectMigrationTimerStart) with the plan id in
// the payload, so a single consumer handles every fired timer.
type migrationScheduler struct {
	sm natsModel.ScheduleManager
}

// NewMigrationScheduler returns a MigrationSchedulerPort backed by the NATS
// ScheduleManager.
func NewMigrationScheduler(sm natsModel.ScheduleManager) ports.MigrationSchedulerPort {
	return &migrationScheduler{sm: sm}
}

// scheduleSubject is the per-plan ADR-51 schedule subject: the plan id MUST be part
// of the subject because ADR-51 keys (and replaces) a schedule by its subject.
func scheduleSubject(planID string) string {
	return migrationMsg.SubjectMigrationSchedule + "." + planID
}

func (s *migrationScheduler) ScheduleStart(planID string, at time.Time) error {
	return s.sm.PublishScheduled(natsModel.ScheduledPublishConfig{
		Subject:       scheduleSubject(planID),
		TargetSubject: migrationMsg.SubjectMigrationTimerStart,
		ScheduleAt:    at,
		Data:          map[string]string{"planId": planID},
		MsgId:         migrationMsg.MigrationStartMsgId(planID),
	})
}

// CancelStart is a no-op. Cancellation is enforced by the runner's status guard: a
// fired-but-cancelled plan is a no-op. Now that each plan owns its schedule subject,
// this could purge that subject to drop the pending timer early, but the status guard
// already makes cancellation correct.
func (s *migrationScheduler) CancelStart(planID string) error {
	return nil
}

// Compile-time check that migrationScheduler implements the port.
var _ ports.MigrationSchedulerPort = (*migrationScheduler)(nil)

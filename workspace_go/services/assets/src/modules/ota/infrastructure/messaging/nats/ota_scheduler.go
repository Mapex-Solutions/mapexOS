package nats

import (
	"time"

	"assets/src/modules/ota/application/ports"
	otaMsg "assets/src/modules/ota/interfaces/message"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// otaScheduler implements ports.OTASchedulerPort using native NATS message
// scheduling. Every timer is published to the SINGLE static SubjectOTASchedule;
// the plan id travels in the payload + MsgId, never in the subject.
type otaScheduler struct {
	sm natsModel.ScheduleManager
}

// NewOTAScheduler returns an OTASchedulerPort backed by the NATS ScheduleManager.
func NewOTAScheduler(sm natsModel.ScheduleManager) ports.OTASchedulerPort {
	return &otaScheduler{sm: sm}
}

func (s *otaScheduler) ScheduleStart(planID string, at time.Time) error {
	return s.sm.PublishScheduled(natsModel.ScheduledPublishConfig{
		Subject:       otaMsg.SubjectOTASchedule,
		TargetSubject: otaMsg.SubjectOTATimerStart,
		ScheduleAt:    at,
		Data:          map[string]string{"planId": planID},
		MsgId:         otaMsg.StartMsgId(planID),
	})
}

func (s *otaScheduler) ScheduleClose(planID string, at time.Time) error {
	return s.sm.PublishScheduled(natsModel.ScheduledPublishConfig{
		Subject:       otaMsg.SubjectOTASchedule,
		TargetSubject: otaMsg.SubjectOTATimerClose,
		ScheduleAt:    at,
		Data:          map[string]string{"planId": planID},
		MsgId:         otaMsg.CloseMsgId(planID),
	})
}

func (s *otaScheduler) ScheduleScan(at time.Time) error {
	return s.sm.PublishScheduled(natsModel.ScheduledPublishConfig{
		Subject:       otaMsg.SubjectOTASchedule,
		TargetSubject: otaMsg.SubjectOTATimerScan,
		ScheduleAt:    at,
		Data:          map[string]string{"trigger": "scheduled"},
		MsgId:         otaMsg.ScanMsgId,
	})
}

// Compile-time check that otaScheduler implements the port.
var _ ports.OTASchedulerPort = (*otaScheduler)(nil)

package nats

import (
	"time"

	"assets/src/modules/ota/application/ports"
	otaMsg "assets/src/modules/ota/interfaces/message"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// firmwareScheduler implements ports.FirmwareSchedulerPort using native NATS
// scheduling — the abandon-check @at fires to the static ota.timer.abandon
// subject (firmwareId in the payload).
type firmwareScheduler struct {
	sm natsModel.ScheduleManager
}

// NewFirmwareScheduler returns a FirmwareSchedulerPort over the NATS
// ScheduleManager.
func NewFirmwareScheduler(sm natsModel.ScheduleManager) ports.FirmwareSchedulerPort {
	return &firmwareScheduler{sm: sm}
}

func (s *firmwareScheduler) ScheduleAbandonCheck(firmwareID string, at time.Time) error {
	return s.sm.PublishScheduled(natsModel.ScheduledPublishConfig{
		Subject:       otaMsg.SubjectOTASchedule,
		TargetSubject: otaMsg.SubjectOTATimerAbandon,
		ScheduleAt:    at,
		Data:          map[string]string{"firmwareId": firmwareID},
		MsgId:         otaMsg.AbandonMsgId(firmwareID),
	})
}

// PurgeAbandonCheck is a no-op: HandleFirmwareAbandon is idempotent (it no-ops
// once the firmware is finalized), so a pending abandon-check fires harmlessly.
func (s *firmwareScheduler) PurgeAbandonCheck(firmwareID string) error { return nil }

var _ ports.FirmwareSchedulerPort = (*firmwareScheduler)(nil)

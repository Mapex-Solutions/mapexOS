package services

import (
	"fmt"
	"time"

	hmMessage "assets/src/modules/healthmonitor/interfaces/message"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// scheduleNextScan publishes a scheduled scan message if none is pending.
// Layer 1: HasPendingMessages check (covers 99% of cases).
// Layer 2: Duplicates=10s + fixed MsgId on stream (catches race window between pods).
func (s *HealthMonitorService) scheduleNextScan() {
	logger.Debug(fmt.Sprintf("[SERVICE:HealthMonitor] [SCHEDULER] checking for pending scan: stream=%s subject=%s",
		hmMessage.ScanStreamName, hmMessage.ScanScheduleSubject))

	pending, err := s.deps.ScheduleManager.HasPendingMessages(hmMessage.ScanStreamName, hmMessage.ScanScheduleSubject)
	if err != nil {
		// Check failed — proceed with publish anyway to keep the loop alive.
		// MsgId dedup on the stream prevents actual duplicates.
		logger.Warn(fmt.Sprintf("[SERVICE:HealthMonitor] [SCHEDULER] pending check failed (publishing anyway to keep loop alive): err=%s", err))
	} else if pending {
		logger.Info("[SERVICE:HealthMonitor] [SCHEDULER] scan already pending in NATS — skipping publish")
		return
	}

	interval, _ := config.GetIntValue("health_monitor_scan_interval")
	if interval <= 0 {
		interval = 600
	}

	scheduleAt := time.Now().Add(time.Duration(interval) * time.Second)

	logger.Debug(fmt.Sprintf("[SERVICE:HealthMonitor] [SCHEDULER] publishing next scan: interval=%ds scheduleAt=%s msgId=%s",
		interval, scheduleAt.Format(time.RFC3339), hmMessage.ScanMsgId))

	if err := s.deps.ScheduleManager.PublishScheduled(natsModel.ScheduledPublishConfig{
		Subject:       hmMessage.ScanScheduleSubject,
		TargetSubject: hmMessage.ScanSubject,
		ScheduleAt:    scheduleAt,
		Data:          map[string]string{"trigger": "scheduled"},
		MsgId:         hmMessage.ScanMsgId,
	}); err != nil {
		logger.Error(err, "[SERVICE:HealthMonitor] [SCHEDULER] failed to publish scheduled scan")
		return
	}

	logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [SCHEDULER] next scan scheduled: interval=%ds scheduleAt=%s",
		interval, scheduleAt.Format(time.RFC3339)))
}

package ports

import "time"

// OTASchedulerPort schedules a plan's two timers (startAt / maxTime) via native
// NATS message scheduling. The NATS implementation lives in the messaging layer
// (it is defined here as the application's driven port).
type OTASchedulerPort interface {
	// ScheduleStart schedules a plan's start timer (the planId travels in the
	// message payload, never in the subject).
	ScheduleStart(planID string, at time.Time) error
	// ScheduleClose schedules a plan's close (maxTime) timer.
	ScheduleClose(planID string, at time.Time) error
}

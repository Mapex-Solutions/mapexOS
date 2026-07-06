package ports

import "time"

// MigrationSchedulerPort schedules a migration plan's start timer via native NATS
// message scheduling. The NATS implementation lives in the messaging layer (it is
// defined here as the application's driven port).
type MigrationSchedulerPort interface {
	// ScheduleStart schedules a plan's start timer (the planId travels in the
	// message payload, never in the subject).
	ScheduleStart(planID string, at time.Time) error
	// CancelStart cancels a plan's pending start timer.
	CancelStart(planID string) error
}

package message

import config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"

// Migration timer subjects. Convention (two kinds, different rules):
//   - A SCHEDULE subject (ADR-51) MUST carry an identifier, because in ADR-51 the
//     subject IS the schedule's identity — a static one makes concurrent schedules
//     overwrite each other. SubjectMigrationSchedule is the BASE; the scheduler
//     appends ".{planId}" per plan (see migrationScheduler.scheduleSubject).
//   - Every other (internal/delivery) subject stays STATIC with the id in the
//     payload — e.g. the delivery target below, read by one consumer.
var (
	// SubjectMigrationSchedule is the BASE schedule subject; the scheduler publishes
	// each plan's @at timer to SubjectMigrationSchedule + "." + planID so every plan
	// owns its own ADR-51 schedule slot.
	SubjectMigrationSchedule = config.Subject("assettemplates", "migration.schedule")

	// SubjectMigrationTimerStart is the STATIC target subject the scheduled message
	// is delivered to when it fires; the plan id travels in the payload.
	SubjectMigrationTimerStart = config.Subject("assettemplates", "migration.timer.start")

	// MigrationScheduleStream carries the scheduled + delivered migration timer
	// messages. Subjects: assettemplates.migration.> (covers migration.schedule.* and
	// migration.timer.start).
	MigrationScheduleStream = config.StreamName("ASSETTEMPLATES", "MIGRATION")
)

// MigrationStartMsgId dedups a plan's start timer within the stream's Duplicates
// window. Unique per plan (the schedule subject already isolates plans; this guards
// against a double publish of the same plan).
func MigrationStartMsgId(planID string) string { return "tmig-start-" + planID }

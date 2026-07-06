package message

import config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"

// Migration timer scheduling subjects. The subjects are STATIC — the plan id
// travels in the message payload (and the MsgId), NEVER in the subject.
var (
	// SubjectMigrationSchedule is the single static subject where scheduled
	// migration timer messages wait until their @at time.
	SubjectMigrationSchedule = config.Subject("assettemplates", "migration.schedule")

	// SubjectMigrationTimerStart is the target subject the scheduled message is
	// delivered to when it fires.
	SubjectMigrationTimerStart = config.Subject("assettemplates", "migration.timer.start")

	// MigrationScheduleStream carries the scheduled + delivered migration timer
	// messages. Subjects: assettemplates.migration.>.
	MigrationScheduleStream = config.StreamName("ASSETTEMPLATES", "MIGRATION")
)

// MigrationStartMsgId dedups a plan's start timer. The plan id lives in the MsgId
// metadata + payload — NOT the subject.
func MigrationStartMsgId(planID string) string { return "tmig-start-" + planID }

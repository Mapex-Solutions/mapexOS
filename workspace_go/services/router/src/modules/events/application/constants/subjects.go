package constants

import (
	"fmt"

	routerEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/router/events"
)

// KindToSubjectMap maps router kind to NATS subject.
// Subject strings are sourced from the cross-service contract
// (packages/contracts/services/router/events) — never redefined here.
var KindToSubjectMap = map[string]string{
	"save_event":   routerEvents.SubjectSaveEvent,
	"lake_house":   routerEvents.SubjectLakeHouse,
	"notification": routerEvents.SubjectNotification,
}

// GetSubjectByKind returns the NATS subject for a given router kind.
func GetSubjectByKind(kind string) (string, error) {
	subject, exists := KindToSubjectMap[kind]
	if !exists {
		return "", fmt.Errorf("unknown or dynamic router kind: %s", kind)
	}
	return subject, nil
}

// GetTriggerSubject returns the subject for trigger kind.
// Uses a fixed subject — triggerId is sent in the payload.
func GetTriggerSubject(triggerId string) string {
	return routerEvents.SubjectTriggerRouterExecute
}

package services

import (
	"encoding/json"
	"fmt"

	appConstants "workflow/src/modules/runtime/application/constants"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// parseScheduleFireBody decodes the schedule body and extracts the instanceId
// used to address the resume.timer subject. Reject codes failures here as
// permanent: a malformed schedule body or a missing instanceId means there is
// no execution to deliver to, retrying will not help.
func (s *RuntimeService) parseScheduleFireBody(msg *natsModel.Message) (map[string]interface{}, string, bool) {
	var body map[string]interface{}
	if err := json.Unmarshal(msg.Data, &body); err != nil {
		msg.Reject(fmt.Sprintf("invalid schedule fire message: %s", err))
		return nil, "", false
	}

	instanceId, _ := body[appConstants.ExecDataKeyInstanceID].(string)
	if instanceId == "" {
		msg.Reject("schedule fire message missing instanceId")
		return nil, "", false
	}
	return body, instanceId, true
}

// republishScheduleResume forwards the schedule body to the WORKFLOW-RESUME
// stream via mapexos.workflow.resume.timer.{instanceId}. A publish error
// triggers a Nack so the schedule fire is retried; success acks the message.
func (s *RuntimeService) republishScheduleResume(msg *natsModel.Message, instanceId string, body map[string]interface{}) {
	if err := s.deps.RuntimePublisher.PublishResumeTimer(instanceId, body); err != nil {
		logger.Error(err, fmt.Sprintf("[CONSUMER:ScheduleFire] Failed to re-publish for instanceId=%s", instanceId))
		msg.Nack(err)
		return
	}
	logger.Debug(fmt.Sprintf("[CONSUMER:ScheduleFire] Re-published (instanceId=%s)", instanceId))
	msg.Ack()
}

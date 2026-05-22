package services

import (
	"fmt"

	appConstants "workflow/src/modules/runtime/application/constants"
	"workflow/src/modules/runtime/domain/constants"
	"workflow/src/modules/runtime/domain/entities"

	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/executions"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// handleSignalMode delivers a signal to a waiting execution.
func (s *RuntimeService) handleSignalMode(msg *natsModel.Message, execMsg *v1.WorkflowExecutionMessage) {
	workflowUUID, _ := execMsg.Data[appConstants.ExecDataKeyWorkflowUUID].(string)
	signalName, _ := execMsg.Data[appConstants.ExecDataKeySignalName].(string)
	signalData, _ := execMsg.Data[appConstants.ExecDataKeySignalData].(map[string]interface{})

	if workflowUUID == "" {
		msg.Reject("data.workflowUUID is required for mode 'signal'")
		return
	}
	if signalName == "" {
		msg.Reject("data.signalName is required for mode 'signal'")
		return
	}

	if err := s.deliverSignal(workflowUUID, signalName, signalData); err != nil {
		msg.Reject(fmt.Sprintf("signal delivery failed: %s", err))
		return
	}

	logger.Info(fmt.Sprintf("[SERVICE:Runtime] Execution signal → uuid=%s signal=%s", workflowUUID, signalName))
	msg.Ack()
}

// handleSignalOrStart tries signal delivery first, falls back to newInstance.
func (s *RuntimeService) handleSignalOrStart(msg *natsModel.Message, execMsg *v1.WorkflowExecutionMessage) {
	instanceID, _ := execMsg.Data[appConstants.ExecDataKeyInstanceID].(string)
	workflowUUID, _ := execMsg.Data[appConstants.ExecDataKeyWorkflowUUID].(string)
	signalName, _ := execMsg.Data[appConstants.ExecDataKeySignalName].(string)
	signalData, _ := execMsg.Data[appConstants.ExecDataKeySignalData].(map[string]interface{})

	if instanceID == "" || workflowUUID == "" || signalName == "" {
		msg.Reject("data.instanceId, data.workflowUUID, and data.signalName are required for mode 'signalOrStart'")
		return
	}

	err := s.deliverSignal(workflowUUID, signalName, signalData)
	if err == nil {
		logger.Info(fmt.Sprintf("[SERVICE:Runtime] Execution signalOrStart → signal delivered uuid=%s signal=%s", workflowUUID, signalName))
		msg.Ack()
		return
	}

	logger.Info(fmt.Sprintf("[SERVICE:Runtime] Execution signalOrStart → signal failed (%s), falling back to newInstance", err))
	s.handleNewInstance(msg, execMsg)
}

// deliverSignal reads execution from KV, validates waiting + signalName match,
// and publishes a resume message to WORKFLOW-RESUME.
func (s *RuntimeService) deliverSignal(executionID string, signalName string, data map[string]interface{}) error {
	execution, err := s.deps.ExecutionStateRepo.Get(executionID)
	if err != nil {
		return fmt.Errorf("execution not found in KV: %w", err)
	}

	if execution.Status != entities.ExecStatusWaiting || execution.NodeStates == nil {
		return fmt.Errorf("execution is not waiting for a signal (status=%s)", execution.Status)
	}

	var waitNodeID string
	for nid, ns := range execution.NodeStates {
		if ns[appConstants.NodeStateKeyWaitType] == constants.WaitTypeSignal {
			nodeSigName, _ := ns[appConstants.NodeStateKeySignalName].(string)
			if nodeSigName == signalName {
				waitNodeID = nid
				break
			}
		}
	}
	if waitNodeID == "" {
		return fmt.Errorf("no node waiting for signal %q", signalName)
	}

	// Purge pending schedule for the signal node (signal arrived before timeout)
	if err := s.deps.RuntimePublisher.PurgeSchedule(executionID, waitNodeID); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Runtime] Failed to purge schedule for signal %s node %s: %s",
			executionID, waitNodeID, err))
	}

	return s.deps.RuntimePublisher.PublishSignalResume(executionID, waitNodeID, data)
}

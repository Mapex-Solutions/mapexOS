package dagwalker

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/shutdown"

	defEntities "workflow/src/modules/definitions/domain/entities"
	instanceEntities "workflow/src/modules/instances/domain/entities"
	"workflow/src/modules/runtime/application/di"
	runtimePorts "workflow/src/modules/runtime/application/ports"
	"workflow/src/modules/runtime/application/services"
	"workflow/src/modules/runtime/domain/entities"
	sharedTypes "workflow/src/shared/types"

	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/executions"

	// Real engine components (all 22 operators registered)
	engineServices "workflow/src/modules/engine/application/services"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

const maxResumeIterations = 500 // Safety limit for async callback loop

// WorkflowTestHarness orchestrates DAG walker E2E tests.
// Constructs a real RuntimeService with mocked dependencies.
type WorkflowTestHarness struct {
	t          *testing.T
	service    runtimePorts.RuntimeServicePort
	StateRepo  *InMemoryStateRepo
	Publisher  *CapturingPublisher
	DefLoader  *StaticDefinitionLoader
	InstLoader *StaticInstanceLoader
	ShutdownMgr *shutdown.ShutdownManager

	definition *defEntities.WorkflowDefinition
}

// NewHarness creates a fully wired test harness for a given workflow definition.
func NewHarness(t *testing.T, def *defEntities.WorkflowDefinition) *WorkflowTestHarness {
	t.Helper()

	stateRepo := NewInMemoryStateRepo()
	publisher := NewCapturingPublisher()

	// Create a default test instance
	instance := &instanceEntities.WorkflowInstance{
		ID:           model.NewObjectID(),
		DefinitionID: def.ID,
		Name:         def.Name,
		OrgID:        def.OrgID,
	}

	defLoader := &StaticDefinitionLoader{Definition: def}
	instLoader := &StaticInstanceLoader{Instance: instance}

	// Real engine components with all 22 operators registered
	condEvaluator, valueResolver := engineServices.New()

	// Construct DI struct manually (dig.In is a zero-size marker, works fine)
	sm := shutdown.New()

	deps := di.RuntimeServiceDependenciesInjection{
		ExecutionStateRepo: stateRepo,
		RuntimePublisher:   publisher,
		ConditionEvaluator: condEvaluator,
		ValueResolver:      valueResolver,
		DefinitionLoader:   defLoader,
		InstanceLoader:     instLoader,
		VaultService:       &NoopVaultService{},
		PluginRepo:         &NoopPluginRepo{},
		ShutdownManager:    sm,
		Metrics:            NewNoopMetrics(),
	}

	svc := services.New(deps)

	return &WorkflowTestHarness{
		t:           t,
		service:     svc,
		StateRepo:   stateRepo,
		Publisher:    publisher,
		DefLoader:   defLoader,
		InstLoader:  instLoader,
		ShutdownMgr: sm,
		definition:  def,
	}
}

// RunSync triggers a workflow execution and expects it to complete synchronously
// (no async nodes). Returns the final WorkflowExecution.
func (h *WorkflowTestHarness) RunSync(eventPayload map[string]interface{}) *entities.WorkflowExecution {
	h.t.Helper()
	h.trigger(eventPayload)
	exec := h.StateRepo.GetLatest()
	if exec == nil {
		h.t.Fatal("no execution found after trigger")
	}
	return exec
}

// RunSyncAllowNack triggers a workflow and tolerates Nack (for dispatch failure tests).
func (h *WorkflowTestHarness) RunSyncAllowNack(eventPayload map[string]interface{}) *entities.WorkflowExecution {
	h.t.Helper()
	h.triggerAllowNack(eventPayload)
	exec := h.StateRepo.GetLatest()
	if exec == nil {
		h.t.Fatal("no execution found after trigger")
	}
	return exec
}

// RunWithCallbacks triggers a workflow and processes async callbacks until completion.
// The callbacks map is keyed by nodeID — when that node suspends, the callback is invoked
// to generate a ResumeMessage.
func (h *WorkflowTestHarness) RunWithCallbacks(
	eventPayload map[string]interface{},
	callbacks map[string]CallbackFunc,
) *entities.WorkflowExecution {
	h.t.Helper()
	h.trigger(eventPayload)

	for i := 0; i < maxResumeIterations; i++ {
		exec := h.StateRepo.GetLatest()
		if exec == nil {
			h.t.Fatal("execution disappeared from state repo")
		}

		// Terminal state — done
		if exec.Status.IsTerminal() {
			return exec
		}

		// Not waiting — might be running (re-enqueue case)
		if exec.Status != entities.ExecStatusWaiting {
			h.t.Fatalf("unexpected status %q after iteration %d (expected waiting or terminal)", exec.Status, i)
		}

		// Process each waiting node
		resumed := false
		for _, nodeID := range exec.ActiveNodeIDs {
			ns := exec.NodeStates[nodeID]
			if ns == nil || ns["waitType"] == nil {
				continue
			}

			// Retry timer: simulate NATS Schedule timeout resume
			if ns["waitType"] == "retryTimer" {
				resume := sharedTypes.ResumeMessage{
					InstanceID: exec.WorkflowUUID,
					NodeID:     nodeID,
					IsTimeout:  true,
				}
				h.sendResume(resume)
				resumed = true
				continue
			}

			cb, ok := callbacks[nodeID]
			if !ok {
				// Use default callback based on waitType
				cb = DefaultAsyncCallback()
			}

			resume := cb(nodeID, ns)
			resume.InstanceID = exec.WorkflowUUID
			// Echo execution token from NodeState (G5 backward compat — if present)
			if token, ok := ns["executionToken"].(string); ok {
				resume.ExecutionToken = token
			}
			h.sendResume(resume)
			resumed = true
		}

		if !resumed {
			h.t.Fatalf("execution stuck in waiting with no resumable nodes (activeNodes=%v)", exec.ActiveNodeIDs)
		}
	}

	exec := h.StateRepo.GetLatest()
	h.t.Fatalf("exceeded max resume iterations (%d), execution still in status %q", maxResumeIterations, exec.Status)
	return nil
}

// GetExecution returns the current execution by UUID.
func (h *WorkflowTestHarness) GetExecution(uuid string) *entities.WorkflowExecution {
	exec, err := h.StateRepo.Get(uuid)
	if err != nil {
		h.t.Fatalf("execution %s not found: %s", uuid, err)
	}
	return exec
}

// DebugPrintPath prints the execution path for debugging.
func (h *WorkflowTestHarness) DebugPrintPath(exec *entities.WorkflowExecution) {
	for i, entry := range exec.ExecutionPath {
		fmt.Printf("  [%d] %s (%s) status=%s handle=%s\n",
			i, entry.NodeID, entry.NodeType, entry.Status, entry.OutputHandle)
	}
}

// trigger sends a WorkflowExecutionMessage (mode=newInstance) to HandleExecution.
func (h *WorkflowTestHarness) trigger(eventPayload map[string]interface{}) {
	h.t.Helper()

	instID := h.InstLoader.Instance.ID.Hex()

	execMsg := v1.WorkflowExecutionMessage{
		Mode:  "newInstance",
		Event: eventPayload,
		Data: map[string]interface{}{
			"instanceId": instID,
		},
	}

	data, err := json.Marshal(execMsg)
	if err != nil {
		h.t.Fatalf("failed to marshal execution message: %s", err)
	}

	acked := false
	msg := natsModel.NewTestMessage(data, 0, &natsModel.TestMessageCallbacks{
		OnAck: func() error {
			acked = true
			return nil
		},
		OnNack: func(err error) error {
			h.t.Fatalf("execution message NACKed: %s", err)
			return nil
		},
		OnReject: func(reason string) error {
			h.t.Fatalf("execution message rejected: %s", reason)
			return nil
		},
	})

	h.service.HandleExecution(msg)

	if !acked {
		h.t.Log("warning: execution message was not ACKed")
	}
}

// triggerAllowNack sends an execution message tolerating Nack (for dispatch failure tests).
func (h *WorkflowTestHarness) triggerAllowNack(eventPayload map[string]interface{}) {
	h.t.Helper()

	execMsg := v1.WorkflowExecutionMessage{
		Mode: "newInstance",
		Data: map[string]interface{}{
			"instanceId": h.InstLoader.Instance.ID.Hex(),
		},
		Event: eventPayload,
	}

	data, err := json.Marshal(execMsg)
	if err != nil {
		h.t.Fatalf("failed to marshal execution message: %s", err)
	}

	msg := natsModel.NewTestMessage(data, 0, &natsModel.TestMessageCallbacks{
		OnAck:    func() error { return nil },
		OnNack:   func(err error) error { return nil }, // Tolerate Nack
		OnReject: func(reason string) error { return nil },
	})

	h.service.HandleExecution(msg)
}

// ResumeExecution sends a re-enqueue resume for an existing execution (simulates pod restart recovery).
func (h *WorkflowTestHarness) ResumeExecution(exec *entities.WorkflowExecution) {
	h.t.Helper()
	activeNode := ""
	if len(exec.ActiveNodeIDs) > 0 {
		activeNode = exec.ActiveNodeIDs[0]
	}
	h.sendResume(sharedTypes.ResumeMessage{
		InstanceID: exec.WorkflowUUID,
		NodeID:     activeNode,
	})
}

// SendResumeWithToken sends a resume with a specific execution token (for token validation tests).
func (h *WorkflowTestHarness) SendResumeWithToken(exec *entities.WorkflowExecution, nodeID, token string, resume sharedTypes.ResumeMessage) {
	h.t.Helper()
	resume.InstanceID = exec.WorkflowUUID
	resume.NodeID = nodeID
	resume.ExecutionToken = token

	data, err := json.Marshal(resume)
	if err != nil {
		h.t.Fatalf("failed to marshal resume message: %s", err)
	}

	msg := natsModel.NewTestMessage(data, 0, &natsModel.TestMessageCallbacks{
		OnAck:    func() error { return nil },
		OnNack:   func(err error) error { return nil },
		OnReject: func(reason string) error { return nil },
	})

	h.service.HandleResume(msg)
}

// sendResume sends a ResumeMessage to HandleResume.
func (h *WorkflowTestHarness) sendResume(resume sharedTypes.ResumeMessage) {
	h.t.Helper()

	data, err := json.Marshal(resume)
	if err != nil {
		h.t.Fatalf("failed to marshal resume message: %s", err)
	}

	msg := natsModel.NewTestMessage(data, 0, &natsModel.TestMessageCallbacks{
		OnAck: func() error { return nil },
		OnNack: func(err error) error {
			h.t.Fatalf("resume message NACKed: %s", err)
			return nil
		},
		OnReject: func(reason string) error {
			h.t.Fatalf("resume message rejected: %s", reason)
			return nil
		},
	})

	h.service.HandleResume(msg)
}

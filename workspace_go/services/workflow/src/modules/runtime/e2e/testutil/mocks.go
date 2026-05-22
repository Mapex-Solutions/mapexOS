package dagwalker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	defEntities "workflow/src/modules/definitions/domain/entities"
	instanceEntities "workflow/src/modules/instances/domain/entities"
	instancePorts "workflow/src/modules/instances/application/ports"
	pluginEntities "workflow/src/modules/plugins/domain/entities"
	pluginPorts "workflow/src/modules/plugins/application/ports"
	runtimePorts "workflow/src/modules/runtime/application/ports"
	"workflow/src/modules/runtime/domain/entities"
	"workflow/src/modules/runtime/domain/repositories"
	sharedTypes "workflow/src/shared/types"

	"workflow/src/bootstrap"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"github.com/prometheus/client_golang/prometheus"
)

/**
 * InMemoryStateRepo
 * Implements ExecutionStateRepository with an in-memory map.
 * Thread-safe for concurrent access during fanout goroutines.
 */
type InMemoryStateRepo struct {
	mu           sync.RWMutex
	executions   map[string]*entities.WorkflowExecution
	revision     uint64
	saveCount    int
	OnCheckpoint func(count int) // called after each Save with the cumulative count
}

var _ repositories.ExecutionStateRepository = (*InMemoryStateRepo)(nil)

func NewInMemoryStateRepo() *InMemoryStateRepo {
	return &InMemoryStateRepo{
		executions: make(map[string]*entities.WorkflowExecution),
	}
}

func (r *InMemoryStateRepo) Create(execution *entities.WorkflowExecution) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.executions[execution.WorkflowUUID]; exists {
		return fmt.Errorf("execution %s already exists", execution.WorkflowUUID)
	}
	r.executions[execution.WorkflowUUID] = r.deepCopy(execution)
	return nil
}

func (r *InMemoryStateRepo) Get(executionID string) (*entities.WorkflowExecution, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	exec, exists := r.executions[executionID]
	if !exists {
		return nil, fmt.Errorf("execution %s not found", executionID)
	}
	return r.deepCopy(exec), nil
}

func (r *InMemoryStateRepo) Save(execution *entities.WorkflowExecution) error {
	r.mu.Lock()
	r.executions[execution.WorkflowUUID] = r.deepCopy(execution)
	r.saveCount++
	count := r.saveCount
	cb := r.OnCheckpoint
	r.mu.Unlock()
	if cb != nil {
		cb(count)
	}
	return nil
}

// GetWithRevision retrieves the execution state plus revision for CAS operations.
func (r *InMemoryStateRepo) GetWithRevision(executionID string) (*entities.WorkflowExecution, uint64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	exec, exists := r.executions[executionID]
	if !exists {
		return nil, 0, fmt.Errorf("execution %s not found", executionID)
	}
	return r.deepCopy(exec), r.revision, nil
}

// SaveWithRevision checkpoints using CAS. Mock does not enforce revision check.
func (r *InMemoryStateRepo) SaveWithRevision(execution *entities.WorkflowExecution, revision uint64) error {
	return r.Save(execution)
}

// GetLatest returns the most recently saved execution (for single-execution tests).
func (r *InMemoryStateRepo) GetLatest() *entities.WorkflowExecution {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, exec := range r.executions {
		return r.deepCopy(exec)
	}
	return nil
}

func (r *InMemoryStateRepo) deepCopy(exec *entities.WorkflowExecution) *entities.WorkflowExecution {
	data, _ := json.Marshal(exec)
	var copy entities.WorkflowExecution
	_ = json.Unmarshal(data, &copy)
	// Preserve ObjectId fields that may not serialize cleanly
	copy.ID = exec.ID
	copy.InstanceID = exec.InstanceID
	copy.DefinitionID = exec.DefinitionID
	copy.OrgID = exec.OrgID
	return &copy
}

/**
 * CapturingPublisher
 * Implements RuntimePublisherPort by recording all calls for test assertions.
 */

// PublishedEvent records a single call to the publisher.
type PublishedEvent struct {
	Method      string
	ExecutionID string
	NodeID      string
	Status      string
	Mode        string
	Data        map[string]interface{}
}

type CapturingPublisher struct {
	mu          sync.Mutex
	Events      []PublishedEvent
	failMethods map[string]bool
}

// FailOnMethod configures the publisher to return an error when the given method is called.
func (p *CapturingPublisher) FailOnMethod(method string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.failMethods == nil {
		p.failMethods = make(map[string]bool)
	}
	p.failMethods[method] = true
}

var _ runtimePorts.RuntimePublisherPort = (*CapturingPublisher)(nil)

func NewCapturingPublisher() *CapturingPublisher {
	return &CapturingPublisher{}
}

func (p *CapturingPublisher) PublishStateEvent(execution *entities.WorkflowExecution, status string) error {
	p.record("PublishStateEvent", execution.WorkflowUUID, "", status, "", nil)
	return nil
}

func (p *CapturingPublisher) PublishResumeMessage(executionID string, nodeID string, status string) error {
	p.record("PublishResumeMessage", executionID, nodeID, status, "", nil)
	return nil
}

func (p *CapturingPublisher) PublishResumeTimer(instanceID string, body map[string]interface{}) error {
	p.record("PublishResumeTimer", instanceID, "", "", "", body)
	return nil
}

func (p *CapturingPublisher) DispatchCodeExecution(execution *entities.WorkflowExecution, nodeID string, nodeState map[string]interface{}, executionToken, msgId string) error {
	p.record("DispatchCodeExecution", execution.WorkflowUUID, nodeID, "", "", nodeState)
	if p.failMethods["DispatchCodeExecution"] {
		return fmt.Errorf("simulated dispatch failure for DispatchCodeExecution")
	}
	return nil
}

func (p *CapturingPublisher) DispatchSubworkflowExecution(execution *entities.WorkflowExecution, nodeID string, nodeState map[string]interface{}, executionToken, msgId string) error {
	p.record("DispatchSubworkflowExecution", execution.WorkflowUUID, nodeID, "", "", nodeState)
	if p.failMethods["DispatchSubworkflowExecution"] {
		return fmt.Errorf("simulated dispatch failure for DispatchSubworkflowExecution")
	}
	return nil
}

func (p *CapturingPublisher) DispatchWorkflowTrigger(execution *entities.WorkflowExecution, nodeID string, mode string, data map[string]interface{}, executionToken, msgId string) error {
	p.record("DispatchWorkflowTrigger", execution.WorkflowUUID, nodeID, "", mode, data)
	if p.failMethods["DispatchWorkflowTrigger"] {
		return fmt.Errorf("simulated dispatch failure for DispatchWorkflowTrigger")
	}
	return nil
}

func (p *CapturingPublisher) PublishSignalResume(executionID string, nodeID string, signalData map[string]interface{}) error {
	p.record("PublishSignalResume", executionID, nodeID, "", "", signalData)
	return nil
}

func (p *CapturingPublisher) PublishCallbackResume(subject string, resume sharedTypes.ResumeMessage) error {
	p.record("PublishCallbackResume", resume.InstanceID, resume.NodeID, resume.Status, subject, nil)
	return nil
}

func (p *CapturingPublisher) PublishSchedule(wfUUID, nodeID string, expiresAt time.Time, waitType string, enableOutput bool) error {
	p.record("PublishSchedule", wfUUID, nodeID, waitType, "", map[string]interface{}{
		"expiresAt": expiresAt, "enableOutput": enableOutput,
	})
	return nil
}

func (p *CapturingPublisher) PurgeSchedule(wfUUID, nodeID string) error {
	p.record("PurgeSchedule", wfUUID, nodeID, "", "", nil)
	return nil
}

func (p *CapturingPublisher) PurgeAllSchedules(wfUUID string) error {
	p.record("PurgeAllSchedules", wfUUID, "", "", "", nil)
	return nil
}

func (p *CapturingPublisher) record(method, execID, nodeID, status, mode string, data map[string]interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Events = append(p.Events, PublishedEvent{
		Method:      method,
		ExecutionID: execID,
		NodeID:      nodeID,
		Status:      status,
		Mode:        mode,
		Data:        data,
	})
}

// FindDispatch returns the first event matching the method and nodeID.
func (p *CapturingPublisher) FindDispatch(method, nodeID string) *PublishedEvent {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, e := range p.Events {
		if e.Method == method && e.NodeID == nodeID {
			return &p.Events[i]
		}
	}
	return nil
}

// CountMethod returns how many times a method was called.
func (p *CapturingPublisher) CountMethod(method string) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	count := 0
	for _, e := range p.Events {
		if e.Method == method {
			count++
		}
	}
	return count
}

// NewNoopMetrics creates a WorkflowMetrics with all fields initialized but unregistered.
// Safe for e2e tests — prevents nil panics without polluting a real Prometheus registry.
func NewNoopMetrics() *bootstrap.WorkflowMetrics {
	return &bootstrap.WorkflowMetrics{
		DefinitionOperations:        prometheus.NewCounterVec(prometheus.CounterOpts{Name: "noop_def_ops"}, []string{"operation", "status"}),
		DefinitionOperationDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "noop_def_dur"}, []string{"operation"}),
		DefinitionListResultsCount:  prometheus.NewHistogram(prometheus.HistogramOpts{Name: "noop_def_list"}),
		DefinitionCacheTotal:        prometheus.NewCounterVec(prometheus.CounterOpts{Name: "noop_def_cache"}, []string{"result"}),
		PluginOperations:            prometheus.NewCounterVec(prometheus.CounterOpts{Name: "noop_plugin_ops"}, []string{"operation", "status"}),
		PluginOperationDuration:     prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "noop_plugin_dur"}, []string{"operation"}),
		PluginListResultsCount:      prometheus.NewHistogram(prometheus.HistogramOpts{Name: "noop_plugin_list"}),
		CacheInvalidationsTotal:     prometheus.NewCounterVec(prometheus.CounterOpts{Name: "noop_cache_inv"}, []string{"status"}),
		ExecutionCompletedTotal:     prometheus.NewCounterVec(prometheus.CounterOpts{Name: "noop_exec_completed"}, []string{"trigger"}),
		ExecutionFailedTotal:        prometheus.NewCounterVec(prometheus.CounterOpts{Name: "noop_exec_failed"}, []string{"trigger"}),
		ExecutionStartedTotal:       prometheus.NewCounterVec(prometheus.CounterOpts{Name: "noop_exec_started"}, []string{"trigger"}),
		ExecutionDuration:           prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "noop_exec_duration"}, []string{"trigger"}),
		ExecutionSteps:              prometheus.NewHistogram(prometheus.HistogramOpts{Name: "noop_exec_steps"}),
		ExecutionActive:             prometheus.NewGauge(prometheus.GaugeOpts{Name: "noop_exec_active"}),
		CheckpointDuration:          prometheus.NewHistogram(prometheus.HistogramOpts{Name: "noop_checkpoint_dur"}),
		DispatchTotal:               prometheus.NewCounterVec(prometheus.CounterOpts{Name: "noop_dispatch"}, []string{"type", "result"}),
		CASRetriesTotal:             prometheus.NewCounter(prometheus.CounterOpts{Name: "noop_cas_retries"}),
		TokenRejectionsTotal:        prometheus.NewCounter(prometheus.CounterOpts{Name: "noop_token_rej"}),
		ShutdownStopsTotal:          prometheus.NewCounter(prometheus.CounterOpts{Name: "noop_shutdown_stops"}),
	}
}

// StaticDefinitionLoader returns a fixed definition for all calls.
type StaticDefinitionLoader struct {
	Definition *defEntities.WorkflowDefinition
}

var _ runtimePorts.DefinitionLoaderPort = (*StaticDefinitionLoader)(nil)

func (l *StaticDefinitionLoader) GetDefinition(_ context.Context, _ string, _ *model.ObjectId) (*defEntities.WorkflowDefinition, error) {
	if l.Definition == nil {
		return nil, fmt.Errorf("definition not found")
	}
	return l.Definition, nil
}

/**
 * StaticInstanceLoader
 * Returns a fixed instance for all calls.
 */
type StaticInstanceLoader struct {
	Instance *instanceEntities.WorkflowInstance
}

var _ instancePorts.InstanceLoaderPort = (*StaticInstanceLoader)(nil)

func (l *StaticInstanceLoader) GetInstance(_ context.Context, _ string) (*instanceEntities.WorkflowInstance, error) {
	if l.Instance == nil {
		return nil, fmt.Errorf("instance not found")
	}
	return l.Instance, nil
}

func (l *StaticInstanceLoader) Invalidate(_ context.Context, _ string) error {
	return nil
}

/**
 * NoopVaultService
 * Returns empty credentials for all calls.
 */
type NoopVaultService struct{}

var _ runtimePorts.VaultPort = (*NoopVaultService)(nil)

func (s *NoopVaultService) DecryptCredential(_ context.Context, _ string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

/**
 * NoopPluginRepo
 * Returns nil for all plugin lookups.
 */
type NoopPluginRepo struct{}

var _ pluginPorts.PluginManifestRepository = (*NoopPluginRepo)(nil)

func (r *NoopPluginRepo) Create(_ context.Context, _ *pluginEntities.PluginManifest) (*pluginEntities.PluginManifest, error) {
	return nil, nil
}
func (r *NoopPluginRepo) FindById(_ context.Context, _ *string) (*pluginEntities.PluginManifest, error) {
	return nil, nil
}
func (r *NoopPluginRepo) FindByPluginId(_ context.Context, _ string) (*pluginEntities.PluginManifest, error) {
	return nil, fmt.Errorf("plugin not found")
}
func (r *NoopPluginRepo) FindByIdAndUpdate(_ context.Context, _ *string, _ map[string]any) (*pluginEntities.PluginManifest, error) {
	return nil, nil
}
func (r *NoopPluginRepo) DeleteById(_ context.Context, _ *string) error {
	return nil
}
func (r *NoopPluginRepo) FindWithFilters(_ context.Context, _ model.Map, _ *model.PaginationOpts, _ model.Map) (*model.PaginatedResult[pluginEntities.PluginManifest], error) {
	return nil, nil
}
func (r *NoopPluginRepo) CountDocuments(_ context.Context, _ model.Map) (int64, error) {
	return 0, nil
}

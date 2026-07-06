package services

import (
	"context"
	"time"

	"assets/src/modules/ota/application/ports"
	"assets/src/modules/ota/domain/entities"

	otaEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/events"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// Inline port mocks shared by the package tests. Each mock records calls and
// returns configurable results so table tests stay declarative.

// newTestObjectID keeps test tables terse.
func newTestObjectID() model.ObjectId { return model.NewObjectID() }

type mockFirmwareRepo struct {
	byID    map[string]*entities.Firmware
	created []*entities.Firmware
	updates []map[string]any
	findErr error
}

func (m *mockFirmwareRepo) Create(_ context.Context, f *entities.Firmware) (*entities.Firmware, error) {
	m.created = append(m.created, f)
	return f, nil
}

func (m *mockFirmwareRepo) FindById(_ context.Context, id *string) (*entities.Firmware, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.byID[*id], nil
}

func (m *mockFirmwareRepo) FindByIdAndUpdate(_ context.Context, id *string, payload map[string]any) (*entities.Firmware, error) {
	m.updates = append(m.updates, payload)
	return m.byID[*id], nil
}

func (m *mockFirmwareRepo) DeleteById(_ context.Context, _ *string) error { return nil }

func (m *mockFirmwareRepo) FindWithFilters(_ context.Context, _ model.Map, _ *model.PaginationOpts, _ model.Map) (*model.PaginatedResult[entities.Firmware], error) {
	return &model.PaginatedResult[entities.Firmware]{}, nil
}

type mockPlanRepo struct {
	byID       map[string]*entities.OTAPlan
	created    []*entities.OTAPlan
	updates    []map[string]any
	found      []entities.OTAPlan
	increments []string
}

func (m *mockPlanRepo) Create(_ context.Context, p *entities.OTAPlan) (*entities.OTAPlan, error) {
	m.created = append(m.created, p)
	return p, nil
}

func (m *mockPlanRepo) FindById(_ context.Context, id *string) (*entities.OTAPlan, error) {
	return m.byID[*id], nil
}

func (m *mockPlanRepo) FindByIdAndUpdate(_ context.Context, id *string, payload map[string]any) (*entities.OTAPlan, error) {
	m.updates = append(m.updates, payload)
	return m.byID[*id], nil
}

func (m *mockPlanRepo) DeleteById(_ context.Context, _ *string) error { return nil }

func (m *mockPlanRepo) FindWithFilters(_ context.Context, _ model.Map, _ *model.PaginationOpts, _ model.Map) (*model.PaginatedResult[entities.OTAPlan], error) {
	return &model.PaginatedResult[entities.OTAPlan]{Items: m.found}, nil
}

func (m *mockPlanRepo) IncrementCounter(_ context.Context, _ string, field string, _ int) error {
	m.increments = append(m.increments, field)
	return nil
}

type mockExecutionRepo struct {
	byID        map[string]*entities.OTAExecution
	bulk        [][]*entities.OTAExecution
	updates     []map[string]any
	byPlan      []entities.OTAExecution
	byFilters   []entities.OTAExecution
	updatedMany []model.Map
}

func (m *mockExecutionRepo) BulkInsert(_ context.Context, execs []*entities.OTAExecution) (int64, error) {
	m.bulk = append(m.bulk, execs)
	return int64(len(execs)), nil
}

func (m *mockExecutionRepo) FindById(_ context.Context, id *string) (*entities.OTAExecution, error) {
	return m.byID[*id], nil
}

func (m *mockExecutionRepo) FindByIdAndUpdate(_ context.Context, id *string, payload map[string]any) (*entities.OTAExecution, error) {
	m.updates = append(m.updates, payload)
	return m.byID[*id], nil
}

func (m *mockExecutionRepo) FindByPlan(_ context.Context, _ *string, _ model.Map, _ *model.PaginationOpts) (*model.PaginatedResult[entities.OTAExecution], error) {
	return &model.PaginatedResult[entities.OTAExecution]{Items: m.byPlan}, nil
}

func (m *mockExecutionRepo) FindWithFilters(_ context.Context, _ model.Map, _ *model.PaginationOpts, _ model.Map) (*model.PaginatedResult[entities.OTAExecution], error) {
	return &model.PaginatedResult[entities.OTAExecution]{Items: m.byFilters}, nil
}

func (m *mockExecutionRepo) CountByState(_ context.Context, _ *string, _ entities.ExecutionState) (int64, error) {
	return 0, nil
}

func (m *mockExecutionRepo) UpdateMany(_ context.Context, filter model.Map, _ model.Map) (int64, error) {
	m.updatedMany = append(m.updatedMany, filter)
	return 1, nil
}

type mockStore struct {
	putURL   string
	getURLs  []string
	getCalls int
	statSize int64
	statSum  string
	statErr  error
	deleted  []string
}

func (m *mockStore) PresignPut(_ context.Context, _, _ string, _ time.Duration) (string, error) {
	return m.putURL, nil
}

func (m *mockStore) PresignGet(_ context.Context, _ string, _ time.Duration) (string, error) {
	url := m.putURL
	if m.getCalls < len(m.getURLs) {
		url = m.getURLs[m.getCalls]
	}
	m.getCalls++
	return url, nil
}

func (m *mockStore) Stat(_ context.Context, _ string) (int64, string, error) {
	return m.statSize, m.statSum, m.statErr
}

func (m *mockStore) Delete(_ context.Context, key string) error {
	m.deleted = append(m.deleted, key)
	return nil
}

type mockFirmwareScheduler struct {
	scheduled []string
	purged    []string
}

func (m *mockFirmwareScheduler) ScheduleAbandonCheck(id string, _ time.Time) error {
	m.scheduled = append(m.scheduled, id)
	return nil
}

func (m *mockFirmwareScheduler) PurgeAbandonCheck(id string) error {
	m.purged = append(m.purged, id)
	return nil
}

type mockOTAScheduler struct {
	starts []string
	closes []string
	scans  int
}

func (m *mockOTAScheduler) ScheduleStart(planID string, _ time.Time) error {
	m.starts = append(m.starts, planID)
	return nil
}

func (m *mockOTAScheduler) ScheduleClose(planID string, _ time.Time) error {
	m.closes = append(m.closes, planID)
	return nil
}

func (m *mockOTAScheduler) ScheduleScan(_ time.Time) error {
	m.scans++
	return nil
}

type mockPresence struct {
	online map[string]bool
}

func (m *mockPresence) IsOnline(_ context.Context, _, assetUUID string) (bool, error) {
	return m.online[assetUUID], nil
}

type mockEdge struct {
	dispatched []ports.OTACommandPayload
}

func (m *mockEdge) Dispatch(_ context.Context, _, _, _ string, payload ports.OTACommandPayload) error {
	m.dispatched = append(m.dispatched, payload)
	return nil
}

type mockAssetReader struct {
	infoByID map[string]*ports.OTAAssetInfo
}

func (m *mockAssetReader) GetAssetInfo(_ context.Context, assetID string) (*ports.OTAAssetInfo, error) {
	return m.infoByID[assetID], nil
}

type mockLiveState struct {
	set []*entities.OTAExecution
}

func (m *mockLiveState) SetExecutionLive(_ context.Context, exec *entities.OTAExecution) error {
	m.set = append(m.set, exec)
	return nil
}

type mockHistory struct {
	published []otaEvents.OTAStatusAdvisory
}

func (m *mockHistory) Publish(_ context.Context, adv otaEvents.OTAStatusAdvisory) error {
	m.published = append(m.published, adv)
	return nil
}

type mockTemplateSwitcher struct {
	switched [][2]string
}

func (m *mockTemplateSwitcher) SwitchTemplate(_ context.Context, assetID, targetTemplateID string) error {
	m.switched = append(m.switched, [2]string{assetID, targetTemplateID})
	return nil
}

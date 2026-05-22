package services

import (
	"context"
	"testing"
	"time"

	"mapexVault/src/modules/credentials/application/constants"
	"mapexVault/src/modules/credentials/application/di"
	"mapexVault/src/modules/credentials/domain/entities"
	"mapexVault/src/modules/credentials/domain/repositories"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/**
 * Mock Types
 */

type mockScheduleManager struct {
	publishedSchedules []natsModel.ScheduledPublishConfig
	purgedSubjects     []string
}

func (m *mockScheduleManager) PublishScheduled(config natsModel.ScheduledPublishConfig) error {
	m.publishedSchedules = append(m.publishedSchedules, config)
	return nil
}

func (m *mockScheduleManager) PurgeStreamSubject(stream, subject string) error {
	m.purgedSubjects = append(m.purgedSubjects, subject)
	return nil
}

func (m *mockScheduleManager) HasPendingMessages(stream, subject string) (bool, error) {
	return false, nil
}

type mockCredentialRepo struct {
	credentials []entities.Credential
	err         error
}

func (m *mockCredentialRepo) Create(_ context.Context, _ *entities.Credential) (*entities.Credential, error) {
	return nil, nil
}
func (m *mockCredentialRepo) FindById(_ context.Context, _ *string) (*entities.Credential, error) {
	return nil, nil
}
func (m *mockCredentialRepo) FindByIdAndUpdate(_ context.Context, _ *string, _ map[string]any) (*entities.Credential, error) {
	return nil, nil
}
func (m *mockCredentialRepo) DeleteById(_ context.Context, _ *string) error { return nil }
func (m *mockCredentialRepo) FindWithFilters(_ context.Context, _ model.Map, _ *model.PaginationOpts, _ model.Map) (*model.PaginatedResult[entities.Credential], error) {
	return nil, nil
}
func (m *mockCredentialRepo) FindActiveWithTokenExpiry(_ context.Context) ([]entities.Credential, error) {
	return m.credentials, m.err
}
func (m *mockCredentialRepo) CountDocuments(_ context.Context, _ model.Map) (int64, error) {
	return 0, nil
}

type mockPublisher struct{}

func (m *mockPublisher) Publish(_ natsModel.PublishConfig) error { return nil }

/**
 * publishRefreshSchedule Tests
 */

func newTestService(sm natsModel.ScheduleManager) *CredentialService {
	return &CredentialService{
		deps: di.CredentialServiceDependenciesInjection{
			ScheduleManager: sm,
			Publisher:       &mockPublisher{},
		},
	}
}

// makeServiceForSeed builds a CredentialService wired with the given schedule
// manager and credential repository, for bootstrap-seed tests.
func makeServiceForSeed(sm natsModel.ScheduleManager, repo repositories.CredentialRepository) *CredentialService {
	return &CredentialService{
		deps: di.CredentialServiceDependenciesInjection{
			ScheduleManager: sm,
			CredentialRepo:  repo,
			Publisher:       &mockPublisher{},
		},
	}
}

func TestPublishRefreshSchedule_PublishesWhenFuture(t *testing.T) {
	mock := &mockScheduleManager{}
	svc := newTestService(mock)

	expiresAt := time.Now().Add(1 * time.Hour)
	svc.publishRefreshSchedule("cred123", entities.CredentialOAuth2, &expiresAt)

	if len(mock.publishedSchedules) != 1 {
		t.Fatalf("expected 1 published schedule, got %d", len(mock.publishedSchedules))
	}

	published := mock.publishedSchedules[0]
	expectedScheduleAt := expiresAt.Add(-time.Duration(constants.RefreshBufferMinutes) * time.Minute)
	diff := published.ScheduleAt.Sub(expectedScheduleAt)
	if diff < -time.Second || diff > time.Second {
		t.Fatalf("expected scheduleAt ~%v, got %v", expectedScheduleAt, published.ScheduleAt)
	}

	if published.TargetSubject != constants.VaultScheduleFiredSubject {
		t.Fatalf("expected target %s, got %s", constants.VaultScheduleFiredSubject, published.TargetSubject)
	}
}

func TestPublishRefreshSchedule_SkipsWhenNilExpiry(t *testing.T) {
	mock := &mockScheduleManager{}
	svc := newTestService(mock)

	svc.publishRefreshSchedule("cred123", entities.CredentialOAuth2, nil)

	if len(mock.publishedSchedules) != 0 {
		t.Fatalf("expected 0 published schedules, got %d", len(mock.publishedSchedules))
	}
}

func TestPublishRefreshSchedule_SkipsWhenAlreadyExpired(t *testing.T) {
	mock := &mockScheduleManager{}
	svc := newTestService(mock)

	expiresAt := time.Now().Add(-5 * time.Minute)
	svc.publishRefreshSchedule("cred123", entities.CredentialOAuth2, &expiresAt)

	if len(mock.publishedSchedules) != 0 {
		t.Fatalf("expected 0 published schedules for expired token, got %d", len(mock.publishedSchedules))
	}
}

func TestPublishRefreshSchedule_PurgesBeforePublishing(t *testing.T) {
	mock := &mockScheduleManager{}
	svc := newTestService(mock)

	expiresAt := time.Now().Add(1 * time.Hour)
	svc.publishRefreshSchedule("cred123", entities.CredentialOAuth2, &expiresAt)

	if len(mock.purgedSubjects) != 1 {
		t.Fatalf("expected 1 purge call, got %d", len(mock.purgedSubjects))
	}

	expectedSubject := constants.VaultScheduleSubjectPrefix + ".cred123"
	if mock.purgedSubjects[0] != expectedSubject {
		t.Fatalf("expected purge subject %s, got %s", expectedSubject, mock.purgedSubjects[0])
	}
}

/**
 * BootstrapSeed Tests
 */

func TestBootstrapSeed_FutureCredentialScheduledNormally(t *testing.T) {
	mock := &mockScheduleManager{}
	expiresAt := time.Now().Add(2 * time.Hour)
	repo := &mockCredentialRepo{
		credentials: []entities.Credential{
			{
				ID:             model.NewObjectID(),
				Type:           entities.CredentialOAuth2,
				Status:         entities.CredentialStatusActive,
				TokenExpiresAt: &expiresAt,
			},
		},
	}

	makeServiceForSeed(mock, repo).bootstrapSeed()

	if len(mock.publishedSchedules) != 1 {
		t.Fatalf("expected 1 schedule, got %d", len(mock.publishedSchedules))
	}

	expectedAt := expiresAt.Add(-time.Duration(constants.RefreshBufferMinutes) * time.Minute)
	diff := mock.publishedSchedules[0].ScheduleAt.Sub(expectedAt)
	if diff < -time.Second || diff > time.Second {
		t.Fatalf("expected scheduleAt ~%v, got %v", expectedAt, mock.publishedSchedules[0].ScheduleAt)
	}
}

func TestBootstrapSeed_ExpiredCredentialScheduledImmediately(t *testing.T) {
	mock := &mockScheduleManager{}
	expiresAt := time.Now().Add(-30 * time.Minute)
	repo := &mockCredentialRepo{
		credentials: []entities.Credential{
			{
				ID:             model.NewObjectID(),
				Type:           entities.CredentialUserAndPass,
				Status:         entities.CredentialStatusActive,
				TokenExpiresAt: &expiresAt,
			},
		},
	}

	now := time.Now()
	makeServiceForSeed(mock, repo).bootstrapSeed()

	if len(mock.publishedSchedules) != 1 {
		t.Fatalf("expected 1 schedule, got %d", len(mock.publishedSchedules))
	}

	// Should be ~now + 30s
	expectedAt := now.Add(30 * time.Second)
	diff := mock.publishedSchedules[0].ScheduleAt.Sub(expectedAt)
	if diff < -2*time.Second || diff > 2*time.Second {
		t.Fatalf("expected scheduleAt ~%v, got %v", expectedAt, mock.publishedSchedules[0].ScheduleAt)
	}
}

func TestBootstrapSeed_EmptyResultNoSchedules(t *testing.T) {
	mock := &mockScheduleManager{}
	repo := &mockCredentialRepo{
		credentials: []entities.Credential{},
	}

	makeServiceForSeed(mock, repo).bootstrapSeed()

	if len(mock.publishedSchedules) != 0 {
		t.Fatalf("expected 0 schedules for empty result, got %d", len(mock.publishedSchedules))
	}
}

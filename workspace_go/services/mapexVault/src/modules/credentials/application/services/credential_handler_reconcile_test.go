package services

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"mapexVault/src/modules/credentials/application/constants"
	"mapexVault/src/modules/credentials/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// TestMain initializes config so RunReconcile can read vault_reconcile_interval
// without panicking. Tests rely on the default value when the env is not set.
func TestMain(m *testing.M) {
	config.InitConfig([]config.ConfigDefinition{
		{Key: "vault_reconcile_interval", Type: "int", Default: 3600},
	})
	os.Exit(m.Run())
}

/**
 * Reconcile Mock Types
 */

// mockReconcileScheduleManager extends mockScheduleManager with per-subject
// pending lookup so reconcile tests can assert reseed vs skip behavior.
type mockReconcileScheduleManager struct {
	publishedSchedules []natsModel.ScheduledPublishConfig
	purgedSubjects     []string

	// pendingBySubject controls HasPendingMessages per subject.
	// Defaults to false (i.e., no pending timer) when the key is absent.
	pendingBySubject map[string]bool
	hasPendingErr    error
}

func (m *mockReconcileScheduleManager) PublishScheduled(config natsModel.ScheduledPublishConfig) error {
	m.publishedSchedules = append(m.publishedSchedules, config)
	return nil
}

func (m *mockReconcileScheduleManager) PurgeStreamSubject(stream, subject string) error {
	m.purgedSubjects = append(m.purgedSubjects, subject)
	return nil
}

func (m *mockReconcileScheduleManager) HasPendingMessages(stream, subject string) (bool, error) {
	if m.hasPendingErr != nil {
		return false, m.hasPendingErr
	}
	return m.pendingBySubject[subject], nil
}

/**
 * RunReconcile Tests
 */

func TestRunReconcile_ReseedsCredentialsWithMissingTimers(t *testing.T) {
	sm := &mockReconcileScheduleManager{pendingBySubject: map[string]bool{}}
	expiresAt := time.Now().Add(30 * time.Minute)
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

	svc := makeServiceForSeed(sm, repo)
	svc.RunReconcile(context.Background())

	// One reseed (the missing per-credential refresh timer); RunReconcile schedules nothing else.
	if len(sm.publishedSchedules) != 1 {
		t.Fatalf("expected 1 published schedule (reseed only), got %d", len(sm.publishedSchedules))
	}
}

func TestRunReconcile_SkipsCredentialsWithPendingTimer(t *testing.T) {
	credId := model.NewObjectID()
	expectedSubject := constants.VaultScheduleSubjectPrefix + "." + credId.Hex()

	sm := &mockReconcileScheduleManager{
		pendingBySubject: map[string]bool{
			expectedSubject: true, // timer already exists
		},
	}
	expiresAt := time.Now().Add(30 * time.Minute)
	repo := &mockCredentialRepo{
		credentials: []entities.Credential{
			{
				ID:             credId,
				Type:           entities.CredentialOAuth2,
				Status:         entities.CredentialStatusActive,
				TokenExpiresAt: &expiresAt,
			},
		},
	}

	svc := makeServiceForSeed(sm, repo)
	svc.RunReconcile(context.Background())

	// The timer already exists, so nothing is reseeded and RunReconcile schedules nothing.
	if len(sm.publishedSchedules) != 0 {
		t.Fatalf("expected 0 published schedules (pending timer, no reseed), got %d", len(sm.publishedSchedules))
	}
}

func TestRunReconcile_HandlesRepositoryError(t *testing.T) {
	sm := &mockReconcileScheduleManager{pendingBySubject: map[string]bool{}}
	repo := &mockCredentialRepo{err: errors.New("mongo down")}

	svc := makeServiceForSeed(sm, repo)
	svc.RunReconcile(context.Background())

	// On a repo error nothing is reseeded and nothing is scheduled.
	if len(sm.publishedSchedules) != 0 {
		t.Fatalf("expected 0 published schedules on repo error, got %d", len(sm.publishedSchedules))
	}
}

func TestRunReconcile_SkipsCredentialWithNilExpiry(t *testing.T) {
	sm := &mockReconcileScheduleManager{pendingBySubject: map[string]bool{}}
	repo := &mockCredentialRepo{
		credentials: []entities.Credential{
			{
				ID:             model.NewObjectID(),
				Type:           entities.CredentialOAuth2,
				Status:         entities.CredentialStatusActive,
				TokenExpiresAt: nil,
			},
		},
	}

	svc := makeServiceForSeed(sm, repo)
	svc.RunReconcile(context.Background())

	// Nil expiry is skipped and RunReconcile schedules nothing.
	if len(sm.publishedSchedules) != 0 {
		t.Fatalf("expected 0 published schedules (nil expiry skipped), got %d", len(sm.publishedSchedules))
	}
}

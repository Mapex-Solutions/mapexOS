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

	// 1 reseed (refresh schedule) + 1 reconcile timer = 2 published
	if len(sm.publishedSchedules) != 2 {
		t.Fatalf("expected 2 published schedules (reseed + next reconcile), got %d", len(sm.publishedSchedules))
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

	// Only the next-reconcile timer should have been published (no reseed).
	if len(sm.publishedSchedules) != 1 {
		t.Fatalf("expected 1 published schedule (next reconcile only), got %d", len(sm.publishedSchedules))
	}
	if sm.publishedSchedules[0].Subject != constants.VaultReconcileScheduleSubject {
		t.Fatalf("expected subject %s, got %s", constants.VaultReconcileScheduleSubject, sm.publishedSchedules[0].Subject)
	}
}

func TestRunReconcile_HandlesRepositoryError(t *testing.T) {
	sm := &mockReconcileScheduleManager{pendingBySubject: map[string]bool{}}
	repo := &mockCredentialRepo{err: errors.New("mongo down")}

	svc := makeServiceForSeed(sm, repo)
	svc.RunReconcile(context.Background())

	// Repo error still must re-arm the next timer to keep the loop alive.
	if len(sm.publishedSchedules) != 1 {
		t.Fatalf("expected 1 published schedule (next reconcile despite repo error), got %d", len(sm.publishedSchedules))
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

	// No reseed (nil expiry) + 1 reconcile timer = 1 published
	if len(sm.publishedSchedules) != 1 {
		t.Fatalf("expected 1 published schedule (next reconcile only), got %d", len(sm.publishedSchedules))
	}
}

/**
 * scheduleNextReconcile Tests
 */

func TestScheduleNextReconcile_SkipsWhenAlreadyPending(t *testing.T) {
	sm := &mockReconcileScheduleManager{
		pendingBySubject: map[string]bool{
			constants.VaultReconcileScheduleSubject: true,
		},
	}

	svc := makeServiceForSeed(sm, &mockCredentialRepo{})
	svc.scheduleNextReconcile()

	if len(sm.publishedSchedules) != 0 {
		t.Fatalf("expected 0 published schedules when already pending, got %d", len(sm.publishedSchedules))
	}
}

func TestScheduleNextReconcile_PublishesWithCorrectConfig(t *testing.T) {
	sm := &mockReconcileScheduleManager{pendingBySubject: map[string]bool{}}

	svc := makeServiceForSeed(sm, &mockCredentialRepo{})
	svc.scheduleNextReconcile()

	if len(sm.publishedSchedules) != 1 {
		t.Fatalf("expected 1 published schedule, got %d", len(sm.publishedSchedules))
	}

	published := sm.publishedSchedules[0]
	if published.Subject != constants.VaultReconcileScheduleSubject {
		t.Fatalf("expected subject %s, got %s", constants.VaultReconcileScheduleSubject, published.Subject)
	}
	if published.TargetSubject != constants.VaultReconcileFiredSubject {
		t.Fatalf("expected target %s, got %s", constants.VaultReconcileFiredSubject, published.TargetSubject)
	}
	if published.MsgId != constants.VaultReconcileMsgId {
		t.Fatalf("expected MsgId %s, got %s", constants.VaultReconcileMsgId, published.MsgId)
	}

	expectedAt := time.Now().Add(time.Duration(constants.VaultReconcileDefaultIntervalSeconds) * time.Second)
	diff := published.ScheduleAt.Sub(expectedAt)
	if diff < -2*time.Second || diff > 2*time.Second {
		t.Fatalf("expected scheduleAt ~%v, got %v", expectedAt, published.ScheduleAt)
	}
}

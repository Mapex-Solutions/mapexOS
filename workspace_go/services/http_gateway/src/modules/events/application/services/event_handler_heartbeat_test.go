package services

import (
	"context"
	"errors"
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	"http_gateway/src/bootstrap"
	dsDto "http_gateway/src/modules/datasources/application/dtos"
	"http_gateway/src/modules/events/application/di"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// newTestService builds an EventService with a minimal Metrics struct using
// fresh Prometheus collectors so calls to HeartbeatsTotal/HeartbeatDuration
// don't panic. NatsBus is left nil — these tests only exercise validation
// paths that fail BEFORE any publish.
func newTestService() *EventService {
	hbTotal := prometheus.NewCounterVec(prometheus.CounterOpts{Name: "test_heartbeats_total"}, []string{"status"})
	hbDuration := prometheus.NewHistogram(prometheus.HistogramOpts{Name: "test_heartbeat_duration_seconds"})
	return &EventService{
		deps: di.EventServiceDependenciesInjection{
			Metrics: &bootstrap.HttpGatewayMetrics{
				HeartbeatsTotal:   hbTotal,
				HeartbeatDuration: hbDuration,
			},
		},
	}
}

// These tests exercise the validation paths of EventService.ProcessHeartbeat —
// paths that return an error BEFORE touching NatsBus, so no real broker is
// needed. Happy-path publish is exercised by the e2e scripts in
// scripts/heartbeat_e2e/ against a running stack.
//
// TKT-2026-0036 reformulation: the legacy AssetBind.Type='fixedAssetId'
// constraint is REMOVED — body now carries the assetUUID, so any DataSource
// shape works.

// stringPtr returns a pointer to the given string (test helper).
func stringPtr(s string) *string { return &s }

func TestProcessHeartbeat_NilDataSource(t *testing.T) {
	s := newTestService()
	err := s.ProcessHeartbeat(context.Background(), nil, "12345678")
	if err == nil {
		t.Fatalf("expected error when dataSource is nil; got nil")
	}
	var sc *customErrors.ServerCustomError
	if !errors.As(err, &sc) {
		t.Fatalf("expected *ServerCustomError; got %T", err)
	}
	if sc.Code != status.NOT_FOUND {
		t.Fatalf("expected NOT_FOUND code; got %d", sc.Code)
	}
}

func TestProcessHeartbeat_MissingOrgId(t *testing.T) {
	s := newTestService()
	ds := &dsDto.DataSourceResponse{}
	err := s.ProcessHeartbeat(context.Background(), ds, "12345678")
	if err == nil {
		t.Fatalf("expected error when OrgId is missing; got nil")
	}
	var sc *customErrors.ServerCustomError
	if !errors.As(err, &sc) {
		t.Fatalf("expected *ServerCustomError; got %T", err)
	}
	if sc.Code != status.NOT_FOUND {
		t.Fatalf("expected NOT_FOUND code; got %d", sc.Code)
	}
}

func TestProcessHeartbeat_EmptyAssetUUID(t *testing.T) {
	s := newTestService()

	orgId, err := model.ToObjectID("0000000000000000000aa001")
	if err != nil {
		t.Fatalf("could not build ObjectID: %v", err)
	}

	ds := &dsDto.DataSourceResponse{
		OrgId:   &orgId,
		PathKey: stringPtr("000001"),
	}

	hbErr := s.ProcessHeartbeat(context.Background(), ds, "")
	if hbErr == nil {
		t.Fatalf("expected error when assetUUID is empty; got nil")
	}
	var sc *customErrors.ServerCustomError
	if !errors.As(hbErr, &sc) {
		t.Fatalf("expected *ServerCustomError; got %T", hbErr)
	}
	if sc.Code != status.UNPROCESSABLE_ENTITY {
		t.Fatalf("expected UNPROCESSABLE_ENTITY code; got %d", sc.Code)
	}
}

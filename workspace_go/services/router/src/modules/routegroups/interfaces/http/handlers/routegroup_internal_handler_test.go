package handlers

import (
	"errors"
	"testing"

	"net/http/httptest"

	"router/src/modules/routegroups/application/dtos"
	"router/src/modules/routegroups/application/services/mocks"

	"github.com/gofiber/fiber/v2"
)

/** GetRouteGroupsByIds */

func TestGetRouteGroupsByIds_Success(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.GetByIdsResponse = []dtos.RouteGroupResponse{
		{Version: ptrString("1.0.0"), Name: ptrString("Group 1")},
		{Version: ptrString("2.0.0"), Name: ptrString("Group 2")},
	}

	app := fiber.New()
	app.Get("/test", injectLocals(map[string]interface{}{
		"queryDTO": &dtos.RouteGroupInternalIdsQuery{
			Ids: "507f1f77bcf86cd799439011,507f1f77bcf86cd799439012",
		},
	}), GetRouteGroupsByIds(svc))

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	r := parseResponseBody(t, resp)
	if r.Status != 200 {
		t.Errorf("expected body status 200, got %d", r.Status)
	}
	if r.Data == nil {
		t.Error("expected data to be non-nil")
	}

	if len(svc.GetByIdsCalls) != 1 {
		t.Fatalf("expected 1 GetByIds call, got %d", len(svc.GetByIdsCalls))
	}
	if len(svc.GetByIdsCalls[0]) != 2 {
		t.Errorf("expected 2 IDs passed to service, got %d", len(svc.GetByIdsCalls[0]))
	}
}

func TestGetRouteGroupsByIds_EmptyIds(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()

	app := fiber.New()
	app.Get("/test", injectLocals(map[string]interface{}{
		"queryDTO": &dtos.RouteGroupInternalIdsQuery{
			Ids: "",
		},
	}), GetRouteGroupsByIds(svc))

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 500 {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}

	r := parseResponseBody(t, resp)
	if r.Status != 500 {
		t.Errorf("expected body status 500, got %d", r.Status)
	}
	if len(r.Errors) == 0 || r.Errors[0] != "ids parameter is required" {
		t.Errorf("expected 'ids parameter is required' error, got %v", r.Errors)
	}
}

func TestGetRouteGroupsByIds_ServiceError(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.GetByIdsError = errors.New("service failure")

	app := fiber.New()
	app.Get("/test", injectLocals(map[string]interface{}{
		"queryDTO": &dtos.RouteGroupInternalIdsQuery{
			Ids: "507f1f77bcf86cd799439011",
		},
	}), GetRouteGroupsByIds(svc))

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 500 {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}

	r := parseResponseBody(t, resp)
	if len(r.Errors) == 0 || r.Errors[0] != "service failure" {
		t.Errorf("expected 'service failure' error, got %v", r.Errors)
	}
}

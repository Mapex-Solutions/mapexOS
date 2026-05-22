package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"router/src/modules/routegroups/application/dtos"
	"router/src/modules/routegroups/application/services/mocks"

	"github.com/gofiber/fiber/v2"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

func TestMain(m *testing.M) {
	logger.InitLogger(logger.LoggerOptions{
		ServiceName: "test",
		Environment: "test",
		Level:       logger.ErrorLevel,
	})
	os.Exit(m.Run())
}

// parseResponseBody unmarshals Fiber response into response.Response struct.
func parseResponseBody(t *testing.T, resp *http.Response) response.Response {
	t.Helper()
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	var r response.Response
	if err := json.Unmarshal(body, &r); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	return r
}

// injectLocals creates a middleware that injects key-value pairs into c.Locals().
func injectLocals(locals map[string]interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		for k, v := range locals {
			c.Locals(k, v)
		}
		return c.Next()
	}
}

// createTestRequestContext creates a RequestContext for handler testing.
func createTestRequestContext() *reqCtx.RequestContext {
	orgId := "507f1f77bcf86cd799439022"
	return &reqCtx.RequestContext{
		ScopedOrgIds: []string{orgId},
		OrgContext:   &orgId,
		OrgContextData: &reqCtx.CoverageOrg{
			ID:      orgId,
			Name:    "Test Org",
			Type:    "customer",
			PathKey: "000001/0001",
		},
		UserId: "user-1",
	}
}

// ptrString returns a pointer to the given string.
func ptrString(s string) *string { return &s }

// ptrBool returns a pointer to the given bool.
func ptrBool(b bool) *bool { return &b }

/** CreateRouteGroup */

func TestCreateRouteGroup_Success(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.CreateResponse = &dtos.RouteGroupResponse{
		Version: ptrString("1.0.0"),
		Name:    ptrString("Test Group"),
		Enabled: ptrBool(true),
	}

	app := fiber.New()
	app.Post("/test", injectLocals(map[string]interface{}{
		"requestContext": createTestRequestContext(),
		"bodyDTO": &dtos.RouteGroupCreateDTO{
			Version: "1.0.0",
			Name:    "Test Group",
		},
	}), CreateRouteGroup(svc))

	req := httptest.NewRequest("POST", "/test", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 201 {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	r := parseResponseBody(t, resp)
	if r.Status != 201 {
		t.Errorf("expected body status 201, got %d", r.Status)
	}
	if r.Data == nil {
		t.Error("expected data to be non-nil")
	}
	if len(svc.CreateCalls) != 1 {
		t.Errorf("expected 1 create call, got %d", len(svc.CreateCalls))
	}
}

func TestCreateRouteGroup_ServiceError(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.CreateError = errors.New("service error")

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(response.Response{
				Status: 500,
				Errors: []string{err.Error()},
			})
		},
	})
	app.Post("/test", injectLocals(map[string]interface{}{
		"requestContext": createTestRequestContext(),
		"bodyDTO": &dtos.RouteGroupCreateDTO{
			Version: "1.0.0",
			Name:    "Test Group",
		},
	}), CreateRouteGroup(svc))

	req := httptest.NewRequest("POST", "/test", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 500 {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}

	r := parseResponseBody(t, resp)
	if len(r.Errors) == 0 {
		t.Error("expected errors in response")
	}
}

func TestCreateRouteGroup_MissingRequestContext(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()

	app := fiber.New()
	// No requestContext injected
	app.Post("/test", injectLocals(map[string]interface{}{
		"bodyDTO": &dtos.RouteGroupCreateDTO{
			Version: "1.0.0",
			Name:    "Test Group",
		},
	}), CreateRouteGroup(svc))

	req := httptest.NewRequest("POST", "/test", nil)
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
}

/** GetRouteGroupById */

func TestGetRouteGroupById_Success(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.GetByIdResponse = &dtos.RouteGroupResponse{
		Version: ptrString("1.0.0"),
		Name:    ptrString("Found Group"),
	}

	app := fiber.New()
	app.Get("/test", injectLocals(map[string]interface{}{
		"paramsDTO": &dtos.RouteGroupIdDTO{RouteGroupId: "507f1f77bcf86cd799439011"},
	}), GetRouteGroupById(svc))

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
}

func TestGetRouteGroupById_ServiceError(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.GetByIdError = errors.New("not found")

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(response.Response{
				Status: 500,
				Errors: []string{err.Error()},
			})
		},
	})
	app.Get("/test", injectLocals(map[string]interface{}{
		"paramsDTO": &dtos.RouteGroupIdDTO{RouteGroupId: "507f1f77bcf86cd799439011"},
	}), GetRouteGroupById(svc))

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 500 {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}

func TestGetRouteGroupById_CallsServiceWithCorrectId(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.GetByIdResponse = &dtos.RouteGroupResponse{}

	expectedId := "507f1f77bcf86cd799439099"

	app := fiber.New()
	app.Get("/test", injectLocals(map[string]interface{}{
		"paramsDTO": &dtos.RouteGroupIdDTO{RouteGroupId: expectedId},
	}), GetRouteGroupById(svc))

	req := httptest.NewRequest("GET", "/test", nil)
	_, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(svc.GetByIdCalls) != 1 {
		t.Fatalf("expected 1 GetById call, got %d", len(svc.GetByIdCalls))
	}
	if svc.GetByIdCalls[0] != expectedId {
		t.Errorf("expected service called with ID %q, got %q", expectedId, svc.GetByIdCalls[0])
	}
}

/** UpdateRouteGroupById */

func TestUpdateRouteGroupById_Success(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.UpdateResponse = &dtos.RouteGroupResponse{
		Version: ptrString("2.0.0"),
		Name:    ptrString("Updated Group"),
	}

	app := fiber.New()
	app.Put("/test", injectLocals(map[string]interface{}{
		"requestContext": createTestRequestContext(),
		"paramsDTO":      &dtos.RouteGroupIdDTO{RouteGroupId: "507f1f77bcf86cd799439011"},
		"bodyDTO": &dtos.RouteGroupUpdateDTO{
			Name: ptrString("Updated Group"),
		},
	}), UpdateRouteGroupById(svc))

	req := httptest.NewRequest("PUT", "/test", nil)
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
}

func TestUpdateRouteGroupById_ServiceError(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.UpdateError = errors.New("update failed")

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(response.Response{
				Status: 500,
				Errors: []string{err.Error()},
			})
		},
	})
	app.Put("/test", injectLocals(map[string]interface{}{
		"requestContext": createTestRequestContext(),
		"paramsDTO":      &dtos.RouteGroupIdDTO{RouteGroupId: "507f1f77bcf86cd799439011"},
		"bodyDTO": &dtos.RouteGroupUpdateDTO{
			Name: ptrString("Updated Group"),
		},
	}), UpdateRouteGroupById(svc))

	req := httptest.NewRequest("PUT", "/test", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 500 {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}

func TestUpdateRouteGroupById_OverridesOrgIdFromRequestContext(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.UpdateResponse = &dtos.RouteGroupResponse{}

	orgIdFromBody := "should-be-overridden"
	orgIdFromContext := "org-from-context"
	requestContext := &reqCtx.RequestContext{
		ScopedOrgIds: []string{orgIdFromContext},
		OrgContext:   &orgIdFromContext,
		OrgContextData: &reqCtx.CoverageOrg{
			ID:      orgIdFromContext,
			Name:    "Test Org",
			Type:    "customer",
			PathKey: "000001/0001",
		},
		UserId: "user-1",
	}

	app := fiber.New()
	app.Put("/test", injectLocals(map[string]interface{}{
		"requestContext": requestContext,
		"paramsDTO":      &dtos.RouteGroupIdDTO{RouteGroupId: "507f1f77bcf86cd799439011"},
		"bodyDTO": &dtos.RouteGroupUpdateDTO{
			Name:  ptrString("Updated"),
			OrgId: &orgIdFromBody,
		},
	}), UpdateRouteGroupById(svc))

	req := httptest.NewRequest("PUT", "/test", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	if len(svc.UpdateCalls) != 1 {
		t.Fatalf("expected 1 update call, got %d", len(svc.UpdateCalls))
	}

	// Verify OrgId was overridden with requestContext.OrgContext value
	updatedOrgId := svc.UpdateCalls[0].DTO.OrgId
	if updatedOrgId == nil {
		t.Fatal("expected OrgId to be non-nil")
	}
	if *updatedOrgId != orgIdFromContext {
		t.Errorf("expected OrgId %q, got %q", orgIdFromContext, *updatedOrgId)
	}
}

/** GetRouteGroups */

func TestGetRouteGroups_Success(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.GetListResponse = &model.PaginatedResult[dtos.RouteGroupResponse]{
		Items: []dtos.RouteGroupResponse{
			{Version: ptrString("1.0.0"), Name: ptrString("Group 1")},
		},
		Pagination: model.Pagination{
			Page:       1,
			PerPage:    20,
			TotalItems: 1,
			TotalPages: 1,
		},
	}

	app := fiber.New()
	app.Get("/test", injectLocals(map[string]interface{}{
		"requestContext": createTestRequestContext(),
		"queryDTO":       &dtos.RouteGroupQueryDTO{},
	}), GetRouteGroups(svc))

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
}

func TestGetRouteGroups_MissingRequestContext(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()

	app := fiber.New()
	// No requestContext injected
	app.Get("/test", injectLocals(map[string]interface{}{
		"queryDTO": &dtos.RouteGroupQueryDTO{},
	}), GetRouteGroups(svc))

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
}

/** DeleteRouteGroupById */

func TestDeleteRouteGroupById_Success(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.DeleteResponse = map[string]bool{"deleted": true}

	app := fiber.New()
	app.Delete("/test", injectLocals(map[string]interface{}{
		"paramsDTO": &dtos.RouteGroupIdDTO{RouteGroupId: "507f1f77bcf86cd799439011"},
	}), DeleteRouteGroupById(svc))

	req := httptest.NewRequest("DELETE", "/test", nil)
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
}

func TestDeleteRouteGroupById_ServiceError(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.DeleteError = errors.New("delete failed")

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(response.Response{
				Status: 500,
				Errors: []string{err.Error()},
			})
		},
	})
	app.Delete("/test", injectLocals(map[string]interface{}{
		"paramsDTO": &dtos.RouteGroupIdDTO{RouteGroupId: "507f1f77bcf86cd799439011"},
	}), DeleteRouteGroupById(svc))

	req := httptest.NewRequest("DELETE", "/test", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 500 {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}

/** GetRouteGroupCount */

func TestGetRouteGroupCount_Success(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.CountResponse = 42

	app := fiber.New()
	app.Get("/test", injectLocals(map[string]interface{}{
		"requestContext": createTestRequestContext(),
	}), GetRouteGroupCount(svc))

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

	// Verify count value in data
	dataMap, ok := r.Data.(map[string]interface{})
	if !ok {
		t.Fatal("expected data to be a map")
	}
	count, ok := dataMap["count"].(float64) // JSON numbers are float64
	if !ok {
		t.Fatal("expected count to be a number")
	}
	if int64(count) != 42 {
		t.Errorf("expected count 42, got %v", count)
	}

	if svc.CountCalls != 1 {
		t.Errorf("expected 1 count call, got %d", svc.CountCalls)
	}
}

func TestGetRouteGroupCount_ServiceError(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()
	svc.CountError = errors.New("count failed")

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(response.Response{
				Status: 500,
				Errors: []string{err.Error()},
			})
		},
	})
	app.Get("/test", injectLocals(map[string]interface{}{
		"requestContext": createTestRequestContext(),
	}), GetRouteGroupCount(svc))

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 500 {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}

func TestGetRouteGroupCount_MissingRequestContext(t *testing.T) {
	svc := mocks.NewMockRouteGroupService()

	app := fiber.New()
	// No requestContext injected
	app.Get("/test", GetRouteGroupCount(svc))

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
}

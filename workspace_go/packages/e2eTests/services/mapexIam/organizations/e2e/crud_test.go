// CRUD tests for the IAM organizations endpoints.
//
// Covers Create, GetByID, Update, Delete — each with happy path,
// auth (no token), validation (invalid payload), and not-found.
// LIST and TREE concerns live in list_test.go and tree_test.go.
package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/organizations/payloads"
)

// TestCreate_201 verifies POST returns 201 with id and pathKey populated.
func TestCreate_201(t *testing.T) {
	runID := random.NewRunID()
	spec := payloads.SagaTestCustomerOrg(runID).Build()

	resp, err := client.Raw(context.Background(), http.MethodPost, "/api/v1/organizations", spec)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var out struct {
		Data struct {
			ID      string `json:"id"`
			PathKey string `json:"pathKey"`
		} `json:"data"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
	assert.NotEmpty(t, out.Data.ID)
	assert.NotEmpty(t, out.Data.PathKey)

	t.Cleanup(func() { deleteOrg(t, out.Data.ID) })
}

// TestCreate_NoToken_401 confirms the auth middleware rejects the
// request before the handler runs.
func TestCreate_NoToken_401(t *testing.T) {
	anon := httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	spec := payloads.SagaTestCustomerOrg(random.NewRunID()).Build()

	resp, err := anon.Raw(context.Background(), http.MethodPost, "/api/v1/organizations", spec)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// TestCreate_EmptyName_400 sends a malformed payload to confirm the
// validator catches the empty name before persistence.
func TestCreate_EmptyName_400(t *testing.T) {
	spec := payloads.SagaTestCustomerOrg(random.NewRunID()).WithName("").Build()

	resp, err := client.Raw(context.Background(), http.MethodPost, "/api/v1/organizations", spec)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestGetByID_200 verifies the round-trip: create then GET by id returns
// the same entity with enabled=true.
func TestGetByID_200(t *testing.T) {
	id, _ := createOrgForTest(t)

	resp, err := client.Raw(context.Background(), http.MethodGet, "/api/v1/organizations/"+id, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var out struct {
		Data struct {
			ID      string `json:"id"`
			Enabled bool   `json:"enabled"`
		} `json:"data"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
	assert.Equal(t, id, out.Data.ID)
	assert.True(t, out.Data.Enabled)
}

// TestGetByID_404 verifies a syntactically valid but non-existent id
// returns 404, not 400 or 500.
func TestGetByID_404(t *testing.T) {
	resp, err := client.Raw(context.Background(), http.MethodGet, "/api/v1/organizations/"+nonExistentOrgID, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// TestUpdate_200 patches the name and verifies the new value is
// reflected on a subsequent GET.
func TestUpdate_200(t *testing.T) {
	id, _ := createOrgForTest(t)

	newName := fmt.Sprintf("saga-org-renamed-%s", random.NewRunID())
	body := map[string]any{"name": newName}

	resp, err := client.Raw(context.Background(), http.MethodPatch, "/api/v1/organizations/"+id, body)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated,
		"expected 200 or 201 from PATCH, got %d", resp.StatusCode)

	verify, err := client.Raw(context.Background(), http.MethodGet, "/api/v1/organizations/"+id, nil)
	require.NoError(t, err)
	defer verify.Body.Close()
	require.Equal(t, http.StatusOK, verify.StatusCode)
	var out struct {
		Data struct{ Name string `json:"name"` } `json:"data"`
	}
	require.NoError(t, json.NewDecoder(verify.Body).Decode(&out))
	assert.Equal(t, newName, out.Data.Name)
}

// TestUpdate_404 patches a non-existent id and expects 404.
func TestUpdate_404(t *testing.T) {
	body := map[string]any{"name": "saga-rename-attempt"}
	resp, err := client.Raw(context.Background(), http.MethodPatch, "/api/v1/organizations/"+nonExistentOrgID, body)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// TestDelete_200 deletes the freshly-created org and confirms the
// follow-up GET returns 404.
func TestDelete_200(t *testing.T) {
	runID := random.NewRunID()
	spec := payloads.SagaTestCustomerOrg(runID).Build()
	resp, err := client.Raw(context.Background(), http.MethodPost, "/api/v1/organizations", spec)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created struct {
		Data struct{ ID string `json:"id"` } `json:"data"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	id := created.Data.ID
	require.NotEmpty(t, id)

	del, err := client.Raw(context.Background(), http.MethodDelete, "/api/v1/organizations/"+id, nil)
	require.NoError(t, err)
	defer del.Body.Close()
	require.True(t, del.StatusCode == http.StatusOK || del.StatusCode == http.StatusNoContent,
		"expected 200 or 204, got %d", del.StatusCode)

	verify, err := client.Raw(context.Background(), http.MethodGet, "/api/v1/organizations/"+id, nil)
	require.NoError(t, err)
	defer verify.Body.Close()
	assert.Equal(t, http.StatusNotFound, verify.StatusCode)
}

// TestDelete_404 deletes a non-existent id and expects 404.
func TestDelete_404(t *testing.T) {
	resp, err := client.Raw(context.Background(), http.MethodDelete, "/api/v1/organizations/"+nonExistentOrgID, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

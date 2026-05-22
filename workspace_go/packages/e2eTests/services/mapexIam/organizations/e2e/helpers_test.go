// Helpers shared across the organizations module e2e tests.
//
// Each test file (crud_test.go, list_test.go, tree_test.go) holds only
// the Test* functions for its concern. Cross-cutting plumbing —
// TestMain, login, CRUD fixtures, list HTTP plumbing, tree HTTP
// plumbing — lives here so the test files stay focused on assertions.
package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/utils"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/mapexIam/organizations/payloads"
)

// client is the package-wide HTTP client. Lives across every Test*;
// carries the bearer + X-Org-Context attached once in TestMain. Tests
// build throwaway clients only when they need a different auth posture
// (no token, wrong token).
var client *httpclient.HTTPClient


// TestMain wires the package-wide client. The stack must be up
// (SetupE2EEnvironment hits /health) and the seed admin must be able to
// log in. If either fails, the package exits without running the tests
// rather than reporting cascade failures.
func TestMain(m *testing.M) {
	if err := utils.SetupE2EEnvironment(); err != nil {
		os.Stderr.WriteString("e2e environment not ready: " + err.Error() + "\n")
		os.Exit(0)
	}

	client = httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	token, err := loginAsSeedAdmin(client)
	if err != nil {
		os.Stderr.WriteString("seed admin login failed: " + err.Error() + "\n")
		os.Exit(1)
	}
	client.SetHeader("Authorization", "Bearer "+token)
	client.SetHeader("X-Org-Context", constants.MapexosOrgID)

	os.Exit(m.Run())
}

// loginAsSeedAdmin performs POST /auth/login with the canonical seed
// admin credentials and returns the bearer. Lives at the e2e helpers
// scope because TestMain needs it before any saga step has run; the
// saga equivalent lives in services/mapexIam/auth/steps/SeedAdminLogin.
func loginAsSeedAdmin(c *httpclient.HTTPClient) (string, error) {
	payload := map[string]any{
		"email":         constants.RootUserEmail,
		"password":      constants.RootUserPassword,
		"keepConnected": false,
	}
	resp, err := c.Raw(context.Background(), http.MethodPost, "/auth/login", payload)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("login: unexpected status %d", resp.StatusCode)
	}
	var out struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Data.AccessToken, nil
}


// createOrgForTest is the convenience helper for happy-path tests. It
// returns the new id and a cleanup closure (also auto-registered via
// t.Cleanup so callers may ignore the second return when convenient).
func createOrgForTest(t *testing.T) (string, func()) {
	t.Helper()
	spec := payloads.SagaTestCustomerOrg(random.NewRunID()).Build()
	resp, err := client.Raw(context.Background(), http.MethodPost, "/api/v1/organizations", spec)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode, "fixture create must succeed")

	var out struct {
		Data struct{ ID string `json:"id"` } `json:"data"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
	require.NotEmpty(t, out.Data.ID)

	id := out.Data.ID
	cleanup := func() { deleteOrg(t, id) }
	t.Cleanup(cleanup)
	return id, cleanup
}

// deleteOrg runs DELETE; tolerates 404 to keep cleanups idempotent when
// a previous step already removed the row.
func deleteOrg(t *testing.T, id string) {
	t.Helper()
	resp, err := client.Raw(context.Background(), http.MethodDelete, "/api/v1/organizations/"+id, nil)
	if err != nil {
		t.Logf("cleanup delete org %s failed: %v", id, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		t.Logf("cleanup delete org %s: unexpected status %d", id, resp.StatusCode)
	}
}


// setupListFixtures creates `count` orgs sharing the same runID prefix.
// Returns the runID and the slice of created ids; registers t.Cleanup
// to delete every fixture in reverse order.
//
// Naming pattern: "<orgNamePrefix>-<runID>-<i>" — same prefix
// "<orgNamePrefix>-<runID>" can be passed as ?name= to scope subsequent
// list queries to ONLY this test's universe.
func setupListFixtures(t *testing.T, count int) (string, []string) {
	t.Helper()
	runID := random.NewRunID()
	ids := make([]string, 0, count)

	for i := range count {
		name := orgName(runID, i+1)
		spec := payloads.SagaTestCustomerOrg(runID).WithName(name).Build()
		resp, err := client.Raw(context.Background(), http.MethodPost, "/api/v1/organizations", spec)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode, "fixture %d create must succeed", i+1)

		var out struct {
			Data struct{ ID string `json:"id"` } `json:"data"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
		_ = resp.Body.Close()
		require.NotEmpty(t, out.Data.ID)
		ids = append(ids, out.Data.ID)
	}

	t.Cleanup(func() {
		for i := len(ids) - 1; i >= 0; i-- {
			deleteOrg(t, ids[i])
		}
	})
	return runID, ids
}

// orgName is the deterministic naming function fixtures share with the
// list filters. Centralising it here keeps the prefix in one place.
func orgName(runID string, index int) string {
	return fmt.Sprintf("%s-%s-%d", orgNamePrefix, runID, index)
}


// listResponse mirrors the gokit PaginatedResult envelope wrapped in the
// service's StandardResponse.
type listResponse struct {
	Data struct {
		Items      []orgListItem  `json:"items"`
		Pagination paginationMeta `json:"pagination"`
	} `json:"data"`
}

type orgListItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type paginationMeta struct {
	Page       int64 `json:"page"`
	PerPage    int64 `json:"perPage"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int64 `json:"totalPages"`
}

// listPage is the unwrapped envelope tests assert against.
type listPage struct {
	Items      []orgListItem
	Pagination paginationMeta
}

// fetchPage executes GET /api/v1/organizations?<query> and decodes the
// envelope into a flat shape so tests stay terse.
func fetchPage(t *testing.T, query string) listPage {
	t.Helper()
	resp, err := client.Raw(context.Background(), http.MethodGet, "/api/v1/organizations?"+query, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "list query %q must return 200", query)

	var out listResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
	return listPage{Items: out.Data.Items, Pagination: out.Data.Pagination}
}

// listQuery composes the canonical query string for prefix-scoped tests.
// All list tests scope the universe to the runID prefix to eliminate
// interference from pre-seeded data and from parallel test runs.
func listQuery(runID string, page, perPage int) string {
	return fmt.Sprintf("page=%d&perPage=%d&name=%s-%s", page, perPage, orgNamePrefix, runID)
}

// walkPagesByPrefix iterates pages [from..to] (inclusive) in the given
// step direction with perPage=1 and returns every id seen. Used by
// Page1_To_Page15 (step=+1) and Page15_To_Page1 (step=-1).
func walkPagesByPrefix(t *testing.T, runID string, from, to, step int) []string {
	t.Helper()
	got := make([]string, 0, listFixtureCount)
	for p := from; (step > 0 && p <= to) || (step < 0 && p >= to); p += step {
		page := fetchPage(t, listQuery(runID, p, 1))
		require.Lenf(t, page.Items, 1, "page %d must return exactly one item with perPage=1", p)
		got = append(got, page.Items[0].ID)
	}
	return got
}

// idsOf extracts the id slice from a list page in the order the
// backend returned the items.
func idsOf(items []orgListItem) []string {
	out := make([]string, 0, len(items))
	for _, it := range items {
		out = append(out, it.ID)
	}
	return out
}


// treeResponse mirrors the /tree endpoint envelope.
type treeResponse struct {
	Data struct {
		Items  []treeItem `json:"items"`
		Cursor cursorInfo `json:"cursor"`
	} `json:"data"`
}

type treeItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type cursorInfo struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	HasNext  bool   `json:"hasNext"`
	HasPrev  bool   `json:"hasPrev"`
}

// treePage is the unwrapped /tree response tests assert against.
type treePage struct {
	Items  []treeItem
	Cursor cursorInfo
}

// fetchTreePage executes GET /api/v1/organizations/tree?<query> and
// decodes the envelope.
func fetchTreePage(t *testing.T, query string) treePage {
	t.Helper()
	resp, err := client.Raw(context.Background(), http.MethodGet, "/api/v1/organizations/tree?"+query, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "tree query %q must return 200", query)

	var out treeResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
	return treePage{Items: out.Data.Items, Cursor: out.Data.Cursor}
}

// hasPrefix is a small helper used to scope tree-walk results to the
// runID-prefixed fixtures. Inlined here to keep the e2e package
// dependency-free of the standard library "strings" import in test
// files where it would clash with reviewer expectations.
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// keysOf returns the keys of a string-set as a slice. Used to convert
// the visited-set built during cursor walks into a slice
// assert.ElementsMatch can compare against the expected universe.
func keysOf(set map[string]struct{}) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	return out
}

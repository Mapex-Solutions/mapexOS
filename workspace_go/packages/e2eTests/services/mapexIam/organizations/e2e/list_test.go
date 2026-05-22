// Page-based list tests for the IAM organizations endpoint.
//
// Covers pagination walks (forward + backward), out-of-range pages, the
// totalItems envelope guarantee, query filters (name partial, type,
// enabled, combined), and auth on the list endpoint. Cursor-based tree
// concerns live in tree_test.go.
package e2e

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
)

// TestList_NoToken_401 makes sure the GET / endpoint also enforces auth.
func TestList_NoToken_401(t *testing.T) {
	anon := httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	resp, err := anon.Raw(context.Background(), http.MethodGet, "/api/v1/organizations?perPage=1", nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// TestList_Pagination_Page1_To_Page15 creates 15 orgs sharing the same
// runID prefix, then walks pages 1..15 with perPage=1 in forward order.
// The union of ids across pages must equal the 15 created — no
// duplicates, no missing — proving the API's order is stable forward.
func TestList_Pagination_Page1_To_Page15(t *testing.T) {
	runID, want := setupListFixtures(t, listFixtureCount)
	got := walkPagesByPrefix(t, runID, 1, listFixtureCount, 1)
	assert.ElementsMatch(t, want, got, "forward walk must visit every fixture exactly once")
}

// TestList_Pagination_Page15_To_Page1 walks the same 15-page set in
// reverse order. The union must still equal the 15 created — proves
// order stability under any traversal direction.
func TestList_Pagination_Page15_To_Page1(t *testing.T) {
	runID, want := setupListFixtures(t, listFixtureCount)
	got := walkPagesByPrefix(t, runID, listFixtureCount, 1, -1)
	assert.ElementsMatch(t, want, got, "backward walk must visit every fixture exactly once")
}

// TestList_PageBeyondLast_ReturnsEmpty asks for the page immediately
// after the last valid one (totalPages + 1) and expects items=[] with
// 200, not 404 or a clamped page-1 echo.
func TestList_PageBeyondLast_ReturnsEmpty(t *testing.T) {
	runID, _ := setupListFixtures(t, listFixtureCount)
	first := fetchPage(t, listQuery(runID, 1, 10))
	require.Greater(t, first.Pagination.TotalPages, int64(0), "fixture must produce at least one page")

	beyond := fetchPage(t, listQuery(runID, int(first.Pagination.TotalPages)+1, 10))
	assert.Empty(t, beyond.Items,
		"page %d (totalPages+1) must come back empty; got %d items — backend may be clamping out-of-range pages",
		first.Pagination.TotalPages+1, len(beyond.Items))
}

// TestList_TotalItems_15 verifies pagination.totalItems matches the
// universe size regardless of perPage. Tries perPage=1 and perPage=100.
func TestList_TotalItems_15(t *testing.T) {
	runID, _ := setupListFixtures(t, listFixtureCount)

	smallPage := fetchPage(t, listQuery(runID, 1, 1))
	assert.EqualValues(t, listFixtureCount, smallPage.Pagination.TotalItems, "totalItems must equal universe size with perPage=1")

	bigPage := fetchPage(t, listQuery(runID, 1, 100))
	assert.EqualValues(t, listFixtureCount, bigPage.Pagination.TotalItems, "totalItems must equal universe size with perPage=100")
}

// TestList_Filter_Name_Match creates fixtures, then narrows the universe
// to a single fixture by exact name, expecting 1 item back.
func TestList_Filter_Name_Match(t *testing.T) {
	runID, _ := setupListFixtures(t, 3)
	targetName := orgName(runID, 1)
	page := fetchPage(t, fmt.Sprintf("perPage=10&name=%s", targetName))
	require.Len(t, page.Items, 1, "exact name filter must return one fixture")
	assert.Equal(t, targetName, page.Items[0].Name)
}

// TestList_Filter_Name_NoMatch expects an empty result when nothing
// matches the filter — endpoint must answer 200 with items=[] rather
// than 404.
func TestList_Filter_Name_NoMatch(t *testing.T) {
	page := fetchPage(t, "perPage=10&name=saga-name-does-not-exist-"+random.NewRunID())
	assert.Empty(t, page.Items)
}

// TestList_Filter_Type_Match creates 3 customer-type fixtures and
// confirms ?type=customer brings them back among the response items.
func TestList_Filter_Type_Match(t *testing.T) {
	runID, want := setupListFixtures(t, 3)
	page := fetchPage(t, fmt.Sprintf("perPage=100&type=customer&name=%s-%s", orgNamePrefix, runID))
	got := idsOf(page.Items)
	assert.Subset(t, got, want, "every fixture should appear in type=customer filter")
}

// TestList_Filter_Enabled_Match exercises the boolean filter; fixtures
// default to enabled=true so all 3 must show up.
func TestList_Filter_Enabled_Match(t *testing.T) {
	runID, want := setupListFixtures(t, 3)
	page := fetchPage(t, fmt.Sprintf("perPage=100&enabled=true&name=%s-%s", orgNamePrefix, runID))
	got := idsOf(page.Items)
	assert.Subset(t, got, want, "every fixture should appear in enabled=true filter")
}

// TestList_Filter_Combined applies type + enabled + name together. The
// AND must keep all 3 fixtures in the response.
func TestList_Filter_Combined(t *testing.T) {
	runID, want := setupListFixtures(t, 3)
	page := fetchPage(t, fmt.Sprintf("perPage=100&type=customer&enabled=true&name=%s-%s", orgNamePrefix, runID))
	got := idsOf(page.Items)
	assert.Subset(t, got, want, "AND of filters must keep every fixture")
}

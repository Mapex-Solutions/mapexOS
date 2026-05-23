// Cursor-based tests for the IAM organizations /tree endpoint.
//
// /tree returns a flat list of organizations under the active
// X-Org-Context (set in TestMain to MapexosOrgID), navigated via Next /
// Previous cursors. Tests filter the response stream by the runID
// prefix carried in the org name, the same isolation strategy used by
// the page-based list tests.
package e2e

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTree_CursorNext_FirstToLast walks the /tree endpoint forward using
// the Next cursor. The walk MUST cover every fixture exactly once and
// stop when the response signals hasNext=false (or returns empty Next).
func TestTree_CursorNext_FirstToLast(t *testing.T) {
	runID, want := setupListFixtures(t, listFixtureCount)

	visited := make(map[string]struct{}, listFixtureCount)
	cursor := ""
	for range 200 {
		page := fetchTreePage(t, fmt.Sprintf("limit=1&direction=next&cursor=%s", cursor))
		for _, item := range page.Items {
			if hasPrefix(item.Name, fmt.Sprintf("%s-%s", orgNamePrefix, runID)) {
				visited[item.ID] = struct{}{}
			}
		}
		if !page.Cursor.HasNext || page.Cursor.Next == "" {
			break
		}
		cursor = page.Cursor.Next
	}

	got := keysOf(visited)
	assert.ElementsMatch(t, want, got, "forward cursor walk must visit every fixture once")
}


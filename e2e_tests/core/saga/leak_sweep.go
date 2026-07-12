package saga

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
)

// Leak is one entity that outlived its run — its name still carries the RunID after
// the run's step Compensates were expected to remove it.
type Leak struct {
	Kind string
	ID   string
	Name string
}

// LeakSpec points the sweep at one public list endpoint: the entity Kind (for the
// report), the client that reaches it, the list path, and the JSON field holding the
// entity name the RunID is matched against.
type LeakSpec struct {
	Kind      string
	Client    *httpclient.HTTPClient
	ListPath  string
	NameField string
}

// leakItem is the decoded (id, name) pair from one list entry.
type leakItem struct {
	id   string
	name string
}

// SweepLeaks queries each spec's public list endpoint and returns items whose name
// still contains c.RunID — entities a run created but did not clean up. It observes
// ONLY the public API (never infra), mirroring the assert rule. An empty return means
// the run cleaned up after itself. A GET or decode failure on any spec is folded into
// a combined error while the other specs still run — the sweep is best-effort
// reporting, never a mutation.
func SweepLeaks(c *Context, specs []LeakSpec) ([]Leak, error) {
	var leaks []Leak
	var errs []string

	for _, spec := range specs {
		items, err := listLeakItems(c, spec)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", spec.Kind, err))
			continue
		}
		for _, it := range items {
			if strings.Contains(it.name, c.RunID) {
				leaks = append(leaks, Leak{Kind: spec.Kind, ID: it.id, Name: it.name})
			}
		}
	}

	if len(errs) > 0 {
		return leaks, fmt.Errorf("leak sweep: %s", strings.Join(errs, "; "))
	}
	return leaks, nil
}

// listLeakItems GETs one spec's list path and decodes the standard
// {data:{items:[{id, <nameField>}]}} envelope into (id, name) pairs.
func listLeakItems(c *Context, spec LeakSpec) ([]leakItem, error) {
	resp, err := spec.Client.Raw(c.Stdctx, http.MethodGet, spec.ListPath, nil)
	if err != nil {
		return nil, fmt.Errorf("get %s: %w", spec.ListPath, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get %s: unexpected status %d", spec.ListPath, resp.StatusCode)
	}

	var env struct {
		Data struct {
			Items []map[string]any `json:"items"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, fmt.Errorf("decode %s: %w", spec.ListPath, err)
	}

	out := make([]leakItem, 0, len(env.Data.Items))
	for _, m := range env.Data.Items {
		id, _ := m["id"].(string)
		name, _ := m[spec.NameField].(string)
		out = append(out, leakItem{id: id, name: name})
	}
	return out, nil
}

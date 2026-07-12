package steps

import (
	"net/http"
	"net/http/httptest"
	"sync"
)

// recordedRequest is one (method, path) the recorder server observed.
type recordedRequest struct {
	method string
	path   string
}

// recorderServer is an httptest server that records every request and replies with a
// caller-chosen status, so a step's Compensate can be exercised against it without a
// live stack.
type recorderServer struct {
	*httptest.Server
	mu       sync.Mutex
	requests []recordedRequest
}

// newRecorderServer starts a recorder replying with status to every request.
func newRecorderServer(status int) *recorderServer {
	rs := &recorderServer{}
	rs.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rs.mu.Lock()
		rs.requests = append(rs.requests, recordedRequest{r.Method, r.URL.Path})
		rs.mu.Unlock()
		w.WriteHeader(status)
	}))
	return rs
}

// count returns how many requests the server has seen.
func (rs *recorderServer) count() int {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return len(rs.requests)
}

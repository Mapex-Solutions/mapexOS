package utils

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertStatusCode checks response status code
func AssertStatusCode(t *testing.T, resp *http.Response, expected int) {
	require.NotNil(t, resp, "Response should not be nil")
	assert.Equal(t, expected, resp.StatusCode, "Unexpected status code")
}

// AssertSuccess checks if response is successful (2xx)
func AssertSuccess(t *testing.T, resp *http.Response) {
	require.NotNil(t, resp, "Response should not be nil")
	assert.GreaterOrEqual(t, resp.StatusCode, 200, "Status should be >= 200")
	assert.Less(t, resp.StatusCode, 300, "Status should be < 300")
}

// AssertCreated checks if resource was created (201)
func AssertCreated(t *testing.T, resp *http.Response) {
	AssertStatusCode(t, resp, http.StatusCreated)
}

// AssertOK checks if request was successful (200)
func AssertOK(t *testing.T, resp *http.Response) {
	AssertStatusCode(t, resp, http.StatusOK)
}

// AssertNotFound checks if resource was not found (404)
func AssertNotFound(t *testing.T, resp *http.Response) {
	AssertStatusCode(t, resp, http.StatusNotFound)
}

// AssertBadRequest checks if request was invalid (400)
func AssertBadRequest(t *testing.T, resp *http.Response) {
	AssertStatusCode(t, resp, http.StatusBadRequest)
}

// AssertForbidden checks if access was forbidden (403)
func AssertForbidden(t *testing.T, resp *http.Response) {
	AssertStatusCode(t, resp, http.StatusForbidden)
}

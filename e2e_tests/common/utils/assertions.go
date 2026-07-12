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

package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test_api_key")
	assert.NotNil(t, client)
	assert.Equal(t, "https://api.incld.dev/v1", client.BaseURL)
	assert.Equal(t, "test_api_key", client.APIKey)
	assert.NotNil(t, client.HTTPClient)
	assert.NotNil(t, client.Schedules)
	assert.NotNil(t, client.Actions)
	assert.NotNil(t, client.Runs)
	assert.NotNil(t, client.Sessions)
	assert.NotNil(t, client.Approvals)
}

func TestClientRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test_api_key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client := NewClient("test_api_key")
	client.BaseURL = server.URL

	var result map[string]bool
	err := client.request("GET", "/", nil, nil, &result)
	assert.NoError(t, err)
	assert.True(t, result["success"])
}

func TestClientRequest_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "invalid request"}`))
	}))
	defer server.Close()

	client := NewClient("test_api_key")
	client.BaseURL = server.URL

	err := client.request("GET", "/", nil, nil, nil)
	assert.Error(t, err)

	apiErr, ok := err.(*APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	assert.Equal(t, "invalid request", apiErr.Message)
}

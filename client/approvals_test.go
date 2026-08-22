package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApprovalsApproveSendsAttributedDecision(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/approvals/apr_123/decisions", r.URL.Path)

		var payload map[string]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Equal(t, "approved", payload["decision"])
		assert.Equal(t, "usr_reviewer_123", payload["approver_id"])
		assert.Equal(t, "Reviewed deploy diff", payload["reason"])

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"apr_123","project_id":"prj_123","resource_id":"document:plan","status":"approved","approver_id":"usr_reviewer_123","decisions":[],"events":[]}}`))
	}))
	defer server.Close()

	client := NewClient("test_api_key")
	client.BaseURL = server.URL

	approval, err := client.Approvals.Approve("apr_123", "usr_reviewer_123", "Reviewed deploy diff")
	require.NoError(t, err)
	assert.Equal(t, ApprovalStatusApproved, approval.Status)
	assert.Equal(t, "usr_reviewer_123", approval.ApproverID)
}

func TestApprovalsRejectSendsApproverID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Equal(t, "rejected", payload["decision"])
		assert.Equal(t, "usr_reviewer_456", payload["approver_id"])

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"apr_456","project_id":"prj_123","resource_id":"document:plan","status":"rejected","approver_id":"usr_reviewer_456","decisions":[],"events":[]}}`))
	}))
	defer server.Close()

	client := NewClient("test_api_key")
	client.BaseURL = server.URL

	approval, err := client.Approvals.Reject("apr_456", "usr_reviewer_456")
	require.NoError(t, err)
	assert.Equal(t, ApprovalStatusRejected, approval.Status)
	assert.Equal(t, "usr_reviewer_456", approval.ApproverID)
}

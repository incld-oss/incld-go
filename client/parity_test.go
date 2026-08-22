package client

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionsRunsAndSessionsParity(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/actions":
			json.NewEncoder(w).Encode(map[string]interface{}{"data": Action{ID: "act_1", Identifier: "sync_contacts"}})
		case r.Method == http.MethodGet && r.URL.Path == "/runs":
			assert.Equal(t, "sch_1", r.URL.Query().Get("schedule_id"))
			json.NewEncoder(w).Encode(RunsPage{Data: []Run{{ID: "run_1"}}, Meta: PageMeta{HasMore: false}})
		case r.Method == http.MethodGet && r.URL.Path == "/runs/run_1":
			json.NewEncoder(w).Encode(map[string]interface{}{"data": Run{ID: "run_1", ScheduleID: "sch_1"}})
		case r.Method == http.MethodPost && r.URL.Path == "/sessions":
			var payload CreateSessionParams
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			assert.Equal(t, "user_1", payload.Claims["user_id"])
			json.NewEncoder(w).Encode(SessionToken{Token: "session_1"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer server.Close()

	action, err := client.Actions.Define(DefineActionParams{Identifier: "sync_contacts", DisplayName: "Sync contacts"})
	require.NoError(t, err)
	assert.Equal(t, "act_1", action.ID)

	runs, err := client.Runs.List(&ListRunsParams{ScheduleID: "sch_1"})
	require.NoError(t, err)
	assert.Equal(t, "run_1", runs.Data[0].ID)

	run, err := client.Runs.Get("run_1")
	require.NoError(t, err)
	assert.Equal(t, "sch_1", run.ScheduleID)

	token, err := client.Sessions.Create(map[string]interface{}{"user_id": "user_1"})
	require.NoError(t, err)
	assert.Equal(t, "session_1", token.Token)
}

func TestApprovalPolicyGetAndPartialUpdate(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		policy := ApprovalPolicy{ID: "pol_1", ResourcePattern: "deploy:*", Mode: ApprovalPolicyModeAny}

		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(map[string]interface{}{"data": policy})
		case http.MethodPatch:
			var payload map[string]interface{}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			assert.Equal(t, "quorum", payload["mode"])
			assert.NotContains(t, payload, "resource_pattern")
			policy.Mode = ApprovalPolicyModeQuorum
			json.NewEncoder(w).Encode(map[string]interface{}{"data": policy})
		}
	})
	defer server.Close()

	policy, err := client.Approvals.GetPolicy("pol_1")
	require.NoError(t, err)
	assert.Equal(t, "pol_1", policy.ID)

	mode := ApprovalPolicyModeQuorum
	policy, err = client.Approvals.UpdatePolicy("pol_1", UpdateApprovalPolicyParams{Mode: &mode})
	require.NoError(t, err)
	assert.Equal(t, ApprovalPolicyModeQuorum, policy.Mode)
}

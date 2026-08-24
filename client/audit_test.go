package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditTombstone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/audit-events/event-1/tombstone", r.URL.Path)

		var body TombstoneAuditEventParams
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, AuditTombstoneDataSubjectErasure, body.Reason)
		assert.Equal(t, "privacy_admin", body.ActorID)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"event-1","component":"custom","type":"customer.exported","source":"manual","visibility":"project","data":{},"tombstoned_at":"2026-08-24T00:00:00Z","tombstone_reason":"data_subject_erasure","tombstone_event_id":"event-2","occurred_at":"2026-08-20T00:00:00Z","inserted_at":"2026-08-20T00:00:00Z"}}`))
	}))
	defer server.Close()

	client := NewClient("test_api_key")
	client.BaseURL = server.URL
	event, err := client.Audit.Tombstone("event-1", TombstoneAuditEventParams{
		Reason:  AuditTombstoneDataSubjectErasure,
		ActorID: "privacy_admin",
	}, nil)

	require.NoError(t, err)
	assert.Equal(t, "data_subject_erasure", event.TombstoneReason)
	assert.Equal(t, "event-2", event.TombstoneEventID)
}

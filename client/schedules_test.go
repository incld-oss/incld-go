package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestClient(handler http.HandlerFunc) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	client := NewClient("test_key")
	client.BaseURL = server.URL
	return client, server
}

func scheduleEnvelope(schedule Schedule) map[string]interface{} {
	return map[string]interface{}{"data": schedule}
}

func TestSchedulesService_Create(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/schedules", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(scheduleEnvelope(Schedule{ID: "sch_123"}))
	})
	defer server.Close()

	sched, err := client.Schedules.Create(CreateScheduleParams{ExternalUserID: "u1"})
	assert.NoError(t, err)
	assert.Equal(t, "sch_123", sched.ID)
}

func TestSchedulesService_List(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/schedules", r.URL.Path)
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "u1", r.URL.Query().Get("external_user_id"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []Schedule{{ID: "sch_1"}, {ID: "sch_2"}},
		})
	})
	defer server.Close()

	scheds, err := client.Schedules.List(&ListParams{Limit: 10, ExternalUserID: "u1"})
	assert.NoError(t, err)
	assert.Len(t, scheds.Data, 2)
}

func TestSchedulesService_Get(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/schedules/sch_123", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(scheduleEnvelope(Schedule{ID: "sch_123"}))
	})
	defer server.Close()

	sched, err := client.Schedules.Get("sch_123")
	assert.NoError(t, err)
	assert.Equal(t, "sch_123", sched.ID)
}

func TestSchedulesService_Update(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "/schedules/sch_123", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(scheduleEnvelope(Schedule{ID: "sch_123"}))
	})
	defer server.Close()

	sched, err := client.Schedules.Update("sch_123", UpdateScheduleParams{})
	assert.NoError(t, err)
	assert.Equal(t, "sch_123", sched.ID)
}

func TestSchedulesService_Delete(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/schedules/sch_123", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(scheduleEnvelope(Schedule{ID: "sch_123", Status: StatusDeleted}))
	})
	defer server.Close()

	_, err := client.Schedules.Delete("sch_123")
	assert.NoError(t, err)
}

func TestSchedulesService_Pause(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/schedules/sch_123/pause", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(scheduleEnvelope(Schedule{ID: "sch_123", Status: StatusPaused}))
	})
	defer server.Close()

	sched, err := client.Schedules.Pause("sch_123")
	assert.NoError(t, err)
	assert.Equal(t, StatusPaused, sched.Status)
}

func TestSchedulesService_Resume(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/schedules/sch_123/resume", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(scheduleEnvelope(Schedule{ID: "sch_123", Status: StatusActive}))
	})
	defer server.Close()

	sched, err := client.Schedules.Resume("sch_123")
	assert.NoError(t, err)
	assert.Equal(t, StatusActive, sched.Status)
}

func TestSchedulesService_ListRuns(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/schedules/sch_123/runs", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []Run{{ID: "run_1"}, {ID: "run_2"}},
		})
	})
	defer server.Close()

	runs, err := client.Schedules.ListRuns("sch_123")
	assert.NoError(t, err)
	assert.Len(t, runs.Data, 2)
}

func TestSchedulesService_ListEvents(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/schedules/sch_123/events", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []AuditEvent{{ID: "evt_1", Component: "schedules", Type: "schedule.created"}},
			"meta": PageMeta{HasMore: false},
		})
	})
	defer server.Close()

	events, err := client.Schedules.ListEvents("sch_123")
	assert.NoError(t, err)
	assert.Equal(t, "evt_1", events.Data[0].ID)
}

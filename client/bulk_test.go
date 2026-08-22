package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBulkCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bulk-operations" || r.Method != http.MethodPost {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["action"] != "contacts.sync" {
			t.Fatalf("unexpected action: %v", body["action"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"bulk_1","project_id":"project_1","action":"contacts.sync","status":"queued","metadata":{},"chunk_size":100,"progress":{"total_items":2,"total_chunks":1,"completed_chunks":0,"succeeded_chunks":0,"failed_chunks":0,"percentage":0},"inserted_at":"2026-08-18T00:00:00Z","updated_at":"2026-08-18T00:00:00Z"}}`))
	}))
	defer server.Close()

	client := NewClient("test")
	client.BaseURL = server.URL
	operation, err := client.Bulk.Create(&CreateBulkOperationParams{
		Action:    "contacts.sync",
		Items:     []map[string]interface{}{{"id": "one"}, {"id": "two"}},
		ChunkSize: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if operation.Progress.TotalItems != 2 {
		t.Fatalf("expected 2 items, got %d", operation.Progress.TotalItems)
	}
}

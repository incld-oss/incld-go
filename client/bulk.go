package client

import (
	"fmt"
	"net/url"
	"strconv"
)

type BulkProgress struct {
	TotalItems      int     `json:"total_items"`
	TotalChunks     int     `json:"total_chunks"`
	CompletedChunks int     `json:"completed_chunks"`
	SucceededChunks int     `json:"succeeded_chunks"`
	FailedChunks    int     `json:"failed_chunks"`
	Percentage      float64 `json:"percentage"`
}

type BulkOperation struct {
	ID                     string                 `json:"id"`
	ProjectID              string                 `json:"project_id"`
	ExternalOrganizationID string                 `json:"external_organization_id,omitempty"`
	Action                 string                 `json:"action"`
	Status                 string                 `json:"status"`
	Metadata               map[string]interface{} `json:"metadata"`
	ChunkSize              int                    `json:"chunk_size"`
	Progress               BulkProgress           `json:"progress"`
	StartedAt              *string                `json:"started_at"`
	CompletedAt            *string                `json:"completed_at"`
	CancelledAt            *string                `json:"cancelled_at"`
	InsertedAt             string                 `json:"inserted_at"`
	UpdatedAt              string                 `json:"updated_at"`
}

type BulkChunk struct {
	ID             string                   `json:"id"`
	OperationID    string                   `json:"operation_id"`
	Index          int                      `json:"index"`
	Items          []map[string]interface{} `json:"items"`
	Status         string                   `json:"status"`
	AttemptCount   int                      `json:"attempt_count"`
	ResponseStatus *int                     `json:"response_status"`
	ResponseBody   *string                  `json:"response_body"`
	Error          *string                  `json:"error"`
	StartedAt      *string                  `json:"started_at"`
	CompletedAt    *string                  `json:"completed_at"`
	InsertedAt     string                   `json:"inserted_at"`
	UpdatedAt      string                   `json:"updated_at"`
}

type BulkEvent struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	ActorID    *string                `json:"actor_id"`
	Data       map[string]interface{} `json:"data"`
	InsertedAt string                 `json:"inserted_at"`
}

type CreateBulkOperationParams struct {
	ExternalOrganizationID string                   `json:"external_organization_id,omitempty"`
	Action                 string                   `json:"action"`
	Items                  []map[string]interface{} `json:"items"`
	ChunkSize              int                      `json:"chunk_size,omitempty"`
	Metadata               map[string]interface{}   `json:"metadata,omitempty"`
}

type ListBulkOperationsParams struct {
	ExternalOrganizationID string
	Status                 string
	Action                 string
	Limit                  int
	Cursor                 string
}

type BulkService struct{ client *Client }

func (s *BulkService) List(params *ListBulkOperationsParams) (*Page[BulkOperation], error) {
	query := url.Values{}
	if params != nil {
		if params.ExternalOrganizationID != "" {
			query.Set("external_organization_id", params.ExternalOrganizationID)
		}
		if params.Status != "" {
			query.Set("status", params.Status)
		}
		if params.Action != "" {
			query.Set("action", params.Action)
		}
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Cursor != "" {
			query.Set("cursor", params.Cursor)
		}
	}

	var resp Page[BulkOperation]
	if err := s.client.request("GET", "/bulk-operations", nil, query, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *BulkService) Get(id string) (*BulkOperation, error) {
	var resp struct {
		Data BulkOperation `json:"data"`
	}
	if err := s.client.request("GET", fmt.Sprintf("/bulk-operations/%s", id), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *BulkService) Create(params *CreateBulkOperationParams) (*BulkOperation, error) {
	return s.CreateWithOptions(params, nil)
}

func (s *BulkService) CreateWithOptions(params *CreateBulkOperationParams, options *RequestOptions) (*BulkOperation, error) {
	var resp struct {
		Data BulkOperation `json:"data"`
	}
	if err := s.client.requestWithOptions("POST", "/bulk-operations", params, nil, &resp, options); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *BulkService) Chunks(id string) (*Page[BulkChunk], error) {
	var resp Page[BulkChunk]
	if err := s.client.request("GET", fmt.Sprintf("/bulk-operations/%s/chunks", id), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *BulkService) Events(id string) (*Page[BulkEvent], error) {
	var resp Page[BulkEvent]
	if err := s.client.request("GET", fmt.Sprintf("/bulk-operations/%s/events", id), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *BulkService) Cancel(id, actorID, reason string) (*BulkOperation, error) {
	body := map[string]string{"actor_id": actorID, "reason": reason}
	var resp struct {
		Data BulkOperation `json:"data"`
	}
	if err := s.client.request("POST", fmt.Sprintf("/bulk-operations/%s/cancel", id), body, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

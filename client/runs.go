package client

import (
	"fmt"
	"net/url"
	"strconv"
)

type ListRunsParams struct {
	Limit          int
	Cursor         string
	Status         string
	Action         string
	ScheduleID     string
	ExternalUserID string
}

type RunsPage = Page[Run]

type RunsService struct{ client *Client }

func (s *RunsService) Get(id string) (*Run, error) {
	var resp struct {
		Data Run `json:"data"`
	}
	if err := s.client.request("GET", fmt.Sprintf("/runs/%s", id), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *RunsService) List(params *ListRunsParams) (*RunsPage, error) {
	query := url.Values{}
	if params != nil {
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Cursor != "" {
			query.Set("cursor", params.Cursor)
		}
		if params.Status != "" {
			query.Set("status", params.Status)
		}
		if params.Action != "" {
			query.Set("action", params.Action)
		}
		if params.ScheduleID != "" {
			query.Set("schedule_id", params.ScheduleID)
		}
		if params.ExternalUserID != "" {
			query.Set("external_user_id", params.ExternalUserID)
		}
	}

	var resp RunsPage
	if err := s.client.request("GET", "/runs", nil, query, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListPage is retained as a source-compatible alias for List.
func (s *RunsService) ListPage(params *ListRunsParams) (*RunsPage, error) { return s.List(params) }

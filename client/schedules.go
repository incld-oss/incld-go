package client

import (
	"fmt"
	"net/url"
	"strconv"
)

type SchedulesService struct {
	client *Client
}

type scheduleResponse struct {
	Data Schedule `json:"data"`
}

func (s *SchedulesService) Create(params CreateScheduleParams) (*Schedule, error) {
	return s.CreateWithOptions(params, nil)
}

func (s *SchedulesService) CreateWithOptions(params CreateScheduleParams, options *RequestOptions) (*Schedule, error) {
	var resp scheduleResponse
	if err := s.client.requestWithOptions("POST", "/schedules", params, nil, &resp, options); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *SchedulesService) List(params *ListParams) (*Page[Schedule], error) {
	var resp Page[Schedule]
	var query url.Values
	if params != nil {
		query = url.Values{}
		if params.ExternalUserID != "" {
			query.Add("external_user_id", params.ExternalUserID)
		}
		if params.ExternalOrganizationID != "" {
			query.Add("external_organization_id", params.ExternalOrganizationID)
		}
		if params.Action != "" {
			query.Add("action", params.Action)
		}
		if params.Status != "" {
			query.Add("status", params.Status)
		}
		if params.Limit > 0 {
			query.Add("limit", strconv.Itoa(params.Limit))
		}
		if params.Cursor != "" {
			query.Add("cursor", params.Cursor)
		}
	}
	if err := s.client.request("GET", "/schedules", nil, query, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *SchedulesService) Get(id string) (*Schedule, error) {
	var resp scheduleResponse
	if err := s.client.request("GET", fmt.Sprintf("/schedules/%s", id), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *SchedulesService) Update(id string, params UpdateScheduleParams) (*Schedule, error) {
	var resp scheduleResponse
	if err := s.client.request("PATCH", fmt.Sprintf("/schedules/%s", id), params, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *SchedulesService) Delete(id string) (*Schedule, error) {
	var resp scheduleResponse
	if err := s.client.request("DELETE", fmt.Sprintf("/schedules/%s", id), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *SchedulesService) Pause(id string) (*Schedule, error) {
	var resp scheduleResponse
	if err := s.client.request("POST", fmt.Sprintf("/schedules/%s/pause", id), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *SchedulesService) Resume(id string) (*Schedule, error) {
	var resp scheduleResponse
	if err := s.client.request("POST", fmt.Sprintf("/schedules/%s/resume", id), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *SchedulesService) ListRuns(scheduleID string) (*Page[Run], error) {
	var resp Page[Run]
	if err := s.client.request("GET", fmt.Sprintf("/schedules/%s/runs", scheduleID), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListEvents returns the append-only Audit lifecycle history for one schedule.
func (s *SchedulesService) ListEvents(scheduleID string) (*Page[AuditEvent], error) {
	var resp Page[AuditEvent]
	if err := s.client.request("GET", fmt.Sprintf("/schedules/%s/events", scheduleID), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

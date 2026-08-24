package client

import (
	"fmt"
	"net/url"
	"strconv"
)

type AuditEvent struct {
	ID                     string                 `json:"id"`
	Component              string                 `json:"component"`
	Type                   string                 `json:"type"`
	ActorID                string                 `json:"actor_id,omitempty"`
	ExternalOrganizationID string                 `json:"external_organization_id,omitempty"`
	SubjectType            string                 `json:"subject_type,omitempty"`
	SubjectID              string                 `json:"subject_id,omitempty"`
	Source                 string                 `json:"source"`
	Visibility             string                 `json:"visibility"`
	Data                   map[string]interface{} `json:"data"`
	TombstonedAt           string                 `json:"tombstoned_at,omitempty"`
	TombstoneReason        string                 `json:"tombstone_reason,omitempty"`
	TombstoneEventID       string                 `json:"tombstone_event_id,omitempty"`
	OccurredAt             string                 `json:"occurred_at"`
	InsertedAt             string                 `json:"inserted_at"`
}

type AuditTombstoneReason string

const (
	AuditTombstoneDataSubjectErasure AuditTombstoneReason = "data_subject_erasure"
	AuditTombstoneSensitiveData      AuditTombstoneReason = "sensitive_data"
	AuditTombstoneCustomerRequest    AuditTombstoneReason = "customer_request"
)

type TombstoneAuditEventParams struct {
	Reason  AuditTombstoneReason `json:"reason"`
	ActorID string               `json:"actor_id,omitempty"`
}

type ListAuditEventsParams struct {
	ExternalOrganizationID string
	Component              string
	Type                   string
	TypePrefix             string
	ActorID                string
	SubjectType            string
	SubjectID              string
	Since                  string
	Until                  string
	Limit                  int
	Cursor                 string
}

type CreateAuditEventParams struct {
	ExternalOrganizationID string                 `json:"external_organization_id,omitempty"`
	Type                   string                 `json:"type"`
	ActorID                string                 `json:"actor_id,omitempty"`
	SubjectType            string                 `json:"subject_type,omitempty"`
	SubjectID              string                 `json:"subject_id,omitempty"`
	Visibility             string                 `json:"visibility,omitempty"`
	ParticipantIDs         []string               `json:"participant_ids,omitempty"`
	AllowedViewerIDs       []string               `json:"allowed_viewer_ids,omitempty"`
	Data                   map[string]interface{} `json:"data,omitempty"`
	OccurredAt             string                 `json:"occurred_at,omitempty"`
}

type AuditService struct{ client *Client }

func (s *AuditService) List(params *ListAuditEventsParams) (*Page[AuditEvent], error) {
	query := url.Values{}
	if params != nil {
		if params.ExternalOrganizationID != "" {
			query.Set("external_organization_id", params.ExternalOrganizationID)
		}
		if params.Component != "" {
			query.Set("component", params.Component)
		}
		if params.Type != "" {
			query.Set("type", params.Type)
		}
		if params.TypePrefix != "" {
			query.Set("type_prefix", params.TypePrefix)
		}
		if params.ActorID != "" {
			query.Set("actor_id", params.ActorID)
		}
		if params.SubjectType != "" {
			query.Set("subject_type", params.SubjectType)
		}
		if params.SubjectID != "" {
			query.Set("subject_id", params.SubjectID)
		}
		if params.Since != "" {
			query.Set("since", params.Since)
		}
		if params.Until != "" {
			query.Set("until", params.Until)
		}
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Cursor != "" {
			query.Set("cursor", params.Cursor)
		}
	}
	var response Page[AuditEvent]
	if err := s.client.request("GET", "/audit-events", nil, query, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *AuditService) Get(id string) (*AuditEvent, error) {
	var response struct {
		Data AuditEvent `json:"data"`
	}
	if err := s.client.request("GET", fmt.Sprintf("/audit-events/%s", id), nil, nil, &response); err != nil {
		return nil, err
	}
	return &response.Data, nil
}

func (s *AuditService) Create(params CreateAuditEventParams, options *RequestOptions) (*AuditEvent, error) {
	var response struct {
		Data AuditEvent `json:"data"`
	}
	if err := s.client.requestWithOptions("POST", "/audit-events", params, nil, &response, options); err != nil {
		return nil, err
	}
	return &response.Data, nil
}

// Tombstone erases PII-bearing fields while preserving the event envelope and
// appending a durable reason-coded record. It requires a trusted server client.
func (s *AuditService) Tombstone(id string, params TombstoneAuditEventParams, options *RequestOptions) (*AuditEvent, error) {
	var response struct {
		Data AuditEvent `json:"data"`
	}
	path := fmt.Sprintf("/audit-events/%s/tombstone", url.PathEscape(id))
	if err := s.client.requestWithOptions("POST", path, params, nil, &response, options); err != nil {
		return nil, err
	}
	return &response.Data, nil
}

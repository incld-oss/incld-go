package client

import (
	"fmt"
	"net/url"
)

type ApprovalStatus string
type ApprovalDecisionValue string
type ApprovalPolicyMode string

const (
	ApprovalStatusPending   ApprovalStatus = "pending"
	ApprovalStatusApproved  ApprovalStatus = "approved"
	ApprovalStatusRejected  ApprovalStatus = "rejected"
	ApprovalStatusCancelled ApprovalStatus = "cancelled"
	ApprovalStatusExpired   ApprovalStatus = "expired"
	ApprovalStatusRevoked   ApprovalStatus = "revoked"

	ApprovalDecisionApproved ApprovalDecisionValue = "approved"
	ApprovalDecisionRejected ApprovalDecisionValue = "rejected"

	ApprovalPolicyModeAny    ApprovalPolicyMode = "any"
	ApprovalPolicyModeAll    ApprovalPolicyMode = "all"
	ApprovalPolicyModeQuorum ApprovalPolicyMode = "quorum"
)

type ApprovalDecision struct {
	ID         string `json:"id"`
	Decision   string `json:"decision"`
	ActorID    string `json:"actor_id"`
	Reason     string `json:"reason,omitempty"`
	InsertedAt string `json:"inserted_at"`
}

type ApprovalEvent struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	ActorID    string                 `json:"actor_id,omitempty"`
	Data       map[string]interface{} `json:"data"`
	InsertedAt string                 `json:"inserted_at"`
}

type ApprovalPolicySnapshot struct {
	ID                string             `json:"id,omitempty"`
	Revision          int                `json:"revision,omitempty"`
	ResourcePattern   string             `json:"resource_pattern,omitempty"`
	AllowedApprovers  []string           `json:"allowed_approvers,omitempty"`
	Mode              ApprovalPolicyMode `json:"mode,omitempty"`
	RequiredApprovals int                `json:"required_approvals,omitempty"`
	AllowSelfApproval bool               `json:"allow_self_approval,omitempty"`
}

type Approval struct {
	ID             string                 `json:"id"`
	ProjectID      string                 `json:"project_id"`
	ResourceType   string                 `json:"resource_type"`
	ResourceID     string                 `json:"resource_id"`
	Action         string                 `json:"action"`
	RequesterID    string                 `json:"requester_id,omitempty"`
	ApproverID     string                 `json:"approver_id,omitempty"`
	Status         ApprovalStatus         `json:"status"`
	Title          string                 `json:"title,omitempty"`
	Description    string                 `json:"description,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	Revision       int                    `json:"revision,omitempty"`
	PolicyID       string                 `json:"policy_id,omitempty"`
	PolicySnapshot ApprovalPolicySnapshot `json:"policy_snapshot,omitempty"`
	IdempotencyKey string                 `json:"idempotency_key,omitempty"`
	ExpiresAt      string                 `json:"expires_at,omitempty"`
	ResolvedAt     string                 `json:"resolved_at,omitempty"`
	CancelledAt    string                 `json:"cancelled_at,omitempty"`
	RevokedAt      string                 `json:"revoked_at,omitempty"`
	Decisions      []ApprovalDecision     `json:"decisions,omitempty"`
	Events         []ApprovalEvent        `json:"events,omitempty"`
	InsertedAt     string                 `json:"inserted_at,omitempty"`
	UpdatedAt      string                 `json:"updated_at,omitempty"`
}

type ApprovalPolicy struct {
	ID                string             `json:"id"`
	ProjectID         string             `json:"project_id"`
	ResourcePattern   string             `json:"resource_pattern"`
	AllowedApprovers  []string           `json:"allowed_approvers"`
	Mode              ApprovalPolicyMode `json:"mode"`
	RequiredApprovals int                `json:"required_approvals"`
	AllowSelfApproval bool               `json:"allow_self_approval"`
	Revision          int                `json:"revision"`
	InsertedAt        string             `json:"inserted_at,omitempty"`
	UpdatedAt         string             `json:"updated_at,omitempty"`
}

type CreateApprovalParams struct {
	ResourceType   string                 `json:"resource_type"`
	ResourceID     string                 `json:"resource_id"`
	Action         string                 `json:"action"`
	RequesterID    string                 `json:"requester_id,omitempty"`
	Title          string                 `json:"title,omitempty"`
	Description    string                 `json:"description,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	ExpiresAt      string                 `json:"expires_at,omitempty"`
	IdempotencyKey string                 `json:"idempotency_key,omitempty"`
}

type UpdateApprovalParams struct {
	Title       *string                 `json:"title,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Metadata    *map[string]interface{} `json:"metadata,omitempty"`
	ExpiresAt   *string                 `json:"expires_at,omitempty"`
}

type ListApprovalsParams struct {
	ResourceID  string
	RequesterID string
	ApproverID  string
	Status      ApprovalStatus
	Limit       int
	Cursor      string
}

type ApprovalCheck struct {
	Approved   bool   `json:"approved"`
	Status     string `json:"status"`
	ApprovalID string `json:"approval_id,omitempty"`
}

type CreateApprovalPolicyParams struct {
	ResourcePattern   string             `json:"resource_pattern"`
	AllowedApprovers  []string           `json:"allowed_approvers"`
	Mode              ApprovalPolicyMode `json:"mode,omitempty"`
	RequiredApprovals int                `json:"required_approvals,omitempty"`
	AllowSelfApproval bool               `json:"allow_self_approval"`
}

type UpdateApprovalPolicyParams struct {
	ResourcePattern   *string             `json:"resource_pattern,omitempty"`
	AllowedApprovers  *[]string           `json:"allowed_approvers,omitempty"`
	Mode              *ApprovalPolicyMode `json:"mode,omitempty"`
	RequiredApprovals *int                `json:"required_approvals,omitempty"`
	AllowSelfApproval *bool               `json:"allow_self_approval,omitempty"`
}

type ApprovalsService struct{ client *Client }

func (s *ApprovalsService) List(params *ListApprovalsParams) (*Page[Approval], error) {
	query := url.Values{}
	if params != nil {
		if params.ResourceID != "" {
			query.Set("resource_id", params.ResourceID)
		}
		if params.RequesterID != "" {
			query.Set("requester_id", params.RequesterID)
		}
		if params.ApproverID != "" {
			query.Set("approver_id", params.ApproverID)
		}
		if params.Status != "" {
			query.Set("status", string(params.Status))
		}
		if params.Limit > 0 {
			query.Set("limit", fmt.Sprint(params.Limit))
		}
		if params.Cursor != "" {
			query.Set("cursor", params.Cursor)
		}
	}
	var resp Page[Approval]
	if err := s.client.request("GET", "/approvals", nil, query, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *ApprovalsService) Get(id string) (*Approval, error) {
	var resp struct {
		Data Approval `json:"data"`
	}
	if err := s.client.request("GET", fmt.Sprintf("/approvals/%s", id), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *ApprovalsService) Check(resourceType, resourceID, action, requesterID string) (*ApprovalCheck, error) {
	body := map[string]string{"resource_type": resourceType, "resource_id": resourceID, "action": action, "requester_id": requesterID}
	var resp struct {
		Data ApprovalCheck `json:"data"`
	}
	if err := s.client.request("POST", "/approvals/check", body, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *ApprovalsService) Create(params CreateApprovalParams) (*Approval, error) {
	return s.CreateWithOptions(params, nil)
}

func (s *ApprovalsService) CreateWithOptions(params CreateApprovalParams, options *RequestOptions) (*Approval, error) {
	var resp struct {
		Data Approval `json:"data"`
	}
	if err := s.client.requestWithOptions("POST", "/approvals", params, nil, &resp, options); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *ApprovalsService) Update(id string, params UpdateApprovalParams) (*Approval, error) {
	var resp struct {
		Data Approval `json:"data"`
	}
	if err := s.client.request("PATCH", fmt.Sprintf("/approvals/%s", id), params, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *ApprovalsService) Delete(id string) (*Approval, error) {
	var resp struct {
		Data Approval `json:"data"`
	}
	if err := s.client.request("DELETE", fmt.Sprintf("/approvals/%s", id), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

type approvalDecisionParams struct {
	Decision   ApprovalDecisionValue `json:"decision"`
	ApproverID string                `json:"approver_id,omitempty"`
	Reason     string                `json:"reason,omitempty"`
}

type approvalCommandParams struct {
	ActorID string `json:"actor_id,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

func (s *ApprovalsService) Decide(id string, decision ApprovalDecisionValue, approverID string, reasons ...string) (*Approval, error) {
	reason := ""
	if len(reasons) > 0 {
		reason = reasons[0]
	}
	return s.command(fmt.Sprintf("/approvals/%s/decisions", id), approvalDecisionParams{Decision: decision, ApproverID: approverID, Reason: reason})
}

func (s *ApprovalsService) DecideWithOptions(id string, decision ApprovalDecisionValue, approverID, reason string, options *RequestOptions) (*Approval, error) {
	return s.commandWithOptions(fmt.Sprintf("/approvals/%s/decisions", id), approvalDecisionParams{Decision: decision, ApproverID: approverID, Reason: reason}, options)
}

func (s *ApprovalsService) Approve(id, approverID string, reasons ...string) (*Approval, error) {
	return s.Decide(id, ApprovalDecisionApproved, approverID, reasons...)
}

func (s *ApprovalsService) Reject(id, approverID string, reasons ...string) (*Approval, error) {
	return s.Decide(id, ApprovalDecisionRejected, approverID, reasons...)
}

func (s *ApprovalsService) Cancel(id, actorID string, reasons ...string) (*Approval, error) {
	reason := ""
	if len(reasons) > 0 {
		reason = reasons[0]
	}
	return s.command(fmt.Sprintf("/approvals/%s/cancel", id), approvalCommandParams{ActorID: actorID, Reason: reason})
}

func (s *ApprovalsService) Revoke(id, actorID string, reasons ...string) (*Approval, error) {
	reason := ""
	if len(reasons) > 0 {
		reason = reasons[0]
	}
	return s.command(fmt.Sprintf("/approvals/%s/revoke", id), approvalCommandParams{ActorID: actorID, Reason: reason})
}

func (s *ApprovalsService) Events(id string) (*Page[ApprovalEvent], error) {
	var resp Page[ApprovalEvent]
	if err := s.client.request("GET", fmt.Sprintf("/approvals/%s/events", id), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *ApprovalsService) ListPolicies() (*Page[ApprovalPolicy], error) {
	var resp Page[ApprovalPolicy]
	if err := s.client.request("GET", "/approval-policies", nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *ApprovalsService) GetPolicy(id string) (*ApprovalPolicy, error) {
	var resp struct {
		Data ApprovalPolicy `json:"data"`
	}
	if err := s.client.request("GET", fmt.Sprintf("/approval-policies/%s", id), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *ApprovalsService) CreatePolicy(params CreateApprovalPolicyParams) (*ApprovalPolicy, error) {
	var resp struct {
		Data ApprovalPolicy `json:"data"`
	}
	if err := s.client.request("POST", "/approval-policies", params, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *ApprovalsService) UpdatePolicy(id string, params UpdateApprovalPolicyParams) (*ApprovalPolicy, error) {
	var resp struct {
		Data ApprovalPolicy `json:"data"`
	}
	if err := s.client.request("PATCH", fmt.Sprintf("/approval-policies/%s", id), params, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *ApprovalsService) DeletePolicy(id string) (*ApprovalPolicy, error) {
	var resp struct {
		Data ApprovalPolicy `json:"data"`
	}
	if err := s.client.request("DELETE", fmt.Sprintf("/approval-policies/%s", id), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *ApprovalsService) command(path string, body interface{}) (*Approval, error) {
	return s.commandWithOptions(path, body, nil)
}

func (s *ApprovalsService) commandWithOptions(path string, body interface{}, options *RequestOptions) (*Approval, error) {
	var resp struct {
		Data Approval `json:"data"`
	}
	if err := s.client.requestWithOptions("POST", path, body, nil, &resp, options); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

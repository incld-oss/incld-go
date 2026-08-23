package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client is the main interface for interacting with the incld API.
type Client struct {
	BaseURL        string
	APIKey         string
	HTTPClient     *http.Client
	organizationID string
	userID         string

	Schedules *SchedulesService
	Actions   *ActionsService
	Runs      *RunsService
	Sessions  *SessionsService
	Approvals *ApprovalsService
	Bulk      *BulkService
	Audit     *AuditService
}

// WithScope returns a client copy whose requests are constrained to one external
// organization and, when userID is non-empty, one end user.
func (c *Client) WithScope(organizationID, userID string) *Client {
	if strings.TrimSpace(organizationID) == "" {
		panic("incld: WithScope requires a non-empty organization ID")
	}
	if userID != "" && strings.TrimSpace(userID) == "" {
		panic("incld: WithScope user ID must not be whitespace")
	}

	scoped := NewClient(c.APIKey)
	scoped.BaseURL = c.BaseURL
	scoped.HTTPClient = c.HTTPClient
	scoped.organizationID = organizationID
	scoped.userID = userID
	return scoped
}

// NewClient creates a new Client with the given API key.
func NewClient(apiKey string) *Client {
	c := &Client{
		BaseURL:    "https://api.incld.dev/v1",
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
	}
	c.Schedules = &SchedulesService{client: c}
	c.Actions = &ActionsService{client: c}
	c.Runs = &RunsService{client: c}
	c.Sessions = &SessionsService{client: c}
	c.Approvals = &ApprovalsService{client: c}
	c.Bulk = &BulkService{client: c}
	c.Audit = &AuditService{client: c}
	return c
}

// APIError represents an error returned by the API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Fields     map[string][]string
	RequestID  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("incld API error %d (%s): %s", e.StatusCode, e.Code, e.Message)
}

type RequestOptions struct {
	Context        context.Context
	IdempotencyKey string
}

type PageMeta struct {
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}

type Page[T any] struct {
	Data []T      `json:"data"`
	Meta PageMeta `json:"meta"`
}

func (c *Client) request(method, path string, body interface{}, query url.Values, v interface{}) error {
	return c.requestWithOptions(method, path, body, query, v, nil)
}

func (c *Client) requestWithOptions(method, path string, body interface{}, query url.Values, v interface{}, options *RequestOptions) error {
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return err
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(b)
	}

	ctx := context.Background()
	if options != nil && options.Context != nil {
		ctx = options.Context
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	if c.organizationID != "" {
		req.Header.Set("Incld-Organization-Id", c.organizationID)
	}
	if c.userID != "" {
		req.Header.Set("Incld-User-Id", c.userID)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if options != nil && options.IdempotencyKey != "" {
		req.Header.Set("Idempotency-Key", options.IdempotencyKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		var envelope struct {
			Error struct {
				Code      string              `json:"code"`
				Message   string              `json:"message"`
				Fields    map[string][]string `json:"fields"`
				RequestID string              `json:"request_id"`
			} `json:"error"`
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		_ = json.Unmarshal(b, &envelope)
		message := envelope.Error.Message
		code := envelope.Error.Code
		if message == "" {
			message = envelope.Message
		}
		if code == "" {
			code = envelope.Code
		}
		if message == "" {
			message = resp.Status
		}
		if code == "" {
			code = "unknown_error"
		}
		requestID := envelope.Error.RequestID
		if requestID == "" {
			requestID = resp.Header.Get("X-Request-ID")
		}
		return &APIError{StatusCode: resp.StatusCode, Code: code, Message: message, Fields: envelope.Error.Fields, RequestID: requestID}
	}

	if v != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}

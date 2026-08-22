package client

import "fmt"

type Action struct {
	ID            string                 `json:"id"`
	Identifier    string                 `json:"identifier"`
	DisplayName   string                 `json:"display_name"`
	Description   string                 `json:"description,omitempty"`
	PayloadSchema map[string]interface{} `json:"payload_schema,omitempty"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
	InsertedAt    string                 `json:"inserted_at,omitempty"`
	UpdatedAt     string                 `json:"updated_at,omitempty"`
}

type DefineActionParams struct {
	Identifier    string                 `json:"identifier"`
	DisplayName   string                 `json:"display_name"`
	Description   string                 `json:"description,omitempty"`
	PayloadSchema map[string]interface{} `json:"payload_schema,omitempty"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
}

type ActionsService struct{ client *Client }

func (s *ActionsService) Define(params DefineActionParams) (*Action, error) {
	var resp struct {
		Data Action `json:"data"`
	}
	if err := s.client.request("POST", "/actions", params, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *ActionsService) List() (*Page[Action], error) {
	var resp Page[Action]
	if err := s.client.request("GET", "/actions", nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *ActionsService) Get(identifier string) (*Action, error) {
	var resp struct {
		Data Action `json:"data"`
	}
	if err := s.client.request("GET", fmt.Sprintf("/actions/%s", identifier), nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

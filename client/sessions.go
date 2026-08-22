package client

type CreateSessionParams struct {
	Claims map[string]interface{} `json:"claims"`
}

type SessionToken struct {
	Token string `json:"token"`
}

type SessionsService struct{ client *Client }

func (s *SessionsService) Create(claims map[string]interface{}) (*SessionToken, error) {
	var token SessionToken
	if err := s.client.request("POST", "/sessions", CreateSessionParams{Claims: claims}, nil, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

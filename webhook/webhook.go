package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// VerifyWebhookSignature verifies the signature of an incoming webhook payload.
// If toleranceSeconds is 0, it defaults to 300 (5 minutes).
func VerifyWebhookSignature(payload []byte, signature string, secret string, toleranceSeconds int64) bool {
	if toleranceSeconds == 0 {
		toleranceSeconds = 300
	}

	parts := strings.Split(signature, ",")
	var timestampStr, providedSignature string
	for _, part := range parts {
		if strings.HasPrefix(part, "t=") {
			timestampStr = strings.TrimPrefix(part, "t=")
		} else if strings.HasPrefix(part, "v1=") {
			providedSignature = strings.TrimPrefix(part, "v1=")
		}
	}

	if timestampStr == "" || providedSignature == "" {
		return false
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return false
	}

	now := time.Now().Unix()
	diff := now - timestamp
	if diff < 0 {
		diff = -diff
	}
	if diff > toleranceSeconds {
		return false
	}

	signedPayload := fmt.Sprintf("%d.%s", timestamp, string(payload))

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return subtle.ConstantTimeCompare([]byte(expectedSignature), []byte(providedSignature)) == 1
}

// Event represents a strongly typed webhook event.
type Event struct {
	Type string `json:"type"`
	Data struct {
		Run struct {
			Schedule struct {
				Action struct {
					Identifier  string `json:"identifier"`
					DisplayName string `json:"display_name"`
				} `json:"action"`
			} `json:"schedule"`
		} `json:"run"`
	} `json:"data"`
}

// VerifyAndParse reads the request body, verifies the HMAC signature using the existing logic,
// and unmarshals the JSON into a strongly typed Event struct.
func VerifyAndParse(r *http.Request, secret string) (*Event, error) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	signature := r.Header.Get("Incld-Signature")
	if signature == "" {
		return nil, errors.New("missing webhook signature")
	}

	if !VerifyWebhookSignature(payload, signature, secret, 300) {
		return nil, errors.New("invalid webhook signature")
	}

	var event Event
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

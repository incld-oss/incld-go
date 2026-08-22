package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVerifyWebhookSignature(t *testing.T) {
	secret := "test_secret"
	payload := []byte(`{"id": "ev_123", "type": "schedule.created"}`)
	timestamp := time.Now().Unix()

	// Generate valid signature
	signedPayload := fmt.Sprintf("%d.%s", timestamp, string(payload))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	validSignatureHeader := fmt.Sprintf("t=%d,v1=%s", timestamp, expectedSignature)

	t.Run("Valid signature", func(t *testing.T) {
		assert.True(t, VerifyWebhookSignature(payload, validSignatureHeader, secret, 300))
	})

	t.Run("Valid signature with default tolerance", func(t *testing.T) {
		assert.True(t, VerifyWebhookSignature(payload, validSignatureHeader, secret, 0))
	})

	t.Run("Invalid signature", func(t *testing.T) {
		invalidSignatureHeader := fmt.Sprintf("t=%d,v1=invalid_signature", timestamp)
		assert.False(t, VerifyWebhookSignature(payload, invalidSignatureHeader, secret, 300))
	})

	t.Run("Expired signature", func(t *testing.T) {
		oldTimestamp := time.Now().Add(-10 * time.Minute).Unix()
		oldSignedPayload := fmt.Sprintf("%d.%s", oldTimestamp, string(payload))
		oldMac := hmac.New(sha256.New, []byte(secret))
		oldMac.Write([]byte(oldSignedPayload))
		oldExpectedSignature := hex.EncodeToString(oldMac.Sum(nil))

		expiredSignatureHeader := fmt.Sprintf("t=%d,v1=%s", oldTimestamp, oldExpectedSignature)
		assert.False(t, VerifyWebhookSignature(payload, expiredSignatureHeader, secret, 300))
	})

	t.Run("Malformed signature header", func(t *testing.T) {
		assert.False(t, VerifyWebhookSignature(payload, "invalid_header", secret, 300))
	})
}

func TestVerifyAndParse(t *testing.T) {
	secret := "test_secret"
	payload := `{"type": "run.created", "data": {"run": {"schedule": {"action": {"identifier": "test_action", "display_name": "Test action"}}}}}`
	timestamp := time.Now().Unix()

	signedPayload := fmt.Sprintf("%d.%s", timestamp, payload)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	signature := fmt.Sprintf("t=%d,v1=%s", timestamp, hex.EncodeToString(mac.Sum(nil)))

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	req.Header.Set("Incld-Signature", signature)

	event, err := VerifyAndParse(req, secret)
	assert.NoError(t, err)
	assert.Equal(t, "run.created", event.Type)
	assert.Equal(t, "test_action", event.Data.Run.Schedule.Action.Identifier)
}

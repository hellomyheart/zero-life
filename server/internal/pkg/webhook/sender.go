package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	maxRetries   = 3
	initialDelay = 1 * time.Second
)

// SendResult holds the result of a webhook send attempt.
type SendResult struct {
	ResponseCode int
	ResponseBody string
	Success      bool
}

// Send delivers the payload to the given URL with optional HMAC-SHA256 signature.
// It retries up to 3 times with exponential backoff (1s, 2s, 4s).
func Send(url string, secret string, payload interface{}) SendResult {
	body, err := json.Marshal(payload)
	if err != nil {
		return SendResult{Success: false}
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			delay := initialDelay * (1 << (attempt - 1)) // 1s, 2s, 4s
			time.Sleep(delay)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		// Sign with HMAC-SHA256 if secret is provided
		if secret != "" {
			mac := hmac.New(sha256.New, []byte(secret))
			mac.Write(body)
			sig := hex.EncodeToString(mac.Sum(nil))
			req.Header.Set("X-Webhook-Signature", sig)
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return SendResult{
				ResponseCode: resp.StatusCode,
				ResponseBody: string(respBody),
				Success:      true,
			}
		}

		// Non-2xx response: retry
		if attempt == maxRetries-1 {
			return SendResult{
				ResponseCode: resp.StatusCode,
				ResponseBody: string(respBody),
				Success:      false,
			}
		}
	}

	return SendResult{Success: false}
}

// BuildPayload constructs a standard webhook payload.
func BuildPayload(trigger string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"event":     trigger,
		"data":      data,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
}

// ValidateURL checks that the URL starts with https://.
func ValidateURL(url string) error {
	if len(url) < 8 || url[:8] != "https://" {
		return fmt.Errorf("webhook URL must start with https://")
	}
	return nil
}

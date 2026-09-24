package cliq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/model"
)

// MaskWebhookURL keeps scheme+host and redacts the path/query for logs.
func MaskWebhookURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "***"
	}
	return u.Scheme + "://" + u.Host + "/***"
}

// Send posts the Cliq payload to the incoming webhook URL.
func Send(client *http.Client, webhookURL string, payload model.CliqPayload) (status int, body string, err error) {
	if client == nil {
		client = http.DefaultClient
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return 0, "", err
	}
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader(encoded))
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("network error while calling Zoho Cliq webhook: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	return resp.StatusCode, string(respBody), nil
}

// IsSuccessStatus reports whether status is 2xx.
func IsSuccessStatus(status int) bool {
	return status >= 200 && status < 300
}

// TruncateBody limits error body logging.
func TruncateBody(body string, max int) string {
	body = strings.TrimSpace(body)
	if len(body) <= max {
		return body
	}
	return body[:max]
}

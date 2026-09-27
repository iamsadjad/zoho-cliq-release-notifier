package httpclient

import (
	"net/http"

	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/constants"
)

// New returns an HTTP client with a bounded request timeout.
// A nil client at a call site must use this instead of http.DefaultClient,
// which has no timeout.
func New() *http.Client {
	return &http.Client{Timeout: constants.HTTPRequestTimeout}
}

// CallerOrNew returns client when the caller supplied one, otherwise New.
func CallerOrNew(client *http.Client) *http.Client {
	if client != nil {
		return client
	}
	return New()
}

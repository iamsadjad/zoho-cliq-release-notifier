package action_test

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/action"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestRun_SkipMissingWebhook(t *testing.T) {
	t.Setenv("INPUT_WEBHOOK_URL", "")
	t.Setenv("INPUT_REPOSITORY", "acme/demo")
	t.Setenv("GITHUB_OUTPUT", filepath.Join(t.TempDir(), "out"))

	runner := &action.Runner{}
	result, err := runner.Run()
	if err != nil {
		t.Fatal(err)
	}
	if result.Notified || string(result.SkippedReason) != "missing_webhook" {
		t.Fatalf("%+v", result)
	}
}

func TestRun_SkipPrerelease(t *testing.T) {
	dir := t.TempDir()
	event := filepath.Join(dir, "event.json")
	_ = os.WriteFile(event, []byte(`{"release":{"tag_name":"v1.0.0-rc.1","prerelease":true,"html_url":"https://example.com","author":{"login":"a"}}}`), 0o644)

	t.Setenv("INPUT_WEBHOOK_URL", "https://cliq.example/hook")
	t.Setenv("INPUT_NOTIFY_PRERELEASE", "false")
	t.Setenv("INPUT_REPOSITORY", "acme/demo")
	t.Setenv("GITHUB_EVENT_PATH", event)
	t.Setenv("GITHUB_OUTPUT", filepath.Join(dir, "out"))

	runner := &action.Runner{}
	result, err := runner.Run()
	if err != nil {
		t.Fatal(err)
	}
	if result.Notified || string(result.SkippedReason) != "prerelease_skipped" {
		t.Fatalf("%+v", result)
	}
}

func TestRun_Notify(t *testing.T) {
	dir := t.TempDir()
	event := filepath.Join(dir, "event.json")
	_ = os.WriteFile(event, []byte(`{"release":{"tag_name":"v1.0.0","name":"One","body":"* shipping","prerelease":false,"html_url":"https://example.com","author":{"login":"bob"}}}`), 0o644)

	t.Setenv("INPUT_WEBHOOK_URL", "https://cliq.example/hook?token=secret")
	t.Setenv("INPUT_NOTIFY_PRERELEASE", "false")
	t.Setenv("INPUT_REPOSITORY", "acme/demo")
	t.Setenv("GITHUB_EVENT_PATH", event)
	t.Setenv("GITHUB_OUTPUT", filepath.Join(dir, "out"))

	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host != "cliq.example" {
			t.Fatalf("unexpected host %s", r.URL.Host)
		}
		return &http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader("ok")),
			Header:     make(http.Header),
		}, nil
	})}

	runner := &action.Runner{HTTPClient: client}
	result, err := runner.Run()
	if err != nil {
		t.Fatal(err)
	}
	if !result.Notified || result.HTTPStatus != 201 || result.ReleaseTag != "v1.0.0" {
		t.Fatalf("%+v", result)
	}
}

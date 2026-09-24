package release_test

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/release"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestResolve_FromEvent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "event.json")
	payload := map[string]any{
		"release": map[string]any{
			"tag_name":   "v9.9.9",
			"name":       "Nine",
			"body":       "* hello",
			"html_url":   "https://example.com/r",
			"author":     map[string]any{"login": "alice"},
			"prerelease": false,
		},
	}
	raw, _ := json.Marshal(payload)
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	r := &release.Resolver{}
	info, err := r.Resolve("", "acme/demo", "", path)
	if err != nil {
		t.Fatal(err)
	}
	if info.TagName != "v9.9.9" || info.AuthorLogin != "alice" {
		t.Fatalf("%+v", info)
	}
}

func TestResolve_MissingEvent(t *testing.T) {
	r := &release.Resolver{}
	_, err := r.Resolve("", "acme/demo", "", "/tmp/does-not-exist-zoho-cliq.json")
	if err == nil || !strings.Contains(err.Error(), "missing release payload") {
		t.Fatalf("err=%v", err)
	}
}

func TestResolve_FetchByTag(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.Path, "/releases/tags/v2.0.0") {
			t.Fatalf("path %s", req.URL.Path)
		}
		if req.Header.Get("Authorization") != "Bearer token" {
			t.Fatal("missing auth")
		}
		body := `{"tag_name":"v2.0.0","name":null,"body":null,"html_url":"https://example.com/v2","author":null,"prerelease":true}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})}

	r := &release.Resolver{HTTPClient: client}
	info, err := r.Resolve("v2.0.0", "acme/demo", "token", "")
	if err != nil {
		t.Fatal(err)
	}
	if !info.Prerelease || info.TagName != "v2.0.0" {
		t.Fatalf("%+v", info)
	}
}

func TestResolve_RequiresTokenForTag(t *testing.T) {
	r := &release.Resolver{}
	_, err := r.Resolve("v1", "acme/demo", "", "")
	if err == nil || !strings.Contains(err.Error(), "github_token is required") {
		t.Fatalf("err=%v", err)
	}
}

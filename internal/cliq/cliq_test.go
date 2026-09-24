package cliq_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/cliq"
	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/model"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestMaskWebhookURL(t *testing.T) {
	got := cliq.MaskWebhookURL("https://cliq.zoho.com/api/v2/channelsbyname/foo/message?zapikey=SECRET")
	if got != "https://cliq.zoho.com/***" {
		t.Fatalf("got %q", got)
	}
	if cliq.MaskWebhookURL("not-a-url") != "***" {
		t.Fatal("expected *** for invalid URL")
	}
}

func TestSend_Success(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Fatalf("method %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatal("missing content-type")
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("ok")),
			Header:     make(http.Header),
		}, nil
	})}

	status, body, err := cliq.Send(client, "https://cliq.example/hook", model.CliqPayload{Text: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if status != 200 || body != "ok" {
		t.Fatalf("status=%d body=%q", status, body)
	}
}

func TestIsSuccessStatus(t *testing.T) {
	if !cliq.IsSuccessStatus(204) || cliq.IsSuccessStatus(500) {
		t.Fatal("unexpected success classification")
	}
}

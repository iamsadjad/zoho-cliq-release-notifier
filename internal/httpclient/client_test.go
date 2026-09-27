package httpclient_test

import (
	"net/http"
	"testing"

	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/constants"
	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/httpclient"
)

func TestNewHasBoundedTimeout(t *testing.T) {
	client := httpclient.New()
	if client.Timeout != constants.HTTPRequestTimeout {
		t.Fatalf("timeout = %s", client.Timeout)
	}
}

func TestCallerOrNewPreservesCallerClient(t *testing.T) {
	caller := &http.Client{}
	if got := httpclient.CallerOrNew(caller); got != caller {
		t.Fatal("expected the caller client to be preserved")
	}
	if got := httpclient.CallerOrNew(nil); got == nil || got.Timeout != constants.HTTPRequestTimeout {
		t.Fatal("expected a bounded client when none is supplied")
	}
}

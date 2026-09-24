package payload_test

import (
	"strings"
	"testing"

	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/constants"
	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/model"
	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/payload"
)

func TestBuild_Production(t *testing.T) {
	p := payload.Build(model.ReleaseInfo{
		TagName:     "v1.2.3",
		Name:        "Release 1.2.3",
		Body:        "* Added docs\n- Fixed bug",
		HTMLURL:     "https://github.com/acme/demo/releases/tag/v1.2.3",
		AuthorLogin: "octocat",
		Prerelease:  false,
	}, "acme/demo")

	if !strings.Contains(p.Text, "✅ _Production release_") {
		t.Fatalf("missing production badge: %s", p.Text)
	}
	if p.Card.Title != "🎉 v1.2.3" || p.Card.Theme != constants.CardTheme {
		t.Fatalf("unexpected card: %+v", p.Card)
	}
	if len(p.Slides) != 2 {
		t.Fatalf("expected 2 slides, got %d", len(p.Slides))
	}
	if len(p.Buttons) != 1 || p.Buttons[0].Action.Data.Web == "" {
		t.Fatalf("unexpected buttons: %+v", p.Buttons)
	}
}

func TestBuild_PrereleaseEmptyBody(t *testing.T) {
	p := payload.Build(model.ReleaseInfo{
		TagName:    "v1.2.3",
		Prerelease: true,
	}, "acme/demo")
	if !strings.Contains(p.Text, "🧪 _Pre-release_") {
		t.Fatalf("missing prerelease badge")
	}
	if !strings.Contains(p.Text, "✨ *v1.2.3*") {
		t.Fatalf("expected tag fallback for empty name")
	}
	if len(p.Slides) != 1 {
		t.Fatalf("expected only details slide, got %d", len(p.Slides))
	}
}

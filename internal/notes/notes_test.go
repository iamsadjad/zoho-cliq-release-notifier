package notes_test

import (
	"strings"
	"testing"

	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/constants"
	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/notes"
)

func TestExtractWhatsNewNotes_Bullets(t *testing.T) {
	body := strings.Join([]string{"## Changes", "* First", "- Second", "plain", "* Third"}, "\n")
	got := notes.ExtractWhatsNewNotes(body)
	want := []string{"First", "Second", "Third"}
	if len(got) != len(want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v want %v", got, want)
		}
	}
}

func TestExtractWhatsNewNotes_Limit(t *testing.T) {
	var lines []string
	for i := 1; i <= 15; i++ {
		lines = append(lines, "* Item "+itoa(i))
	}
	got := notes.ExtractWhatsNewNotes(strings.Join(lines, "\n"))
	if len(got) != constants.MaxWhatsNewNotes {
		t.Fatalf("len=%d want %d", len(got), constants.MaxWhatsNewNotes)
	}
}

func TestExtractWhatsNewNotes_Fallback(t *testing.T) {
	got := notes.ExtractWhatsNewNotes("Line one.\n\nLine   two.")
	if len(got) != 1 || got[0] != "Line one. Line two." {
		t.Fatalf("unexpected fallback: %#v", got)
	}
}

func TestExtractWhatsNewNotes_Empty(t *testing.T) {
	if notes.ExtractWhatsNewNotes("") != nil {
		t.Fatal("expected nil")
	}
}

func TestExtractWhatsNewNotes_IgnoresStarWithoutSpace(t *testing.T) {
	got := notes.ExtractWhatsNewNotes("*nospace\n* with space")
	if len(got) != 1 || got[0] != "with space" {
		t.Fatalf("got %#v", got)
	}
}

func itoa(n int) string {
	const digits = "0123456789"
	if n < 10 {
		return digits[n : n+1]
	}
	return itoa(n/10) + digits[n%10:n%10+1]
}

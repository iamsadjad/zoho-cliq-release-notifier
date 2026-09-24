package notes

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/constants"
)

var bulletLine = regexp.MustCompile(`^(?:\* |- )(.*)$`)

// ExtractWhatsNewNotes extracts up to MaxWhatsNewNotes changelog bullets.
// Falls back to a whitespace-collapsed summary when no bullets exist.
func ExtractWhatsNewNotes(releaseBody string) []string {
	var bullets []string
	for _, line := range strings.Split(releaseBody, "\n") {
		line = strings.TrimRight(line, "\r")
		m := bulletLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		note := strings.TrimSpace(m[1])
		if note == "" {
			continue
		}
		bullets = append(bullets, note)
		if len(bullets) >= constants.MaxWhatsNewNotes {
			break
		}
	}
	if len(bullets) > 0 {
		return bullets
	}

	formatted := releaseBody
	if len(formatted) > constants.BodyTruncateLength {
		formatted = formatted[:constants.BodyTruncateLength]
	}
	fallback := collapseWhitespace(formatted)
	if len(fallback) > constants.FallbackNoteMaxLength {
		fallback = fallback[:constants.FallbackNoteMaxLength]
	}
	if fallback == "" {
		return nil
	}
	return []string{fallback}
}

func collapseWhitespace(value string) string {
	var b strings.Builder
	b.Grow(len(value))
	prevSpace := false
	for _, r := range value {
		if unicode.IsSpace(r) {
			if prevSpace {
				continue
			}
			b.WriteByte(' ')
			prevSpace = true
			continue
		}
		prevSpace = false
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

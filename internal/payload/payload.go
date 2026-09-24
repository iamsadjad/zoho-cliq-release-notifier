package payload

import (
	"fmt"
	"strings"

	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/constants"
	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/model"
	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/notes"
)

// Build builds the Zoho Cliq incoming-webhook payload.
func Build(release model.ReleaseInfo, repository string) model.CliqPayload {
	displayName := strings.TrimSpace(release.Name)
	if displayName == "" {
		displayName = release.TagName
	}

	badge := "✅ _Production release_"
	if release.Prerelease {
		badge = "🧪 _Pre-release_"
	}

	messageText := strings.Join([]string{
		"🚀 *New Release Published!*",
		"",
		fmt.Sprintf("✨ *%s*", displayName),
		fmt.Sprintf("🏷️ `%s`  ·  👤 _%s_", release.TagName, release.AuthorLogin),
		fmt.Sprintf("📦 `%s`", repository),
		"",
		badge,
	}, "\n")

	if len(messageText) > constants.MessageTextMaxLength {
		messageText = messageText[:constants.MessageTextMaxLength] + "..."
	}

	slides := []model.CliqSlide{
		{
			Type:  "label",
			Title: constants.ReleaseDetailsTitle,
			Data: []map[string]string{
				{"Release": displayName},
				{"Tag": release.TagName},
				{"Author": release.AuthorLogin},
				{"Repository": repository},
			},
		},
	}

	whatsNew := notes.ExtractWhatsNewNotes(release.Body)
	if len(whatsNew) > 0 {
		slides = append(slides, model.CliqSlide{
			Type:  "list",
			Title: constants.WhatsNewTitle,
			Data:  whatsNew,
		})
	}

	return model.CliqPayload{
		Text: messageText,
		Card: model.CliqCard{
			Title:     "🎉 " + release.TagName,
			Theme:     constants.CardTheme,
			Thumbnail: constants.DefaultThumbnailURL,
		},
		Slides: slides,
		Buttons: []model.CliqButton{
			{
				Label: constants.ViewOnGitHubLabel,
				Type:  "+",
				Action: model.CliqButtonAction{
					Type: "open.url",
					Data: model.CliqButtonActionData{Web: release.HTMLURL},
				},
			},
		},
	}
}

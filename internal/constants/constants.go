package constants

import "time"

const (
	// HTTPRequestTimeout bounds each outbound webhook and GitHub API call.
	HTTPRequestTimeout    = 30 * time.Second
	BodyTruncateLength    = 1800
	MessageTextMaxLength  = 3500
	MaxWhatsNewNotes      = 10
	FallbackNoteMaxLength = 500
	DefaultThumbnailURL   = "https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png"

	WhatsNewTitle       = "📝 What's New"
	ReleaseDetailsTitle = "📋 Release Details"
	ViewOnGitHubLabel   = "📦 View on GitHub"
	CardTheme           = "modern-inline"
)

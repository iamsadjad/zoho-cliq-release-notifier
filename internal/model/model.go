package model

// ReleaseInfo is the subset of GitHub release metadata used for notifications.
type ReleaseInfo struct {
	TagName     string
	Name        string
	Body        string
	HTMLURL     string
	AuthorLogin string
	Prerelease  bool
}

// Inputs are the GitHub Action inputs.
type Inputs struct {
	WebhookURL       string
	Tag              string
	NotifyPrerelease bool
	GitHubToken      string
	Repository       string
	EventPath        string
}

// SkipReason explains why a notification was not sent.
type SkipReason string

const (
	SkipNone              SkipReason = ""
	SkipMissingWebhook    SkipReason = "missing_webhook"
	SkipPrereleaseSkipped SkipReason = "prerelease_skipped"
)

// Result is written to GitHub Action outputs.
type Result struct {
	Notified      bool
	SkippedReason SkipReason
	HTTPStatus    int // 0 when not called
	ReleaseTag    string
	Prerelease    *bool
}

// CliqPayload is the Zoho Cliq incoming-webhook JSON body.
type CliqPayload struct {
	Text    string       `json:"text"`
	Card    CliqCard     `json:"card"`
	Slides  []CliqSlide  `json:"slides"`
	Buttons []CliqButton `json:"buttons"`
}

type CliqCard struct {
	Title     string `json:"title"`
	Theme     string `json:"theme"`
	Thumbnail string `json:"thumbnail"`
}

// CliqSlide is either a label or list slide.
type CliqSlide struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Data  any    `json:"data"`
}

type CliqButton struct {
	Label  string           `json:"label"`
	Type   string           `json:"type"`
	Action CliqButtonAction `json:"action"`
}

type CliqButtonAction struct {
	Type string               `json:"type"`
	Data CliqButtonActionData `json:"data"`
}

type CliqButtonActionData struct {
	Web string `json:"web"`
}

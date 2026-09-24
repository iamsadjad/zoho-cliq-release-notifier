package action

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/cliq"
	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/gha"
	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/inputs"
	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/model"
	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/payload"
	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/release"
)

// Runner executes the notification action.
type Runner struct {
	HTTPClient *http.Client
}

// Run loads inputs, resolves the release, and notifies Zoho Cliq.
func (r *Runner) Run() (model.Result, error) {
	in, err := inputs.Load()
	if err != nil {
		return model.Result{}, err
	}

	if in.WebhookURL != "" {
		gha.SetSecret(in.WebhookURL)
	}

	if in.WebhookURL == "" {
		gha.Warning("webhook_url is not set. Skipping notification.")
		result := model.Result{
			Notified:      false,
			SkippedReason: model.SkipMissingWebhook,
		}
		return result, writeOutputs(result)
	}

	resolver := &release.Resolver{HTTPClient: r.HTTPClient}
	rel, err := resolver.Resolve(in.Tag, in.Repository, in.GitHubToken, in.EventPath)
	if err != nil {
		return model.Result{}, err
	}

	if rel.Prerelease && !in.NotifyPrerelease {
		gha.Info("Skipping prerelease notification for tag %s", rel.TagName)
		prerelease := true
		result := model.Result{
			Notified:      false,
			SkippedReason: model.SkipPrereleaseSkipped,
			ReleaseTag:    rel.TagName,
			Prerelease:    &prerelease,
		}
		return result, writeOutputs(result)
	}

	body := payload.Build(rel, in.Repository)
	gha.Info("Sending notification to Zoho Cliq (%s)...", cliq.MaskWebhookURL(in.WebhookURL))

	status, respBody, err := cliq.Send(r.HTTPClient, in.WebhookURL, body)
	if err != nil {
		return model.Result{}, err
	}
	if !cliq.IsSuccessStatus(status) {
		gha.Error("Failed to send notification to Zoho Cliq (HTTP %d)", status)
		if respBody != "" {
			gha.Error("%s", cliq.TruncateBody(respBody, 2000))
		}
		return model.Result{}, fmt.Errorf("Zoho Cliq webhook returned HTTP %d", status)
	}

	gha.Info("Successfully sent notification to Zoho Cliq (HTTP %d)", status)
	prerelease := rel.Prerelease
	result := model.Result{
		Notified:   true,
		HTTPStatus: status,
		ReleaseTag: rel.TagName,
		Prerelease: &prerelease,
	}
	return result, writeOutputs(result)
}

func writeOutputs(result model.Result) error {
	if err := gha.SetOutput("notified", strconv.FormatBool(result.Notified)); err != nil {
		return err
	}
	if err := gha.SetOutput("skipped_reason", string(result.SkippedReason)); err != nil {
		return err
	}
	httpStatus := ""
	if result.HTTPStatus != 0 {
		httpStatus = strconv.Itoa(result.HTTPStatus)
	}
	if err := gha.SetOutput("http_status", httpStatus); err != nil {
		return err
	}
	if err := gha.SetOutput("release_tag", result.ReleaseTag); err != nil {
		return err
	}
	prerelease := ""
	if result.Prerelease != nil {
		prerelease = strconv.FormatBool(*result.Prerelease)
	}
	return gha.SetOutput("prerelease", prerelease)
}
